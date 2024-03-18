package types

import (
	sdkmath "cosmossdk.io/math"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ OracleKeeper = MockOracle{}
	_ AvsKeeper    = MockAvs{}
)

type AssetsKeeper interface {
	GetStakingAssetInfo(ctx sdk.Context, assetID string) (info *assetstype.StakingAssetInfo, err error)
	IteratorOperatorAssetState(ctx sdk.Context, f func(operatorAddr, assetID string, state *assetstype.OperatorAssetInfo) error) error
	AppChainInfoIsExist(ctx sdk.Context, chainID string) bool
	GetOperatorAssetInfos(ctx sdk.Context, operatorAddr sdk.Address, _ map[string]interface{}) (assetsInfo map[string]*assetstype.OperatorAssetInfo, err error)
	UpdateStakerAssetState(ctx sdk.Context, stakerID string, assetID string, changeAmount assetstype.StakerSingleAssetChangeInfo) (err error)
	UpdateOperatorAssetState(ctx sdk.Context, operatorAddr sdk.Address, assetID string, changeAmount assetstype.OperatorSingleAssetChangeInfo) (err error)
}

type DelegationKeeper interface {
	DelegationStateByOperatorAssets(ctx sdk.Context, operatorAddr string, assetsFilter map[string]interface{}) (map[string]map[string]delegationtype.DelegationAmounts, error)
	IterateDelegationState(ctx sdk.Context, f func(restakerID, assetID, operatorAddr string, state *delegationtype.DelegationAmounts) error) error
	UpdateDelegationState(ctx sdk.Context, stakerID string, assetID string, delegationAmounts map[string]*delegationtype.DelegationAmounts) (err error)

	UpdateStakerDelegationTotalAmount(ctx sdk.Context, stakerID string, assetID string, opAmount sdkmath.Int) error
}

type PriceChange struct {
	OriginalPrice sdkmath.Int
	NewPrice      sdkmath.Int
	Decimal       uint8
}

// OracleKeeper is the oracle interface expected by operator module
// These functions need to be implemented by the oracle module
type OracleKeeper interface {
	// GetSpecifiedAssetsPrice is a function to retrieve the asset price according to the assetID
	// the first return value is price, and the second return value is decimal of the price.
	GetSpecifiedAssetsPrice(ctx sdk.Context, assetID string) (sdkmath.Int, uint8, error)
	// GetPriceChangeAssets the operator module expect a function that can retrieve all information
	// about assets price change. Then it can update the USD share state according to the change
	// information. This function need to return a map, the key is assetID and the value is PriceChange
	GetPriceChangeAssets(ctx sdk.Context) (map[string]*PriceChange, error)
}

type MockOracle struct{}

func (MockOracle) GetSpecifiedAssetsPrice(_ sdk.Context, _ string) (sdkmath.Int, uint8, error) {
	return sdkmath.NewInt(1), 0, nil
}

func (MockOracle) GetPriceChangeAssets(_ sdk.Context) (map[string]*PriceChange, error) {
	// use USDT as the mock asset
	ret := make(map[string]*PriceChange, 0)
	usdtAssetID := "0xdac17f958d2ee523a2206206994597c13d831ec7_0x65"
	ret[usdtAssetID] = &PriceChange{
		NewPrice:      sdkmath.NewInt(1),
		OriginalPrice: sdkmath.NewInt(1),
		Decimal:       0,
	}
	return nil, nil
}

type MockAvs struct{}

func (MockAvs) GetAvsSupportedAssets(_ sdk.Context, _ string) (map[string]interface{}, error) {
	// set USDT as the default asset supported by AVS
	ret := make(map[string]interface{})
	usdtAssetID := "0xdac17f958d2ee523a2206206994597c13d831ec7_0x65"
	ret[usdtAssetID] = nil
	return ret, nil
}

func (MockAvs) GetAvsSlashContract(_ sdk.Context, _ string) (string, error) {
	return "", nil
}

type AvsKeeper interface {
	// GetAvsSupportedAssets The ctx can be historical or current, depending on the state you wish to retrieve.
	// If the caller want to retrieve a historical assets info supported by Avs, it needs to generate a historical
	// context through calling `ContextForHistoricalState` implemented in x/assets/types/general.go
	GetAvsSupportedAssets(ctx sdk.Context, avsAddr string) (map[string]interface{}, error)
	GetAvsSlashContract(ctx sdk.Context, avsAddr string) (string, error)
}

// add for dogfood

type SlashKeeper interface {
	IsOperatorFrozen(ctx sdk.Context, addr sdk.AccAddress) bool
}

type OperatorConsentHooks interface {
	// This hook is called when an operator opts in to a chain.
	AfterOperatorOptIn(
		ctx sdk.Context,
		addr sdk.AccAddress,
		chainID string,
		pubKey tmprotocrypto.PublicKey,
	)
	// This hook is called when an operator's consensus key is replaced for
	// a chain.
	AfterOperatorKeyReplacement(
		ctx sdk.Context,
		addr sdk.AccAddress,
		oldKey tmprotocrypto.PublicKey,
		newKey tmprotocrypto.PublicKey,
		chainID string,
	)
	// This hook is called when an operator opts out of a chain.
	AfterOperatorOptOutInitiated(
		ctx sdk.Context,
		addr sdk.AccAddress,
		chainID string,
		key tmprotocrypto.PublicKey,
	)
}