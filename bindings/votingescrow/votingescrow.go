// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package votingescrow

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

// IVotingEscrowPoint is an auto generated low-level Go binding around an user-defined struct.
type IVotingEscrowPoint struct {
	Bias  *big.Int
	Slope *big.Int
	Ts    *big.Int
	Blk   *big.Int
}

// VotingEscrowMetaData contains all meta data concerning the VotingEscrow contract.
var VotingEscrowMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"abstain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"attach\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"balanceOfNFT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"create_lock_for\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"deposit_for\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"detach\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"isApprovedOrOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"loc\",\"type\":\"uint256\"}],\"name\":\"point_history\",\"outputs\":[{\"components\":[{\"internalType\":\"int128\",\"name\":\"bias\",\"type\":\"int128\"},{\"internalType\":\"int128\",\"name\":\"slope\",\"type\":\"int128\"},{\"internalType\":\"uint256\",\"name\":\"ts\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blk\",\"type\":\"uint256\"}],\"internalType\":\"structIVotingEscrow.Point\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"team\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tokenOfOwnerByIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"user_point_epoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"loc\",\"type\":\"uint256\"}],\"name\":\"user_point_history\",\"outputs\":[{\"components\":[{\"internalType\":\"int128\",\"name\":\"bias\",\"type\":\"int128\"},{\"internalType\":\"int128\",\"name\":\"slope\",\"type\":\"int128\"},{\"internalType\":\"uint256\",\"name\":\"ts\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blk\",\"type\":\"uint256\"}],\"internalType\":\"structIVotingEscrow.Point\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"voting\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// VotingEscrowABI is the input ABI used to generate the binding from.
// Deprecated: Use VotingEscrowMetaData.ABI instead.
var VotingEscrowABI = VotingEscrowMetaData.ABI

// VotingEscrow is an auto generated Go binding around an Ethereum contract.
type VotingEscrow struct {
	VotingEscrowCaller     // Read-only binding to the contract
	VotingEscrowTransactor // Write-only binding to the contract
	VotingEscrowFilterer   // Log filterer for contract events
}

