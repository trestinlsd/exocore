package keeper

import (
	"testing"
	"time"

	"github.com/ExocoreNetwork/exocore/x/oracle/keeper"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"

	assetskeeper "github.com/ExocoreNetwork/exocore/x/assets/keeper"
	delegationkeeper "github.com/ExocoreNetwork/exocore/x/delegation/keeper"
	dogfoodkeeper "github.com/ExocoreNetwork/exocore/x/dogfood/keeper"
	"github.com/stretchr/testify/require"
)

func OracleKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"OracleParams",
	)
	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		dogfoodkeeper.Keeper{},
		delegationkeeper.Keeper{},
		assetskeeper.Keeper{},
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{
		Time: time.Now().UTC(),
	}, false, log.NewNopLogger())

	// Initialize params
	p4Test := types.DefaultParams()
	p4Test.TokenFeeders[1].StartBaseBlock = 1
	k.SetParams(ctx, p4Test)

	return &k, ctx
}
