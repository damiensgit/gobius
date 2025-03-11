// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package engine

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// Model is an auto generated low-level Go binding around an user-defined struct.
type Model struct {
	Fee  *big.Int
	Addr common.Address
	Rate *big.Int
	Cid  []byte
}

// Task is an auto generated low-level Go binding around an user-defined struct.
type Task struct {
	Model     [32]byte
	Fee       *big.Int
	Owner     common.Address
	Blocktime uint64
	Version   uint8
	Cid       []byte
}

// EngineMetaData contains all meta data concerning the Engine contract.
var EngineMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv18_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"denominator\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv_Overflow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PRBMath_SD59x18_Abs_MinSD59x18\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PRBMath_SD59x18_Div_InputTooSmall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"SD59x18\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"SD59x18\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"PRBMath_SD59x18_Div_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"SD59x18\",\"name\":\"x\",\"type\":\"int256\"}],\"name\":\"PRBMath_SD59x18_Exp2_InputTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PRBMath_SD59x18_Mul_InputTooSmall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"SD59x18\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"SD59x18\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"PRBMath_SD59x18_Mul_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"UD60x18\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"PRBMath_UD60x18_Exp2_InputTooBig\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"task\",\"type\":\"bytes32\"}],\"name\":\"ContestationSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"task\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"yea\",\"type\":\"bool\"}],\"name\":\"ContestationVote\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"start_idx\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"end_idx\",\"type\":\"uint32\"}],\"name\":\"ContestationVoteFinish\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"ModelAddrChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"ModelFeeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ModelRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"name\":\"PausedChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PauserTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"}],\"name\":\"SignalCommitment\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"task\",\"type\":\"bytes32\"}],\"name\":\"SolutionClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"SolutionMineableRateChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"task\",\"type\":\"bytes32\"}],\"name\":\"SolutionSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"startBlockTime\",\"type\":\"uint64\"}],\"name\":\"StartBlockTimeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"model\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"TaskSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"TreasuryTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ValidatorDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ValidatorWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"ValidatorWithdrawCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unlockTime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ValidatorWithdrawInitiated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"}],\"name\":\"VersionChanged\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"accruedFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"baseToken\",\"outputs\":[{\"internalType\":\"contractIBaseToken\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"taskids_\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes[]\",\"name\":\"cids_\",\"type\":\"bytes[]\"}],\"name\":\"bulkSubmitSolution\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"version_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"owner_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"model_\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"input_\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"n_\",\"type\":\"uint256\"}],\"name\":\"bulkSubmitTask\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"count_\",\"type\":\"uint256\"}],\"name\":\"cancelValidatorWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"}],\"name\":\"claimSolution\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"commitments\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"contestationVoteExtensionTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"amnt_\",\"type\":\"uint32\"}],\"name\":\"contestationVoteFinish\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"contestationVoteNays\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"contestationVoteYeas\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"contestationVoted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"contestationVotedIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"contestations\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"blocktime\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"finish_start_index\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"slashAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"t\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ts\",\"type\":\"uint256\"}],\"name\":\"diffMul\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"exitValidatorMinUnlockTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"cid_\",\"type\":\"bytes\"}],\"name\":\"generateCommitment\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"content_\",\"type\":\"bytes\"}],\"name\":\"generateIPFSCID\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPsuedoTotalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSlashAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidatorMinimum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cid\",\"type\":\"bytes\"}],\"internalType\":\"structModel\",\"name\":\"o_\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender_\",\"type\":\"address\"}],\"name\":\"hashModel\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"model\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"blocktime\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"cid\",\"type\":\"bytes\"}],\"internalType\":\"structTask\",\"name\":\"o_\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"prevhash_\",\"type\":\"bytes32\"}],\"name\":\"hashTask\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"}],\"name\":\"initiateValidatorWithdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"lastContestationLossTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"lastSolutionSubmission\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxContestationValidatorStakeSince\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minClaimSolutionTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minContestationVotePeriodTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minRetractionWaitTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"models\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cid\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauser\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"pendingValidatorWithdrawRequests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"unlockTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"pendingValidatorWithdrawRequestsCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"prevhash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"template_\",\"type\":\"bytes\"}],\"name\":\"registerModel\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"retractionFeePercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"t\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ts\",\"type\":\"uint256\"}],\"name\":\"reward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"model_\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"name\":\"setModelAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"model_\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"}],\"name\":\"setModelFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"paused_\",\"type\":\"bool\"}],\"name\":\"setPaused\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"model_\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"rate_\",\"type\":\"uint256\"}],\"name\":\"setSolutionMineableRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"startBlockTime_\",\"type\":\"uint64\"}],\"name\":\"setStartBlockTime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"veStaking_\",\"type\":\"address\"}],\"name\":\"setVeStaking\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"version_\",\"type\":\"uint256\"}],\"name\":\"setVersion\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"voter_\",\"type\":\"address\"}],\"name\":\"setVoter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitment_\",\"type\":\"bytes32\"}],\"name\":\"signalCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"slashAmountPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"solutionFeePercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"solutionRateLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"solutions\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"blocktime\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"claimed\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"cid\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"solutionsStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"solutionsStakeAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"startBlockTime\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"}],\"name\":\"submitContestation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"cid_\",\"type\":\"bytes\"}],\"name\":\"submitSolution\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"version_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"owner_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"model_\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"input_\",\"type\":\"bytes\"}],\"name\":\"submitTask\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"t\",\"type\":\"uint256\"}],\"name\":\"targetTs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"taskOwnerRewardPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"tasks\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"model\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"blocktime\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"cid\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalHeld\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to_\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to_\",\"type\":\"address\"}],\"name\":\"transferPauser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to_\",\"type\":\"address\"}],\"name\":\"transferTreasury\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"treasury\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"treasuryRewardPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"}],\"name\":\"validatorCanVote\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"}],\"name\":\"validatorDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validatorMinimumPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"count_\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to_\",\"type\":\"address\"}],\"name\":\"validatorWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validatorWithdrawPendingAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"staked\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"since\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"veRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"veStaking\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"yea_\",\"type\":\"bool\"}],\"name\":\"voteOnContestation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"voter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"}],\"name\":\"votingPeriodEnded\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawAccruedFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// EngineABI is the input ABI used to generate the binding from.
// Deprecated: Use EngineMetaData.ABI instead.
var EngineABI = EngineMetaData.ABI

// Engine is an auto generated Go binding around an Ethereum contract.
type Engine struct {
	EngineCaller     // Read-only binding to the contract
	EngineTransactor // Write-only binding to the contract
	EngineFilterer   // Log filterer for contract events
}

