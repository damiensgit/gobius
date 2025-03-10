// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package enginev5

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

// EngineV5MetaData contains all meta data concerning the EngineV5 contract.
var EngineV5MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv18_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"denominator\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv_Overflow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PRBMath_SD59x18_Abs_MinSD59x18\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PRBMath_SD59x18_Div_InputTooSmall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"SD59x18\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"SD59x18\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"PRBMath_SD59x18_Div_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"SD59x18\",\"name\":\"x\",\"type\":\"int256\"}],\"name\":\"PRBMath_SD59x18_Exp2_InputTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PRBMath_SD59x18_Mul_InputTooSmall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"SD59x18\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"SD59x18\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"PRBMath_SD59x18_Mul_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"UD60x18\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"PRBMath_UD60x18_Exp2_InputTooBig\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"task\",\"type\":\"bytes32\"}],\"name\":\"ContestationSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"task\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"yea\",\"type\":\"bool\"}],\"name\":\"ContestationVote\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"start_idx\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"end_idx\",\"type\":\"uint32\"}],\"name\":\"ContestationVoteFinish\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"ModelRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"name\":\"PausedChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"PauserTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"}],\"name\":\"SignalCommitment\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"task\",\"type\":\"bytes32\"}],\"name\":\"SolutionClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"SolutionMineableRateChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"task\",\"type\":\"bytes32\"}],\"name\":\"SolutionSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"startBlockTime\",\"type\":\"uint64\"}],\"name\":\"StartBlockTimeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"model\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"TaskSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"TreasuryTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ValidatorDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ValidatorWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"ValidatorWithdrawCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unlockTime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ValidatorWithdrawInitiated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"version\",\"type\":\"uint256\"}],\"name\":\"VersionChanged\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"accruedFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"baseToken\",\"outputs\":[{\"internalType\":\"contractIBaseToken\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"taskids_\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes[]\",\"name\":\"cids_\",\"type\":\"bytes[]\"}],\"name\":\"bulkSubmitSolution\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"version_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"owner_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"model_\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"input_\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"n_\",\"type\":\"uint256\"}],\"name\":\"bulkSubmitTask\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"count_\",\"type\":\"uint256\"}],\"name\":\"cancelValidatorWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"}],\"name\":\"claimSolution\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"commitments\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"contestationVoteExtensionTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"amnt_\",\"type\":\"uint32\"}],\"name\":\"contestationVoteFinish\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"contestationVoteNays\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"contestationVoteYeas\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"contestationVoted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"contestationVotedIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"contestations\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"blocktime\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"finish_start_index\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"slashAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"t\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ts\",\"type\":\"uint256\"}],\"name\":\"diffMul\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"exitValidatorMinUnlockTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"cid_\",\"type\":\"bytes\"}],\"name\":\"generateCommitment\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"content_\",\"type\":\"bytes\"}],\"name\":\"generateIPFSCID\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPsuedoTotalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSlashAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidatorMinimum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cid\",\"type\":\"bytes\"}],\"internalType\":\"structModel\",\"name\":\"o_\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender_\",\"type\":\"address\"}],\"name\":\"hashModel\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"model\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"blocktime\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"cid\",\"type\":\"bytes\"}],\"internalType\":\"structTask\",\"name\":\"o_\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sender_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"prevhash_\",\"type\":\"bytes32\"}],\"name\":\"hashTask\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"}],\"name\":\"initiateValidatorWithdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"lastContestationLossTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"lastSolutionSubmission\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxContestationValidatorStakeSince\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minClaimSolutionTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minContestationVotePeriodTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minRetractionWaitTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"models\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"cid\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauser\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"pendingValidatorWithdrawRequests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"unlockTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"pendingValidatorWithdrawRequestsCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"prevhash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"template_\",\"type\":\"bytes\"}],\"name\":\"registerModel\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"retractionFeePercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"t\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ts\",\"type\":\"uint256\"}],\"name\":\"reward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"paused_\",\"type\":\"bool\"}],\"name\":\"setPaused\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"model_\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"rate_\",\"type\":\"uint256\"}],\"name\":\"setSolutionMineableRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"startBlockTime_\",\"type\":\"uint64\"}],\"name\":\"setStartBlockTime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"veStaking_\",\"type\":\"address\"}],\"name\":\"setVeStaking\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"version_\",\"type\":\"uint256\"}],\"name\":\"setVersion\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"voter_\",\"type\":\"address\"}],\"name\":\"setVoter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitment_\",\"type\":\"bytes32\"}],\"name\":\"signalCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"slashAmountPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"solutionFeePercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"solutionRateLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"solutions\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"blocktime\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"claimed\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"cid\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"solutionsStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"solutionsStakeAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"startBlockTime\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"}],\"name\":\"submitContestation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"cid_\",\"type\":\"bytes\"}],\"name\":\"submitSolution\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"version_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"owner_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"model_\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"input_\",\"type\":\"bytes\"}],\"name\":\"submitTask\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"t\",\"type\":\"uint256\"}],\"name\":\"targetTs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"taskOwnerRewardPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"tasks\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"model\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"blocktime\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"cid\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalHeld\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to_\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to_\",\"type\":\"address\"}],\"name\":\"transferPauser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to_\",\"type\":\"address\"}],\"name\":\"transferTreasury\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"treasury\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"treasuryRewardPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"}],\"name\":\"validatorCanVote\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"}],\"name\":\"validatorDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validatorMinimumPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"count_\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to_\",\"type\":\"address\"}],\"name\":\"validatorWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validatorWithdrawPendingAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"staked\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"since\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"veRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"veStaking\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"yea_\",\"type\":\"bool\"}],\"name\":\"voteOnContestation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"voter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"}],\"name\":\"votingPeriodEnded\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawAccruedFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50620000226200002860201b60201c565b620001d2565b600060019054906101000a900460ff16156200007b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620000729062000175565b60405180910390fd5b60ff801660008054906101000a900460ff1660ff1614620000ec5760ff6000806101000a81548160ff021916908360ff1602179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb384740249860ff604051620000e39190620001b5565b60405180910390a15b565b600082825260208201905092915050565b7f496e697469616c697a61626c653a20636f6e747261637420697320696e69746960008201527f616c697a696e6700000000000000000000000000000000000000000000000000602082015250565b60006200015d602783620000ee565b91506200016a82620000ff565b604082019050919050565b6000602082019050818103600083015262000190816200014e565b9050919050565b600060ff82169050919050565b620001af8162000197565b82525050565b6000602082019050620001cc6000830184620001a4565b92915050565b61ab3480620001e26000396000f3fe608060405234801561001057600080fd5b50600436106104955760003560e01c80637b36006a11610262578063b4dc35b711610151578063d2780940116100ce578063e236f46b11610092578063e236f46b14610ed5578063e579f50014610f08578063f1b8989d14610f3d578063f2fde38b14610f5b578063f43cc77314610f77578063fa52c7d814610f9557610495565b8063d278094014610e07578063d2992baa14610e37578063d33b2ef514610e68578063d8a6021c14610e9b578063dc06a89f14610eb757610495565b8063c55dae6311610115578063c55dae6314610d3d578063cbd2422d14610d5b578063cf596e4514610d77578063d1f0c94114610da7578063d2307ae414610dd757610495565b8063b4dc35b714610c83578063c17ddb2a14610cb3578063c1f7272314610cd1578063c24b563114610cef578063c31784be14610d1f57610495565b806393f1f8ac116101df578063a2492a90116101a3578063a2492a9014610bf1578063a4fa8d5714610c0f578063a53e252514610c3f578063a8f837f314610c5d578063ada82c7d14610c7957610495565b806393f1f8ac14610b5d57806396bb02c314610b795780639b97511914610b975780639fd0506d14610bb5578063a1975adf14610bd357610495565b80638b4d7b35116102265780638b4d7b3514610acb5780638da5cb5b14610ae75780638e6d86fd14610b055780639280944414610b2357806393a090ec14610b4157610495565b80637b36006a14610a255780638129fc1c14610a4357806382b5077f14610a4d5780638365779514610a6b578063839df94514610a9b57610495565b80633d57f5d9116103895780635c975abb11610306578063715018a6116102ca578063715018a61461097457806372dc0ee11461097e57806375c705091461099c578063763253bb146109cf57806377286d17146109eb5780637881c5e614610a0757610495565b80635c975abb146108e257806361d027b31461090057806365d445fb1461091e578063671f81521461093a578063682c20581461095657610495565b80634bc2a6571161034d5780634bc2a657146108405780634ff03efa1461085c578063506ea7de1461088c57806354fd4d50146108a857806356914caf146108c657610495565b80633d57f5d91461079c578063408def1e146107ba57806340e8c56d146107d65780634421ea211461080657806346c96aac1461082257610495565b80631825c20e116104175780632258d105116103db5780632258d105146106e45780632943a49014610702578063303fb0d61461071e578063393cb1c71461074e5780633d18b9121461077e57610495565b80631825c20e146106085780631b75c43e146106245780631f88ea1c14610654578063218a304814610684578063218e6859146106b457610495565b80630c18d4ce1161045e5780630c18d4ce146105505780630d468d951461056e5780631466b63a1461058c57806316c38b3c146105bc57806317f3e041146105d857610495565b8062fd70821461049a57806305d1bc26146104b857806308745dd1146104e857806308afe0eb146105045780630a98573714610520575b600080fd5b6104a2610fc7565b6040516104af9190617ef7565b60405180910390f35b6104d260048036038101906104cd9190617f5c565b610fcd565b6040516104df9190617fa4565b60405180910390f35b61050260048036038101906104fd91906180e7565b611069565b005b61051e60048036038101906105199190618181565b611297565b005b61053a60048036038101906105359190618230565b6114f9565b6040516105479190617ef7565b60405180910390f35b6105586117f4565b6040516105659190618280565b60405180910390f35b61057661180e565b6040516105839190617ef7565b60405180910390f35b6105a660048036038101906105a191906184ce565b611814565b6040516105b3919061854c565b60405180910390f35b6105d660048036038101906105d19190618593565b61185a565b005b6105f260048036038101906105ed9190617f5c565b611936565b6040516105ff9190617ef7565b60405180910390f35b610622600480360381019061061d91906185c0565b61194e565b005b61063e60048036038101906106399190618600565b611a41565b60405161064b9190617ef7565b60405180910390f35b61066e6004803603810190610669919061862d565b611a59565b60405161067b9190617ef7565b60405180910390f35b61069e60048036038101906106999190618701565b611c1d565b6040516106ab919061854c565b60405180910390f35b6106ce60048036038101906106c99190618600565b611c60565b6040516106db9190617ef7565b60405180910390f35b6106ec611c78565b6040516106f99190617ef7565b60405180910390f35b61071c60048036038101906107179190618600565b611cca565b005b6107386004803603810190610733919061875d565b611d16565b60405161074591906187ac565b60405180910390f35b610768600480360381019061076391906187c7565b611d64565b604051610775919061854c565b60405180910390f35b610786611d9d565b6040516107939190617ef7565b60405180910390f35b6107a4611de0565b6040516107b19190617ef7565b60405180910390f35b6107d460048036038101906107cf9190618230565b611e32565b005b6107f060048036038101906107eb919061883b565b611e7b565b6040516107fd9190618907565b60405180910390f35b610820600480360381019061081b9190618600565b611e8f565b005b61082a611f1e565b60405161083791906187ac565b60405180910390f35b61085a60048036038101906108559190618600565b611f44565b005b61087660048036038101906108719190618929565b611f90565b604051610883919061854c565b60405180910390f35b6108a660048036038101906108a19190617f5c565b612211565b005b6108b061231d565b6040516108bd9190617ef7565b60405180910390f35b6108e060048036038101906108db919061899d565b612323565b005b6108ea61238d565b6040516108f79190617fa4565b60405180910390f35b6109086123a0565b60405161091591906187ac565b60405180910390f35b61093860048036038101906109339190618aa9565b6123c6565b005b610954600480360381019061094f9190617f5c565b61248f565b005b61095e61297f565b60405161096b9190617ef7565b60405180910390f35b61097c612985565b005b610986612999565b6040516109939190617ef7565b60405180910390f35b6109b660048036038101906109b19190617f5c565b61299f565b6040516109c69493929190618b2a565b60405180910390f35b6109e960048036038101906109e49190618b76565b612a98565b005b610a056004803603810190610a009190617f5c565b612ebd565b005b610a0f613440565b604051610a1c9190617ef7565b60405180910390f35b610a2d613539565b604051610a3a9190617ef7565b60405180910390f35b610a4b61353f565b005b610a55613638565b604051610a6291906187ac565b60405180910390f35b610a856004803603810190610a809190618bb6565b61365e565b604051610a929190617ef7565b60405180910390f35b610ab56004803603810190610ab09190617f5c565b61389c565b604051610ac29190617ef7565b60405180910390f35b610ae56004803603810190610ae09190618c32565b6138b4565b005b610aef614431565b604051610afc91906187ac565b60405180910390f35b610b0d61445b565b604051610b1a9190617ef7565b60405180910390f35b610b2b614461565b604051610b389190617ef7565b60405180910390f35b610b5b6004803603810190610b569190618c72565b614467565b005b610b776004803603810190610b72919061875d565b61482c565b005b610b81614930565b604051610b8e9190617ef7565b60405180910390f35b610b9f614936565b604051610bac9190617ef7565b60405180910390f35b610bbd61493c565b604051610bca91906187ac565b60405180910390f35b610bdb614962565b604051610be89190617ef7565b60405180910390f35b610bf9614968565b604051610c069190617ef7565b60405180910390f35b610c296004803603810190610c24919061862d565b61496e565b604051610c369190617ef7565b60405180910390f35b610c476149f5565b604051610c549190617ef7565b60405180910390f35b610c776004803603810190610c729190618cb2565b6149fb565b005b610c81614a66565b005b610c9d6004803603810190610c989190617f5c565b614ba0565b604051610caa9190617ef7565b60405180910390f35b610cbb614bb8565b604051610cc8919061854c565b60405180910390f35b610cd9614bbe565b604051610ce69190617ef7565b60405180910390f35b610d096004803603810190610d049190618600565b614bc4565b604051610d169190617ef7565b60405180910390f35b610d27614bdc565b604051610d349190617ef7565b60405180910390f35b610d45614be2565b604051610d529190618d3e565b60405180910390f35b610d756004803603810190610d709190618230565b614c08565b005b610d916004803603810190610d8c9190618230565b614e17565b604051610d9e9190617ef7565b60405180910390f35b610dc16004803603810190610dbc919061875d565b614ed7565b604051610dce91906187ac565b60405180910390f35b610df16004803603810190610dec9190618600565b614f25565b604051610dfe9190617ef7565b60405180910390f35b610e216004803603810190610e1c9190618d59565b614f3d565b604051610e2e9190617fa4565b60405180910390f35b610e516004803603810190610e4c9190618c72565b614f6c565b604051610e5f929190618d99565b60405180910390f35b610e826004803603810190610e7d9190617f5c565b614f9d565b604051610e929493929190618dd1565b60405180910390f35b610eb56004803603810190610eb09190618600565b615011565b005b610ebf6150a0565b604051610ecc9190617ef7565b60405180910390f35b610eef6004803603810190610eea9190617f5c565b6150a6565b604051610eff9493929190618e16565b60405180910390f35b610f226004803603810190610f1d9190617f5c565b61517e565b604051610f3496959493929190618e71565b60405180910390f35b610f45615283565b604051610f529190617ef7565b60405180910390f35b610f756004803603810190610f709190618600565b615289565b005b610f7f61529d565b604051610f8c9190617ef7565b60405180910390f35b610faf6004803603810190610faa9190618600565b6152a3565b604051610fbe93929190618ed9565b60405180910390f35b60725481565b6000608a546082600084815260200190815260200160002080549050608160008581526020019081526020016000208054905061100a9190618f3f565b6110149190618f73565b607354607e600085815260200190815260200160002060000160149054906101000a900467ffffffffffffffff1667ffffffffffffffff166110569190618f3f565b6110609190618f3f565b42119050919050565b606760149054906101000a900460ff16156110b9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110b090619012565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff166076600086815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160361115e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016111559061907e565b60405180910390fd5b60766000858152602001908152602001600020600001548310156111b7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016111ae906190ea565b60405180910390fd5b60006111c383836152ed565b90506111d28787878785615415565b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd3330876040518463ffffffff1660e01b81526004016112319392919061910a565b6020604051808303816000875af1158015611250573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112749190619156565b5083608660008282546112879190618f3f565b9250508190555050505050505050565b606760149054906101000a900460ff16156112e7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112de90619012565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff166076600087815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160361138c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016113839061907e565b60405180910390fd5b60766000868152602001908152602001600020600001548410156113e5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016113dc906190ea565b60405180910390fd5b60006113f184846152ed565b905060005b8281101561141c5761140b8989898986615415565b8061141590619183565b90506113f6565b50606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd333085896114699190618f73565b6040518463ffffffff1660e01b81526004016114879392919061910a565b6020604051808303816000875af11580156114a6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114ca9190619156565b5081856114d79190618f73565b608660008282546114e89190618f3f565b925050819055505050505050505050565b6000606760149054906101000a900460ff161561154b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161154290619012565b60405180910390fd5b81607a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054607760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001546115d991906191cb565b101561161a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161161190619225565b60405180910390fd5b60006075544261162a9190618f3f565b90506001607860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461167c9190618f3f565b925050819055506000607860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050604051806040016040528083815260200185815250607960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000838152602001908152602001600020600082015181600001556020820151816001015590505083607a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546117929190618f3f565b92505081905550803373ffffffffffffffffffffffffffffffffffffffff167fcc7e0e76b20394ef965d71c4111d21e1b322bb8775dc2c7acd29eb0d3c3dd96d84876040516117e2929190618d99565b60405180910390a38092505050919050565b606a60009054906101000a900467ffffffffffffffff1681565b608c5481565b60008282856000015186602001518760a0015160405160200161183b959493929190619245565b6040516020818303038152906040528051906020012090509392505050565b606760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146118ea576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016118e1906192eb565b60405180910390fd5b80606760146101000a81548160ff0219169083151502179055508015157fd83d5281277e107f080e362699d46082adb74e7dc6a9bccbc87d8ae9533add4460405160405180910390a250565b60806020528060005260406000206000915090505481565b606760149054906101000a900460ff161561199e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161199590619012565b60405180910390fd5b6119a7336155b2565b6119e6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016119dd90619357565b60405180910390fd5b60006119f2338461365e565b14611a32576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611a29906193c3565b60405180910390fd5b611a3d828233615652565b5050565b607a6020528060005260406000206000915090505481565b60008083118015611a6a5750600082115b611aa9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611aa09061942f565b60405180910390fd5b6000611ab484614e17565b90506000611adb611ac483615872565b611acd86615872565b61587c90919063ffffffff16565b9050670cf4ad2e86166c9c611aef82615a15565b1215611b085768056bc75e2d6310000092505050611c17565b6000611b1b670de0b6b3a7640000615872565b90506000611b3168056bc75e2d63100000615872565b90506000611b7e611b6f84611b6185611b53888a615a1f90919063ffffffff16565b615a4d90919063ffffffff16565b615a1f90919063ffffffff16565b84615bdd90919063ffffffff16565b90506801158e460913d00000611b9382615a15565b12611ba657600095505050505050611c17565b611bc2611bb36000615872565b82615c0b90919063ffffffff16565b15611bec57611be0611bdb611bd683615c28565b615cb8565b615a15565b95505050505050611c17565b611c0f611c0a611bfb83615cb8565b8561587c90919063ffffffff16565b615a15565b955050505050505b92915050565b600081836020015184600001518560600151604051602001611c42949392919061944f565b60405160208183030381529060405280519060200120905092915050565b60856020528060005260406000206000915090505481565b600080611c83613440565b9050670de0b6b3a7640000606c54670de0b6b3a7640000611ca491906191cb565b82611caf9190618f73565b611cb991906194ca565b81611cc491906191cb565b91505090565b611cd2615dd6565b80608b60006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60826020528160005260406000208181548110611d3257600080fd5b906000526020600020016000915091509054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600084848484604051602001611d7d9493929190619528565b604051602081830303815290604052805190602001209050949350505050565b6000611ddb606a60009054906101000a900467ffffffffffffffff1667ffffffffffffffff1642611dce91906191cb565b611dd6613440565b61496e565b905090565b600080611deb613440565b9050670de0b6b3a7640000606d54670de0b6b3a7640000611e0c91906191cb565b82611e179190618f73565b611e2191906194ca565b81611e2c91906191cb565b91505090565b611e3a615dd6565b80606b819055507f8c854a81cb5c93e7e482d30fb9c6f88fdbdb320f10f7a853c2263659b54e563f81604051611e709190617ef7565b60405180910390a150565b6060611e8783836152ed565b905092915050565b611e97615dd6565b80606760006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff167f5a85b4270fc1e75035e6cd505418ce65e0bcc36cc7eb9ce9e6f8c6181d4cb57760405160405180910390a250565b608d60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b611f4c615dd6565b80608d60006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6000606760149054906101000a900460ff1615611fe2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611fd990619012565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff1603612051576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612048906195b4565b60405180910390fd5b600061205d84846152ed565b9050600060405180608001604052808781526020018873ffffffffffffffffffffffffffffffffffffffff1681526020016000815260200183815250905060006120a78233611c1d565b9050600073ffffffffffffffffffffffffffffffffffffffff166076600083815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461214e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161214590619620565b60405180910390fd5b81607660008381526020019081526020016000206000820151816000015560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040820151816002015560608201518160030190816121d29190619842565b50905050807fa4b0af38d049ba81703a0d0e46cc2ff39681210302134046237111a8fb7dee7260405160405180910390a2809350505050949350505050565b606760149054906101000a900460ff1615612261576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161225890619012565b60405180910390fd5b6000607c600083815260200190815260200160002054146122b7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016122ae90619960565b60405180910390fd5b6122bf615e54565b607c600083815260200190815260200160002081905550803373ffffffffffffffffffffffffffffffffffffffff167f09b4c028a2e50fec6f1c6a0163c59e8fbe92b231e5c03ef3adec585e63a14b9260405160405180910390a350565b606b5481565b606760149054906101000a900460ff1615612373576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161236a90619012565b60405180910390fd5b61237d6001615efe565b61238883838361613a565b505050565b606760149054906101000a900460ff1681565b606660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b606760149054906101000a900460ff1615612416576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161240d90619012565b60405180910390fd5b61242284849050615efe565b60005b848490508110156124885761247785858381811061244657612445619980565b5b905060200201358484848181106124605761245f619980565b5b905060200281019061247291906199be565b61613a565b8061248190619183565b9050612425565b5050505050565b606760149054906101000a900460ff16156124df576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016124d690619012565b60405180910390fd5b6124e8336155b2565b612527576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161251e90619357565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16607d600083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16036125cc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016125c390619a6d565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16607e600083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614612671576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161266890619ad9565b60405180910390fd5b607154607d600083815260200190815260200160002060000160149054906101000a900467ffffffffffffffff1667ffffffffffffffff166126b39190618f3f565b42106126f4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016126eb90619b45565b60405180910390fd5b607d6000828152602001908152602001600020600001601c9054906101000a900460ff1615612758576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161274f90619bb1565b60405180910390fd5b6000612762611de0565b905060405180608001604052803373ffffffffffffffffffffffffffffffffffffffff1681526020014267ffffffffffffffff168152602001600063ffffffff16815260200182815250607e600084815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160000160146101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550604082015181600001601c6101000a81548163ffffffff021916908363ffffffff16021790555060608201518160010155905050813373ffffffffffffffffffffffffffffffffffffffff167f6958c989e915d3e41a35076e3c480363910055408055ad86ae1ee13d41c4064060405160405180910390a36128b982600133615652565b8060776000607d600086815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001541061297b5761297a826000607d600086815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16615652565b5b5050565b60685481565b61298d615dd6565b61299760006164aa565b565b606f5481565b607d6020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060000160149054906101000a900467ffffffffffffffff169080600001601c9054906101000a900460ff1690806001018054612a159061966f565b80601f0160208091040260200160405190810160405280929190818152602001828054612a419061966f565b8015612a8e5780601f10612a6357610100808354040283529160200191612a8e565b820191906000526020600020905b815481529060010190602001808311612a7157829003601f168201915b5050505050905084565b606760149054906101000a900460ff1615612ae8576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612adf90619012565b60405180910390fd5b6000607960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008481526020019081526020016000206040518060400160405290816000820154815260200160018201548152505090506000816000015111612ba2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612b9990619c1d565b60405180910390fd5b8060000151421015612be9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612be090619c89565b60405180910390fd5b8060200151607760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001541015612c72576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612c6990619cf5565b60405180910390fd5b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb8383602001516040518363ffffffff1660e01b8152600401612cd3929190619d15565b6020604051808303816000875af1158015612cf2573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d169190619156565b508060200151607760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000016000828254612d6d91906191cb565b92505081905550806020015160866000828254612d8a91906191cb565b925050819055508060200151607a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254612de491906191cb565b92505081905550607960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600084815260200190815260200160002060008082016000905560018201600090555050828273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f109aeff667601aad33abcb7c5df1617754eefd18253d586a95c198cb479b5bd48460200151604051612eb09190617ef7565b60405180910390a4505050565b606760149054906101000a900460ff1615612f0d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612f0490619012565b60405180910390fd5b612f15611c78565b607a6000607d600085815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205460776000607d600086815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000015461300e91906191cb565b101561304f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161304690619d8a565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16607d600083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16036130f4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016130eb90619df6565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16607e600083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614613199576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161319090619e62565b60405180910390fd5b607154426131a791906191cb565b607d600083815260200190815260200160002060000160149054906101000a900467ffffffffffffffff1667ffffffffffffffff161061321c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161321390619ece565b60405180910390fd5b60735460715460856000607d600086815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546132a29190618f3f565b6132ac9190618f3f565b607d600083815260200190815260200160002060000160149054906101000a900467ffffffffffffffff1667ffffffffffffffff1611613321576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161331890619f60565b60405180910390fd5b60001515607d6000838152602001908152602001600020600001601c9054906101000a900460ff1615151461338b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161338290619fcc565b60405180910390fd5b6001607d6000838152602001908152602001600020600001601c6101000a81548160ff02191690831515021790555080607d600083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f0b76b4ae356796814d36b46f7c500bbd27b2cce1e6059a6fa2bebfd5a389b19060405160405180910390a361343d81616570565b50565b600080606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b815260040161349e91906187ac565b602060405180830381865afa1580156134bb573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134df919061a001565b9050697f0e10af47c1c700000081106134fc576000915050613536565b6000608c546086548361350f91906191cb565b61351991906191cb565b905080697f0e10af47c1c700000061353191906191cb565b925050505b90565b60735481565b6005600060019054906101000a900460ff1615801561357057508060ff1660008054906101000a900460ff1660ff16105b6135af576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016135a69061a0a0565b60405180910390fd5b806000806101000a81548160ff021916908360ff1602179055506001600060016101000a81548160ff02191690831515021790555060008060016101000a81548160ff0219169083151502179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024988160405161362d919061a0c0565b60405180910390a150565b608b60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008073ffffffffffffffffffffffffffffffffffffffff16607e600084815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16036136d25760019050613896565b6136db82610fcd565b156136e95760029050613896565b607f600083815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16156137555760039050613896565b6000607760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010154036137a85760049050613896565b607454607760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206001015410156137fd5760059050613896565b607e600083815260200190815260200160002060000160149054906101000a900467ffffffffffffffff1667ffffffffffffffff16607454607760008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206001015461388291906191cb565b11156138915760069050613896565b600090505b92915050565b607c6020528060005260406000206000915090505481565b606760149054906101000a900460ff1615613904576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016138fb90619012565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16607e600084815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16036139a9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016139a09061a127565b60405180910390fd5b6139b282610fcd565b6139f1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016139e89061a193565b60405180910390fd5b60008163ffffffff1611613a3a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613a319061a1ff565b60405180910390fd5b6000608160008481526020019081526020016000208054905090506000608260008581526020019081526020016000208054905090506000607e6000868152602001908152602001600020600001601c9054906101000a900463ffffffff169050600084607e6000888152602001908152602001600020600001601c9054906101000a900463ffffffff16613acf919061a21f565b90506000607e60008881526020019081526020016000206001015490508363ffffffff168563ffffffff161115614086576000818563ffffffff16613b149190618f73565b9050600060018763ffffffff1614613b4357600282613b3391906194ca565b82613b3e91906191cb565b613b45565b815b9050600060018863ffffffff1614613b8557600188613b64919061a257565b63ffffffff168284613b7691906191cb565b613b8091906194ca565b613b88565b60005b905060008663ffffffff1690505b8563ffffffff16811015613df6578863ffffffff16811015613de3576000608160008d81526020019081526020016000208281548110613bd957613bd8619980565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905085607760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000016000828254613c589190618f3f565b9250508190555060008203613d2657606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb82866040518363ffffffff1660e01b8152600401613cc4929190619d15565b6020604051808303816000875af1158015613ce3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613d079190619156565b508360866000828254613d1a91906191cb565b92505081905550613de1565b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb82856040518363ffffffff1660e01b8152600401613d83929190619d15565b6020604051808303816000875af1158015613da2573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613dc69190619156565b508260866000828254613dd991906191cb565b925050819055505b505b8080613dee90619183565b915050613b96565b506000607e60008c8152602001908152602001600020600001601c9054906101000a900463ffffffff1663ffffffff160361407e57606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb607b60008d815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16607b60008e8152602001908152602001600020600101546040518363ffffffff1660e01b8152600401613ed4929190619d15565b6020604051808303816000875af1158015613ef3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613f179190619156565b50607b60008b81526020019081526020016000206001015460866000828254613f4091906191cb565b92505081905550608360008b81526020019081526020016000205460776000608160008e8152602001908152602001600020600081548110613f8557613f84619980565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000016000828254613ffc9190618f3f565b925050819055504260856000607d60008e815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055505b5050506143b4565b6000818663ffffffff1661409a9190618f73565b9050600060018663ffffffff16146140be576002826140b991906194ca565b6140c0565b815b9050600060018763ffffffff1614614100576001876140df919061a257565b63ffffffff1682846140f191906191cb565b6140fb91906194ca565b614103565b60005b905060008663ffffffff1690505b8563ffffffff16811015614371578763ffffffff1681101561435e576000608260008d8152602001908152602001600020828154811061415457614153619980565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905085607760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160008282546141d39190618f3f565b92505081905550600082036142a157606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb82866040518363ffffffff1660e01b815260040161423f929190619d15565b6020604051808303816000875af115801561425e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906142829190619156565b50836086600082825461429591906191cb565b9250508190555061435c565b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb82856040518363ffffffff1660e01b81526004016142fe929190619d15565b6020604051808303816000875af115801561431d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906143419190619156565b50826086600082825461435491906191cb565b925050819055505b505b808061436990619183565b915050614111565b506000607e60008c8152602001908152602001600020600001601c9054906101000a900463ffffffff1663ffffffff16036143b0576143af8a616570565b5b5050505b81607e6000898152602001908152602001600020600001601c6101000a81548163ffffffff021916908363ffffffff1602179055508263ffffffff16877f71d8c71303e35a39162e33a402c9897bf9848388537bac7d5e1b0d202eca4e6684604051614420919061a28f565b60405180910390a350505050505050565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60745481565b60715481565b606760149054906101000a900460ff16156144b7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016144ae90619012565b60405180910390fd5b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd3330846040518463ffffffff1660e01b81526004016145169392919061910a565b6020604051808303816000875af1158015614535573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906145599190619156565b50806086600082825461456c9190618f3f565b92505081905550600061457d611c78565b905080607760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000015411614665578082607760008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001546146179190618f3f565b106146645742607760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600101819055505b5b604051806060016040528083607760008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001546146be9190618f3f565b8152602001607760008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206001015481526020018473ffffffffffffffffffffffffffffffffffffffff16815250607760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082015181600001556020820151816001015560408201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055509050508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8d4844488c19a90828439e71d14ebad860806d04f8ef8b25a82179fab2699b898460405161481f9190617ef7565b60405180910390a3505050565b614834615dd6565b600073ffffffffffffffffffffffffffffffffffffffff166076600084815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16036148d9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016148d09061907e565b60405180910390fd5b806076600084815260200190815260200160002060020181905550817f0321e8a918b4c47bf3677852c070983825f30a47bc8d9416691454fa6a727d63826040516149249190617ef7565b60405180910390a25050565b606c5481565b60845481565b606760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60875481565b608a5481565b600080820361498757670de0b6b3a764000090506149ef565b670de0b6b3a7640000697f0e10af47c1c70000006149a58585611a59565b670de0b6b3a764000085697f0e10af47c1c70000006149c491906191cb565b6149ce9190618f73565b6149d89190618f73565b6149e291906194ca565b6149ec91906194ca565b90505b92915050565b60755481565b614a03615dd6565b80606a60006101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055508067ffffffffffffffff167fa15d6b0930a82638ac4775bfd1b2e9f1e86be67e1bd3a09fa8a77a8f079769d760405160405180910390a250565b606760149054906101000a900460ff1615614ab6576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614aad90619012565b60405180910390fd5b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb606660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff166068546040518363ffffffff1660e01b8152600401614b37929190619d15565b6020604051808303816000875af1158015614b56573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614b7a9190619156565b5060685460866000828254614b8f91906191cb565b925050819055506000606881905550565b60836020528060005260406000206000915090505481565b60695481565b60895481565b60886020528060005260406000206000915090505481565b60705481565b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b606760149054906101000a900460ff1615614c58576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614c4f90619012565b60405180910390fd5b6000607960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008381526020019081526020016000206040518060400160405290816000820154815260200160018201548152505090506000816000015111614d12576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614d0990619c1d565b60405180910390fd5b8060200151607a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254614d6591906191cb565b92505081905550607960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600083815260200190815260200160002060008082016000905560018201600090555050813373ffffffffffffffffffffffffffffffffffffffff167ff9d5cfcdfe9803069225971ba315f4302add9b477c3441cc363bc38fdb065b9060405160405180910390a35050565b600063bbf81e00821115614e3757697f0e10af47c1c70000009050614ed2565b6000614e70614e6b614e66614e4f6301e13380617051565b614e5887617051565b61705b90919063ffffffff16565b617090565b617122565b9050670de0b6b3a764000081670de0b6b3a764000080697f0e10af47c1c7000000614e9b9190618f73565b614ea59190618f73565b614eaf91906194ca565b614eb991906194ca565b697f0e10af47c1c7000000614ece91906191cb565b9150505b919050565b60816020528160005260406000208181548110614ef357600080fd5b906000526020600020016000915091509054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60786020528060005260406000206000915090505481565b607f6020528160005260406000206020528060005260406000206000915091509054906101000a900460ff1681565b6079602052816000526040600020602052806000526040600020600091509150508060000154908060010154905082565b607e6020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060000160149054906101000a900467ffffffffffffffff169080600001601c9054906101000a900463ffffffff16908060010154905084565b615019615dd6565b80606660006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff167f6bdb9ceff405d990c9b60be9f719fbb80889d5f064e8fd76efd5bea353e9071260405160405180910390a250565b606d5481565b60766020528060005260406000206000915090508060000154908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154908060030180546150fb9061966f565b80601f01602080910402602001604051908101604052809291908181526020018280546151279061966f565b80156151745780601f1061514957610100808354040283529160200191615174565b820191906000526020600020905b81548152906001019060200180831161515757829003601f168201915b5050505050905084565b607b6020528060005260406000206000915090508060000154908060010154908060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020160149054906101000a900467ffffffffffffffff169080600201601c9054906101000a900460ff16908060030180546152009061966f565b80601f016020809104026020016040519081016040528092919081815260200182805461522c9061966f565b80156152795780601f1061524e57610100808354040283529160200191615279565b820191906000526020600020905b81548152906001019060200180831161525c57829003601f168201915b5050505050905086565b606e5481565b615291615dd6565b61529a8161712c565b50565b60865481565b60776020528060005260406000206000915090508060000154908060010154908060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905083565b606062010000838390501115615338576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161532f9061a2f6565b60405180910390fd5b6000615346848490506171af565b9050600081858584604051602001615361949392919061a3c3565b6040516020818303038152906040529050600261537e82516171af565b8260405160200161539092919061a43a565b6040516020818303038152906040526040516153ac919061a46d565b602060405180830381855afa1580156153c9573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906153ec919061a499565b6040516020016153fc919061a50d565b6040516020818303038152906040529250505092915050565b60006040518060c001604052808581526020018481526020018673ffffffffffffffffffffffffffffffffffffffff1681526020014267ffffffffffffffff1681526020018760ff16815260200183815250905060006154788233606954611814565b905081607b6000838152602001908152602001600020600082015181600001556020820151816001015560408201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060608201518160020160146101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550608082015181600201601c6101000a81548160ff021916908360ff16021790555060a082015181600301908161554e9190619842565b509050503373ffffffffffffffffffffffffffffffffffffffff1685827fc3d3e0544c80e3bb83f62659259ae1574f72a91515ab3cae3dd75cf77e1b0aea8760405161559a9190617ef7565b60405180910390a48060698190555050505050505050565b60006155bc611c78565b607a60008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054607760008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000015461564991906191cb565b10159050919050565b6001607f600085815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550811561573a5760816000848152602001908152602001600020819080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506157af565b60826000848152602001908152602001600020819080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505b607e600084815260200190815260200160002060010154607760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001600082825461581791906191cb565b92505081905550828173ffffffffffffffffffffffffffffffffffffffff167f1aa9e4be46e24e1f2e7eeb1613c01629213cd42965d2716e18531b63e552e411846040516158659190617fa4565b60405180910390a3505050565b6000819050919050565b60008061588884615a15565b9050600061589584615a15565b90507f80000000000000000000000000000000000000000000000000000000000000008214806158e457507f800000000000000000000000000000000000000000000000000000000000000081145b1561591b576040517f9fe2b45000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000841261592c5783615931565b836000035b9150600083126159415782615946565b826000035b9050600061595d83670de0b6b3a76400008461735e565b90507f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8111156159c65787876040517fd49c26b30000000000000000000000000000000000000000000000000000000081526004016159bd92919061a572565b60405180910390fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff858718139050615a0781615a005782600003615a02565b825b617478565b965050505050505092915050565b6000819050919050565b6000615a45615a2d83615a15565b615a3685615a15565b615a40919061a59b565b617478565b905092915050565b600080615a5984615a15565b90506000615a6684615a15565b90507f8000000000000000000000000000000000000000000000000000000000000000821480615ab557507f800000000000000000000000000000000000000000000000000000000000000081145b15615aec576040517fa6070c2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008060008412615afd5783615b02565b836000035b915060008312615b125782615b17565b826000035b90506000615b258383617482565b90507f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff811115615b8e5787876040517f120b5b43000000000000000000000000000000000000000000000000000000008152600401615b8592919061a572565b60405180910390fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff858718139050615bcf81615bc85782600003615bca565b825b617478565b965050505050505092915050565b6000615c03615beb83615a15565b615bf485615a15565b615bfe919061a5de565b617478565b905092915050565b6000615c1682615a15565b615c1f84615a15565b12905092915050565b600080615c3483615a15565b90507f80000000000000000000000000000000000000000000000000000000000000008103615c8f576040517fec2b9e6700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008112615c9d5782615cb0565b615caf81615caa9061a622565b617478565b5b915050919050565b600080615cc483615a15565b90506000811215615d4b577ffffffffffffffffffffffffffffffffffffffffffffffffcc22e87f6eb468eeb811215615d01576000915050615dd1565b615d44615d20615d1b615d1684600003617478565b615cb8565b615a15565b6ec097ce7bc90715b34b9f100000000081615d3e57615d3d61949b565b5b05617478565b9150615dcf565b680a688906bd8affffff811315615d9957826040517f0360d028000000000000000000000000000000000000000000000000000000008152600401615d90919061a66a565b60405180910390fd5b6000670de0b6b3a7640000604083901b81615db757615db661949b565b5b059050615dcb615dc68261756c565b617478565b9250505b505b919050565b615dde617ecc565b73ffffffffffffffffffffffffffffffffffffffff16615dfc614431565b73ffffffffffffffffffffffffffffffffffffffff1614615e52576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401615e499061a6d1565b60405180910390fd5b565b60008046905061a4ba811480615e6c575062066eed81145b80615e79575062066eee81145b15615ef657606473ffffffffffffffffffffffffffffffffffffffff1663a3b1b31d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015615eca573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190615eee919061a001565b915050615efb565b439150505b90565b607354607154608560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054615f4e9190618f3f565b615f589190618f3f565b4211615f99576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401615f909061a763565b60405180910390fd5b80608454615fa79190618f73565b607760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000016000828254615ff891906191cb565b92505081905550670de0b6b3a7640000816087546160169190618f73565b61602091906194ca565b608860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020544261606b91906191cb565b116160ab576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016160a29061a7cf565b60405180910390fd5b42608860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506160f8336155b2565b616137576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161612e90619357565b60405180910390fd5b50565b6000801b607b60008581526020019081526020016000206000015403616195576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161618c9061a83b565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16607d600085815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461623a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016162319061a8a7565b60405180910390fd5b600061624833858585611d64565b90506000607c600083815260200190815260200160002054116162a0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016162979061a913565b60405180910390fd5b6162a8615e54565b607c600083815260200190815260200160002054106162fc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016162f39061a97f565b60405180910390fd5b60405180608001604052803373ffffffffffffffffffffffffffffffffffffffff1681526020014267ffffffffffffffff16815260200160001515815260200184848080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050815250607d600086815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160000160146101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550604082015181600001601c6101000a81548160ff02191690831515021790555060608201518160010190816164429190619842565b509050506084546083600086815260200190815260200160002081905550833373ffffffffffffffffffffffffffffffffffffffff167f957c18b5af8413899ea8a576a4d3fb16839a02c9fccfdce098b6d59ef248525b60405160405180910390a350505050565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081603360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000607b6000838152602001908152602001600020600001549050600060766000838152602001908152602001600020600001549050607b6000848152602001908152602001600020600101548111156165c957600090505b60008111156166aa57606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb6076600085815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16836040518363ffffffff1660e01b8152600401616665929190619d15565b6020604051808303816000875af1158015616684573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906166a89190619156565b505b600081607b6000868152602001908152602001600020600101546166ce91906191cb565b90506000670de0b6b3a7640000606e54670de0b6b3a76400006166f191906191cb565b836166fc9190618f73565b61670691906194ca565b8261671191906191cb565b905080606860008282546167259190618f3f565b925050819055506000818361673a91906191cb565b9050600081111561682857606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb607d600089815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1684866167c691906191cb565b6040518363ffffffff1660e01b81526004016167e3929190619d15565b6020604051808303816000875af1158015616802573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906168269190619156565b505b608b60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ebe2b12b6040518163ffffffff1660e01b8152600401602060405180830381865afa158015616895573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906168b9919061a001565b421115616a1d57606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb608b60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16608c546040518363ffffffff1660e01b8152600401616941929190619d15565b6020604051808303816000875af1158015616960573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906169849190619156565b50608b60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633c6b16ab608c546040518263ffffffff1660e01b81526004016169e29190617ef7565b600060405180830381600087803b1580156169fc57600080fd5b505af1158015616a10573d6000803e3d6000fd5b505050506000608c819055505b6000607660008781526020019081526020016000206002015490506000670de0b6b3a76400009050600073ffffffffffffffffffffffffffffffffffffffff16608d60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614616be557608d60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b65a78b5886040518263ffffffff1660e01b8152600401616af6919061854c565b602060405180830381865afa158015616b13573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190616b379190619156565b15616bdf57608d60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16637eca406f886040518263ffffffff1660e01b8152600401616b97919061854c565b602060405180830381865afa158015616bb4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190616bd8919061a001565b9050616be4565b600090505b5b600082118015616bf55750600081115b15616f6b5760006f01812f9cf7920e2b66973e20000000008284616c17611d9d565b616c219190618f73565b616c2b9190618f73565b616c3591906194ca565b905080608c6000828254616c499190618f3f565b925050819055506000811115616f69576000670de0b6b3a7640000607054670de0b6b3a7640000616c7a91906191cb565b83616c859190618f73565b616c8f91906194ca565b82616c9a91906191cb565b90506000670de0b6b3a7640000608954670de0b6b3a7640000616cbd91906191cb565b84616cc89190618f73565b616cd291906194ca565b83616cdd91906191cb565b9050606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb607d60008e815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16838587616d6191906191cb565b616d6b91906191cb565b6040518363ffffffff1660e01b8152600401616d88929190619d15565b6020604051808303816000875af1158015616da7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190616dcb9190619156565b50606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb606660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16846040518363ffffffff1660e01b8152600401616e4b929190619d15565b6020604051808303816000875af1158015616e6a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190616e8e9190619156565b50606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb607b60008e815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16836040518363ffffffff1660e01b8152600401616f22929190619d15565b6020604051808303816000875af1158015616f41573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190616f659190619156565b5050505b505b608360008981526020019081526020016000205460776000607d60008c815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160008282546170069190618f3f565b9250508190555083607b60008a81526020019081526020016000206001015461702f91906191cb565b6086600082825461704091906191cb565b925050819055505050505050505050565b6000819050919050565b600061708861708361706c85617122565b670de0b6b3a764000061707e86617122565b61735e565b617ed4565b905092915050565b60008061709c83617122565b9050680a688906bd8affffff8111156170ec57826040517fb3b6ba1f0000000000000000000000000000000000000000000000000000000081526004016170e3919061a9ae565b60405180910390fd5b6000670de0b6b3a7640000604083901b61710691906194ca565b90506171196171148261756c565b617ed4565b92505050919050565b6000819050919050565b617134615dd6565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036171a3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161719a9061aa3b565b60405180910390fd5b6171ac816164aa565b50565b606060008290506000600190505b607f8267ffffffffffffffff1611156171f25760078267ffffffffffffffff16901c9150806171eb9061aa5b565b90506171bd565b60008167ffffffffffffffff1667ffffffffffffffff811115617218576172176182b1565b5b6040519080825280601f01601f19166020018201604052801561724a5781602001600182028036833780820191505090505b50905084925060005b8267ffffffffffffffff168167ffffffffffffffff1610156172e957607f841660801760f81b828267ffffffffffffffff168151811061729657617295619980565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535060078467ffffffffffffffff16901c935080806172e19061aa5b565b915050617253565b50607f60f81b816001846172fd919061aa8b565b67ffffffffffffffff168151811061731857617317619980565b5b6020010181815160f81c60f81b169150907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350809350505050919050565b60008060008019858709858702925082811083820303915050600081036173995783828161738f5761738e61949b565b5b0492505050617471565b8381106173e1578585856040517f63a057780000000000000000000000000000000000000000000000000000000081526004016173d89392919061aac7565b60405180910390fd5b60008486880990508281118203915080830392506000600186190186169050600081870496508185049450600182836000030401905080840285179450600060028860030218905080880260020381029050808802600203810290508088026002038102905080880260020381029050808802600203810290508088026002038102905080860296505050505050505b9392505050565b6000819050919050565b60008060008019848609848602925082811083820303915050600081036174c557670de0b6b3a764000082816174bb576174ba61949b565b5b0492505050617566565b670de0b6b3a764000081106175135784846040517f5173648d00000000000000000000000000000000000000000000000000000000815260040161750a929190618d99565b60405180910390fd5b6000670de0b6b3a764000085870990507faccb18165bd6fe31ae1cf318dc5b51eee0e1ba569b88cd74c1773b91fac106696001620400008060000304018483118403026204000083860304170293505050505b92915050565b6000778000000000000000000000000000000000000000000000009050600067ff00000000000000831611156176cd576000678000000000000000831611156175c257604068016a09e667f3bcc9098202901c90505b6000674000000000000000831611156175e85760406801306fe0a31b7152df8202901c90505b60006720000000000000008316111561760e5760406801172b83c7d517adce8202901c90505b60006710000000000000008316111561763457604068010b5586cf9890f62a8202901c90505b60006708000000000000008316111561765a5760406801059b0d31585743ae8202901c90505b600067040000000000000083161115617680576040680102c9a3e778060ee78202901c90505b6000670200000000000000831611156176a657604068010163da9fb33356d88202901c90505b6000670100000000000000831611156176cc576040680100b1afa5abcbed618202901c90505b5b600066ff0000000000008316111561780857600066800000000000008316111561770457604068010058c86da1c09ea28202901c90505b60006640000000000000831611156177295760406801002c605e2e8cec508202901c90505b600066200000000000008316111561774e576040680100162f3904051fa18202901c90505b60006610000000000000831611156177735760406801000b175effdc76ba8202901c90505b6000660800000000000083161115617798576040680100058ba01fb9f96d8202901c90505b60006604000000000000831611156177bd57604068010002c5cc37da94928202901c90505b60006602000000000000831611156177e25760406801000162e525ee05478202901c90505b600066010000000000008316111561780757604068010000b17255775c048202901c90505b5b600065ff00000000008316111561793a576000658000000000008316111561783d5760406801000058b91b5bc9ae8202901c90505b60006540000000000083161115617861576040680100002c5c89d5ec6d8202901c90505b6000652000000000008316111561788557604068010000162e43f4f8318202901c90505b600065100000000000831611156178a9576040680100000b1721bcfc9a8202901c90505b600065080000000000831611156178cd57604068010000058b90cf1e6e8202901c90505b600065040000000000831611156178f15760406801000002c5c863b73f8202901c90505b60006502000000000083161115617915576040680100000162e430e5a28202901c90505b600065010000000000831611156179395760406801000000b1721835518202901c90505b5b600064ff0000000083161115617a635760006480000000008316111561796d576040680100000058b90c0b498202901c90505b60006440000000008316111561799057604068010000002c5c8601cc8202901c90505b6000642000000000831611156179b35760406801000000162e42fff08202901c90505b6000641000000000831611156179d657604068010000000b17217fbb8202901c90505b6000640800000000831611156179f95760406801000000058b90bfce8202901c90505b600064040000000083161115617a1c576040680100000002c5c85fe38202901c90505b600064020000000083161115617a3f57604068010000000162e42ff18202901c90505b600064010000000083161115617a62576040680100000000b17217f88202901c90505b5b600063ff00000083161115617b83576000638000000083161115617a9457604068010000000058b90bfc8202901c90505b6000634000000083161115617ab65760406801000000002c5c85fe8202901c90505b6000632000000083161115617ad8576040680100000000162e42ff8202901c90505b6000631000000083161115617afa5760406801000000000b17217f8202901c90505b6000630800000083161115617b1c576040680100000000058b90c08202901c90505b6000630400000083161115617b3e57604068010000000002c5c8608202901c90505b6000630200000083161115617b605760406801000000000162e4308202901c90505b6000630100000083161115617b8257604068010000000000b172188202901c90505b5b600062ff000083161115617c9a5760006280000083161115617bb25760406801000000000058b90c8202901c90505b60006240000083161115617bd3576040680100000000002c5c868202901c90505b60006220000083161115617bf457604068010000000000162e438202901c90505b60006210000083161115617c15576040680100000000000b17218202901c90505b60006208000083161115617c3657604068010000000000058b918202901c90505b60006204000083161115617c575760406801000000000002c5c88202901c90505b60006202000083161115617c78576040680100000000000162e48202901c90505b60006201000083161115617c995760406801000000000000b1728202901c90505b5b600061ff0083161115617da857600061800083161115617cc7576040680100000000000058b98202901c90505b600061400083161115617ce757604068010000000000002c5d8202901c90505b600061200083161115617d075760406801000000000000162e8202901c90505b600061100083161115617d2757604068010000000000000b178202901c90505b600061080083161115617d475760406801000000000000058c8202901c90505b600061040083161115617d67576040680100000000000002c68202901c90505b600061020083161115617d87576040680100000000000001638202901c90505b600061010083161115617da7576040680100000000000000b18202901c90505b5b600060ff83161115617ead576000608083161115617dd3576040680100000000000000598202901c90505b6000604083161115617df25760406801000000000000002c8202901c90505b6000602083161115617e11576040680100000000000000168202901c90505b6000601083161115617e305760406801000000000000000b8202901c90505b6000600883161115617e4f576040680100000000000000068202901c90505b6000600483161115617e6e576040680100000000000000038202901c90505b6000600283161115617e8d576040680100000000000000018202901c90505b6000600183161115617eac576040680100000000000000018202901c90505b5b670de0b6b3a764000081029050604082901c60bf0381901c9050919050565b600033905090565b6000819050919050565b6000819050919050565b617ef181617ede565b82525050565b6000602082019050617f0c6000830184617ee8565b92915050565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b617f3981617f26565b8114617f4457600080fd5b50565b600081359050617f5681617f30565b92915050565b600060208284031215617f7257617f71617f1c565b5b6000617f8084828501617f47565b91505092915050565b60008115159050919050565b617f9e81617f89565b82525050565b6000602082019050617fb96000830184617f95565b92915050565b600060ff82169050919050565b617fd581617fbf565b8114617fe057600080fd5b50565b600081359050617ff281617fcc565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061802382617ff8565b9050919050565b61803381618018565b811461803e57600080fd5b50565b6000813590506180508161802a565b92915050565b61805f81617ede565b811461806a57600080fd5b50565b60008135905061807c81618056565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f8401126180a7576180a6618082565b5b8235905067ffffffffffffffff8111156180c4576180c3618087565b5b6020830191508360018202830111156180e0576180df61808c565b5b9250929050565b60008060008060008060a0878903121561810457618103617f1c565b5b600061811289828a01617fe3565b965050602061812389828a01618041565b955050604061813489828a01617f47565b945050606061814589828a0161806d565b935050608087013567ffffffffffffffff81111561816657618165617f21565b5b61817289828a01618091565b92509250509295509295509295565b600080600080600080600060c0888a0312156181a05761819f617f1c565b5b60006181ae8a828b01617fe3565b97505060206181bf8a828b01618041565b96505060406181d08a828b01617f47565b95505060606181e18a828b0161806d565b945050608088013567ffffffffffffffff81111561820257618201617f21565b5b61820e8a828b01618091565b935093505060a06182218a828b0161806d565b91505092959891949750929550565b60006020828403121561824657618245617f1c565b5b60006182548482850161806d565b91505092915050565b600067ffffffffffffffff82169050919050565b61827a8161825d565b82525050565b60006020820190506182956000830184618271565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6182e9826182a0565b810181811067ffffffffffffffff82111715618308576183076182b1565b5b80604052505050565b600061831b617f12565b905061832782826182e0565b919050565b600080fd5b61833a8161825d565b811461834557600080fd5b50565b60008135905061835781618331565b92915050565b600080fd5b600067ffffffffffffffff82111561837d5761837c6182b1565b5b618386826182a0565b9050602081019050919050565b82818337600083830152505050565b60006183b56183b084618362565b618311565b9050828152602081018484840111156183d1576183d061835d565b5b6183dc848285618393565b509392505050565b600082601f8301126183f9576183f8618082565b5b81356184098482602086016183a2565b91505092915050565b600060c082840312156184285761842761829b565b5b61843260c0618311565b9050600061844284828501617f47565b60008301525060206184568482850161806d565b602083015250604061846a84828501618041565b604083015250606061847e84828501618348565b606083015250608061849284828501617fe3565b60808301525060a082013567ffffffffffffffff8111156184b6576184b561832c565b5b6184c2848285016183e4565b60a08301525092915050565b6000806000606084860312156184e7576184e6617f1c565b5b600084013567ffffffffffffffff81111561850557618504617f21565b5b61851186828701618412565b935050602061852286828701618041565b925050604061853386828701617f47565b9150509250925092565b61854681617f26565b82525050565b6000602082019050618561600083018461853d565b92915050565b61857081617f89565b811461857b57600080fd5b50565b60008135905061858d81618567565b92915050565b6000602082840312156185a9576185a8617f1c565b5b60006185b78482850161857e565b91505092915050565b600080604083850312156185d7576185d6617f1c565b5b60006185e585828601617f47565b92505060206185f68582860161857e565b9150509250929050565b60006020828403121561861657618615617f1c565b5b600061862484828501618041565b91505092915050565b6000806040838503121561864457618643617f1c565b5b60006186528582860161806d565b92505060206186638582860161806d565b9150509250929050565b6000608082840312156186835761868261829b565b5b61868d6080618311565b9050600061869d8482850161806d565b60008301525060206186b184828501618041565b60208301525060406186c58482850161806d565b604083015250606082013567ffffffffffffffff8111156186e9576186e861832c565b5b6186f5848285016183e4565b60608301525092915050565b6000806040838503121561871857618717617f1c565b5b600083013567ffffffffffffffff81111561873657618735617f21565b5b6187428582860161866d565b925050602061875385828601618041565b9150509250929050565b6000806040838503121561877457618773617f1c565b5b600061878285828601617f47565b92505060206187938582860161806d565b9150509250929050565b6187a681618018565b82525050565b60006020820190506187c1600083018461879d565b92915050565b600080600080606085870312156187e1576187e0617f1c565b5b60006187ef87828801618041565b945050602061880087828801617f47565b935050604085013567ffffffffffffffff81111561882157618820617f21565b5b61882d87828801618091565b925092505092959194509250565b6000806020838503121561885257618851617f1c565b5b600083013567ffffffffffffffff8111156188705761886f617f21565b5b61887c85828601618091565b92509250509250929050565b600081519050919050565b600082825260208201905092915050565b60005b838110156188c25780820151818401526020810190506188a7565b60008484015250505050565b60006188d982618888565b6188e38185618893565b93506188f38185602086016188a4565b6188fc816182a0565b840191505092915050565b6000602082019050818103600083015261892181846188ce565b905092915050565b6000806000806060858703121561894357618942617f1c565b5b600061895187828801618041565b94505060206189628782880161806d565b935050604085013567ffffffffffffffff81111561898357618982617f21565b5b61898f87828801618091565b925092505092959194509250565b6000806000604084860312156189b6576189b5617f1c565b5b60006189c486828701617f47565b935050602084013567ffffffffffffffff8111156189e5576189e4617f21565b5b6189f186828701618091565b92509250509250925092565b60008083601f840112618a1357618a12618082565b5b8235905067ffffffffffffffff811115618a3057618a2f618087565b5b602083019150836020820283011115618a4c57618a4b61808c565b5b9250929050565b60008083601f840112618a6957618a68618082565b5b8235905067ffffffffffffffff811115618a8657618a85618087565b5b602083019150836020820283011115618aa257618aa161808c565b5b9250929050565b60008060008060408587031215618ac357618ac2617f1c565b5b600085013567ffffffffffffffff811115618ae157618ae0617f21565b5b618aed878288016189fd565b9450945050602085013567ffffffffffffffff811115618b1057618b0f617f21565b5b618b1c87828801618a53565b925092505092959194509250565b6000608082019050618b3f600083018761879d565b618b4c6020830186618271565b618b596040830185617f95565b8181036060830152618b6b81846188ce565b905095945050505050565b60008060408385031215618b8d57618b8c617f1c565b5b6000618b9b8582860161806d565b9250506020618bac85828601618041565b9150509250929050565b60008060408385031215618bcd57618bcc617f1c565b5b6000618bdb85828601618041565b9250506020618bec85828601617f47565b9150509250929050565b600063ffffffff82169050919050565b618c0f81618bf6565b8114618c1a57600080fd5b50565b600081359050618c2c81618c06565b92915050565b60008060408385031215618c4957618c48617f1c565b5b6000618c5785828601617f47565b9250506020618c6885828601618c1d565b9150509250929050565b60008060408385031215618c8957618c88617f1c565b5b6000618c9785828601618041565b9250506020618ca88582860161806d565b9150509250929050565b600060208284031215618cc857618cc7617f1c565b5b6000618cd684828501618348565b91505092915050565b6000819050919050565b6000618d04618cff618cfa84617ff8565b618cdf565b617ff8565b9050919050565b6000618d1682618ce9565b9050919050565b6000618d2882618d0b565b9050919050565b618d3881618d1d565b82525050565b6000602082019050618d536000830184618d2f565b92915050565b60008060408385031215618d7057618d6f617f1c565b5b6000618d7e85828601617f47565b9250506020618d8f85828601618041565b9150509250929050565b6000604082019050618dae6000830185617ee8565b618dbb6020830184617ee8565b9392505050565b618dcb81618bf6565b82525050565b6000608082019050618de6600083018761879d565b618df36020830186618271565b618e006040830185618dc2565b618e0d6060830184617ee8565b95945050505050565b6000608082019050618e2b6000830187617ee8565b618e38602083018661879d565b618e456040830185617ee8565b8181036060830152618e5781846188ce565b905095945050505050565b618e6b81617fbf565b82525050565b600060c082019050618e86600083018961853d565b618e936020830188617ee8565b618ea0604083018761879d565b618ead6060830186618271565b618eba6080830185618e62565b81810360a0830152618ecc81846188ce565b9050979650505050505050565b6000606082019050618eee6000830186617ee8565b618efb6020830185617ee8565b618f08604083018461879d565b949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000618f4a82617ede565b9150618f5583617ede565b9250828201905080821115618f6d57618f6c618f10565b5b92915050565b6000618f7e82617ede565b9150618f8983617ede565b9250828202618f9781617ede565b91508282048414831517618fae57618fad618f10565b5b5092915050565b600082825260208201905092915050565b7f7061757365640000000000000000000000000000000000000000000000000000600082015250565b6000618ffc600683618fb5565b915061900782618fc6565b602082019050919050565b6000602082019050818103600083015261902b81618fef565b9050919050565b7f6d6f64656c20646f6573206e6f74206578697374000000000000000000000000600082015250565b6000619068601483618fb5565b915061907382619032565b602082019050919050565b600060208201905081810360008301526190978161905b565b9050919050565b7f6c6f77657220666565207468616e206d6f64656c206665650000000000000000600082015250565b60006190d4601883618fb5565b91506190df8261909e565b602082019050919050565b60006020820190508181036000830152619103816190c7565b9050919050565b600060608201905061911f600083018661879d565b61912c602083018561879d565b6191396040830184617ee8565b949350505050565b60008151905061915081618567565b92915050565b60006020828403121561916c5761916b617f1c565b5b600061917a84828501619141565b91505092915050565b600061918e82617ede565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036191c0576191bf618f10565b5b600182019050919050565b60006191d682617ede565b91506191e183617ede565b92508282039050818111156191f9576191f8618f10565b5b92915050565b50565b600061920f600083618fb5565b915061921a826191ff565b600082019050919050565b6000602082019050818103600083015261923e81619202565b9050919050565b600060a08201905061925a600083018861879d565b619267602083018761853d565b619274604083018661853d565b6192816060830185617ee8565b818103608083015261929381846188ce565b90509695505050505050565b7f6e6f742070617573657200000000000000000000000000000000000000000000600082015250565b60006192d5600a83618fb5565b91506192e08261929f565b602082019050919050565b60006020820190508181036000830152619304816192c8565b9050919050565b7f6d696e207374616b656420746f6f206c6f770000000000000000000000000000600082015250565b6000619341601283618fb5565b915061934c8261930b565b602082019050919050565b6000602082019050818103600083015261937081619334565b9050919050565b7f6e6f7420616c6c6f776564000000000000000000000000000000000000000000600082015250565b60006193ad600b83618fb5565b91506193b882619377565b602082019050919050565b600060208201905081810360008301526193dc816193a0565b9050919050565b7f6d696e2076616c73000000000000000000000000000000000000000000000000600082015250565b6000619419600883618fb5565b9150619424826193e3565b602082019050919050565b600060208201905081810360008301526194488161940c565b9050919050565b6000608082019050619464600083018761879d565b619471602083018661879d565b61947e6040830185617ee8565b818103606083015261949081846188ce565b905095945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006194d582617ede565b91506194e083617ede565b9250826194f0576194ef61949b565b5b828204905092915050565b60006195078385618893565b9350619514838584618393565b61951d836182a0565b840190509392505050565b600060608201905061953d600083018761879d565b61954a602083018661853d565b818103604083015261955d8184866194fb565b905095945050505050565b7f61646472657373206d757374206265206e6f6e2d7a65726f0000000000000000600082015250565b600061959e601883618fb5565b91506195a982619568565b602082019050919050565b600060208201905081810360008301526195cd81619591565b9050919050565b7f6d6f64656c20616c726561647920726567697374657265640000000000000000600082015250565b600061960a601883618fb5565b9150619615826195d4565b602082019050919050565b60006020820190508181036000830152619639816195fd565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061968757607f821691505b60208210810361969a57619699619640565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026197027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826196c5565b61970c86836196c5565b95508019841693508086168417925050509392505050565b600061973f61973a61973584617ede565b618cdf565b617ede565b9050919050565b6000819050919050565b61975983619724565b61976d61976582619746565b8484546196d2565b825550505050565b600090565b619782619775565b61978d818484619750565b505050565b5b818110156197b1576197a660008261977a565b600181019050619793565b5050565b601f8211156197f6576197c7816196a0565b6197d0846196b5565b810160208510156197df578190505b6197f36197eb856196b5565b830182619792565b50505b505050565b600082821c905092915050565b6000619819600019846008026197fb565b1980831691505092915050565b60006198328383619808565b9150826002028217905092915050565b61984b82618888565b67ffffffffffffffff811115619864576198636182b1565b5b61986e825461966f565b6198798282856197b5565b600060209050601f8311600181146198ac576000841561989a578287015190505b6198a48582619826565b86555061990c565b601f1984166198ba866196a0565b60005b828110156198e2578489015182556001820191506020850194506020810190506198bd565b868310156198ff57848901516198fb601f891682619808565b8355505b6001600288020188555050505b505050505050565b7f636f6d6d69746d656e7420657869737473000000000000000000000000000000600082015250565b600061994a601183618fb5565b915061995582619914565b602082019050919050565b600060208201905081810360008301526199798161993d565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600080fd5b600080fd5b600080fd5b600080833560016020038436030381126199db576199da6199af565b5b80840192508235915067ffffffffffffffff8211156199fd576199fc6199b4565b5b602083019250600182023603831315619a1957619a186199b9565b5b509250929050565b7f736f6c7574696f6e20646f6573206e6f74206578697374000000000000000000600082015250565b6000619a57601783618fb5565b9150619a6282619a21565b602082019050919050565b60006020820190508181036000830152619a8681619a4a565b9050919050565b7f636f6e746573746174696f6e20616c7265616479206578697374730000000000600082015250565b6000619ac3601b83618fb5565b9150619ace82619a8d565b602082019050919050565b60006020820190508181036000830152619af281619ab6565b9050919050565b7f746f6f206c617465000000000000000000000000000000000000000000000000600082015250565b6000619b2f600883618fb5565b9150619b3a82619af9565b602082019050919050565b60006020820190508181036000830152619b5e81619b22565b9050919050565b7f7774660000000000000000000000000000000000000000000000000000000000600082015250565b6000619b9b600383618fb5565b9150619ba682619b65565b602082019050919050565b60006020820190508181036000830152619bca81619b8e565b9050919050565b7f72657175657374206e6f74206578697374000000000000000000000000000000600082015250565b6000619c07601183618fb5565b9150619c1282619bd1565b602082019050919050565b60006020820190508181036000830152619c3681619bfa565b9050919050565b7f77616974206c6f6e676572000000000000000000000000000000000000000000600082015250565b6000619c73600b83618fb5565b9150619c7e82619c3d565b602082019050919050565b60006020820190508181036000830152619ca281619c66565b9050919050565b7f7374616b6520696e73756666696369656e740000000000000000000000000000600082015250565b6000619cdf601283618fb5565b9150619cea82619ca9565b602082019050919050565b60006020820190508181036000830152619d0e81619cd2565b9050919050565b6000604082019050619d2a600083018561879d565b619d376020830184617ee8565b9392505050565b7f76616c696461746f72206d696e207374616b656420746f6f206c6f7700000000600082015250565b6000619d74601c83618fb5565b9150619d7f82619d3e565b602082019050919050565b60006020820190508181036000830152619da381619d67565b9050919050565b7f736f6c7574696f6e206e6f7420666f756e640000000000000000000000000000600082015250565b6000619de0601283618fb5565b9150619deb82619daa565b602082019050919050565b60006020820190508181036000830152619e0f81619dd3565b9050919050565b7f68617320636f6e746573746174696f6e00000000000000000000000000000000600082015250565b6000619e4c601083618fb5565b9150619e5782619e16565b602082019050919050565b60006020820190508181036000830152619e7b81619e3f565b9050919050565b7f6e6f7420656e6f7567682064656c617900000000000000000000000000000000600082015250565b6000619eb8601083618fb5565b9150619ec382619e82565b602082019050919050565b60006020820190508181036000830152619ee781619eab565b9050919050565b7f636c61696d536f6c7574696f6e20636f6f6c646f776e206166746572206c6f7360008201527f7420636f6e746573746174696f6e000000000000000000000000000000000000602082015250565b6000619f4a602e83618fb5565b9150619f5582619eee565b604082019050919050565b60006020820190508181036000830152619f7981619f3d565b9050919050565b7f616c726561647920636c61696d65640000000000000000000000000000000000600082015250565b6000619fb6600f83618fb5565b9150619fc182619f80565b602082019050919050565b60006020820190508181036000830152619fe581619fa9565b9050919050565b600081519050619ffb81618056565b92915050565b60006020828403121561a0175761a016617f1c565b5b600061a02584828501619fec565b91505092915050565b7f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160008201527f647920696e697469616c697a6564000000000000000000000000000000000000602082015250565b600061a08a602e83618fb5565b915061a0958261a02e565b604082019050919050565b6000602082019050818103600083015261a0b98161a07d565b9050919050565b600060208201905061a0d56000830184618e62565b92915050565b7f636f6e746573746174696f6e20646f65736e2774206578697374000000000000600082015250565b600061a111601a83618fb5565b915061a11c8261a0db565b602082019050919050565b6000602082019050818103600083015261a1408161a104565b9050919050565b7f766f74696e6720706572696f64206e6f7420656e646564000000000000000000600082015250565b600061a17d601783618fb5565b915061a1888261a147565b602082019050919050565b6000602082019050818103600083015261a1ac8161a170565b9050919050565b7f616d6e7420746f6f20736d616c6c000000000000000000000000000000000000600082015250565b600061a1e9600e83618fb5565b915061a1f48261a1b3565b602082019050919050565b6000602082019050818103600083015261a2188161a1dc565b9050919050565b600061a22a82618bf6565b915061a23583618bf6565b9250828201905063ffffffff81111561a2515761a250618f10565b5b92915050565b600061a26282618bf6565b915061a26d83618bf6565b9250828203905063ffffffff81111561a2895761a288618f10565b5b92915050565b600060208201905061a2a46000830184618dc2565b92915050565b7f4d617820636f6e74656e742073697a6520697320363535333620627974657300600082015250565b600061a2e0601f83618fb5565b915061a2eb8261a2aa565b602082019050919050565b6000602082019050818103600083015261a30f8161a2d3565b9050919050565b7f0802120000000000000000000000000000000000000000000000000000000000815250565b600081905092915050565b600061a35282618888565b61a35c818561a33c565b935061a36c8185602086016188a4565b80840191505092915050565b600061a384838561a33c565b935061a391838584618393565b82840190509392505050565b7f1800000000000000000000000000000000000000000000000000000000000000815250565b600061a3ce8261a316565b60038201915061a3de828761a347565b915061a3eb82858761a378565b915061a3f68261a39d565b60018201915061a406828461a347565b915081905095945050505050565b7f0a00000000000000000000000000000000000000000000000000000000000000815250565b600061a4458261a414565b60018201915061a455828561a347565b915061a461828461a347565b91508190509392505050565b600061a479828461a347565b915081905092915050565b60008151905061a49381617f30565b92915050565b60006020828403121561a4af5761a4ae617f1c565b5b600061a4bd8482850161a484565b91505092915050565b7f1220000000000000000000000000000000000000000000000000000000000000815250565b6000819050919050565b61a50761a50282617f26565b61a4ec565b82525050565b600061a5188261a4c6565b60028201915061a528828461a4f6565b60208201915081905092915050565b6000819050919050565b600061a55c61a55761a5528461a537565b618cdf565b61a537565b9050919050565b61a56c8161a541565b82525050565b600060408201905061a587600083018561a563565b61a594602083018461a563565b9392505050565b600061a5a68261a537565b915061a5b18361a537565b925082820390508181126000841216828213600085121516171561a5d85761a5d7618f10565b5b92915050565b600061a5e98261a537565b915061a5f48361a537565b92508282019050828112156000831216838212600084121516171561a61c5761a61b618f10565b5b92915050565b600061a62d8261a537565b91507f8000000000000000000000000000000000000000000000000000000000000000820361a65f5761a65e618f10565b5b816000039050919050565b600060208201905061a67f600083018461a563565b92915050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b600061a6bb602083618fb5565b915061a6c68261a685565b602082019050919050565b6000602082019050818103600083015261a6ea8161a6ae565b9050919050565b7f7375626d6974536f6c7574696f6e20636f6f6c646f776e206166746572206c6f60008201527f737420636f6e746573746174696f6e0000000000000000000000000000000000602082015250565b600061a74d602f83618fb5565b915061a7588261a6f1565b604082019050919050565b6000602082019050818103600083015261a77c8161a740565b9050919050565b7f736f6c7574696f6e2072617465206c696d697400000000000000000000000000600082015250565b600061a7b9601383618fb5565b915061a7c48261a783565b602082019050919050565b6000602082019050818103600083015261a7e88161a7ac565b9050919050565b7f7461736b20646f6573206e6f7420657869737400000000000000000000000000600082015250565b600061a825601383618fb5565b915061a8308261a7ef565b602082019050919050565b6000602082019050818103600083015261a8548161a818565b9050919050565b7f736f6c7574696f6e20616c7265616479207375626d6974746564000000000000600082015250565b600061a891601a83618fb5565b915061a89c8261a85b565b602082019050919050565b6000602082019050818103600083015261a8c08161a884565b9050919050565b7f6e6f6e206578697374656e7420636f6d6d69746d656e74000000000000000000600082015250565b600061a8fd601783618fb5565b915061a9088261a8c7565b602082019050919050565b6000602082019050818103600083015261a92c8161a8f0565b9050919050565b7f636f6d6d69746d656e74206d75737420626520696e2070617374000000000000600082015250565b600061a969601a83618fb5565b915061a9748261a933565b602082019050919050565b6000602082019050818103600083015261a9988161a95c565b9050919050565b61a9a881619724565b82525050565b600060208201905061a9c3600083018461a99f565b92915050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b600061aa25602683618fb5565b915061aa308261a9c9565b604082019050919050565b6000602082019050818103600083015261aa548161aa18565b9050919050565b600061aa668261825d565b915067ffffffffffffffff820361aa805761aa7f618f10565b5b600182019050919050565b600061aa968261825d565b915061aaa18361825d565b9250828203905067ffffffffffffffff81111561aac15761aac0618f10565b5b92915050565b600060608201905061aadc6000830186617ee8565b61aae96020830185617ee8565b61aaf66040830184617ee8565b94935050505056fea2646970667358221220b6131463cc38c3956d894312a407fb29ed646fda486c28aeeaec54d9f168a5dc64736f6c63430008130033",
}

