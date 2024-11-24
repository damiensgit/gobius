// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bulkclaim

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

// BulkclaimMetaData contains all meta data concerning the Bulkclaim contract.
var BulkclaimMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIArbius\",\"name\":\"_arbius\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_taskids\",\"type\":\"bytes32[]\"}],\"name\":\"claimSolutions\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161040f38038061040f833981810160405281019061003291906100ed565b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505061011a565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100a88261007d565b9050919050565b60006100ba8261009d565b9050919050565b6100ca816100af565b81146100d557600080fd5b50565b6000815190506100e7816100c1565b92915050565b60006020828403121561010357610102610078565b5b6000610111848285016100d8565b91505092915050565b6102e6806101296000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80631023ad7c14610030575b600080fd5b61004a6004803603810190610045919061017f565b61004c565b005b60005b8282905081101561010b5760008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166377286d178484848181106100a9576100a86101cc565b5b905060200201356040518263ffffffff1660e01b81526004016100cc9190610214565b600060405180830381600087803b1580156100e657600080fd5b505af19250505080156100f7575060015b50808061010390610268565b91505061004f565b505050565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f84011261013f5761013e61011a565b5b8235905067ffffffffffffffff81111561015c5761015b61011f565b5b60208301915083602082028301111561017857610177610124565b5b9250929050565b6000806020838503121561019657610195610110565b5b600083013567ffffffffffffffff8111156101b4576101b3610115565b5b6101c085828601610129565b92509250509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000819050919050565b61020e816101fb565b82525050565b60006020820190506102296000830184610205565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000819050919050565b60006102738261025e565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036102a5576102a461022f565b5b60018201905091905056fea26469706673582212204d5f77ddddef647aeb5cac513f31db5de1537749cf9b1d145dae1d58367dc58b64736f6c63430008150033",
}

// BulkclaimABI is the input ABI used to generate the binding from.
// Deprecated: Use BulkclaimMetaData.ABI instead.
var BulkclaimABI = BulkclaimMetaData.ABI

// BulkclaimBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BulkclaimMetaData.Bin instead.
var BulkclaimBin = BulkclaimMetaData.Bin

// DeployBulkclaim deploys a new Ethereum contract, binding an instance of Bulkclaim to it.
func DeployBulkclaim(auth *bind.TransactOpts, backend bind.ContractBackend, _arbius common.Address) (common.Address, *types.Transaction, *Bulkclaim, error) {
	parsed, err := BulkclaimMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BulkclaimBin), backend, _arbius)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Bulkclaim{BulkclaimCaller: BulkclaimCaller{contract: contract}, BulkclaimTransactor: BulkclaimTransactor{contract: contract}, BulkclaimFilterer: BulkclaimFilterer{contract: contract}}, nil
}

// Bulkclaim is an auto generated Go binding around an Ethereum contract.
type Bulkclaim struct {
	BulkclaimCaller     // Read-only binding to the contract
	BulkclaimTransactor // Write-only binding to the contract
	BulkclaimFilterer   // Log filterer for contract events
}

// BulkclaimCaller is an auto generated read-only Go binding around an Ethereum contract.
type BulkclaimCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BulkclaimTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BulkclaimTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BulkclaimFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BulkclaimFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BulkclaimSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BulkclaimSession struct {
	Contract     *Bulkclaim        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BulkclaimCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BulkclaimCallerSession struct {
	Contract *BulkclaimCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// BulkclaimTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BulkclaimTransactorSession struct {
	Contract     *BulkclaimTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// BulkclaimRaw is an auto generated low-level Go binding around an Ethereum contract.
type BulkclaimRaw struct {
	Contract *Bulkclaim // Generic contract binding to access the raw methods on
}

// BulkclaimCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BulkclaimCallerRaw struct {
	Contract *BulkclaimCaller // Generic read-only contract binding to access the raw methods on
}

// BulkclaimTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BulkclaimTransactorRaw struct {
	Contract *BulkclaimTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBulkclaim creates a new instance of Bulkclaim, bound to a specific deployed contract.
func NewBulkclaim(address common.Address, backend bind.ContractBackend) (*Bulkclaim, error) {
	contract, err := bindBulkclaim(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bulkclaim{BulkclaimCaller: BulkclaimCaller{contract: contract}, BulkclaimTransactor: BulkclaimTransactor{contract: contract}, BulkclaimFilterer: BulkclaimFilterer{contract: contract}}, nil
}

// NewBulkclaimCaller creates a new read-only instance of Bulkclaim, bound to a specific deployed contract.
func NewBulkclaimCaller(address common.Address, caller bind.ContractCaller) (*BulkclaimCaller, error) {
	contract, err := bindBulkclaim(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BulkclaimCaller{contract: contract}, nil
}

// NewBulkclaimTransactor creates a new write-only instance of Bulkclaim, bound to a specific deployed contract.
func NewBulkclaimTransactor(address common.Address, transactor bind.ContractTransactor) (*BulkclaimTransactor, error) {
	contract, err := bindBulkclaim(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BulkclaimTransactor{contract: contract}, nil
}

// NewBulkclaimFilterer creates a new log filterer instance of Bulkclaim, bound to a specific deployed contract.
func NewBulkclaimFilterer(address common.Address, filterer bind.ContractFilterer) (*BulkclaimFilterer, error) {
	contract, err := bindBulkclaim(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BulkclaimFilterer{contract: contract}, nil
}

// bindBulkclaim binds a generic wrapper to an already deployed contract.
func bindBulkclaim(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BulkclaimMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bulkclaim *BulkclaimRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bulkclaim.Contract.BulkclaimCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bulkclaim *BulkclaimRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bulkclaim.Contract.BulkclaimTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bulkclaim *BulkclaimRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bulkclaim.Contract.BulkclaimTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bulkclaim *BulkclaimCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bulkclaim.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bulkclaim *BulkclaimTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bulkclaim.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bulkclaim *BulkclaimTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bulkclaim.Contract.contract.Transact(opts, method, params...)
}

// ClaimSolutions is a paid mutator transaction binding the contract method 0x1023ad7c.
//
// Solidity: function claimSolutions(bytes32[] _taskids) returns()
func (_Bulkclaim *BulkclaimTransactor) ClaimSolutions(opts *bind.TransactOpts, _taskids [][32]byte) (*types.Transaction, error) {
	return _Bulkclaim.contract.Transact(opts, "claimSolutions", _taskids)
}

// ClaimSolutions is a paid mutator transaction binding the contract method 0x1023ad7c.
//
// Solidity: function claimSolutions(bytes32[] _taskids) returns()
func (_Bulkclaim *BulkclaimSession) ClaimSolutions(_taskids [][32]byte) (*types.Transaction, error) {
	return _Bulkclaim.Contract.ClaimSolutions(&_Bulkclaim.TransactOpts, _taskids)
}

// ClaimSolutions is a paid mutator transaction binding the contract method 0x1023ad7c.
//
// Solidity: function claimSolutions(bytes32[] _taskids) returns()
func (_Bulkclaim *BulkclaimTransactorSession) ClaimSolutions(_taskids [][32]byte) (*types.Transaction, error) {
	return _Bulkclaim.Contract.ClaimSolutions(&_Bulkclaim.TransactOpts, _taskids)
}
