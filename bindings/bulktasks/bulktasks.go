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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"basetoken_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"engine_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"basetoken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"commitments_\",\"type\":\"bytes32[]\"}],\"name\":\"bulkSignalCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_taskids\",\"type\":\"bytes32[]\"}],\"name\":\"claimSolutions\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"engine\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608080604052346100f057604081610321803803809161001f82856100f5565b8339810103126100f05761003e60206100378361012e565b920161012e565b602060018060a01b03809316604460018060a01b0319600095869584838854161787551680916001541617600155604051948593849263095ea7b360e01b8452600484015260001960248401525af180156100e5576100a7575b6040516101de90816101438239f35b6020813d82116100dd575b816100bf602093836100f5565b810103126100d95751801515036100d65780610098565b80fd5b5080fd5b3d91506100b2565b6040513d84823e3d90fd5b600080fd5b601f909101601f19168101906001600160401b0382119082101761011857604052565b634e487b7160e01b600052604160045260246000fd5b51906001600160a01b03821682036100f05756fe6080604052600436101561001257600080fd5b6000803560e01c90816302546ffc146100b15781631023ad7c146100555750806350ec55e1146100505763c9d4623f1461004b57600080fd5b61017f565b610127565b346100ae57610063366100d5565b50600180546040516377286d1760e01b815236909201926001600160a01b039091169060048301905b848110610097578580f35b80602080928437868060248782885af1500161008c565b80fd5b346100ae57806003193601126100ae57546001600160a01b03166080908152602090f35b906020600319830112610122576001600160401b03916004359083821161012257806023830112156101225781600401359384116101225760248460051b83010111610122576024019190565b600080fd5b3461012257610135366100d5565b506001805460405163283753ef60e11b81526004810193369093019290916001600160a01b0316905b83811061016757005b8060208092873760008060248682875af1500161015e565b34610122576000366003190112610122576001546040516001600160a01b039091168152602090f3fea26469706673582212205ac14e7a0bc9b26ff5e1aa16faefc0da18920397854999a41f73ce8d3accabbb64736f6c63430008130033",
}

// BulkTasksABI is the input ABI used to generate the binding from.
// Deprecated: Use BulkTasksMetaData.ABI instead.
var BulkTasksABI = BulkTasksMetaData.ABI

// BulkTasksBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BulkTasksMetaData.Bin instead.
var BulkTasksBin = BulkTasksMetaData.Bin

// DeployBulkTasks deploys a new Ethereum contract, binding an instance of BulkTasks to it.
func DeployBulkTasks(auth *bind.TransactOpts, backend bind.ContractBackend, basetoken_ common.Address, engine_ common.Address) (common.Address, *types.Transaction, *BulkTasks, error) {
	parsed, err := BulkTasksMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BulkTasksBin), backend, basetoken_, engine_)
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

// Basetoken is a free data retrieval call binding the contract method 0x02546ffc.
//
// Solidity: function basetoken() view returns(address)
func (_BulkTasks *BulkTasksCaller) Basetoken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BulkTasks.contract.Call(opts, &out, "basetoken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Basetoken is a free data retrieval call binding the contract method 0x02546ffc.
//
// Solidity: function basetoken() view returns(address)
func (_BulkTasks *BulkTasksSession) Basetoken() (common.Address, error) {
	return _BulkTasks.Contract.Basetoken(&_BulkTasks.CallOpts)
}

// Basetoken is a free data retrieval call binding the contract method 0x02546ffc.
//
// Solidity: function basetoken() view returns(address)
func (_BulkTasks *BulkTasksCallerSession) Basetoken() (common.Address, error) {
	return _BulkTasks.Contract.Basetoken(&_BulkTasks.CallOpts)
}

// Engine is a free data retrieval call binding the contract method 0xc9d4623f.
//
// Solidity: function engine() view returns(address)
func (_BulkTasks *BulkTasksCaller) Engine(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BulkTasks.contract.Call(opts, &out, "engine")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Engine is a free data retrieval call binding the contract method 0xc9d4623f.
//
// Solidity: function engine() view returns(address)
func (_BulkTasks *BulkTasksSession) Engine() (common.Address, error) {
	return _BulkTasks.Contract.Engine(&_BulkTasks.CallOpts)
}

// Engine is a free data retrieval call binding the contract method 0xc9d4623f.
//
// Solidity: function engine() view returns(address)
func (_BulkTasks *BulkTasksCallerSession) Engine() (common.Address, error) {
	return _BulkTasks.Contract.Engine(&_BulkTasks.CallOpts)
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