// EngineV5ABI is the input ABI used to generate the binding from.
// Deprecated: Use EngineV5MetaData.ABI instead.
var EngineV5ABI = EngineV5MetaData.ABI

// EngineV5Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EngineV5MetaData.Bin instead.
var EngineV5Bin = EngineV5MetaData.Bin

// DeployEngineV5 deploys a new Ethereum contract, binding an instance of EngineV5 to it.
func DeployEngineV5(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EngineV5, error) {
	parsed, err := EngineV5MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EngineV5Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EngineV5{EngineV5Caller: EngineV5Caller{contract: contract}, EngineV5Transactor: EngineV5Transactor{contract: contract}, EngineV5Filterer: EngineV5Filterer{contract: contract}}, nil
}

// EngineV5 is an auto generated Go binding around an Ethereum contract.
type EngineV5 struct {
	EngineV5Caller     // Read-only binding to the contract
	EngineV5Transactor // Write-only binding to the contract
	EngineV5Filterer   // Log filterer for contract events
}

// EngineV5Caller is an auto generated read-only Go binding around an Ethereum contract.
type EngineV5Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EngineV5Transactor is an auto generated write-only Go binding around an Ethereum contract.
type EngineV5Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EngineV5Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EngineV5Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EngineV5Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EngineV5Session struct {
	Contract     *EngineV5         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EngineV5CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EngineV5CallerSession struct {
	Contract *EngineV5Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// EngineV5TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EngineV5TransactorSession struct {
	Contract     *EngineV5Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// EngineV5Raw is an auto generated low-level Go binding around an Ethereum contract.
type EngineV5Raw struct {
	Contract *EngineV5 // Generic contract binding to access the raw methods on
}

// EngineV5CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EngineV5CallerRaw struct {
	Contract *EngineV5Caller // Generic read-only contract binding to access the raw methods on
}

// EngineV5TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EngineV5TransactorRaw struct {
	Contract *EngineV5Transactor // Generic write-only contract binding to access the raw methods on
}