// EngineCaller is an auto generated read-only Go binding around an Ethereum contract.
type EngineCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EngineTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EngineTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EngineFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EngineFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EngineSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EngineSession struct {
	Contract     *Engine           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EngineCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EngineCallerSession struct {
	Contract *EngineCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// EngineTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EngineTransactorSession struct {
	Contract     *EngineTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EngineRaw is an auto generated low-level Go binding around an Ethereum contract.
type EngineRaw struct {
	Contract *Engine // Generic contract binding to access the raw methods on
}

// EngineCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EngineCallerRaw struct {
	Contract *EngineCaller // Generic read-only contract binding to access the raw methods on
}

// EngineTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EngineTransactorRaw struct {
	Contract *EngineTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEngine creates a new instance of Engine, bound to a specific deployed contract.
func NewEngine(address common.Address, backend bind.ContractBackend) (*Engine, error) {
	contract, err := bindEngine(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Engine{EngineCaller: EngineCaller{contract: contract}, EngineTransactor: EngineTransactor{contract: contract}, EngineFilterer: EngineFilterer{contract: contract}}, nil
}

// NewEngineCaller creates a new read-only instance of Engine, bound to a specific deployed contract.
func NewEngineCaller(address common.Address, caller bind.ContractCaller) (*EngineCaller, error) {
	contract, err := bindEngine(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EngineCaller{contract: contract}, nil
}

// NewEngineTransactor creates a new write-only instance of Engine, bound to a specific deployed contract.
func NewEngineTransactor(address common.Address, transactor bind.ContractTransactor) (*EngineTransactor, error) {
	contract, err := bindEngine(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EngineTransactor{contract: contract}, nil
}

// NewEngineFilterer creates a new log filterer instance of Engine, bound to a specific deployed contract.
func NewEngineFilterer(address common.Address, filterer bind.ContractFilterer) (*EngineFilterer, error) {
	contract, err := bindEngine(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EngineFilterer{contract: contract}, nil
}

// bindEngine binds a generic wrapper to an already deployed contract.
func bindEngine(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EngineMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Engine *EngineRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Engine.Contract.EngineCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Engine *EngineRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Engine.Contract.EngineTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Engine *EngineRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Engine.Contract.EngineTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Engine *EngineCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Engine.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Engine *EngineTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Engine.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Engine *EngineTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Engine.Contract.contract.Transact(opts, method, params...)
}

// AccruedFees is a free data retrieval call binding the contract method 0x682c2058.
//
// Solidity: function accruedFees() view returns(uint256)
func (_Engine *EngineCaller) AccruedFees(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "accruedFees")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AccruedFees is a free data retrieval call binding the contract method 0x682c2058.
//
// Solidity: function accruedFees() view returns(uint256)
func (_Engine *EngineSession) AccruedFees() (*big.Int, error) {
	return _Engine.Contract.AccruedFees(&_Engine.CallOpts)
}

// AccruedFees is a free data retrieval call binding the contract method 0x682c2058.
//
// Solidity: function accruedFees() view returns(uint256)
func (_Engine *EngineCallerSession) AccruedFees() (*big.Int, error) {
	return _Engine.Contract.AccruedFees(&_Engine.CallOpts)
}

// BaseToken is a free data retrieval call binding the contract method 0xc55dae63.
//
// Solidity: function baseToken() view returns(address)
func (_Engine *EngineCaller) BaseToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "baseToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BaseToken is a free data retrieval call binding the contract method 0xc55dae63.
//
// Solidity: function baseToken() view returns(address)
func (_Engine *EngineSession) BaseToken() (common.Address, error) {
	return _Engine.Contract.BaseToken(&_Engine.CallOpts)
}

// BaseToken is a free data retrieval call binding the contract method 0xc55dae63.
//
// Solidity: function baseToken() view returns(address)
func (_Engine *EngineCallerSession) BaseToken() (common.Address, error) {
	return _Engine.Contract.BaseToken(&_Engine.CallOpts)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(uint256)
func (_Engine *EngineCaller) Commitments(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "commitments", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(uint256)
func (_Engine *EngineSession) Commitments(arg0 [32]byte) (*big.Int, error) {
	return _Engine.Contract.Commitments(&_Engine.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(uint256)
func (_Engine *EngineCallerSession) Commitments(arg0 [32]byte) (*big.Int, error) {
	return _Engine.Contract.Commitments(&_Engine.CallOpts, arg0)
}

// ContestationVoteExtensionTime is a free data retrieval call binding the contract method 0xa2492a90.
//
// Solidity: function contestationVoteExtensionTime() view returns(uint256)
func (_Engine *EngineCaller) ContestationVoteExtensionTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "contestationVoteExtensionTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ContestationVoteExtensionTime is a free data retrieval call binding the contract method 0xa2492a90.
//
// Solidity: function contestationVoteExtensionTime() view returns(uint256)
func (_Engine *EngineSession) ContestationVoteExtensionTime() (*big.Int, error) {
	return _Engine.Contract.ContestationVoteExtensionTime(&_Engine.CallOpts)
}

// ContestationVoteExtensionTime is a free data retrieval call binding the contract method 0xa2492a90.
//
// Solidity: function contestationVoteExtensionTime() view returns(uint256)
func (_Engine *EngineCallerSession) ContestationVoteExtensionTime() (*big.Int, error) {
	return _Engine.Contract.ContestationVoteExtensionTime(&_Engine.CallOpts)
}

// ContestationVoteNays is a free data retrieval call binding the contract method 0x303fb0d6.
//
// Solidity: function contestationVoteNays(bytes32 , uint256 ) view returns(address)
func (_Engine *EngineCaller) ContestationVoteNays(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "contestationVoteNays", arg0, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ContestationVoteNays is a free data retrieval call binding the contract method 0x303fb0d6.
//
// Solidity: function contestationVoteNays(bytes32 , uint256 ) view returns(address)
func (_Engine *EngineSession) ContestationVoteNays(arg0 [32]byte, arg1 *big.Int) (common.Address, error) {
	return _Engine.Contract.ContestationVoteNays(&_Engine.CallOpts, arg0, arg1)
}

// ContestationVoteNays is a free data retrieval call binding the contract method 0x303fb0d6.
//
// Solidity: function contestationVoteNays(bytes32 , uint256 ) view returns(address)
func (_Engine *EngineCallerSession) ContestationVoteNays(arg0 [32]byte, arg1 *big.Int) (common.Address, error) {
	return _Engine.Contract.ContestationVoteNays(&_Engine.CallOpts, arg0, arg1)
}

// ContestationVoteYeas is a free data retrieval call binding the contract method 0xd1f0c941.
//
// Solidity: function contestationVoteYeas(bytes32 , uint256 ) view returns(address)
func (_Engine *EngineCaller) ContestationVoteYeas(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "contestationVoteYeas", arg0, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ContestationVoteYeas is a free data retrieval call binding the contract method 0xd1f0c941.
//
// Solidity: function contestationVoteYeas(bytes32 , uint256 ) view returns(address)
func (_Engine *EngineSession) ContestationVoteYeas(arg0 [32]byte, arg1 *big.Int) (common.Address, error) {
	return _Engine.Contract.ContestationVoteYeas(&_Engine.CallOpts, arg0, arg1)
}

// ContestationVoteYeas is a free data retrieval call binding the contract method 0xd1f0c941.
//
// Solidity: function contestationVoteYeas(bytes32 , uint256 ) view returns(address)
func (_Engine *EngineCallerSession) ContestationVoteYeas(arg0 [32]byte, arg1 *big.Int) (common.Address, error) {
	return _Engine.Contract.ContestationVoteYeas(&_Engine.CallOpts, arg0, arg1)
}

// ContestationVoted is a free data retrieval call binding the contract method 0xd2780940.
//
// Solidity: function contestationVoted(bytes32 , address ) view returns(bool)
func (_Engine *EngineCaller) ContestationVoted(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "contestationVoted", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ContestationVoted is a free data retrieval call binding the contract method 0xd2780940.
//
// Solidity: function contestationVoted(bytes32 , address ) view returns(bool)
func (_Engine *EngineSession) ContestationVoted(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _Engine.Contract.ContestationVoted(&_Engine.CallOpts, arg0, arg1)
}

// ContestationVoted is a free data retrieval call binding the contract method 0xd2780940.
//
// Solidity: function contestationVoted(bytes32 , address ) view returns(bool)
func (_Engine *EngineCallerSession) ContestationVoted(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _Engine.Contract.ContestationVoted(&_Engine.CallOpts, arg0, arg1)
}

// ContestationVotedIndex is a free data retrieval call binding the contract method 0x17f3e041.
//
// Solidity: function contestationVotedIndex(bytes32 ) view returns(uint256)
func (_Engine *EngineCaller) ContestationVotedIndex(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "contestationVotedIndex", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ContestationVotedIndex is a free data retrieval call binding the contract method 0x17f3e041.
//
// Solidity: function contestationVotedIndex(bytes32 ) view returns(uint256)
func (_Engine *EngineSession) ContestationVotedIndex(arg0 [32]byte) (*big.Int, error) {
	return _Engine.Contract.ContestationVotedIndex(&_Engine.CallOpts, arg0)
}

// ContestationVotedIndex is a free data retrieval call binding the contract method 0x17f3e041.
//
// Solidity: function contestationVotedIndex(bytes32 ) view returns(uint256)
func (_Engine *EngineCallerSession) ContestationVotedIndex(arg0 [32]byte) (*big.Int, error) {
	return _Engine.Contract.ContestationVotedIndex(&_Engine.CallOpts, arg0)
}

// Contestations is a free data retrieval call binding the contract method 0xd33b2ef5.
//
// Solidity: function contestations(bytes32 ) view returns(address validator, uint64 blocktime, uint32 finish_start_index, uint256 slashAmount)
func (_Engine *EngineCaller) Contestations(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Validator        common.Address
	Blocktime        uint64
	FinishStartIndex uint32
	SlashAmount      *big.Int
}, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "contestations", arg0)

	outstruct := new(struct {
		Validator        common.Address
		Blocktime        uint64
		FinishStartIndex uint32
		SlashAmount      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Blocktime = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.FinishStartIndex = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.SlashAmount = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Contestations is a free data retrieval call binding the contract method 0xd33b2ef5.
//
// Solidity: function contestations(bytes32 ) view returns(address validator, uint64 blocktime, uint32 finish_start_index, uint256 slashAmount)
func (_Engine *EngineSession) Contestations(arg0 [32]byte) (struct {
	Validator        common.Address
	Blocktime        uint64
	FinishStartIndex uint32
	SlashAmount      *big.Int
}, error) {
	return _Engine.Contract.Contestations(&_Engine.CallOpts, arg0)
}

// Contestations is a free data retrieval call binding the contract method 0xd33b2ef5.
//
// Solidity: function contestations(bytes32 ) view returns(address validator, uint64 blocktime, uint32 finish_start_index, uint256 slashAmount)
func (_Engine *EngineCallerSession) Contestations(arg0 [32]byte) (struct {
	Validator        common.Address
	Blocktime        uint64
	FinishStartIndex uint32
	SlashAmount      *big.Int
}, error) {
	return _Engine.Contract.Contestations(&_Engine.CallOpts, arg0)
}

// DiffMul is a free data retrieval call binding the contract method 0x1f88ea1c.
//
// Solidity: function diffMul(uint256 t, uint256 ts) pure returns(uint256)
func (_Engine *EngineCaller) DiffMul(opts *bind.CallOpts, t *big.Int, ts *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "diffMul", t, ts)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DiffMul is a free data retrieval call binding the contract method 0x1f88ea1c.
//
// Solidity: function diffMul(uint256 t, uint256 ts) pure returns(uint256)
func (_Engine *EngineSession) DiffMul(t *big.Int, ts *big.Int) (*big.Int, error) {
	return _Engine.Contract.DiffMul(&_Engine.CallOpts, t, ts)
}

// DiffMul is a free data retrieval call binding the contract method 0x1f88ea1c.
//
// Solidity: function diffMul(uint256 t, uint256 ts) pure returns(uint256)
func (_Engine *EngineCallerSession) DiffMul(t *big.Int, ts *big.Int) (*big.Int, error) {
	return _Engine.Contract.DiffMul(&_Engine.CallOpts, t, ts)
}

// ExitValidatorMinUnlockTime is a free data retrieval call binding the contract method 0xa53e2525.
//
// Solidity: function exitValidatorMinUnlockTime() view returns(uint256)
func (_Engine *EngineCaller) ExitValidatorMinUnlockTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "exitValidatorMinUnlockTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ExitValidatorMinUnlockTime is a free data retrieval call binding the contract method 0xa53e2525.
//
// Solidity: function exitValidatorMinUnlockTime() view returns(uint256)
func (_Engine *EngineSession) ExitValidatorMinUnlockTime() (*big.Int, error) {
	return _Engine.Contract.ExitValidatorMinUnlockTime(&_Engine.CallOpts)
}

// ExitValidatorMinUnlockTime is a free data retrieval call binding the contract method 0xa53e2525.
//
// Solidity: function exitValidatorMinUnlockTime() view returns(uint256)
func (_Engine *EngineCallerSession) ExitValidatorMinUnlockTime() (*big.Int, error) {
	return _Engine.Contract.ExitValidatorMinUnlockTime(&_Engine.CallOpts)
}

// GenerateCommitment is a free data retrieval call binding the contract method 0x393cb1c7.
//
// Solidity: function generateCommitment(address sender_, bytes32 taskid_, bytes cid_) pure returns(bytes32)
func (_Engine *EngineCaller) GenerateCommitment(opts *bind.CallOpts, sender_ common.Address, taskid_ [32]byte, cid_ []byte) ([32]byte, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "generateCommitment", sender_, taskid_, cid_)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GenerateCommitment is a free data retrieval call binding the contract method 0x393cb1c7.
//
// Solidity: function generateCommitment(address sender_, bytes32 taskid_, bytes cid_) pure returns(bytes32)
func (_Engine *EngineSession) GenerateCommitment(sender_ common.Address, taskid_ [32]byte, cid_ []byte) ([32]byte, error) {
	return _Engine.Contract.GenerateCommitment(&_Engine.CallOpts, sender_, taskid_, cid_)
}

// GenerateCommitment is a free data retrieval call binding the contract method 0x393cb1c7.
//
// Solidity: function generateCommitment(address sender_, bytes32 taskid_, bytes cid_) pure returns(bytes32)
func (_Engine *EngineCallerSession) GenerateCommitment(sender_ common.Address, taskid_ [32]byte, cid_ []byte) ([32]byte, error) {
	return _Engine.Contract.GenerateCommitment(&_Engine.CallOpts, sender_, taskid_, cid_)
}

// GenerateIPFSCID is a free data retrieval call binding the contract method 0x40e8c56d.
//
// Solidity: function generateIPFSCID(bytes content_) pure returns(bytes)
func (_Engine *EngineCaller) GenerateIPFSCID(opts *bind.CallOpts, content_ []byte) ([]byte, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "generateIPFSCID", content_)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GenerateIPFSCID is a free data retrieval call binding the contract method 0x40e8c56d.
//
// Solidity: function generateIPFSCID(bytes content_) pure returns(bytes)
func (_Engine *EngineSession) GenerateIPFSCID(content_ []byte) ([]byte, error) {
	return _Engine.Contract.GenerateIPFSCID(&_Engine.CallOpts, content_)
}

// GenerateIPFSCID is a free data retrieval call binding the contract method 0x40e8c56d.
//
// Solidity: function generateIPFSCID(bytes content_) pure returns(bytes)
func (_Engine *EngineCallerSession) GenerateIPFSCID(content_ []byte) ([]byte, error) {
	return _Engine.Contract.GenerateIPFSCID(&_Engine.CallOpts, content_)
}

// GetPsuedoTotalSupply is a free data retrieval call binding the contract method 0x7881c5e6.
//
// Solidity: function getPsuedoTotalSupply() view returns(uint256)
func (_Engine *EngineCaller) GetPsuedoTotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "getPsuedoTotalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPsuedoTotalSupply is a free data retrieval call binding the contract method 0x7881c5e6.
//
// Solidity: function getPsuedoTotalSupply() view returns(uint256)
func (_Engine *EngineSession) GetPsuedoTotalSupply() (*big.Int, error) {
	return _Engine.Contract.GetPsuedoTotalSupply(&_Engine.CallOpts)
}

// GetPsuedoTotalSupply is a free data retrieval call binding the contract method 0x7881c5e6.
//
// Solidity: function getPsuedoTotalSupply() view returns(uint256)
func (_Engine *EngineCallerSession) GetPsuedoTotalSupply() (*big.Int, error) {
	return _Engine.Contract.GetPsuedoTotalSupply(&_Engine.CallOpts)
}

// GetReward is a free data retrieval call binding the contract method 0x3d18b912.
//
// Solidity: function getReward() view returns(uint256)
func (_Engine *EngineCaller) GetReward(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "getReward")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetReward is a free data retrieval call binding the contract method 0x3d18b912.
//
// Solidity: function getReward() view returns(uint256)
func (_Engine *EngineSession) GetReward() (*big.Int, error) {
	return _Engine.Contract.GetReward(&_Engine.CallOpts)
}

// GetReward is a free data retrieval call binding the contract method 0x3d18b912.
//
// Solidity: function getReward() view returns(uint256)
func (_Engine *EngineCallerSession) GetReward() (*big.Int, error) {
	return _Engine.Contract.GetReward(&_Engine.CallOpts)
}

// GetSlashAmount is a free data retrieval call binding the contract method 0x3d57f5d9.
//
// Solidity: function getSlashAmount() view returns(uint256)
func (_Engine *EngineCaller) GetSlashAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "getSlashAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetSlashAmount is a free data retrieval call binding the contract method 0x3d57f5d9.
//
// Solidity: function getSlashAmount() view returns(uint256)
func (_Engine *EngineSession) GetSlashAmount() (*big.Int, error) {
	return _Engine.Contract.GetSlashAmount(&_Engine.CallOpts)
}

// GetSlashAmount is a free data retrieval call binding the contract method 0x3d57f5d9.
//
// Solidity: function getSlashAmount() view returns(uint256)
func (_Engine *EngineCallerSession) GetSlashAmount() (*big.Int, error) {
	return _Engine.Contract.GetSlashAmount(&_Engine.CallOpts)
}

// GetValidatorMinimum is a free data retrieval call binding the contract method 0x2258d105.
//
// Solidity: function getValidatorMinimum() view returns(uint256)
func (_Engine *EngineCaller) GetValidatorMinimum(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "getValidatorMinimum")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetValidatorMinimum is a free data retrieval call binding the contract method 0x2258d105.
//
// Solidity: function getValidatorMinimum() view returns(uint256)
func (_Engine *EngineSession) GetValidatorMinimum() (*big.Int, error) {
	return _Engine.Contract.GetValidatorMinimum(&_Engine.CallOpts)
}

// GetValidatorMinimum is a free data retrieval call binding the contract method 0x2258d105.
//
// Solidity: function getValidatorMinimum() view returns(uint256)
func (_Engine *EngineCallerSession) GetValidatorMinimum() (*big.Int, error) {
	return _Engine.Contract.GetValidatorMinimum(&_Engine.CallOpts)
}

// HashModel is a free data retrieval call binding the contract method 0x218a3048.
//
// Solidity: function hashModel((uint256,address,uint256,bytes) o_, address sender_) pure returns(bytes32)
func (_Engine *EngineCaller) HashModel(opts *bind.CallOpts, o_ Model, sender_ common.Address) ([32]byte, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "hashModel", o_, sender_)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashModel is a free data retrieval call binding the contract method 0x218a3048.
//
// Solidity: function hashModel((uint256,address,uint256,bytes) o_, address sender_) pure returns(bytes32)
func (_Engine *EngineSession) HashModel(o_ Model, sender_ common.Address) ([32]byte, error) {
	return _Engine.Contract.HashModel(&_Engine.CallOpts, o_, sender_)
}

// HashModel is a free data retrieval call binding the contract method 0x218a3048.
//
// Solidity: function hashModel((uint256,address,uint256,bytes) o_, address sender_) pure returns(bytes32)
func (_Engine *EngineCallerSession) HashModel(o_ Model, sender_ common.Address) ([32]byte, error) {
	return _Engine.Contract.HashModel(&_Engine.CallOpts, o_, sender_)
}

// HashTask is a free data retrieval call binding the contract method 0x1466b63a.
//
// Solidity: function hashTask((bytes32,uint256,address,uint64,uint8,bytes) o_, address sender_, bytes32 prevhash_) pure returns(bytes32)
func (_Engine *EngineCaller) HashTask(opts *bind.CallOpts, o_ Task, sender_ common.Address, prevhash_ [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "hashTask", o_, sender_, prevhash_)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashTask is a free data retrieval call binding the contract method 0x1466b63a.
//
// Solidity: function hashTask((bytes32,uint256,address,uint64,uint8,bytes) o_, address sender_, bytes32 prevhash_) pure returns(bytes32)
func (_Engine *EngineSession) HashTask(o_ Task, sender_ common.Address, prevhash_ [32]byte) ([32]byte, error) {
	return _Engine.Contract.HashTask(&_Engine.CallOpts, o_, sender_, prevhash_)
}

// HashTask is a free data retrieval call binding the contract method 0x1466b63a.
//
// Solidity: function hashTask((bytes32,uint256,address,uint64,uint8,bytes) o_, address sender_, bytes32 prevhash_) pure returns(bytes32)
func (_Engine *EngineCallerSession) HashTask(o_ Task, sender_ common.Address, prevhash_ [32]byte) ([32]byte, error) {
	return _Engine.Contract.HashTask(&_Engine.CallOpts, o_, sender_, prevhash_)
}

// LastContestationLossTime is a free data retrieval call binding the contract method 0x218e6859.
//
// Solidity: function lastContestationLossTime(address ) view returns(uint256)
func (_Engine *EngineCaller) LastContestationLossTime(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "lastContestationLossTime", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastContestationLossTime is a free data retrieval call binding the contract method 0x218e6859.
//
// Solidity: function lastContestationLossTime(address ) view returns(uint256)
func (_Engine *EngineSession) LastContestationLossTime(arg0 common.Address) (*big.Int, error) {
	return _Engine.Contract.LastContestationLossTime(&_Engine.CallOpts, arg0)
}

// LastContestationLossTime is a free data retrieval call binding the contract method 0x218e6859.
//
// Solidity: function lastContestationLossTime(address ) view returns(uint256)
func (_Engine *EngineCallerSession) LastContestationLossTime(arg0 common.Address) (*big.Int, error) {
	return _Engine.Contract.LastContestationLossTime(&_Engine.CallOpts, arg0)
}

// LastSolutionSubmission is a free data retrieval call binding the contract method 0xc24b5631.
//
// Solidity: function lastSolutionSubmission(address ) view returns(uint256)
func (_Engine *EngineCaller) LastSolutionSubmission(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "lastSolutionSubmission", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastSolutionSubmission is a free data retrieval call binding the contract method 0xc24b5631.
//
// Solidity: function lastSolutionSubmission(address ) view returns(uint256)
func (_Engine *EngineSession) LastSolutionSubmission(arg0 common.Address) (*big.Int, error) {
	return _Engine.Contract.LastSolutionSubmission(&_Engine.CallOpts, arg0)
}

// LastSolutionSubmission is a free data retrieval call binding the contract method 0xc24b5631.
//
// Solidity: function lastSolutionSubmission(address ) view returns(uint256)
func (_Engine *EngineCallerSession) LastSolutionSubmission(arg0 common.Address) (*big.Int, error) {
	return _Engine.Contract.LastSolutionSubmission(&_Engine.CallOpts, arg0)
}

// MaxContestationValidatorStakeSince is a free data retrieval call binding the contract method 0x8e6d86fd.
//
// Solidity: function maxContestationValidatorStakeSince() view returns(uint256)
func (_Engine *EngineCaller) MaxContestationValidatorStakeSince(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "maxContestationValidatorStakeSince")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxContestationValidatorStakeSince is a free data retrieval call binding the contract method 0x8e6d86fd.
//
// Solidity: function maxContestationValidatorStakeSince() view returns(uint256)
func (_Engine *EngineSession) MaxContestationValidatorStakeSince() (*big.Int, error) {
	return _Engine.Contract.MaxContestationValidatorStakeSince(&_Engine.CallOpts)
}

// MaxContestationValidatorStakeSince is a free data retrieval call binding the contract method 0x8e6d86fd.
//
// Solidity: function maxContestationValidatorStakeSince() view returns(uint256)
func (_Engine *EngineCallerSession) MaxContestationValidatorStakeSince() (*big.Int, error) {
	return _Engine.Contract.MaxContestationValidatorStakeSince(&_Engine.CallOpts)
}

// MinClaimSolutionTime is a free data retrieval call binding the contract method 0x92809444.
//
// Solidity: function minClaimSolutionTime() view returns(uint256)
func (_Engine *EngineCaller) MinClaimSolutionTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "minClaimSolutionTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinClaimSolutionTime is a free data retrieval call binding the contract method 0x92809444.
//
// Solidity: function minClaimSolutionTime() view returns(uint256)
func (_Engine *EngineSession) MinClaimSolutionTime() (*big.Int, error) {
	return _Engine.Contract.MinClaimSolutionTime(&_Engine.CallOpts)
}

// MinClaimSolutionTime is a free data retrieval call binding the contract method 0x92809444.
//
// Solidity: function minClaimSolutionTime() view returns(uint256)
func (_Engine *EngineCallerSession) MinClaimSolutionTime() (*big.Int, error) {
	return _Engine.Contract.MinClaimSolutionTime(&_Engine.CallOpts)
}

// MinContestationVotePeriodTime is a free data retrieval call binding the contract method 0x7b36006a.
//
// Solidity: function minContestationVotePeriodTime() view returns(uint256)
func (_Engine *EngineCaller) MinContestationVotePeriodTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "minContestationVotePeriodTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinContestationVotePeriodTime is a free data retrieval call binding the contract method 0x7b36006a.
//
// Solidity: function minContestationVotePeriodTime() view returns(uint256)
func (_Engine *EngineSession) MinContestationVotePeriodTime() (*big.Int, error) {
	return _Engine.Contract.MinContestationVotePeriodTime(&_Engine.CallOpts)
}

// MinContestationVotePeriodTime is a free data retrieval call binding the contract method 0x7b36006a.
//
// Solidity: function minContestationVotePeriodTime() view returns(uint256)
func (_Engine *EngineCallerSession) MinContestationVotePeriodTime() (*big.Int, error) {
	return _Engine.Contract.MinContestationVotePeriodTime(&_Engine.CallOpts)
}

// MinRetractionWaitTime is a free data retrieval call binding the contract method 0x00fd7082.
//
// Solidity: function minRetractionWaitTime() view returns(uint256)
func (_Engine *EngineCaller) MinRetractionWaitTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "minRetractionWaitTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinRetractionWaitTime is a free data retrieval call binding the contract method 0x00fd7082.
//
// Solidity: function minRetractionWaitTime() view returns(uint256)
func (_Engine *EngineSession) MinRetractionWaitTime() (*big.Int, error) {
	return _Engine.Contract.MinRetractionWaitTime(&_Engine.CallOpts)
}

// MinRetractionWaitTime is a free data retrieval call binding the contract method 0x00fd7082.
//
// Solidity: function minRetractionWaitTime() view returns(uint256)
func (_Engine *EngineCallerSession) MinRetractionWaitTime() (*big.Int, error) {
	return _Engine.Contract.MinRetractionWaitTime(&_Engine.CallOpts)
}

// Models is a free data retrieval call binding the contract method 0xe236f46b.
//
// Solidity: function models(bytes32 ) view returns(uint256 fee, address addr, uint256 rate, bytes cid)
func (_Engine *EngineCaller) Models(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Fee  *big.Int
	Addr common.Address
	Rate *big.Int
	Cid  []byte
}, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "models", arg0)

	outstruct := new(struct {
		Fee  *big.Int
		Addr common.Address
		Rate *big.Int
		Cid  []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fee = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Addr = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Rate = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Cid = *abi.ConvertType(out[3], new([]byte)).(*[]byte)

	return *outstruct, err

}

// Models is a free data retrieval call binding the contract method 0xe236f46b.
//
// Solidity: function models(bytes32 ) view returns(uint256 fee, address addr, uint256 rate, bytes cid)
func (_Engine *EngineSession) Models(arg0 [32]byte) (struct {
	Fee  *big.Int
	Addr common.Address
	Rate *big.Int
	Cid  []byte
}, error) {
	return _Engine.Contract.Models(&_Engine.CallOpts, arg0)
}

// Models is a free data retrieval call binding the contract method 0xe236f46b.
//
// Solidity: function models(bytes32 ) view returns(uint256 fee, address addr, uint256 rate, bytes cid)
func (_Engine *EngineCallerSession) Models(arg0 [32]byte) (struct {
	Fee  *big.Int
	Addr common.Address
	Rate *big.Int
	Cid  []byte
}, error) {
	return _Engine.Contract.Models(&_Engine.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Engine *EngineCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Engine *EngineSession) Owner() (common.Address, error) {
	return _Engine.Contract.Owner(&_Engine.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Engine *EngineCallerSession) Owner() (common.Address, error) {
	return _Engine.Contract.Owner(&_Engine.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Engine *EngineCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Engine *EngineSession) Paused() (bool, error) {
	return _Engine.Contract.Paused(&_Engine.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Engine *EngineCallerSession) Paused() (bool, error) {
	return _Engine.Contract.Paused(&_Engine.CallOpts)
}

// Pauser is a free data retrieval call binding the contract method 0x9fd0506d.
//
// Solidity: function pauser() view returns(address)
func (_Engine *EngineCaller) Pauser(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "pauser")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Pauser is a free data retrieval call binding the contract method 0x9fd0506d.
//
// Solidity: function pauser() view returns(address)
func (_Engine *EngineSession) Pauser() (common.Address, error) {
	return _Engine.Contract.Pauser(&_Engine.CallOpts)
}

// Pauser is a free data retrieval call binding the contract method 0x9fd0506d.
//
// Solidity: function pauser() view returns(address)
func (_Engine *EngineCallerSession) Pauser() (common.Address, error) {
	return _Engine.Contract.Pauser(&_Engine.CallOpts)
}

// PendingValidatorWithdrawRequests is a free data retrieval call binding the contract method 0xd2992baa.
//
// Solidity: function pendingValidatorWithdrawRequests(address , uint256 ) view returns(uint256 unlockTime, uint256 amount)
func (_Engine *EngineCaller) PendingValidatorWithdrawRequests(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (struct {
	UnlockTime *big.Int
	Amount     *big.Int
}, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "pendingValidatorWithdrawRequests", arg0, arg1)

	outstruct := new(struct {
		UnlockTime *big.Int
		Amount     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.UnlockTime = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// PendingValidatorWithdrawRequests is a free data retrieval call binding the contract method 0xd2992baa.
//
// Solidity: function pendingValidatorWithdrawRequests(address , uint256 ) view returns(uint256 unlockTime, uint256 amount)
func (_Engine *EngineSession) PendingValidatorWithdrawRequests(arg0 common.Address, arg1 *big.Int) (struct {
	UnlockTime *big.Int
	Amount     *big.Int
}, error) {
	return _Engine.Contract.PendingValidatorWithdrawRequests(&_Engine.CallOpts, arg0, arg1)
}

// PendingValidatorWithdrawRequests is a free data retrieval call binding the contract method 0xd2992baa.
//
// Solidity: function pendingValidatorWithdrawRequests(address , uint256 ) view returns(uint256 unlockTime, uint256 amount)
func (_Engine *EngineCallerSession) PendingValidatorWithdrawRequests(arg0 common.Address, arg1 *big.Int) (struct {
	UnlockTime *big.Int
	Amount     *big.Int
}, error) {
	return _Engine.Contract.PendingValidatorWithdrawRequests(&_Engine.CallOpts, arg0, arg1)
}

// PendingValidatorWithdrawRequestsCount is a free data retrieval call binding the contract method 0xd2307ae4.
//
// Solidity: function pendingValidatorWithdrawRequestsCount(address ) view returns(uint256)
func (_Engine *EngineCaller) PendingValidatorWithdrawRequestsCount(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "pendingValidatorWithdrawRequestsCount", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PendingValidatorWithdrawRequestsCount is a free data retrieval call binding the contract method 0xd2307ae4.
//
// Solidity: function pendingValidatorWithdrawRequestsCount(address ) view returns(uint256)
func (_Engine *EngineSession) PendingValidatorWithdrawRequestsCount(arg0 common.Address) (*big.Int, error) {
	return _Engine.Contract.PendingValidatorWithdrawRequestsCount(&_Engine.CallOpts, arg0)
}

// PendingValidatorWithdrawRequestsCount is a free data retrieval call binding the contract method 0xd2307ae4.
//
// Solidity: function pendingValidatorWithdrawRequestsCount(address ) view returns(uint256)
func (_Engine *EngineCallerSession) PendingValidatorWithdrawRequestsCount(arg0 common.Address) (*big.Int, error) {
	return _Engine.Contract.PendingValidatorWithdrawRequestsCount(&_Engine.CallOpts, arg0)
}

// Prevhash is a free data retrieval call binding the contract method 0xc17ddb2a.
//
// Solidity: function prevhash() view returns(bytes32)
func (_Engine *EngineCaller) Prevhash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "prevhash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Prevhash is a free data retrieval call binding the contract method 0xc17ddb2a.
//
// Solidity: function prevhash() view returns(bytes32)
func (_Engine *EngineSession) Prevhash() ([32]byte, error) {
	return _Engine.Contract.Prevhash(&_Engine.CallOpts)
}

// Prevhash is a free data retrieval call binding the contract method 0xc17ddb2a.
//
// Solidity: function prevhash() view returns(bytes32)
func (_Engine *EngineCallerSession) Prevhash() ([32]byte, error) {
	return _Engine.Contract.Prevhash(&_Engine.CallOpts)
}

// RetractionFeePercentage is a free data retrieval call binding the contract method 0x72dc0ee1.
//
// Solidity: function retractionFeePercentage() view returns(uint256)
func (_Engine *EngineCaller) RetractionFeePercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "retractionFeePercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RetractionFeePercentage is a free data retrieval call binding the contract method 0x72dc0ee1.
//
// Solidity: function retractionFeePercentage() view returns(uint256)
func (_Engine *EngineSession) RetractionFeePercentage() (*big.Int, error) {
	return _Engine.Contract.RetractionFeePercentage(&_Engine.CallOpts)
}

// RetractionFeePercentage is a free data retrieval call binding the contract method 0x72dc0ee1.
//
// Solidity: function retractionFeePercentage() view returns(uint256)
func (_Engine *EngineCallerSession) RetractionFeePercentage() (*big.Int, error) {
	return _Engine.Contract.RetractionFeePercentage(&_Engine.CallOpts)
}

// Reward is a free data retrieval call binding the contract method 0xa4fa8d57.
//
// Solidity: function reward(uint256 t, uint256 ts) pure returns(uint256)
func (_Engine *EngineCaller) Reward(opts *bind.CallOpts, t *big.Int, ts *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "reward", t, ts)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Reward is a free data retrieval call binding the contract method 0xa4fa8d57.
//
// Solidity: function reward(uint256 t, uint256 ts) pure returns(uint256)
func (_Engine *EngineSession) Reward(t *big.Int, ts *big.Int) (*big.Int, error) {
	return _Engine.Contract.Reward(&_Engine.CallOpts, t, ts)
}

// Reward is a free data retrieval call binding the contract method 0xa4fa8d57.
//
// Solidity: function reward(uint256 t, uint256 ts) pure returns(uint256)
func (_Engine *EngineCallerSession) Reward(t *big.Int, ts *big.Int) (*big.Int, error) {
	return _Engine.Contract.Reward(&_Engine.CallOpts, t, ts)
}

// SlashAmountPercentage is a free data retrieval call binding the contract method 0xdc06a89f.
//
// Solidity: function slashAmountPercentage() view returns(uint256)
func (_Engine *EngineCaller) SlashAmountPercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "slashAmountPercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SlashAmountPercentage is a free data retrieval call binding the contract method 0xdc06a89f.
//
// Solidity: function slashAmountPercentage() view returns(uint256)
func (_Engine *EngineSession) SlashAmountPercentage() (*big.Int, error) {
	return _Engine.Contract.SlashAmountPercentage(&_Engine.CallOpts)
}

// SlashAmountPercentage is a free data retrieval call binding the contract method 0xdc06a89f.
//
// Solidity: function slashAmountPercentage() view returns(uint256)
func (_Engine *EngineCallerSession) SlashAmountPercentage() (*big.Int, error) {
	return _Engine.Contract.SlashAmountPercentage(&_Engine.CallOpts)
}

// SolutionFeePercentage is a free data retrieval call binding the contract method 0xf1b8989d.
//
// Solidity: function solutionFeePercentage() view returns(uint256)
func (_Engine *EngineCaller) SolutionFeePercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "solutionFeePercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SolutionFeePercentage is a free data retrieval call binding the contract method 0xf1b8989d.
//
// Solidity: function solutionFeePercentage() view returns(uint256)
func (_Engine *EngineSession) SolutionFeePercentage() (*big.Int, error) {
	return _Engine.Contract.SolutionFeePercentage(&_Engine.CallOpts)
}

// SolutionFeePercentage is a free data retrieval call binding the contract method 0xf1b8989d.
//
// Solidity: function solutionFeePercentage() view returns(uint256)
func (_Engine *EngineCallerSession) SolutionFeePercentage() (*big.Int, error) {
	return _Engine.Contract.SolutionFeePercentage(&_Engine.CallOpts)
}

// SolutionRateLimit is a free data retrieval call binding the contract method 0xa1975adf.
//
// Solidity: function solutionRateLimit() view returns(uint256)
func (_Engine *EngineCaller) SolutionRateLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "solutionRateLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SolutionRateLimit is a free data retrieval call binding the contract method 0xa1975adf.
//
// Solidity: function solutionRateLimit() view returns(uint256)
func (_Engine *EngineSession) SolutionRateLimit() (*big.Int, error) {
	return _Engine.Contract.SolutionRateLimit(&_Engine.CallOpts)
}

// SolutionRateLimit is a free data retrieval call binding the contract method 0xa1975adf.
//
// Solidity: function solutionRateLimit() view returns(uint256)
func (_Engine *EngineCallerSession) SolutionRateLimit() (*big.Int, error) {
	return _Engine.Contract.SolutionRateLimit(&_Engine.CallOpts)
}

// Solutions is a free data retrieval call binding the contract method 0x75c70509.
//
// Solidity: function solutions(bytes32 ) view returns(address validator, uint64 blocktime, bool claimed, bytes cid)
func (_Engine *EngineCaller) Solutions(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Validator common.Address
	Blocktime uint64
	Claimed   bool
	Cid       []byte
}, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "solutions", arg0)

	outstruct := new(struct {
		Validator common.Address
		Blocktime uint64
		Claimed   bool
		Cid       []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Blocktime = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.Claimed = *abi.ConvertType(out[2], new(bool)).(*bool)
	outstruct.Cid = *abi.ConvertType(out[3], new([]byte)).(*[]byte)

	return *outstruct, err

}

// Solutions is a free data retrieval call binding the contract method 0x75c70509.
//
// Solidity: function solutions(bytes32 ) view returns(address validator, uint64 blocktime, bool claimed, bytes cid)
func (_Engine *EngineSession) Solutions(arg0 [32]byte) (struct {
	Validator common.Address
	Blocktime uint64
	Claimed   bool
	Cid       []byte
}, error) {
	return _Engine.Contract.Solutions(&_Engine.CallOpts, arg0)
}

// Solutions is a free data retrieval call binding the contract method 0x75c70509.
//
// Solidity: function solutions(bytes32 ) view returns(address validator, uint64 blocktime, bool claimed, bytes cid)
func (_Engine *EngineCallerSession) Solutions(arg0 [32]byte) (struct {
	Validator common.Address
	Blocktime uint64
	Claimed   bool
	Cid       []byte
}, error) {
	return _Engine.Contract.Solutions(&_Engine.CallOpts, arg0)
}

// SolutionsStake is a free data retrieval call binding the contract method 0xb4dc35b7.
//
// Solidity: function solutionsStake(bytes32 ) view returns(uint256)
func (_Engine *EngineCaller) SolutionsStake(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "solutionsStake", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SolutionsStake is a free data retrieval call binding the contract method 0xb4dc35b7.
//
// Solidity: function solutionsStake(bytes32 ) view returns(uint256)
func (_Engine *EngineSession) SolutionsStake(arg0 [32]byte) (*big.Int, error) {
	return _Engine.Contract.SolutionsStake(&_Engine.CallOpts, arg0)
}

// SolutionsStake is a free data retrieval call binding the contract method 0xb4dc35b7.
//
// Solidity: function solutionsStake(bytes32 ) view returns(uint256)
func (_Engine *EngineCallerSession) SolutionsStake(arg0 [32]byte) (*big.Int, error) {
	return _Engine.Contract.SolutionsStake(&_Engine.CallOpts, arg0)
}

// SolutionsStakeAmount is a free data retrieval call binding the contract method 0x9b975119.
//
// Solidity: function solutionsStakeAmount() view returns(uint256)
func (_Engine *EngineCaller) SolutionsStakeAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "solutionsStakeAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SolutionsStakeAmount is a free data retrieval call binding the contract method 0x9b975119.
//
// Solidity: function solutionsStakeAmount() view returns(uint256)
func (_Engine *EngineSession) SolutionsStakeAmount() (*big.Int, error) {
	return _Engine.Contract.SolutionsStakeAmount(&_Engine.CallOpts)
}

// SolutionsStakeAmount is a free data retrieval call binding the contract method 0x9b975119.
//
// Solidity: function solutionsStakeAmount() view returns(uint256)
func (_Engine *EngineCallerSession) SolutionsStakeAmount() (*big.Int, error) {
	return _Engine.Contract.SolutionsStakeAmount(&_Engine.CallOpts)
}

// StartBlockTime is a free data retrieval call binding the contract method 0x0c18d4ce.
//
// Solidity: function startBlockTime() view returns(uint64)
func (_Engine *EngineCaller) StartBlockTime(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "startBlockTime")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// StartBlockTime is a free data retrieval call binding the contract method 0x0c18d4ce.
//
// Solidity: function startBlockTime() view returns(uint64)
func (_Engine *EngineSession) StartBlockTime() (uint64, error) {
	return _Engine.Contract.StartBlockTime(&_Engine.CallOpts)
}

// StartBlockTime is a free data retrieval call binding the contract method 0x0c18d4ce.
//
// Solidity: function startBlockTime() view returns(uint64)
func (_Engine *EngineCallerSession) StartBlockTime() (uint64, error) {
	return _Engine.Contract.StartBlockTime(&_Engine.CallOpts)
}

// TargetTs is a free data retrieval call binding the contract method 0xcf596e45.
//
// Solidity: function targetTs(uint256 t) pure returns(uint256)
func (_Engine *EngineCaller) TargetTs(opts *bind.CallOpts, t *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "targetTs", t)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TargetTs is a free data retrieval call binding the contract method 0xcf596e45.
//
// Solidity: function targetTs(uint256 t) pure returns(uint256)
func (_Engine *EngineSession) TargetTs(t *big.Int) (*big.Int, error) {
	return _Engine.Contract.TargetTs(&_Engine.CallOpts, t)
}

// TargetTs is a free data retrieval call binding the contract method 0xcf596e45.
//
// Solidity: function targetTs(uint256 t) pure returns(uint256)
func (_Engine *EngineCallerSession) TargetTs(t *big.Int) (*big.Int, error) {
	return _Engine.Contract.TargetTs(&_Engine.CallOpts, t)
}

// TaskOwnerRewardPercentage is a free data retrieval call binding the contract method 0xc1f72723.
//
// Solidity: function taskOwnerRewardPercentage() view returns(uint256)
func (_Engine *EngineCaller) TaskOwnerRewardPercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "taskOwnerRewardPercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TaskOwnerRewardPercentage is a free data retrieval call binding the contract method 0xc1f72723.
//
// Solidity: function taskOwnerRewardPercentage() view returns(uint256)
func (_Engine *EngineSession) TaskOwnerRewardPercentage() (*big.Int, error) {
	return _Engine.Contract.TaskOwnerRewardPercentage(&_Engine.CallOpts)
}

// TaskOwnerRewardPercentage is a free data retrieval call binding the contract method 0xc1f72723.
//
// Solidity: function taskOwnerRewardPercentage() view returns(uint256)
func (_Engine *EngineCallerSession) TaskOwnerRewardPercentage() (*big.Int, error) {
	return _Engine.Contract.TaskOwnerRewardPercentage(&_Engine.CallOpts)
}

// Tasks is a free data retrieval call binding the contract method 0xe579f500.
//
// Solidity: function tasks(bytes32 ) view returns(bytes32 model, uint256 fee, address owner, uint64 blocktime, uint8 version, bytes cid)
func (_Engine *EngineCaller) Tasks(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Model     [32]byte
	Fee       *big.Int
	Owner     common.Address
	Blocktime uint64
	Version   uint8
	Cid       []byte
}, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "tasks", arg0)

	outstruct := new(struct {
		Model     [32]byte
		Fee       *big.Int
		Owner     common.Address
		Blocktime uint64
		Version   uint8
		Cid       []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Model = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Fee = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Owner = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Blocktime = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.Version = *abi.ConvertType(out[4], new(uint8)).(*uint8)
	outstruct.Cid = *abi.ConvertType(out[5], new([]byte)).(*[]byte)

	return *outstruct, err

}

// Tasks is a free data retrieval call binding the contract method 0xe579f500.
//
// Solidity: function tasks(bytes32 ) view returns(bytes32 model, uint256 fee, address owner, uint64 blocktime, uint8 version, bytes cid)
func (_Engine *EngineSession) Tasks(arg0 [32]byte) (struct {
	Model     [32]byte
	Fee       *big.Int
	Owner     common.Address
	Blocktime uint64
	Version   uint8
	Cid       []byte
}, error) {
	return _Engine.Contract.Tasks(&_Engine.CallOpts, arg0)
}

// Tasks is a free data retrieval call binding the contract method 0xe579f500.
//
// Solidity: function tasks(bytes32 ) view returns(bytes32 model, uint256 fee, address owner, uint64 blocktime, uint8 version, bytes cid)
func (_Engine *EngineCallerSession) Tasks(arg0 [32]byte) (struct {
	Model     [32]byte
	Fee       *big.Int
	Owner     common.Address
	Blocktime uint64
	Version   uint8
	Cid       []byte
}, error) {
	return _Engine.Contract.Tasks(&_Engine.CallOpts, arg0)
}

// TotalHeld is a free data retrieval call binding the contract method 0xf43cc773.
//
// Solidity: function totalHeld() view returns(uint256)
func (_Engine *EngineCaller) TotalHeld(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "totalHeld")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalHeld is a free data retrieval call binding the contract method 0xf43cc773.
//
// Solidity: function totalHeld() view returns(uint256)
func (_Engine *EngineSession) TotalHeld() (*big.Int, error) {
	return _Engine.Contract.TotalHeld(&_Engine.CallOpts)
}

// TotalHeld is a free data retrieval call binding the contract method 0xf43cc773.
//
// Solidity: function totalHeld() view returns(uint256)
func (_Engine *EngineCallerSession) TotalHeld() (*big.Int, error) {
	return _Engine.Contract.TotalHeld(&_Engine.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_Engine *EngineCaller) Treasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "treasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_Engine *EngineSession) Treasury() (common.Address, error) {
	return _Engine.Contract.Treasury(&_Engine.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_Engine *EngineCallerSession) Treasury() (common.Address, error) {
	return _Engine.Contract.Treasury(&_Engine.CallOpts)
}

// TreasuryRewardPercentage is a free data retrieval call binding the contract method 0xc31784be.
//
// Solidity: function treasuryRewardPercentage() view returns(uint256)
func (_Engine *EngineCaller) TreasuryRewardPercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "treasuryRewardPercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TreasuryRewardPercentage is a free data retrieval call binding the contract method 0xc31784be.
//
// Solidity: function treasuryRewardPercentage() view returns(uint256)
func (_Engine *EngineSession) TreasuryRewardPercentage() (*big.Int, error) {
	return _Engine.Contract.TreasuryRewardPercentage(&_Engine.CallOpts)
}

// TreasuryRewardPercentage is a free data retrieval call binding the contract method 0xc31784be.
//
// Solidity: function treasuryRewardPercentage() view returns(uint256)
func (_Engine *EngineCallerSession) TreasuryRewardPercentage() (*big.Int, error) {
	return _Engine.Contract.TreasuryRewardPercentage(&_Engine.CallOpts)
}

// ValidatorCanVote is a free data retrieval call binding the contract method 0x83657795.
//
// Solidity: function validatorCanVote(address addr_, bytes32 taskid_) view returns(uint256)
func (_Engine *EngineCaller) ValidatorCanVote(opts *bind.CallOpts, addr_ common.Address, taskid_ [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "validatorCanVote", addr_, taskid_)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorCanVote is a free data retrieval call binding the contract method 0x83657795.
//
// Solidity: function validatorCanVote(address addr_, bytes32 taskid_) view returns(uint256)
func (_Engine *EngineSession) ValidatorCanVote(addr_ common.Address, taskid_ [32]byte) (*big.Int, error) {
	return _Engine.Contract.ValidatorCanVote(&_Engine.CallOpts, addr_, taskid_)
}

// ValidatorCanVote is a free data retrieval call binding the contract method 0x83657795.
//
// Solidity: function validatorCanVote(address addr_, bytes32 taskid_) view returns(uint256)
func (_Engine *EngineCallerSession) ValidatorCanVote(addr_ common.Address, taskid_ [32]byte) (*big.Int, error) {
	return _Engine.Contract.ValidatorCanVote(&_Engine.CallOpts, addr_, taskid_)
}

// ValidatorMinimumPercentage is a free data retrieval call binding the contract method 0x96bb02c3.
//
// Solidity: function validatorMinimumPercentage() view returns(uint256)
func (_Engine *EngineCaller) ValidatorMinimumPercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "validatorMinimumPercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorMinimumPercentage is a free data retrieval call binding the contract method 0x96bb02c3.
//
// Solidity: function validatorMinimumPercentage() view returns(uint256)
func (_Engine *EngineSession) ValidatorMinimumPercentage() (*big.Int, error) {
	return _Engine.Contract.ValidatorMinimumPercentage(&_Engine.CallOpts)
}

// ValidatorMinimumPercentage is a free data retrieval call binding the contract method 0x96bb02c3.
//
// Solidity: function validatorMinimumPercentage() view returns(uint256)
func (_Engine *EngineCallerSession) ValidatorMinimumPercentage() (*big.Int, error) {
	return _Engine.Contract.ValidatorMinimumPercentage(&_Engine.CallOpts)
}

// ValidatorWithdrawPendingAmount is a free data retrieval call binding the contract method 0x1b75c43e.
//
// Solidity: function validatorWithdrawPendingAmount(address ) view returns(uint256)
func (_Engine *EngineCaller) ValidatorWithdrawPendingAmount(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "validatorWithdrawPendingAmount", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorWithdrawPendingAmount is a free data retrieval call binding the contract method 0x1b75c43e.
//
// Solidity: function validatorWithdrawPendingAmount(address ) view returns(uint256)
func (_Engine *EngineSession) ValidatorWithdrawPendingAmount(arg0 common.Address) (*big.Int, error) {
	return _Engine.Contract.ValidatorWithdrawPendingAmount(&_Engine.CallOpts, arg0)
}

// ValidatorWithdrawPendingAmount is a free data retrieval call binding the contract method 0x1b75c43e.
//
// Solidity: function validatorWithdrawPendingAmount(address ) view returns(uint256)
func (_Engine *EngineCallerSession) ValidatorWithdrawPendingAmount(arg0 common.Address) (*big.Int, error) {
	return _Engine.Contract.ValidatorWithdrawPendingAmount(&_Engine.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(uint256 staked, uint256 since, address addr)
func (_Engine *EngineCaller) Validators(opts *bind.CallOpts, arg0 common.Address) (struct {
	Staked *big.Int
	Since  *big.Int
	Addr   common.Address
}, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "validators", arg0)

	outstruct := new(struct {
		Staked *big.Int
		Since  *big.Int
		Addr   common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Staked = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Since = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Addr = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(uint256 staked, uint256 since, address addr)
func (_Engine *EngineSession) Validators(arg0 common.Address) (struct {
	Staked *big.Int
	Since  *big.Int
	Addr   common.Address
}, error) {
	return _Engine.Contract.Validators(&_Engine.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(uint256 staked, uint256 since, address addr)
func (_Engine *EngineCallerSession) Validators(arg0 common.Address) (struct {
	Staked *big.Int
	Since  *big.Int
	Addr   common.Address
}, error) {
	return _Engine.Contract.Validators(&_Engine.CallOpts, arg0)
}

// VeRewards is a free data retrieval call binding the contract method 0x0d468d95.
//
// Solidity: function veRewards() view returns(uint256)
func (_Engine *EngineCaller) VeRewards(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "veRewards")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VeRewards is a free data retrieval call binding the contract method 0x0d468d95.
//
// Solidity: function veRewards() view returns(uint256)
func (_Engine *EngineSession) VeRewards() (*big.Int, error) {
	return _Engine.Contract.VeRewards(&_Engine.CallOpts)
}

// VeRewards is a free data retrieval call binding the contract method 0x0d468d95.
//
// Solidity: function veRewards() view returns(uint256)
func (_Engine *EngineCallerSession) VeRewards() (*big.Int, error) {
	return _Engine.Contract.VeRewards(&_Engine.CallOpts)
}

// VeStaking is a free data retrieval call binding the contract method 0x82b5077f.
//
// Solidity: function veStaking() view returns(address)
func (_Engine *EngineCaller) VeStaking(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "veStaking")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VeStaking is a free data retrieval call binding the contract method 0x82b5077f.
//
// Solidity: function veStaking() view returns(address)
func (_Engine *EngineSession) VeStaking() (common.Address, error) {
	return _Engine.Contract.VeStaking(&_Engine.CallOpts)
}

// VeStaking is a free data retrieval call binding the contract method 0x82b5077f.
//
// Solidity: function veStaking() view returns(address)
func (_Engine *EngineCallerSession) VeStaking() (common.Address, error) {
	return _Engine.Contract.VeStaking(&_Engine.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_Engine *EngineCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_Engine *EngineSession) Version() (*big.Int, error) {
	return _Engine.Contract.Version(&_Engine.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_Engine *EngineCallerSession) Version() (*big.Int, error) {
	return _Engine.Contract.Version(&_Engine.CallOpts)
}

// Voter is a free data retrieval call binding the contract method 0x46c96aac.
//
// Solidity: function voter() view returns(address)
func (_Engine *EngineCaller) Voter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "voter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Voter is a free data retrieval call binding the contract method 0x46c96aac.
//
// Solidity: function voter() view returns(address)
func (_Engine *EngineSession) Voter() (common.Address, error) {
	return _Engine.Contract.Voter(&_Engine.CallOpts)
}

// Voter is a free data retrieval call binding the contract method 0x46c96aac.
//
// Solidity: function voter() view returns(address)
func (_Engine *EngineCallerSession) Voter() (common.Address, error) {
	return _Engine.Contract.Voter(&_Engine.CallOpts)
}

// VotingPeriodEnded is a free data retrieval call binding the contract method 0x05d1bc26.
//
// Solidity: function votingPeriodEnded(bytes32 taskid_) view returns(bool)
func (_Engine *EngineCaller) VotingPeriodEnded(opts *bind.CallOpts, taskid_ [32]byte) (bool, error) {
	var out []interface{}
	err := _Engine.contract.Call(opts, &out, "votingPeriodEnded", taskid_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VotingPeriodEnded is a free data retrieval call binding the contract method 0x05d1bc26.
//
// Solidity: function votingPeriodEnded(bytes32 taskid_) view returns(bool)
func (_Engine *EngineSession) VotingPeriodEnded(taskid_ [32]byte) (bool, error) {
	return _Engine.Contract.VotingPeriodEnded(&_Engine.CallOpts, taskid_)
}

// VotingPeriodEnded is a free data retrieval call binding the contract method 0x05d1bc26.
//
// Solidity: function votingPeriodEnded(bytes32 taskid_) view returns(bool)
func (_Engine *EngineCallerSession) VotingPeriodEnded(taskid_ [32]byte) (bool, error) {
	return _Engine.Contract.VotingPeriodEnded(&_Engine.CallOpts, taskid_)
}

// BulkSubmitSolution is a paid mutator transaction binding the contract method 0x65d445fb.
//
// Solidity: function bulkSubmitSolution(bytes32[] taskids_, bytes[] cids_) returns()
func (_Engine *EngineTransactor) BulkSubmitSolution(opts *bind.TransactOpts, taskids_ [][32]byte, cids_ [][]byte) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "bulkSubmitSolution", taskids_, cids_)
}

// BulkSubmitSolution is a paid mutator transaction binding the contract method 0x65d445fb.
//
// Solidity: function bulkSubmitSolution(bytes32[] taskids_, bytes[] cids_) returns()
func (_Engine *EngineSession) BulkSubmitSolution(taskids_ [][32]byte, cids_ [][]byte) (*types.Transaction, error) {
	return _Engine.Contract.BulkSubmitSolution(&_Engine.TransactOpts, taskids_, cids_)
}

// BulkSubmitSolution is a paid mutator transaction binding the contract method 0x65d445fb.
//
// Solidity: function bulkSubmitSolution(bytes32[] taskids_, bytes[] cids_) returns()
func (_Engine *EngineTransactorSession) BulkSubmitSolution(taskids_ [][32]byte, cids_ [][]byte) (*types.Transaction, error) {
	return _Engine.Contract.BulkSubmitSolution(&_Engine.TransactOpts, taskids_, cids_)
}

// BulkSubmitTask is a paid mutator transaction binding the contract method 0x08afe0eb.
//
// Solidity: function bulkSubmitTask(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_, uint256 n_) returns()
func (_Engine *EngineTransactor) BulkSubmitTask(opts *bind.TransactOpts, version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte, n_ *big.Int) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "bulkSubmitTask", version_, owner_, model_, fee_, input_, n_)
}

// BulkSubmitTask is a paid mutator transaction binding the contract method 0x08afe0eb.
//
// Solidity: function bulkSubmitTask(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_, uint256 n_) returns()
func (_Engine *EngineSession) BulkSubmitTask(version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte, n_ *big.Int) (*types.Transaction, error) {
	return _Engine.Contract.BulkSubmitTask(&_Engine.TransactOpts, version_, owner_, model_, fee_, input_, n_)
}

// BulkSubmitTask is a paid mutator transaction binding the contract method 0x08afe0eb.
//
// Solidity: function bulkSubmitTask(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_, uint256 n_) returns()
func (_Engine *EngineTransactorSession) BulkSubmitTask(version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte, n_ *big.Int) (*types.Transaction, error) {
	return _Engine.Contract.BulkSubmitTask(&_Engine.TransactOpts, version_, owner_, model_, fee_, input_, n_)
}

// CancelValidatorWithdraw is a paid mutator transaction binding the contract method 0xcbd2422d.
//
// Solidity: function cancelValidatorWithdraw(uint256 count_) returns()
func (_Engine *EngineTransactor) CancelValidatorWithdraw(opts *bind.TransactOpts, count_ *big.Int) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "cancelValidatorWithdraw", count_)
}

// CancelValidatorWithdraw is a paid mutator transaction binding the contract method 0xcbd2422d.
//
// Solidity: function cancelValidatorWithdraw(uint256 count_) returns()
func (_Engine *EngineSession) CancelValidatorWithdraw(count_ *big.Int) (*types.Transaction, error) {
	return _Engine.Contract.CancelValidatorWithdraw(&_Engine.TransactOpts, count_)
}

// CancelValidatorWithdraw is a paid mutator transaction binding the contract method 0xcbd2422d.
//
// Solidity: function cancelValidatorWithdraw(uint256 count_) returns()
func (_Engine *EngineTransactorSession) CancelValidatorWithdraw(count_ *big.Int) (*types.Transaction, error) {
	return _Engine.Contract.CancelValidatorWithdraw(&_Engine.TransactOpts, count_)
}

// ClaimSolution is a paid mutator transaction binding the contract method 0x77286d17.
//
// Solidity: function claimSolution(bytes32 taskid_) returns()
func (_Engine *EngineTransactor) ClaimSolution(opts *bind.TransactOpts, taskid_ [32]byte) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "claimSolution", taskid_)
}

// ClaimSolution is a paid mutator transaction binding the contract method 0x77286d17.
//
// Solidity: function claimSolution(bytes32 taskid_) returns()
func (_Engine *EngineSession) ClaimSolution(taskid_ [32]byte) (*types.Transaction, error) {
	return _Engine.Contract.ClaimSolution(&_Engine.TransactOpts, taskid_)
}

// ClaimSolution is a paid mutator transaction binding the contract method 0x77286d17.
//
// Solidity: function claimSolution(bytes32 taskid_) returns()
func (_Engine *EngineTransactorSession) ClaimSolution(taskid_ [32]byte) (*types.Transaction, error) {
	return _Engine.Contract.ClaimSolution(&_Engine.TransactOpts, taskid_)
}

// ContestationVoteFinish is a paid mutator transaction binding the contract method 0x8b4d7b35.
//
// Solidity: function contestationVoteFinish(bytes32 taskid_, uint32 amnt_) returns()
func (_Engine *EngineTransactor) ContestationVoteFinish(opts *bind.TransactOpts, taskid_ [32]byte, amnt_ uint32) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "contestationVoteFinish", taskid_, amnt_)
}

// ContestationVoteFinish is a paid mutator transaction binding the contract method 0x8b4d7b35.
//
// Solidity: function contestationVoteFinish(bytes32 taskid_, uint32 amnt_) returns()
func (_Engine *EngineSession) ContestationVoteFinish(taskid_ [32]byte, amnt_ uint32) (*types.Transaction, error) {
	return _Engine.Contract.ContestationVoteFinish(&_Engine.TransactOpts, taskid_, amnt_)
}

// ContestationVoteFinish is a paid mutator transaction binding the contract method 0x8b4d7b35.
//
// Solidity: function contestationVoteFinish(bytes32 taskid_, uint32 amnt_) returns()
func (_Engine *EngineTransactorSession) ContestationVoteFinish(taskid_ [32]byte, amnt_ uint32) (*types.Transaction, error) {
	return _Engine.Contract.ContestationVoteFinish(&_Engine.TransactOpts, taskid_, amnt_)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Engine *EngineTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Engine *EngineSession) Initialize() (*types.Transaction, error) {
	return _Engine.Contract.Initialize(&_Engine.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Engine *EngineTransactorSession) Initialize() (*types.Transaction, error) {
	return _Engine.Contract.Initialize(&_Engine.TransactOpts)
}

// InitiateValidatorWithdraw is a paid mutator transaction binding the contract method 0x0a985737.
//
// Solidity: function initiateValidatorWithdraw(uint256 amount_) returns(uint256)
func (_Engine *EngineTransactor) InitiateValidatorWithdraw(opts *bind.TransactOpts, amount_ *big.Int) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "initiateValidatorWithdraw", amount_)
}

// InitiateValidatorWithdraw is a paid mutator transaction binding the contract method 0x0a985737.
//
// Solidity: function initiateValidatorWithdraw(uint256 amount_) returns(uint256)
func (_Engine *EngineSession) InitiateValidatorWithdraw(amount_ *big.Int) (*types.Transaction, error) {
	return _Engine.Contract.InitiateValidatorWithdraw(&_Engine.TransactOpts, amount_)
}

// InitiateValidatorWithdraw is a paid mutator transaction binding the contract method 0x0a985737.
//
// Solidity: function initiateValidatorWithdraw(uint256 amount_) returns(uint256)
func (_Engine *EngineTransactorSession) InitiateValidatorWithdraw(amount_ *big.Int) (*types.Transaction, error) {
	return _Engine.Contract.InitiateValidatorWithdraw(&_Engine.TransactOpts, amount_)
}

// RegisterModel is a paid mutator transaction binding the contract method 0x4ff03efa.
//
// Solidity: function registerModel(address addr_, uint256 fee_, bytes template_) returns(bytes32)
func (_Engine *EngineTransactor) RegisterModel(opts *bind.TransactOpts, addr_ common.Address, fee_ *big.Int, template_ []byte) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "registerModel", addr_, fee_, template_)
}

// RegisterModel is a paid mutator transaction binding the contract method 0x4ff03efa.
//
// Solidity: function registerModel(address addr_, uint256 fee_, bytes template_) returns(bytes32)
func (_Engine *EngineSession) RegisterModel(addr_ common.Address, fee_ *big.Int, template_ []byte) (*types.Transaction, error) {
	return _Engine.Contract.RegisterModel(&_Engine.TransactOpts, addr_, fee_, template_)
}

// RegisterModel is a paid mutator transaction binding the contract method 0x4ff03efa.
//
// Solidity: function registerModel(address addr_, uint256 fee_, bytes template_) returns(bytes32)
func (_Engine *EngineTransactorSession) RegisterModel(addr_ common.Address, fee_ *big.Int, template_ []byte) (*types.Transaction, error) {
	return _Engine.Contract.RegisterModel(&_Engine.TransactOpts, addr_, fee_, template_)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Engine *EngineTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Engine *EngineSession) RenounceOwnership() (*types.Transaction, error) {
	return _Engine.Contract.RenounceOwnership(&_Engine.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Engine *EngineTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Engine.Contract.RenounceOwnership(&_Engine.TransactOpts)
}

// SetModelAddr is a paid mutator transaction binding the contract method 0xdafb4322.
//
// Solidity: function setModelAddr(bytes32 model_, address addr_) returns()
func (_Engine *EngineTransactor) SetModelAddr(opts *bind.TransactOpts, model_ [32]byte, addr_ common.Address) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "setModelAddr", model_, addr_)
}

// SetModelAddr is a paid mutator transaction binding the contract method 0xdafb4322.
//
// Solidity: function setModelAddr(bytes32 model_, address addr_) returns()
func (_Engine *EngineSession) SetModelAddr(model_ [32]byte, addr_ common.Address) (*types.Transaction, error) {
	return _Engine.Contract.SetModelAddr(&_Engine.TransactOpts, model_, addr_)
}

// SetModelAddr is a paid mutator transaction binding the contract method 0xdafb4322.
//
// Solidity: function setModelAddr(bytes32 model_, address addr_) returns()
func (_Engine *EngineTransactorSession) SetModelAddr(model_ [32]byte, addr_ common.Address) (*types.Transaction, error) {
	return _Engine.Contract.SetModelAddr(&_Engine.TransactOpts, model_, addr_)
}

// SetModelFee is a paid mutator transaction binding the contract method 0xf029ac3e.
//
// Solidity: function setModelFee(bytes32 model_, uint256 fee_) returns()
func (_Engine *EngineTransactor) SetModelFee(opts *bind.TransactOpts, model_ [32]byte, fee_ *big.Int) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "setModelFee", model_, fee_)
}

// SetModelFee is a paid mutator transaction binding the contract method 0xf029ac3e.
//
// Solidity: function setModelFee(bytes32 model_, uint256 fee_) returns()
func (_Engine *EngineSession) SetModelFee(model_ [32]byte, fee_ *big.Int) (*types.Transaction, error) {
	return _Engine.Contract.SetModelFee(&_Engine.TransactOpts, model_, fee_)
}

// SetModelFee is a paid mutator transaction binding the contract method 0xf029ac3e.
//
// Solidity: function setModelFee(bytes32 model_, uint256 fee_) returns()
func (_Engine *EngineTransactorSession) SetModelFee(model_ [32]byte, fee_ *big.Int) (*types.Transaction, error) {
	return _Engine.Contract.SetModelFee(&_Engine.TransactOpts, model_, fee_)
}

// SetPaused is a paid mutator transaction binding the contract method 0x16c38b3c.
//
// Solidity: function setPaused(bool paused_) returns()
func (_Engine *EngineTransactor) SetPaused(opts *bind.TransactOpts, paused_ bool) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "setPaused", paused_)
}

// SetPaused is a paid mutator transaction binding the contract method 0x16c38b3c.
//
// Solidity: function setPaused(bool paused_) returns()
func (_Engine *EngineSession) SetPaused(paused_ bool) (*types.Transaction, error) {
	return _Engine.Contract.SetPaused(&_Engine.TransactOpts, paused_)
}

// SetPaused is a paid mutator transaction binding the contract method 0x16c38b3c.
//
// Solidity: function setPaused(bool paused_) returns()
func (_Engine *EngineTransactorSession) SetPaused(paused_ bool) (*types.Transaction, error) {
	return _Engine.Contract.SetPaused(&_Engine.TransactOpts, paused_)
}

// SetSolutionMineableRate is a paid mutator transaction binding the contract method 0x93f1f8ac.
//
// Solidity: function setSolutionMineableRate(bytes32 model_, uint256 rate_) returns()
func (_Engine *EngineTransactor) SetSolutionMineableRate(opts *bind.TransactOpts, model_ [32]byte, rate_ *big.Int) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "setSolutionMineableRate", model_, rate_)
}

// SetSolutionMineableRate is a paid mutator transaction binding the contract method 0x93f1f8ac.
//
// Solidity: function setSolutionMineableRate(bytes32 model_, uint256 rate_) returns()
func (_Engine *EngineSession) SetSolutionMineableRate(model_ [32]byte, rate_ *big.Int) (*types.Transaction, error) {
	return _Engine.Contract.SetSolutionMineableRate(&_Engine.TransactOpts, model_, rate_)
}

// SetSolutionMineableRate is a paid mutator transaction binding the contract method 0x93f1f8ac.
//
// Solidity: function setSolutionMineableRate(bytes32 model_, uint256 rate_) returns()
func (_Engine *EngineTransactorSession) SetSolutionMineableRate(model_ [32]byte, rate_ *big.Int) (*types.Transaction, error) {
	return _Engine.Contract.SetSolutionMineableRate(&_Engine.TransactOpts, model_, rate_)
}

// SetStartBlockTime is a paid mutator transaction binding the contract method 0xa8f837f3.
//
// Solidity: function setStartBlockTime(uint64 startBlockTime_) returns()
func (_Engine *EngineTransactor) SetStartBlockTime(opts *bind.TransactOpts, startBlockTime_ uint64) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "setStartBlockTime", startBlockTime_)
}

// SetStartBlockTime is a paid mutator transaction binding the contract method 0xa8f837f3.
//
// Solidity: function setStartBlockTime(uint64 startBlockTime_) returns()
func (_Engine *EngineSession) SetStartBlockTime(startBlockTime_ uint64) (*types.Transaction, error) {
	return _Engine.Contract.SetStartBlockTime(&_Engine.TransactOpts, startBlockTime_)
}

// SetStartBlockTime is a paid mutator transaction binding the contract method 0xa8f837f3.
//
// Solidity: function setStartBlockTime(uint64 startBlockTime_) returns()
func (_Engine *EngineTransactorSession) SetStartBlockTime(startBlockTime_ uint64) (*types.Transaction, error) {
	return _Engine.Contract.SetStartBlockTime(&_Engine.TransactOpts, startBlockTime_)
}

// SetVeStaking is a paid mutator transaction binding the contract method 0x2943a490.
//
// Solidity: function setVeStaking(address veStaking_) returns()
func (_Engine *EngineTransactor) SetVeStaking(opts *bind.TransactOpts, veStaking_ common.Address) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "setVeStaking", veStaking_)
}

// SetVeStaking is a paid mutator transaction binding the contract method 0x2943a490.
//
// Solidity: function setVeStaking(address veStaking_) returns()
func (_Engine *EngineSession) SetVeStaking(veStaking_ common.Address) (*types.Transaction, error) {
	return _Engine.Contract.SetVeStaking(&_Engine.TransactOpts, veStaking_)
}

// SetVeStaking is a paid mutator transaction binding the contract method 0x2943a490.
//
// Solidity: function setVeStaking(address veStaking_) returns()
func (_Engine *EngineTransactorSession) SetVeStaking(veStaking_ common.Address) (*types.Transaction, error) {
	return _Engine.Contract.SetVeStaking(&_Engine.TransactOpts, veStaking_)
}

// SetVersion is a paid mutator transaction binding the contract method 0x408def1e.
//
// Solidity: function setVersion(uint256 version_) returns()
func (_Engine *EngineTransactor) SetVersion(opts *bind.TransactOpts, version_ *big.Int) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "setVersion", version_)
}

// SetVersion is a paid mutator transaction binding the contract method 0x408def1e.
//
// Solidity: function setVersion(uint256 version_) returns()
func (_Engine *EngineSession) SetVersion(version_ *big.Int) (*types.Transaction, error) {
	return _Engine.Contract.SetVersion(&_Engine.TransactOpts, version_)
}

// SetVersion is a paid mutator transaction binding the contract method 0x408def1e.
//
// Solidity: function setVersion(uint256 version_) returns()
func (_Engine *EngineTransactorSession) SetVersion(version_ *big.Int) (*types.Transaction, error) {
	return _Engine.Contract.SetVersion(&_Engine.TransactOpts, version_)
}

// SetVoter is a paid mutator transaction binding the contract method 0x4bc2a657.
//
// Solidity: function setVoter(address voter_) returns()
func (_Engine *EngineTransactor) SetVoter(opts *bind.TransactOpts, voter_ common.Address) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "setVoter", voter_)
}

// SetVoter is a paid mutator transaction binding the contract method 0x4bc2a657.
//
// Solidity: function setVoter(address voter_) returns()
func (_Engine *EngineSession) SetVoter(voter_ common.Address) (*types.Transaction, error) {
	return _Engine.Contract.SetVoter(&_Engine.TransactOpts, voter_)
}

// SetVoter is a paid mutator transaction binding the contract method 0x4bc2a657.
//
// Solidity: function setVoter(address voter_) returns()
func (_Engine *EngineTransactorSession) SetVoter(voter_ common.Address) (*types.Transaction, error) {
	return _Engine.Contract.SetVoter(&_Engine.TransactOpts, voter_)
}

// SignalCommitment is a paid mutator transaction binding the contract method 0x506ea7de.
//
// Solidity: function signalCommitment(bytes32 commitment_) returns()
func (_Engine *EngineTransactor) SignalCommitment(opts *bind.TransactOpts, commitment_ [32]byte) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "signalCommitment", commitment_)
}

// SignalCommitment is a paid mutator transaction binding the contract method 0x506ea7de.
//
// Solidity: function signalCommitment(bytes32 commitment_) returns()
func (_Engine *EngineSession) SignalCommitment(commitment_ [32]byte) (*types.Transaction, error) {
	return _Engine.Contract.SignalCommitment(&_Engine.TransactOpts, commitment_)
}

// SignalCommitment is a paid mutator transaction binding the contract method 0x506ea7de.
//
// Solidity: function signalCommitment(bytes32 commitment_) returns()
func (_Engine *EngineTransactorSession) SignalCommitment(commitment_ [32]byte) (*types.Transaction, error) {
	return _Engine.Contract.SignalCommitment(&_Engine.TransactOpts, commitment_)
}

// SubmitContestation is a paid mutator transaction binding the contract method 0x671f8152.
//
// Solidity: function submitContestation(bytes32 taskid_) returns()
func (_Engine *EngineTransactor) SubmitContestation(opts *bind.TransactOpts, taskid_ [32]byte) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "submitContestation", taskid_)
}

// SubmitContestation is a paid mutator transaction binding the contract method 0x671f8152.
//
// Solidity: function submitContestation(bytes32 taskid_) returns()
func (_Engine *EngineSession) SubmitContestation(taskid_ [32]byte) (*types.Transaction, error) {
	return _Engine.Contract.SubmitContestation(&_Engine.TransactOpts, taskid_)
}

// SubmitContestation is a paid mutator transaction binding the contract method 0x671f8152.
//
// Solidity: function submitContestation(bytes32 taskid_) returns()
func (_Engine *EngineTransactorSession) SubmitContestation(taskid_ [32]byte) (*types.Transaction, error) {
	return _Engine.Contract.SubmitContestation(&_Engine.TransactOpts, taskid_)
}

// SubmitSolution is a paid mutator transaction binding the contract method 0x56914caf.
//
// Solidity: function submitSolution(bytes32 taskid_, bytes cid_) returns()
func (_Engine *EngineTransactor) SubmitSolution(opts *bind.TransactOpts, taskid_ [32]byte, cid_ []byte) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "submitSolution", taskid_, cid_)
}

// SubmitSolution is a paid mutator transaction binding the contract method 0x56914caf.
//
// Solidity: function submitSolution(bytes32 taskid_, bytes cid_) returns()
func (_Engine *EngineSession) SubmitSolution(taskid_ [32]byte, cid_ []byte) (*types.Transaction, error) {
	return _Engine.Contract.SubmitSolution(&_Engine.TransactOpts, taskid_, cid_)
}

// SubmitSolution is a paid mutator transaction binding the contract method 0x56914caf.
//
// Solidity: function submitSolution(bytes32 taskid_, bytes cid_) returns()
func (_Engine *EngineTransactorSession) SubmitSolution(taskid_ [32]byte, cid_ []byte) (*types.Transaction, error) {
	return _Engine.Contract.SubmitSolution(&_Engine.TransactOpts, taskid_, cid_)
}

// SubmitTask is a paid mutator transaction binding the contract method 0x08745dd1.
//
// Solidity: function submitTask(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_) returns()
func (_Engine *EngineTransactor) SubmitTask(opts *bind.TransactOpts, version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "submitTask", version_, owner_, model_, fee_, input_)
}

// SubmitTask is a paid mutator transaction binding the contract method 0x08745dd1.
//
// Solidity: function submitTask(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_) returns()
func (_Engine *EngineSession) SubmitTask(version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte) (*types.Transaction, error) {
	return _Engine.Contract.SubmitTask(&_Engine.TransactOpts, version_, owner_, model_, fee_, input_)
}

// SubmitTask is a paid mutator transaction binding the contract method 0x08745dd1.
//
// Solidity: function submitTask(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_) returns()
func (_Engine *EngineTransactorSession) SubmitTask(version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte) (*types.Transaction, error) {
	return _Engine.Contract.SubmitTask(&_Engine.TransactOpts, version_, owner_, model_, fee_, input_)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to_) returns()
func (_Engine *EngineTransactor) TransferOwnership(opts *bind.TransactOpts, to_ common.Address) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "transferOwnership", to_)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to_) returns()
func (_Engine *EngineSession) TransferOwnership(to_ common.Address) (*types.Transaction, error) {
	return _Engine.Contract.TransferOwnership(&_Engine.TransactOpts, to_)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to_) returns()
func (_Engine *EngineTransactorSession) TransferOwnership(to_ common.Address) (*types.Transaction, error) {
	return _Engine.Contract.TransferOwnership(&_Engine.TransactOpts, to_)
}

// TransferPauser is a paid mutator transaction binding the contract method 0x4421ea21.
//
// Solidity: function transferPauser(address to_) returns()
func (_Engine *EngineTransactor) TransferPauser(opts *bind.TransactOpts, to_ common.Address) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "transferPauser", to_)
}

