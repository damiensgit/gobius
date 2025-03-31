// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vestaking

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

// VeStakingMetaData contains all meta data concerning the VeStaking contract.
var VeStakingMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"_stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"}],\"name\":\"_updateBalance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"_withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"account\",\"type\":\"uint256\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"account\",\"type\":\"uint256\"}],\"name\":\"earned\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getReward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRewardForDuration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastTimeRewardApplicable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastUpdateTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"notifyRewardAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"periodFinish\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardPerToken\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rewardPerTokenPaid\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardPerTokenStored\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardsDuration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// VeStakingABI is the input ABI used to generate the binding from.
// Deprecated: Use VeStakingMetaData.ABI instead.
var VeStakingABI = VeStakingMetaData.ABI

// VeStaking is an auto generated Go binding around an Ethereum contract.
type VeStaking struct {
	VeStakingCaller     // Read-only binding to the contract
	VeStakingTransactor // Write-only binding to the contract
	VeStakingFilterer   // Log filterer for contract events
}

// VeStakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type VeStakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VeStakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VeStakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VeStakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VeStakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VeStakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VeStakingSession struct {
	Contract     *VeStaking        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VeStakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VeStakingCallerSession struct {
	Contract *VeStakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// VeStakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VeStakingTransactorSession struct {
	Contract     *VeStakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// VeStakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type VeStakingRaw struct {
	Contract *VeStaking // Generic contract binding to access the raw methods on
}

// VeStakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VeStakingCallerRaw struct {
	Contract *VeStakingCaller // Generic read-only contract binding to access the raw methods on
}

// VeStakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VeStakingTransactorRaw struct {
	Contract *VeStakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVeStaking creates a new instance of VeStaking, bound to a specific deployed contract.
func NewVeStaking(address common.Address, backend bind.ContractBackend) (*VeStaking, error) {
	contract, err := bindVeStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VeStaking{VeStakingCaller: VeStakingCaller{contract: contract}, VeStakingTransactor: VeStakingTransactor{contract: contract}, VeStakingFilterer: VeStakingFilterer{contract: contract}}, nil
}

// NewVeStakingCaller creates a new read-only instance of VeStaking, bound to a specific deployed contract.
func NewVeStakingCaller(address common.Address, caller bind.ContractCaller) (*VeStakingCaller, error) {
	contract, err := bindVeStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VeStakingCaller{contract: contract}, nil
}

// NewVeStakingTransactor creates a new write-only instance of VeStaking, bound to a specific deployed contract.
func NewVeStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*VeStakingTransactor, error) {
	contract, err := bindVeStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VeStakingTransactor{contract: contract}, nil
}

// NewVeStakingFilterer creates a new log filterer instance of VeStaking, bound to a specific deployed contract.
func NewVeStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*VeStakingFilterer, error) {
	contract, err := bindVeStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VeStakingFilterer{contract: contract}, nil
}