// NewEngineV5 creates a new instance of EngineV5, bound to a specific deployed contract.
func NewEngineV5(address common.Address, backend bind.ContractBackend) (*EngineV5, error) {
	contract, err := bindEngineV5(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EngineV5{EngineV5Caller: EngineV5Caller{contract: contract}, EngineV5Transactor: EngineV5Transactor{contract: contract}, EngineV5Filterer: EngineV5Filterer{contract: contract}}, nil
}

// NewEngineV5Caller creates a new read-only instance of EngineV5, bound to a specific deployed contract.
func NewEngineV5Caller(address common.Address, caller bind.ContractCaller) (*EngineV5Caller, error) {
	contract, err := bindEngineV5(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EngineV5Caller{contract: contract}, nil
}

// NewEngineV5Transactor creates a new write-only instance of EngineV5, bound to a specific deployed contract.
func NewEngineV5Transactor(address common.Address, transactor bind.ContractTransactor) (*EngineV5Transactor, error) {
	contract, err := bindEngineV5(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EngineV5Transactor{contract: contract}, nil
}

// NewEngineV5Filterer creates a new log filterer instance of EngineV5, bound to a specific deployed contract.
func NewEngineV5Filterer(address common.Address, filterer bind.ContractFilterer) (*EngineV5Filterer, error) {
	contract, err := bindEngineV5(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EngineV5Filterer{contract: contract}, nil
}

// bindEngineV5 binds a generic wrapper to an already deployed contract.
func bindEngineV5(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EngineV5MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EngineV5 *EngineV5Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EngineV5.Contract.EngineV5Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EngineV5 *EngineV5Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EngineV5.Contract.EngineV5Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EngineV5 *EngineV5Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EngineV5.Contract.EngineV5Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EngineV5 *EngineV5CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EngineV5.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EngineV5 *EngineV5TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EngineV5.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EngineV5 *EngineV5TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EngineV5.Contract.contract.Transact(opts, method, params...)
}

// AccruedFees is a free data retrieval call binding the contract method 0x682c2058.
//
// Solidity: function accruedFees() view returns(uint256)
func (_EngineV5 *EngineV5Caller) AccruedFees(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "accruedFees")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AccruedFees is a free data retrieval call binding the contract method 0x682c2058.
//
// Solidity: function accruedFees() view returns(uint256)
func (_EngineV5 *EngineV5Session) AccruedFees() (*big.Int, error) {
	return _EngineV5.Contract.AccruedFees(&_EngineV5.CallOpts)
}

// AccruedFees is a free data retrieval call binding the contract method 0x682c2058.
//
// Solidity: function accruedFees() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) AccruedFees() (*big.Int, error) {
	return _EngineV5.Contract.AccruedFees(&_EngineV5.CallOpts)
}

// BaseToken is a free data retrieval call binding the contract method 0xc55dae63.
//
// Solidity: function baseToken() view returns(address)
func (_EngineV5 *EngineV5Caller) BaseToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "baseToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BaseToken is a free data retrieval call binding the contract method 0xc55dae63.
//
// Solidity: function baseToken() view returns(address)
func (_EngineV5 *EngineV5Session) BaseToken() (common.Address, error) {
	return _EngineV5.Contract.BaseToken(&_EngineV5.CallOpts)
}

// BaseToken is a free data retrieval call binding the contract method 0xc55dae63.
//
// Solidity: function baseToken() view returns(address)
func (_EngineV5 *EngineV5CallerSession) BaseToken() (common.Address, error) {
	return _EngineV5.Contract.BaseToken(&_EngineV5.CallOpts)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(uint256)
func (_EngineV5 *EngineV5Caller) Commitments(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "commitments", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(uint256)
func (_EngineV5 *EngineV5Session) Commitments(arg0 [32]byte) (*big.Int, error) {
	return _EngineV5.Contract.Commitments(&_EngineV5.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) Commitments(arg0 [32]byte) (*big.Int, error) {
	return _EngineV5.Contract.Commitments(&_EngineV5.CallOpts, arg0)
}

// ContestationVoteExtensionTime is a free data retrieval call binding the contract method 0xa2492a90.
//
// Solidity: function contestationVoteExtensionTime() view returns(uint256)
func (_EngineV5 *EngineV5Caller) ContestationVoteExtensionTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "contestationVoteExtensionTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ContestationVoteExtensionTime is a free data retrieval call binding the contract method 0xa2492a90.
//
// Solidity: function contestationVoteExtensionTime() view returns(uint256)
func (_EngineV5 *EngineV5Session) ContestationVoteExtensionTime() (*big.Int, error) {
	return _EngineV5.Contract.ContestationVoteExtensionTime(&_EngineV5.CallOpts)
}

// ContestationVoteExtensionTime is a free data retrieval call binding the contract method 0xa2492a90.
//
// Solidity: function contestationVoteExtensionTime() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) ContestationVoteExtensionTime() (*big.Int, error) {
	return _EngineV5.Contract.ContestationVoteExtensionTime(&_EngineV5.CallOpts)
}

// ContestationVoteNays is a free data retrieval call binding the contract method 0x303fb0d6.
//
// Solidity: function contestationVoteNays(bytes32 , uint256 ) view returns(address)
func (_EngineV5 *EngineV5Caller) ContestationVoteNays(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "contestationVoteNays", arg0, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ContestationVoteNays is a free data retrieval call binding the contract method 0x303fb0d6.
//
// Solidity: function contestationVoteNays(bytes32 , uint256 ) view returns(address)
func (_EngineV5 *EngineV5Session) ContestationVoteNays(arg0 [32]byte, arg1 *big.Int) (common.Address, error) {
	return _EngineV5.Contract.ContestationVoteNays(&_EngineV5.CallOpts, arg0, arg1)
}

// ContestationVoteNays is a free data retrieval call binding the contract method 0x303fb0d6.
//
// Solidity: function contestationVoteNays(bytes32 , uint256 ) view returns(address)
func (_EngineV5 *EngineV5CallerSession) ContestationVoteNays(arg0 [32]byte, arg1 *big.Int) (common.Address, error) {
	return _EngineV5.Contract.ContestationVoteNays(&_EngineV5.CallOpts, arg0, arg1)
}

// ContestationVoteYeas is a free data retrieval call binding the contract method 0xd1f0c941.
//
// Solidity: function contestationVoteYeas(bytes32 , uint256 ) view returns(address)
func (_EngineV5 *EngineV5Caller) ContestationVoteYeas(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "contestationVoteYeas", arg0, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ContestationVoteYeas is a free data retrieval call binding the contract method 0xd1f0c941.
//
// Solidity: function contestationVoteYeas(bytes32 , uint256 ) view returns(address)
func (_EngineV5 *EngineV5Session) ContestationVoteYeas(arg0 [32]byte, arg1 *big.Int) (common.Address, error) {
	return _EngineV5.Contract.ContestationVoteYeas(&_EngineV5.CallOpts, arg0, arg1)
}

// ContestationVoteYeas is a free data retrieval call binding the contract method 0xd1f0c941.
//
// Solidity: function contestationVoteYeas(bytes32 , uint256 ) view returns(address)
func (_EngineV5 *EngineV5CallerSession) ContestationVoteYeas(arg0 [32]byte, arg1 *big.Int) (common.Address, error) {
	return _EngineV5.Contract.ContestationVoteYeas(&_EngineV5.CallOpts, arg0, arg1)
}

// ContestationVoted is a free data retrieval call binding the contract method 0xd2780940.
//
// Solidity: function contestationVoted(bytes32 , address ) view returns(bool)
func (_EngineV5 *EngineV5Caller) ContestationVoted(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "contestationVoted", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ContestationVoted is a free data retrieval call binding the contract method 0xd2780940.
//
// Solidity: function contestationVoted(bytes32 , address ) view returns(bool)
func (_EngineV5 *EngineV5Session) ContestationVoted(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _EngineV5.Contract.ContestationVoted(&_EngineV5.CallOpts, arg0, arg1)
}

// ContestationVoted is a free data retrieval call binding the contract method 0xd2780940.
//
// Solidity: function contestationVoted(bytes32 , address ) view returns(bool)
func (_EngineV5 *EngineV5CallerSession) ContestationVoted(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _EngineV5.Contract.ContestationVoted(&_EngineV5.CallOpts, arg0, arg1)
}

// ContestationVotedIndex is a free data retrieval call binding the contract method 0x17f3e041.
//
// Solidity: function contestationVotedIndex(bytes32 ) view returns(uint256)
func (_EngineV5 *EngineV5Caller) ContestationVotedIndex(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "contestationVotedIndex", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ContestationVotedIndex is a free data retrieval call binding the contract method 0x17f3e041.
//
// Solidity: function contestationVotedIndex(bytes32 ) view returns(uint256)
func (_EngineV5 *EngineV5Session) ContestationVotedIndex(arg0 [32]byte) (*big.Int, error) {
	return _EngineV5.Contract.ContestationVotedIndex(&_EngineV5.CallOpts, arg0)
}

// ContestationVotedIndex is a free data retrieval call binding the contract method 0x17f3e041.
//
// Solidity: function contestationVotedIndex(bytes32 ) view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) ContestationVotedIndex(arg0 [32]byte) (*big.Int, error) {
	return _EngineV5.Contract.ContestationVotedIndex(&_EngineV5.CallOpts, arg0)
}

// Contestations is a free data retrieval call binding the contract method 0xd33b2ef5.
//
// Solidity: function contestations(bytes32 ) view returns(address validator, uint64 blocktime, uint32 finish_start_index, uint256 slashAmount)
func (_EngineV5 *EngineV5Caller) Contestations(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Validator        common.Address
	Blocktime        uint64
	FinishStartIndex uint32
	SlashAmount      *big.Int
}, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "contestations", arg0)

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
func (_EngineV5 *EngineV5Session) Contestations(arg0 [32]byte) (struct {
	Validator        common.Address
	Blocktime        uint64
	FinishStartIndex uint32
	SlashAmount      *big.Int
}, error) {
	return _EngineV5.Contract.Contestations(&_EngineV5.CallOpts, arg0)
}

// Contestations is a free data retrieval call binding the contract method 0xd33b2ef5.
//
// Solidity: function contestations(bytes32 ) view returns(address validator, uint64 blocktime, uint32 finish_start_index, uint256 slashAmount)
func (_EngineV5 *EngineV5CallerSession) Contestations(arg0 [32]byte) (struct {
	Validator        common.Address
	Blocktime        uint64
	FinishStartIndex uint32
	SlashAmount      *big.Int
}, error) {
	return _EngineV5.Contract.Contestations(&_EngineV5.CallOpts, arg0)
}

// DiffMul is a free data retrieval call binding the contract method 0x1f88ea1c.
//
// Solidity: function diffMul(uint256 t, uint256 ts) pure returns(uint256)
func (_EngineV5 *EngineV5Caller) DiffMul(opts *bind.CallOpts, t *big.Int, ts *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "diffMul", t, ts)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DiffMul is a free data retrieval call binding the contract method 0x1f88ea1c.
//
// Solidity: function diffMul(uint256 t, uint256 ts) pure returns(uint256)
func (_EngineV5 *EngineV5Session) DiffMul(t *big.Int, ts *big.Int) (*big.Int, error) {
	return _EngineV5.Contract.DiffMul(&_EngineV5.CallOpts, t, ts)
}

// DiffMul is a free data retrieval call binding the contract method 0x1f88ea1c.
//
// Solidity: function diffMul(uint256 t, uint256 ts) pure returns(uint256)
func (_EngineV5 *EngineV5CallerSession) DiffMul(t *big.Int, ts *big.Int) (*big.Int, error) {
	return _EngineV5.Contract.DiffMul(&_EngineV5.CallOpts, t, ts)
}

// ExitValidatorMinUnlockTime is a free data retrieval call binding the contract method 0xa53e2525.
//
// Solidity: function exitValidatorMinUnlockTime() view returns(uint256)
func (_EngineV5 *EngineV5Caller) ExitValidatorMinUnlockTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "exitValidatorMinUnlockTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ExitValidatorMinUnlockTime is a free data retrieval call binding the contract method 0xa53e2525.
//
// Solidity: function exitValidatorMinUnlockTime() view returns(uint256)
func (_EngineV5 *EngineV5Session) ExitValidatorMinUnlockTime() (*big.Int, error) {
	return _EngineV5.Contract.ExitValidatorMinUnlockTime(&_EngineV5.CallOpts)
}

// ExitValidatorMinUnlockTime is a free data retrieval call binding the contract method 0xa53e2525.
//
// Solidity: function exitValidatorMinUnlockTime() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) ExitValidatorMinUnlockTime() (*big.Int, error) {
	return _EngineV5.Contract.ExitValidatorMinUnlockTime(&_EngineV5.CallOpts)
}

// GenerateCommitment is a free data retrieval call binding the contract method 0x393cb1c7.
//
// Solidity: function generateCommitment(address sender_, bytes32 taskid_, bytes cid_) pure returns(bytes32)
func (_EngineV5 *EngineV5Caller) GenerateCommitment(opts *bind.CallOpts, sender_ common.Address, taskid_ [32]byte, cid_ []byte) ([32]byte, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "generateCommitment", sender_, taskid_, cid_)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GenerateCommitment is a free data retrieval call binding the contract method 0x393cb1c7.
//
// Solidity: function generateCommitment(address sender_, bytes32 taskid_, bytes cid_) pure returns(bytes32)
func (_EngineV5 *EngineV5Session) GenerateCommitment(sender_ common.Address, taskid_ [32]byte, cid_ []byte) ([32]byte, error) {
	return _EngineV5.Contract.GenerateCommitment(&_EngineV5.CallOpts, sender_, taskid_, cid_)
}

// GenerateCommitment is a free data retrieval call binding the contract method 0x393cb1c7.
//
// Solidity: function generateCommitment(address sender_, bytes32 taskid_, bytes cid_) pure returns(bytes32)
func (_EngineV5 *EngineV5CallerSession) GenerateCommitment(sender_ common.Address, taskid_ [32]byte, cid_ []byte) ([32]byte, error) {
	return _EngineV5.Contract.GenerateCommitment(&_EngineV5.CallOpts, sender_, taskid_, cid_)
}

// GenerateIPFSCID is a free data retrieval call binding the contract method 0x40e8c56d.
//
// Solidity: function generateIPFSCID(bytes content_) pure returns(bytes)
func (_EngineV5 *EngineV5Caller) GenerateIPFSCID(opts *bind.CallOpts, content_ []byte) ([]byte, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "generateIPFSCID", content_)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GenerateIPFSCID is a free data retrieval call binding the contract method 0x40e8c56d.
//
// Solidity: function generateIPFSCID(bytes content_) pure returns(bytes)
func (_EngineV5 *EngineV5Session) GenerateIPFSCID(content_ []byte) ([]byte, error) {
	return _EngineV5.Contract.GenerateIPFSCID(&_EngineV5.CallOpts, content_)
}

// GenerateIPFSCID is a free data retrieval call binding the contract method 0x40e8c56d.
//
// Solidity: function generateIPFSCID(bytes content_) pure returns(bytes)
func (_EngineV5 *EngineV5CallerSession) GenerateIPFSCID(content_ []byte) ([]byte, error) {
	return _EngineV5.Contract.GenerateIPFSCID(&_EngineV5.CallOpts, content_)
}

// GetPsuedoTotalSupply is a free data retrieval call binding the contract method 0x7881c5e6.
//
// Solidity: function getPsuedoTotalSupply() view returns(uint256)
func (_EngineV5 *EngineV5Caller) GetPsuedoTotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "getPsuedoTotalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPsuedoTotalSupply is a free data retrieval call binding the contract method 0x7881c5e6.
//
// Solidity: function getPsuedoTotalSupply() view returns(uint256)
func (_EngineV5 *EngineV5Session) GetPsuedoTotalSupply() (*big.Int, error) {
	return _EngineV5.Contract.GetPsuedoTotalSupply(&_EngineV5.CallOpts)
}

// GetPsuedoTotalSupply is a free data retrieval call binding the contract method 0x7881c5e6.
//
// Solidity: function getPsuedoTotalSupply() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) GetPsuedoTotalSupply() (*big.Int, error) {
	return _EngineV5.Contract.GetPsuedoTotalSupply(&_EngineV5.CallOpts)
}

// GetReward is a free data retrieval call binding the contract method 0x3d18b912.
//
// Solidity: function getReward() view returns(uint256)
func (_EngineV5 *EngineV5Caller) GetReward(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "getReward")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetReward is a free data retrieval call binding the contract method 0x3d18b912.
//
// Solidity: function getReward() view returns(uint256)
func (_EngineV5 *EngineV5Session) GetReward() (*big.Int, error) {
	return _EngineV5.Contract.GetReward(&_EngineV5.CallOpts)
}

// GetReward is a free data retrieval call binding the contract method 0x3d18b912.
//
// Solidity: function getReward() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) GetReward() (*big.Int, error) {
	return _EngineV5.Contract.GetReward(&_EngineV5.CallOpts)
}

// GetSlashAmount is a free data retrieval call binding the contract method 0x3d57f5d9.
//
// Solidity: function getSlashAmount() view returns(uint256)
func (_EngineV5 *EngineV5Caller) GetSlashAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "getSlashAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetSlashAmount is a free data retrieval call binding the contract method 0x3d57f5d9.
//
// Solidity: function getSlashAmount() view returns(uint256)
func (_EngineV5 *EngineV5Session) GetSlashAmount() (*big.Int, error) {
	return _EngineV5.Contract.GetSlashAmount(&_EngineV5.CallOpts)
}

// GetSlashAmount is a free data retrieval call binding the contract method 0x3d57f5d9.
//
// Solidity: function getSlashAmount() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) GetSlashAmount() (*big.Int, error) {
	return _EngineV5.Contract.GetSlashAmount(&_EngineV5.CallOpts)
}

// GetValidatorMinimum is a free data retrieval call binding the contract method 0x2258d105.
//
// Solidity: function getValidatorMinimum() view returns(uint256)
func (_EngineV5 *EngineV5Caller) GetValidatorMinimum(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "getValidatorMinimum")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetValidatorMinimum is a free data retrieval call binding the contract method 0x2258d105.
//
// Solidity: function getValidatorMinimum() view returns(uint256)
func (_EngineV5 *EngineV5Session) GetValidatorMinimum() (*big.Int, error) {
	return _EngineV5.Contract.GetValidatorMinimum(&_EngineV5.CallOpts)
}

// GetValidatorMinimum is a free data retrieval call binding the contract method 0x2258d105.
//
// Solidity: function getValidatorMinimum() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) GetValidatorMinimum() (*big.Int, error) {
	return _EngineV5.Contract.GetValidatorMinimum(&_EngineV5.CallOpts)
}

// HashModel is a free data retrieval call binding the contract method 0x218a3048.
//
// Solidity: function hashModel((uint256,address,uint256,bytes) o_, address sender_) pure returns(bytes32)
func (_EngineV5 *EngineV5Caller) HashModel(opts *bind.CallOpts, o_ Model, sender_ common.Address) ([32]byte, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "hashModel", o_, sender_)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashModel is a free data retrieval call binding the contract method 0x218a3048.
//
// Solidity: function hashModel((uint256,address,uint256,bytes) o_, address sender_) pure returns(bytes32)
func (_EngineV5 *EngineV5Session) HashModel(o_ Model, sender_ common.Address) ([32]byte, error) {
	return _EngineV5.Contract.HashModel(&_EngineV5.CallOpts, o_, sender_)
}

// HashModel is a free data retrieval call binding the contract method 0x218a3048.
//
// Solidity: function hashModel((uint256,address,uint256,bytes) o_, address sender_) pure returns(bytes32)
func (_EngineV5 *EngineV5CallerSession) HashModel(o_ Model, sender_ common.Address) ([32]byte, error) {
	return _EngineV5.Contract.HashModel(&_EngineV5.CallOpts, o_, sender_)
}

// HashTask is a free data retrieval call binding the contract method 0x1466b63a.
//
// Solidity: function hashTask((bytes32,uint256,address,uint64,uint8,bytes) o_, address sender_, bytes32 prevhash_) pure returns(bytes32)
func (_EngineV5 *EngineV5Caller) HashTask(opts *bind.CallOpts, o_ Task, sender_ common.Address, prevhash_ [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "hashTask", o_, sender_, prevhash_)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashTask is a free data retrieval call binding the contract method 0x1466b63a.
//
// Solidity: function hashTask((bytes32,uint256,address,uint64,uint8,bytes) o_, address sender_, bytes32 prevhash_) pure returns(bytes32)
func (_EngineV5 *EngineV5Session) HashTask(o_ Task, sender_ common.Address, prevhash_ [32]byte) ([32]byte, error) {
	return _EngineV5.Contract.HashTask(&_EngineV5.CallOpts, o_, sender_, prevhash_)
}

// HashTask is a free data retrieval call binding the contract method 0x1466b63a.
//
// Solidity: function hashTask((bytes32,uint256,address,uint64,uint8,bytes) o_, address sender_, bytes32 prevhash_) pure returns(bytes32)
func (_EngineV5 *EngineV5CallerSession) HashTask(o_ Task, sender_ common.Address, prevhash_ [32]byte) ([32]byte, error) {
	return _EngineV5.Contract.HashTask(&_EngineV5.CallOpts, o_, sender_, prevhash_)
}

// LastContestationLossTime is a free data retrieval call binding the contract method 0x218e6859.
//
// Solidity: function lastContestationLossTime(address ) view returns(uint256)
func (_EngineV5 *EngineV5Caller) LastContestationLossTime(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "lastContestationLossTime", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastContestationLossTime is a free data retrieval call binding the contract method 0x218e6859.
//
// Solidity: function lastContestationLossTime(address ) view returns(uint256)
func (_EngineV5 *EngineV5Session) LastContestationLossTime(arg0 common.Address) (*big.Int, error) {
	return _EngineV5.Contract.LastContestationLossTime(&_EngineV5.CallOpts, arg0)
}

// LastContestationLossTime is a free data retrieval call binding the contract method 0x218e6859.
//
// Solidity: function lastContestationLossTime(address ) view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) LastContestationLossTime(arg0 common.Address) (*big.Int, error) {
	return _EngineV5.Contract.LastContestationLossTime(&_EngineV5.CallOpts, arg0)
}

// LastSolutionSubmission is a free data retrieval call binding the contract method 0xc24b5631.
//
// Solidity: function lastSolutionSubmission(address ) view returns(uint256)
func (_EngineV5 *EngineV5Caller) LastSolutionSubmission(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "lastSolutionSubmission", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastSolutionSubmission is a free data retrieval call binding the contract method 0xc24b5631.
//
// Solidity: function lastSolutionSubmission(address ) view returns(uint256)
func (_EngineV5 *EngineV5Session) LastSolutionSubmission(arg0 common.Address) (*big.Int, error) {
	return _EngineV5.Contract.LastSolutionSubmission(&_EngineV5.CallOpts, arg0)
}

// LastSolutionSubmission is a free data retrieval call binding the contract method 0xc24b5631.
//
// Solidity: function lastSolutionSubmission(address ) view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) LastSolutionSubmission(arg0 common.Address) (*big.Int, error) {
	return _EngineV5.Contract.LastSolutionSubmission(&_EngineV5.CallOpts, arg0)
}

// MaxContestationValidatorStakeSince is a free data retrieval call binding the contract method 0x8e6d86fd.
//
// Solidity: function maxContestationValidatorStakeSince() view returns(uint256)
func (_EngineV5 *EngineV5Caller) MaxContestationValidatorStakeSince(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "maxContestationValidatorStakeSince")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxContestationValidatorStakeSince is a free data retrieval call binding the contract method 0x8e6d86fd.
//
// Solidity: function maxContestationValidatorStakeSince() view returns(uint256)
func (_EngineV5 *EngineV5Session) MaxContestationValidatorStakeSince() (*big.Int, error) {
	return _EngineV5.Contract.MaxContestationValidatorStakeSince(&_EngineV5.CallOpts)
}

// MaxContestationValidatorStakeSince is a free data retrieval call binding the contract method 0x8e6d86fd.
//
// Solidity: function maxContestationValidatorStakeSince() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) MaxContestationValidatorStakeSince() (*big.Int, error) {
	return _EngineV5.Contract.MaxContestationValidatorStakeSince(&_EngineV5.CallOpts)
}

// MinClaimSolutionTime is a free data retrieval call binding the contract method 0x92809444.
//
// Solidity: function minClaimSolutionTime() view returns(uint256)
func (_EngineV5 *EngineV5Caller) MinClaimSolutionTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "minClaimSolutionTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinClaimSolutionTime is a free data retrieval call binding the contract method 0x92809444.
//
// Solidity: function minClaimSolutionTime() view returns(uint256)
func (_EngineV5 *EngineV5Session) MinClaimSolutionTime() (*big.Int, error) {
	return _EngineV5.Contract.MinClaimSolutionTime(&_EngineV5.CallOpts)
}

// MinClaimSolutionTime is a free data retrieval call binding the contract method 0x92809444.
//
// Solidity: function minClaimSolutionTime() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) MinClaimSolutionTime() (*big.Int, error) {
	return _EngineV5.Contract.MinClaimSolutionTime(&_EngineV5.CallOpts)
}

// MinContestationVotePeriodTime is a free data retrieval call binding the contract method 0x7b36006a.
//
// Solidity: function minContestationVotePeriodTime() view returns(uint256)
func (_EngineV5 *EngineV5Caller) MinContestationVotePeriodTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "minContestationVotePeriodTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinContestationVotePeriodTime is a free data retrieval call binding the contract method 0x7b36006a.
//
// Solidity: function minContestationVotePeriodTime() view returns(uint256)
func (_EngineV5 *EngineV5Session) MinContestationVotePeriodTime() (*big.Int, error) {
	return _EngineV5.Contract.MinContestationVotePeriodTime(&_EngineV5.CallOpts)
}

// MinContestationVotePeriodTime is a free data retrieval call binding the contract method 0x7b36006a.
//
// Solidity: function minContestationVotePeriodTime() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) MinContestationVotePeriodTime() (*big.Int, error) {
	return _EngineV5.Contract.MinContestationVotePeriodTime(&_EngineV5.CallOpts)
}