// TransferPauser is a paid mutator transaction binding the contract method 0x4421ea21.
//
// Solidity: function transferPauser(address to_) returns()
func (_Engine *EngineSession) TransferPauser(to_ common.Address) (*types.Transaction, error) {
	return _Engine.Contract.TransferPauser(&_Engine.TransactOpts, to_)
}

// TransferPauser is a paid mutator transaction binding the contract method 0x4421ea21.
//
// Solidity: function transferPauser(address to_) returns()
func (_Engine *EngineTransactorSession) TransferPauser(to_ common.Address) (*types.Transaction, error) {
	return _Engine.Contract.TransferPauser(&_Engine.TransactOpts, to_)
}

// TransferTreasury is a paid mutator transaction binding the contract method 0xd8a6021c.
//
// Solidity: function transferTreasury(address to_) returns()
func (_Engine *EngineTransactor) TransferTreasury(opts *bind.TransactOpts, to_ common.Address) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "transferTreasury", to_)
}

// TransferTreasury is a paid mutator transaction binding the contract method 0xd8a6021c.
//
// Solidity: function transferTreasury(address to_) returns()
func (_Engine *EngineSession) TransferTreasury(to_ common.Address) (*types.Transaction, error) {
	return _Engine.Contract.TransferTreasury(&_Engine.TransactOpts, to_)
}

