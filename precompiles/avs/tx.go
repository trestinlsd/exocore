//nolint:dupl
package avs

import (
	"fmt"
	"slices"

	errorsmod "cosmossdk.io/errors"

	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	cmn "github.com/evmos/evmos/v16/precompiles/common"
)

const (
	MethodRegisterAVS               = "registerAVS"
	MethodUpdateAVS                 = "updateAVS"
	MethodDeregisterAVS             = "deregisterAVS"
	MethodRegisterOperatorToAVS     = "registerOperatorToAVS"
	MethodDeregisterOperatorFromAVS = "deregisterOperatorFromAVS"
	MethodCreateAVSTask             = "createTask"
	MethodRegisterBLSPublicKey      = "registerBLSPublicKey"
	MethodChallenge                 = "challenge"
	MethodOperatorSubmitTask        = "operatorSubmitTask"
)

// RegisterAVS AVSInfoRegister register the avs related information and change the state in avs keeper module.
func (p Precompile) RegisterAVS(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// parse the avs input params first.
	avsParams, err := p.GetAVSParamsFromInputs(ctx, args)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse args error")
	}
	// verification of the calling address to ensure it is avs contract owner
	if !slices.Contains(avsParams.AvsOwnerAddress, avsParams.CallerAddress) {
		return nil, errorsmod.Wrap(err, "not qualified to registerOrDeregister")
	}
	// The AVS registration is done by the calling contract.
	avsParams.AvsAddress = contract.CallerAddress.String()
	avsParams.Action = avstypes.RegisterAction
	// Finally, update the AVS information in the keeper.
	err = p.avsKeeper.UpdateAVSInfo(ctx, avsParams)
	if err != nil {
		fmt.Println("Failed to update AVS info", err)
		return nil, err
	}
	if err = p.EmitAVSRegistered(ctx, stateDB, avsParams); err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

func (p Precompile) DeregisterAVS(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != len(p.ABI.Methods[MethodDeregisterAVS].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, len(p.ABI.Methods[MethodDeregisterAVS].Inputs), len(args))
	}
	avsParams := &avstypes.AVSRegisterOrDeregisterParams{}
	callerAddress, ok := args[0].(common.Address)
	if !ok || (callerAddress == common.Address{}) {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "common.Address", callerAddress)
	}
	avsParams.CallerAddress = sdk.AccAddress(callerAddress[:]).String()
	avsName, ok := args[1].(string)
	if !ok || avsName == "" {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, "string", avsName)
	}
	avsParams.AvsName = avsName

	avsParams.AvsAddress = contract.CallerAddress.String()
	avsParams.Action = avstypes.DeRegisterAction
	// validates that this is owner

	err := p.avsKeeper.UpdateAVSInfo(ctx, avsParams)
	if err != nil {
		return nil, err
	}
	if err = p.EmitAVSDeregistered(ctx, stateDB, avsParams); err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

func (p Precompile) UpdateAVS(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// parse the avs input params first.
	avsParams, err := p.GetAVSParamsFromUpdateInputs(ctx, args)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse args error")
	}

	avsParams.AvsAddress = contract.CallerAddress.String()
	avsParams.Action = avstypes.UpdateAction
	previousAVSInfo, err := p.avsKeeper.GetAVSInfo(ctx, avsParams.AvsAddress)
	if err != nil {
		return nil, err
	}
	// If avs UpdateAction check CallerAddress
	if !slices.Contains(previousAVSInfo.Info.AvsOwnerAddress, avsParams.CallerAddress) {
		return nil, fmt.Errorf("this caller not qualified to update %s", avsParams.CallerAddress)
	}
	err = p.avsKeeper.UpdateAVSInfo(ctx, avsParams)
	if err != nil {
		return nil, err
	}

	if err = p.EmitAVSUpdated(ctx, stateDB, avsParams); err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

func (p Precompile) BindOperatorToAVS(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != len(p.ABI.Methods[MethodRegisterOperatorToAVS].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, len(p.ABI.Methods[MethodRegisterOperatorToAVS].Inputs), len(args))
	}
	callerAddress, ok := args[0].(common.Address)
	if !ok || (callerAddress == common.Address{}) {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "common.Address", callerAddress)
	}

	operatorParams := &avstypes.OperatorOptParams{}
	operatorParams.OperatorAddress = sdk.AccAddress(callerAddress[:]).String()
	operatorParams.AvsAddress = contract.CallerAddress.String()
	operatorParams.Action = avstypes.RegisterAction
	err := p.avsKeeper.OperatorOptAction(ctx, operatorParams)
	if err != nil {
		return nil, err
	}
	if err = p.EmitOperatorJoined(ctx, stateDB, operatorParams); err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

func (p Precompile) UnbindOperatorToAVS(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != len(p.ABI.Methods[MethodRegisterOperatorToAVS].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, len(p.ABI.Methods[MethodRegisterOperatorToAVS].Inputs), len(args))
	}
	callerAddress, ok := args[0].(common.Address)
	if !ok || (callerAddress == common.Address{}) {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "common.Address", callerAddress)
	}
	operatorParams := &avstypes.OperatorOptParams{}
	operatorParams.OperatorAddress = sdk.AccAddress(callerAddress[:]).String()
	operatorParams.AvsAddress = contract.CallerAddress.String()
	operatorParams.Action = avstypes.DeRegisterAction
	err := p.avsKeeper.OperatorOptAction(ctx, operatorParams)
	if err != nil {
		return nil, err
	}
	if err = p.EmitOperatorOuted(ctx, stateDB, operatorParams); err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

