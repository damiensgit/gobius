// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package quoter

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

// QuoterMetaData contains all meta data concerning the Quoter contract.
var QuoterMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_profitLevel\",\"type\":\"uint256\"}],\"name\":\"ProfitLevelChanged\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"GetAIUSPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"profitLevel\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// QuoterABI is the input ABI used to generate the binding from.
// Deprecated: Use QuoterMetaData.ABI instead.
var QuoterABI = QuoterMetaData.ABI

// Quoter is an auto generated Go binding around an Ethereum contract.
type Quoter struct {
	QuoterCaller     // Read-only binding to the contract
	QuoterTransactor // Write-only binding to the contract
	QuoterFilterer   // Log filterer for contract events
}

// QuoterCaller is an auto generated read-only Go binding around an Ethereum contract.
type QuoterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuoterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type QuoterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuoterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type QuoterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuoterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type QuoterSession struct {
	Contract     *Quoter           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// QuoterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type QuoterCallerSession struct {
	Contract *QuoterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// QuoterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type QuoterTransactorSession struct {
	Contract     *QuoterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// QuoterRaw is an auto generated low-level Go binding around an Ethereum contract.
type QuoterRaw struct {
	Contract *Quoter // Generic contract binding to access the raw methods on
}

// QuoterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type QuoterCallerRaw struct {
	Contract *QuoterCaller // Generic read-only contract binding to access the raw methods on
}

// QuoterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type QuoterTransactorRaw struct {
	Contract *QuoterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewQuoter creates a new instance of Quoter, bound to a specific deployed contract.
func NewQuoter(address common.Address, backend bind.ContractBackend) (*Quoter, error) {
	contract, err := bindQuoter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Quoter{QuoterCaller: QuoterCaller{contract: contract}, QuoterTransactor: QuoterTransactor{contract: contract}, QuoterFilterer: QuoterFilterer{contract: contract}}, nil
}

// NewQuoterCaller creates a new read-only instance of Quoter, bound to a specific deployed contract.
func NewQuoterCaller(address common.Address, caller bind.ContractCaller) (*QuoterCaller, error) {
	contract, err := bindQuoter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &QuoterCaller{contract: contract}, nil
}

// NewQuoterTransactor creates a new write-only instance of Quoter, bound to a specific deployed contract.
func NewQuoterTransactor(address common.Address, transactor bind.ContractTransactor) (*QuoterTransactor, error) {
	contract, err := bindQuoter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &QuoterTransactor{contract: contract}, nil
}

// NewQuoterFilterer creates a new log filterer instance of Quoter, bound to a specific deployed contract.
func NewQuoterFilterer(address common.Address, filterer bind.ContractFilterer) (*QuoterFilterer, error) {
	contract, err := bindQuoter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &QuoterFilterer{contract: contract}, nil
}

// bindQuoter binds a generic wrapper to an already deployed contract.
func bindQuoter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := QuoterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Quoter *QuoterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Quoter.Contract.QuoterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Quoter *QuoterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Quoter.Contract.QuoterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Quoter *QuoterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Quoter.Contract.QuoterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Quoter *QuoterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Quoter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Quoter *QuoterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Quoter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Quoter *QuoterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Quoter.Contract.contract.Transact(opts, method, params...)
}

