package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/pflag"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

const (
	FlagOperatorAddress     = "operator-address"
	FlagTaskResponse        = "task-response"
	FlagBlsSignature        = "bls-signature"
	FlagTaskContractAddress = "task-contract-address"
	FlagTaskID              = "task-id"
	FlagPhase               = "phase"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		CmdSubmitTaskResult(),
	)
	return txCmd
}

// CmdSubmitTaskResult returns a CLI command handler for submit  a TaskResult
// transaction.
func CmdSubmitTaskResult() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-task-result",
		Short: "submit task result",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			msg := newBuildMsg(clientCtx, cmd.Flags())

			// this calls ValidateBasic internally so we don't need to do that.
			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	f := cmd.Flags()
	f.String(
		FlagOperatorAddress, "", "The address of the operator being queried "+
			" If not provided, it will default to the sender's address.",
	)
	f.String(
		FlagTaskResponse, "", "The task response data",
	)
	f.String(
		FlagBlsSignature, "", "The operator bls sig info",
	)
	f.String(
		FlagTaskContractAddress, "", "The contract address of task",
	)
	f.Uint64(
		FlagTaskID, 1, "The  task id",
	)
	f.Uint32(
		FlagPhase, 0, "The phase is a two-phase submission with two values, 0 and 1",
	)
	// #nosec G703 // this only errors if the flag isn't defined.
	_ = cmd.MarkFlagRequired(FlagTaskID)
	_ = cmd.MarkFlagRequired(FlagBlsSignature)
	_ = cmd.MarkFlagRequired(FlagTaskContractAddress)
	_ = cmd.MarkFlagRequired(FlagPhase)

	// transaction level flags from the SDK
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func newBuildMsg(
	clientCtx client.Context, fs *pflag.FlagSet,
) *types.SubmitTaskResultReq {
	sender := clientCtx.GetFromAddress()
	operatorAddress, _ := fs.GetString(FlagOperatorAddress)
	if operatorAddress == "" {
		operatorAddress = sender.String()
	}
	taskResponse, _ := fs.GetString(FlagTaskResponse)
	taskRes, _ := hex.DecodeString(taskResponse)
	blsSignature, _ := fs.GetString(FlagBlsSignature)
	sig, _ := hex.DecodeString(blsSignature)
	taskContractAddress, _ := fs.GetString(FlagTaskContractAddress)

	taskID, _ := fs.GetUint64(FlagTaskID)
	phase, _ := fs.GetUint32(FlagPhase)

	msg := &types.SubmitTaskResultReq{
		FromAddress: sender.String(),
		Info: &types.TaskResultInfo{
			OperatorAddress:     operatorAddress,
			TaskResponse:        taskRes,
			BlsSignature:        sig,
			TaskContractAddress: taskContractAddress,
			TaskId:              taskID,
			Phase:               phase,
		},
	}
	return msg
}
