// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bulktasks

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

// BulkTasksMetaData contains all meta data concerning the BulkTasks contract.
var BulkTasksMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"commitments_\",\"type\":\"bytes32[]\"}],\"name\":\"bulkSignalCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_taskids\",\"type\":\"bytes32[]\"}],\"name\":\"claimSolutions\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50734a24b101728e07a52053c13fb4db2bcf490cabc373ffffffffffffffffffffffffffffffffffffffff1663095ea7b3739b51ef044d3486a1fb0a2d55a6e0ceeadd323e667fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6040518363ffffffff1660e01b8152600401610094929190610137565b6020604051808303816000875af11580156100b3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100d7919061019d565b506101ca565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610108826100dd565b9050919050565b610118816100fd565b82525050565b6000819050919050565b6101318161011e565b82525050565b600060408201905061014c600083018561010f565b6101596020830184610128565b9392505050565b600080fd5b60008115159050919050565b61017a81610165565b811461018557600080fd5b50565b60008151905061019781610171565b92915050565b6000602082840312156101b3576101b2610160565b5b60006101c184828501610188565b91505092915050565b610241806101d96000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80631023ad7c1461003b57806350ec55e114610057575b600080fd5b610055600480360381019061005091906101be565b610073565b005b610071600480360381019061006c91906101be565b6100e1565b005b6040517f77286d1700000000000000000000000000000000000000000000000000000000815236600101835b818110156100da5760208160048501376000806024856000739b51ef044d3486a1fb0a2d55a6e0ceeadd323e665af15060208101905061009f565b5050505050565b6040517f506ea7de00000000000000000000000000000000000000000000000000000000815236600101835b818110156101485760208160048501376000806024856000739b51ef044d3486a1fb0a2d55a6e0ceeadd323e665af15060208101905061010d565b5050505050565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f84011261017e5761017d610159565b5b8235905067ffffffffffffffff81111561019b5761019a61015e565b5b6020830191508360208202830111156101b7576101b6610163565b5b9250929050565b600080602083850312156101d5576101d461014f565b5b600083013567ffffffffffffffff8111156101f3576101f2610154565b5b6101ff85828601610168565b9250925050925092905056fea26469706673582212209ac2b3860861b35960caca6e35e86149b527ba7995d5c085517c7bd49e4997b864736f6c63430008130033",
}

// BulkTasksABI is the input ABI used to generate the binding from.
// Deprecated: Use BulkTasksMetaData.ABI instead.
var BulkTasksABI = BulkTasksMetaData.ABI

// BulkTasksBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BulkTasksMetaData.Bin instead.
var BulkTasksBin = BulkTasksMetaData.Bin

// DeployBulkTasks deploys a new Ethereum contract, binding an instance of BulkTasks to it.
func DeployBulkTasks(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BulkTasks, error) {
	parsed, err := BulkTasksMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BulkTasksBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BulkTasks{BulkTasksCaller: BulkTasksCaller{contract: contract}, BulkTasksTransactor: BulkTasksTransactor{contract: contract}, BulkTasksFilterer: BulkTasksFilterer{contract: contract}}, nil
}

// BulkTasks is an auto generated Go binding around an Ethereum contract.
type BulkTasks struct {
	BulkTasksCaller     // Read-only binding to the contract
	BulkTasksTransactor // Write-only binding to the contract
	BulkTasksFilterer   // Log filterer for contract events
}