// VotingEscrowCaller is an auto generated read-only Go binding around an Ethereum contract.
type VotingEscrowCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VotingEscrowTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VotingEscrowTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VotingEscrowFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VotingEscrowFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VotingEscrowSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VotingEscrowSession struct {
	Contract     *VotingEscrow     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VotingEscrowCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VotingEscrowCallerSession struct {
	Contract *VotingEscrowCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// VotingEscrowTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VotingEscrowTransactorSession struct {
	Contract     *VotingEscrowTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// VotingEscrowRaw is an auto generated low-level Go binding around an Ethereum contract.
type VotingEscrowRaw struct {
	Contract *VotingEscrow // Generic contract binding to access the raw methods on
}

// VotingEscrowCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VotingEscrowCallerRaw struct {
	Contract *VotingEscrowCaller // Generic read-only contract binding to access the raw methods on
}

// VotingEscrowTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VotingEscrowTransactorRaw struct {
	Contract *VotingEscrowTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVotingEscrow creates a new instance of VotingEscrow, bound to a specific deployed contract.
func NewVotingEscrow(address common.Address, backend bind.ContractBackend) (*VotingEscrow, error) {
	contract, err := bindVotingEscrow(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VotingEscrow{VotingEscrowCaller: VotingEscrowCaller{contract: contract}, VotingEscrowTransactor: VotingEscrowTransactor{contract: contract}, VotingEscrowFilterer: VotingEscrowFilterer{contract: contract}}, nil
}

// NewVotingEscrowCaller creates a new read-only instance of VotingEscrow, bound to a specific deployed contract.
func NewVotingEscrowCaller(address common.Address, caller bind.ContractCaller) (*VotingEscrowCaller, error) {
	contract, err := bindVotingEscrow(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VotingEscrowCaller{contract: contract}, nil
}

// NewVotingEscrowTransactor creates a new write-only instance of VotingEscrow, bound to a specific deployed contract.
func NewVotingEscrowTransactor(address common.Address, transactor bind.ContractTransactor) (*VotingEscrowTransactor, error) {
	contract, err := bindVotingEscrow(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VotingEscrowTransactor{contract: contract}, nil
}

// NewVotingEscrowFilterer creates a new log filterer instance of VotingEscrow, bound to a specific deployed contract.
func NewVotingEscrowFilterer(address common.Address, filterer bind.ContractFilterer) (*VotingEscrowFilterer, error) {
	contract, err := bindVotingEscrow(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VotingEscrowFilterer{contract: contract}, nil
}

// bindVotingEscrow binds a generic wrapper to an already deployed contract.
func bindVotingEscrow(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VotingEscrowMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VotingEscrow *VotingEscrowRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VotingEscrow.Contract.VotingEscrowCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VotingEscrow *VotingEscrowRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VotingEscrow.Contract.VotingEscrowTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VotingEscrow *VotingEscrowRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VotingEscrow.Contract.VotingEscrowTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VotingEscrow *VotingEscrowCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VotingEscrow.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VotingEscrow *VotingEscrowTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VotingEscrow.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VotingEscrow *VotingEscrowTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VotingEscrow.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_VotingEscrow *VotingEscrowCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VotingEscrow.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_VotingEscrow *VotingEscrowSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _VotingEscrow.Contract.BalanceOf(&_VotingEscrow.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_VotingEscrow *VotingEscrowCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _VotingEscrow.Contract.BalanceOf(&_VotingEscrow.CallOpts, arg0)
}

// BalanceOfNFT is a free data retrieval call binding the contract method 0xe7e242d4.
//
// Solidity: function balanceOfNFT(uint256 ) view returns(uint256)
func (_VotingEscrow *VotingEscrowCaller) BalanceOfNFT(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VotingEscrow.contract.Call(opts, &out, "balanceOfNFT", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOfNFT is a free data retrieval call binding the contract method 0xe7e242d4.
//
// Solidity: function balanceOfNFT(uint256 ) view returns(uint256)
func (_VotingEscrow *VotingEscrowSession) BalanceOfNFT(arg0 *big.Int) (*big.Int, error) {
	return _VotingEscrow.Contract.BalanceOfNFT(&_VotingEscrow.CallOpts, arg0)
}

// BalanceOfNFT is a free data retrieval call binding the contract method 0xe7e242d4.
//
// Solidity: function balanceOfNFT(uint256 ) view returns(uint256)
func (_VotingEscrow *VotingEscrowCallerSession) BalanceOfNFT(arg0 *big.Int) (*big.Int, error) {
	return _VotingEscrow.Contract.BalanceOfNFT(&_VotingEscrow.CallOpts, arg0)
}

// Epoch is a free data retrieval call binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() view returns(uint256)
func (_VotingEscrow *VotingEscrowCaller) Epoch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VotingEscrow.contract.Call(opts, &out, "epoch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Epoch is a free data retrieval call binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() view returns(uint256)
func (_VotingEscrow *VotingEscrowSession) Epoch() (*big.Int, error) {
	return _VotingEscrow.Contract.Epoch(&_VotingEscrow.CallOpts)
}

// Epoch is a free data retrieval call binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() view returns(uint256)
func (_VotingEscrow *VotingEscrowCallerSession) Epoch() (*big.Int, error) {
	return _VotingEscrow.Contract.Epoch(&_VotingEscrow.CallOpts)
}

// IsApprovedOrOwner is a free data retrieval call binding the contract method 0x430c2081.
//
// Solidity: function isApprovedOrOwner(address , uint256 ) view returns(bool)
func (_VotingEscrow *VotingEscrowCaller) IsApprovedOrOwner(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (bool, error) {
	var out []interface{}
	err := _VotingEscrow.contract.Call(opts, &out, "isApprovedOrOwner", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedOrOwner is a free data retrieval call binding the contract method 0x430c2081.
//
// Solidity: function isApprovedOrOwner(address , uint256 ) view returns(bool)
func (_VotingEscrow *VotingEscrowSession) IsApprovedOrOwner(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _VotingEscrow.Contract.IsApprovedOrOwner(&_VotingEscrow.CallOpts, arg0, arg1)
}

// IsApprovedOrOwner is a free data retrieval call binding the contract method 0x430c2081.
//
// Solidity: function isApprovedOrOwner(address , uint256 ) view returns(bool)
func (_VotingEscrow *VotingEscrowCallerSession) IsApprovedOrOwner(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _VotingEscrow.Contract.IsApprovedOrOwner(&_VotingEscrow.CallOpts, arg0, arg1)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 ) view returns(address)
func (_VotingEscrow *VotingEscrowCaller) OwnerOf(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _VotingEscrow.contract.Call(opts, &out, "ownerOf", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 ) view returns(address)
func (_VotingEscrow *VotingEscrowSession) OwnerOf(arg0 *big.Int) (common.Address, error) {
	return _VotingEscrow.Contract.OwnerOf(&_VotingEscrow.CallOpts, arg0)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 ) view returns(address)
func (_VotingEscrow *VotingEscrowCallerSession) OwnerOf(arg0 *big.Int) (common.Address, error) {
	return _VotingEscrow.Contract.OwnerOf(&_VotingEscrow.CallOpts, arg0)
}

// PointHistory is a free data retrieval call binding the contract method 0xd1febfb9.
//
// Solidity: function point_history(uint256 loc) view returns((int128,int128,uint256,uint256))
func (_VotingEscrow *VotingEscrowCaller) PointHistory(opts *bind.CallOpts, loc *big.Int) (IVotingEscrowPoint, error) {
	var out []interface{}
	err := _VotingEscrow.contract.Call(opts, &out, "point_history", loc)

	if err != nil {
		return *new(IVotingEscrowPoint), err
	}

	out0 := *abi.ConvertType(out[0], new(IVotingEscrowPoint)).(*IVotingEscrowPoint)

	return out0, err

}

// PointHistory is a free data retrieval call binding the contract method 0xd1febfb9.
//
// Solidity: function point_history(uint256 loc) view returns((int128,int128,uint256,uint256))
func (_VotingEscrow *VotingEscrowSession) PointHistory(loc *big.Int) (IVotingEscrowPoint, error) {
	return _VotingEscrow.Contract.PointHistory(&_VotingEscrow.CallOpts, loc)
}

// PointHistory is a free data retrieval call binding the contract method 0xd1febfb9.
//
// Solidity: function point_history(uint256 loc) view returns((int128,int128,uint256,uint256))
func (_VotingEscrow *VotingEscrowCallerSession) PointHistory(loc *big.Int) (IVotingEscrowPoint, error) {
	return _VotingEscrow.Contract.PointHistory(&_VotingEscrow.CallOpts, loc)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_VotingEscrow *VotingEscrowCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VotingEscrow.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_VotingEscrow *VotingEscrowSession) Token() (common.Address, error) {
	return _VotingEscrow.Contract.Token(&_VotingEscrow.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_VotingEscrow *VotingEscrowCallerSession) Token() (common.Address, error) {
	return _VotingEscrow.Contract.Token(&_VotingEscrow.CallOpts)
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address , uint256 ) view returns(uint256)
func (_VotingEscrow *VotingEscrowCaller) TokenOfOwnerByIndex(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VotingEscrow.contract.Call(opts, &out, "tokenOfOwnerByIndex", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address , uint256 ) view returns(uint256)
func (_VotingEscrow *VotingEscrowSession) TokenOfOwnerByIndex(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _VotingEscrow.Contract.TokenOfOwnerByIndex(&_VotingEscrow.CallOpts, arg0, arg1)
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address , uint256 ) view returns(uint256)
func (_VotingEscrow *VotingEscrowCallerSession) TokenOfOwnerByIndex(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _VotingEscrow.Contract.TokenOfOwnerByIndex(&_VotingEscrow.CallOpts, arg0, arg1)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_VotingEscrow *VotingEscrowCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VotingEscrow.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_VotingEscrow *VotingEscrowSession) TotalSupply() (*big.Int, error) {
	return _VotingEscrow.Contract.TotalSupply(&_VotingEscrow.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_VotingEscrow *VotingEscrowCallerSession) TotalSupply() (*big.Int, error) {
	return _VotingEscrow.Contract.TotalSupply(&_VotingEscrow.CallOpts)
}

// UserPointEpoch is a free data retrieval call binding the contract method 0xe441135c.
//
// Solidity: function user_point_epoch(uint256 tokenId) view returns(uint256)
func (_VotingEscrow *VotingEscrowCaller) UserPointEpoch(opts *bind.CallOpts, tokenId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VotingEscrow.contract.Call(opts, &out, "user_point_epoch", tokenId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UserPointEpoch is a free data retrieval call binding the contract method 0xe441135c.
//
// Solidity: function user_point_epoch(uint256 tokenId) view returns(uint256)
func (_VotingEscrow *VotingEscrowSession) UserPointEpoch(tokenId *big.Int) (*big.Int, error) {
	return _VotingEscrow.Contract.UserPointEpoch(&_VotingEscrow.CallOpts, tokenId)
}

// UserPointEpoch is a free data retrieval call binding the contract method 0xe441135c.
//
// Solidity: function user_point_epoch(uint256 tokenId) view returns(uint256)
func (_VotingEscrow *VotingEscrowCallerSession) UserPointEpoch(tokenId *big.Int) (*big.Int, error) {
	return _VotingEscrow.Contract.UserPointEpoch(&_VotingEscrow.CallOpts, tokenId)
}

// UserPointHistory is a free data retrieval call binding the contract method 0x1376f3da.
//
// Solidity: function user_point_history(uint256 tokenId, uint256 loc) view returns((int128,int128,uint256,uint256))
func (_VotingEscrow *VotingEscrowCaller) UserPointHistory(opts *bind.CallOpts, tokenId *big.Int, loc *big.Int) (IVotingEscrowPoint, error) {
	var out []interface{}
	err := _VotingEscrow.contract.Call(opts, &out, "user_point_history", tokenId, loc)

	if err != nil {
		return *new(IVotingEscrowPoint), err
	}

	out0 := *abi.ConvertType(out[0], new(IVotingEscrowPoint)).(*IVotingEscrowPoint)

	return out0, err

}

// UserPointHistory is a free data retrieval call binding the contract method 0x1376f3da.
//
// Solidity: function user_point_history(uint256 tokenId, uint256 loc) view returns((int128,int128,uint256,uint256))
func (_VotingEscrow *VotingEscrowSession) UserPointHistory(tokenId *big.Int, loc *big.Int) (IVotingEscrowPoint, error) {
	return _VotingEscrow.Contract.UserPointHistory(&_VotingEscrow.CallOpts, tokenId, loc)
}

// UserPointHistory is a free data retrieval call binding the contract method 0x1376f3da.
//
// Solidity: function user_point_history(uint256 tokenId, uint256 loc) view returns((int128,int128,uint256,uint256))
func (_VotingEscrow *VotingEscrowCallerSession) UserPointHistory(tokenId *big.Int, loc *big.Int) (IVotingEscrowPoint, error) {
	return _VotingEscrow.Contract.UserPointHistory(&_VotingEscrow.CallOpts, tokenId, loc)
}

// Abstain is a paid mutator transaction binding the contract method 0xc1f0fb9f.
//
// Solidity: function abstain(uint256 tokenId) returns()
func (_VotingEscrow *VotingEscrowTransactor) Abstain(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.contract.Transact(opts, "abstain", tokenId)
}

// Abstain is a paid mutator transaction binding the contract method 0xc1f0fb9f.
//
// Solidity: function abstain(uint256 tokenId) returns()
func (_VotingEscrow *VotingEscrowSession) Abstain(tokenId *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.Contract.Abstain(&_VotingEscrow.TransactOpts, tokenId)
}

// Abstain is a paid mutator transaction binding the contract method 0xc1f0fb9f.
//
// Solidity: function abstain(uint256 tokenId) returns()
func (_VotingEscrow *VotingEscrowTransactorSession) Abstain(tokenId *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.Contract.Abstain(&_VotingEscrow.TransactOpts, tokenId)
}

// Attach is a paid mutator transaction binding the contract method 0xfbd3a29d.
//
// Solidity: function attach(uint256 tokenId) returns()
func (_VotingEscrow *VotingEscrowTransactor) Attach(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.contract.Transact(opts, "attach", tokenId)
}

// Attach is a paid mutator transaction binding the contract method 0xfbd3a29d.
//
// Solidity: function attach(uint256 tokenId) returns()
func (_VotingEscrow *VotingEscrowSession) Attach(tokenId *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.Contract.Attach(&_VotingEscrow.TransactOpts, tokenId)
}

// Attach is a paid mutator transaction binding the contract method 0xfbd3a29d.
//
// Solidity: function attach(uint256 tokenId) returns()
func (_VotingEscrow *VotingEscrowTransactorSession) Attach(tokenId *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.Contract.Attach(&_VotingEscrow.TransactOpts, tokenId)
}

// Checkpoint is a paid mutator transaction binding the contract method 0xc2c4c5c1.
//
// Solidity: function checkpoint() returns()
func (_VotingEscrow *VotingEscrowTransactor) Checkpoint(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VotingEscrow.contract.Transact(opts, "checkpoint")
}

// Checkpoint is a paid mutator transaction binding the contract method 0xc2c4c5c1.
//
// Solidity: function checkpoint() returns()
func (_VotingEscrow *VotingEscrowSession) Checkpoint() (*types.Transaction, error) {
	return _VotingEscrow.Contract.Checkpoint(&_VotingEscrow.TransactOpts)
}

// Checkpoint is a paid mutator transaction binding the contract method 0xc2c4c5c1.
//
// Solidity: function checkpoint() returns()
func (_VotingEscrow *VotingEscrowTransactorSession) Checkpoint() (*types.Transaction, error) {
	return _VotingEscrow.Contract.Checkpoint(&_VotingEscrow.TransactOpts)
}

// CreateLockFor is a paid mutator transaction binding the contract method 0xd4e54c3b.
//
// Solidity: function create_lock_for(uint256 , uint256 , address ) returns(uint256)
func (_VotingEscrow *VotingEscrowTransactor) CreateLockFor(opts *bind.TransactOpts, arg0 *big.Int, arg1 *big.Int, arg2 common.Address) (*types.Transaction, error) {
	return _VotingEscrow.contract.Transact(opts, "create_lock_for", arg0, arg1, arg2)
}

// CreateLockFor is a paid mutator transaction binding the contract method 0xd4e54c3b.
//
// Solidity: function create_lock_for(uint256 , uint256 , address ) returns(uint256)
func (_VotingEscrow *VotingEscrowSession) CreateLockFor(arg0 *big.Int, arg1 *big.Int, arg2 common.Address) (*types.Transaction, error) {
	return _VotingEscrow.Contract.CreateLockFor(&_VotingEscrow.TransactOpts, arg0, arg1, arg2)
}

// CreateLockFor is a paid mutator transaction binding the contract method 0xd4e54c3b.
//
// Solidity: function create_lock_for(uint256 , uint256 , address ) returns(uint256)
func (_VotingEscrow *VotingEscrowTransactorSession) CreateLockFor(arg0 *big.Int, arg1 *big.Int, arg2 common.Address) (*types.Transaction, error) {
	return _VotingEscrow.Contract.CreateLockFor(&_VotingEscrow.TransactOpts, arg0, arg1, arg2)
}

// DepositFor is a paid mutator transaction binding the contract method 0xee99fe28.
//
// Solidity: function deposit_for(uint256 tokenId, uint256 value) returns()
func (_VotingEscrow *VotingEscrowTransactor) DepositFor(opts *bind.TransactOpts, tokenId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.contract.Transact(opts, "deposit_for", tokenId, value)
}

// DepositFor is a paid mutator transaction binding the contract method 0xee99fe28.
//
// Solidity: function deposit_for(uint256 tokenId, uint256 value) returns()
func (_VotingEscrow *VotingEscrowSession) DepositFor(tokenId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.Contract.DepositFor(&_VotingEscrow.TransactOpts, tokenId, value)
}

// DepositFor is a paid mutator transaction binding the contract method 0xee99fe28.
//
// Solidity: function deposit_for(uint256 tokenId, uint256 value) returns()
func (_VotingEscrow *VotingEscrowTransactorSession) DepositFor(tokenId *big.Int, value *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.Contract.DepositFor(&_VotingEscrow.TransactOpts, tokenId, value)
}

// Detach is a paid mutator transaction binding the contract method 0x986b7d8a.
//
// Solidity: function detach(uint256 tokenId) returns()
func (_VotingEscrow *VotingEscrowTransactor) Detach(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.contract.Transact(opts, "detach", tokenId)
}

// Detach is a paid mutator transaction binding the contract method 0x986b7d8a.
//
// Solidity: function detach(uint256 tokenId) returns()
func (_VotingEscrow *VotingEscrowSession) Detach(tokenId *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.Contract.Detach(&_VotingEscrow.TransactOpts, tokenId)
}

// Detach is a paid mutator transaction binding the contract method 0x986b7d8a.
//
// Solidity: function detach(uint256 tokenId) returns()
func (_VotingEscrow *VotingEscrowTransactorSession) Detach(tokenId *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.Contract.Detach(&_VotingEscrow.TransactOpts, tokenId)
}

// Team is a paid mutator transaction binding the contract method 0x85f2aef2.
//
// Solidity: function team() returns(address)
func (_VotingEscrow *VotingEscrowTransactor) Team(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VotingEscrow.contract.Transact(opts, "team")
}

// Team is a paid mutator transaction binding the contract method 0x85f2aef2.
//
// Solidity: function team() returns(address)
func (_VotingEscrow *VotingEscrowSession) Team() (*types.Transaction, error) {
	return _VotingEscrow.Contract.Team(&_VotingEscrow.TransactOpts)
}

// Team is a paid mutator transaction binding the contract method 0x85f2aef2.
//
// Solidity: function team() returns(address)
func (_VotingEscrow *VotingEscrowTransactorSession) Team() (*types.Transaction, error) {
	return _VotingEscrow.Contract.Team(&_VotingEscrow.TransactOpts)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) returns()
func (_VotingEscrow *VotingEscrowTransactor) TransferFrom(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.contract.Transact(opts, "transferFrom", arg0, arg1, arg2)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) returns()
func (_VotingEscrow *VotingEscrowSession) TransferFrom(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.Contract.TransferFrom(&_VotingEscrow.TransactOpts, arg0, arg1, arg2)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) returns()
func (_VotingEscrow *VotingEscrowTransactorSession) TransferFrom(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.Contract.TransferFrom(&_VotingEscrow.TransactOpts, arg0, arg1, arg2)
}

// Voting is a paid mutator transaction binding the contract method 0xfd4a77f1.
//
// Solidity: function voting(uint256 tokenId) returns()
func (_VotingEscrow *VotingEscrowTransactor) Voting(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.contract.Transact(opts, "voting", tokenId)
}

// Voting is a paid mutator transaction binding the contract method 0xfd4a77f1.
//
// Solidity: function voting(uint256 tokenId) returns()
func (_VotingEscrow *VotingEscrowSession) Voting(tokenId *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.Contract.Voting(&_VotingEscrow.TransactOpts, tokenId)
}

// Voting is a paid mutator transaction binding the contract method 0xfd4a77f1.
//
// Solidity: function voting(uint256 tokenId) returns()
func (_VotingEscrow *VotingEscrowTransactorSession) Voting(tokenId *big.Int) (*types.Transaction, error) {
	return _VotingEscrow.Contract.Voting(&_VotingEscrow.TransactOpts, tokenId)
}