// TransferTreasury is a paid mutator transaction binding the contract method 0xd8a6021c.
//
// Solidity: function transferTreasury(address to_) returns()
func (_Engine *EngineTransactorSession) TransferTreasury(to_ common.Address) (*types.Transaction, error) {
	return _Engine.Contract.TransferTreasury(&_Engine.TransactOpts, to_)
}

// ValidatorDeposit is a paid mutator transaction binding the contract method 0x93a090ec.
//
// Solidity: function validatorDeposit(address validator_, uint256 amount_) returns()
func (_Engine *EngineTransactor) ValidatorDeposit(opts *bind.TransactOpts, validator_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "validatorDeposit", validator_, amount_)
}

// ValidatorDeposit is a paid mutator transaction binding the contract method 0x93a090ec.
//
// Solidity: function validatorDeposit(address validator_, uint256 amount_) returns()
func (_Engine *EngineSession) ValidatorDeposit(validator_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _Engine.Contract.ValidatorDeposit(&_Engine.TransactOpts, validator_, amount_)
}

// ValidatorDeposit is a paid mutator transaction binding the contract method 0x93a090ec.
//
// Solidity: function validatorDeposit(address validator_, uint256 amount_) returns()
func (_Engine *EngineTransactorSession) ValidatorDeposit(validator_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _Engine.Contract.ValidatorDeposit(&_Engine.TransactOpts, validator_, amount_)
}