// MinRetractionWaitTime is a free data retrieval call binding the contract method 0x00fd7082.
//
// Solidity: function minRetractionWaitTime() view returns(uint256)
func (_EngineV5 *EngineV5Caller) MinRetractionWaitTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "minRetractionWaitTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinRetractionWaitTime is a free data retrieval call binding the contract method 0x00fd7082.
//
// Solidity: function minRetractionWaitTime() view returns(uint256)
func (_EngineV5 *EngineV5Session) MinRetractionWaitTime() (*big.Int, error) {
	return _EngineV5.Contract.MinRetractionWaitTime(&_EngineV5.CallOpts)
}

// MinRetractionWaitTime is a free data retrieval call binding the contract method 0x00fd7082.
//
// Solidity: function minRetractionWaitTime() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) MinRetractionWaitTime() (*big.Int, error) {
	return _EngineV5.Contract.MinRetractionWaitTime(&_EngineV5.CallOpts)
}

// Models is a free data retrieval call binding the contract method 0xe236f46b.
//
// Solidity: function models(bytes32 ) view returns(uint256 fee, address addr, uint256 rate, bytes cid)
func (_EngineV5 *EngineV5Caller) Models(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Fee  *big.Int
	Addr common.Address
	Rate *big.Int
	Cid  []byte
}, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "models", arg0)

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
func (_EngineV5 *EngineV5Session) Models(arg0 [32]byte) (struct {
	Fee  *big.Int
	Addr common.Address
	Rate *big.Int
	Cid  []byte
}, error) {
	return _EngineV5.Contract.Models(&_EngineV5.CallOpts, arg0)
}

