package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "dogfood"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// InitialValidatorSetID is the initial validator set id.
const InitialValidatorSetID = uint64(1)

const (
	// ExocoreValidatorBytePrefix is the prefix for the validator store.
	ExocoreValidatorBytePrefix byte = iota + 1

	// QueuedOperationsByte is the byte used to store the queue of operations.
	QueuedOperationsByte

	// OptOutsToFinishBytePrefix is the byte used to store the list of operator addresses whose
	// opt outs are maturing at the provided epoch.
	OptOutsToFinishBytePrefix

	// OperatorOptOutFinishEpochBytePrefix is the byte prefix to store the epoch at which an
	// operator's opt out will mature.
	OperatorOptOutFinishEpochBytePrefix

	// ConsensusAddrsToPruneBytePrefix is the byte prefix to store the list of consensus
	// addresses that can be pruned from the operator module at the provided epoch.
	ConsensusAddrsToPruneBytePrefix

	// UnbondingReleaseMaturityBytePrefix is the byte prefix to store the list of undelegations
	// that will mature at the provided epoch.
	UnbondingReleaseMaturityBytePrefix

	// PendingOperationsByte is the byte used to store the list of operations to be applied at
	// the end of the current block.
	PendingOperationsByte

	// PendingOptOutsByte is the byte used to store the list of operator addresses whose opt
	// outs will be made effective at the end of the current block.
	PendingOptOutsByte

	// PendingConsensusAddrsByte is the byte used to store the list of consensus addresses to be
	// pruned at the end of the block.
	PendingConsensusAddrsByte

	// PendingUndelegationsByte is the byte used to store the list of undelegations that will
	// mature at the end of the current block.
	PendingUndelegationsByte
	// ValidatorSetBytePrefix is the prefix for the historical validator set store.
	ValidatorSetBytePrefix

	// ValidatorSetIDBytePrefix is the prefix for the validator set id store.
	ValidatorSetIDBytePrefix

	// HeaderBytePrefix is the prefix for the header store.
	HeaderBytePrefix
)

// ExocoreValidatorKey returns the key for the validator store.
func ExocoreValidatorKey(address sdk.AccAddress) []byte {
	return append([]byte{ExocoreValidatorBytePrefix}, address.Bytes()...)
}

// ValidatorSetKey returns the key for the historical validator set store.
func ValidatorSetKey(id uint64) []byte {
	bz := sdk.Uint64ToBigEndian(id)
	return append([]byte{ValidatorSetBytePrefix}, bz...)
}

// ValidatorSetIDKey returns the key for the validator set id store.
func ValidatorSetIDKey(height int64) ([]byte, bool) {
	uheight, ok := SafeInt64ToUint64(height)
	if !ok {
		return nil, false
	}
	return append(
		[]byte{ValidatorSetIDBytePrefix},
		sdk.Uint64ToBigEndian(uheight)...,
	), true
}

// HeaderKey returns the key for the header store.
func HeaderKey(height int64) ([]byte, bool) {
	uheight, ok := SafeInt64ToUint64(height)
	if !ok {
		return nil, false
	}
	return append(
		[]byte{HeaderBytePrefix},
		sdk.Uint64ToBigEndian(uheight)...,
	), true
}

// QueuedOperationsKey returns the key for the queued operations store.
func QueuedOperationsKey() []byte {
	return []byte{QueuedOperationsByte}
}

// OptOutsToFinishKey returns the key for the operator opt out maturity store (epoch -> list of
// addresses).
func OptOutsToFinishKey(epoch int64) ([]byte, bool) {
	uepoch, ok := SafeInt64ToUint64(epoch)
	if !ok {
		return nil, false
	}
	return append(
		[]byte{OptOutsToFinishBytePrefix},
		sdk.Uint64ToBigEndian(uepoch)...,
	), true
}

// OperatorOptOutFinishEpochKey is the key for the operator opt out maturity store
// (sdk.AccAddress -> epoch)
func OperatorOptOutFinishEpochKey(address sdk.AccAddress) []byte {
	return append([]byte{OperatorOptOutFinishEpochBytePrefix}, address.Bytes()...)
}

// ConsensusAddrsToPruneKey is the key to lookup the list of operator consensus addresses that
// can be pruned from the operator module at the provided epoch.
func ConsensusAddrsToPruneKey(epoch int64) ([]byte, bool) {
	uepoch, ok := SafeInt64ToUint64(epoch)
	if !ok {
		return nil, false
	}
	return append(
		[]byte{ConsensusAddrsToPruneBytePrefix},
		sdk.Uint64ToBigEndian(uepoch)...,
	), true
}

// UnbondingReleaseMaturityKey is the key to lookup the list of undelegations that will mature
// at the provided epoch.
func UnbondingReleaseMaturityKey(epoch int64) ([]byte, bool) {
	uepoch, ok := SafeInt64ToUint64(epoch)
	if !ok {
		return nil, false
	}
	return append(
		[]byte{UnbondingReleaseMaturityBytePrefix},
		sdk.Uint64ToBigEndian(uepoch)...,
	), true
}

// PendingOperationsKey returns the key for the pending operations store.
func PendingOperationsKey() []byte {
	return []byte{PendingOperationsByte}
}

// PendingOptOutsKey returns the key for the pending opt-outs store.
func PendingOptOutsKey() []byte {
	return []byte{PendingOptOutsByte}
}

// PendingConsensusAddrsByte is the byte used to store the list of consensus addresses to be
// pruned at the end of the block.
func PendingConsensusAddrsKey() []byte {
	return []byte{PendingConsensusAddrsByte}
}

// PendingUndelegationsKey returns the key for the pending undelegations store.
func PendingUndelegationsKey() []byte {
	return []byte{PendingUndelegationsByte}
}

// SafeInt64ToUint64 is a wrapper function to convert an int64
// to a uint64. It returns (0, false) if the conversion is not possible.
// This is safe as long as the int64 is non-negative.
func SafeInt64ToUint64(id int64) (uint64, bool) {
	if id < 0 {
		return 0, false
	}
	return uint64(id), true // #nosec G701 // already checked.
}