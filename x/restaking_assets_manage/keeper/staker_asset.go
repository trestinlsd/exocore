package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
)

func (k Keeper) GetStakerAssetInfos(ctx sdk.Context, stakerId string) (assetsInfo map[string]*types2.StakerSingleAssetOrChangeInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixReStakerAssetInfos)
	iterator := sdk.KVStorePrefixIterator(store, []byte(stakerId))
	defer iterator.Close()

	ret := make(map[string]*types2.StakerSingleAssetOrChangeInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var stateInfo types2.StakerSingleAssetOrChangeInfo
		k.cdc.MustUnmarshal(iterator.Value(), &stateInfo)
		_, assetId, err := types2.ParseStakerAndAssetIdFromKey(iterator.Key())
		if err != nil {
			return nil, err
		}
		ret[assetId] = &stateInfo
	}
	return ret, nil
}

func (k Keeper) GetStakerSpecifiedAssetInfo(ctx sdk.Context, stakerId string, assetId string) (info *types2.StakerSingleAssetOrChangeInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixReStakerAssetInfos)
	key := types2.GetAssetStateKey(stakerId, assetId)
	ifExist := store.Has(key)
	if !ifExist {
		return nil, types2.ErrNoStakerAssetKey
	}

	value := store.Get(key)

	ret := types2.StakerSingleAssetOrChangeInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

func (k Keeper) UpdateStakerAssetState(ctx sdk.Context, stakerId string, assetId string, changeAmount types2.StakerSingleAssetOrChangeInfo) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixReStakerAssetInfos)

	key := types2.GetAssetStateKey(stakerId, assetId)
	assetState := types2.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue:     math.NewInt(0),
		CanWithdrawAmountOrWantChangeValue:      math.NewInt(0),
		WaitUnDelegationAmountOrWantChangeValue: math.NewInt(0),
	}
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &assetState)
	}

	if !changeAmount.TotalDepositAmountOrWantChangeValue.IsNil() {
		if changeAmount.TotalDepositAmountOrWantChangeValue.IsNegative() {
			if assetState.TotalDepositAmountOrWantChangeValue.LT(changeAmount.TotalDepositAmountOrWantChangeValue.Abs()) {
				return errorsmod.Wrap(types2.ErrSubAmountIsMoreThanOrigin, fmt.Sprintf("TotalDepositAmount:%s,changeValue:%s", assetState.TotalDepositAmountOrWantChangeValue, changeAmount.TotalDepositAmountOrWantChangeValue))
			}
		}
		if !changeAmount.TotalDepositAmountOrWantChangeValue.IsZero() {
			assetState.TotalDepositAmountOrWantChangeValue = assetState.TotalDepositAmountOrWantChangeValue.Add(changeAmount.TotalDepositAmountOrWantChangeValue)
		}
	}

	if !changeAmount.CanWithdrawAmountOrWantChangeValue.IsNil() {
		if changeAmount.CanWithdrawAmountOrWantChangeValue.IsNegative() {
			if assetState.CanWithdrawAmountOrWantChangeValue.LT(changeAmount.CanWithdrawAmountOrWantChangeValue.Abs()) {
				return errorsmod.Wrap(types2.ErrSubAmountIsMoreThanOrigin, fmt.Sprintf("CanWithdrawAmount:%s,changeValue:%s", assetState.CanWithdrawAmountOrWantChangeValue, changeAmount.CanWithdrawAmountOrWantChangeValue))
			}
		}

		if !changeAmount.CanWithdrawAmountOrWantChangeValue.IsZero() {
			assetState.CanWithdrawAmountOrWantChangeValue = assetState.CanWithdrawAmountOrWantChangeValue.Add(changeAmount.CanWithdrawAmountOrWantChangeValue)
		}
	}

	if !changeAmount.WaitUnDelegationAmountOrWantChangeValue.IsNil() {
		if changeAmount.WaitUnDelegationAmountOrWantChangeValue.IsNegative() {
			if assetState.WaitUnDelegationAmountOrWantChangeValue.LT(changeAmount.WaitUnDelegationAmountOrWantChangeValue.Abs()) {
				return errorsmod.Wrap(types2.ErrSubAmountIsMoreThanOrigin, fmt.Sprintf("WaitUndelegationAmount:%s,changeValue:%s", assetState.WaitUnDelegationAmountOrWantChangeValue, changeAmount.WaitUnDelegationAmountOrWantChangeValue))
			}
		}

		if !changeAmount.WaitUnDelegationAmountOrWantChangeValue.IsZero() {
			assetState.WaitUnDelegationAmountOrWantChangeValue = assetState.WaitUnDelegationAmountOrWantChangeValue.Add(changeAmount.WaitUnDelegationAmountOrWantChangeValue)
		}
	}

	bz := k.cdc.MustMarshal(&assetState)
	store.Set(key, bz)

	return nil
}