// ValidatorWithdraw is a paid mutator transaction binding the contract method 0x763253bb.
//
// Solidity: function validatorWithdraw(uint256 count_, address to_) returns()
func (_Engine *EngineTransactor) ValidatorWithdraw(opts *bind.TransactOpts, count_ *big.Int, to_ common.Address) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "validatorWithdraw", count_, to_)
}

// ValidatorWithdraw is a paid mutator transaction binding the contract method 0x763253bb.
//
// Solidity: function validatorWithdraw(uint256 count_, address to_) returns()
func (_Engine *EngineSession) ValidatorWithdraw(count_ *big.Int, to_ common.Address) (*types.Transaction, error) {
	return _Engine.Contract.ValidatorWithdraw(&_Engine.TransactOpts, count_, to_)
}

// ValidatorWithdraw is a paid mutator transaction binding the contract method 0x763253bb.
//
// Solidity: function validatorWithdraw(uint256 count_, address to_) returns()
func (_Engine *EngineTransactorSession) ValidatorWithdraw(count_ *big.Int, to_ common.Address) (*types.Transaction, error) {
	return _Engine.Contract.ValidatorWithdraw(&_Engine.TransactOpts, count_, to_)
}

// VoteOnContestation is a paid mutator transaction binding the contract method 0x1825c20e.
//
// Solidity: function voteOnContestation(bytes32 taskid_, bool yea_) returns()
func (_Engine *EngineTransactor) VoteOnContestation(opts *bind.TransactOpts, taskid_ [32]byte, yea_ bool) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "voteOnContestation", taskid_, yea_)
}