// Models is a free data retrieval call binding the contract method 0xe236f46b.
//
// Solidity: function models(bytes32 ) view returns(uint256 fee, address addr, uint256 rate, bytes cid)
func (_EngineV5 *EngineV5CallerSession) Models(arg0 [32]byte) (struct {
	Fee  *big.Int
	Addr common.Address
	Rate *big.Int
	Cid  []byte
}, error) {
	return _EngineV5.Contract.Models(&_EngineV5.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_EngineV5 *EngineV5Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_EngineV5 *EngineV5Session) Owner() (common.Address, error) {
	return _EngineV5.Contract.Owner(&_EngineV5.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_EngineV5 *EngineV5CallerSession) Owner() (common.Address, error) {
	return _EngineV5.Contract.Owner(&_EngineV5.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_EngineV5 *EngineV5Caller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_EngineV5 *EngineV5Session) Paused() (bool, error) {
	return _EngineV5.Contract.Paused(&_EngineV5.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_EngineV5 *EngineV5CallerSession) Paused() (bool, error) {
	return _EngineV5.Contract.Paused(&_EngineV5.CallOpts)
}

// Pauser is a free data retrieval call binding the contract method 0x9fd0506d.
//
// Solidity: function pauser() view returns(address)
func (_EngineV5 *EngineV5Caller) Pauser(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "pauser")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Pauser is a free data retrieval call binding the contract method 0x9fd0506d.
//
// Solidity: function pauser() view returns(address)
func (_EngineV5 *EngineV5Session) Pauser() (common.Address, error) {
	return _EngineV5.Contract.Pauser(&_EngineV5.CallOpts)
}

// Pauser is a free data retrieval call binding the contract method 0x9fd0506d.
//
// Solidity: function pauser() view returns(address)
func (_EngineV5 *EngineV5CallerSession) Pauser() (common.Address, error) {
	return _EngineV5.Contract.Pauser(&_EngineV5.CallOpts)
}

// PendingValidatorWithdrawRequests is a free data retrieval call binding the contract method 0xd2992baa.
//
// Solidity: function pendingValidatorWithdrawRequests(address , uint256 ) view returns(uint256 unlockTime, uint256 amount)
func (_EngineV5 *EngineV5Caller) PendingValidatorWithdrawRequests(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (struct {
	UnlockTime *big.Int
	Amount     *big.Int
}, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "pendingValidatorWithdrawRequests", arg0, arg1)

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
func (_EngineV5 *EngineV5Session) PendingValidatorWithdrawRequests(arg0 common.Address, arg1 *big.Int) (struct {
	UnlockTime *big.Int
	Amount     *big.Int
}, error) {
	return _EngineV5.Contract.PendingValidatorWithdrawRequests(&_EngineV5.CallOpts, arg0, arg1)
}

// PendingValidatorWithdrawRequests is a free data retrieval call binding the contract method 0xd2992baa.
//
// Solidity: function pendingValidatorWithdrawRequests(address , uint256 ) view returns(uint256 unlockTime, uint256 amount)
func (_EngineV5 *EngineV5CallerSession) PendingValidatorWithdrawRequests(arg0 common.Address, arg1 *big.Int) (struct {
	UnlockTime *big.Int
	Amount     *big.Int
}, error) {
	return _EngineV5.Contract.PendingValidatorWithdrawRequests(&_EngineV5.CallOpts, arg0, arg1)
}

// PendingValidatorWithdrawRequestsCount is a free data retrieval call binding the contract method 0xd2307ae4.
//
// Solidity: function pendingValidatorWithdrawRequestsCount(address ) view returns(uint256)
func (_EngineV5 *EngineV5Caller) PendingValidatorWithdrawRequestsCount(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "pendingValidatorWithdrawRequestsCount", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PendingValidatorWithdrawRequestsCount is a free data retrieval call binding the contract method 0xd2307ae4.
//
// Solidity: function pendingValidatorWithdrawRequestsCount(address ) view returns(uint256)
func (_EngineV5 *EngineV5Session) PendingValidatorWithdrawRequestsCount(arg0 common.Address) (*big.Int, error) {
	return _EngineV5.Contract.PendingValidatorWithdrawRequestsCount(&_EngineV5.CallOpts, arg0)
}

// PendingValidatorWithdrawRequestsCount is a free data retrieval call binding the contract method 0xd2307ae4.
//
// Solidity: function pendingValidatorWithdrawRequestsCount(address ) view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) PendingValidatorWithdrawRequestsCount(arg0 common.Address) (*big.Int, error) {
	return _EngineV5.Contract.PendingValidatorWithdrawRequestsCount(&_EngineV5.CallOpts, arg0)
}

// Prevhash is a free data retrieval call binding the contract method 0xc17ddb2a.
//
// Solidity: function prevhash() view returns(bytes32)
func (_EngineV5 *EngineV5Caller) Prevhash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "prevhash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Prevhash is a free data retrieval call binding the contract method 0xc17ddb2a.
//
// Solidity: function prevhash() view returns(bytes32)
func (_EngineV5 *EngineV5Session) Prevhash() ([32]byte, error) {
	return _EngineV5.Contract.Prevhash(&_EngineV5.CallOpts)
}

// Prevhash is a free data retrieval call binding the contract method 0xc17ddb2a.
//
// Solidity: function prevhash() view returns(bytes32)
func (_EngineV5 *EngineV5CallerSession) Prevhash() ([32]byte, error) {
	return _EngineV5.Contract.Prevhash(&_EngineV5.CallOpts)
}

// RetractionFeePercentage is a free data retrieval call binding the contract method 0x72dc0ee1.
//
// Solidity: function retractionFeePercentage() view returns(uint256)
func (_EngineV5 *EngineV5Caller) RetractionFeePercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "retractionFeePercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RetractionFeePercentage is a free data retrieval call binding the contract method 0x72dc0ee1.
//
// Solidity: function retractionFeePercentage() view returns(uint256)
func (_EngineV5 *EngineV5Session) RetractionFeePercentage() (*big.Int, error) {
	return _EngineV5.Contract.RetractionFeePercentage(&_EngineV5.CallOpts)
}

// RetractionFeePercentage is a free data retrieval call binding the contract method 0x72dc0ee1.
//
// Solidity: function retractionFeePercentage() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) RetractionFeePercentage() (*big.Int, error) {
	return _EngineV5.Contract.RetractionFeePercentage(&_EngineV5.CallOpts)
}

// Reward is a free data retrieval call binding the contract method 0xa4fa8d57.
//
// Solidity: function reward(uint256 t, uint256 ts) pure returns(uint256)
func (_EngineV5 *EngineV5Caller) Reward(opts *bind.CallOpts, t *big.Int, ts *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "reward", t, ts)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Reward is a free data retrieval call binding the contract method 0xa4fa8d57.
//
// Solidity: function reward(uint256 t, uint256 ts) pure returns(uint256)
func (_EngineV5 *EngineV5Session) Reward(t *big.Int, ts *big.Int) (*big.Int, error) {
	return _EngineV5.Contract.Reward(&_EngineV5.CallOpts, t, ts)
}

// Reward is a free data retrieval call binding the contract method 0xa4fa8d57.
//
// Solidity: function reward(uint256 t, uint256 ts) pure returns(uint256)
func (_EngineV5 *EngineV5CallerSession) Reward(t *big.Int, ts *big.Int) (*big.Int, error) {
	return _EngineV5.Contract.Reward(&_EngineV5.CallOpts, t, ts)
}

// SlashAmountPercentage is a free data retrieval call binding the contract method 0xdc06a89f.
//
// Solidity: function slashAmountPercentage() view returns(uint256)
func (_EngineV5 *EngineV5Caller) SlashAmountPercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "slashAmountPercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SlashAmountPercentage is a free data retrieval call binding the contract method 0xdc06a89f.
//
// Solidity: function slashAmountPercentage() view returns(uint256)
func (_EngineV5 *EngineV5Session) SlashAmountPercentage() (*big.Int, error) {
	return _EngineV5.Contract.SlashAmountPercentage(&_EngineV5.CallOpts)
}

// SlashAmountPercentage is a free data retrieval call binding the contract method 0xdc06a89f.
//
// Solidity: function slashAmountPercentage() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) SlashAmountPercentage() (*big.Int, error) {
	return _EngineV5.Contract.SlashAmountPercentage(&_EngineV5.CallOpts)
}

// SolutionFeePercentage is a free data retrieval call binding the contract method 0xf1b8989d.
//
// Solidity: function solutionFeePercentage() view returns(uint256)
func (_EngineV5 *EngineV5Caller) SolutionFeePercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "solutionFeePercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SolutionFeePercentage is a free data retrieval call binding the contract method 0xf1b8989d.
//
// Solidity: function solutionFeePercentage() view returns(uint256)
func (_EngineV5 *EngineV5Session) SolutionFeePercentage() (*big.Int, error) {
	return _EngineV5.Contract.SolutionFeePercentage(&_EngineV5.CallOpts)
}

// SolutionFeePercentage is a free data retrieval call binding the contract method 0xf1b8989d.
//
// Solidity: function solutionFeePercentage() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) SolutionFeePercentage() (*big.Int, error) {
	return _EngineV5.Contract.SolutionFeePercentage(&_EngineV5.CallOpts)
}

// SolutionRateLimit is a free data retrieval call binding the contract method 0xa1975adf.
//
// Solidity: function solutionRateLimit() view returns(uint256)
func (_EngineV5 *EngineV5Caller) SolutionRateLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "solutionRateLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SolutionRateLimit is a free data retrieval call binding the contract method 0xa1975adf.
//
// Solidity: function solutionRateLimit() view returns(uint256)
func (_EngineV5 *EngineV5Session) SolutionRateLimit() (*big.Int, error) {
	return _EngineV5.Contract.SolutionRateLimit(&_EngineV5.CallOpts)
}

// SolutionRateLimit is a free data retrieval call binding the contract method 0xa1975adf.
//
// Solidity: function solutionRateLimit() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) SolutionRateLimit() (*big.Int, error) {
	return _EngineV5.Contract.SolutionRateLimit(&_EngineV5.CallOpts)
}

// Solutions is a free data retrieval call binding the contract method 0x75c70509.
//
// Solidity: function solutions(bytes32 ) view returns(address validator, uint64 blocktime, bool claimed, bytes cid)
func (_EngineV5 *EngineV5Caller) Solutions(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Validator common.Address
	Blocktime uint64
	Claimed   bool
	Cid       []byte
}, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "solutions", arg0)

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
func (_EngineV5 *EngineV5Session) Solutions(arg0 [32]byte) (struct {
	Validator common.Address
	Blocktime uint64
	Claimed   bool
	Cid       []byte
}, error) {
	return _EngineV5.Contract.Solutions(&_EngineV5.CallOpts, arg0)
}

// Solutions is a free data retrieval call binding the contract method 0x75c70509.
//
// Solidity: function solutions(bytes32 ) view returns(address validator, uint64 blocktime, bool claimed, bytes cid)
func (_EngineV5 *EngineV5CallerSession) Solutions(arg0 [32]byte) (struct {
	Validator common.Address
	Blocktime uint64
	Claimed   bool
	Cid       []byte
}, error) {
	return _EngineV5.Contract.Solutions(&_EngineV5.CallOpts, arg0)
}

// SolutionsStake is a free data retrieval call binding the contract method 0xb4dc35b7.
//
// Solidity: function solutionsStake(bytes32 ) view returns(uint256)
func (_EngineV5 *EngineV5Caller) SolutionsStake(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "solutionsStake", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SolutionsStake is a free data retrieval call binding the contract method 0xb4dc35b7.
//
// Solidity: function solutionsStake(bytes32 ) view returns(uint256)
func (_EngineV5 *EngineV5Session) SolutionsStake(arg0 [32]byte) (*big.Int, error) {
	return _EngineV5.Contract.SolutionsStake(&_EngineV5.CallOpts, arg0)
}

// SolutionsStake is a free data retrieval call binding the contract method 0xb4dc35b7.
//
// Solidity: function solutionsStake(bytes32 ) view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) SolutionsStake(arg0 [32]byte) (*big.Int, error) {
	return _EngineV5.Contract.SolutionsStake(&_EngineV5.CallOpts, arg0)
}

// SolutionsStakeAmount is a free data retrieval call binding the contract method 0x9b975119.
//
// Solidity: function solutionsStakeAmount() view returns(uint256)
func (_EngineV5 *EngineV5Caller) SolutionsStakeAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "solutionsStakeAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SolutionsStakeAmount is a free data retrieval call binding the contract method 0x9b975119.
//
// Solidity: function solutionsStakeAmount() view returns(uint256)
func (_EngineV5 *EngineV5Session) SolutionsStakeAmount() (*big.Int, error) {
	return _EngineV5.Contract.SolutionsStakeAmount(&_EngineV5.CallOpts)
}

// SolutionsStakeAmount is a free data retrieval call binding the contract method 0x9b975119.
//
// Solidity: function solutionsStakeAmount() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) SolutionsStakeAmount() (*big.Int, error) {
	return _EngineV5.Contract.SolutionsStakeAmount(&_EngineV5.CallOpts)
}

// StartBlockTime is a free data retrieval call binding the contract method 0x0c18d4ce.
//
// Solidity: function startBlockTime() view returns(uint64)
func (_EngineV5 *EngineV5Caller) StartBlockTime(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "startBlockTime")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// StartBlockTime is a free data retrieval call binding the contract method 0x0c18d4ce.
//
// Solidity: function startBlockTime() view returns(uint64)
func (_EngineV5 *EngineV5Session) StartBlockTime() (uint64, error) {
	return _EngineV5.Contract.StartBlockTime(&_EngineV5.CallOpts)
}

// StartBlockTime is a free data retrieval call binding the contract method 0x0c18d4ce.
//
// Solidity: function startBlockTime() view returns(uint64)
func (_EngineV5 *EngineV5CallerSession) StartBlockTime() (uint64, error) {
	return _EngineV5.Contract.StartBlockTime(&_EngineV5.CallOpts)
}

// TargetTs is a free data retrieval call binding the contract method 0xcf596e45.
//
// Solidity: function targetTs(uint256 t) pure returns(uint256)
func (_EngineV5 *EngineV5Caller) TargetTs(opts *bind.CallOpts, t *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "targetTs", t)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TargetTs is a free data retrieval call binding the contract method 0xcf596e45.
//
// Solidity: function targetTs(uint256 t) pure returns(uint256)
func (_EngineV5 *EngineV5Session) TargetTs(t *big.Int) (*big.Int, error) {
	return _EngineV5.Contract.TargetTs(&_EngineV5.CallOpts, t)
}

// TargetTs is a free data retrieval call binding the contract method 0xcf596e45.
//
// Solidity: function targetTs(uint256 t) pure returns(uint256)
func (_EngineV5 *EngineV5CallerSession) TargetTs(t *big.Int) (*big.Int, error) {
	return _EngineV5.Contract.TargetTs(&_EngineV5.CallOpts, t)
}

// TaskOwnerRewardPercentage is a free data retrieval call binding the contract method 0xc1f72723.
//
// Solidity: function taskOwnerRewardPercentage() view returns(uint256)
func (_EngineV5 *EngineV5Caller) TaskOwnerRewardPercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "taskOwnerRewardPercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TaskOwnerRewardPercentage is a free data retrieval call binding the contract method 0xc1f72723.
//
// Solidity: function taskOwnerRewardPercentage() view returns(uint256)
func (_EngineV5 *EngineV5Session) TaskOwnerRewardPercentage() (*big.Int, error) {
	return _EngineV5.Contract.TaskOwnerRewardPercentage(&_EngineV5.CallOpts)
}

// TaskOwnerRewardPercentage is a free data retrieval call binding the contract method 0xc1f72723.
//
// Solidity: function taskOwnerRewardPercentage() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) TaskOwnerRewardPercentage() (*big.Int, error) {
	return _EngineV5.Contract.TaskOwnerRewardPercentage(&_EngineV5.CallOpts)
}

// Tasks is a free data retrieval call binding the contract method 0xe579f500.
//
// Solidity: function tasks(bytes32 ) view returns(bytes32 model, uint256 fee, address owner, uint64 blocktime, uint8 version, bytes cid)
func (_EngineV5 *EngineV5Caller) Tasks(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Model     [32]byte
	Fee       *big.Int
	Owner     common.Address
	Blocktime uint64
	Version   uint8
	Cid       []byte
}, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "tasks", arg0)

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
func (_EngineV5 *EngineV5Session) Tasks(arg0 [32]byte) (struct {
	Model     [32]byte
	Fee       *big.Int
	Owner     common.Address
	Blocktime uint64
	Version   uint8
	Cid       []byte
}, error) {
	return _EngineV5.Contract.Tasks(&_EngineV5.CallOpts, arg0)
}

// Tasks is a free data retrieval call binding the contract method 0xe579f500.
//
// Solidity: function tasks(bytes32 ) view returns(bytes32 model, uint256 fee, address owner, uint64 blocktime, uint8 version, bytes cid)
func (_EngineV5 *EngineV5CallerSession) Tasks(arg0 [32]byte) (struct {
	Model     [32]byte
	Fee       *big.Int
	Owner     common.Address
	Blocktime uint64
	Version   uint8
	Cid       []byte
}, error) {
	return _EngineV5.Contract.Tasks(&_EngineV5.CallOpts, arg0)
}

// TotalHeld is a free data retrieval call binding the contract method 0xf43cc773.
//
// Solidity: function totalHeld() view returns(uint256)
func (_EngineV5 *EngineV5Caller) TotalHeld(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "totalHeld")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalHeld is a free data retrieval call binding the contract method 0xf43cc773.
//
// Solidity: function totalHeld() view returns(uint256)
func (_EngineV5 *EngineV5Session) TotalHeld() (*big.Int, error) {
	return _EngineV5.Contract.TotalHeld(&_EngineV5.CallOpts)
}

// TotalHeld is a free data retrieval call binding the contract method 0xf43cc773.
//
// Solidity: function totalHeld() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) TotalHeld() (*big.Int, error) {
	return _EngineV5.Contract.TotalHeld(&_EngineV5.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_EngineV5 *EngineV5Caller) Treasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "treasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_EngineV5 *EngineV5Session) Treasury() (common.Address, error) {
	return _EngineV5.Contract.Treasury(&_EngineV5.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_EngineV5 *EngineV5CallerSession) Treasury() (common.Address, error) {
	return _EngineV5.Contract.Treasury(&_EngineV5.CallOpts)
}

// TreasuryRewardPercentage is a free data retrieval call binding the contract method 0xc31784be.
//
// Solidity: function treasuryRewardPercentage() view returns(uint256)
func (_EngineV5 *EngineV5Caller) TreasuryRewardPercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "treasuryRewardPercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TreasuryRewardPercentage is a free data retrieval call binding the contract method 0xc31784be.
//
// Solidity: function treasuryRewardPercentage() view returns(uint256)
func (_EngineV5 *EngineV5Session) TreasuryRewardPercentage() (*big.Int, error) {
	return _EngineV5.Contract.TreasuryRewardPercentage(&_EngineV5.CallOpts)
}

// TreasuryRewardPercentage is a free data retrieval call binding the contract method 0xc31784be.
//
// Solidity: function treasuryRewardPercentage() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) TreasuryRewardPercentage() (*big.Int, error) {
	return _EngineV5.Contract.TreasuryRewardPercentage(&_EngineV5.CallOpts)
}

// ValidatorCanVote is a free data retrieval call binding the contract method 0x83657795.
//
// Solidity: function validatorCanVote(address addr_, bytes32 taskid_) view returns(uint256)
func (_EngineV5 *EngineV5Caller) ValidatorCanVote(opts *bind.CallOpts, addr_ common.Address, taskid_ [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "validatorCanVote", addr_, taskid_)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorCanVote is a free data retrieval call binding the contract method 0x83657795.
//
// Solidity: function validatorCanVote(address addr_, bytes32 taskid_) view returns(uint256)
func (_EngineV5 *EngineV5Session) ValidatorCanVote(addr_ common.Address, taskid_ [32]byte) (*big.Int, error) {
	return _EngineV5.Contract.ValidatorCanVote(&_EngineV5.CallOpts, addr_, taskid_)
}

// ValidatorCanVote is a free data retrieval call binding the contract method 0x83657795.
//
// Solidity: function validatorCanVote(address addr_, bytes32 taskid_) view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) ValidatorCanVote(addr_ common.Address, taskid_ [32]byte) (*big.Int, error) {
	return _EngineV5.Contract.ValidatorCanVote(&_EngineV5.CallOpts, addr_, taskid_)
}

// ValidatorMinimumPercentage is a free data retrieval call binding the contract method 0x96bb02c3.
//
// Solidity: function validatorMinimumPercentage() view returns(uint256)
func (_EngineV5 *EngineV5Caller) ValidatorMinimumPercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "validatorMinimumPercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorMinimumPercentage is a free data retrieval call binding the contract method 0x96bb02c3.
//
// Solidity: function validatorMinimumPercentage() view returns(uint256)
func (_EngineV5 *EngineV5Session) ValidatorMinimumPercentage() (*big.Int, error) {
	return _EngineV5.Contract.ValidatorMinimumPercentage(&_EngineV5.CallOpts)
}

// ValidatorMinimumPercentage is a free data retrieval call binding the contract method 0x96bb02c3.
//
// Solidity: function validatorMinimumPercentage() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) ValidatorMinimumPercentage() (*big.Int, error) {
	return _EngineV5.Contract.ValidatorMinimumPercentage(&_EngineV5.CallOpts)
}