// ProfitLevel is a free data retrieval call binding the contract method 0x659963ab.
//
// Solidity: function profitLevel() view returns(uint256)
func (_Quoter *QuoterCaller) ProfitLevel(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Quoter.contract.Call(opts, &out, "profitLevel")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProfitLevel is a free data retrieval call binding the contract method 0x659963ab.
//
// Solidity: function profitLevel() view returns(uint256)
func (_Quoter *QuoterSession) ProfitLevel() (*big.Int, error) {
	return _Quoter.Contract.ProfitLevel(&_Quoter.CallOpts)
}

// ProfitLevel is a free data retrieval call binding the contract method 0x659963ab.
//
// Solidity: function profitLevel() view returns(uint256)
func (_Quoter *QuoterCallerSession) ProfitLevel() (*big.Int, error) {
	return _Quoter.Contract.ProfitLevel(&_Quoter.CallOpts)
}

// GetAIUSPrice is a paid mutator transaction binding the contract method 0x3585b2d2.
//
// Solidity: function GetAIUSPrice() returns(uint256, uint256)
func (_Quoter *QuoterTransactor) GetAIUSPrice(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Quoter.contract.Transact(opts, "GetAIUSPrice")
}

// GetAIUSPrice is a paid mutator transaction binding the contract method 0x3585b2d2.
//
// Solidity: function GetAIUSPrice() returns(uint256, uint256)
func (_Quoter *QuoterSession) GetAIUSPrice() (*types.Transaction, error) {
	return _Quoter.Contract.GetAIUSPrice(&_Quoter.TransactOpts)
}

// GetAIUSPrice is a paid mutator transaction binding the contract method 0x3585b2d2.
//
// Solidity: function GetAIUSPrice() returns(uint256, uint256)
func (_Quoter *QuoterTransactorSession) GetAIUSPrice() (*types.Transaction, error) {
	return _Quoter.Contract.GetAIUSPrice(&_Quoter.TransactOpts)
}

// QuoterProfitLevelChangedIterator is returned from FilterProfitLevelChanged and is used to iterate over the raw logs and unpacked data for ProfitLevelChanged events raised by the Quoter contract.
type QuoterProfitLevelChangedIterator struct {
	Event *QuoterProfitLevelChanged // Event containing the contract specifics and raw log

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
func (it *QuoterProfitLevelChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(QuoterProfitLevelChanged)
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
		it.Event = new(QuoterProfitLevelChanged)
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
func (it *QuoterProfitLevelChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *QuoterProfitLevelChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// QuoterProfitLevelChanged represents a ProfitLevelChanged event raised by the Quoter contract.
type QuoterProfitLevelChanged struct {
	ProfitLevel *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterProfitLevelChanged is a free log retrieval operation binding the contract event 0x75f45b324d8a27b3aa9d81842214231f946d89ac77f70889ab5b9f701e541172.
//
// Solidity: event ProfitLevelChanged(uint256 _profitLevel)
func (_Quoter *QuoterFilterer) FilterProfitLevelChanged(opts *bind.FilterOpts) (*QuoterProfitLevelChangedIterator, error) {

	logs, sub, err := _Quoter.contract.FilterLogs(opts, "ProfitLevelChanged")
	if err != nil {
		return nil, err
	}
	return &QuoterProfitLevelChangedIterator{contract: _Quoter.contract, event: "ProfitLevelChanged", logs: logs, sub: sub}, nil
}

// WatchProfitLevelChanged is a free log subscription operation binding the contract event 0x75f45b324d8a27b3aa9d81842214231f946d89ac77f70889ab5b9f701e541172.
//
// Solidity: event ProfitLevelChanged(uint256 _profitLevel)
func (_Quoter *QuoterFilterer) WatchProfitLevelChanged(opts *bind.WatchOpts, sink chan<- *QuoterProfitLevelChanged) (event.Subscription, error) {

	logs, sub, err := _Quoter.contract.WatchLogs(opts, "ProfitLevelChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(QuoterProfitLevelChanged)
				if err := _Quoter.contract.UnpackLog(event, "ProfitLevelChanged", log); err != nil {
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

// ParseProfitLevelChanged is a log parse operation binding the contract event 0x75f45b324d8a27b3aa9d81842214231f946d89ac77f70889ab5b9f701e541172.
//
// Solidity: event ProfitLevelChanged(uint256 _profitLevel)
func (_Quoter *QuoterFilterer) ParseProfitLevelChanged(log types.Log) (*QuoterProfitLevelChanged, error) {
	event := new(QuoterProfitLevelChanged)
	if err := _Quoter.contract.UnpackLog(event, "ProfitLevelChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