// bindVeStaking binds a generic wrapper to an already deployed contract.
func bindVeStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VeStakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VeStaking *VeStakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VeStaking.Contract.VeStakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VeStaking *VeStakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VeStaking.Contract.VeStakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VeStaking *VeStakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VeStaking.Contract.VeStakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VeStaking *VeStakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VeStaking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VeStaking *VeStakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VeStaking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VeStaking *VeStakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VeStaking.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x9cc7f708.
//
// Solidity: function balanceOf(uint256 account) view returns(uint256)
func (_VeStaking *VeStakingCaller) BalanceOf(opts *bind.CallOpts, account *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VeStaking.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x9cc7f708.
//
// Solidity: function balanceOf(uint256 account) view returns(uint256)
func (_VeStaking *VeStakingSession) BalanceOf(account *big.Int) (*big.Int, error) {
	return _VeStaking.Contract.BalanceOf(&_VeStaking.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x9cc7f708.
//
// Solidity: function balanceOf(uint256 account) view returns(uint256)
func (_VeStaking *VeStakingCallerSession) BalanceOf(account *big.Int) (*big.Int, error) {
	return _VeStaking.Contract.BalanceOf(&_VeStaking.CallOpts, account)
}

// Earned is a free data retrieval call binding the contract method 0x4d6ed8c4.
//
// Solidity: function earned(uint256 account) view returns(uint256)
func (_VeStaking *VeStakingCaller) Earned(opts *bind.CallOpts, account *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VeStaking.contract.Call(opts, &out, "earned", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Earned is a free data retrieval call binding the contract method 0x4d6ed8c4.
//
// Solidity: function earned(uint256 account) view returns(uint256)
func (_VeStaking *VeStakingSession) Earned(account *big.Int) (*big.Int, error) {
	return _VeStaking.Contract.Earned(&_VeStaking.CallOpts, account)
}

// Earned is a free data retrieval call binding the contract method 0x4d6ed8c4.
//
// Solidity: function earned(uint256 account) view returns(uint256)
func (_VeStaking *VeStakingCallerSession) Earned(account *big.Int) (*big.Int, error) {
	return _VeStaking.Contract.Earned(&_VeStaking.CallOpts, account)
}

// GetRewardForDuration is a free data retrieval call binding the contract method 0x1c1f78eb.
//
// Solidity: function getRewardForDuration() view returns(uint256)
func (_VeStaking *VeStakingCaller) GetRewardForDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeStaking.contract.Call(opts, &out, "getRewardForDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRewardForDuration is a free data retrieval call binding the contract method 0x1c1f78eb.
//
// Solidity: function getRewardForDuration() view returns(uint256)
func (_VeStaking *VeStakingSession) GetRewardForDuration() (*big.Int, error) {
	return _VeStaking.Contract.GetRewardForDuration(&_VeStaking.CallOpts)
}

// GetRewardForDuration is a free data retrieval call binding the contract method 0x1c1f78eb.
//
// Solidity: function getRewardForDuration() view returns(uint256)
func (_VeStaking *VeStakingCallerSession) GetRewardForDuration() (*big.Int, error) {
	return _VeStaking.Contract.GetRewardForDuration(&_VeStaking.CallOpts)
}

// LastTimeRewardApplicable is a free data retrieval call binding the contract method 0x80faa57d.
//
// Solidity: function lastTimeRewardApplicable() view returns(uint256)
func (_VeStaking *VeStakingCaller) LastTimeRewardApplicable(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeStaking.contract.Call(opts, &out, "lastTimeRewardApplicable")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastTimeRewardApplicable is a free data retrieval call binding the contract method 0x80faa57d.
//
// Solidity: function lastTimeRewardApplicable() view returns(uint256)
func (_VeStaking *VeStakingSession) LastTimeRewardApplicable() (*big.Int, error) {
	return _VeStaking.Contract.LastTimeRewardApplicable(&_VeStaking.CallOpts)
}

// LastTimeRewardApplicable is a free data retrieval call binding the contract method 0x80faa57d.
//
// Solidity: function lastTimeRewardApplicable() view returns(uint256)
func (_VeStaking *VeStakingCallerSession) LastTimeRewardApplicable() (*big.Int, error) {
	return _VeStaking.Contract.LastTimeRewardApplicable(&_VeStaking.CallOpts)
}

// LastUpdateTime is a free data retrieval call binding the contract method 0xc8f33c91.
//
// Solidity: function lastUpdateTime() view returns(uint256)
func (_VeStaking *VeStakingCaller) LastUpdateTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeStaking.contract.Call(opts, &out, "lastUpdateTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastUpdateTime is a free data retrieval call binding the contract method 0xc8f33c91.
//
// Solidity: function lastUpdateTime() view returns(uint256)
func (_VeStaking *VeStakingSession) LastUpdateTime() (*big.Int, error) {
	return _VeStaking.Contract.LastUpdateTime(&_VeStaking.CallOpts)
}

// LastUpdateTime is a free data retrieval call binding the contract method 0xc8f33c91.
//
// Solidity: function lastUpdateTime() view returns(uint256)
func (_VeStaking *VeStakingCallerSession) LastUpdateTime() (*big.Int, error) {
	return _VeStaking.Contract.LastUpdateTime(&_VeStaking.CallOpts)
}

// PeriodFinish is a free data retrieval call binding the contract method 0xebe2b12b.
//
// Solidity: function periodFinish() view returns(uint256)
func (_VeStaking *VeStakingCaller) PeriodFinish(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeStaking.contract.Call(opts, &out, "periodFinish")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PeriodFinish is a free data retrieval call binding the contract method 0xebe2b12b.
//
// Solidity: function periodFinish() view returns(uint256)
func (_VeStaking *VeStakingSession) PeriodFinish() (*big.Int, error) {
	return _VeStaking.Contract.PeriodFinish(&_VeStaking.CallOpts)
}

// PeriodFinish is a free data retrieval call binding the contract method 0xebe2b12b.
//
// Solidity: function periodFinish() view returns(uint256)
func (_VeStaking *VeStakingCallerSession) PeriodFinish() (*big.Int, error) {
	return _VeStaking.Contract.PeriodFinish(&_VeStaking.CallOpts)
}

// RewardPerToken is a free data retrieval call binding the contract method 0xcd3daf9d.
//
// Solidity: function rewardPerToken() view returns(uint256)
func (_VeStaking *VeStakingCaller) RewardPerToken(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeStaking.contract.Call(opts, &out, "rewardPerToken")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RewardPerToken is a free data retrieval call binding the contract method 0xcd3daf9d.
//
// Solidity: function rewardPerToken() view returns(uint256)
func (_VeStaking *VeStakingSession) RewardPerToken() (*big.Int, error) {
	return _VeStaking.Contract.RewardPerToken(&_VeStaking.CallOpts)
}

// RewardPerToken is a free data retrieval call binding the contract method 0xcd3daf9d.
//
// Solidity: function rewardPerToken() view returns(uint256)
func (_VeStaking *VeStakingCallerSession) RewardPerToken() (*big.Int, error) {
	return _VeStaking.Contract.RewardPerToken(&_VeStaking.CallOpts)
}

// RewardPerTokenPaid is a free data retrieval call binding the contract method 0xd0779da8.
//
// Solidity: function rewardPerTokenPaid(uint256 ) view returns(uint256)
func (_VeStaking *VeStakingCaller) RewardPerTokenPaid(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VeStaking.contract.Call(opts, &out, "rewardPerTokenPaid", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RewardPerTokenPaid is a free data retrieval call binding the contract method 0xd0779da8.
//
// Solidity: function rewardPerTokenPaid(uint256 ) view returns(uint256)
func (_VeStaking *VeStakingSession) RewardPerTokenPaid(arg0 *big.Int) (*big.Int, error) {
	return _VeStaking.Contract.RewardPerTokenPaid(&_VeStaking.CallOpts, arg0)
}

// RewardPerTokenPaid is a free data retrieval call binding the contract method 0xd0779da8.
//
// Solidity: function rewardPerTokenPaid(uint256 ) view returns(uint256)
func (_VeStaking *VeStakingCallerSession) RewardPerTokenPaid(arg0 *big.Int) (*big.Int, error) {
	return _VeStaking.Contract.RewardPerTokenPaid(&_VeStaking.CallOpts, arg0)
}

// RewardPerTokenStored is a free data retrieval call binding the contract method 0xdf136d65.
//
// Solidity: function rewardPerTokenStored() view returns(uint256)
func (_VeStaking *VeStakingCaller) RewardPerTokenStored(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeStaking.contract.Call(opts, &out, "rewardPerTokenStored")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RewardPerTokenStored is a free data retrieval call binding the contract method 0xdf136d65.
//
// Solidity: function rewardPerTokenStored() view returns(uint256)
func (_VeStaking *VeStakingSession) RewardPerTokenStored() (*big.Int, error) {
	return _VeStaking.Contract.RewardPerTokenStored(&_VeStaking.CallOpts)
}

// RewardPerTokenStored is a free data retrieval call binding the contract method 0xdf136d65.
//
// Solidity: function rewardPerTokenStored() view returns(uint256)
func (_VeStaking *VeStakingCallerSession) RewardPerTokenStored() (*big.Int, error) {
	return _VeStaking.Contract.RewardPerTokenStored(&_VeStaking.CallOpts)
}

// RewardRate is a free data retrieval call binding the contract method 0x7b0a47ee.
//
// Solidity: function rewardRate() view returns(uint256)
func (_VeStaking *VeStakingCaller) RewardRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeStaking.contract.Call(opts, &out, "rewardRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RewardRate is a free data retrieval call binding the contract method 0x7b0a47ee.
//
// Solidity: function rewardRate() view returns(uint256)
func (_VeStaking *VeStakingSession) RewardRate() (*big.Int, error) {
	return _VeStaking.Contract.RewardRate(&_VeStaking.CallOpts)
}

// RewardRate is a free data retrieval call binding the contract method 0x7b0a47ee.
//
// Solidity: function rewardRate() view returns(uint256)
func (_VeStaking *VeStakingCallerSession) RewardRate() (*big.Int, error) {
	return _VeStaking.Contract.RewardRate(&_VeStaking.CallOpts)
}

// Rewards is a free data retrieval call binding the contract method 0xf301af42.
//
// Solidity: function rewards(uint256 ) view returns(uint256)
func (_VeStaking *VeStakingCaller) Rewards(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VeStaking.contract.Call(opts, &out, "rewards", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Rewards is a free data retrieval call binding the contract method 0xf301af42.
//
// Solidity: function rewards(uint256 ) view returns(uint256)
func (_VeStaking *VeStakingSession) Rewards(arg0 *big.Int) (*big.Int, error) {
	return _VeStaking.Contract.Rewards(&_VeStaking.CallOpts, arg0)
}

// Rewards is a free data retrieval call binding the contract method 0xf301af42.
//
// Solidity: function rewards(uint256 ) view returns(uint256)
func (_VeStaking *VeStakingCallerSession) Rewards(arg0 *big.Int) (*big.Int, error) {
	return _VeStaking.Contract.Rewards(&_VeStaking.CallOpts, arg0)
}

// RewardsDuration is a free data retrieval call binding the contract method 0x386a9525.
//
// Solidity: function rewardsDuration() view returns(uint256)
func (_VeStaking *VeStakingCaller) RewardsDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeStaking.contract.Call(opts, &out, "rewardsDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RewardsDuration is a free data retrieval call binding the contract method 0x386a9525.
//
// Solidity: function rewardsDuration() view returns(uint256)
func (_VeStaking *VeStakingSession) RewardsDuration() (*big.Int, error) {
	return _VeStaking.Contract.RewardsDuration(&_VeStaking.CallOpts)
}

// RewardsDuration is a free data retrieval call binding the contract method 0x386a9525.
//
// Solidity: function rewardsDuration() view returns(uint256)
func (_VeStaking *VeStakingCallerSession) RewardsDuration() (*big.Int, error) {
	return _VeStaking.Contract.RewardsDuration(&_VeStaking.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_VeStaking *VeStakingCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeStaking.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_VeStaking *VeStakingSession) TotalSupply() (*big.Int, error) {
	return _VeStaking.Contract.TotalSupply(&_VeStaking.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_VeStaking *VeStakingCallerSession) TotalSupply() (*big.Int, error) {
	return _VeStaking.Contract.TotalSupply(&_VeStaking.CallOpts)
}

// Stake is a paid mutator transaction binding the contract method 0x61dc2c36.
//
// Solidity: function _stake(uint256 tokenId, uint256 amount) returns()
func (_VeStaking *VeStakingTransactor) Stake(opts *bind.TransactOpts, tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VeStaking.contract.Transact(opts, "_stake", tokenId, amount)
}

// Stake is a paid mutator transaction binding the contract method 0x61dc2c36.
//
// Solidity: function _stake(uint256 tokenId, uint256 amount) returns()
func (_VeStaking *VeStakingSession) Stake(tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VeStaking.Contract.Stake(&_VeStaking.TransactOpts, tokenId, amount)
}

// Stake is a paid mutator transaction binding the contract method 0x61dc2c36.
//
// Solidity: function _stake(uint256 tokenId, uint256 amount) returns()
func (_VeStaking *VeStakingTransactorSession) Stake(tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VeStaking.Contract.Stake(&_VeStaking.TransactOpts, tokenId, amount)
}

// UpdateBalance is a paid mutator transaction binding the contract method 0xdad3aab6.
//
// Solidity: function _updateBalance(uint256 tokenId, uint256 newAmount) returns()
func (_VeStaking *VeStakingTransactor) UpdateBalance(opts *bind.TransactOpts, tokenId *big.Int, newAmount *big.Int) (*types.Transaction, error) {
	return _VeStaking.contract.Transact(opts, "_updateBalance", tokenId, newAmount)
}

// UpdateBalance is a paid mutator transaction binding the contract method 0xdad3aab6.
//
// Solidity: function _updateBalance(uint256 tokenId, uint256 newAmount) returns()
func (_VeStaking *VeStakingSession) UpdateBalance(tokenId *big.Int, newAmount *big.Int) (*types.Transaction, error) {
	return _VeStaking.Contract.UpdateBalance(&_VeStaking.TransactOpts, tokenId, newAmount)
}

// UpdateBalance is a paid mutator transaction binding the contract method 0xdad3aab6.
//
// Solidity: function _updateBalance(uint256 tokenId, uint256 newAmount) returns()
func (_VeStaking *VeStakingTransactorSession) UpdateBalance(tokenId *big.Int, newAmount *big.Int) (*types.Transaction, error) {
	return _VeStaking.Contract.UpdateBalance(&_VeStaking.TransactOpts, tokenId, newAmount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xac6a2b5d.
//
// Solidity: function _withdraw(uint256 amount) returns()
func (_VeStaking *VeStakingTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VeStaking.contract.Transact(opts, "_withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xac6a2b5d.
//
// Solidity: function _withdraw(uint256 amount) returns()
func (_VeStaking *VeStakingSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _VeStaking.Contract.Withdraw(&_VeStaking.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xac6a2b5d.
//
// Solidity: function _withdraw(uint256 amount) returns()
func (_VeStaking *VeStakingTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _VeStaking.Contract.Withdraw(&_VeStaking.TransactOpts, amount)
}

// GetReward is a paid mutator transaction binding the contract method 0x1c4b774b.
//
// Solidity: function getReward(uint256 tokenId) returns()
func (_VeStaking *VeStakingTransactor) GetReward(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _VeStaking.contract.Transact(opts, "getReward", tokenId)
}

// GetReward is a paid mutator transaction binding the contract method 0x1c4b774b.
//
// Solidity: function getReward(uint256 tokenId) returns()
func (_VeStaking *VeStakingSession) GetReward(tokenId *big.Int) (*types.Transaction, error) {
	return _VeStaking.Contract.GetReward(&_VeStaking.TransactOpts, tokenId)
}

// GetReward is a paid mutator transaction binding the contract method 0x1c4b774b.
//
// Solidity: function getReward(uint256 tokenId) returns()
func (_VeStaking *VeStakingTransactorSession) GetReward(tokenId *big.Int) (*types.Transaction, error) {
	return _VeStaking.Contract.GetReward(&_VeStaking.TransactOpts, tokenId)
}

// NotifyRewardAmount is a paid mutator transaction binding the contract method 0x3c6b16ab.
//
// Solidity: function notifyRewardAmount(uint256 reward) returns()
func (_VeStaking *VeStakingTransactor) NotifyRewardAmount(opts *bind.TransactOpts, reward *big.Int) (*types.Transaction, error) {
	return _VeStaking.contract.Transact(opts, "notifyRewardAmount", reward)
}

// NotifyRewardAmount is a paid mutator transaction binding the contract method 0x3c6b16ab.
//
// Solidity: function notifyRewardAmount(uint256 reward) returns()
func (_VeStaking *VeStakingSession) NotifyRewardAmount(reward *big.Int) (*types.Transaction, error) {
	return _VeStaking.Contract.NotifyRewardAmount(&_VeStaking.TransactOpts, reward)
}

// NotifyRewardAmount is a paid mutator transaction binding the contract method 0x3c6b16ab.
//
// Solidity: function notifyRewardAmount(uint256 reward) returns()
func (_VeStaking *VeStakingTransactorSession) NotifyRewardAmount(reward *big.Int) (*types.Transaction, error) {
	return _VeStaking.Contract.NotifyRewardAmount(&_VeStaking.TransactOpts, reward)
}