// ValidatorWithdrawPendingAmount is a free data retrieval call binding the contract method 0x1b75c43e.
//
// Solidity: function validatorWithdrawPendingAmount(address ) view returns(uint256)
func (_EngineV5 *EngineV5Caller) ValidatorWithdrawPendingAmount(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "validatorWithdrawPendingAmount", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorWithdrawPendingAmount is a free data retrieval call binding the contract method 0x1b75c43e.
//
// Solidity: function validatorWithdrawPendingAmount(address ) view returns(uint256)
func (_EngineV5 *EngineV5Session) ValidatorWithdrawPendingAmount(arg0 common.Address) (*big.Int, error) {
	return _EngineV5.Contract.ValidatorWithdrawPendingAmount(&_EngineV5.CallOpts, arg0)
}

// ValidatorWithdrawPendingAmount is a free data retrieval call binding the contract method 0x1b75c43e.
//
// Solidity: function validatorWithdrawPendingAmount(address ) view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) ValidatorWithdrawPendingAmount(arg0 common.Address) (*big.Int, error) {
	return _EngineV5.Contract.ValidatorWithdrawPendingAmount(&_EngineV5.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(uint256 staked, uint256 since, address addr)
func (_EngineV5 *EngineV5Caller) Validators(opts *bind.CallOpts, arg0 common.Address) (struct {
	Staked *big.Int
	Since  *big.Int
	Addr   common.Address
}, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "validators", arg0)

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
func (_EngineV5 *EngineV5Session) Validators(arg0 common.Address) (struct {
	Staked *big.Int
	Since  *big.Int
	Addr   common.Address
}, error) {
	return _EngineV5.Contract.Validators(&_EngineV5.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(uint256 staked, uint256 since, address addr)
func (_EngineV5 *EngineV5CallerSession) Validators(arg0 common.Address) (struct {
	Staked *big.Int
	Since  *big.Int
	Addr   common.Address
}, error) {
	return _EngineV5.Contract.Validators(&_EngineV5.CallOpts, arg0)
}

// VeRewards is a free data retrieval call binding the contract method 0x0d468d95.
//
// Solidity: function veRewards() view returns(uint256)
func (_EngineV5 *EngineV5Caller) VeRewards(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "veRewards")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VeRewards is a free data retrieval call binding the contract method 0x0d468d95.
//
// Solidity: function veRewards() view returns(uint256)
func (_EngineV5 *EngineV5Session) VeRewards() (*big.Int, error) {
	return _EngineV5.Contract.VeRewards(&_EngineV5.CallOpts)
}

// VeRewards is a free data retrieval call binding the contract method 0x0d468d95.
//
// Solidity: function veRewards() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) VeRewards() (*big.Int, error) {
	return _EngineV5.Contract.VeRewards(&_EngineV5.CallOpts)
}

// VeStaking is a free data retrieval call binding the contract method 0x82b5077f.
//
// Solidity: function veStaking() view returns(address)
func (_EngineV5 *EngineV5Caller) VeStaking(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "veStaking")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VeStaking is a free data retrieval call binding the contract method 0x82b5077f.
//
// Solidity: function veStaking() view returns(address)
func (_EngineV5 *EngineV5Session) VeStaking() (common.Address, error) {
	return _EngineV5.Contract.VeStaking(&_EngineV5.CallOpts)
}

// VeStaking is a free data retrieval call binding the contract method 0x82b5077f.
//
// Solidity: function veStaking() view returns(address)
func (_EngineV5 *EngineV5CallerSession) VeStaking() (common.Address, error) {
	return _EngineV5.Contract.VeStaking(&_EngineV5.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_EngineV5 *EngineV5Caller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_EngineV5 *EngineV5Session) Version() (*big.Int, error) {
	return _EngineV5.Contract.Version(&_EngineV5.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (_EngineV5 *EngineV5CallerSession) Version() (*big.Int, error) {
	return _EngineV5.Contract.Version(&_EngineV5.CallOpts)
}

// Voter is a free data retrieval call binding the contract method 0x46c96aac.
//
// Solidity: function voter() view returns(address)
func (_EngineV5 *EngineV5Caller) Voter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "voter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Voter is a free data retrieval call binding the contract method 0x46c96aac.
//
// Solidity: function voter() view returns(address)
func (_EngineV5 *EngineV5Session) Voter() (common.Address, error) {
	return _EngineV5.Contract.Voter(&_EngineV5.CallOpts)
}

// Voter is a free data retrieval call binding the contract method 0x46c96aac.
//
// Solidity: function voter() view returns(address)
func (_EngineV5 *EngineV5CallerSession) Voter() (common.Address, error) {
	return _EngineV5.Contract.Voter(&_EngineV5.CallOpts)
}

// VotingPeriodEnded is a free data retrieval call binding the contract method 0x05d1bc26.
//
// Solidity: function votingPeriodEnded(bytes32 taskid_) view returns(bool)
func (_EngineV5 *EngineV5Caller) VotingPeriodEnded(opts *bind.CallOpts, taskid_ [32]byte) (bool, error) {
	var out []interface{}
	err := _EngineV5.contract.Call(opts, &out, "votingPeriodEnded", taskid_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VotingPeriodEnded is a free data retrieval call binding the contract method 0x05d1bc26.
//
// Solidity: function votingPeriodEnded(bytes32 taskid_) view returns(bool)
func (_EngineV5 *EngineV5Session) VotingPeriodEnded(taskid_ [32]byte) (bool, error) {
	return _EngineV5.Contract.VotingPeriodEnded(&_EngineV5.CallOpts, taskid_)
}

// VotingPeriodEnded is a free data retrieval call binding the contract method 0x05d1bc26.
//
// Solidity: function votingPeriodEnded(bytes32 taskid_) view returns(bool)
func (_EngineV5 *EngineV5CallerSession) VotingPeriodEnded(taskid_ [32]byte) (bool, error) {
	return _EngineV5.Contract.VotingPeriodEnded(&_EngineV5.CallOpts, taskid_)
}

// BulkSubmitSolution is a paid mutator transaction binding the contract method 0x65d445fb.
//
// Solidity: function bulkSubmitSolution(bytes32[] taskids_, bytes[] cids_) returns()
func (_EngineV5 *EngineV5Transactor) BulkSubmitSolution(opts *bind.TransactOpts, taskids_ [][32]byte, cids_ [][]byte) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "bulkSubmitSolution", taskids_, cids_)
}

// BulkSubmitSolution is a paid mutator transaction binding the contract method 0x65d445fb.
//
// Solidity: function bulkSubmitSolution(bytes32[] taskids_, bytes[] cids_) returns()
func (_EngineV5 *EngineV5Session) BulkSubmitSolution(taskids_ [][32]byte, cids_ [][]byte) (*types.Transaction, error) {
	return _EngineV5.Contract.BulkSubmitSolution(&_EngineV5.TransactOpts, taskids_, cids_)
}

// BulkSubmitSolution is a paid mutator transaction binding the contract method 0x65d445fb.
//
// Solidity: function bulkSubmitSolution(bytes32[] taskids_, bytes[] cids_) returns()
func (_EngineV5 *EngineV5TransactorSession) BulkSubmitSolution(taskids_ [][32]byte, cids_ [][]byte) (*types.Transaction, error) {
	return _EngineV5.Contract.BulkSubmitSolution(&_EngineV5.TransactOpts, taskids_, cids_)
}

// BulkSubmitTask is a paid mutator transaction binding the contract method 0x08afe0eb.
//
// Solidity: function bulkSubmitTask(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_, uint256 n_) returns()
func (_EngineV5 *EngineV5Transactor) BulkSubmitTask(opts *bind.TransactOpts, version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte, n_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "bulkSubmitTask", version_, owner_, model_, fee_, input_, n_)
}

// BulkSubmitTask is a paid mutator transaction binding the contract method 0x08afe0eb.
//
// Solidity: function bulkSubmitTask(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_, uint256 n_) returns()
func (_EngineV5 *EngineV5Session) BulkSubmitTask(version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte, n_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.Contract.BulkSubmitTask(&_EngineV5.TransactOpts, version_, owner_, model_, fee_, input_, n_)
}

// BulkSubmitTask is a paid mutator transaction binding the contract method 0x08afe0eb.
//
// Solidity: function bulkSubmitTask(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_, uint256 n_) returns()
func (_EngineV5 *EngineV5TransactorSession) BulkSubmitTask(version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte, n_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.Contract.BulkSubmitTask(&_EngineV5.TransactOpts, version_, owner_, model_, fee_, input_, n_)
}

// CancelValidatorWithdraw is a paid mutator transaction binding the contract method 0xcbd2422d.
//
// Solidity: function cancelValidatorWithdraw(uint256 count_) returns()
func (_EngineV5 *EngineV5Transactor) CancelValidatorWithdraw(opts *bind.TransactOpts, count_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "cancelValidatorWithdraw", count_)
}

// CancelValidatorWithdraw is a paid mutator transaction binding the contract method 0xcbd2422d.
//
// Solidity: function cancelValidatorWithdraw(uint256 count_) returns()
func (_EngineV5 *EngineV5Session) CancelValidatorWithdraw(count_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.Contract.CancelValidatorWithdraw(&_EngineV5.TransactOpts, count_)
}

// CancelValidatorWithdraw is a paid mutator transaction binding the contract method 0xcbd2422d.
//
// Solidity: function cancelValidatorWithdraw(uint256 count_) returns()
func (_EngineV5 *EngineV5TransactorSession) CancelValidatorWithdraw(count_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.Contract.CancelValidatorWithdraw(&_EngineV5.TransactOpts, count_)
}

// ClaimSolution is a paid mutator transaction binding the contract method 0x77286d17.
//
// Solidity: function claimSolution(bytes32 taskid_) returns()
func (_EngineV5 *EngineV5Transactor) ClaimSolution(opts *bind.TransactOpts, taskid_ [32]byte) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "claimSolution", taskid_)
}

// ClaimSolution is a paid mutator transaction binding the contract method 0x77286d17.
//
// Solidity: function claimSolution(bytes32 taskid_) returns()
func (_EngineV5 *EngineV5Session) ClaimSolution(taskid_ [32]byte) (*types.Transaction, error) {
	return _EngineV5.Contract.ClaimSolution(&_EngineV5.TransactOpts, taskid_)
}

// ClaimSolution is a paid mutator transaction binding the contract method 0x77286d17.
//
// Solidity: function claimSolution(bytes32 taskid_) returns()
func (_EngineV5 *EngineV5TransactorSession) ClaimSolution(taskid_ [32]byte) (*types.Transaction, error) {
	return _EngineV5.Contract.ClaimSolution(&_EngineV5.TransactOpts, taskid_)
}

// ContestationVoteFinish is a paid mutator transaction binding the contract method 0x8b4d7b35.
//
// Solidity: function contestationVoteFinish(bytes32 taskid_, uint32 amnt_) returns()
func (_EngineV5 *EngineV5Transactor) ContestationVoteFinish(opts *bind.TransactOpts, taskid_ [32]byte, amnt_ uint32) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "contestationVoteFinish", taskid_, amnt_)
}

// ContestationVoteFinish is a paid mutator transaction binding the contract method 0x8b4d7b35.
//
// Solidity: function contestationVoteFinish(bytes32 taskid_, uint32 amnt_) returns()
func (_EngineV5 *EngineV5Session) ContestationVoteFinish(taskid_ [32]byte, amnt_ uint32) (*types.Transaction, error) {
	return _EngineV5.Contract.ContestationVoteFinish(&_EngineV5.TransactOpts, taskid_, amnt_)
}

// ContestationVoteFinish is a paid mutator transaction binding the contract method 0x8b4d7b35.
//
// Solidity: function contestationVoteFinish(bytes32 taskid_, uint32 amnt_) returns()
func (_EngineV5 *EngineV5TransactorSession) ContestationVoteFinish(taskid_ [32]byte, amnt_ uint32) (*types.Transaction, error) {
	return _EngineV5.Contract.ContestationVoteFinish(&_EngineV5.TransactOpts, taskid_, amnt_)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_EngineV5 *EngineV5Transactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_EngineV5 *EngineV5Session) Initialize() (*types.Transaction, error) {
	return _EngineV5.Contract.Initialize(&_EngineV5.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_EngineV5 *EngineV5TransactorSession) Initialize() (*types.Transaction, error) {
	return _EngineV5.Contract.Initialize(&_EngineV5.TransactOpts)
}

// InitiateValidatorWithdraw is a paid mutator transaction binding the contract method 0x0a985737.
//
// Solidity: function initiateValidatorWithdraw(uint256 amount_) returns(uint256)
func (_EngineV5 *EngineV5Transactor) InitiateValidatorWithdraw(opts *bind.TransactOpts, amount_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "initiateValidatorWithdraw", amount_)
}

// InitiateValidatorWithdraw is a paid mutator transaction binding the contract method 0x0a985737.
//
// Solidity: function initiateValidatorWithdraw(uint256 amount_) returns(uint256)
func (_EngineV5 *EngineV5Session) InitiateValidatorWithdraw(amount_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.Contract.InitiateValidatorWithdraw(&_EngineV5.TransactOpts, amount_)
}

// InitiateValidatorWithdraw is a paid mutator transaction binding the contract method 0x0a985737.
//
// Solidity: function initiateValidatorWithdraw(uint256 amount_) returns(uint256)
func (_EngineV5 *EngineV5TransactorSession) InitiateValidatorWithdraw(amount_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.Contract.InitiateValidatorWithdraw(&_EngineV5.TransactOpts, amount_)
}

// RegisterModel is a paid mutator transaction binding the contract method 0x4ff03efa.
//
// Solidity: function registerModel(address addr_, uint256 fee_, bytes template_) returns(bytes32)
func (_EngineV5 *EngineV5Transactor) RegisterModel(opts *bind.TransactOpts, addr_ common.Address, fee_ *big.Int, template_ []byte) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "registerModel", addr_, fee_, template_)
}

// RegisterModel is a paid mutator transaction binding the contract method 0x4ff03efa.
//
// Solidity: function registerModel(address addr_, uint256 fee_, bytes template_) returns(bytes32)
func (_EngineV5 *EngineV5Session) RegisterModel(addr_ common.Address, fee_ *big.Int, template_ []byte) (*types.Transaction, error) {
	return _EngineV5.Contract.RegisterModel(&_EngineV5.TransactOpts, addr_, fee_, template_)
}

// RegisterModel is a paid mutator transaction binding the contract method 0x4ff03efa.
//
// Solidity: function registerModel(address addr_, uint256 fee_, bytes template_) returns(bytes32)
func (_EngineV5 *EngineV5TransactorSession) RegisterModel(addr_ common.Address, fee_ *big.Int, template_ []byte) (*types.Transaction, error) {
	return _EngineV5.Contract.RegisterModel(&_EngineV5.TransactOpts, addr_, fee_, template_)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_EngineV5 *EngineV5Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_EngineV5 *EngineV5Session) RenounceOwnership() (*types.Transaction, error) {
	return _EngineV5.Contract.RenounceOwnership(&_EngineV5.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_EngineV5 *EngineV5TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _EngineV5.Contract.RenounceOwnership(&_EngineV5.TransactOpts)
}

// SetPaused is a paid mutator transaction binding the contract method 0x16c38b3c.
//
// Solidity: function setPaused(bool paused_) returns()
func (_EngineV5 *EngineV5Transactor) SetPaused(opts *bind.TransactOpts, paused_ bool) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "setPaused", paused_)
}

// SetPaused is a paid mutator transaction binding the contract method 0x16c38b3c.
//
// Solidity: function setPaused(bool paused_) returns()
func (_EngineV5 *EngineV5Session) SetPaused(paused_ bool) (*types.Transaction, error) {
	return _EngineV5.Contract.SetPaused(&_EngineV5.TransactOpts, paused_)
}

// SetPaused is a paid mutator transaction binding the contract method 0x16c38b3c.
//
// Solidity: function setPaused(bool paused_) returns()
func (_EngineV5 *EngineV5TransactorSession) SetPaused(paused_ bool) (*types.Transaction, error) {
	return _EngineV5.Contract.SetPaused(&_EngineV5.TransactOpts, paused_)
}

// SetSolutionMineableRate is a paid mutator transaction binding the contract method 0x93f1f8ac.
//
// Solidity: function setSolutionMineableRate(bytes32 model_, uint256 rate_) returns()
func (_EngineV5 *EngineV5Transactor) SetSolutionMineableRate(opts *bind.TransactOpts, model_ [32]byte, rate_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "setSolutionMineableRate", model_, rate_)
}

// SetSolutionMineableRate is a paid mutator transaction binding the contract method 0x93f1f8ac.
//
// Solidity: function setSolutionMineableRate(bytes32 model_, uint256 rate_) returns()
func (_EngineV5 *EngineV5Session) SetSolutionMineableRate(model_ [32]byte, rate_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.Contract.SetSolutionMineableRate(&_EngineV5.TransactOpts, model_, rate_)
}

// SetSolutionMineableRate is a paid mutator transaction binding the contract method 0x93f1f8ac.
//
// Solidity: function setSolutionMineableRate(bytes32 model_, uint256 rate_) returns()
func (_EngineV5 *EngineV5TransactorSession) SetSolutionMineableRate(model_ [32]byte, rate_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.Contract.SetSolutionMineableRate(&_EngineV5.TransactOpts, model_, rate_)
}

// SetStartBlockTime is a paid mutator transaction binding the contract method 0xa8f837f3.
//
// Solidity: function setStartBlockTime(uint64 startBlockTime_) returns()
func (_EngineV5 *EngineV5Transactor) SetStartBlockTime(opts *bind.TransactOpts, startBlockTime_ uint64) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "setStartBlockTime", startBlockTime_)
}

// SetStartBlockTime is a paid mutator transaction binding the contract method 0xa8f837f3.
//
// Solidity: function setStartBlockTime(uint64 startBlockTime_) returns()
func (_EngineV5 *EngineV5Session) SetStartBlockTime(startBlockTime_ uint64) (*types.Transaction, error) {
	return _EngineV5.Contract.SetStartBlockTime(&_EngineV5.TransactOpts, startBlockTime_)
}

// SetStartBlockTime is a paid mutator transaction binding the contract method 0xa8f837f3.
//
// Solidity: function setStartBlockTime(uint64 startBlockTime_) returns()
func (_EngineV5 *EngineV5TransactorSession) SetStartBlockTime(startBlockTime_ uint64) (*types.Transaction, error) {
	return _EngineV5.Contract.SetStartBlockTime(&_EngineV5.TransactOpts, startBlockTime_)
}

// SetVeStaking is a paid mutator transaction binding the contract method 0x2943a490.
//
// Solidity: function setVeStaking(address veStaking_) returns()
func (_EngineV5 *EngineV5Transactor) SetVeStaking(opts *bind.TransactOpts, veStaking_ common.Address) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "setVeStaking", veStaking_)
}

// SetVeStaking is a paid mutator transaction binding the contract method 0x2943a490.
//
// Solidity: function setVeStaking(address veStaking_) returns()
func (_EngineV5 *EngineV5Session) SetVeStaking(veStaking_ common.Address) (*types.Transaction, error) {
	return _EngineV5.Contract.SetVeStaking(&_EngineV5.TransactOpts, veStaking_)
}

// SetVeStaking is a paid mutator transaction binding the contract method 0x2943a490.
//
// Solidity: function setVeStaking(address veStaking_) returns()
func (_EngineV5 *EngineV5TransactorSession) SetVeStaking(veStaking_ common.Address) (*types.Transaction, error) {
	return _EngineV5.Contract.SetVeStaking(&_EngineV5.TransactOpts, veStaking_)
}

// SetVersion is a paid mutator transaction binding the contract method 0x408def1e.
//
// Solidity: function setVersion(uint256 version_) returns()
func (_EngineV5 *EngineV5Transactor) SetVersion(opts *bind.TransactOpts, version_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "setVersion", version_)
}

// SetVersion is a paid mutator transaction binding the contract method 0x408def1e.
//
// Solidity: function setVersion(uint256 version_) returns()
func (_EngineV5 *EngineV5Session) SetVersion(version_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.Contract.SetVersion(&_EngineV5.TransactOpts, version_)
}

// SetVersion is a paid mutator transaction binding the contract method 0x408def1e.
//
// Solidity: function setVersion(uint256 version_) returns()
func (_EngineV5 *EngineV5TransactorSession) SetVersion(version_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.Contract.SetVersion(&_EngineV5.TransactOpts, version_)
}

// SetVoter is a paid mutator transaction binding the contract method 0x4bc2a657.
//
// Solidity: function setVoter(address voter_) returns()
func (_EngineV5 *EngineV5Transactor) SetVoter(opts *bind.TransactOpts, voter_ common.Address) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "setVoter", voter_)
}

// SetVoter is a paid mutator transaction binding the contract method 0x4bc2a657.
//
// Solidity: function setVoter(address voter_) returns()
func (_EngineV5 *EngineV5Session) SetVoter(voter_ common.Address) (*types.Transaction, error) {
	return _EngineV5.Contract.SetVoter(&_EngineV5.TransactOpts, voter_)
}

// SetVoter is a paid mutator transaction binding the contract method 0x4bc2a657.
//
// Solidity: function setVoter(address voter_) returns()
func (_EngineV5 *EngineV5TransactorSession) SetVoter(voter_ common.Address) (*types.Transaction, error) {
	return _EngineV5.Contract.SetVoter(&_EngineV5.TransactOpts, voter_)
}

// SignalCommitment is a paid mutator transaction binding the contract method 0x506ea7de.
//
// Solidity: function signalCommitment(bytes32 commitment_) returns()
func (_EngineV5 *EngineV5Transactor) SignalCommitment(opts *bind.TransactOpts, commitment_ [32]byte) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "signalCommitment", commitment_)
}

// SignalCommitment is a paid mutator transaction binding the contract method 0x506ea7de.
//
// Solidity: function signalCommitment(bytes32 commitment_) returns()
func (_EngineV5 *EngineV5Session) SignalCommitment(commitment_ [32]byte) (*types.Transaction, error) {
	return _EngineV5.Contract.SignalCommitment(&_EngineV5.TransactOpts, commitment_)
}

// SignalCommitment is a paid mutator transaction binding the contract method 0x506ea7de.
//
// Solidity: function signalCommitment(bytes32 commitment_) returns()
func (_EngineV5 *EngineV5TransactorSession) SignalCommitment(commitment_ [32]byte) (*types.Transaction, error) {
	return _EngineV5.Contract.SignalCommitment(&_EngineV5.TransactOpts, commitment_)
}

// SubmitContestation is a paid mutator transaction binding the contract method 0x671f8152.
//
// Solidity: function submitContestation(bytes32 taskid_) returns()
func (_EngineV5 *EngineV5Transactor) SubmitContestation(opts *bind.TransactOpts, taskid_ [32]byte) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "submitContestation", taskid_)
}

// SubmitContestation is a paid mutator transaction binding the contract method 0x671f8152.
//
// Solidity: function submitContestation(bytes32 taskid_) returns()
func (_EngineV5 *EngineV5Session) SubmitContestation(taskid_ [32]byte) (*types.Transaction, error) {
	return _EngineV5.Contract.SubmitContestation(&_EngineV5.TransactOpts, taskid_)
}

// SubmitContestation is a paid mutator transaction binding the contract method 0x671f8152.
//
// Solidity: function submitContestation(bytes32 taskid_) returns()
func (_EngineV5 *EngineV5TransactorSession) SubmitContestation(taskid_ [32]byte) (*types.Transaction, error) {
	return _EngineV5.Contract.SubmitContestation(&_EngineV5.TransactOpts, taskid_)
}

// SubmitSolution is a paid mutator transaction binding the contract method 0x56914caf.
//
// Solidity: function submitSolution(bytes32 taskid_, bytes cid_) returns()
func (_EngineV5 *EngineV5Transactor) SubmitSolution(opts *bind.TransactOpts, taskid_ [32]byte, cid_ []byte) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "submitSolution", taskid_, cid_)
}

// SubmitSolution is a paid mutator transaction binding the contract method 0x56914caf.
//
// Solidity: function submitSolution(bytes32 taskid_, bytes cid_) returns()
func (_EngineV5 *EngineV5Session) SubmitSolution(taskid_ [32]byte, cid_ []byte) (*types.Transaction, error) {
	return _EngineV5.Contract.SubmitSolution(&_EngineV5.TransactOpts, taskid_, cid_)
}

// SubmitSolution is a paid mutator transaction binding the contract method 0x56914caf.
//
// Solidity: function submitSolution(bytes32 taskid_, bytes cid_) returns()
func (_EngineV5 *EngineV5TransactorSession) SubmitSolution(taskid_ [32]byte, cid_ []byte) (*types.Transaction, error) {
	return _EngineV5.Contract.SubmitSolution(&_EngineV5.TransactOpts, taskid_, cid_)
}

// SubmitTask is a paid mutator transaction binding the contract method 0x08745dd1.
//
// Solidity: function submitTask(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_) returns()
func (_EngineV5 *EngineV5Transactor) SubmitTask(opts *bind.TransactOpts, version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "submitTask", version_, owner_, model_, fee_, input_)
}

// SubmitTask is a paid mutator transaction binding the contract method 0x08745dd1.
//
// Solidity: function submitTask(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_) returns()
func (_EngineV5 *EngineV5Session) SubmitTask(version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte) (*types.Transaction, error) {
	return _EngineV5.Contract.SubmitTask(&_EngineV5.TransactOpts, version_, owner_, model_, fee_, input_)
}

// SubmitTask is a paid mutator transaction binding the contract method 0x08745dd1.
//
// Solidity: function submitTask(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_) returns()
func (_EngineV5 *EngineV5TransactorSession) SubmitTask(version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte) (*types.Transaction, error) {
	return _EngineV5.Contract.SubmitTask(&_EngineV5.TransactOpts, version_, owner_, model_, fee_, input_)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to_) returns()
func (_EngineV5 *EngineV5Transactor) TransferOwnership(opts *bind.TransactOpts, to_ common.Address) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "transferOwnership", to_)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to_) returns()
func (_EngineV5 *EngineV5Session) TransferOwnership(to_ common.Address) (*types.Transaction, error) {
	return _EngineV5.Contract.TransferOwnership(&_EngineV5.TransactOpts, to_)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address to_) returns()
