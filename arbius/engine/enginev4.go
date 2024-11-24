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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv18_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"denominator\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv_Overflow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PRBMath_SD59x18_Abs_MinSD59x18\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PRBMath_SD59x18_Div_InputTooSmall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"SD59x18\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"SD59x18\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"PRBMath_SD59x18_Div_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"SD59x18\",\"name\":\"x\",\"type\":\"int256\"}],\"name\":\"PRBMath_SD59x18_Exp2_InputTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PRBMath_SD59x18_Mul_InputTooSmall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"SD59x18\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"SD59x18\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"PRBMath_SD59x18_Mul_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"UD60x18\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"PRBMath_UD60x18_Exp2_InputTooBig\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"task\",\"type\":\"bytes32\"}],\"name\":\"ContestationSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"task\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"yea\",\"type\":\"bool\"}],\"name\":\"ContestationVote\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"start_idx\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"end_idx\",\"type\":\"uint32\"}],\"name\":\"ContestationVoteFinish\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ModelRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"name\":\"PausedChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PauserTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"}],\"name\":\"SignalCommitment\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"task\",\"type\":\"bytes32\"}],\"name\":\"SolutionClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"SolutionMineableRateChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"task\",\"type\":\"bytes32\"}],\"name\":\"SolutionSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"startBlockTime\",\"type\":\"uint64\"}],\"name\":\"StartBlockTimeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"model\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"TaskSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"TreasuryTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ValidatorDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ValidatorWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"ValidatorWithdrawCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unlockTime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ValidatorWithdrawInitiated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"}],\"name\":\"VersionChanged\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"accruedFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"baseToken\",\"outputs\":[{\"internalType\":\"contractIBaseToken\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"taskids_\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes[]\",\"name\":\"cids_\",\"type\":\"bytes[]\"}],\"name\":\"bulkSubmitSolution\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"version_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"owner_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"model_\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"input_\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"n_\",\"type\":\"uint256\"}],\"name\":\"bulkSubmitTask\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"count_\",\"type\":\"uint256\"}],\"name\":\"cancelValidatorWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"}],\"name\":\"claimSolution\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"commitments\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"contestationVoteExtensionTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"amnt_\",\"type\":\"uint32\"}],\"name\":\"contestationVoteFinish\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"contestationVoteNays\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"contestationVoteYeas\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"contestationVoted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"contestationVotedIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"contestations\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"blocktime\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"finish_start_index\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"slashAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"t\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ts\",\"type\":\"uint256\"}],\"name\":\"diffMul\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"exitValidatorMinUnlockTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"cid_\",\"type\":\"bytes\"}],\"name\":\"generateCommitment\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"content_\",\"type\":\"bytes\"}],\"name\":\"generateIPFSCID\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPsuedoTotalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSlashAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidatorMinimum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cid\",\"type\":\"bytes\"}],\"internalType\":\"structModel\",\"name\":\"o_\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender_\",\"type\":\"address\"}],\"name\":\"hashModel\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"model\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"blocktime\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"cid\",\"type\":\"bytes\"}],\"internalType\":\"structTask\",\"name\":\"o_\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"prevhash_\",\"type\":\"bytes32\"}],\"name\":\"hashTask\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"}],\"name\":\"initiateValidatorWithdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"lastContestationLossTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"lastSolutionSubmission\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxContestationValidatorStakeSince\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minClaimSolutionTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minContestationVotePeriodTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minRetractionWaitTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"models\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cid\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauser\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"pendingValidatorWithdrawRequests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"unlockTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"pendingValidatorWithdrawRequestsCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"prevhash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"template_\",\"type\":\"bytes\"}],\"name\":\"registerModel\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"retractionFeePercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"t\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ts\",\"type\":\"uint256\"}],\"name\":\"reward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"paused_\",\"type\":\"bool\"}],\"name\":\"setPaused\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"model_\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"rate_\",\"type\":\"uint256\"}],\"name\":\"setSolutionMineableRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"startBlockTime_\",\"type\":\"uint64\"}],\"name\":\"setStartBlockTime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"veStaking_\",\"type\":\"address\"}],\"name\":\"setVeStaking\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"version_\",\"type\":\"uint256\"}],\"name\":\"setVersion\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitment_\",\"type\":\"bytes32\"}],\"name\":\"signalCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"slashAmountPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"solutionFeePercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"solutionRateLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"solutions\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"blocktime\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"claimed\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"cid\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"solutionsStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"solutionsStakeAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"startBlockTime\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"}],\"name\":\"submitContestation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"cid_\",\"type\":\"bytes\"}],\"name\":\"submitSolution\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"version_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"owner_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"model_\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"input_\",\"type\":\"bytes\"}],\"name\":\"submitTask\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"t\",\"type\":\"uint256\"}],\"name\":\"targetTs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"taskOwnerRewardPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"tasks\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"model\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"blocktime\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"cid\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalHeld\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to_\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to_\",\"type\":\"address\"}],\"name\":\"transferPauser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to_\",\"type\":\"address\"}],\"name\":\"transferTreasury\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"treasury\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"treasuryRewardPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"}],\"name\":\"validatorCanVote\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"}],\"name\":\"validatorDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validatorMinimumPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"count_\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to_\",\"type\":\"address\"}],\"name\":\"validatorWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validatorWithdrawPendingAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"staked\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"since\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"veRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"veStaking\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"yea_\",\"type\":\"bool\"}],\"name\":\"voteOnContestation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"}],\"name\":\"votingPeriodEnded\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawAccruedFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50620000226200002860201b60201c565b620001d2565b600060019054906101000a900460ff16156200007b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620000729062000175565b60405180910390fd5b60ff801660008054906101000a900460ff1660ff1614620000ec5760ff6000806101000a81548160ff021916908360ff1602179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb384740249860ff604051620000e39190620001b5565b60405180910390a15b565b600082825260208201905092915050565b7f496e697469616c697a61626c653a20636f6e747261637420697320696e69746960008201527f616c697a696e6700000000000000000000000000000000000000000000000000602082015250565b60006200015d602783620000ee565b91506200016a82620000ff565b604082019050919050565b6000602082019050818103600083015262000190816200014e565b9050919050565b600060ff82169050919050565b620001af8162000197565b82525050565b6000602082019050620001cc6000830184620001a4565b92915050565b61a8cb80620001e26000396000f3fe608060405234801561001057600080fd5b506004361061047f5760003560e01c80638129fc1c11610257578063c17ddb2a11610146578063d2992baa116100c3578063e579f50011610087578063e579f50014610eb8578063f1b8989d14610eed578063f2fde38b14610f0b578063f43cc77314610f27578063fa52c7d814610f455761047f565b8063d2992baa14610de7578063d33b2ef514610e18578063d8a6021c14610e4b578063dc06a89f14610e67578063e236f46b14610e855761047f565b8063cbd2422d1161010a578063cbd2422d14610d0b578063cf596e4514610d27578063d1f0c94114610d57578063d2307ae414610d87578063d278094014610db75761047f565b8063c17ddb2a14610c63578063c1f7272314610c81578063c24b563114610c9f578063c31784be14610ccf578063c55dae6314610ced5761047f565b806396bb02c3116101d4578063a4fa8d5711610198578063a4fa8d5714610bbf578063a53e252514610bef578063a8f837f314610c0d578063ada82c7d14610c29578063b4dc35b714610c335761047f565b806396bb02c314610b295780639b97511914610b475780639fd0506d14610b65578063a1975adf14610b83578063a2492a9014610ba15761047f565b80638da5cb5b1161021b5780638da5cb5b14610a975780638e6d86fd14610ab55780639280944414610ad357806393a090ec14610af157806393f1f8ac14610b0d5761047f565b80638129fc1c146109f357806382b5077f146109fd5780638365779514610a1b578063839df94514610a4b5780638b4d7b3514610a7b5761047f565b80633d57f5d91161037357806365d445fb116102f057806375c70509116102b457806375c705091461094c578063763253bb1461097f57806377286d171461099b5780637881c5e6146109b75780637b36006a146109d55761047f565b806365d445fb146108ce578063671f8152146108ea578063682c205814610906578063715018a61461092457806372dc0ee11461092e5761047f565b8063506ea7de11610337578063506ea7de1461083c57806354fd4d501461085857806356914caf146108765780635c975abb1461089257806361d027b3146108b05761047f565b80633d57f5d914610786578063408def1e146107a457806340e8c56d146107c05780634421ea21146107f05780634ff03efa1461080c5761047f565b80631825c20e116104015780632258d105116103c55780632258d105146106ce5780632943a490146106ec578063303fb0d614610708578063393cb1c7146107385780633d18b912146107685761047f565b80631825c20e146105f25780631b75c43e1461060e5780631f88ea1c1461063e578063218a30481461066e578063218e68591461069e5761047f565b80630c18d4ce116104485780630c18d4ce1461053a5780630d468d95146105585780631466b63a1461057657806316c38b3c146105a657806317f3e041146105c25761047f565b8062fd70821461048457806305d1bc26146104a257806308745dd1146104d257806308afe0eb146104ee5780630a9857371461050a575b600080fd5b61048c610f77565b6040516104999190617c8e565b60405180910390f35b6104bc60048036038101906104b79190617cf3565b610f7d565b6040516104c99190617d3b565b60405180910390f35b6104ec60048036038101906104e79190617e7e565b611019565b005b61050860048036038101906105039190617f18565b611247565b005b610524600480360381019061051f9190617fc7565b6114a9565b6040516105319190617c8e565b60405180910390f35b6105426117a4565b60405161054f9190618017565b60405180910390f35b6105606117be565b60405161056d9190617c8e565b60405180910390f35b610590600480360381019061058b9190618265565b6117c4565b60405161059d91906182e3565b60405180910390f35b6105c060048036038101906105bb919061832a565b61180a565b005b6105dc60048036038101906105d79190617cf3565b6118e6565b6040516105e99190617c8e565b60405180910390f35b61060c60048036038101906106079190618357565b6118fe565b005b61062860048036038101906106239190618397565b6119f1565b6040516106359190617c8e565b60405180910390f35b610658600480360381019061065391906183c4565b611a09565b6040516106659190617c8e565b60405180910390f35b61068860048036038101906106839190618498565b611bcd565b60405161069591906182e3565b60405180910390f35b6106b860048036038101906106b39190618397565b611c10565b6040516106c59190617c8e565b60405180910390f35b6106d6611c28565b6040516106e39190617c8e565b60405180910390f35b61070660048036038101906107019190618397565b611c7a565b005b610722600480360381019061071d91906184f4565b611cc6565b60405161072f9190618543565b60405180910390f35b610752600480360381019061074d919061855e565b611d14565b60405161075f91906182e3565b60405180910390f35b610770611d4d565b60405161077d9190617c8e565b60405180910390f35b61078e611d90565b60405161079b9190617c8e565b60405180910390f35b6107be60048036038101906107b99190617fc7565b611de2565b005b6107da60048036038101906107d591906185d2565b611e2b565b6040516107e7919061869e565b60405180910390f35b61080a60048036038101906108059190618397565b611e3f565b005b610826600480360381019061082191906186c0565b611ece565b60405161083391906182e3565b60405180910390f35b61085660048036038101906108519190617cf3565b61214f565b005b61086061225b565b60405161086d9190617c8e565b60405180910390f35b610890600480360381019061088b9190618734565b612261565b005b61089a6122cb565b6040516108a79190617d3b565b60405180910390f35b6108b86122de565b6040516108c59190618543565b60405180910390f35b6108e860048036038101906108e39190618840565b612304565b005b61090460048036038101906108ff9190617cf3565b6123cd565b005b61090e6128bd565b60405161091b9190617c8e565b60405180910390f35b61092c6128c3565b005b6109366128d7565b6040516109439190617c8e565b60405180910390f35b61096660048036038101906109619190617cf3565b6128dd565b60405161097694939291906188c1565b60405180910390f35b6109996004803603810190610994919061890d565b6129d6565b005b6109b560048036038101906109b09190617cf3565b612dfb565b005b6109bf61337e565b6040516109cc9190617c8e565b60405180910390f35b6109dd613477565b6040516109ea9190617c8e565b60405180910390f35b6109fb61347d565b005b610a05613576565b604051610a129190618543565b60405180910390f35b610a356004803603810190610a30919061894d565b61359c565b604051610a429190617c8e565b60405180910390f35b610a656004803603810190610a609190617cf3565b6137da565b604051610a729190617c8e565b60405180910390f35b610a956004803603810190610a9091906189c9565b6137f2565b005b610a9f61436f565b604051610aac9190618543565b60405180910390f35b610abd614399565b604051610aca9190617c8e565b60405180910390f35b610adb61439f565b604051610ae89190617c8e565b60405180910390f35b610b0b6004803603810190610b069190618a09565b6143a5565b005b610b276004803603810190610b2291906184f4565b61476a565b005b610b3161486e565b604051610b3e9190617c8e565b60405180910390f35b610b4f614874565b604051610b5c9190617c8e565b60405180910390f35b610b6d61487a565b604051610b7a9190618543565b60405180910390f35b610b8b6148a0565b604051610b989190617c8e565b60405180910390f35b610ba96148a6565b604051610bb69190617c8e565b60405180910390f35b610bd96004803603810190610bd491906183c4565b6148ac565b604051610be69190617c8e565b60405180910390f35b610bf7614933565b604051610c049190617c8e565b60405180910390f35b610c276004803603810190610c229190618a49565b614939565b005b610c316149a4565b005b610c4d6004803603810190610c489190617cf3565b614ade565b604051610c5a9190617c8e565b60405180910390f35b610c6b614af6565b604051610c7891906182e3565b60405180910390f35b610c89614afc565b604051610c969190617c8e565b60405180910390f35b610cb96004803603810190610cb49190618397565b614b02565b604051610cc69190617c8e565b60405180910390f35b610cd7614b1a565b604051610ce49190617c8e565b60405180910390f35b610cf5614b20565b604051610d029190618ad5565b60405180910390f35b610d256004803603810190610d209190617fc7565b614b46565b005b610d416004803603810190610d3c9190617fc7565b614d55565b604051610d4e9190617c8e565b60405180910390f35b610d716004803603810190610d6c91906184f4565b614e15565b604051610d7e9190618543565b60405180910390f35b610da16004803603810190610d9c9190618397565b614e63565b604051610dae9190617c8e565b60405180910390f35b610dd16004803603810190610dcc9190618af0565b614e7b565b604051610dde9190617d3b565b60405180910390f35b610e016004803603810190610dfc9190618a09565b614eaa565b604051610e0f929190618b30565b60405180910390f35b610e326004803603810190610e2d9190617cf3565b614edb565b604051610e429493929190618b68565b60405180910390f35b610e656004803603810190610e609190618397565b614f4f565b005b610e6f614fde565b604051610e7c9190617c8e565b60405180910390f35b610e9f6004803603810190610e9a9190617cf3565b614fe4565b604051610eaf9493929190618bad565b60405180910390f35b610ed26004803603810190610ecd9190617cf3565b6150bc565b604051610ee496959493929190618c08565b60405180910390f35b610ef56151c1565b604051610f029190617c8e565b60405180910390f35b610f256004803603810190610f209190618397565b6151c7565b005b610f2f6151db565b604051610f3c9190617c8e565b60405180910390f35b610f5f6004803603810190610f5a9190618397565b6151e1565b604051610f6e93929190618c70565b60405180910390f35b60725481565b6000608a5460826000848152602001908152602001600020805490506081600085815260200190815260200160002080549050610fba9190618cd6565b610fc49190618d0a565b607354607e600085815260200190815260200160002060000160149054906101000a900467ffffffffffffffff1667ffffffffffffffff166110069190618cd6565b6110109190618cd6565b42119050919050565b606760149054906101000a900460ff1615611069576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161106090618da9565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff166076600086815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160361110e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161110590618e15565b60405180910390fd5b6076600085815260200190815260200160002060000154831015611167576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161115e90618e81565b60405180910390fd5b6000611173838361522b565b90506111828787878785615353565b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd3330876040518463ffffffff1660e01b81526004016111e193929190618ea1565b6020604051808303816000875af1158015611200573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112249190618eed565b5083608660008282546112379190618cd6565b9250508190555050505050505050565b606760149054906101000a900460ff1615611297576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161128e90618da9565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff166076600087815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160361133c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161133390618e15565b60405180910390fd5b6076600086815260200190815260200160002060000154841015611395576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161138c90618e81565b60405180910390fd5b60006113a1848461522b565b905060005b828110156113cc576113bb8989898986615353565b806113c590618f1a565b90506113a6565b50606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd333085896114199190618d0a565b6040518463ffffffff1660e01b815260040161143793929190618ea1565b6020604051808303816000875af1158015611456573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061147a9190618eed565b5081856114879190618d0a565b608660008282546114989190618cd6565b925050819055505050505050505050565b6000606760149054906101000a900460ff16156114fb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016114f290618da9565b60405180910390fd5b81607a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054607760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001546115899190618f62565b10156115ca576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016115c190618fbc565b60405180910390fd5b6000607554426115da9190618cd6565b90506001607860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461162c9190618cd6565b925050819055506000607860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050604051806040016040528083815260200185815250607960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000838152602001908152602001600020600082015181600001556020820151816001015590505083607a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546117429190618cd6565b92505081905550803373ffffffffffffffffffffffffffffffffffffffff167fcc7e0e76b20394ef965d71c4111d21e1b322bb8775dc2c7acd29eb0d3c3dd96d8487604051611792929190618b30565b60405180910390a38092505050919050565b606a60009054906101000a900467ffffffffffffffff1681565b608c5481565b60008282856000015186602001518760a001516040516020016117eb959493929190618fdc565b6040516020818303038152906040528051906020012090509392505050565b606760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461189a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161189190619082565b60405180910390fd5b80606760146101000a81548160ff0219169083151502179055508015157fd83d5281277e107f080e362699d46082adb74e7dc6a9bccbc87d8ae9533add4460405160405180910390a250565b60806020528060005260406000206000915090505481565b606760149054906101000a900460ff161561194e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161194590618da9565b60405180910390fd5b611957336154f0565b611996576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161198d906190ee565b60405180910390fd5b60006119a2338461359c565b146119e2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016119d99061915a565b60405180910390fd5b6119ed828233615590565b5050565b607a6020528060005260406000206000915090505481565b60008083118015611a1a5750600082115b611a59576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611a50906191c6565b60405180910390fd5b6000611a6484614d55565b90506000611a8b611a74836157b0565b611a7d866157b0565b6157ba90919063ffffffff16565b9050670cf4ad2e86166c9c611a9f82615953565b1215611ab85768056bc75e2d6310000092505050611bc7565b6000611acb670de0b6b3a76400006157b0565b90506000611ae168056bc75e2d631000006157b0565b90506000611b2e611b1f84611b1185611b03888a61595d90919063ffffffff16565b61598b90919063ffffffff16565b61595d90919063ffffffff16565b84615b1b90919063ffffffff16565b90506801158e460913d00000611b4382615953565b12611b5657600095505050505050611bc7565b611b72611b6360006157b0565b82615b4990919063ffffffff16565b15611b9c57611b90611b8b611b8683615b66565b615bf6565b615953565b95505050505050611bc7565b611bbf611bba611bab83615bf6565b856157ba90919063ffffffff16565b615953565b955050505050505b92915050565b600081836020015184600001518560600151604051602001611bf294939291906191e6565b60405160208183030381529060405280519060200120905092915050565b60856020528060005260406000206000915090505481565b600080611c3361337e565b9050670de0b6b3a7640000606c54670de0b6b3a7640000611c549190618f62565b82611c5f9190618d0a565b611c699190619261565b81611c749190618f62565b91505090565b611c82615d14565b80608b60006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60826020528160005260406000208181548110611ce257600080fd5b906000526020600020016000915091509054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600084848484604051602001611d2d94939291906192bf565b604051602081830303815290604052805190602001209050949350505050565b6000611d8b606a60009054906101000a900467ffffffffffffffff1667ffffffffffffffff1642611d7e9190618f62565b611d8661337e565b6148ac565b905090565b600080611d9b61337e565b9050670de0b6b3a7640000606d54670de0b6b3a7640000611dbc9190618f62565b82611dc79190618d0a565b611dd19190619261565b81611ddc9190618f62565b91505090565b611dea615d14565b80606b819055507f8c854a81cb5c93e7e482d30fb9c6f88fdbdb320f10f7a853c2263659b54e563f81604051611e209190617c8e565b60405180910390a150565b6060611e37838361522b565b905092915050565b611e47615d14565b80606760006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff167f5a85b4270fc1e75035e6cd505418ce65e0bcc36cc7eb9ce9e6f8c6181d4cb57760405160405180910390a250565b6000606760149054906101000a900460ff1615611f20576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611f1790618da9565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff1603611f8f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611f869061934b565b60405180910390fd5b6000611f9b848461522b565b9050600060405180608001604052808781526020018873ffffffffffffffffffffffffffffffffffffffff168152602001600081526020018381525090506000611fe58233611bcd565b9050600073ffffffffffffffffffffffffffffffffffffffff166076600083815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461208c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612083906193b7565b60405180910390fd5b81607660008381526020019081526020016000206000820151816000015560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160020155606082015181600301908161211091906195d9565b50905050807fa4b0af38d049ba81703a0d0e46cc2ff39681210302134046237111a8fb7dee7260405160405180910390a2809350505050949350505050565b606760149054906101000a900460ff161561219f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161219690618da9565b60405180910390fd5b6000607c600083815260200190815260200160002054146121f5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016121ec906196f7565b60405180910390fd5b6121fd615d92565b607c600083815260200190815260200160002081905550803373ffffffffffffffffffffffffffffffffffffffff167f09b4c028a2e50fec6f1c6a0163c59e8fbe92b231e5c03ef3adec585e63a14b9260405160405180910390a350565b606b5481565b606760149054906101000a900460ff16156122b1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016122a890618da9565b60405180910390fd5b6122bb6001615e3c565b6122c6838383616078565b505050565b606760149054906101000a900460ff1681565b606660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b606760149054906101000a900460ff1615612354576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161234b90618da9565b60405180910390fd5b61236084849050615e3c565b60005b848490508110156123c6576123b585858381811061238457612383619717565b5b9050602002013584848481811061239e5761239d619717565b5b90506020028101906123b09190619755565b616078565b806123bf90618f1a565b9050612363565b5050505050565b606760149054906101000a900460ff161561241d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161241490618da9565b60405180910390fd5b612426336154f0565b612465576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161245c906190ee565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16607d600083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160361250a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161250190619804565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16607e600083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146125af576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016125a690619870565b60405180910390fd5b607154607d600083815260200190815260200160002060000160149054906101000a900467ffffffffffffffff1667ffffffffffffffff166125f19190618cd6565b4210612632576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612629906198dc565b60405180910390fd5b607d6000828152602001908152602001600020600001601c9054906101000a900460ff1615612696576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161268d90619948565b60405180910390fd5b60006126a0611d90565b905060405180608001604052803373ffffffffffffffffffffffffffffffffffffffff1681526020014267ffffffffffffffff168152602001600063ffffffff16815260200182815250607e600084815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160000160146101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550604082015181600001601c6101000a81548163ffffffff021916908363ffffffff16021790555060608201518160010155905050813373ffffffffffffffffffffffffffffffffffffffff167f6958c989e915d3e41a35076e3c480363910055408055ad86ae1ee13d41c4064060405160405180910390a36127f782600133615590565b8060776000607d600086815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000154106128b9576128b8826000607d600086815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16615590565b5b5050565b60685481565b6128cb615d14565b6128d560006163e8565b565b606f5481565b607d6020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060000160149054906101000a900467ffffffffffffffff169080600001601c9054906101000a900460ff169080600101805461295390619406565b80601f016020809104026020016040519081016040528092919081815260200182805461297f90619406565b80156129cc5780601f106129a1576101008083540402835291602001916129cc565b820191906000526020600020905b8154815290600101906020018083116129af57829003601f168201915b5050505050905084565b606760149054906101000a900460ff1615612a26576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612a1d90618da9565b60405180910390fd5b6000607960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008481526020019081526020016000206040518060400160405290816000820154815260200160018201548152505090506000816000015111612ae0576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612ad7906199b4565b60405180910390fd5b8060000151421015612b27576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612b1e90619a20565b60405180910390fd5b8060200151607760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001541015612bb0576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612ba790619a8c565b60405180910390fd5b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb8383602001516040518363ffffffff1660e01b8152600401612c11929190619aac565b6020604051808303816000875af1158015612c30573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612c549190618eed565b508060200151607760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000016000828254612cab9190618f62565b92505081905550806020015160866000828254612cc89190618f62565b925050819055508060200151607a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254612d229190618f62565b92505081905550607960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600084815260200190815260200160002060008082016000905560018201600090555050828273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f109aeff667601aad33abcb7c5df1617754eefd18253d586a95c198cb479b5bd48460200151604051612dee9190617c8e565b60405180910390a4505050565b606760149054906101000a900460ff1615612e4b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612e4290618da9565b60405180910390fd5b612e53611c28565b607a6000607d600085815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205460776000607d600086815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000154612f4c9190618f62565b1015612f8d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612f8490619b21565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16607d600083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603613032576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161302990619b8d565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16607e600083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146130d7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016130ce90619bf9565b60405180910390fd5b607154426130e59190618f62565b607d600083815260200190815260200160002060000160149054906101000a900467ffffffffffffffff1667ffffffffffffffff161061315a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161315190619c65565b60405180910390fd5b60735460715460856000607d600086815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546131e09190618cd6565b6131ea9190618cd6565b607d600083815260200190815260200160002060000160149054906101000a900467ffffffffffffffff1667ffffffffffffffff161161325f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161325690619cf7565b60405180910390fd5b60001515607d6000838152602001908152602001600020600001601c9054906101000a900460ff161515146132c9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016132c090619d63565b60405180910390fd5b6001607d6000838152602001908152602001600020600001601c6101000a81548160ff02191690831515021790555080607d600083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f0b76b4ae356796814d36b46f7c500bbd27b2cce1e6059a6fa2bebfd5a389b19060405160405180910390a361337b816164ae565b50565b600080606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016133dc9190618543565b602060405180830381865afa1580156133f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061341d9190619d98565b9050697f0e10af47c1c7000000811061343a576000915050613474565b6000608c546086548361344d9190618f62565b6134579190618f62565b905080697f0e10af47c1c700000061346f9190618f62565b925050505b90565b60735481565b6004600060019054906101000a900460ff161580156134ae57508060ff1660008054906101000a900460ff1660ff16105b6134ed576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016134e490619e37565b60405180910390fd5b806000806101000a81548160ff021916908360ff1602179055506001600060016101000a81548160ff02191690831515021790555060008060016101000a81548160ff0219169083151502179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024988160405161356b9190619e57565b60405180910390a150565b608b60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008073ffffffffffffffffffffffffffffffffffffffff16607e600084815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160361361057600190506137d4565b61361982610f7d565b1561362757600290506137d4565b607f600083815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161561369357600390506137d4565b6000607760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010154036136e657600490506137d4565b607454607760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010154101561373b57600590506137d4565b607e600083815260200190815260200160002060000160149054906101000a900467ffffffffffffffff1667ffffffffffffffff16607454607760008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600101546137c09190618f62565b11156137cf57600690506137d4565b600090505b92915050565b607c6020528060005260406000206000915090505481565b606760149054906101000a900460ff1615613842576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161383990618da9565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16607e600084815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16036138e7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016138de90619ebe565b60405180910390fd5b6138f082610f7d565b61392f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161392690619f2a565b60405180910390fd5b60008163ffffffff1611613978576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161396f90619f96565b60405180910390fd5b6000608160008481526020019081526020016000208054905090506000608260008581526020019081526020016000208054905090506000607e6000868152602001908152602001600020600001601c9054906101000a900463ffffffff169050600084607e6000888152602001908152602001600020600001601c9054906101000a900463ffffffff16613a0d9190619fb6565b90506000607e60008881526020019081526020016000206001015490508363ffffffff168563ffffffff161115613fc4576000818563ffffffff16613a529190618d0a565b9050600060018763ffffffff1614613a8157600282613a719190619261565b82613a7c9190618f62565b613a83565b815b9050600060018863ffffffff1614613ac357600188613aa29190619fee565b63ffffffff168284613ab49190618f62565b613abe9190619261565b613ac6565b60005b905060008663ffffffff1690505b8563ffffffff16811015613d34578863ffffffff16811015613d21576000608160008d81526020019081526020016000208281548110613b1757613b16619717565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905085607760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000016000828254613b969190618cd6565b9250508190555060008203613c6457606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb82866040518363ffffffff1660e01b8152600401613c02929190619aac565b6020604051808303816000875af1158015613c21573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613c459190618eed565b508360866000828254613c589190618f62565b92505081905550613d1f565b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb82856040518363ffffffff1660e01b8152600401613cc1929190619aac565b6020604051808303816000875af1158015613ce0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613d049190618eed565b508260866000828254613d179190618f62565b925050819055505b505b8080613d2c90618f1a565b915050613ad4565b506000607e60008c8152602001908152602001600020600001601c9054906101000a900463ffffffff1663ffffffff1603613fbc57606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb607b60008d815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16607b60008e8152602001908152602001600020600101546040518363ffffffff1660e01b8152600401613e12929190619aac565b6020604051808303816000875af1158015613e31573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613e559190618eed565b50607b60008b81526020019081526020016000206001015460866000828254613e7e9190618f62565b92505081905550608360008b81526020019081526020016000205460776000608160008e8152602001908152602001600020600081548110613ec357613ec2619717565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000016000828254613f3a9190618cd6565b925050819055504260856000607d60008e815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055505b5050506142f2565b6000818663ffffffff16613fd89190618d0a565b9050600060018663ffffffff1614613ffc57600282613ff79190619261565b613ffe565b815b9050600060018763ffffffff161461403e5760018761401d9190619fee565b63ffffffff16828461402f9190618f62565b6140399190619261565b614041565b60005b905060008663ffffffff1690505b8563ffffffff168110156142af578763ffffffff1681101561429c576000608260008d8152602001908152602001600020828154811061409257614091619717565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905085607760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160008282546141119190618cd6565b92505081905550600082036141df57606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb82866040518363ffffffff1660e01b815260040161417d929190619aac565b6020604051808303816000875af115801561419c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906141c09190618eed565b5083608660008282546141d39190618f62565b9250508190555061429a565b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb82856040518363ffffffff1660e01b815260040161423c929190619aac565b6020604051808303816000875af115801561425b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061427f9190618eed565b5082608660008282546142929190618f62565b925050819055505b505b80806142a790618f1a565b91505061404f565b506000607e60008c8152602001908152602001600020600001601c9054906101000a900463ffffffff1663ffffffff16036142ee576142ed8a6164ae565b5b5050505b81607e6000898152602001908152602001600020600001601c6101000a81548163ffffffff021916908363ffffffff1602179055508263ffffffff16877f71d8c71303e35a39162e33a402c9897bf9848388537bac7d5e1b0d202eca4e668460405161435e919061a026565b60405180910390a350505050505050565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60745481565b60715481565b606760149054906101000a900460ff16156143f5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016143ec90618da9565b60405180910390fd5b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd3330846040518463ffffffff1660e01b815260040161445493929190618ea1565b6020604051808303816000875af1158015614473573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906144979190618eed565b5080608660008282546144aa9190618cd6565b9250508190555060006144bb611c28565b905080607760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000154116145a3578082607760008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001546145559190618cd6565b106145a25742607760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600101819055505b5b604051806060016040528083607760008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001546145fc9190618cd6565b8152602001607760008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206001015481526020018473ffffffffffffffffffffffffffffffffffffffff16815250607760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082015181600001556020820151816001015560408201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055509050508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8d4844488c19a90828439e71d14ebad860806d04f8ef8b25a82179fab2699b898460405161475d9190617c8e565b60405180910390a3505050565b614772615d14565b600073ffffffffffffffffffffffffffffffffffffffff166076600084815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603614817576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161480e90618e15565b60405180910390fd5b806076600084815260200190815260200160002060020181905550817f0321e8a918b4c47bf3677852c070983825f30a47bc8d9416691454fa6a727d63826040516148629190617c8e565b60405180910390a25050565b606c5481565b60845481565b606760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60875481565b608a5481565b60008082036148c557670de0b6b3a7640000905061492d565b670de0b6b3a7640000697f0e10af47c1c70000006148e38585611a09565b670de0b6b3a764000085697f0e10af47c1c70000006149029190618f62565b61490c9190618d0a565b6149169190618d0a565b6149209190619261565b61492a9190619261565b90505b92915050565b60755481565b614941615d14565b80606a60006101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055508067ffffffffffffffff167fa15d6b0930a82638ac4775bfd1b2e9f1e86be67e1bd3a09fa8a77a8f079769d760405160405180910390a250565b606760149054906101000a900460ff16156149f4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016149eb90618da9565b60405180910390fd5b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb606660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff166068546040518363ffffffff1660e01b8152600401614a75929190619aac565b6020604051808303816000875af1158015614a94573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614ab89190618eed565b5060685460866000828254614acd9190618f62565b925050819055506000606881905550565b60836020528060005260406000206000915090505481565b60695481565b60895481565b60886020528060005260406000206000915090505481565b60705481565b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b606760149054906101000a900460ff1615614b96576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614b8d90618da9565b60405180910390fd5b6000607960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008381526020019081526020016000206040518060400160405290816000820154815260200160018201548152505090506000816000015111614c50576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614c47906199b4565b60405180910390fd5b8060200151607a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254614ca39190618f62565b92505081905550607960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600083815260200190815260200160002060008082016000905560018201600090555050813373ffffffffffffffffffffffffffffffffffffffff167ff9d5cfcdfe9803069225971ba315f4302add9b477c3441cc363bc38fdb065b9060405160405180910390a35050565b600063bbf81e00821115614d7557697f0e10af47c1c70000009050614e10565b6000614dae614da9614da4614d8d6301e13380616de8565b614d9687616de8565b616df290919063ffffffff16565b616e27565b616eb9565b9050670de0b6b3a764000081670de0b6b3a764000080697f0e10af47c1c7000000614dd99190618d0a565b614de39190618d0a565b614ded9190619261565b614df79190619261565b697f0e10af47c1c7000000614e0c9190618f62565b9150505b919050565b60816020528160005260406000208181548110614e3157600080fd5b906000526020600020016000915091509054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60786020528060005260406000206000915090505481565b607f6020528160005260406000206020528060005260406000206000915091509054906101000a900460ff1681565b6079602052816000526040600020602052806000526040600020600091509150508060000154908060010154905082565b607e6020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060000160149054906101000a900467ffffffffffffffff169080600001601c9054906101000a900463ffffffff16908060010154905084565b614f57615d14565b80606660006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff167f6bdb9ceff405d990c9b60be9f719fbb80889d5f064e8fd76efd5bea353e9071260405160405180910390a250565b606d5481565b60766020528060005260406000206000915090508060000154908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169080600201549080600301805461503990619406565b80601f016020809104026020016040519081016040528092919081815260200182805461506590619406565b80156150b25780601f10615087576101008083540402835291602001916150b2565b820191906000526020600020905b81548152906001019060200180831161509557829003601f168201915b5050505050905084565b607b6020528060005260406000206000915090508060000154908060010154908060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020160149054906101000a900467ffffffffffffffff169080600201601c9054906101000a900460ff169080600301805461513e90619406565b80601f016020809104026020016040519081016040528092919081815260200182805461516a90619406565b80156151b75780601f1061518c576101008083540402835291602001916151b7565b820191906000526020600020905b81548152906001019060200180831161519a57829003601f168201915b5050505050905086565b606e5481565b6151cf615d14565b6151d881616ec3565b50565b60865481565b60776020528060005260406000206000915090508060000154908060010154908060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905083565b606062010000838390501115615276576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161526d9061a08d565b60405180910390fd5b600061528484849050616f46565b905060008185858460405160200161529f949392919061a15a565b604051602081830303815290604052905060026152bc8251616f46565b826040516020016152ce92919061a1d1565b6040516020818303038152906040526040516152ea919061a204565b602060405180830381855afa158015615307573d6000803e3d6000fd5b5050506040513d601f19601f8201168201806040525081019061532a919061a230565b60405160200161533a919061a2a4565b6040516020818303038152906040529250505092915050565b60006040518060c001604052808581526020018481526020018673ffffffffffffffffffffffffffffffffffffffff1681526020014267ffffffffffffffff1681526020018760ff16815260200183815250905060006153b682336069546117c4565b905081607b6000838152602001908152602001600020600082015181600001556020820151816001015560408201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060608201518160020160146101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550608082015181600201601c6101000a81548160ff021916908360ff16021790555060a082015181600301908161548c91906195d9565b509050503373ffffffffffffffffffffffffffffffffffffffff1685827fc3d3e0544c80e3bb83f62659259ae1574f72a91515ab3cae3dd75cf77e1b0aea876040516154d89190617c8e565b60405180910390a48060698190555050505050505050565b60006154fa611c28565b607a60008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054607760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001546155879190618f62565b10159050919050565b6001607f600085815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555081156156785760816000848152602001908152602001600020819080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506156ed565b60826000848152602001908152602001600020819080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505b607e600084815260200190815260200160002060010154607760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160008282546157559190618f62565b92505081905550828173ffffffffffffffffffffffffffffffffffffffff167f1aa9e4be46e24e1f2e7eeb1613c01629213cd42965d2716e18531b63e552e411846040516157a39190617d3b565b60405180910390a3505050565b6000819050919050565b6000806157c684615953565b905060006157d384615953565b90507f800000000000000000000000000000000000000000000000000000000000000082148061582257507f800000000000000000000000000000000000000000000000000000000000000081145b15615859576040517f9fe2b45000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000841261586a578361586f565b836000035b91506000831261587f5782615884565b826000035b9050600061589b83670de0b6b3a7640000846170f5565b90507f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8111156159045787876040517fd49c26b30000000000000000000000000000000000000000000000000000000081526004016158fb92919061a309565b60405180910390fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8587181390506159458161593e5782600003615940565b825b61720f565b965050505050505092915050565b6000819050919050565b600061598361596b83615953565b61597485615953565b61597e919061a332565b61720f565b905092915050565b60008061599784615953565b905060006159a484615953565b90507f80000000000000000000000000000000000000000000000000000000000000008214806159f357507f800000000000000000000000000000000000000000000000000000000000000081145b15615a2a576040517fa6070c2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008060008412615a3b5783615a40565b836000035b915060008312615a505782615a55565b826000035b90506000615a638383617219565b90507f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff811115615acc5787876040517f120b5b43000000000000000000000000000000000000000000000000000000008152600401615ac392919061a309565b60405180910390fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff858718139050615b0d81615b065782600003615b08565b825b61720f565b965050505050505092915050565b6000615b41615b2983615953565b615b3285615953565b615b3c919061a375565b61720f565b905092915050565b6000615b5482615953565b615b5d84615953565b12905092915050565b600080615b7283615953565b90507f80000000000000000000000000000000000000000000000000000000000000008103615bcd576040517fec2b9e6700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008112615bdb5782615bee565b615bed81615be89061a3b9565b61720f565b5b915050919050565b600080615c0283615953565b90506000811215615c89577ffffffffffffffffffffffffffffffffffffffffffffffffcc22e87f6eb468eeb811215615c3f576000915050615d0f565b615c82615c5e615c59615c548460000361720f565b615bf6565b615953565b6ec097ce7bc90715b34b9f100000000081615c7c57615c7b619232565b5b0561720f565b9150615d0d565b680a688906bd8affffff811315615cd757826040517f0360d028000000000000000000000000000000000000000000000000000000008152600401615cce919061a401565b60405180910390fd5b6000670de0b6b3a7640000604083901b81615cf557615cf4619232565b5b059050615d09615d0482617303565b61720f565b9250505b505b919050565b615d1c617c63565b73ffffffffffffffffffffffffffffffffffffffff16615d3a61436f565b73ffffffffffffffffffffffffffffffffffffffff1614615d90576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401615d879061a468565b60405180910390fd5b565b60008046905061a4ba811480615daa575062066eed81145b80615db7575062066eee81145b15615e3457606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015615e08573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190615e2c9190619d98565b915050615e39565b439150505b90565b607354607154608560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054615e8c9190618cd6565b615e969190618cd6565b4211615ed7576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401615ece9061a4fa565b60405180910390fd5b80608454615ee59190618d0a565b607760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000016000828254615f369190618f62565b92505081905550670de0b6b3a764000081608754615f549190618d0a565b615f5e9190619261565b608860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205442615fa99190618f62565b11615fe9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401615fe09061a566565b60405180910390fd5b42608860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550616036336154f0565b616075576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161606c906190ee565b60405180910390fd5b50565b6000801b607b600085815260200190815260200160002060000154036160d3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016160ca9061a5d2565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16607d600085815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614616178576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161616f9061a63e565b60405180910390fd5b600061618633858585611d14565b90506000607c600083815260200190815260200160002054116161de576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016161d59061a6aa565b60405180910390fd5b6161e6615d92565b607c6000838152602001908152602001600020541061623a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016162319061a716565b60405180910390fd5b60405180608001604052803373ffffffffffffffffffffffffffffffffffffffff1681526020014267ffffffffffffffff16815260200160001515815260200184848080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050815250607d600086815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160000160146101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550604082015181600001601c6101000a81548160ff021916908315150217905550606082015181600101908161638091906195d9565b509050506084546083600086815260200190815260200160002081905550833373ffffffffffffffffffffffffffffffffffffffff167f957c18b5af8413899ea8a576a4d3fb16839a02c9fccfdce098b6d59ef248525b60405160405180910390a350505050565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081603360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600060766000607b6000858152602001908152602001600020600001548152602001908152602001600020600001549050607b60008381526020019081526020016000206001015481111561650257600090505b60008111156165f957606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb60766000607b600087815260200190815260200160002060000154815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16836040518363ffffffff1660e01b81526004016165b4929190619aac565b6020604051808303816000875af11580156165d3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906165f79190618eed565b505b600081607b60008581526020019081526020016000206001015461661d9190618f62565b90506000670de0b6b3a7640000606e54670de0b6b3a76400006166409190618f62565b8361664b9190618d0a565b6166559190619261565b826166609190618f62565b905080606860008282546166749190618cd6565b92505081905550600081836166899190618f62565b9050600081111561677757606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb607d600088815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1684866167159190618f62565b6040518363ffffffff1660e01b8152600401616732929190619aac565b6020604051808303816000875af1158015616751573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906167759190618eed565b505b608b60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ebe2b12b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156167e4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906168089190619d98565b42111561696c57606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb608b60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16608c546040518363ffffffff1660e01b8152600401616890929190619aac565b6020604051808303816000875af11580156168af573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906168d39190618eed565b50608b60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633c6b16ab608c546040518263ffffffff1660e01b81526004016169319190617c8e565b600060405180830381600087803b15801561694b57600080fd5b505af115801561695f573d6000803e3d6000fd5b505050506000608c819055505b600060766000607b60008981526020019081526020016000206000015481526020019081526020016000206002015490506000811115616d04576000671bc16d674ec80000826169ba611d4d565b6169c49190618d0a565b6169ce9190619261565b905080608c60008282546169e29190618cd6565b925050819055506000811115616d02576000670de0b6b3a7640000607054670de0b6b3a7640000616a139190618f62565b83616a1e9190618d0a565b616a289190619261565b82616a339190618f62565b90506000670de0b6b3a7640000608954670de0b6b3a7640000616a569190618f62565b84616a619190618d0a565b616a6b9190619261565b83616a769190618f62565b9050606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb607d60008c815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16838587616afa9190618f62565b616b049190618f62565b6040518363ffffffff1660e01b8152600401616b21929190619aac565b6020604051808303816000875af1158015616b40573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190616b649190618eed565b50606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb606660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16846040518363ffffffff1660e01b8152600401616be4929190619aac565b6020604051808303816000875af1158015616c03573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190616c279190618eed565b50606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb607b60008c815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16836040518363ffffffff1660e01b8152600401616cbb929190619aac565b6020604051808303816000875af1158015616cda573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190616cfe9190618eed565b5050505b505b608360008781526020019081526020016000205460776000607d60008a815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000016000828254616d9f9190618cd6565b9250508190555082607b600088815260200190815260200160002060010154616dc89190618f62565b60866000828254616dd99190618f62565b92505081905550505050505050565b6000819050919050565b6000616e1f616e1a616e0385616eb9565b670de0b6b3a7640000616e1586616eb9565b6170f5565b617c6b565b905092915050565b600080616e3383616eb9565b9050680a688906bd8affffff811115616e8357826040517fb3b6ba1f000000000000000000000000000000000000000000000000000000008152600401616e7a919061a745565b60405180910390fd5b6000670de0b6b3a7640000604083901b616e9d9190619261565b9050616eb0616eab82617303565b617c6b565b92505050919050565b6000819050919050565b616ecb615d14565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603616f3a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401616f319061a7d2565b60405180910390fd5b616f43816163e8565b50565b606060008290506000600190505b607f8267ffffffffffffffff161115616f895760078267ffffffffffffffff16901c915080616f829061a7f2565b9050616f54565b60008167ffffffffffffffff1667ffffffffffffffff811115616faf57616fae618048565b5b6040519080825280601f01601f191660200182016040528015616fe15781602001600182028036833780820191505090505b50905084925060005b8267ffffffffffffffff168167ffffffffffffffff16101561708057607f841660801760f81b828267ffffffffffffffff168151811061702d5761702c619717565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535060078467ffffffffffffffff16901c935080806170789061a7f2565b915050616fea565b50607f60f81b81600184617094919061a822565b67ffffffffffffffff16815181106170af576170ae619717565b5b6020010181815160f81c60f81b169150907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350809350505050919050565b60008060008019858709858702925082811083820303915050600081036171305783828161712657617125619232565b5b0492505050617208565b838110617178578585856040517f63a0577800000000000000000000000000000000000000000000000000000000815260040161716f9392919061a85e565b60405180910390fd5b60008486880990508281118203915080830392506000600186190186169050600081870496508185049450600182836000030401905080840285179450600060028860030218905080880260020381029050808802600203810290508088026002038102905080880260020381029050808802600203810290508088026002038102905080860296505050505050505b9392505050565b6000819050919050565b600080600080198486098486029250828110838203039150506000810361725c57670de0b6b3a7640000828161725257617251619232565b5b04925050506172fd565b670de0b6b3a764000081106172aa5784846040517f5173648d0000000000000000000000000000000000000000000000000000000081526004016172a1929190618b30565b60405180910390fd5b6000670de0b6b3a764000085870990507faccb18165bd6fe31ae1cf318dc5b51eee0e1ba569b88cd74c1773b91fac106696001620400008060000304018483118403026204000083860304170293505050505b92915050565b6000778000000000000000000000000000000000000000000000009050600067ff00000000000000831611156174645760006780000000000000008316111561735957604068016a09e667f3bcc9098202901c90505b60006740000000000000008316111561737f5760406801306fe0a31b7152df8202901c90505b6000672000000000000000831611156173a55760406801172b83c7d517adce8202901c90505b6000671000000000000000831611156173cb57604068010b5586cf9890f62a8202901c90505b6000670800000000000000831611156173f15760406801059b0d31585743ae8202901c90505b600067040000000000000083161115617417576040680102c9a3e778060ee78202901c90505b60006702000000000000008316111561743d57604068010163da9fb33356d88202901c90505b600067010000000000000083161115617463576040680100b1afa5abcbed618202901c90505b5b600066ff0000000000008316111561759f57600066800000000000008316111561749b57604068010058c86da1c09ea28202901c90505b60006640000000000000831611156174c05760406801002c605e2e8cec508202901c90505b60006620000000000000831611156174e5576040680100162f3904051fa18202901c90505b600066100000000000008316111561750a5760406801000b175effdc76ba8202901c90505b600066080000000000008316111561752f576040680100058ba01fb9f96d8202901c90505b600066040000000000008316111561755457604068010002c5cc37da94928202901c90505b60006602000000000000831611156175795760406801000162e525ee05478202901c90505b600066010000000000008316111561759e57604068010000b17255775c048202901c90505b5b600065ff0000000000831611156176d157600065800000000000831611156175d45760406801000058b91b5bc9ae8202901c90505b600065400000000000831611156175f8576040680100002c5c89d5ec6d8202901c90505b6000652000000000008316111561761c57604068010000162e43f4f8318202901c90505b60006510000000000083161115617640576040680100000b1721bcfc9a8202901c90505b6000650800000000008316111561766457604068010000058b90cf1e6e8202901c90505b600065040000000000831611156176885760406801000002c5c863b73f8202901c90505b600065020000000000831611156176ac576040680100000162e430e5a28202901c90505b600065010000000000831611156176d05760406801000000b1721835518202901c90505b5b600064ff00000000831611156177fa57600064800000000083161115617704576040680100000058b90c0b498202901c90505b60006440000000008316111561772757604068010000002c5c8601cc8202901c90505b60006420000000008316111561774a5760406801000000162e42fff08202901c90505b60006410000000008316111561776d57604068010000000b17217fbb8202901c90505b6000640800000000831611156177905760406801000000058b90bfce8202901c90505b6000640400000000831611156177b3576040680100000002c5c85fe38202901c90505b6000640200000000831611156177d657604068010000000162e42ff18202901c90505b6000640100000000831611156177f9576040680100000000b17217f88202901c90505b5b600063ff0000008316111561791a57600063800000008316111561782b57604068010000000058b90bfc8202901c90505b600063400000008316111561784d5760406801000000002c5c85fe8202901c90505b600063200000008316111561786f576040680100000000162e42ff8202901c90505b60006310000000831611156178915760406801000000000b17217f8202901c90505b60006308000000831611156178b3576040680100000000058b90c08202901c90505b60006304000000831611156178d557604068010000000002c5c8608202901c90505b60006302000000831611156178f75760406801000000000162e4308202901c90505b600063010000008316111561791957604068010000000000b172188202901c90505b5b600062ff000083161115617a3157600062800000831611156179495760406801000000000058b90c8202901c90505b6000624000008316111561796a576040680100000000002c5c868202901c90505b6000622000008316111561798b57604068010000000000162e438202901c90505b600062100000831611156179ac576040680100000000000b17218202901c90505b600062080000831611156179cd57604068010000000000058b918202901c90505b600062040000831611156179ee5760406801000000000002c5c88202901c90505b60006202000083161115617a0f576040680100000000000162e48202901c90505b60006201000083161115617a305760406801000000000000b1728202901c90505b5b600061ff0083161115617b3f57600061800083161115617a5e576040680100000000000058b98202901c90505b600061400083161115617a7e57604068010000000000002c5d8202901c90505b600061200083161115617a9e5760406801000000000000162e8202901c90505b600061100083161115617abe57604068010000000000000b178202901c90505b600061080083161115617ade5760406801000000000000058c8202901c90505b600061040083161115617afe576040680100000000000002c68202901c90505b600061020083161115617b1e576040680100000000000001638202901c90505b600061010083161115617b3e576040680100000000000000b18202901c90505b5b600060ff83161115617c44576000608083161115617b6a576040680100000000000000598202901c90505b6000604083161115617b895760406801000000000000002c8202901c90505b6000602083161115617ba8576040680100000000000000168202901c90505b6000601083161115617bc75760406801000000000000000b8202901c90505b6000600883161115617be6576040680100000000000000068202901c90505b6000600483161115617c05576040680100000000000000038202901c90505b6000600283161115617c24576040680100000000000000018202901c90505b6000600183161115617c43576040680100000000000000018202901c90505b5b670de0b6b3a764000081029050604082901c60bf0381901c9050919050565b600033905090565b6000819050919050565b6000819050919050565b617c8881617c75565b82525050565b6000602082019050617ca36000830184617c7f565b92915050565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b617cd081617cbd565b8114617cdb57600080fd5b50565b600081359050617ced81617cc7565b92915050565b600060208284031215617d0957617d08617cb3565b5b6000617d1784828501617cde565b91505092915050565b60008115159050919050565b617d3581617d20565b82525050565b6000602082019050617d506000830184617d2c565b92915050565b600060ff82169050919050565b617d6c81617d56565b8114617d7757600080fd5b50565b600081359050617d8981617d63565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000617dba82617d8f565b9050919050565b617dca81617daf565b8114617dd557600080fd5b50565b600081359050617de781617dc1565b92915050565b617df681617c75565b8114617e0157600080fd5b50565b600081359050617e1381617ded565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f840112617e3e57617e3d617e19565b5b8235905067ffffffffffffffff811115617e5b57617e5a617e1e565b5b602083019150836001820283011115617e7757617e76617e23565b5b9250929050565b60008060008060008060a08789031215617e9b57617e9a617cb3565b5b6000617ea989828a01617d7a565b9650506020617eba89828a01617dd8565b9550506040617ecb89828a01617cde565b9450506060617edc89828a01617e04565b935050608087013567ffffffffffffffff811115617efd57617efc617cb8565b5b617f0989828a01617e28565b92509250509295509295509295565b600080600080600080600060c0888a031215617f3757617f36617cb3565b5b6000617f458a828b01617d7a565b9750506020617f568a828b01617dd8565b9650506040617f678a828b01617cde565b9550506060617f788a828b01617e04565b945050608088013567ffffffffffffffff811115617f9957617f98617cb8565b5b617fa58a828b01617e28565b935093505060a0617fb88a828b01617e04565b91505092959891949750929550565b600060208284031215617fdd57617fdc617cb3565b5b6000617feb84828501617e04565b91505092915050565b600067ffffffffffffffff82169050919050565b61801181617ff4565b82525050565b600060208201905061802c6000830184618008565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61808082618037565b810181811067ffffffffffffffff8211171561809f5761809e618048565b5b80604052505050565b60006180b2617ca9565b90506180be8282618077565b919050565b600080fd5b6180d181617ff4565b81146180dc57600080fd5b50565b6000813590506180ee816180c8565b92915050565b600080fd5b600067ffffffffffffffff82111561811457618113618048565b5b61811d82618037565b9050602081019050919050565b82818337600083830152505050565b600061814c618147846180f9565b6180a8565b905082815260208101848484011115618168576181676180f4565b5b61817384828561812a565b509392505050565b600082601f8301126181905761818f617e19565b5b81356181a0848260208601618139565b91505092915050565b600060c082840312156181bf576181be618032565b5b6181c960c06180a8565b905060006181d984828501617cde565b60008301525060206181ed84828501617e04565b602083015250604061820184828501617dd8565b6040830152506060618215848285016180df565b606083015250608061822984828501617d7a565b60808301525060a082013567ffffffffffffffff81111561824d5761824c6180c3565b5b6182598482850161817b565b60a08301525092915050565b60008060006060848603121561827e5761827d617cb3565b5b600084013567ffffffffffffffff81111561829c5761829b617cb8565b5b6182a8868287016181a9565b93505060206182b986828701617dd8565b92505060406182ca86828701617cde565b9150509250925092565b6182dd81617cbd565b82525050565b60006020820190506182f860008301846182d4565b92915050565b61830781617d20565b811461831257600080fd5b50565b600081359050618324816182fe565b92915050565b6000602082840312156183405761833f617cb3565b5b600061834e84828501618315565b91505092915050565b6000806040838503121561836e5761836d617cb3565b5b600061837c85828601617cde565b925050602061838d85828601618315565b9150509250929050565b6000602082840312156183ad576183ac617cb3565b5b60006183bb84828501617dd8565b91505092915050565b600080604083850312156183db576183da617cb3565b5b60006183e985828601617e04565b92505060206183fa85828601617e04565b9150509250929050565b60006080828403121561841a57618419618032565b5b61842460806180a8565b9050600061843484828501617e04565b600083015250602061844884828501617dd8565b602083015250604061845c84828501617e04565b604083015250606082013567ffffffffffffffff8111156184805761847f6180c3565b5b61848c8482850161817b565b60608301525092915050565b600080604083850312156184af576184ae617cb3565b5b600083013567ffffffffffffffff8111156184cd576184cc617cb8565b5b6184d985828601618404565b92505060206184ea85828601617dd8565b9150509250929050565b6000806040838503121561850b5761850a617cb3565b5b600061851985828601617cde565b925050602061852a85828601617e04565b9150509250929050565b61853d81617daf565b82525050565b60006020820190506185586000830184618534565b92915050565b6000806000806060858703121561857857618577617cb3565b5b600061858687828801617dd8565b945050602061859787828801617cde565b935050604085013567ffffffffffffffff8111156185b8576185b7617cb8565b5b6185c487828801617e28565b925092505092959194509250565b600080602083850312156185e9576185e8617cb3565b5b600083013567ffffffffffffffff81111561860757618606617cb8565b5b61861385828601617e28565b92509250509250929050565b600081519050919050565b600082825260208201905092915050565b60005b8381101561865957808201518184015260208101905061863e565b60008484015250505050565b60006186708261861f565b61867a818561862a565b935061868a81856020860161863b565b61869381618037565b840191505092915050565b600060208201905081810360008301526186b88184618665565b905092915050565b600080600080606085870312156186da576186d9617cb3565b5b60006186e887828801617dd8565b94505060206186f987828801617e04565b935050604085013567ffffffffffffffff81111561871a57618719617cb8565b5b61872687828801617e28565b925092505092959194509250565b60008060006040848603121561874d5761874c617cb3565b5b600061875b86828701617cde565b935050602084013567ffffffffffffffff81111561877c5761877b617cb8565b5b61878886828701617e28565b92509250509250925092565b60008083601f8401126187aa576187a9617e19565b5b8235905067ffffffffffffffff8111156187c7576187c6617e1e565b5b6020830191508360208202830111156187e3576187e2617e23565b5b9250929050565b60008083601f840112618800576187ff617e19565b5b8235905067ffffffffffffffff81111561881d5761881c617e1e565b5b60208301915083602082028301111561883957618838617e23565b5b9250929050565b6000806000806040858703121561885a57618859617cb3565b5b600085013567ffffffffffffffff81111561887857618877617cb8565b5b61888487828801618794565b9450945050602085013567ffffffffffffffff8111156188a7576188a6617cb8565b5b6188b3878288016187ea565b925092505092959194509250565b60006080820190506188d66000830187618534565b6188e36020830186618008565b6188f06040830185617d2c565b81810360608301526189028184618665565b905095945050505050565b6000806040838503121561892457618923617cb3565b5b600061893285828601617e04565b925050602061894385828601617dd8565b9150509250929050565b6000806040838503121561896457618963617cb3565b5b600061897285828601617dd8565b925050602061898385828601617cde565b9150509250929050565b600063ffffffff82169050919050565b6189a68161898d565b81146189b157600080fd5b50565b6000813590506189c38161899d565b92915050565b600080604083850312156189e0576189df617cb3565b5b60006189ee85828601617cde565b92505060206189ff858286016189b4565b9150509250929050565b60008060408385031215618a2057618a1f617cb3565b5b6000618a2e85828601617dd8565b9250506020618a3f85828601617e04565b9150509250929050565b600060208284031215618a5f57618a5e617cb3565b5b6000618a6d848285016180df565b91505092915050565b6000819050919050565b6000618a9b618a96618a9184617d8f565b618a76565b617d8f565b9050919050565b6000618aad82618a80565b9050919050565b6000618abf82618aa2565b9050919050565b618acf81618ab4565b82525050565b6000602082019050618aea6000830184618ac6565b92915050565b60008060408385031215618b0757618b06617cb3565b5b6000618b1585828601617cde565b9250506020618b2685828601617dd8565b9150509250929050565b6000604082019050618b456000830185617c7f565b618b526020830184617c7f565b9392505050565b618b628161898d565b82525050565b6000608082019050618b7d6000830187618534565b618b8a6020830186618008565b618b976040830185618b59565b618ba46060830184617c7f565b95945050505050565b6000608082019050618bc26000830187617c7f565b618bcf6020830186618534565b618bdc6040830185617c7f565b8181036060830152618bee8184618665565b905095945050505050565b618c0281617d56565b82525050565b600060c082019050618c1d60008301896182d4565b618c2a6020830188617c7f565b618c376040830187618534565b618c446060830186618008565b618c516080830185618bf9565b81810360a0830152618c638184618665565b9050979650505050505050565b6000606082019050618c856000830186617c7f565b618c926020830185617c7f565b618c9f6040830184618534565b949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000618ce182617c75565b9150618cec83617c75565b9250828201905080821115618d0457618d03618ca7565b5b92915050565b6000618d1582617c75565b9150618d2083617c75565b9250828202618d2e81617c75565b91508282048414831517618d4557618d44618ca7565b5b5092915050565b600082825260208201905092915050565b7f7061757365640000000000000000000000000000000000000000000000000000600082015250565b6000618d93600683618d4c565b9150618d9e82618d5d565b602082019050919050565b60006020820190508181036000830152618dc281618d86565b9050919050565b7f6d6f64656c20646f6573206e6f74206578697374000000000000000000000000600082015250565b6000618dff601483618d4c565b9150618e0a82618dc9565b602082019050919050565b60006020820190508181036000830152618e2e81618df2565b9050919050565b7f6c6f77657220666565207468616e206d6f64656c206665650000000000000000600082015250565b6000618e6b601883618d4c565b9150618e7682618e35565b602082019050919050565b60006020820190508181036000830152618e9a81618e5e565b9050919050565b6000606082019050618eb66000830186618534565b618ec36020830185618534565b618ed06040830184617c7f565b949350505050565b600081519050618ee7816182fe565b92915050565b600060208284031215618f0357618f02617cb3565b5b6000618f1184828501618ed8565b91505092915050565b6000618f2582617c75565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203618f5757618f56618ca7565b5b600182019050919050565b6000618f6d82617c75565b9150618f7883617c75565b9250828203905081811115618f9057618f8f618ca7565b5b92915050565b50565b6000618fa6600083618d4c565b9150618fb182618f96565b600082019050919050565b60006020820190508181036000830152618fd581618f99565b9050919050565b600060a082019050618ff16000830188618534565b618ffe60208301876182d4565b61900b60408301866182d4565b6190186060830185617c7f565b818103608083015261902a8184618665565b90509695505050505050565b7f6e6f742070617573657200000000000000000000000000000000000000000000600082015250565b600061906c600a83618d4c565b915061907782619036565b602082019050919050565b6000602082019050818103600083015261909b8161905f565b9050919050565b7f6d696e207374616b656420746f6f206c6f770000000000000000000000000000600082015250565b60006190d8601283618d4c565b91506190e3826190a2565b602082019050919050565b60006020820190508181036000830152619107816190cb565b9050919050565b7f6e6f7420616c6c6f776564000000000000000000000000000000000000000000600082015250565b6000619144600b83618d4c565b915061914f8261910e565b602082019050919050565b6000602082019050818103600083015261917381619137565b9050919050565b7f6d696e2076616c73000000000000000000000000000000000000000000000000600082015250565b60006191b0600883618d4c565b91506191bb8261917a565b602082019050919050565b600060208201905081810360008301526191df816191a3565b9050919050565b60006080820190506191fb6000830187618534565b6192086020830186618534565b6192156040830185617c7f565b81810360608301526192278184618665565b905095945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600061926c82617c75565b915061927783617c75565b92508261928757619286619232565b5b828204905092915050565b600061929e838561862a565b93506192ab83858461812a565b6192b483618037565b840190509392505050565b60006060820190506192d46000830187618534565b6192e160208301866182d4565b81810360408301526192f4818486619292565b905095945050505050565b7f61646472657373206d757374206265206e6f6e2d7a65726f0000000000000000600082015250565b6000619335601883618d4c565b9150619340826192ff565b602082019050919050565b6000602082019050818103600083015261936481619328565b9050919050565b7f6d6f64656c20616c726561647920726567697374657265640000000000000000600082015250565b60006193a1601883618d4c565b91506193ac8261936b565b602082019050919050565b600060208201905081810360008301526193d081619394565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061941e57607f821691505b602082108103619431576194306193d7565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026194997fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8261945c565b6194a3868361945c565b95508019841693508086168417925050509392505050565b60006194d66194d16194cc84617c75565b618a76565b617c75565b9050919050565b6000819050919050565b6194f0836194bb565b6195046194fc826194dd565b848454619469565b825550505050565b600090565b61951961950c565b6195248184846194e7565b505050565b5b818110156195485761953d600082619511565b60018101905061952a565b5050565b601f82111561958d5761955e81619437565b6195678461944c565b81016020851015619576578190505b61958a6195828561944c565b830182619529565b50505b505050565b600082821c905092915050565b60006195b060001984600802619592565b1980831691505092915050565b60006195c9838361959f565b9150826002028217905092915050565b6195e28261861f565b67ffffffffffffffff8111156195fb576195fa618048565b5b6196058254619406565b61961082828561954c565b600060209050601f8311600181146196435760008415619631578287015190505b61963b85826195bd565b8655506196a3565b601f19841661965186619437565b60005b8281101561967957848901518255600182019150602085019450602081019050619654565b868310156196965784890151619692601f89168261959f565b8355505b6001600288020188555050505b505050505050565b7f636f6d6d69746d656e7420657869737473000000000000000000000000000000600082015250565b60006196e1601183618d4c565b91506196ec826196ab565b602082019050919050565b60006020820190508181036000830152619710816196d4565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600080fd5b600080fd5b600080fd5b6000808335600160200384360303811261977257619771619746565b5b80840192508235915067ffffffffffffffff8211156197945761979361974b565b5b6020830192506001820236038313156197b0576197af619750565b5b509250929050565b7f736f6c7574696f6e20646f6573206e6f74206578697374000000000000000000600082015250565b60006197ee601783618d4c565b91506197f9826197b8565b602082019050919050565b6000602082019050818103600083015261981d816197e1565b9050919050565b7f636f6e746573746174696f6e20616c7265616479206578697374730000000000600082015250565b600061985a601b83618d4c565b915061986582619824565b602082019050919050565b600060208201905081810360008301526198898161984d565b9050919050565b7f746f6f206c617465000000000000000000000000000000000000000000000000600082015250565b60006198c6600883618d4c565b91506198d182619890565b602082019050919050565b600060208201905081810360008301526198f5816198b9565b9050919050565b7f7774660000000000000000000000000000000000000000000000000000000000600082015250565b6000619932600383618d4c565b915061993d826198fc565b602082019050919050565b6000602082019050818103600083015261996181619925565b9050919050565b7f72657175657374206e6f74206578697374000000000000000000000000000000600082015250565b600061999e601183618d4c565b91506199a982619968565b602082019050919050565b600060208201905081810360008301526199cd81619991565b9050919050565b7f77616974206c6f6e676572000000000000000000000000000000000000000000600082015250565b6000619a0a600b83618d4c565b9150619a15826199d4565b602082019050919050565b60006020820190508181036000830152619a39816199fd565b9050919050565b7f7374616b6520696e73756666696369656e740000000000000000000000000000600082015250565b6000619a76601283618d4c565b9150619a8182619a40565b602082019050919050565b60006020820190508181036000830152619aa581619a69565b9050919050565b6000604082019050619ac16000830185618534565b619ace6020830184617c7f565b9392505050565b7f76616c696461746f72206d696e207374616b656420746f6f206c6f7700000000600082015250565b6000619b0b601c83618d4c565b9150619b1682619ad5565b602082019050919050565b60006020820190508181036000830152619b3a81619afe565b9050919050565b7f736f6c7574696f6e206e6f7420666f756e640000000000000000000000000000600082015250565b6000619b77601283618d4c565b9150619b8282619b41565b602082019050919050565b60006020820190508181036000830152619ba681619b6a565b9050919050565b7f68617320636f6e746573746174696f6e00000000000000000000000000000000600082015250565b6000619be3601083618d4c565b9150619bee82619bad565b602082019050919050565b60006020820190508181036000830152619c1281619bd6565b9050919050565b7f6e6f7420656e6f7567682064656c617900000000000000000000000000000000600082015250565b6000619c4f601083618d4c565b9150619c5a82619c19565b602082019050919050565b60006020820190508181036000830152619c7e81619c42565b9050919050565b7f636c61696d536f6c7574696f6e20636f6f6c646f776e206166746572206c6f7360008201527f7420636f6e746573746174696f6e000000000000000000000000000000000000602082015250565b6000619ce1602e83618d4c565b9150619cec82619c85565b604082019050919050565b60006020820190508181036000830152619d1081619cd4565b9050919050565b7f616c726561647920636c61696d65640000000000000000000000000000000000600082015250565b6000619d4d600f83618d4c565b9150619d5882619d17565b602082019050919050565b60006020820190508181036000830152619d7c81619d40565b9050919050565b600081519050619d9281617ded565b92915050565b600060208284031215619dae57619dad617cb3565b5b6000619dbc84828501619d83565b91505092915050565b7f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160008201527f647920696e697469616c697a6564000000000000000000000000000000000000602082015250565b6000619e21602e83618d4c565b9150619e2c82619dc5565b604082019050919050565b60006020820190508181036000830152619e5081619e14565b9050919050565b6000602082019050619e6c6000830184618bf9565b92915050565b7f636f6e746573746174696f6e20646f65736e2774206578697374000000000000600082015250565b6000619ea8601a83618d4c565b9150619eb382619e72565b602082019050919050565b60006020820190508181036000830152619ed781619e9b565b9050919050565b7f766f74696e6720706572696f64206e6f7420656e646564000000000000000000600082015250565b6000619f14601783618d4c565b9150619f1f82619ede565b602082019050919050565b60006020820190508181036000830152619f4381619f07565b9050919050565b7f616d6e7420746f6f20736d616c6c000000000000000000000000000000000000600082015250565b6000619f80600e83618d4c565b9150619f8b82619f4a565b602082019050919050565b60006020820190508181036000830152619faf81619f73565b9050919050565b6000619fc18261898d565b9150619fcc8361898d565b9250828201905063ffffffff811115619fe857619fe7618ca7565b5b92915050565b6000619ff98261898d565b915061a0048361898d565b9250828203905063ffffffff81111561a0205761a01f618ca7565b5b92915050565b600060208201905061a03b6000830184618b59565b92915050565b7f4d617820636f6e74656e742073697a6520697320363535333620627974657300600082015250565b600061a077601f83618d4c565b915061a0828261a041565b602082019050919050565b6000602082019050818103600083015261a0a68161a06a565b9050919050565b7f0802120000000000000000000000000000000000000000000000000000000000815250565b600081905092915050565b600061a0e98261861f565b61a0f3818561a0d3565b935061a10381856020860161863b565b80840191505092915050565b600061a11b838561a0d3565b935061a12883858461812a565b82840190509392505050565b7f1800000000000000000000000000000000000000000000000000000000000000815250565b600061a1658261a0ad565b60038201915061a175828761a0de565b915061a18282858761a10f565b915061a18d8261a134565b60018201915061a19d828461a0de565b915081905095945050505050565b7f0a00000000000000000000000000000000000000000000000000000000000000815250565b600061a1dc8261a1ab565b60018201915061a1ec828561a0de565b915061a1f8828461a0de565b91508190509392505050565b600061a210828461a0de565b915081905092915050565b60008151905061a22a81617cc7565b92915050565b60006020828403121561a2465761a245617cb3565b5b600061a2548482850161a21b565b91505092915050565b7f1220000000000000000000000000000000000000000000000000000000000000815250565b6000819050919050565b61a29e61a29982617cbd565b61a283565b82525050565b600061a2af8261a25d565b60028201915061a2bf828461a28d565b60208201915081905092915050565b6000819050919050565b600061a2f361a2ee61a2e98461a2ce565b618a76565b61a2ce565b9050919050565b61a3038161a2d8565b82525050565b600060408201905061a31e600083018561a2fa565b61a32b602083018461a2fa565b9392505050565b600061a33d8261a2ce565b915061a3488361a2ce565b925082820390508181126000841216828213600085121516171561a36f5761a36e618ca7565b5b92915050565b600061a3808261a2ce565b915061a38b8361a2ce565b92508282019050828112156000831216838212600084121516171561a3b35761a3b2618ca7565b5b92915050565b600061a3c48261a2ce565b91507f8000000000000000000000000000000000000000000000000000000000000000820361a3f65761a3f5618ca7565b5b816000039050919050565b600060208201905061a416600083018461a2fa565b92915050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b600061a452602083618d4c565b915061a45d8261a41c565b602082019050919050565b6000602082019050818103600083015261a4818161a445565b9050919050565b7f7375626d6974536f6c7574696f6e20636f6f6c646f776e206166746572206c6f60008201527f737420636f6e746573746174696f6e0000000000000000000000000000000000602082015250565b600061a4e4602f83618d4c565b915061a4ef8261a488565b604082019050919050565b6000602082019050818103600083015261a5138161a4d7565b9050919050565b7f736f6c7574696f6e2072617465206c696d697400000000000000000000000000600082015250565b600061a550601383618d4c565b915061a55b8261a51a565b602082019050919050565b6000602082019050818103600083015261a57f8161a543565b9050919050565b7f7461736b20646f6573206e6f7420657869737400000000000000000000000000600082015250565b600061a5bc601383618d4c565b915061a5c78261a586565b602082019050919050565b6000602082019050818103600083015261a5eb8161a5af565b9050919050565b7f736f6c7574696f6e20616c7265616479207375626d6974746564000000000000600082015250565b600061a628601a83618d4c565b915061a6338261a5f2565b602082019050919050565b6000602082019050818103600083015261a6578161a61b565b9050919050565b7f6e6f6e206578697374656e7420636f6d6d69746d656e74000000000000000000600082015250565b600061a694601783618d4c565b915061a69f8261a65e565b602082019050919050565b6000602082019050818103600083015261a6c38161a687565b9050919050565b7f636f6d6d69746d656e74206d75737420626520696e2070617374000000000000600082015250565b600061a700601a83618d4c565b915061a70b8261a6ca565b602082019050919050565b6000602082019050818103600083015261a72f8161a6f3565b9050919050565b61a73f816194bb565b82525050565b600060208201905061a75a600083018461a736565b92915050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b600061a7bc602683618d4c565b915061a7c78261a760565b604082019050919050565b6000602082019050818103600083015261a7eb8161a7af565b9050919050565b600061a7fd82617ff4565b915067ffffffffffffffff820361a8175761a816618ca7565b5b600182019050919050565b600061a82d82617ff4565b915061a83883617ff4565b9250828203905067ffffffffffffffff81111561a8585761a857618ca7565b5b92915050565b600060608201905061a8736000830186617c7f565b61a8806020830185617c7f565b61a88d6040830184617c7f565b94935050505056fea26469706673582212202c2b1a3db24234d3807a3f6916dadd58e3fb7a86dda961fe9dc9935a1c97baab64736f6c63430008130033",
}

// EngineABI is the input ABI used to generate the binding from.
// Deprecated: Use EngineMetaData.ABI instead.
var EngineABI = EngineMetaData.ABI

// EngineBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EngineMetaData.Bin instead.
var EngineBin = EngineMetaData.Bin

// DeployEngine deploys a new Ethereum contract, binding an instance of Engine to it.
func DeployEngine(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Engine, error) {
	parsed, err := EngineMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EngineBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Engine{EngineCaller: EngineCaller{contract: contract}, EngineTransactor: EngineTransactor{contract: contract}, EngineFilterer: EngineFilterer{contract: contract}}, nil
}

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
