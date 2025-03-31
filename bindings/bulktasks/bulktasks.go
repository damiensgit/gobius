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
	Bin: "0x608060405234801561001057600080fd5b50604051610664380380610664833981810160405281019061003291906101fe565b816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663095ea7b3600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6040518363ffffffff1660e01b8152600401610150929190610266565b6020604051808303816000875af115801561016f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061019391906102c7565b5050506102f4565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006101cb826101a0565b9050919050565b6101db816101c0565b81146101e657600080fd5b50565b6000815190506101f8816101d2565b92915050565b600080604083850312156102155761021461019b565b5b6000610223858286016101e9565b9250506020610234858286016101e9565b9150509250929050565b610247816101c0565b82525050565b6000819050919050565b6102608161024d565b82525050565b600060408201905061027b600083018561023e565b6102886020830184610257565b9392505050565b60008115159050919050565b6102a48161028f565b81146102af57600080fd5b50565b6000815190506102c18161029b565b92915050565b6000602082840312156102dd576102dc61019b565b5b60006102eb848285016102b2565b91505092915050565b610361806103036000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806302546ffc146100515780631023ad7c1461006f57806350ec55e11461008b578063c9d4623f146100a7575b600080fd5b6100596100c5565b6040516100669190610254565b60405180910390f35b610089600480360381019061008491906102de565b6100e9565b005b6100a560048036038101906100a091906102de565b61016b565b005b6100af6101ed565b6040516100bc9190610254565b60405180910390f35b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690506040517f77286d1700000000000000000000000000000000000000000000000000000000815236600101845b818110156101635760208160048501376000806024856000885af15060208101905061013c565b505050505050565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690506040517f506ea7de00000000000000000000000000000000000000000000000000000000815236600101845b818110156101e55760208160048501376000806024856000885af1506020810190506101be565b505050505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061023e82610213565b9050919050565b61024e81610233565b82525050565b60006020820190506102696000830184610245565b92915050565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f84011261029e5761029d610279565b5b8235905067ffffffffffffffff8111156102bb576102ba61027e565b5b6020830191508360208202830111156102d7576102d6610283565b5b9250929050565b600080602083850312156102f5576102f461026f565b5b600083013567ffffffffffffffff81111561031357610312610274565b5b61031f85828601610288565b9250925050925092905056fea26469706673582212208721b2bc725d99bbbf34a59dde4c824d8a337741b04dd997a8cf91cf4db88cc164736f6c63430008130033",
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