func (_EngineV5 *EngineV5TransactorSession) TransferOwnership(to_ common.Address) (*types.Transaction, error) {
	return _EngineV5.Contract.TransferOwnership(&_EngineV5.TransactOpts, to_)
}

// TransferPauser is a paid mutator transaction binding the contract method 0x4421ea21.
//
// Solidity: function transferPauser(address to_) returns()
func (_EngineV5 *EngineV5Transactor) TransferPauser(opts *bind.TransactOpts, to_ common.Address) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "transferPauser", to_)
}

// TransferPauser is a paid mutator transaction binding the contract method 0x4421ea21.
//
// Solidity: function transferPauser(address to_) returns()
func (_EngineV5 *EngineV5Session) TransferPauser(to_ common.Address) (*types.Transaction, error) {
	return _EngineV5.Contract.TransferPauser(&_EngineV5.TransactOpts, to_)
}

// TransferPauser is a paid mutator transaction binding the contract method 0x4421ea21.
//
// Solidity: function transferPauser(address to_) returns()
func (_EngineV5 *EngineV5TransactorSession) TransferPauser(to_ common.Address) (*types.Transaction, error) {
	return _EngineV5.Contract.TransferPauser(&_EngineV5.TransactOpts, to_)
}

// TransferTreasury is a paid mutator transaction binding the contract method 0xd8a6021c.
//
// Solidity: function transferTreasury(address to_) returns()
func (_EngineV5 *EngineV5Transactor) TransferTreasury(opts *bind.TransactOpts, to_ common.Address) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "transferTreasury", to_)
}

// TransferTreasury is a paid mutator transaction binding the contract method 0xd8a6021c.
//
// Solidity: function transferTreasury(address to_) returns()
func (_EngineV5 *EngineV5Session) TransferTreasury(to_ common.Address) (*types.Transaction, error) {
	return _EngineV5.Contract.TransferTreasury(&_EngineV5.TransactOpts, to_)
}

// TransferTreasury is a paid mutator transaction binding the contract method 0xd8a6021c.
//
// Solidity: function transferTreasury(address to_) returns()
func (_EngineV5 *EngineV5TransactorSession) TransferTreasury(to_ common.Address) (*types.Transaction, error) {
	return _EngineV5.Contract.TransferTreasury(&_EngineV5.TransactOpts, to_)
}

// ValidatorDeposit is a paid mutator transaction binding the contract method 0x93a090ec.
//
// Solidity: function validatorDeposit(address validator_, uint256 amount_) returns()
func (_EngineV5 *EngineV5Transactor) ValidatorDeposit(opts *bind.TransactOpts, validator_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "validatorDeposit", validator_, amount_)
}

// ValidatorDeposit is a paid mutator transaction binding the contract method 0x93a090ec.
//
// Solidity: function validatorDeposit(address validator_, uint256 amount_) returns()
func (_EngineV5 *EngineV5Session) ValidatorDeposit(validator_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.Contract.ValidatorDeposit(&_EngineV5.TransactOpts, validator_, amount_)
}

// ValidatorDeposit is a paid mutator transaction binding the contract method 0x93a090ec.
//
// Solidity: function validatorDeposit(address validator_, uint256 amount_) returns()
func (_EngineV5 *EngineV5TransactorSession) ValidatorDeposit(validator_ common.Address, amount_ *big.Int) (*types.Transaction, error) {
	return _EngineV5.Contract.ValidatorDeposit(&_EngineV5.TransactOpts, validator_, amount_)
}

// ValidatorWithdraw is a paid mutator transaction binding the contract method 0x763253bb.
//
// Solidity: function validatorWithdraw(uint256 count_, address to_) returns()
func (_EngineV5 *EngineV5Transactor) ValidatorWithdraw(opts *bind.TransactOpts, count_ *big.Int, to_ common.Address) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "validatorWithdraw", count_, to_)
}

// ValidatorWithdraw is a paid mutator transaction binding the contract method 0x763253bb.
//
// Solidity: function validatorWithdraw(uint256 count_, address to_) returns()
func (_EngineV5 *EngineV5Session) ValidatorWithdraw(count_ *big.Int, to_ common.Address) (*types.Transaction, error) {
	return _EngineV5.Contract.ValidatorWithdraw(&_EngineV5.TransactOpts, count_, to_)
}

// ValidatorWithdraw is a paid mutator transaction binding the contract method 0x763253bb.
//
// Solidity: function validatorWithdraw(uint256 count_, address to_) returns()
func (_EngineV5 *EngineV5TransactorSession) ValidatorWithdraw(count_ *big.Int, to_ common.Address) (*types.Transaction, error) {
	return _EngineV5.Contract.ValidatorWithdraw(&_EngineV5.TransactOpts, count_, to_)
}

// VoteOnContestation is a paid mutator transaction binding the contract method 0x1825c20e.
//
// Solidity: function voteOnContestation(bytes32 taskid_, bool yea_) returns()
func (_EngineV5 *EngineV5Transactor) VoteOnContestation(opts *bind.TransactOpts, taskid_ [32]byte, yea_ bool) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "voteOnContestation", taskid_, yea_)
}

// VoteOnContestation is a paid mutator transaction binding the contract method 0x1825c20e.
//
// Solidity: function voteOnContestation(bytes32 taskid_, bool yea_) returns()
func (_EngineV5 *EngineV5Session) VoteOnContestation(taskid_ [32]byte, yea_ bool) (*types.Transaction, error) {
	return _EngineV5.Contract.VoteOnContestation(&_EngineV5.TransactOpts, taskid_, yea_)
}

// VoteOnContestation is a paid mutator transaction binding the contract method 0x1825c20e.
//
// Solidity: function voteOnContestation(bytes32 taskid_, bool yea_) returns()
func (_EngineV5 *EngineV5TransactorSession) VoteOnContestation(taskid_ [32]byte, yea_ bool) (*types.Transaction, error) {
	return _EngineV5.Contract.VoteOnContestation(&_EngineV5.TransactOpts, taskid_, yea_)
}

// WithdrawAccruedFees is a paid mutator transaction binding the contract method 0xada82c7d.
//
// Solidity: function withdrawAccruedFees() returns()
func (_EngineV5 *EngineV5Transactor) WithdrawAccruedFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EngineV5.contract.Transact(opts, "withdrawAccruedFees")
}

// WithdrawAccruedFees is a paid mutator transaction binding the contract method 0xada82c7d.
//
// Solidity: function withdrawAccruedFees() returns()
func (_EngineV5 *EngineV5Session) WithdrawAccruedFees() (*types.Transaction, error) {
	return _EngineV5.Contract.WithdrawAccruedFees(&_EngineV5.TransactOpts)
}

// WithdrawAccruedFees is a paid mutator transaction binding the contract method 0xada82c7d.
//
// Solidity: function withdrawAccruedFees() returns()
func (_EngineV5 *EngineV5TransactorSession) WithdrawAccruedFees() (*types.Transaction, error) {
	return _EngineV5.Contract.WithdrawAccruedFees(&_EngineV5.TransactOpts)
}

