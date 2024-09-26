package keeper

import (
	"strings"

	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EpochsHooksWrapper is the wrapper structure that implements the epochs hooks for the avs
// keeper.
type EpochsHooksWrapper struct {
	keeper *Keeper
}

// Interface guard
var _ epochstypes.EpochHooks = EpochsHooksWrapper{}

// EpochsHooks returns the epochs hooks wrapper.
func (k *Keeper) EpochsHooks() EpochsHooksWrapper {
	return EpochsHooksWrapper{k}
}

// BeforeEpochStart: noop, We don't need to do anything here
func (wrapper EpochsHooksWrapper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) {
}

// AfterEpochEnd mints and allocates coins at the end of each epoch end
func (wrapper EpochsHooksWrapper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, _ int64) {
	expEpochID := wrapper.keeper.GetParams(ctx).EpochIdentifier
	if strings.Compare(epochIdentifier, expEpochID) == 0 {
		// the minted coins generated by minting module will do the token allocation and distribution here
		previousTotalPower := wrapper.keeper.StakingKeeper.GetLastTotalPower(ctx)
		logger := wrapper.keeper.Logger()
		logger.Info(
			"AfterEpochEnd of distribution",
		)
		err := wrapper.keeper.AllocateTokens(ctx, previousTotalPower.Int64())
		if err != nil {
			logger.Error("failed to allocate tokens", "err", err)
			return
		}
	}
}