// VoteOnContestation is a paid mutator transaction binding the contract method 0x1825c20e.
//
// Solidity: function voteOnContestation(bytes32 taskid_, bool yea_) returns()
func (_Engine *EngineSession) VoteOnContestation(taskid_ [32]byte, yea_ bool) (*types.Transaction, error) {
	return _Engine.Contract.VoteOnContestation(&_Engine.TransactOpts, taskid_, yea_)
}

// VoteOnContestation is a paid mutator transaction binding the contract method 0x1825c20e.
//
// Solidity: function voteOnContestation(bytes32 taskid_, bool yea_) returns()
func (_Engine *EngineTransactorSession) VoteOnContestation(taskid_ [32]byte, yea_ bool) (*types.Transaction, error) {
	return _Engine.Contract.VoteOnContestation(&_Engine.TransactOpts, taskid_, yea_)
}

// WithdrawAccruedFees is a paid mutator transaction binding the contract method 0xada82c7d.
//
// Solidity: function withdrawAccruedFees() returns()
func (_Engine *EngineTransactor) WithdrawAccruedFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Engine.contract.Transact(opts, "withdrawAccruedFees")
}

// WithdrawAccruedFees is a paid mutator transaction binding the contract method 0xada82c7d.
//
// Solidity: function withdrawAccruedFees() returns()
func (_Engine *EngineSession) WithdrawAccruedFees() (*types.Transaction, error) {
	return _Engine.Contract.WithdrawAccruedFees(&_Engine.TransactOpts)
}

// WithdrawAccruedFees is a paid mutator transaction binding the contract method 0xada82c7d.
//
// Solidity: function withdrawAccruedFees() returns()
func (_Engine *EngineTransactorSession) WithdrawAccruedFees() (*types.Transaction, error) {
	return _Engine.Contract.WithdrawAccruedFees(&_Engine.TransactOpts)
}