// EngineV5ContestationSubmittedIterator is returned from FilterContestationSubmitted and is used to iterate over the raw logs and unpacked data for ContestationSubmitted events raised by the EngineV5 contract.
type EngineV5ContestationSubmittedIterator struct {
	Event *EngineV5ContestationSubmitted // Event containing the contract specifics and raw log

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
func (it *EngineV5ContestationSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5ContestationSubmitted)
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
		it.Event = new(EngineV5ContestationSubmitted)
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
func (it *EngineV5ContestationSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5ContestationSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5ContestationSubmitted represents a ContestationSubmitted event raised by the EngineV5 contract.
type EngineV5ContestationSubmitted struct {
	Addr common.Address
	Task [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterContestationSubmitted is a free log retrieval operation binding the contract event 0x6958c989e915d3e41a35076e3c480363910055408055ad86ae1ee13d41c40640.
//
// Solidity: event ContestationSubmitted(address indexed addr, bytes32 indexed task)
func (_EngineV5 *EngineV5Filterer) FilterContestationSubmitted(opts *bind.FilterOpts, addr []common.Address, task [][32]byte) (*EngineV5ContestationSubmittedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var taskRule []interface{}
	for _, taskItem := range task {
		taskRule = append(taskRule, taskItem)
	}

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "ContestationSubmitted", addrRule, taskRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5ContestationSubmittedIterator{contract: _EngineV5.contract, event: "ContestationSubmitted", logs: logs, sub: sub}, nil
}

// WatchContestationSubmitted is a free log subscription operation binding the contract event 0x6958c989e915d3e41a35076e3c480363910055408055ad86ae1ee13d41c40640.
//
// Solidity: event ContestationSubmitted(address indexed addr, bytes32 indexed task)
func (_EngineV5 *EngineV5Filterer) WatchContestationSubmitted(opts *bind.WatchOpts, sink chan<- *EngineV5ContestationSubmitted, addr []common.Address, task [][32]byte) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var taskRule []interface{}
	for _, taskItem := range task {
		taskRule = append(taskRule, taskItem)
	}

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "ContestationSubmitted", addrRule, taskRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5ContestationSubmitted)
				if err := _EngineV5.contract.UnpackLog(event, "ContestationSubmitted", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseContestationSubmitted(log types.Log) (*EngineV5ContestationSubmitted, error) {
	event := new(EngineV5ContestationSubmitted)
	if err := _EngineV5.contract.UnpackLog(event, "ContestationSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5ContestationVoteIterator is returned from FilterContestationVote and is used to iterate over the raw logs and unpacked data for ContestationVote events raised by the EngineV5 contract.
type EngineV5ContestationVoteIterator struct {
	Event *EngineV5ContestationVote // Event containing the contract specifics and raw log

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
func (it *EngineV5ContestationVoteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5ContestationVote)
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
		it.Event = new(EngineV5ContestationVote)
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
func (it *EngineV5ContestationVoteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5ContestationVoteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5ContestationVote represents a ContestationVote event raised by the EngineV5 contract.
type EngineV5ContestationVote struct {
	Addr common.Address
	Task [32]byte
	Yea  bool
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterContestationVote is a free log retrieval operation binding the contract event 0x1aa9e4be46e24e1f2e7eeb1613c01629213cd42965d2716e18531b63e552e411.
//
// Solidity: event ContestationVote(address indexed addr, bytes32 indexed task, bool yea)
func (_EngineV5 *EngineV5Filterer) FilterContestationVote(opts *bind.FilterOpts, addr []common.Address, task [][32]byte) (*EngineV5ContestationVoteIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var taskRule []interface{}
	for _, taskItem := range task {
		taskRule = append(taskRule, taskItem)
	}

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "ContestationVote", addrRule, taskRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5ContestationVoteIterator{contract: _EngineV5.contract, event: "ContestationVote", logs: logs, sub: sub}, nil
}

// WatchContestationVote is a free log subscription operation binding the contract event 0x1aa9e4be46e24e1f2e7eeb1613c01629213cd42965d2716e18531b63e552e411.
//
// Solidity: event ContestationVote(address indexed addr, bytes32 indexed task, bool yea)
func (_EngineV5 *EngineV5Filterer) WatchContestationVote(opts *bind.WatchOpts, sink chan<- *EngineV5ContestationVote, addr []common.Address, task [][32]byte) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var taskRule []interface{}
	for _, taskItem := range task {
		taskRule = append(taskRule, taskItem)
	}

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "ContestationVote", addrRule, taskRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5ContestationVote)
				if err := _EngineV5.contract.UnpackLog(event, "ContestationVote", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseContestationVote(log types.Log) (*EngineV5ContestationVote, error) {
	event := new(EngineV5ContestationVote)
	if err := _EngineV5.contract.UnpackLog(event, "ContestationVote", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5ContestationVoteFinishIterator is returned from FilterContestationVoteFinish and is used to iterate over the raw logs and unpacked data for ContestationVoteFinish events raised by the EngineV5 contract.
type EngineV5ContestationVoteFinishIterator struct {
	Event *EngineV5ContestationVoteFinish // Event containing the contract specifics and raw log

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
func (it *EngineV5ContestationVoteFinishIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5ContestationVoteFinish)
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
		it.Event = new(EngineV5ContestationVoteFinish)
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
func (it *EngineV5ContestationVoteFinishIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5ContestationVoteFinishIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5ContestationVoteFinish represents a ContestationVoteFinish event raised by the EngineV5 contract.
type EngineV5ContestationVoteFinish struct {
	Id       [32]byte
	StartIdx uint32
	EndIdx   uint32
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterContestationVoteFinish is a free log retrieval operation binding the contract event 0x71d8c71303e35a39162e33a402c9897bf9848388537bac7d5e1b0d202eca4e66.
//
// Solidity: event ContestationVoteFinish(bytes32 indexed id, uint32 indexed start_idx, uint32 end_idx)
func (_EngineV5 *EngineV5Filterer) FilterContestationVoteFinish(opts *bind.FilterOpts, id [][32]byte, start_idx []uint32) (*EngineV5ContestationVoteFinishIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var start_idxRule []interface{}
	for _, start_idxItem := range start_idx {
		start_idxRule = append(start_idxRule, start_idxItem)
	}

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "ContestationVoteFinish", idRule, start_idxRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5ContestationVoteFinishIterator{contract: _EngineV5.contract, event: "ContestationVoteFinish", logs: logs, sub: sub}, nil
}

// WatchContestationVoteFinish is a free log subscription operation binding the contract event 0x71d8c71303e35a39162e33a402c9897bf9848388537bac7d5e1b0d202eca4e66.
//
// Solidity: event ContestationVoteFinish(bytes32 indexed id, uint32 indexed start_idx, uint32 end_idx)
func (_EngineV5 *EngineV5Filterer) WatchContestationVoteFinish(opts *bind.WatchOpts, sink chan<- *EngineV5ContestationVoteFinish, id [][32]byte, start_idx []uint32) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var start_idxRule []interface{}
	for _, start_idxItem := range start_idx {
		start_idxRule = append(start_idxRule, start_idxItem)
	}

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "ContestationVoteFinish", idRule, start_idxRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5ContestationVoteFinish)
				if err := _EngineV5.contract.UnpackLog(event, "ContestationVoteFinish", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseContestationVoteFinish(log types.Log) (*EngineV5ContestationVoteFinish, error) {
	event := new(EngineV5ContestationVoteFinish)
	if err := _EngineV5.contract.UnpackLog(event, "ContestationVoteFinish", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5InitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the EngineV5 contract.
type EngineV5InitializedIterator struct {
	Event *EngineV5Initialized // Event containing the contract specifics and raw log

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
func (it *EngineV5InitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5Initialized)
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
		it.Event = new(EngineV5Initialized)
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
func (it *EngineV5InitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5InitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5Initialized represents a Initialized event raised by the EngineV5 contract.
type EngineV5Initialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_EngineV5 *EngineV5Filterer) FilterInitialized(opts *bind.FilterOpts) (*EngineV5InitializedIterator, error) {

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &EngineV5InitializedIterator{contract: _EngineV5.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_EngineV5 *EngineV5Filterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *EngineV5Initialized) (event.Subscription, error) {

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5Initialized)
				if err := _EngineV5.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseInitialized(log types.Log) (*EngineV5Initialized, error) {
	event := new(EngineV5Initialized)
	if err := _EngineV5.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5ModelRegisteredIterator is returned from FilterModelRegistered and is used to iterate over the raw logs and unpacked data for ModelRegistered events raised by the EngineV5 contract.
type EngineV5ModelRegisteredIterator struct {
	Event *EngineV5ModelRegistered // Event containing the contract specifics and raw log

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
func (it *EngineV5ModelRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5ModelRegistered)
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
		it.Event = new(EngineV5ModelRegistered)
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
func (it *EngineV5ModelRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5ModelRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5ModelRegistered represents a ModelRegistered event raised by the EngineV5 contract.
type EngineV5ModelRegistered struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterModelRegistered is a free log retrieval operation binding the contract event 0xa4b0af38d049ba81703a0d0e46cc2ff39681210302134046237111a8fb7dee72.
//
// Solidity: event ModelRegistered(bytes32 indexed id)
func (_EngineV5 *EngineV5Filterer) FilterModelRegistered(opts *bind.FilterOpts, id [][32]byte) (*EngineV5ModelRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "ModelRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5ModelRegisteredIterator{contract: _EngineV5.contract, event: "ModelRegistered", logs: logs, sub: sub}, nil
}

// WatchModelRegistered is a free log subscription operation binding the contract event 0xa4b0af38d049ba81703a0d0e46cc2ff39681210302134046237111a8fb7dee72.
//
// Solidity: event ModelRegistered(bytes32 indexed id)
func (_EngineV5 *EngineV5Filterer) WatchModelRegistered(opts *bind.WatchOpts, sink chan<- *EngineV5ModelRegistered, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "ModelRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5ModelRegistered)
				if err := _EngineV5.contract.UnpackLog(event, "ModelRegistered", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseModelRegistered(log types.Log) (*EngineV5ModelRegistered, error) {
	event := new(EngineV5ModelRegistered)
	if err := _EngineV5.contract.UnpackLog(event, "ModelRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the EngineV5 contract.
type EngineV5OwnershipTransferredIterator struct {
	Event *EngineV5OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *EngineV5OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5OwnershipTransferred)
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
		it.Event = new(EngineV5OwnershipTransferred)
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
func (it *EngineV5OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5OwnershipTransferred represents a OwnershipTransferred event raised by the EngineV5 contract.
type EngineV5OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_EngineV5 *EngineV5Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*EngineV5OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5OwnershipTransferredIterator{contract: _EngineV5.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_EngineV5 *EngineV5Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *EngineV5OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5OwnershipTransferred)
				if err := _EngineV5.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseOwnershipTransferred(log types.Log) (*EngineV5OwnershipTransferred, error) {
	event := new(EngineV5OwnershipTransferred)
	if err := _EngineV5.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5PausedChangedIterator is returned from FilterPausedChanged and is used to iterate over the raw logs and unpacked data for PausedChanged events raised by the EngineV5 contract.
type EngineV5PausedChangedIterator struct {
	Event *EngineV5PausedChanged // Event containing the contract specifics and raw log

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
func (it *EngineV5PausedChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5PausedChanged)
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
		it.Event = new(EngineV5PausedChanged)
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
func (it *EngineV5PausedChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5PausedChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5PausedChanged represents a PausedChanged event raised by the EngineV5 contract.
type EngineV5PausedChanged struct {
	Paused bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPausedChanged is a free log retrieval operation binding the contract event 0xd83d5281277e107f080e362699d46082adb74e7dc6a9bccbc87d8ae9533add44.
//
// Solidity: event PausedChanged(bool indexed paused)
func (_EngineV5 *EngineV5Filterer) FilterPausedChanged(opts *bind.FilterOpts, paused []bool) (*EngineV5PausedChangedIterator, error) {

	var pausedRule []interface{}
	for _, pausedItem := range paused {
		pausedRule = append(pausedRule, pausedItem)
	}

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "PausedChanged", pausedRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5PausedChangedIterator{contract: _EngineV5.contract, event: "PausedChanged", logs: logs, sub: sub}, nil
}

// WatchPausedChanged is a free log subscription operation binding the contract event 0xd83d5281277e107f080e362699d46082adb74e7dc6a9bccbc87d8ae9533add44.
//
// Solidity: event PausedChanged(bool indexed paused)
func (_EngineV5 *EngineV5Filterer) WatchPausedChanged(opts *bind.WatchOpts, sink chan<- *EngineV5PausedChanged, paused []bool) (event.Subscription, error) {

	var pausedRule []interface{}
	for _, pausedItem := range paused {
		pausedRule = append(pausedRule, pausedItem)
	}

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "PausedChanged", pausedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5PausedChanged)
				if err := _EngineV5.contract.UnpackLog(event, "PausedChanged", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParsePausedChanged(log types.Log) (*EngineV5PausedChanged, error) {
	event := new(EngineV5PausedChanged)
	if err := _EngineV5.contract.UnpackLog(event, "PausedChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5PauserTransferredIterator is returned from FilterPauserTransferred and is used to iterate over the raw logs and unpacked data for PauserTransferred events raised by the EngineV5 contract.
type EngineV5PauserTransferredIterator struct {
	Event *EngineV5PauserTransferred // Event containing the contract specifics and raw log

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
func (it *EngineV5PauserTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5PauserTransferred)
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
		it.Event = new(EngineV5PauserTransferred)
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
func (it *EngineV5PauserTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5PauserTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5PauserTransferred represents a PauserTransferred event raised by the EngineV5 contract.
type EngineV5PauserTransferred struct {
	To  common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPauserTransferred is a free log retrieval operation binding the contract event 0x5a85b4270fc1e75035e6cd505418ce65e0bcc36cc7eb9ce9e6f8c6181d4cb577.
//
// Solidity: event PauserTransferred(address indexed to)
func (_EngineV5 *EngineV5Filterer) FilterPauserTransferred(opts *bind.FilterOpts, to []common.Address) (*EngineV5PauserTransferredIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "PauserTransferred", toRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5PauserTransferredIterator{contract: _EngineV5.contract, event: "PauserTransferred", logs: logs, sub: sub}, nil
}

// WatchPauserTransferred is a free log subscription operation binding the contract event 0x5a85b4270fc1e75035e6cd505418ce65e0bcc36cc7eb9ce9e6f8c6181d4cb577.
//
// Solidity: event PauserTransferred(address indexed to)
func (_EngineV5 *EngineV5Filterer) WatchPauserTransferred(opts *bind.WatchOpts, sink chan<- *EngineV5PauserTransferred, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "PauserTransferred", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5PauserTransferred)
				if err := _EngineV5.contract.UnpackLog(event, "PauserTransferred", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParsePauserTransferred(log types.Log) (*EngineV5PauserTransferred, error) {
	event := new(EngineV5PauserTransferred)
	if err := _EngineV5.contract.UnpackLog(event, "PauserTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5SignalCommitmentIterator is returned from FilterSignalCommitment and is used to iterate over the raw logs and unpacked data for SignalCommitment events raised by the EngineV5 contract.
type EngineV5SignalCommitmentIterator struct {
	Event *EngineV5SignalCommitment // Event containing the contract specifics and raw log

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
func (it *EngineV5SignalCommitmentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5SignalCommitment)
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
		it.Event = new(EngineV5SignalCommitment)
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
func (it *EngineV5SignalCommitmentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5SignalCommitmentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5SignalCommitment represents a SignalCommitment event raised by the EngineV5 contract.
type EngineV5SignalCommitment struct {
	Addr       common.Address
	Commitment [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSignalCommitment is a free log retrieval operation binding the contract event 0x09b4c028a2e50fec6f1c6a0163c59e8fbe92b231e5c03ef3adec585e63a14b92.
//
// Solidity: event SignalCommitment(address indexed addr, bytes32 indexed commitment)
func (_EngineV5 *EngineV5Filterer) FilterSignalCommitment(opts *bind.FilterOpts, addr []common.Address, commitment [][32]byte) (*EngineV5SignalCommitmentIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var commitmentRule []interface{}
	for _, commitmentItem := range commitment {
		commitmentRule = append(commitmentRule, commitmentItem)
	}

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "SignalCommitment", addrRule, commitmentRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5SignalCommitmentIterator{contract: _EngineV5.contract, event: "SignalCommitment", logs: logs, sub: sub}, nil
}

// WatchSignalCommitment is a free log subscription operation binding the contract event 0x09b4c028a2e50fec6f1c6a0163c59e8fbe92b231e5c03ef3adec585e63a14b92.
//
// Solidity: event SignalCommitment(address indexed addr, bytes32 indexed commitment)
func (_EngineV5 *EngineV5Filterer) WatchSignalCommitment(opts *bind.WatchOpts, sink chan<- *EngineV5SignalCommitment, addr []common.Address, commitment [][32]byte) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var commitmentRule []interface{}
	for _, commitmentItem := range commitment {
		commitmentRule = append(commitmentRule, commitmentItem)
	}

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "SignalCommitment", addrRule, commitmentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5SignalCommitment)
				if err := _EngineV5.contract.UnpackLog(event, "SignalCommitment", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseSignalCommitment(log types.Log) (*EngineV5SignalCommitment, error) {
	event := new(EngineV5SignalCommitment)
	if err := _EngineV5.contract.UnpackLog(event, "SignalCommitment", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5SolutionClaimedIterator is returned from FilterSolutionClaimed and is used to iterate over the raw logs and unpacked data for SolutionClaimed events raised by the EngineV5 contract.
type EngineV5SolutionClaimedIterator struct {
	Event *EngineV5SolutionClaimed // Event containing the contract specifics and raw log

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
func (it *EngineV5SolutionClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5SolutionClaimed)
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
		it.Event = new(EngineV5SolutionClaimed)
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
func (it *EngineV5SolutionClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5SolutionClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5SolutionClaimed represents a SolutionClaimed event raised by the EngineV5 contract.
type EngineV5SolutionClaimed struct {
	Addr common.Address
	Task [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSolutionClaimed is a free log retrieval operation binding the contract event 0x0b76b4ae356796814d36b46f7c500bbd27b2cce1e6059a6fa2bebfd5a389b190.
//
// Solidity: event SolutionClaimed(address indexed addr, bytes32 indexed task)
func (_EngineV5 *EngineV5Filterer) FilterSolutionClaimed(opts *bind.FilterOpts, addr []common.Address, task [][32]byte) (*EngineV5SolutionClaimedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var taskRule []interface{}
	for _, taskItem := range task {
		taskRule = append(taskRule, taskItem)
	}

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "SolutionClaimed", addrRule, taskRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5SolutionClaimedIterator{contract: _EngineV5.contract, event: "SolutionClaimed", logs: logs, sub: sub}, nil
}

// WatchSolutionClaimed is a free log subscription operation binding the contract event 0x0b76b4ae356796814d36b46f7c500bbd27b2cce1e6059a6fa2bebfd5a389b190.
//
// Solidity: event SolutionClaimed(address indexed addr, bytes32 indexed task)
func (_EngineV5 *EngineV5Filterer) WatchSolutionClaimed(opts *bind.WatchOpts, sink chan<- *EngineV5SolutionClaimed, addr []common.Address, task [][32]byte) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var taskRule []interface{}
	for _, taskItem := range task {
		taskRule = append(taskRule, taskItem)
	}

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "SolutionClaimed", addrRule, taskRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5SolutionClaimed)
				if err := _EngineV5.contract.UnpackLog(event, "SolutionClaimed", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseSolutionClaimed(log types.Log) (*EngineV5SolutionClaimed, error) {
	event := new(EngineV5SolutionClaimed)
	if err := _EngineV5.contract.UnpackLog(event, "SolutionClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5SolutionMineableRateChangeIterator is returned from FilterSolutionMineableRateChange and is used to iterate over the raw logs and unpacked data for SolutionMineableRateChange events raised by the EngineV5 contract.
type EngineV5SolutionMineableRateChangeIterator struct {
	Event *EngineV5SolutionMineableRateChange // Event containing the contract specifics and raw log

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
func (it *EngineV5SolutionMineableRateChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5SolutionMineableRateChange)
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
		it.Event = new(EngineV5SolutionMineableRateChange)
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
func (it *EngineV5SolutionMineableRateChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5SolutionMineableRateChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5SolutionMineableRateChange represents a SolutionMineableRateChange event raised by the EngineV5 contract.
type EngineV5SolutionMineableRateChange struct {
	Id   [32]byte
	Rate *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSolutionMineableRateChange is a free log retrieval operation binding the contract event 0x0321e8a918b4c47bf3677852c070983825f30a47bc8d9416691454fa6a727d63.
//
// Solidity: event SolutionMineableRateChange(bytes32 indexed id, uint256 rate)
func (_EngineV5 *EngineV5Filterer) FilterSolutionMineableRateChange(opts *bind.FilterOpts, id [][32]byte) (*EngineV5SolutionMineableRateChangeIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "SolutionMineableRateChange", idRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5SolutionMineableRateChangeIterator{contract: _EngineV5.contract, event: "SolutionMineableRateChange", logs: logs, sub: sub}, nil
}

// WatchSolutionMineableRateChange is a free log subscription operation binding the contract event 0x0321e8a918b4c47bf3677852c070983825f30a47bc8d9416691454fa6a727d63.
//
// Solidity: event SolutionMineableRateChange(bytes32 indexed id, uint256 rate)
func (_EngineV5 *EngineV5Filterer) WatchSolutionMineableRateChange(opts *bind.WatchOpts, sink chan<- *EngineV5SolutionMineableRateChange, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "SolutionMineableRateChange", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5SolutionMineableRateChange)
				if err := _EngineV5.contract.UnpackLog(event, "SolutionMineableRateChange", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseSolutionMineableRateChange(log types.Log) (*EngineV5SolutionMineableRateChange, error) {
	event := new(EngineV5SolutionMineableRateChange)
	if err := _EngineV5.contract.UnpackLog(event, "SolutionMineableRateChange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5SolutionSubmittedIterator is returned from FilterSolutionSubmitted and is used to iterate over the raw logs and unpacked data for SolutionSubmitted events raised by the EngineV5 contract.
type EngineV5SolutionSubmittedIterator struct {
	Event *EngineV5SolutionSubmitted // Event containing the contract specifics and raw log

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
func (it *EngineV5SolutionSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5SolutionSubmitted)
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
		it.Event = new(EngineV5SolutionSubmitted)
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
func (it *EngineV5SolutionSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5SolutionSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5SolutionSubmitted represents a SolutionSubmitted event raised by the EngineV5 contract.
type EngineV5SolutionSubmitted struct {
	Addr common.Address
	Task [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSolutionSubmitted is a free log retrieval operation binding the contract event 0x957c18b5af8413899ea8a576a4d3fb16839a02c9fccfdce098b6d59ef248525b.
//
// Solidity: event SolutionSubmitted(address indexed addr, bytes32 indexed task)
func (_EngineV5 *EngineV5Filterer) FilterSolutionSubmitted(opts *bind.FilterOpts, addr []common.Address, task [][32]byte) (*EngineV5SolutionSubmittedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var taskRule []interface{}
	for _, taskItem := range task {
		taskRule = append(taskRule, taskItem)
	}

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "SolutionSubmitted", addrRule, taskRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5SolutionSubmittedIterator{contract: _EngineV5.contract, event: "SolutionSubmitted", logs: logs, sub: sub}, nil
}

// WatchSolutionSubmitted is a free log subscription operation binding the contract event 0x957c18b5af8413899ea8a576a4d3fb16839a02c9fccfdce098b6d59ef248525b.
//
// Solidity: event SolutionSubmitted(address indexed addr, bytes32 indexed task)
func (_EngineV5 *EngineV5Filterer) WatchSolutionSubmitted(opts *bind.WatchOpts, sink chan<- *EngineV5SolutionSubmitted, addr []common.Address, task [][32]byte) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var taskRule []interface{}
	for _, taskItem := range task {
		taskRule = append(taskRule, taskItem)
	}

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "SolutionSubmitted", addrRule, taskRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5SolutionSubmitted)
				if err := _EngineV5.contract.UnpackLog(event, "SolutionSubmitted", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseSolutionSubmitted(log types.Log) (*EngineV5SolutionSubmitted, error) {
	event := new(EngineV5SolutionSubmitted)
	if err := _EngineV5.contract.UnpackLog(event, "SolutionSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5StartBlockTimeChangedIterator is returned from FilterStartBlockTimeChanged and is used to iterate over the raw logs and unpacked data for StartBlockTimeChanged events raised by the EngineV5 contract.
type EngineV5StartBlockTimeChangedIterator struct {
	Event *EngineV5StartBlockTimeChanged // Event containing the contract specifics and raw log

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
func (it *EngineV5StartBlockTimeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5StartBlockTimeChanged)
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
		it.Event = new(EngineV5StartBlockTimeChanged)
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
func (it *EngineV5StartBlockTimeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5StartBlockTimeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5StartBlockTimeChanged represents a StartBlockTimeChanged event raised by the EngineV5 contract.
type EngineV5StartBlockTimeChanged struct {
	StartBlockTime uint64
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterStartBlockTimeChanged is a free log retrieval operation binding the contract event 0xa15d6b0930a82638ac4775bfd1b2e9f1e86be67e1bd3a09fa8a77a8f079769d7.
//
// Solidity: event StartBlockTimeChanged(uint64 indexed startBlockTime)
func (_EngineV5 *EngineV5Filterer) FilterStartBlockTimeChanged(opts *bind.FilterOpts, startBlockTime []uint64) (*EngineV5StartBlockTimeChangedIterator, error) {

	var startBlockTimeRule []interface{}
	for _, startBlockTimeItem := range startBlockTime {
		startBlockTimeRule = append(startBlockTimeRule, startBlockTimeItem)
	}

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "StartBlockTimeChanged", startBlockTimeRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5StartBlockTimeChangedIterator{contract: _EngineV5.contract, event: "StartBlockTimeChanged", logs: logs, sub: sub}, nil
}

// WatchStartBlockTimeChanged is a free log subscription operation binding the contract event 0xa15d6b0930a82638ac4775bfd1b2e9f1e86be67e1bd3a09fa8a77a8f079769d7.
//
// Solidity: event StartBlockTimeChanged(uint64 indexed startBlockTime)
func (_EngineV5 *EngineV5Filterer) WatchStartBlockTimeChanged(opts *bind.WatchOpts, sink chan<- *EngineV5StartBlockTimeChanged, startBlockTime []uint64) (event.Subscription, error) {

	var startBlockTimeRule []interface{}
	for _, startBlockTimeItem := range startBlockTime {
		startBlockTimeRule = append(startBlockTimeRule, startBlockTimeItem)
	}

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "StartBlockTimeChanged", startBlockTimeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5StartBlockTimeChanged)
				if err := _EngineV5.contract.UnpackLog(event, "StartBlockTimeChanged", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseStartBlockTimeChanged(log types.Log) (*EngineV5StartBlockTimeChanged, error) {
	event := new(EngineV5StartBlockTimeChanged)
	if err := _EngineV5.contract.UnpackLog(event, "StartBlockTimeChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5TaskSubmittedIterator is returned from FilterTaskSubmitted and is used to iterate over the raw logs and unpacked data for TaskSubmitted events raised by the EngineV5 contract.
type EngineV5TaskSubmittedIterator struct {
	Event *EngineV5TaskSubmitted // Event containing the contract specifics and raw log

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
func (it *EngineV5TaskSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5TaskSubmitted)
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
		it.Event = new(EngineV5TaskSubmitted)
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
func (it *EngineV5TaskSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5TaskSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5TaskSubmitted represents a TaskSubmitted event raised by the EngineV5 contract.
type EngineV5TaskSubmitted struct {
	Id     [32]byte
	Model  [32]byte
	Fee    *big.Int
	Sender common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTaskSubmitted is a free log retrieval operation binding the contract event 0xc3d3e0544c80e3bb83f62659259ae1574f72a91515ab3cae3dd75cf77e1b0aea.
//
// Solidity: event TaskSubmitted(bytes32 indexed id, bytes32 indexed model, uint256 fee, address indexed sender)
func (_EngineV5 *EngineV5Filterer) FilterTaskSubmitted(opts *bind.FilterOpts, id [][32]byte, model [][32]byte, sender []common.Address) (*EngineV5TaskSubmittedIterator, error) {

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

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "TaskSubmitted", idRule, modelRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5TaskSubmittedIterator{contract: _EngineV5.contract, event: "TaskSubmitted", logs: logs, sub: sub}, nil
}

// WatchTaskSubmitted is a free log subscription operation binding the contract event 0xc3d3e0544c80e3bb83f62659259ae1574f72a91515ab3cae3dd75cf77e1b0aea.
//
// Solidity: event TaskSubmitted(bytes32 indexed id, bytes32 indexed model, uint256 fee, address indexed sender)
func (_EngineV5 *EngineV5Filterer) WatchTaskSubmitted(opts *bind.WatchOpts, sink chan<- *EngineV5TaskSubmitted, id [][32]byte, model [][32]byte, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "TaskSubmitted", idRule, modelRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5TaskSubmitted)
				if err := _EngineV5.contract.UnpackLog(event, "TaskSubmitted", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseTaskSubmitted(log types.Log) (*EngineV5TaskSubmitted, error) {
	event := new(EngineV5TaskSubmitted)
	if err := _EngineV5.contract.UnpackLog(event, "TaskSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5TreasuryTransferredIterator is returned from FilterTreasuryTransferred and is used to iterate over the raw logs and unpacked data for TreasuryTransferred events raised by the EngineV5 contract.
type EngineV5TreasuryTransferredIterator struct {
	Event *EngineV5TreasuryTransferred // Event containing the contract specifics and raw log

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
func (it *EngineV5TreasuryTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5TreasuryTransferred)
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
		it.Event = new(EngineV5TreasuryTransferred)
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
func (it *EngineV5TreasuryTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5TreasuryTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5TreasuryTransferred represents a TreasuryTransferred event raised by the EngineV5 contract.
type EngineV5TreasuryTransferred struct {
	To  common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterTreasuryTransferred is a free log retrieval operation binding the contract event 0x6bdb9ceff405d990c9b60be9f719fbb80889d5f064e8fd76efd5bea353e90712.
//
// Solidity: event TreasuryTransferred(address indexed to)
func (_EngineV5 *EngineV5Filterer) FilterTreasuryTransferred(opts *bind.FilterOpts, to []common.Address) (*EngineV5TreasuryTransferredIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "TreasuryTransferred", toRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5TreasuryTransferredIterator{contract: _EngineV5.contract, event: "TreasuryTransferred", logs: logs, sub: sub}, nil
}

// WatchTreasuryTransferred is a free log subscription operation binding the contract event 0x6bdb9ceff405d990c9b60be9f719fbb80889d5f064e8fd76efd5bea353e90712.
//
// Solidity: event TreasuryTransferred(address indexed to)
func (_EngineV5 *EngineV5Filterer) WatchTreasuryTransferred(opts *bind.WatchOpts, sink chan<- *EngineV5TreasuryTransferred, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "TreasuryTransferred", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5TreasuryTransferred)
				if err := _EngineV5.contract.UnpackLog(event, "TreasuryTransferred", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseTreasuryTransferred(log types.Log) (*EngineV5TreasuryTransferred, error) {
	event := new(EngineV5TreasuryTransferred)
	if err := _EngineV5.contract.UnpackLog(event, "TreasuryTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5ValidatorDepositIterator is returned from FilterValidatorDeposit and is used to iterate over the raw logs and unpacked data for ValidatorDeposit events raised by the EngineV5 contract.
type EngineV5ValidatorDepositIterator struct {
	Event *EngineV5ValidatorDeposit // Event containing the contract specifics and raw log

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
func (it *EngineV5ValidatorDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5ValidatorDeposit)
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
		it.Event = new(EngineV5ValidatorDeposit)
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
func (it *EngineV5ValidatorDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5ValidatorDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5ValidatorDeposit represents a ValidatorDeposit event raised by the EngineV5 contract.
type EngineV5ValidatorDeposit struct {
	Addr      common.Address
	Validator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterValidatorDeposit is a free log retrieval operation binding the contract event 0x8d4844488c19a90828439e71d14ebad860806d04f8ef8b25a82179fab2699b89.
//
// Solidity: event ValidatorDeposit(address indexed addr, address indexed validator, uint256 amount)
func (_EngineV5 *EngineV5Filterer) FilterValidatorDeposit(opts *bind.FilterOpts, addr []common.Address, validator []common.Address) (*EngineV5ValidatorDepositIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "ValidatorDeposit", addrRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5ValidatorDepositIterator{contract: _EngineV5.contract, event: "ValidatorDeposit", logs: logs, sub: sub}, nil
}

// WatchValidatorDeposit is a free log subscription operation binding the contract event 0x8d4844488c19a90828439e71d14ebad860806d04f8ef8b25a82179fab2699b89.
//
// Solidity: event ValidatorDeposit(address indexed addr, address indexed validator, uint256 amount)
func (_EngineV5 *EngineV5Filterer) WatchValidatorDeposit(opts *bind.WatchOpts, sink chan<- *EngineV5ValidatorDeposit, addr []common.Address, validator []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "ValidatorDeposit", addrRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5ValidatorDeposit)
				if err := _EngineV5.contract.UnpackLog(event, "ValidatorDeposit", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseValidatorDeposit(log types.Log) (*EngineV5ValidatorDeposit, error) {
	event := new(EngineV5ValidatorDeposit)
	if err := _EngineV5.contract.UnpackLog(event, "ValidatorDeposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5ValidatorWithdrawIterator is returned from FilterValidatorWithdraw and is used to iterate over the raw logs and unpacked data for ValidatorWithdraw events raised by the EngineV5 contract.
type EngineV5ValidatorWithdrawIterator struct {
	Event *EngineV5ValidatorWithdraw // Event containing the contract specifics and raw log

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
func (it *EngineV5ValidatorWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5ValidatorWithdraw)
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
		it.Event = new(EngineV5ValidatorWithdraw)
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
func (it *EngineV5ValidatorWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5ValidatorWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5ValidatorWithdraw represents a ValidatorWithdraw event raised by the EngineV5 contract.
type EngineV5ValidatorWithdraw struct {
	Addr   common.Address
	To     common.Address
	Count  *big.Int
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterValidatorWithdraw is a free log retrieval operation binding the contract event 0x109aeff667601aad33abcb7c5df1617754eefd18253d586a95c198cb479b5bd4.
//
// Solidity: event ValidatorWithdraw(address indexed addr, address indexed to, uint256 indexed count, uint256 amount)
func (_EngineV5 *EngineV5Filterer) FilterValidatorWithdraw(opts *bind.FilterOpts, addr []common.Address, to []common.Address, count []*big.Int) (*EngineV5ValidatorWithdrawIterator, error) {

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

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "ValidatorWithdraw", addrRule, toRule, countRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5ValidatorWithdrawIterator{contract: _EngineV5.contract, event: "ValidatorWithdraw", logs: logs, sub: sub}, nil
}

// WatchValidatorWithdraw is a free log subscription operation binding the contract event 0x109aeff667601aad33abcb7c5df1617754eefd18253d586a95c198cb479b5bd4.
//
// Solidity: event ValidatorWithdraw(address indexed addr, address indexed to, uint256 indexed count, uint256 amount)
func (_EngineV5 *EngineV5Filterer) WatchValidatorWithdraw(opts *bind.WatchOpts, sink chan<- *EngineV5ValidatorWithdraw, addr []common.Address, to []common.Address, count []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "ValidatorWithdraw", addrRule, toRule, countRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5ValidatorWithdraw)
				if err := _EngineV5.contract.UnpackLog(event, "ValidatorWithdraw", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseValidatorWithdraw(log types.Log) (*EngineV5ValidatorWithdraw, error) {
	event := new(EngineV5ValidatorWithdraw)
	if err := _EngineV5.contract.UnpackLog(event, "ValidatorWithdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5ValidatorWithdrawCancelledIterator is returned from FilterValidatorWithdrawCancelled and is used to iterate over the raw logs and unpacked data for ValidatorWithdrawCancelled events raised by the EngineV5 contract.
type EngineV5ValidatorWithdrawCancelledIterator struct {
	Event *EngineV5ValidatorWithdrawCancelled // Event containing the contract specifics and raw log

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
func (it *EngineV5ValidatorWithdrawCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5ValidatorWithdrawCancelled)
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
		it.Event = new(EngineV5ValidatorWithdrawCancelled)
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
func (it *EngineV5ValidatorWithdrawCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5ValidatorWithdrawCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5ValidatorWithdrawCancelled represents a ValidatorWithdrawCancelled event raised by the EngineV5 contract.
type EngineV5ValidatorWithdrawCancelled struct {
	Addr  common.Address
	Count *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterValidatorWithdrawCancelled is a free log retrieval operation binding the contract event 0xf9d5cfcdfe9803069225971ba315f4302add9b477c3441cc363bc38fdb065b90.
//
// Solidity: event ValidatorWithdrawCancelled(address indexed addr, uint256 indexed count)
func (_EngineV5 *EngineV5Filterer) FilterValidatorWithdrawCancelled(opts *bind.FilterOpts, addr []common.Address, count []*big.Int) (*EngineV5ValidatorWithdrawCancelledIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var countRule []interface{}
	for _, countItem := range count {
		countRule = append(countRule, countItem)
	}

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "ValidatorWithdrawCancelled", addrRule, countRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5ValidatorWithdrawCancelledIterator{contract: _EngineV5.contract, event: "ValidatorWithdrawCancelled", logs: logs, sub: sub}, nil
}

// WatchValidatorWithdrawCancelled is a free log subscription operation binding the contract event 0xf9d5cfcdfe9803069225971ba315f4302add9b477c3441cc363bc38fdb065b90.
//
// Solidity: event ValidatorWithdrawCancelled(address indexed addr, uint256 indexed count)
func (_EngineV5 *EngineV5Filterer) WatchValidatorWithdrawCancelled(opts *bind.WatchOpts, sink chan<- *EngineV5ValidatorWithdrawCancelled, addr []common.Address, count []*big.Int) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var countRule []interface{}
	for _, countItem := range count {
		countRule = append(countRule, countItem)
	}

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "ValidatorWithdrawCancelled", addrRule, countRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5ValidatorWithdrawCancelled)
				if err := _EngineV5.contract.UnpackLog(event, "ValidatorWithdrawCancelled", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseValidatorWithdrawCancelled(log types.Log) (*EngineV5ValidatorWithdrawCancelled, error) {
	event := new(EngineV5ValidatorWithdrawCancelled)
	if err := _EngineV5.contract.UnpackLog(event, "ValidatorWithdrawCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5ValidatorWithdrawInitiatedIterator is returned from FilterValidatorWithdrawInitiated and is used to iterate over the raw logs and unpacked data for ValidatorWithdrawInitiated events raised by the EngineV5 contract.
type EngineV5ValidatorWithdrawInitiatedIterator struct {
	Event *EngineV5ValidatorWithdrawInitiated // Event containing the contract specifics and raw log

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
func (it *EngineV5ValidatorWithdrawInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5ValidatorWithdrawInitiated)
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
		it.Event = new(EngineV5ValidatorWithdrawInitiated)
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
func (it *EngineV5ValidatorWithdrawInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5ValidatorWithdrawInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5ValidatorWithdrawInitiated represents a ValidatorWithdrawInitiated event raised by the EngineV5 contract.
type EngineV5ValidatorWithdrawInitiated struct {
	Addr       common.Address
	Count      *big.Int
	UnlockTime *big.Int
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterValidatorWithdrawInitiated is a free log retrieval operation binding the contract event 0xcc7e0e76b20394ef965d71c4111d21e1b322bb8775dc2c7acd29eb0d3c3dd96d.
//
// Solidity: event ValidatorWithdrawInitiated(address indexed addr, uint256 indexed count, uint256 unlockTime, uint256 amount)
func (_EngineV5 *EngineV5Filterer) FilterValidatorWithdrawInitiated(opts *bind.FilterOpts, addr []common.Address, count []*big.Int) (*EngineV5ValidatorWithdrawInitiatedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var countRule []interface{}
	for _, countItem := range count {
		countRule = append(countRule, countItem)
	}

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "ValidatorWithdrawInitiated", addrRule, countRule)
	if err != nil {
		return nil, err
	}
	return &EngineV5ValidatorWithdrawInitiatedIterator{contract: _EngineV5.contract, event: "ValidatorWithdrawInitiated", logs: logs, sub: sub}, nil
}

// WatchValidatorWithdrawInitiated is a free log subscription operation binding the contract event 0xcc7e0e76b20394ef965d71c4111d21e1b322bb8775dc2c7acd29eb0d3c3dd96d.
//
// Solidity: event ValidatorWithdrawInitiated(address indexed addr, uint256 indexed count, uint256 unlockTime, uint256 amount)
func (_EngineV5 *EngineV5Filterer) WatchValidatorWithdrawInitiated(opts *bind.WatchOpts, sink chan<- *EngineV5ValidatorWithdrawInitiated, addr []common.Address, count []*big.Int) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var countRule []interface{}
	for _, countItem := range count {
		countRule = append(countRule, countItem)
	}

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "ValidatorWithdrawInitiated", addrRule, countRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5ValidatorWithdrawInitiated)
				if err := _EngineV5.contract.UnpackLog(event, "ValidatorWithdrawInitiated", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseValidatorWithdrawInitiated(log types.Log) (*EngineV5ValidatorWithdrawInitiated, error) {
	event := new(EngineV5ValidatorWithdrawInitiated)
	if err := _EngineV5.contract.UnpackLog(event, "ValidatorWithdrawInitiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EngineV5VersionChangedIterator is returned from FilterVersionChanged and is used to iterate over the raw logs and unpacked data for VersionChanged events raised by the EngineV5 contract.
type EngineV5VersionChangedIterator struct {
	Event *EngineV5VersionChanged // Event containing the contract specifics and raw log

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
func (it *EngineV5VersionChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EngineV5VersionChanged)
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
		it.Event = new(EngineV5VersionChanged)
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
func (it *EngineV5VersionChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EngineV5VersionChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EngineV5VersionChanged represents a VersionChanged event raised by the EngineV5 contract.
type EngineV5VersionChanged struct {
	Version *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterVersionChanged is a free log retrieval operation binding the contract event 0x8c854a81cb5c93e7e482d30fb9c6f88fdbdb320f10f7a853c2263659b54e563f.
//
// Solidity: event VersionChanged(uint256 version)
func (_EngineV5 *EngineV5Filterer) FilterVersionChanged(opts *bind.FilterOpts) (*EngineV5VersionChangedIterator, error) {

	logs, sub, err := _EngineV5.contract.FilterLogs(opts, "VersionChanged")
	if err != nil {
		return nil, err
	}
	return &EngineV5VersionChangedIterator{contract: _EngineV5.contract, event: "VersionChanged", logs: logs, sub: sub}, nil
}

// WatchVersionChanged is a free log subscription operation binding the contract event 0x8c854a81cb5c93e7e482d30fb9c6f88fdbdb320f10f7a853c2263659b54e563f.
//
// Solidity: event VersionChanged(uint256 version)
func (_EngineV5 *EngineV5Filterer) WatchVersionChanged(opts *bind.WatchOpts, sink chan<- *EngineV5VersionChanged) (event.Subscription, error) {

	logs, sub, err := _EngineV5.contract.WatchLogs(opts, "VersionChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EngineV5VersionChanged)
				if err := _EngineV5.contract.UnpackLog(event, "VersionChanged", log); err != nil {
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
func (_EngineV5 *EngineV5Filterer) ParseVersionChanged(log types.Log) (*EngineV5VersionChanged, error) {
	event := new(EngineV5VersionChanged)
	if err := _EngineV5.contract.UnpackLog(event, "VersionChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