// CreateAVSTask Middleware uses exocore's default avstask template to create tasks in avstask module.
func (p Precompile) CreateAVSTask(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	params, err := p.GetTaskParamsFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}
	params.TaskContractAddress = contract.CallerAddress.String()
	taskID, err := p.avsKeeper.CreateAVSTask(ctx, params)
	if err != nil {
		return nil, err
	}
	if err = p.EmitTaskCreated(ctx, stateDB, params); err != nil {
		return nil, err
	}
	return method.Outputs.Pack(taskID)
}

// Challenge Middleware uses exocore's default avstask template to create tasks in avstask module.
func (p Precompile) Challenge(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != len(p.ABI.Methods[MethodChallenge].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, len(p.ABI.Methods[MethodChallenge].Inputs), len(args))
	}
	challengeParams := &avstypes.ChallengeParams{}
	challengeParams.TaskContractAddress = contract.CallerAddress
	callerAddress, ok := args[0].(common.Address)
	if !ok || (callerAddress == common.Address{}) {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "common.Address", callerAddress)
	}
	challengeParams.CallerAddress = sdk.AccAddress(callerAddress[:]).String()

	taskHash, ok := args[1].([]byte)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, "[]byte", taskHash)
	}
	challengeParams.TaskHash = taskHash

	taskID, ok := args[2].(uint64)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 2, "uint64", taskID)
	}
	challengeParams.TaskID = taskID

	taskResponseHash, ok := args[3].([]byte)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 3, "[]byte", taskResponseHash)
	}
	challengeParams.TaskResponseHash = taskResponseHash

	operatorAddress, ok := args[4].(string)
	if !ok || operatorAddress == "" {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 4, "string", operatorAddress)
	}
	operator, err := sdk.AccAddressFromBech32(operatorAddress)
	if err != nil {
		return nil, err
	}

	challengeParams.OperatorAddress = operator
	err = p.avsKeeper.RaiseAndResolveChallenge(ctx, challengeParams)
	if err != nil {
		return nil, err
	}

	if err = p.EmitChallengeInitiated(ctx, stateDB, challengeParams); err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

// RegisterBLSPublicKey
func (p Precompile) RegisterBLSPublicKey(
	ctx sdk.Context,
	_ common.Address,
	_ *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != len(p.ABI.Methods[MethodRegisterBLSPublicKey].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, len(p.ABI.Methods[MethodRegisterBLSPublicKey].Inputs), len(args))
	}
	blsParams := &avstypes.BlsParams{}
	callerAddress, ok := args[0].(common.Address)
	if !ok || (callerAddress == common.Address{}) {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "common.Address", callerAddress)
	}
	blsParams.Operator = sdk.AccAddress(callerAddress[:]).String()
	name, ok := args[1].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, "string", name)
	}
	blsParams.Name = name

	pubkeyBz, ok := args[2].([]byte)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 2, "[]byte", pubkeyBz)
	}
	blsParams.PubKey = pubkeyBz

	pubkeyRegistrationSignature, ok := args[3].([]byte)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 3, "[]byte", pubkeyRegistrationSignature)
	}
	blsParams.PubkeyRegistrationSignature = pubkeyRegistrationSignature

	pubkeyRegistrationMessageHash, ok := args[4].([]byte)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 4, "[]byte", pubkeyRegistrationMessageHash)
	}
	blsParams.PubkeyRegistrationMessageHash = pubkeyRegistrationMessageHash

	err := p.avsKeeper.RegisterBLSPublicKey(ctx, blsParams)
	if err != nil {
		return nil, err
	}

	if err = p.EmitPublicKeyRegistered(ctx, stateDB, blsParams); err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

// OperatorSubmitTask operator submit results
func (p Precompile) OperatorSubmitTask(
	ctx sdk.Context,
	_ common.Address,
	_ *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != len(p.ABI.Methods[MethodOperatorSubmitTask].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, len(p.ABI.Methods[MethodOperatorSubmitTask].Inputs), len(args))
	}
	resultParams := &avstypes.TaskResultParams{}

	callerAddress, ok := args[0].(common.Address)
	if !ok || (callerAddress == common.Address{}) {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "common.Address", callerAddress)
	}
	resultParams.CallerAddress = sdk.AccAddress(callerAddress[:]).String()

	taskID, ok := args[1].(uint64)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, "uint64", args[1])
	}
	resultParams.TaskID = taskID

	taskResponse, ok := args[2].([]byte)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 2, "[]byte", taskResponse)
	}
	resultParams.TaskResponse = taskResponse

	blsSignature, ok := args[3].([]byte)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 3, "[]byte", blsSignature)
	}
	resultParams.BlsSignature = blsSignature

	taskAddr, ok := args[4].(common.Address)
	if !ok || (taskAddr == common.Address{}) {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 4, "common.Address", taskAddr)
	}
	resultParams.TaskContractAddress = taskAddr

	phase, ok := args[5].(uint8)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 5, "uint8", phase)
	}
	resultParams.Phase = phase

	resultParams.OperatorAddress = resultParams.CallerAddress

	result := &avstypes.TaskResultInfo{
		TaskId:              resultParams.TaskID,
		OperatorAddress:     resultParams.OperatorAddress,
		TaskContractAddress: resultParams.TaskContractAddress.String(),
		TaskResponse:        resultParams.TaskResponse,
		BlsSignature:        resultParams.BlsSignature,
		Phase:               uint32(resultParams.Phase),
	}
	err := p.avsKeeper.SetTaskResultInfo(ctx, resultParams.OperatorAddress, result)
	if err != nil {
		return nil, err
	}

	if err := p.EmitTaskSubmittedByOperator(ctx, stateDB, resultParams); err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}