// EngineContestationSubmittedIterator is returned from FilterContestationSubmitted and is used to iterate over the raw logs and unpacked data for ContestationSubmitted events raised by the Engine contract.
type EngineContestationSubmittedIterator struct {
	Event *EngineContestationSubmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineContestationSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineContestationSubmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineContestationSubmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineContestationSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineContestationSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineContestationSubmitted represents a ContestationSubmitted event raised by the Engine contract.
type EngineContestationSubmitted struct {
	Addr common.Address
	Task [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterContestationSubmitted is a free log retrieval operation binding the contract event 0x6958c989e915d3e41a35076e3c480363910055408055ad86ae1ee13d41c40640.
//
// Solidity: event ContestationSubmitted(address indexed addr, bytes32 indexed task)
func (_Engine *EngineFilterer) FilterContestationSubmitted(opts *bind.FilterOpts, addr []common.Address, task [][32]byte) (*EngineContestationSubmittedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var taskRule []interface{}
	for _, taskItem := range task {
		taskRule = append(taskRule, taskItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "ContestationSubmitted", addrRule, taskRule)
	if err != nil {
		return nil, err
	}
	return &EngineContestationSubmittedIterator{contract: _Engine.contract, event: "ContestationSubmitted", logs: logs, sub: sub}, nil
}

// WatchContestationSubmitted is a free log subscription operation binding the contract event 0x6958c989e915d3e41a35076e3c480363910055408055ad86ae1ee13d41c40640.
//
// Solidity: event ContestationSubmitted(address indexed addr, bytes32 indexed task)
func (_Engine *EngineFilterer) WatchContestationSubmitted(opts *bind.WatchOpts, sink chan<- *EngineContestationSubmitted, addr []common.Address, task [][32]byte) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var taskRule []interface{}
	for _, taskItem := range task {
		taskRule = append(taskRule, taskItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "ContestationSubmitted", addrRule, taskRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineContestationSubmitted)
				if err := _Engine.contract.UnpackLog(event, "ContestationSubmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseContestationSubmitted is a log parse operation binding the contract event 0x6958c989e915d3e41a35076e3c480363910055408055ad86ae1ee13d41c40640.
//
// Solidity: event ContestationSubmitted(address indexed addr, bytes32 indexed task)
func (_Engine *EngineFilterer) ParseContestationSubmitted(log types.Log) (*EngineContestationSubmitted, error) {
	event := new(EngineContestationSubmitted)
	if err := _Engine.contract.UnpackLog(event, "ContestationSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineContestationVoteIterator is returned from FilterContestationVote and is used to iterate over the raw logs and unpacked data for ContestationVote events raised by the Engine contract.
type EngineContestationVoteIterator struct {
	Event *EngineContestationVote // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineContestationVoteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineContestationVote)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineContestationVote)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineContestationVoteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineContestationVoteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineContestationVote represents a ContestationVote event raised by the Engine contract.
type EngineContestationVote struct {
	Addr common.Address
	Task [32]byte
	Yea  bool
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterContestationVote is a free log retrieval operation binding the contract event 0x1aa9e4be46e24e1f2e7eeb1613c01629213cd42965d2716e18531b63e552e411.
//
// Solidity: event ContestationVote(address indexed addr, bytes32 indexed task, bool yea)
func (_Engine *EngineFilterer) FilterContestationVote(opts *bind.FilterOpts, addr []common.Address, task [][32]byte) (*EngineContestationVoteIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var taskRule []interface{}
	for _, taskItem := range task {
		taskRule = append(taskRule, taskItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "ContestationVote", addrRule, taskRule)
	if err != nil {
		return nil, err
	}
	return &EngineContestationVoteIterator{contract: _Engine.contract, event: "ContestationVote", logs: logs, sub: sub}, nil
}

// WatchContestationVote is a free log subscription operation binding the contract event 0x1aa9e4be46e24e1f2e7eeb1613c01629213cd42965d2716e18531b63e552e411.
//
// Solidity: event ContestationVote(address indexed addr, bytes32 indexed task, bool yea)
func (_Engine *EngineFilterer) WatchContestationVote(opts *bind.WatchOpts, sink chan<- *EngineContestationVote, addr []common.Address, task [][32]byte) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var taskRule []interface{}
	for _, taskItem := range task {
		taskRule = append(taskRule, taskItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "ContestationVote", addrRule, taskRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineContestationVote)
				if err := _Engine.contract.UnpackLog(event, "ContestationVote", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseContestationVote is a log parse operation binding the contract event 0x1aa9e4be46e24e1f2e7eeb1613c01629213cd42965d2716e18531b63e552e411.
//
// Solidity: event ContestationVote(address indexed addr, bytes32 indexed task, bool yea)
func (_Engine *EngineFilterer) ParseContestationVote(log types.Log) (*EngineContestationVote, error) {
	event := new(EngineContestationVote)
	if err := _Engine.contract.UnpackLog(event, "ContestationVote", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineContestationVoteFinishIterator is returned from FilterContestationVoteFinish and is used to iterate over the raw logs and unpacked data for ContestationVoteFinish events raised by the Engine contract.
type EngineContestationVoteFinishIterator struct {
	Event *EngineContestationVoteFinish // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineContestationVoteFinishIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineContestationVoteFinish)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineContestationVoteFinish)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineContestationVoteFinishIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineContestationVoteFinishIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineContestationVoteFinish represents a ContestationVoteFinish event raised by the Engine contract.
type EngineContestationVoteFinish struct {
	Id       [32]byte
	StartIdx uint32
	EndIdx   uint32
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterContestationVoteFinish is a free log retrieval operation binding the contract event 0x71d8c71303e35a39162e33a402c9897bf9848388537bac7d5e1b0d202eca4e66.
//
// Solidity: event ContestationVoteFinish(bytes32 indexed id, uint32 indexed start_idx, uint32 end_idx)
func (_Engine *EngineFilterer) FilterContestationVoteFinish(opts *bind.FilterOpts, id [][32]byte, start_idx []uint32) (*EngineContestationVoteFinishIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var start_idxRule []interface{}
	for _, start_idxItem := range start_idx {
		start_idxRule = append(start_idxRule, start_idxItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "ContestationVoteFinish", idRule, start_idxRule)
	if err != nil {
		return nil, err
	}
	return &EngineContestationVoteFinishIterator{contract: _Engine.contract, event: "ContestationVoteFinish", logs: logs, sub: sub}, nil
}

// WatchContestationVoteFinish is a free log subscription operation binding the contract event 0x71d8c71303e35a39162e33a402c9897bf9848388537bac7d5e1b0d202eca4e66.
//
// Solidity: event ContestationVoteFinish(bytes32 indexed id, uint32 indexed start_idx, uint32 end_idx)
func (_Engine *EngineFilterer) WatchContestationVoteFinish(opts *bind.WatchOpts, sink chan<- *EngineContestationVoteFinish, id [][32]byte, start_idx []uint32) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var start_idxRule []interface{}
	for _, start_idxItem := range start_idx {
		start_idxRule = append(start_idxRule, start_idxItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "ContestationVoteFinish", idRule, start_idxRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineContestationVoteFinish)
				if err := _Engine.contract.UnpackLog(event, "ContestationVoteFinish", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseContestationVoteFinish is a log parse operation binding the contract event 0x71d8c71303e35a39162e33a402c9897bf9848388537bac7d5e1b0d202eca4e66.
//
// Solidity: event ContestationVoteFinish(bytes32 indexed id, uint32 indexed start_idx, uint32 end_idx)
func (_Engine *EngineFilterer) ParseContestationVoteFinish(log types.Log) (*EngineContestationVoteFinish, error) {
	event := new(EngineContestationVoteFinish)
	if err := _Engine.contract.UnpackLog(event, "ContestationVoteFinish", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Engine contract.
type EngineInitializedIterator struct {
	Event *EngineInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineInitialized represents a Initialized event raised by the Engine contract.
type EngineInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Engine *EngineFilterer) FilterInitialized(opts *bind.FilterOpts) (*EngineInitializedIterator, error) {

	logs, sub, err := _Engine.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &EngineInitializedIterator{contract: _Engine.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Engine *EngineFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *EngineInitialized) (event.Subscription, error) {

	logs, sub, err := _Engine.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineInitialized)
				if err := _Engine.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Engine *EngineFilterer) ParseInitialized(log types.Log) (*EngineInitialized, error) {
	event := new(EngineInitialized)
	if err := _Engine.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineModelAddrChangedIterator is returned from FilterModelAddrChanged and is used to iterate over the raw logs and unpacked data for ModelAddrChanged events raised by the Engine contract.
type EngineModelAddrChangedIterator struct {
	Event *EngineModelAddrChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineModelAddrChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineModelAddrChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineModelAddrChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineModelAddrChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineModelAddrChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineModelAddrChanged represents a ModelAddrChanged event raised by the Engine contract.
type EngineModelAddrChanged struct {
	Id   [32]byte
	Addr common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterModelAddrChanged is a free log retrieval operation binding the contract event 0xeafdd1f0eea366f4be1889897ed391bb67f32c22bb7ba7d9d3db1888acfd7f63.
//
// Solidity: event ModelAddrChanged(bytes32 indexed id, address addr)
func (_Engine *EngineFilterer) FilterModelAddrChanged(opts *bind.FilterOpts, id [][32]byte) (*EngineModelAddrChangedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "ModelAddrChanged", idRule)
	if err != nil {
		return nil, err
	}
	return &EngineModelAddrChangedIterator{contract: _Engine.contract, event: "ModelAddrChanged", logs: logs, sub: sub}, nil
}

// WatchModelAddrChanged is a free log subscription operation binding the contract event 0xeafdd1f0eea366f4be1889897ed391bb67f32c22bb7ba7d9d3db1888acfd7f63.
//
// Solidity: event ModelAddrChanged(bytes32 indexed id, address addr)
func (_Engine *EngineFilterer) WatchModelAddrChanged(opts *bind.WatchOpts, sink chan<- *EngineModelAddrChanged, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "ModelAddrChanged", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineModelAddrChanged)
				if err := _Engine.contract.UnpackLog(event, "ModelAddrChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseModelAddrChanged is a log parse operation binding the contract event 0xeafdd1f0eea366f4be1889897ed391bb67f32c22bb7ba7d9d3db1888acfd7f63.
//
// Solidity: event ModelAddrChanged(bytes32 indexed id, address addr)
func (_Engine *EngineFilterer) ParseModelAddrChanged(log types.Log) (*EngineModelAddrChanged, error) {
	event := new(EngineModelAddrChanged)
	if err := _Engine.contract.UnpackLog(event, "ModelAddrChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineModelFeeChangedIterator is returned from FilterModelFeeChanged and is used to iterate over the raw logs and unpacked data for ModelFeeChanged events raised by the Engine contract.
type EngineModelFeeChangedIterator struct {
	Event *EngineModelFeeChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineModelFeeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineModelFeeChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineModelFeeChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineModelFeeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineModelFeeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineModelFeeChanged represents a ModelFeeChanged event raised by the Engine contract.
type EngineModelFeeChanged struct {
	Id  [32]byte
	Fee *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterModelFeeChanged is a free log retrieval operation binding the contract event 0x5930669abd6e86dc60e9ac97eba465cc7616d921cbd528312d3dfb4c98f38670.
//
// Solidity: event ModelFeeChanged(bytes32 indexed id, uint256 fee)
func (_Engine *EngineFilterer) FilterModelFeeChanged(opts *bind.FilterOpts, id [][32]byte) (*EngineModelFeeChangedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "ModelFeeChanged", idRule)
	if err != nil {
		return nil, err
	}
	return &EngineModelFeeChangedIterator{contract: _Engine.contract, event: "ModelFeeChanged", logs: logs, sub: sub}, nil
}

// WatchModelFeeChanged is a free log subscription operation binding the contract event 0x5930669abd6e86dc60e9ac97eba465cc7616d921cbd528312d3dfb4c98f38670.
//
// Solidity: event ModelFeeChanged(bytes32 indexed id, uint256 fee)
func (_Engine *EngineFilterer) WatchModelFeeChanged(opts *bind.WatchOpts, sink chan<- *EngineModelFeeChanged, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "ModelFeeChanged", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineModelFeeChanged)
				if err := _Engine.contract.UnpackLog(event, "ModelFeeChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseModelFeeChanged is a log parse operation binding the contract event 0x5930669abd6e86dc60e9ac97eba465cc7616d921cbd528312d3dfb4c98f38670.
//
// Solidity: event ModelFeeChanged(bytes32 indexed id, uint256 fee)
func (_Engine *EngineFilterer) ParseModelFeeChanged(log types.Log) (*EngineModelFeeChanged, error) {
	event := new(EngineModelFeeChanged)
	if err := _Engine.contract.UnpackLog(event, "ModelFeeChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineModelRegisteredIterator is returned from FilterModelRegistered and is used to iterate over the raw logs and unpacked data for ModelRegistered events raised by the Engine contract.
type EngineModelRegisteredIterator struct {
	Event *EngineModelRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineModelRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineModelRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineModelRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineModelRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineModelRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineModelRegistered represents a ModelRegistered event raised by the Engine contract.
type EngineModelRegistered struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterModelRegistered is a free log retrieval operation binding the contract event 0xa4b0af38d049ba81703a0d0e46cc2ff39681210302134046237111a8fb7dee72.
//
// Solidity: event ModelRegistered(bytes32 indexed id)
func (_Engine *EngineFilterer) FilterModelRegistered(opts *bind.FilterOpts, id [][32]byte) (*EngineModelRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "ModelRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &EngineModelRegisteredIterator{contract: _Engine.contract, event: "ModelRegistered", logs: logs, sub: sub}, nil
}

// WatchModelRegistered is a free log subscription operation binding the contract event 0xa4b0af38d049ba81703a0d0e46cc2ff39681210302134046237111a8fb7dee72.
//
// Solidity: event ModelRegistered(bytes32 indexed id)
func (_Engine *EngineFilterer) WatchModelRegistered(opts *bind.WatchOpts, sink chan<- *EngineModelRegistered, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "ModelRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineModelRegistered)
				if err := _Engine.contract.UnpackLog(event, "ModelRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseModelRegistered is a log parse operation binding the contract event 0xa4b0af38d049ba81703a0d0e46cc2ff39681210302134046237111a8fb7dee72.
//
// Solidity: event ModelRegistered(bytes32 indexed id)
func (_Engine *EngineFilterer) ParseModelRegistered(log types.Log) (*EngineModelRegistered, error) {
	event := new(EngineModelRegistered)
	if err := _Engine.contract.UnpackLog(event, "ModelRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Engine contract.
type EngineOwnershipTransferredIterator struct {
	Event *EngineOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineOwnershipTransferred represents a OwnershipTransferred event raised by the Engine contract.
type EngineOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Engine *EngineFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*EngineOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &EngineOwnershipTransferredIterator{contract: _Engine.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Engine *EngineFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *EngineOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineOwnershipTransferred)
				if err := _Engine.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Engine *EngineFilterer) ParseOwnershipTransferred(log types.Log) (*EngineOwnershipTransferred, error) {
	event := new(EngineOwnershipTransferred)
	if err := _Engine.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EnginePausedChangedIterator is returned from FilterPausedChanged and is used to iterate over the raw logs and unpacked data for PausedChanged events raised by the Engine contract.
type EnginePausedChangedIterator struct {
	Event *EnginePausedChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EnginePausedChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EnginePausedChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EnginePausedChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EnginePausedChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EnginePausedChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EnginePausedChanged represents a PausedChanged event raised by the Engine contract.
type EnginePausedChanged struct {
	Paused bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPausedChanged is a free log retrieval operation binding the contract event 0xd83d5281277e107f080e362699d46082adb74e7dc6a9bccbc87d8ae9533add44.
//
// Solidity: event PausedChanged(bool indexed paused)
func (_Engine *EngineFilterer) FilterPausedChanged(opts *bind.FilterOpts, paused []bool) (*EnginePausedChangedIterator, error) {

	var pausedRule []interface{}
	for _, pausedItem := range paused {
		pausedRule = append(pausedRule, pausedItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "PausedChanged", pausedRule)
	if err != nil {
		return nil, err
	}
	return &EnginePausedChangedIterator{contract: _Engine.contract, event: "PausedChanged", logs: logs, sub: sub}, nil
}

// WatchPausedChanged is a free log subscription operation binding the contract event 0xd83d5281277e107f080e362699d46082adb74e7dc6a9bccbc87d8ae9533add44.
//
// Solidity: event PausedChanged(bool indexed paused)
func (_Engine *EngineFilterer) WatchPausedChanged(opts *bind.WatchOpts, sink chan<- *EnginePausedChanged, paused []bool) (event.Subscription, error) {

	var pausedRule []interface{}
	for _, pausedItem := range paused {
		pausedRule = append(pausedRule, pausedItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "PausedChanged", pausedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EnginePausedChanged)
				if err := _Engine.contract.UnpackLog(event, "PausedChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePausedChanged is a log parse operation binding the contract event 0xd83d5281277e107f080e362699d46082adb74e7dc6a9bccbc87d8ae9533add44.
//
// Solidity: event PausedChanged(bool indexed paused)
func (_Engine *EngineFilterer) ParsePausedChanged(log types.Log) (*EnginePausedChanged, error) {
	event := new(EnginePausedChanged)
	if err := _Engine.contract.UnpackLog(event, "PausedChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EnginePauserTransferredIterator is returned from FilterPauserTransferred and is used to iterate over the raw logs and unpacked data for PauserTransferred events raised by the Engine contract.
type EnginePauserTransferredIterator struct {
	Event *EnginePauserTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EnginePauserTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EnginePauserTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EnginePauserTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EnginePauserTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EnginePauserTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EnginePauserTransferred represents a PauserTransferred event raised by the Engine contract.
type EnginePauserTransferred struct {
	To  common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPauserTransferred is a free log retrieval operation binding the contract event 0x5a85b4270fc1e75035e6cd505418ce65e0bcc36cc7eb9ce9e6f8c6181d4cb577.
//
// Solidity: event PauserTransferred(address indexed to)
func (_Engine *EngineFilterer) FilterPauserTransferred(opts *bind.FilterOpts, to []common.Address) (*EnginePauserTransferredIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "PauserTransferred", toRule)
	if err != nil {
		return nil, err
	}
	return &EnginePauserTransferredIterator{contract: _Engine.contract, event: "PauserTransferred", logs: logs, sub: sub}, nil
}

// WatchPauserTransferred is a free log subscription operation binding the contract event 0x5a85b4270fc1e75035e6cd505418ce65e0bcc36cc7eb9ce9e6f8c6181d4cb577.
//
// Solidity: event PauserTransferred(address indexed to)
func (_Engine *EngineFilterer) WatchPauserTransferred(opts *bind.WatchOpts, sink chan<- *EnginePauserTransferred, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "PauserTransferred", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EnginePauserTransferred)
				if err := _Engine.contract.UnpackLog(event, "PauserTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePauserTransferred is a log parse operation binding the contract event 0x5a85b4270fc1e75035e6cd505418ce65e0bcc36cc7eb9ce9e6f8c6181d4cb577.
//
// Solidity: event PauserTransferred(address indexed to)
func (_Engine *EngineFilterer) ParsePauserTransferred(log types.Log) (*EnginePauserTransferred, error) {
	event := new(EnginePauserTransferred)
	if err := _Engine.contract.UnpackLog(event, "PauserTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineSignalCommitmentIterator is returned from FilterSignalCommitment and is used to iterate over the raw logs and unpacked data for SignalCommitment events raised by the Engine contract.
type EngineSignalCommitmentIterator struct {
	Event *EngineSignalCommitment // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineSignalCommitmentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineSignalCommitment)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineSignalCommitment)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineSignalCommitmentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineSignalCommitmentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineSignalCommitment represents a SignalCommitment event raised by the Engine contract.
type EngineSignalCommitment struct {
	Addr       common.Address
	Commitment [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSignalCommitment is a free log retrieval operation binding the contract event 0x09b4c028a2e50fec6f1c6a0163c59e8fbe92b231e5c03ef3adec585e63a14b92.
//
// Solidity: event SignalCommitment(address indexed addr, bytes32 indexed commitment)
func (_Engine *EngineFilterer) FilterSignalCommitment(opts *bind.FilterOpts, addr []common.Address, commitment [][32]byte) (*EngineSignalCommitmentIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var commitmentRule []interface{}
	for _, commitmentItem := range commitment {
		commitmentRule = append(commitmentRule, commitmentItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "SignalCommitment", addrRule, commitmentRule)
	if err != nil {
		return nil, err
	}
	return &EngineSignalCommitmentIterator{contract: _Engine.contract, event: "SignalCommitment", logs: logs, sub: sub}, nil
}

// WatchSignalCommitment is a free log subscription operation binding the contract event 0x09b4c028a2e50fec6f1c6a0163c59e8fbe92b231e5c03ef3adec585e63a14b92.
//
// Solidity: event SignalCommitment(address indexed addr, bytes32 indexed commitment)
func (_Engine *EngineFilterer) WatchSignalCommitment(opts *bind.WatchOpts, sink chan<- *EngineSignalCommitment, addr []common.Address, commitment [][32]byte) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var commitmentRule []interface{}
	for _, commitmentItem := range commitment {
		commitmentRule = append(commitmentRule, commitmentItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "SignalCommitment", addrRule, commitmentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineSignalCommitment)
				if err := _Engine.contract.UnpackLog(event, "SignalCommitment", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSignalCommitment is a log parse operation binding the contract event 0x09b4c028a2e50fec6f1c6a0163c59e8fbe92b231e5c03ef3adec585e63a14b92.
//
// Solidity: event SignalCommitment(address indexed addr, bytes32 indexed commitment)
func (_Engine *EngineFilterer) ParseSignalCommitment(log types.Log) (*EngineSignalCommitment, error) {
	event := new(EngineSignalCommitment)
	if err := _Engine.contract.UnpackLog(event, "SignalCommitment", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineSolutionClaimedIterator is returned from FilterSolutionClaimed and is used to iterate over the raw logs and unpacked data for SolutionClaimed events raised by the Engine contract.
type EngineSolutionClaimedIterator struct {
	Event *EngineSolutionClaimed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineSolutionClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineSolutionClaimed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineSolutionClaimed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineSolutionClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineSolutionClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineSolutionClaimed represents a SolutionClaimed event raised by the Engine contract.
type EngineSolutionClaimed struct {
	Addr common.Address
	Task [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSolutionClaimed is a free log retrieval operation binding the contract event 0x0b76b4ae356796814d36b46f7c500bbd27b2cce1e6059a6fa2bebfd5a389b190.
//
// Solidity: event SolutionClaimed(address indexed addr, bytes32 indexed task)
func (_Engine *EngineFilterer) FilterSolutionClaimed(opts *bind.FilterOpts, addr []common.Address, task [][32]byte) (*EngineSolutionClaimedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var taskRule []interface{}
	for _, taskItem := range task {
		taskRule = append(taskRule, taskItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "SolutionClaimed", addrRule, taskRule)
	if err != nil {
		return nil, err
	}
	return &EngineSolutionClaimedIterator{contract: _Engine.contract, event: "SolutionClaimed", logs: logs, sub: sub}, nil
}

// WatchSolutionClaimed is a free log subscription operation binding the contract event 0x0b76b4ae356796814d36b46f7c500bbd27b2cce1e6059a6fa2bebfd5a389b190.
//
// Solidity: event SolutionClaimed(address indexed addr, bytes32 indexed task)
func (_Engine *EngineFilterer) WatchSolutionClaimed(opts *bind.WatchOpts, sink chan<- *EngineSolutionClaimed, addr []common.Address, task [][32]byte) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var taskRule []interface{}
	for _, taskItem := range task {
		taskRule = append(taskRule, taskItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "SolutionClaimed", addrRule, taskRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineSolutionClaimed)
				if err := _Engine.contract.UnpackLog(event, "SolutionClaimed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSolutionClaimed is a log parse operation binding the contract event 0x0b76b4ae356796814d36b46f7c500bbd27b2cce1e6059a6fa2bebfd5a389b190.
//
// Solidity: event SolutionClaimed(address indexed addr, bytes32 indexed task)
func (_Engine *EngineFilterer) ParseSolutionClaimed(log types.Log) (*EngineSolutionClaimed, error) {
	event := new(EngineSolutionClaimed)
	if err := _Engine.contract.UnpackLog(event, "SolutionClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineSolutionMineableRateChangeIterator is returned from FilterSolutionMineableRateChange and is used to iterate over the raw logs and unpacked data for SolutionMineableRateChange events raised by the Engine contract.
type EngineSolutionMineableRateChangeIterator struct {
	Event *EngineSolutionMineableRateChange // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineSolutionMineableRateChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineSolutionMineableRateChange)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineSolutionMineableRateChange)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineSolutionMineableRateChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineSolutionMineableRateChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineSolutionMineableRateChange represents a SolutionMineableRateChange event raised by the Engine contract.
type EngineSolutionMineableRateChange struct {
	Id   [32]byte
	Rate *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSolutionMineableRateChange is a free log retrieval operation binding the contract event 0x0321e8a918b4c47bf3677852c070983825f30a47bc8d9416691454fa6a727d63.
//
// Solidity: event SolutionMineableRateChange(bytes32 indexed id, uint256 rate)
func (_Engine *EngineFilterer) FilterSolutionMineableRateChange(opts *bind.FilterOpts, id [][32]byte) (*EngineSolutionMineableRateChangeIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "SolutionMineableRateChange", idRule)
	if err != nil {
		return nil, err
	}
	return &EngineSolutionMineableRateChangeIterator{contract: _Engine.contract, event: "SolutionMineableRateChange", logs: logs, sub: sub}, nil
}

// WatchSolutionMineableRateChange is a free log subscription operation binding the contract event 0x0321e8a918b4c47bf3677852c070983825f30a47bc8d9416691454fa6a727d63.
//
// Solidity: event SolutionMineableRateChange(bytes32 indexed id, uint256 rate)
func (_Engine *EngineFilterer) WatchSolutionMineableRateChange(opts *bind.WatchOpts, sink chan<- *EngineSolutionMineableRateChange, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "SolutionMineableRateChange", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineSolutionMineableRateChange)
				if err := _Engine.contract.UnpackLog(event, "SolutionMineableRateChange", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSolutionMineableRateChange is a log parse operation binding the contract event 0x0321e8a918b4c47bf3677852c070983825f30a47bc8d9416691454fa6a727d63.
//
// Solidity: event SolutionMineableRateChange(bytes32 indexed id, uint256 rate)
func (_Engine *EngineFilterer) ParseSolutionMineableRateChange(log types.Log) (*EngineSolutionMineableRateChange, error) {
	event := new(EngineSolutionMineableRateChange)
	if err := _Engine.contract.UnpackLog(event, "SolutionMineableRateChange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineSolutionSubmittedIterator is returned from FilterSolutionSubmitted and is used to iterate over the raw logs and unpacked data for SolutionSubmitted events raised by the Engine contract.
type EngineSolutionSubmittedIterator struct {
	Event *EngineSolutionSubmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineSolutionSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineSolutionSubmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineSolutionSubmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineSolutionSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineSolutionSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineSolutionSubmitted represents a SolutionSubmitted event raised by the Engine contract.
type EngineSolutionSubmitted struct {
	Addr common.Address
	Task [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSolutionSubmitted is a free log retrieval operation binding the contract event 0x957c18b5af8413899ea8a576a4d3fb16839a02c9fccfdce098b6d59ef248525b.
//
// Solidity: event SolutionSubmitted(address indexed addr, bytes32 indexed task)
func (_Engine *EngineFilterer) FilterSolutionSubmitted(opts *bind.FilterOpts, addr []common.Address, task [][32]byte) (*EngineSolutionSubmittedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var taskRule []interface{}
	for _, taskItem := range task {
		taskRule = append(taskRule, taskItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "SolutionSubmitted", addrRule, taskRule)
	if err != nil {
		return nil, err
	}
	return &EngineSolutionSubmittedIterator{contract: _Engine.contract, event: "SolutionSubmitted", logs: logs, sub: sub}, nil
}

// WatchSolutionSubmitted is a free log subscription operation binding the contract event 0x957c18b5af8413899ea8a576a4d3fb16839a02c9fccfdce098b6d59ef248525b.
//
// Solidity: event SolutionSubmitted(address indexed addr, bytes32 indexed task)
func (_Engine *EngineFilterer) WatchSolutionSubmitted(opts *bind.WatchOpts, sink chan<- *EngineSolutionSubmitted, addr []common.Address, task [][32]byte) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var taskRule []interface{}
	for _, taskItem := range task {
		taskRule = append(taskRule, taskItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "SolutionSubmitted", addrRule, taskRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineSolutionSubmitted)
				if err := _Engine.contract.UnpackLog(event, "SolutionSubmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSolutionSubmitted is a log parse operation binding the contract event 0x957c18b5af8413899ea8a576a4d3fb16839a02c9fccfdce098b6d59ef248525b.
//
// Solidity: event SolutionSubmitted(address indexed addr, bytes32 indexed task)
func (_Engine *EngineFilterer) ParseSolutionSubmitted(log types.Log) (*EngineSolutionSubmitted, error) {
	event := new(EngineSolutionSubmitted)
	if err := _Engine.contract.UnpackLog(event, "SolutionSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineStartBlockTimeChangedIterator is returned from FilterStartBlockTimeChanged and is used to iterate over the raw logs and unpacked data for StartBlockTimeChanged events raised by the Engine contract.
type EngineStartBlockTimeChangedIterator struct {
	Event *EngineStartBlockTimeChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineStartBlockTimeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineStartBlockTimeChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineStartBlockTimeChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineStartBlockTimeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineStartBlockTimeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineStartBlockTimeChanged represents a StartBlockTimeChanged event raised by the Engine contract.
type EngineStartBlockTimeChanged struct {
	StartBlockTime uint64
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterStartBlockTimeChanged is a free log retrieval operation binding the contract event 0xa15d6b0930a82638ac4775bfd1b2e9f1e86be67e1bd3a09fa8a77a8f079769d7.
//
// Solidity: event StartBlockTimeChanged(uint64 indexed startBlockTime)
func (_Engine *EngineFilterer) FilterStartBlockTimeChanged(opts *bind.FilterOpts, startBlockTime []uint64) (*EngineStartBlockTimeChangedIterator, error) {

	var startBlockTimeRule []interface{}
	for _, startBlockTimeItem := range startBlockTime {
		startBlockTimeRule = append(startBlockTimeRule, startBlockTimeItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "StartBlockTimeChanged", startBlockTimeRule)
	if err != nil {
		return nil, err
	}
	return &EngineStartBlockTimeChangedIterator{contract: _Engine.contract, event: "StartBlockTimeChanged", logs: logs, sub: sub}, nil
}

// WatchStartBlockTimeChanged is a free log subscription operation binding the contract event 0xa15d6b0930a82638ac4775bfd1b2e9f1e86be67e1bd3a09fa8a77a8f079769d7.
//
// Solidity: event StartBlockTimeChanged(uint64 indexed startBlockTime)
func (_Engine *EngineFilterer) WatchStartBlockTimeChanged(opts *bind.WatchOpts, sink chan<- *EngineStartBlockTimeChanged, startBlockTime []uint64) (event.Subscription, error) {

	var startBlockTimeRule []interface{}
	for _, startBlockTimeItem := range startBlockTime {
		startBlockTimeRule = append(startBlockTimeRule, startBlockTimeItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "StartBlockTimeChanged", startBlockTimeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineStartBlockTimeChanged)
				if err := _Engine.contract.UnpackLog(event, "StartBlockTimeChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStartBlockTimeChanged is a log parse operation binding the contract event 0xa15d6b0930a82638ac4775bfd1b2e9f1e86be67e1bd3a09fa8a77a8f079769d7.
//
// Solidity: event StartBlockTimeChanged(uint64 indexed startBlockTime)
func (_Engine *EngineFilterer) ParseStartBlockTimeChanged(log types.Log) (*EngineStartBlockTimeChanged, error) {
	event := new(EngineStartBlockTimeChanged)
	if err := _Engine.contract.UnpackLog(event, "StartBlockTimeChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineTaskSubmittedIterator is returned from FilterTaskSubmitted and is used to iterate over the raw logs and unpacked data for TaskSubmitted events raised by the Engine contract.
type EngineTaskSubmittedIterator struct {
	Event *EngineTaskSubmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineTaskSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineTaskSubmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineTaskSubmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineTaskSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineTaskSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineTaskSubmitted represents a TaskSubmitted event raised by the Engine contract.
type EngineTaskSubmitted struct {
	Id     [32]byte
	Model  [32]byte
	Fee    *big.Int
	Sender common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTaskSubmitted is a free log retrieval operation binding the contract event 0xc3d3e0544c80e3bb83f62659259ae1574f72a91515ab3cae3dd75cf77e1b0aea.
//
// Solidity: event TaskSubmitted(bytes32 indexed id, bytes32 indexed model, uint256 fee, address indexed sender)
func (_Engine *EngineFilterer) FilterTaskSubmitted(opts *bind.FilterOpts, id [][32]byte, model [][32]byte, sender []common.Address) (*EngineTaskSubmittedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var modelRule []interface{}
	for _, modelItem := range model {
		modelRule = append(modelRule, modelItem)
	}

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "TaskSubmitted", idRule, modelRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &EngineTaskSubmittedIterator{contract: _Engine.contract, event: "TaskSubmitted", logs: logs, sub: sub}, nil
}

// WatchTaskSubmitted is a free log subscription operation binding the contract event 0xc3d3e0544c80e3bb83f62659259ae1574f72a91515ab3cae3dd75cf77e1b0aea.
//
// Solidity: event TaskSubmitted(bytes32 indexed id, bytes32 indexed model, uint256 fee, address indexed sender)
func (_Engine *EngineFilterer) WatchTaskSubmitted(opts *bind.WatchOpts, sink chan<- *EngineTaskSubmitted, id [][32]byte, model [][32]byte, sender []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var modelRule []interface{}
	for _, modelItem := range model {
		modelRule = append(modelRule, modelItem)
	}

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "TaskSubmitted", idRule, modelRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineTaskSubmitted)
				if err := _Engine.contract.UnpackLog(event, "TaskSubmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskSubmitted is a log parse operation binding the contract event 0xc3d3e0544c80e3bb83f62659259ae1574f72a91515ab3cae3dd75cf77e1b0aea.
//
// Solidity: event TaskSubmitted(bytes32 indexed id, bytes32 indexed model, uint256 fee, address indexed sender)
func (_Engine *EngineFilterer) ParseTaskSubmitted(log types.Log) (*EngineTaskSubmitted, error) {
	event := new(EngineTaskSubmitted)
	if err := _Engine.contract.UnpackLog(event, "TaskSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineTreasuryTransferredIterator is returned from FilterTreasuryTransferred and is used to iterate over the raw logs and unpacked data for TreasuryTransferred events raised by the Engine contract.
type EngineTreasuryTransferredIterator struct {
	Event *EngineTreasuryTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineTreasuryTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineTreasuryTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineTreasuryTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineTreasuryTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineTreasuryTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineTreasuryTransferred represents a TreasuryTransferred event raised by the Engine contract.
type EngineTreasuryTransferred struct {
	To  common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterTreasuryTransferred is a free log retrieval operation binding the contract event 0x6bdb9ceff405d990c9b60be9f719fbb80889d5f064e8fd76efd5bea353e90712.
//
// Solidity: event TreasuryTransferred(address indexed to)
func (_Engine *EngineFilterer) FilterTreasuryTransferred(opts *bind.FilterOpts, to []common.Address) (*EngineTreasuryTransferredIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "TreasuryTransferred", toRule)
	if err != nil {
		return nil, err
	}
	return &EngineTreasuryTransferredIterator{contract: _Engine.contract, event: "TreasuryTransferred", logs: logs, sub: sub}, nil
}

// WatchTreasuryTransferred is a free log subscription operation binding the contract event 0x6bdb9ceff405d990c9b60be9f719fbb80889d5f064e8fd76efd5bea353e90712.
//
// Solidity: event TreasuryTransferred(address indexed to)
func (_Engine *EngineFilterer) WatchTreasuryTransferred(opts *bind.WatchOpts, sink chan<- *EngineTreasuryTransferred, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "TreasuryTransferred", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineTreasuryTransferred)
				if err := _Engine.contract.UnpackLog(event, "TreasuryTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTreasuryTransferred is a log parse operation binding the contract event 0x6bdb9ceff405d990c9b60be9f719fbb80889d5f064e8fd76efd5bea353e90712.
//
// Solidity: event TreasuryTransferred(address indexed to)
func (_Engine *EngineFilterer) ParseTreasuryTransferred(log types.Log) (*EngineTreasuryTransferred, error) {
	event := new(EngineTreasuryTransferred)
	if err := _Engine.contract.UnpackLog(event, "TreasuryTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineValidatorDepositIterator is returned from FilterValidatorDeposit and is used to iterate over the raw logs and unpacked data for ValidatorDeposit events raised by the Engine contract.
type EngineValidatorDepositIterator struct {
	Event *EngineValidatorDeposit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineValidatorDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineValidatorDeposit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineValidatorDeposit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineValidatorDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineValidatorDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineValidatorDeposit represents a ValidatorDeposit event raised by the Engine contract.
type EngineValidatorDeposit struct {
	Addr      common.Address
	Validator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterValidatorDeposit is a free log retrieval operation binding the contract event 0x8d4844488c19a90828439e71d14ebad860806d04f8ef8b25a82179fab2699b89.
//
// Solidity: event ValidatorDeposit(address indexed addr, address indexed validator, uint256 amount)
func (_Engine *EngineFilterer) FilterValidatorDeposit(opts *bind.FilterOpts, addr []common.Address, validator []common.Address) (*EngineValidatorDepositIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "ValidatorDeposit", addrRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return &EngineValidatorDepositIterator{contract: _Engine.contract, event: "ValidatorDeposit", logs: logs, sub: sub}, nil
}

// WatchValidatorDeposit is a free log subscription operation binding the contract event 0x8d4844488c19a90828439e71d14ebad860806d04f8ef8b25a82179fab2699b89.
//
// Solidity: event ValidatorDeposit(address indexed addr, address indexed validator, uint256 amount)
func (_Engine *EngineFilterer) WatchValidatorDeposit(opts *bind.WatchOpts, sink chan<- *EngineValidatorDeposit, addr []common.Address, validator []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "ValidatorDeposit", addrRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineValidatorDeposit)
				if err := _Engine.contract.UnpackLog(event, "ValidatorDeposit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorDeposit is a log parse operation binding the contract event 0x8d4844488c19a90828439e71d14ebad860806d04f8ef8b25a82179fab2699b89.
//
// Solidity: event ValidatorDeposit(address indexed addr, address indexed validator, uint256 amount)
func (_Engine *EngineFilterer) ParseValidatorDeposit(log types.Log) (*EngineValidatorDeposit, error) {
	event := new(EngineValidatorDeposit)
	if err := _Engine.contract.UnpackLog(event, "ValidatorDeposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineValidatorWithdrawIterator is returned from FilterValidatorWithdraw and is used to iterate over the raw logs and unpacked data for ValidatorWithdraw events raised by the Engine contract.
type EngineValidatorWithdrawIterator struct {
	Event *EngineValidatorWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineValidatorWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineValidatorWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineValidatorWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineValidatorWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineValidatorWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineValidatorWithdraw represents a ValidatorWithdraw event raised by the Engine contract.
type EngineValidatorWithdraw struct {
	Addr   common.Address
	To     common.Address
	Count  *big.Int
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterValidatorWithdraw is a free log retrieval operation binding the contract event 0x109aeff667601aad33abcb7c5df1617754eefd18253d586a95c198cb479b5bd4.
//
// Solidity: event ValidatorWithdraw(address indexed addr, address indexed to, uint256 indexed count, uint256 amount)
func (_Engine *EngineFilterer) FilterValidatorWithdraw(opts *bind.FilterOpts, addr []common.Address, to []common.Address, count []*big.Int) (*EngineValidatorWithdrawIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var countRule []interface{}
	for _, countItem := range count {
		countRule = append(countRule, countItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "ValidatorWithdraw", addrRule, toRule, countRule)
	if err != nil {
		return nil, err
	}
	return &EngineValidatorWithdrawIterator{contract: _Engine.contract, event: "ValidatorWithdraw", logs: logs, sub: sub}, nil
}

// WatchValidatorWithdraw is a free log subscription operation binding the contract event 0x109aeff667601aad33abcb7c5df1617754eefd18253d586a95c198cb479b5bd4.
//
// Solidity: event ValidatorWithdraw(address indexed addr, address indexed to, uint256 indexed count, uint256 amount)
func (_Engine *EngineFilterer) WatchValidatorWithdraw(opts *bind.WatchOpts, sink chan<- *EngineValidatorWithdraw, addr []common.Address, to []common.Address, count []*big.Int) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var countRule []interface{}
	for _, countItem := range count {
		countRule = append(countRule, countItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "ValidatorWithdraw", addrRule, toRule, countRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineValidatorWithdraw)
				if err := _Engine.contract.UnpackLog(event, "ValidatorWithdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorWithdraw is a log parse operation binding the contract event 0x109aeff667601aad33abcb7c5df1617754eefd18253d586a95c198cb479b5bd4.
//
// Solidity: event ValidatorWithdraw(address indexed addr, address indexed to, uint256 indexed count, uint256 amount)
func (_Engine *EngineFilterer) ParseValidatorWithdraw(log types.Log) (*EngineValidatorWithdraw, error) {
	event := new(EngineValidatorWithdraw)
	if err := _Engine.contract.UnpackLog(event, "ValidatorWithdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineValidatorWithdrawCancelledIterator is returned from FilterValidatorWithdrawCancelled and is used to iterate over the raw logs and unpacked data for ValidatorWithdrawCancelled events raised by the Engine contract.
type EngineValidatorWithdrawCancelledIterator struct {
	Event *EngineValidatorWithdrawCancelled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineValidatorWithdrawCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineValidatorWithdrawCancelled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineValidatorWithdrawCancelled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineValidatorWithdrawCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineValidatorWithdrawCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineValidatorWithdrawCancelled represents a ValidatorWithdrawCancelled event raised by the Engine contract.
type EngineValidatorWithdrawCancelled struct {
	Addr  common.Address
	Count *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterValidatorWithdrawCancelled is a free log retrieval operation binding the contract event 0xf9d5cfcdfe9803069225971ba315f4302add9b477c3441cc363bc38fdb065b90.
//
// Solidity: event ValidatorWithdrawCancelled(address indexed addr, uint256 indexed count)
func (_Engine *EngineFilterer) FilterValidatorWithdrawCancelled(opts *bind.FilterOpts, addr []common.Address, count []*big.Int) (*EngineValidatorWithdrawCancelledIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var countRule []interface{}
	for _, countItem := range count {
		countRule = append(countRule, countItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "ValidatorWithdrawCancelled", addrRule, countRule)
	if err != nil {
		return nil, err
	}
	return &EngineValidatorWithdrawCancelledIterator{contract: _Engine.contract, event: "ValidatorWithdrawCancelled", logs: logs, sub: sub}, nil
}

// WatchValidatorWithdrawCancelled is a free log subscription operation binding the contract event 0xf9d5cfcdfe9803069225971ba315f4302add9b477c3441cc363bc38fdb065b90.
//
// Solidity: event ValidatorWithdrawCancelled(address indexed addr, uint256 indexed count)
func (_Engine *EngineFilterer) WatchValidatorWithdrawCancelled(opts *bind.WatchOpts, sink chan<- *EngineValidatorWithdrawCancelled, addr []common.Address, count []*big.Int) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var countRule []interface{}
	for _, countItem := range count {
		countRule = append(countRule, countItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "ValidatorWithdrawCancelled", addrRule, countRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineValidatorWithdrawCancelled)
				if err := _Engine.contract.UnpackLog(event, "ValidatorWithdrawCancelled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorWithdrawCancelled is a log parse operation binding the contract event 0xf9d5cfcdfe9803069225971ba315f4302add9b477c3441cc363bc38fdb065b90.
//
// Solidity: event ValidatorWithdrawCancelled(address indexed addr, uint256 indexed count)
func (_Engine *EngineFilterer) ParseValidatorWithdrawCancelled(log types.Log) (*EngineValidatorWithdrawCancelled, error) {
	event := new(EngineValidatorWithdrawCancelled)
	if err := _Engine.contract.UnpackLog(event, "ValidatorWithdrawCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineValidatorWithdrawInitiatedIterator is returned from FilterValidatorWithdrawInitiated and is used to iterate over the raw logs and unpacked data for ValidatorWithdrawInitiated events raised by the Engine contract.
type EngineValidatorWithdrawInitiatedIterator struct {
	Event *EngineValidatorWithdrawInitiated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineValidatorWithdrawInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineValidatorWithdrawInitiated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineValidatorWithdrawInitiated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineValidatorWithdrawInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineValidatorWithdrawInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineValidatorWithdrawInitiated represents a ValidatorWithdrawInitiated event raised by the Engine contract.
type EngineValidatorWithdrawInitiated struct {
	Addr       common.Address
	Count      *big.Int
	UnlockTime *big.Int
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterValidatorWithdrawInitiated is a free log retrieval operation binding the contract event 0xcc7e0e76b20394ef965d71c4111d21e1b322bb8775dc2c7acd29eb0d3c3dd96d.
//
// Solidity: event ValidatorWithdrawInitiated(address indexed addr, uint256 indexed count, uint256 unlockTime, uint256 amount)
func (_Engine *EngineFilterer) FilterValidatorWithdrawInitiated(opts *bind.FilterOpts, addr []common.Address, count []*big.Int) (*EngineValidatorWithdrawInitiatedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var countRule []interface{}
	for _, countItem := range count {
		countRule = append(countRule, countItem)
	}

	logs, sub, err := _Engine.contract.FilterLogs(opts, "ValidatorWithdrawInitiated", addrRule, countRule)
	if err != nil {
		return nil, err
	}
	return &EngineValidatorWithdrawInitiatedIterator{contract: _Engine.contract, event: "ValidatorWithdrawInitiated", logs: logs, sub: sub}, nil
}

// WatchValidatorWithdrawInitiated is a free log subscription operation binding the contract event 0xcc7e0e76b20394ef965d71c4111d21e1b322bb8775dc2c7acd29eb0d3c3dd96d.
//
// Solidity: event ValidatorWithdrawInitiated(address indexed addr, uint256 indexed count, uint256 unlockTime, uint256 amount)
func (_Engine *EngineFilterer) WatchValidatorWithdrawInitiated(opts *bind.WatchOpts, sink chan<- *EngineValidatorWithdrawInitiated, addr []common.Address, count []*big.Int) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var countRule []interface{}
	for _, countItem := range count {
		countRule = append(countRule, countItem)
	}

	logs, sub, err := _Engine.contract.WatchLogs(opts, "ValidatorWithdrawInitiated", addrRule, countRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineValidatorWithdrawInitiated)
				if err := _Engine.contract.UnpackLog(event, "ValidatorWithdrawInitiated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorWithdrawInitiated is a log parse operation binding the contract event 0xcc7e0e76b20394ef965d71c4111d21e1b322bb8775dc2c7acd29eb0d3c3dd96d.
//
// Solidity: event ValidatorWithdrawInitiated(address indexed addr, uint256 indexed count, uint256 unlockTime, uint256 amount)
func (_Engine *EngineFilterer) ParseValidatorWithdrawInitiated(log types.Log) (*EngineValidatorWithdrawInitiated, error) {
	event := new(EngineValidatorWithdrawInitiated)
	if err := _Engine.contract.UnpackLog(event, "ValidatorWithdrawInitiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineVersionChangedIterator is returned from FilterVersionChanged and is used to iterate over the raw logs and unpacked data for VersionChanged events raised by the Engine contract.
type EngineVersionChangedIterator struct {
	Event *EngineVersionChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EngineVersionChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineVersionChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EngineVersionChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EngineVersionChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineVersionChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineVersionChanged represents a VersionChanged event raised by the Engine contract.
type EngineVersionChanged struct {
	Version *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterVersionChanged is a free log retrieval operation binding the contract event 0x8c854a81cb5c93e7e482d30fb9c6f88fdbdb320f10f7a853c2263659b54e563f.
//
// Solidity: event VersionChanged(uint256 version)
func (_Engine *EngineFilterer) FilterVersionChanged(opts *bind.FilterOpts) (*EngineVersionChangedIterator, error) {

	logs, sub, err := _Engine.contract.FilterLogs(opts, "VersionChanged")
	if err != nil {
		return nil, err
	}
	return &EngineVersionChangedIterator{contract: _Engine.contract, event: "VersionChanged", logs: logs, sub: sub}, nil
}

// WatchVersionChanged is a free log subscription operation binding the contract event 0x8c854a81cb5c93e7e482d30fb9c6f88fdbdb320f10f7a853c2263659b54e563f.
//
// Solidity: event VersionChanged(uint256 version)
func (_Engine *EngineFilterer) WatchVersionChanged(opts *bind.WatchOpts, sink chan<- *EngineVersionChanged) (event.Subscription, error) {

	logs, sub, err := _Engine.contract.WatchLogs(opts, "VersionChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineVersionChanged)
				if err := _Engine.contract.UnpackLog(event, "VersionChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVersionChanged is a log parse operation binding the contract event 0x8c854a81cb5c93e7e482d30fb9c6f88fdbdb320f10f7a853c2263659b54e563f.
//
// Solidity: event VersionChanged(uint256 version)
func (_Engine *EngineFilterer) ParseVersionChanged(log types.Log) (*EngineVersionChanged, error) {
	event := new(EngineVersionChanged)
	if err := _Engine.contract.UnpackLog(event, "VersionChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
