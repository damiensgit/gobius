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

// BulktasksMetaData contains all meta data concerning the Bulktasks contract.
var BulktasksMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"commitments_\",\"type\":\"bytes32[]\"}],\"name\":\"bulkSignalCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_taskids\",\"type\":\"bytes32[]\"}],\"name\":\"claimSolutions\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_tasks\",\"type\":\"uint8\"}],\"name\":\"submitMultipleTasks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_tasks\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"encodedCall\",\"type\":\"bytes\"}],\"name\":\"submitMultipleTasksEncoded\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550738afe4055ebc86bd2afb3940c0095c9aca511d85273ffffffffffffffffffffffffffffffffffffffff1663095ea7b3733bf6050327fa280ee1b5f3e8fd5ea2efe8a6472a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6040518363ffffffff1660e01b81526004016100d4929190610177565b6020604051808303816000875af11580156100f3573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061011791906101dd565b5061020a565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006101488261011d565b9050919050565b6101588161013d565b82525050565b6000819050919050565b6101718161015e565b82525050565b600060408201905061018c600083018561014f565b6101996020830184610168565b9392505050565b600080fd5b60008115159050919050565b6101ba816101a5565b81146101c557600080fd5b50565b6000815190506101d7816101b1565b92915050565b6000602082840312156101f3576101f26101a0565b5b6000610201848285016101c8565b91505092915050565b61063a806102196000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80631023ad7c1461005157806350ec55e11461006d578063dc0143cc14610089578063e5491b6c146100a5575b600080fd5b61006b60048036038101906100669190610465565b6100c1565b005b61008760048036038101906100829190610465565b61012f565b005b6100a3600480360381019061009e91906104eb565b61019d565b005b6100bf60048036038101906100ba91906105a4565b6103b3565b005b6040517f77286d1700000000000000000000000000000000000000000000000000000000815236600101835b818110156101285760208160048501376000806024856000733bf6050327fa280ee1b5f3e8fd5ea2efe8a6472a5af1506020810190506100ed565b5050505050565b6040517f506ea7de00000000000000000000000000000000000000000000000000000000815236600101835b818110156101965760208160048501376000806024856000733bf6050327fa280ee1b5f3e8fd5ea2efe8a6472a5af15060208101905061015b565b5050505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690506040517f08745dd1000000000000000000000000000000000000000000000000000000008152600060048201528160248201527f98617a8cd4a11db63100ad44bea4e5e296aecfd78b2ef06aee3e364c7307f21260448201526000606482015260a06084820152601260a48201527f7b2270726f6d7074223a22626c61636b227d000000000000000000000000000060c482015260005b838110156103ad5760008060e8846000733bf6050327fa280ee1b5f3e8fd5ea2efe8a6472a5af15060008060e8846000733bf6050327fa280ee1b5f3e8fd5ea2efe8a6472a5af15060008060e8846000733bf6050327fa280ee1b5f3e8fd5ea2efe8a6472a5af15060008060e8846000733bf6050327fa280ee1b5f3e8fd5ea2efe8a6472a5af15060008060e8846000733bf6050327fa280ee1b5f3e8fd5ea2efe8a6472a5af15060008060e8846000733bf6050327fa280ee1b5f3e8fd5ea2efe8a6472a5af15060008060e8846000733bf6050327fa280ee1b5f3e8fd5ea2efe8a6472a5af15060008060e8846000733bf6050327fa280ee1b5f3e8fd5ea2efe8a6472a5af15060008060e8846000733bf6050327fa280ee1b5f3e8fd5ea2efe8a6472a5af15060008060e8846000733bf6050327fa280ee1b5f3e8fd5ea2efe8a6472a5af15060018101905061025a565b50505050565b6040518183823760005b848110156103ef5760008084846000733bf6050327fa280ee1b5f3e8fd5ea2efe8a6472a5af1506001810190506103bd565b5050505050565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f84011261042557610424610400565b5b8235905067ffffffffffffffff81111561044257610441610405565b5b60208301915083602082028301111561045e5761045d61040a565b5b9250929050565b6000806020838503121561047c5761047b6103f6565b5b600083013567ffffffffffffffff81111561049a576104996103fb565b5b6104a68582860161040f565b92509250509250929050565b600060ff82169050919050565b6104c8816104b2565b81146104d357600080fd5b50565b6000813590506104e5816104bf565b92915050565b600060208284031215610501576105006103f6565b5b600061050f848285016104d6565b91505092915050565b6000819050919050565b61052b81610518565b811461053657600080fd5b50565b60008135905061054881610522565b92915050565b60008083601f84011261056457610563610400565b5b8235905067ffffffffffffffff81111561058157610580610405565b5b60208301915083600182028301111561059d5761059c61040a565b5b9250929050565b6000806000604084860312156105bd576105bc6103f6565b5b60006105cb86828701610539565b935050602084013567ffffffffffffffff8111156105ec576105eb6103fb565b5b6105f88682870161054e565b9250925050925092509256fea2646970667358221220b26cf8c820d25fdcaaed360634ee9fde127d89f7ba4c8da793ba31f983c5217564736f6c63430008150033",
}

