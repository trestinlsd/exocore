package task_test

import (
	"github.com/ExocoreNetwork/exocore/app"
	"github.com/ExocoreNetwork/exocore/precompiles/avsTask"
	"github.com/ExocoreNetwork/exocore/x/taskmanageravs/keeper"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/evmos/v14/x/evm/statedb"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
	"math/big"
)

func (s *TaskPrecompileTestSuite) TestIsTransaction() {
	testCases := []struct {
		name   string
		method string
		isTx   bool
	}{
		{
			task.MethodCreateNewTask,
			s.precompile.Methods[task.MethodCreateNewTask].Name,
			true,
		},
		{
			"invalid",
			"invalid",
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Require().Equal(s.precompile.IsTransaction(tc.method), tc.isTx)
		})
	}
}

func paddingClientChainAddress(input []byte, outputLength int) []byte {
	if len(input) < outputLength {
		padding := make([]byte, outputLength-len(input))
		return append(input, padding...)
	}
	return input
}

// TestRun tests the precompiles Run method createTask.
func (s *TaskPrecompileTestSuite) TestRunCreateTask() {
	// deposit params for test
	//exoCoreLzAppEventTopic := "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec"
	//	usdtAddress := common.FromHex("0xdAC17F958D2ee523a2206206994597C13D831ec7")

	successRet, err := s.precompile.Methods[task.MethodCreateNewTask].Outputs.Pack(true)
	s.Require().NoError(err)
	testcases := []struct {
		name        string
		malleate    func() (common.Address, []byte)
		readOnly    bool
		expPass     bool
		errContains string
		returnBytes []byte
	}{
		{
			name: "pass - task via pre-compiles",
			malleate: func() (common.Address, []byte) {
				taskParam := &keeper.CreateNewTaskParams{

					NumberToBeSquared:         3,
					QuorumThresholdPercentage: 3,
					QuorumNumbers:             nil,
				}
				_, err = s.App.AvsTaskKeeper.CreateNewTask(s.Ctx, taskParam)
				s.Require().NoError(err)
				return [20]byte{}, nil
			},
			returnBytes: successRet,
			readOnly:    false,
			expPass:     true,
		},
	}
	for _, tc := range testcases {
		tc := tc
		s.Run(tc.name, func() {
			// setup basic test suite
			s.SetupTest()

			baseFee := s.App.FeeMarketKeeper.GetBaseFee(s.Ctx)

			// malleate testcase
			caller, input := tc.malleate()

			contract := vm.NewPrecompile(vm.AccountRef(caller), s.precompile, big.NewInt(0), uint64(1e6))
			contract.Input = input

			contractAddr := contract.Address()
			// Build and sign Ethereum transaction
			txArgs := evmtypes.EvmTxArgs{
				ChainID:   s.App.EvmKeeper.ChainID(),
				Nonce:     0,
				To:        &contractAddr,
				Amount:    nil,
				GasLimit:  100000,
				GasPrice:  app.MainnetMinGasPrices.BigInt(),
				GasFeeCap: baseFee,
				GasTipCap: big.NewInt(1),
				Accesses:  &ethtypes.AccessList{},
			}
			msgEthereumTx := evmtypes.NewTx(&txArgs)

			msgEthereumTx.From = s.Address.String()
			err := msgEthereumTx.Sign(s.EthSigner, s.Signer)
			s.Require().NoError(err, "failed to sign Ethereum message")

			// Instantiate config
			proposerAddress := s.Ctx.BlockHeader().ProposerAddress
			cfg, err := s.App.EvmKeeper.EVMConfig(s.Ctx, proposerAddress, s.App.EvmKeeper.ChainID())
			s.Require().NoError(err, "failed to instantiate EVM config")

			msg, err := msgEthereumTx.AsMessage(s.EthSigner, baseFee)
			s.Require().NoError(err, "failed to instantiate Ethereum message")

			// Create StateDB
			s.StateDB = statedb.New(s.Ctx, s.App.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(s.Ctx.HeaderHash().Bytes())))
			// Instantiate EVM
			evm := s.App.EvmKeeper.NewEVM(
				s.Ctx, msg, cfg, nil, s.StateDB,
			)
			params := s.App.EvmKeeper.GetParams(s.Ctx)
			activePrecompiles := params.GetActivePrecompilesAddrs()
			precompileMap := s.App.EvmKeeper.Precompiles(activePrecompiles...)
			err = vm.ValidatePrecompiles(precompileMap, activePrecompiles)
			s.Require().NoError(err, "invalid precompiles", activePrecompiles)
			evm.WithPrecompiles(precompileMap, activePrecompiles)

			// Run precompiled contract
			bz, err := s.precompile.Run(evm, contract, tc.readOnly)

			// Check results
			if tc.expPass {
				s.Require().NoError(err, "expected no error when running the precompile")
				s.Require().Equal(tc.returnBytes, bz, "the return doesn't match the expected result")
			} else {
				s.Require().Error(err, "expected error to be returned when running the precompile")
				s.Require().Nil(bz, "expected returned bytes to be nil")
				s.Require().ErrorContains(err, tc.errContains)
			}
		})
	}
}