// BulkTasksCaller is an auto generated read-only Go binding around an Ethereum contract.
type BulkTasksCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BulkTasksTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BulkTasksTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BulkTasksFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BulkTasksFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BulkTasksSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BulkTasksSession struct {
	Contract     *BulkTasks        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BulkTasksCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BulkTasksCallerSession struct {
	Contract *BulkTasksCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// BulkTasksTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BulkTasksTransactorSession struct {
	Contract     *BulkTasksTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// BulkTasksRaw is an auto generated low-level Go binding around an Ethereum contract.
type BulkTasksRaw struct {
	Contract *BulkTasks // Generic contract binding to access the raw methods on
}

// BulkTasksCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BulkTasksCallerRaw struct {
	Contract *BulkTasksCaller // Generic read-only contract binding to access the raw methods on
}

// BulkTasksTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BulkTasksTransactorRaw struct {
	Contract *BulkTasksTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBulkTasks creates a new instance of BulkTasks, bound to a specific deployed contract.
func NewBulkTasks(address common.Address, backend bind.ContractBackend) (*BulkTasks, error) {
	contract, err := bindBulkTasks(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BulkTasks{BulkTasksCaller: BulkTasksCaller{contract: contract}, BulkTasksTransactor: BulkTasksTransactor{contract: contract}, BulkTasksFilterer: BulkTasksFilterer{contract: contract}}, nil
}

// NewBulkTasksCaller creates a new read-only instance of BulkTasks, bound to a specific deployed contract.
func NewBulkTasksCaller(address common.Address, caller bind.ContractCaller) (*BulkTasksCaller, error) {
	contract, err := bindBulkTasks(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BulkTasksCaller{contract: contract}, nil
}

// NewBulkTasksTransactor creates a new write-only instance of BulkTasks, bound to a specific deployed contract.
func NewBulkTasksTransactor(address common.Address, transactor bind.ContractTransactor) (*BulkTasksTransactor, error) {
	contract, err := bindBulkTasks(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BulkTasksTransactor{contract: contract}, nil
}

// NewBulkTasksFilterer creates a new log filterer instance of BulkTasks, bound to a specific deployed contract.
func NewBulkTasksFilterer(address common.Address, filterer bind.ContractFilterer) (*BulkTasksFilterer, error) {
	contract, err := bindBulkTasks(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BulkTasksFilterer{contract: contract}, nil
}

// bindBulkTasks binds a generic wrapper to an already deployed contract.
func bindBulkTasks(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BulkTasksMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BulkTasks *BulkTasksRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BulkTasks.Contract.BulkTasksCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BulkTasks *BulkTasksRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BulkTasks.Contract.BulkTasksTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BulkTasks *BulkTasksRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BulkTasks.Contract.BulkTasksTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BulkTasks *BulkTasksCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BulkTasks.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BulkTasks *BulkTasksTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BulkTasks.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BulkTasks *BulkTasksTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BulkTasks.Contract.contract.Transact(opts, method, params...)
}

// BulkSignalCommitment is a paid mutator transaction binding the contract method 0x50ec55e1.
//
// Solidity: function bulkSignalCommitment(bytes32[] commitments_) returns()
func (_BulkTasks *BulkTasksTransactor) BulkSignalCommitment(opts *bind.TransactOpts, commitments_ [][32]byte) (*types.Transaction, error) {
	return _BulkTasks.contract.Transact(opts, "bulkSignalCommitment", commitments_)
}

// BulkSignalCommitment is a paid mutator transaction binding the contract method 0x50ec55e1.
//
// Solidity: function bulkSignalCommitment(bytes32[] commitments_) returns()
func (_BulkTasks *BulkTasksSession) BulkSignalCommitment(commitments_ [][32]byte) (*types.Transaction, error) {
	return _BulkTasks.Contract.BulkSignalCommitment(&_BulkTasks.TransactOpts, commitments_)
}

// BulkSignalCommitment is a paid mutator transaction binding the contract method 0x50ec55e1.
//
// Solidity: function bulkSignalCommitment(bytes32[] commitments_) returns()
func (_BulkTasks *BulkTasksTransactorSession) BulkSignalCommitment(commitments_ [][32]byte) (*types.Transaction, error) {
	return _BulkTasks.Contract.BulkSignalCommitment(&_BulkTasks.TransactOpts, commitments_)
}

// ClaimSolutions is a paid mutator transaction binding the contract method 0x1023ad7c.
//
// Solidity: function claimSolutions(bytes32[] _taskids) returns()
func (_BulkTasks *BulkTasksTransactor) ClaimSolutions(opts *bind.TransactOpts, _taskids [][32]byte) (*types.Transaction, error) {
	return _BulkTasks.contract.Transact(opts, "claimSolutions", _taskids)
}

// ClaimSolutions is a paid mutator transaction binding the contract method 0x1023ad7c.
//
// Solidity: function claimSolutions(bytes32[] _taskids) returns()
func (_BulkTasks *BulkTasksSession) ClaimSolutions(_taskids [][32]byte) (*types.Transaction, error) {
	return _BulkTasks.Contract.ClaimSolutions(&_BulkTasks.TransactOpts, _taskids)
}

// ClaimSolutions is a paid mutator transaction binding the contract method 0x1023ad7c.
//
// Solidity: function claimSolutions(bytes32[] _taskids) returns()
func (_BulkTasks *BulkTasksTransactorSession) ClaimSolutions(_taskids [][32]byte) (*types.Transaction, error) {
	return _BulkTasks.Contract.ClaimSolutions(&_BulkTasks.TransactOpts, _taskids)
}