// BulktasksABI is the input ABI used to generate the binding from.
// Deprecated: Use BulktasksMetaData.ABI instead.
var BulktasksABI = BulktasksMetaData.ABI

// BulktasksBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BulktasksMetaData.Bin instead.
var BulktasksBin = BulktasksMetaData.Bin

// DeployBulktasks deploys a new Ethereum contract, binding an instance of Bulktasks to it.
func DeployBulktasks(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Bulktasks, error) {
	parsed, err := BulktasksMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BulktasksBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Bulktasks{BulktasksCaller: BulktasksCaller{contract: contract}, BulktasksTransactor: BulktasksTransactor{contract: contract}, BulktasksFilterer: BulktasksFilterer{contract: contract}}, nil
}

// Bulktasks is an auto generated Go binding around an Ethereum contract.
type Bulktasks struct {
	BulktasksCaller     // Read-only binding to the contract
	BulktasksTransactor // Write-only binding to the contract
	BulktasksFilterer   // Log filterer for contract events
}

// BulktasksCaller is an auto generated read-only Go binding around an Ethereum contract.
type BulktasksCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BulktasksTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BulktasksTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BulktasksFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BulktasksFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BulktasksSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BulktasksSession struct {
	Contract     *Bulktasks        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BulktasksCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BulktasksCallerSession struct {
	Contract *BulktasksCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// BulktasksTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BulktasksTransactorSession struct {
	Contract     *BulktasksTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// BulktasksRaw is an auto generated low-level Go binding around an Ethereum contract.
type BulktasksRaw struct {
	Contract *Bulktasks // Generic contract binding to access the raw methods on
}

// BulktasksCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BulktasksCallerRaw struct {
	Contract *BulktasksCaller // Generic read-only contract binding to access the raw methods on
}

// BulktasksTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BulktasksTransactorRaw struct {
	Contract *BulktasksTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBulktasks creates a new instance of Bulktasks, bound to a specific deployed contract.
func NewBulktasks(address common.Address, backend bind.ContractBackend) (*Bulktasks, error) {
	contract, err := bindBulktasks(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bulktasks{BulktasksCaller: BulktasksCaller{contract: contract}, BulktasksTransactor: BulktasksTransactor{contract: contract}, BulktasksFilterer: BulktasksFilterer{contract: contract}}, nil
}

// NewBulktasksCaller creates a new read-only instance of Bulktasks, bound to a specific deployed contract.
func NewBulktasksCaller(address common.Address, caller bind.ContractCaller) (*BulktasksCaller, error) {
	contract, err := bindBulktasks(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BulktasksCaller{contract: contract}, nil
}

// NewBulktasksTransactor creates a new write-only instance of Bulktasks, bound to a specific deployed contract.
func NewBulktasksTransactor(address common.Address, transactor bind.ContractTransactor) (*BulktasksTransactor, error) {
	contract, err := bindBulktasks(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BulktasksTransactor{contract: contract}, nil
}

// NewBulktasksFilterer creates a new log filterer instance of Bulktasks, bound to a specific deployed contract.
func NewBulktasksFilterer(address common.Address, filterer bind.ContractFilterer) (*BulktasksFilterer, error) {
	contract, err := bindBulktasks(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BulktasksFilterer{contract: contract}, nil
}

// bindBulktasks binds a generic wrapper to an already deployed contract.
func bindBulktasks(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BulktasksMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bulktasks *BulktasksRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bulktasks.Contract.BulktasksCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bulktasks *BulktasksRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bulktasks.Contract.BulktasksTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bulktasks *BulktasksRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bulktasks.Contract.BulktasksTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bulktasks *BulktasksCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bulktasks.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bulktasks *BulktasksTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bulktasks.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bulktasks *BulktasksTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bulktasks.Contract.contract.Transact(opts, method, params...)
}

// BulkSignalCommitment is a paid mutator transaction binding the contract method 0x50ec55e1.
//
// Solidity: function bulkSignalCommitment(bytes32[] commitments_) returns()
func (_Bulktasks *BulktasksTransactor) BulkSignalCommitment(opts *bind.TransactOpts, commitments_ [][32]byte) (*types.Transaction, error) {
	return _Bulktasks.contract.Transact(opts, "bulkSignalCommitment", commitments_)
}

// BulkSignalCommitment is a paid mutator transaction binding the contract method 0x50ec55e1.
//
// Solidity: function bulkSignalCommitment(bytes32[] commitments_) returns()
func (_Bulktasks *BulktasksSession) BulkSignalCommitment(commitments_ [][32]byte) (*types.Transaction, error) {
	return _Bulktasks.Contract.BulkSignalCommitment(&_Bulktasks.TransactOpts, commitments_)
}

// BulkSignalCommitment is a paid mutator transaction binding the contract method 0x50ec55e1.
//
// Solidity: function bulkSignalCommitment(bytes32[] commitments_) returns()
func (_Bulktasks *BulktasksTransactorSession) BulkSignalCommitment(commitments_ [][32]byte) (*types.Transaction, error) {
	return _Bulktasks.Contract.BulkSignalCommitment(&_Bulktasks.TransactOpts, commitments_)
}

// ClaimSolutions is a paid mutator transaction binding the contract method 0x1023ad7c.
//
// Solidity: function claimSolutions(bytes32[] _taskids) returns()
func (_Bulktasks *BulktasksTransactor) ClaimSolutions(opts *bind.TransactOpts, _taskids [][32]byte) (*types.Transaction, error) {
	return _Bulktasks.contract.Transact(opts, "claimSolutions", _taskids)
}

// ClaimSolutions is a paid mutator transaction binding the contract method 0x1023ad7c.
//
// Solidity: function claimSolutions(bytes32[] _taskids) returns()
func (_Bulktasks *BulktasksSession) ClaimSolutions(_taskids [][32]byte) (*types.Transaction, error) {
	return _Bulktasks.Contract.ClaimSolutions(&_Bulktasks.TransactOpts, _taskids)
}

// ClaimSolutions is a paid mutator transaction binding the contract method 0x1023ad7c.
//
// Solidity: function claimSolutions(bytes32[] _taskids) returns()
func (_Bulktasks *BulktasksTransactorSession) ClaimSolutions(_taskids [][32]byte) (*types.Transaction, error) {
	return _Bulktasks.Contract.ClaimSolutions(&_Bulktasks.TransactOpts, _taskids)
}

// SubmitMultipleTasks is a paid mutator transaction binding the contract method 0xdc0143cc.
//
// Solidity: function submitMultipleTasks(uint8 _tasks) returns()
func (_Bulktasks *BulktasksTransactor) SubmitMultipleTasks(opts *bind.TransactOpts, _tasks uint8) (*types.Transaction, error) {
	return _Bulktasks.contract.Transact(opts, "submitMultipleTasks", _tasks)
}

// SubmitMultipleTasks is a paid mutator transaction binding the contract method 0xdc0143cc.
//
// Solidity: function submitMultipleTasks(uint8 _tasks) returns()
func (_Bulktasks *BulktasksSession) SubmitMultipleTasks(_tasks uint8) (*types.Transaction, error) {
	return _Bulktasks.Contract.SubmitMultipleTasks(&_Bulktasks.TransactOpts, _tasks)
}

// SubmitMultipleTasks is a paid mutator transaction binding the contract method 0xdc0143cc.
//
// Solidity: function submitMultipleTasks(uint8 _tasks) returns()
func (_Bulktasks *BulktasksTransactorSession) SubmitMultipleTasks(_tasks uint8) (*types.Transaction, error) {
	return _Bulktasks.Contract.SubmitMultipleTasks(&_Bulktasks.TransactOpts, _tasks)
}

// SubmitMultipleTasksEncoded is a paid mutator transaction binding the contract method 0xe5491b6c.
//
// Solidity: function submitMultipleTasksEncoded(int256 _tasks, bytes encodedCall) returns()
func (_Bulktasks *BulktasksTransactor) SubmitMultipleTasksEncoded(opts *bind.TransactOpts, _tasks *big.Int, encodedCall []byte) (*types.Transaction, error) {
	return _Bulktasks.contract.Transact(opts, "submitMultipleTasksEncoded", _tasks, encodedCall)
}

// SubmitMultipleTasksEncoded is a paid mutator transaction binding the contract method 0xe5491b6c.
//
// Solidity: function submitMultipleTasksEncoded(int256 _tasks, bytes encodedCall) returns()
func (_Bulktasks *BulktasksSession) SubmitMultipleTasksEncoded(_tasks *big.Int, encodedCall []byte) (*types.Transaction, error) {
	return _Bulktasks.Contract.SubmitMultipleTasksEncoded(&_Bulktasks.TransactOpts, _tasks, encodedCall)
}

// SubmitMultipleTasksEncoded is a paid mutator transaction binding the contract method 0xe5491b6c.
//
// Solidity: function submitMultipleTasksEncoded(int256 _tasks, bytes encodedCall) returns()
func (_Bulktasks *BulktasksTransactorSession) SubmitMultipleTasksEncoded(_tasks *big.Int, encodedCall []byte) (*types.Transaction, error) {
	return _Bulktasks.Contract.SubmitMultipleTasksEncoded(&_Bulktasks.TransactOpts, _tasks, encodedCall)
}
