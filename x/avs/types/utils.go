package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// CommitPhase represents the phases of the Two-Phase Commit protocol
type CommitPhase uint32

const (
	PreparePhase CommitPhase = iota
	DoCommitPhase
)

type OperatorOptParams struct {
	Name            string `json:"name"`
	BlsPublicKey    string `json:"bls_public_key"`
	IsRegistered    bool   `json:"is_registered"`
	Action          uint64 `json:"action"`
	OperatorAddress string `json:"operator_address"`
	Status          string `json:"status"`
	AvsAddress      string `json:"avs_address"`
}

type TaskInfoParams struct {
	TaskContractAddress   string `json:"task_contract_address"`
	TaskName              string `json:"name"`
	Hash                  []byte `json:"hash"`
	TaskID                uint64 `json:"task_id"`
	TaskResponsePeriod    uint64 `json:"task_response_period"`
	TaskStatisticalPeriod uint64 `json:"task_statistical_period"`
	TaskChallengePeriod   uint64 `json:"task_challenge_period"`
	ThresholdPercentage   uint64 `json:"threshold_percentage"`
	StartingEpoch         uint64 `json:"starting_epoch"`
	OperatorAddress       string `json:"operator_address"`
	TaskResponseHash      string `json:"task_response_hash"`
	TaskResponse          []byte `json:"task_response"`
	BlsSignature          []byte `json:"bls_signature"`
	Stage                 string `json:"stage"`
	ActualThreshold       uint64 `json:"actual_threshold"`
	OptInCount            uint64 `json:"opt_in_count"`
	SignedCount           uint64 `json:"signed_count"`
	NoSignedCount         uint64 `json:"no_signed_count"`
	ErrSignedCount        uint64 `json:"err_signed_count"`
	CallerAddress         string `json:"caller_address"`
}
type BlsParams struct {
	Operator                      string
	Name                          string
	PubKey                        []byte
	PubkeyRegistrationSignature   []byte
	PubkeyRegistrationMessageHash []byte
}

type ProofParams struct {
	TaskID              string
	TaskContractAddress string
	AvsAddress          string
	Aggregator          string
	OperatorStatus      []OperatorStatusParams
	CallerAddress       string
}
type OperatorStatusParams struct {
	OperatorAddress string
	Status          string
	ProofData       string
}

const (
	RegisterAction   = 1
	DeRegisterAction = 2
	UpdateAction     = 3
)

type ChallengeParams struct {
	TaskContractAddress common.Address `json:"task_contract_address"`
	TaskHash            []byte         `json:"hash"`
	TaskID              uint64         `json:"task_id"`
	OperatorAddress     sdk.AccAddress `json:"operator_address"`
	TaskResponseHash    []byte         `json:"task_response_hash"`
	CallerAddress       string         `json:"caller_address"`
}

type TaskResultParams struct {
	OperatorAddress     string         `json:"operator_address"`
	TaskResponseHash    string         `json:"task_response_hash"`
	TaskResponse        []byte         `json:"task_response"`
	BlsSignature        []byte         `json:"bls_signature"`
	TaskContractAddress common.Address `json:"task_contract_address"`
	TaskID              uint64         `json:"task_id"`
	Phase               uint8          `json:"phase"`
	CallerAddress       string         `json:"caller_address"`
}

type OperatorParams struct {
	EarningsAddr     string         `json:"earnings_addr"`
	ApproveAddr      string         `json:"approve_addr"`
	OperatorMetaInfo string         `json:"operator_meta_info"`
	CallerAddress    sdk.AccAddress `json:"caller_address"`
}
