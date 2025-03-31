// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package voter

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

// VoterMetaData contains all meta data concerning the Voter contract.
var VoterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_model\",\"type\":\"bytes32\"}],\"name\":\"createGauge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochVoteEnd\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_model\",\"type\":\"bytes32\"}],\"name\":\"getGaugeMultiplier\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"isAlive\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"isGauge\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"isWhitelisted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_model\",\"type\":\"bytes32\"}],\"name\":\"killGauge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lastVoted\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"length\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"modelVote\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"models\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"poke\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_model\",\"type\":\"bytes32\"}],\"name\":\"reviveGauge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalWeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"usedWeights\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"_modelVote\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_weights\",\"type\":\"uint256[]\"}],\"name\":\"vote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"votes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"votingEscrow\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"weights\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_model\",\"type\":\"bytes32\"}],\"name\":\"whitelist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// VoterABI is the input ABI used to generate the binding from.
// Deprecated: Use VoterMetaData.ABI instead.
var VoterABI = VoterMetaData.ABI

// Voter is an auto generated Go binding around an Ethereum contract.
type Voter struct {
	VoterCaller     // Read-only binding to the contract
	VoterTransactor // Write-only binding to the contract
	VoterFilterer   // Log filterer for contract events
}

// VoterCaller is an auto generated read-only Go binding around an Ethereum contract.
type VoterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VoterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VoterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VoterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VoterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VoterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VoterSession struct {
	Contract     *Voter            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VoterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VoterCallerSession struct {
	Contract *VoterCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// VoterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VoterTransactorSession struct {
	Contract     *VoterTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VoterRaw is an auto generated low-level Go binding around an Ethereum contract.
type VoterRaw struct {
	Contract *Voter // Generic contract binding to access the raw methods on
}

// VoterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VoterCallerRaw struct {
	Contract *VoterCaller // Generic read-only contract binding to access the raw methods on
}

// VoterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VoterTransactorRaw struct {
	Contract *VoterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVoter creates a new instance of Voter, bound to a specific deployed contract.
func NewVoter(address common.Address, backend bind.ContractBackend) (*Voter, error) {
	contract, err := bindVoter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Voter{VoterCaller: VoterCaller{contract: contract}, VoterTransactor: VoterTransactor{contract: contract}, VoterFilterer: VoterFilterer{contract: contract}}, nil
}

// NewVoterCaller creates a new read-only instance of Voter, bound to a specific deployed contract.
func NewVoterCaller(address common.Address, caller bind.ContractCaller) (*VoterCaller, error) {
	contract, err := bindVoter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VoterCaller{contract: contract}, nil
}

// NewVoterTransactor creates a new write-only instance of Voter, bound to a specific deployed contract.
func NewVoterTransactor(address common.Address, transactor bind.ContractTransactor) (*VoterTransactor, error) {
	contract, err := bindVoter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VoterTransactor{contract: contract}, nil
}

// NewVoterFilterer creates a new log filterer instance of Voter, bound to a specific deployed contract.
func NewVoterFilterer(address common.Address, filterer bind.ContractFilterer) (*VoterFilterer, error) {
	contract, err := bindVoter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VoterFilterer{contract: contract}, nil
}

// bindVoter binds a generic wrapper to an already deployed contract.
func bindVoter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VoterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Voter *VoterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Voter.Contract.VoterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Voter *VoterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Voter.Contract.VoterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Voter *VoterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Voter.Contract.VoterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Voter *VoterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Voter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Voter *VoterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Voter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Voter *VoterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Voter.Contract.contract.Transact(opts, method, params...)
}

// EpochVoteEnd is a free data retrieval call binding the contract method 0xecc44a38.
//
// Solidity: function epochVoteEnd() view returns(uint256)
func (_Voter *VoterCaller) EpochVoteEnd(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Voter.contract.Call(opts, &out, "epochVoteEnd")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochVoteEnd is a free data retrieval call binding the contract method 0xecc44a38.
//
// Solidity: function epochVoteEnd() view returns(uint256)
func (_Voter *VoterSession) EpochVoteEnd() (*big.Int, error) {
	return _Voter.Contract.EpochVoteEnd(&_Voter.CallOpts)
}

// EpochVoteEnd is a free data retrieval call binding the contract method 0xecc44a38.
//
// Solidity: function epochVoteEnd() view returns(uint256)
func (_Voter *VoterCallerSession) EpochVoteEnd() (*big.Int, error) {
	return _Voter.Contract.EpochVoteEnd(&_Voter.CallOpts)
}

// GetGaugeMultiplier is a free data retrieval call binding the contract method 0x7eca406f.
//
// Solidity: function getGaugeMultiplier(bytes32 _model) view returns(uint256)
func (_Voter *VoterCaller) GetGaugeMultiplier(opts *bind.CallOpts, _model [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Voter.contract.Call(opts, &out, "getGaugeMultiplier", _model)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetGaugeMultiplier is a free data retrieval call binding the contract method 0x7eca406f.
//
// Solidity: function getGaugeMultiplier(bytes32 _model) view returns(uint256)
func (_Voter *VoterSession) GetGaugeMultiplier(_model [32]byte) (*big.Int, error) {
	return _Voter.Contract.GetGaugeMultiplier(&_Voter.CallOpts, _model)
}

// GetGaugeMultiplier is a free data retrieval call binding the contract method 0x7eca406f.
//
// Solidity: function getGaugeMultiplier(bytes32 _model) view returns(uint256)
func (_Voter *VoterCallerSession) GetGaugeMultiplier(_model [32]byte) (*big.Int, error) {
	return _Voter.Contract.GetGaugeMultiplier(&_Voter.CallOpts, _model)
}

// IsAlive is a free data retrieval call binding the contract method 0x1b8dd060.
//
// Solidity: function isAlive(bytes32 ) view returns(bool)
func (_Voter *VoterCaller) IsAlive(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _Voter.contract.Call(opts, &out, "isAlive", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAlive is a free data retrieval call binding the contract method 0x1b8dd060.
//
// Solidity: function isAlive(bytes32 ) view returns(bool)
func (_Voter *VoterSession) IsAlive(arg0 [32]byte) (bool, error) {
	return _Voter.Contract.IsAlive(&_Voter.CallOpts, arg0)
}

// IsAlive is a free data retrieval call binding the contract method 0x1b8dd060.
//
// Solidity: function isAlive(bytes32 ) view returns(bool)
func (_Voter *VoterCallerSession) IsAlive(arg0 [32]byte) (bool, error) {
	return _Voter.Contract.IsAlive(&_Voter.CallOpts, arg0)
}

// IsGauge is a free data retrieval call binding the contract method 0xb65a78b5.
//
// Solidity: function isGauge(bytes32 ) view returns(bool)
func (_Voter *VoterCaller) IsGauge(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _Voter.contract.Call(opts, &out, "isGauge", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsGauge is a free data retrieval call binding the contract method 0xb65a78b5.
//
// Solidity: function isGauge(bytes32 ) view returns(bool)
func (_Voter *VoterSession) IsGauge(arg0 [32]byte) (bool, error) {
	return _Voter.Contract.IsGauge(&_Voter.CallOpts, arg0)
}

// IsGauge is a free data retrieval call binding the contract method 0xb65a78b5.
//
// Solidity: function isGauge(bytes32 ) view returns(bool)
func (_Voter *VoterCallerSession) IsGauge(arg0 [32]byte) (bool, error) {
	return _Voter.Contract.IsGauge(&_Voter.CallOpts, arg0)
}

// IsWhitelisted is a free data retrieval call binding the contract method 0x01a5e3fe.
//
// Solidity: function isWhitelisted(bytes32 ) view returns(bool)
func (_Voter *VoterCaller) IsWhitelisted(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _Voter.contract.Call(opts, &out, "isWhitelisted", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsWhitelisted is a free data retrieval call binding the contract method 0x01a5e3fe.
//
// Solidity: function isWhitelisted(bytes32 ) view returns(bool)
func (_Voter *VoterSession) IsWhitelisted(arg0 [32]byte) (bool, error) {
	return _Voter.Contract.IsWhitelisted(&_Voter.CallOpts, arg0)
}

// IsWhitelisted is a free data retrieval call binding the contract method 0x01a5e3fe.
//
// Solidity: function isWhitelisted(bytes32 ) view returns(bool)
func (_Voter *VoterCallerSession) IsWhitelisted(arg0 [32]byte) (bool, error) {
	return _Voter.Contract.IsWhitelisted(&_Voter.CallOpts, arg0)
}

// LastVoted is a free data retrieval call binding the contract method 0xf3594be0.
//
// Solidity: function lastVoted(uint256 ) view returns(uint256)
func (_Voter *VoterCaller) LastVoted(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Voter.contract.Call(opts, &out, "lastVoted", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastVoted is a free data retrieval call binding the contract method 0xf3594be0.
//
// Solidity: function lastVoted(uint256 ) view returns(uint256)
func (_Voter *VoterSession) LastVoted(arg0 *big.Int) (*big.Int, error) {
	return _Voter.Contract.LastVoted(&_Voter.CallOpts, arg0)
}

// LastVoted is a free data retrieval call binding the contract method 0xf3594be0.
//
// Solidity: function lastVoted(uint256 ) view returns(uint256)
func (_Voter *VoterCallerSession) LastVoted(arg0 *big.Int) (*big.Int, error) {
	return _Voter.Contract.LastVoted(&_Voter.CallOpts, arg0)
}

// Length is a free data retrieval call binding the contract method 0x1f7b6d32.
//
// Solidity: function length() view returns(uint256)
func (_Voter *VoterCaller) Length(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Voter.contract.Call(opts, &out, "length")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Length is a free data retrieval call binding the contract method 0x1f7b6d32.
//
// Solidity: function length() view returns(uint256)
func (_Voter *VoterSession) Length() (*big.Int, error) {
	return _Voter.Contract.Length(&_Voter.CallOpts)
}

// Length is a free data retrieval call binding the contract method 0x1f7b6d32.
//
// Solidity: function length() view returns(uint256)
func (_Voter *VoterCallerSession) Length() (*big.Int, error) {
	return _Voter.Contract.Length(&_Voter.CallOpts)
}

// ModelVote is a free data retrieval call binding the contract method 0x5a1fdbc0.
//
// Solidity: function modelVote(uint256 , uint256 ) view returns(bytes32)
func (_Voter *VoterCaller) ModelVote(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Voter.contract.Call(opts, &out, "modelVote", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ModelVote is a free data retrieval call binding the contract method 0x5a1fdbc0.
//
// Solidity: function modelVote(uint256 , uint256 ) view returns(bytes32)
func (_Voter *VoterSession) ModelVote(arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	return _Voter.Contract.ModelVote(&_Voter.CallOpts, arg0, arg1)
}

// ModelVote is a free data retrieval call binding the contract method 0x5a1fdbc0.
//
// Solidity: function modelVote(uint256 , uint256 ) view returns(bytes32)
func (_Voter *VoterCallerSession) ModelVote(arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	return _Voter.Contract.ModelVote(&_Voter.CallOpts, arg0, arg1)
}

// Models is a free data retrieval call binding the contract method 0x6a030ca9.
//
// Solidity: function models(uint256 ) view returns(bytes32)
func (_Voter *VoterCaller) Models(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Voter.contract.Call(opts, &out, "models", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Models is a free data retrieval call binding the contract method 0x6a030ca9.
//
// Solidity: function models(uint256 ) view returns(bytes32)
func (_Voter *VoterSession) Models(arg0 *big.Int) ([32]byte, error) {
	return _Voter.Contract.Models(&_Voter.CallOpts, arg0)
}

// Models is a free data retrieval call binding the contract method 0x6a030ca9.
//
// Solidity: function models(uint256 ) view returns(bytes32)
func (_Voter *VoterCallerSession) Models(arg0 *big.Int) ([32]byte, error) {
	return _Voter.Contract.Models(&_Voter.CallOpts, arg0)
}

// TotalWeight is a free data retrieval call binding the contract method 0x96c82e57.
//
// Solidity: function totalWeight() view returns(uint256)
func (_Voter *VoterCaller) TotalWeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Voter.contract.Call(opts, &out, "totalWeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalWeight is a free data retrieval call binding the contract method 0x96c82e57.
//
// Solidity: function totalWeight() view returns(uint256)
func (_Voter *VoterSession) TotalWeight() (*big.Int, error) {
	return _Voter.Contract.TotalWeight(&_Voter.CallOpts)
}

// TotalWeight is a free data retrieval call binding the contract method 0x96c82e57.
//
// Solidity: function totalWeight() view returns(uint256)
func (_Voter *VoterCallerSession) TotalWeight() (*big.Int, error) {
	return _Voter.Contract.TotalWeight(&_Voter.CallOpts)
}

// UsedWeights is a free data retrieval call binding the contract method 0x79e93824.
//
// Solidity: function usedWeights(uint256 ) view returns(uint256)
func (_Voter *VoterCaller) UsedWeights(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Voter.contract.Call(opts, &out, "usedWeights", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UsedWeights is a free data retrieval call binding the contract method 0x79e93824.
//
// Solidity: function usedWeights(uint256 ) view returns(uint256)
func (_Voter *VoterSession) UsedWeights(arg0 *big.Int) (*big.Int, error) {
	return _Voter.Contract.UsedWeights(&_Voter.CallOpts, arg0)
}

// UsedWeights is a free data retrieval call binding the contract method 0x79e93824.
//
// Solidity: function usedWeights(uint256 ) view returns(uint256)
func (_Voter *VoterCallerSession) UsedWeights(arg0 *big.Int) (*big.Int, error) {
	return _Voter.Contract.UsedWeights(&_Voter.CallOpts, arg0)
}

// Votes is a free data retrieval call binding the contract method 0x5b0c8a63.
//
// Solidity: function votes(uint256 , bytes32 ) view returns(uint256)
func (_Voter *VoterCaller) Votes(opts *bind.CallOpts, arg0 *big.Int, arg1 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Voter.contract.Call(opts, &out, "votes", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Votes is a free data retrieval call binding the contract method 0x5b0c8a63.
//
// Solidity: function votes(uint256 , bytes32 ) view returns(uint256)
func (_Voter *VoterSession) Votes(arg0 *big.Int, arg1 [32]byte) (*big.Int, error) {
	return _Voter.Contract.Votes(&_Voter.CallOpts, arg0, arg1)
}

// Votes is a free data retrieval call binding the contract method 0x5b0c8a63.
//
// Solidity: function votes(uint256 , bytes32 ) view returns(uint256)
func (_Voter *VoterCallerSession) Votes(arg0 *big.Int, arg1 [32]byte) (*big.Int, error) {
	return _Voter.Contract.Votes(&_Voter.CallOpts, arg0, arg1)
}

// VotingEscrow is a free data retrieval call binding the contract method 0x4f2bfe5b.
//
// Solidity: function votingEscrow() view returns(address)
func (_Voter *VoterCaller) VotingEscrow(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Voter.contract.Call(opts, &out, "votingEscrow")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VotingEscrow is a free data retrieval call binding the contract method 0x4f2bfe5b.
//
// Solidity: function votingEscrow() view returns(address)
func (_Voter *VoterSession) VotingEscrow() (common.Address, error) {
	return _Voter.Contract.VotingEscrow(&_Voter.CallOpts)
}

// VotingEscrow is a free data retrieval call binding the contract method 0x4f2bfe5b.
//
// Solidity: function votingEscrow() view returns(address)
func (_Voter *VoterCallerSession) VotingEscrow() (common.Address, error) {
	return _Voter.Contract.VotingEscrow(&_Voter.CallOpts)
}

// Weights is a free data retrieval call binding the contract method 0x7addf675.
//
// Solidity: function weights(bytes32 ) view returns(uint256)
func (_Voter *VoterCaller) Weights(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Voter.contract.Call(opts, &out, "weights", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Weights is a free data retrieval call binding the contract method 0x7addf675.
//
// Solidity: function weights(bytes32 ) view returns(uint256)
func (_Voter *VoterSession) Weights(arg0 [32]byte) (*big.Int, error) {
	return _Voter.Contract.Weights(&_Voter.CallOpts, arg0)
}

// Weights is a free data retrieval call binding the contract method 0x7addf675.
//
// Solidity: function weights(bytes32 ) view returns(uint256)
func (_Voter *VoterCallerSession) Weights(arg0 [32]byte) (*big.Int, error) {
	return _Voter.Contract.Weights(&_Voter.CallOpts, arg0)
}

// CreateGauge is a paid mutator transaction binding the contract method 0xf9315157.
//
// Solidity: function createGauge(bytes32 _model) returns()
func (_Voter *VoterTransactor) CreateGauge(opts *bind.TransactOpts, _model [32]byte) (*types.Transaction, error) {
	return _Voter.contract.Transact(opts, "createGauge", _model)
}

// CreateGauge is a paid mutator transaction binding the contract method 0xf9315157.
//
// Solidity: function createGauge(bytes32 _model) returns()
func (_Voter *VoterSession) CreateGauge(_model [32]byte) (*types.Transaction, error) {
	return _Voter.Contract.CreateGauge(&_Voter.TransactOpts, _model)
}

// CreateGauge is a paid mutator transaction binding the contract method 0xf9315157.
//
// Solidity: function createGauge(bytes32 _model) returns()
func (_Voter *VoterTransactorSession) CreateGauge(_model [32]byte) (*types.Transaction, error) {
	return _Voter.Contract.CreateGauge(&_Voter.TransactOpts, _model)
}

// KillGauge is a paid mutator transaction binding the contract method 0x82fe6eae.
//
// Solidity: function killGauge(bytes32 _model) returns()
func (_Voter *VoterTransactor) KillGauge(opts *bind.TransactOpts, _model [32]byte) (*types.Transaction, error) {
	return _Voter.contract.Transact(opts, "killGauge", _model)
}

// KillGauge is a paid mutator transaction binding the contract method 0x82fe6eae.
//
// Solidity: function killGauge(bytes32 _model) returns()
func (_Voter *VoterSession) KillGauge(_model [32]byte) (*types.Transaction, error) {
	return _Voter.Contract.KillGauge(&_Voter.TransactOpts, _model)
}

// KillGauge is a paid mutator transaction binding the contract method 0x82fe6eae.
//
// Solidity: function killGauge(bytes32 _model) returns()
func (_Voter *VoterTransactorSession) KillGauge(_model [32]byte) (*types.Transaction, error) {
	return _Voter.Contract.KillGauge(&_Voter.TransactOpts, _model)
}

// Poke is a paid mutator transaction binding the contract method 0x32145f90.
//
// Solidity: function poke(uint256 _tokenId) returns()
func (_Voter *VoterTransactor) Poke(opts *bind.TransactOpts, _tokenId *big.Int) (*types.Transaction, error) {
	return _Voter.contract.Transact(opts, "poke", _tokenId)
}

// Poke is a paid mutator transaction binding the contract method 0x32145f90.
//
// Solidity: function poke(uint256 _tokenId) returns()
func (_Voter *VoterSession) Poke(_tokenId *big.Int) (*types.Transaction, error) {
	return _Voter.Contract.Poke(&_Voter.TransactOpts, _tokenId)
}

// Poke is a paid mutator transaction binding the contract method 0x32145f90.
//
// Solidity: function poke(uint256 _tokenId) returns()
func (_Voter *VoterTransactorSession) Poke(_tokenId *big.Int) (*types.Transaction, error) {
	return _Voter.Contract.Poke(&_Voter.TransactOpts, _tokenId)
}

// Reset is a paid mutator transaction binding the contract method 0x310bd74b.
//
// Solidity: function reset(uint256 _tokenId) returns()
func (_Voter *VoterTransactor) Reset(opts *bind.TransactOpts, _tokenId *big.Int) (*types.Transaction, error) {
	return _Voter.contract.Transact(opts, "reset", _tokenId)
}

// Reset is a paid mutator transaction binding the contract method 0x310bd74b.
//
// Solidity: function reset(uint256 _tokenId) returns()
func (_Voter *VoterSession) Reset(_tokenId *big.Int) (*types.Transaction, error) {
	return _Voter.Contract.Reset(&_Voter.TransactOpts, _tokenId)
}

// Reset is a paid mutator transaction binding the contract method 0x310bd74b.
//
// Solidity: function reset(uint256 _tokenId) returns()
func (_Voter *VoterTransactorSession) Reset(_tokenId *big.Int) (*types.Transaction, error) {
	return _Voter.Contract.Reset(&_Voter.TransactOpts, _tokenId)
}

// ReviveGauge is a paid mutator transaction binding the contract method 0x591e5582.
//
// Solidity: function reviveGauge(bytes32 _model) returns()
func (_Voter *VoterTransactor) ReviveGauge(opts *bind.TransactOpts, _model [32]byte) (*types.Transaction, error) {
	return _Voter.contract.Transact(opts, "reviveGauge", _model)
}

// ReviveGauge is a paid mutator transaction binding the contract method 0x591e5582.
//
// Solidity: function reviveGauge(bytes32 _model) returns()
func (_Voter *VoterSession) ReviveGauge(_model [32]byte) (*types.Transaction, error) {
	return _Voter.Contract.ReviveGauge(&_Voter.TransactOpts, _model)
}

// ReviveGauge is a paid mutator transaction binding the contract method 0x591e5582.
//
// Solidity: function reviveGauge(bytes32 _model) returns()
func (_Voter *VoterTransactorSession) ReviveGauge(_model [32]byte) (*types.Transaction, error) {
	return _Voter.Contract.ReviveGauge(&_Voter.TransactOpts, _model)
}

// Vote is a paid mutator transaction binding the contract method 0x0c295755.
//
// Solidity: function vote(uint256 tokenId, bytes32[] _modelVote, uint256[] _weights) returns()
func (_Voter *VoterTransactor) Vote(opts *bind.TransactOpts, tokenId *big.Int, _modelVote [][32]byte, _weights []*big.Int) (*types.Transaction, error) {
	return _Voter.contract.Transact(opts, "vote", tokenId, _modelVote, _weights)
}

// Vote is a paid mutator transaction binding the contract method 0x0c295755.
//
// Solidity: function vote(uint256 tokenId, bytes32[] _modelVote, uint256[] _weights) returns()
func (_Voter *VoterSession) Vote(tokenId *big.Int, _modelVote [][32]byte, _weights []*big.Int) (*types.Transaction, error) {
	return _Voter.Contract.Vote(&_Voter.TransactOpts, tokenId, _modelVote, _weights)
}

// Vote is a paid mutator transaction binding the contract method 0x0c295755.
//
// Solidity: function vote(uint256 tokenId, bytes32[] _modelVote, uint256[] _weights) returns()
func (_Voter *VoterTransactorSession) Vote(tokenId *big.Int, _modelVote [][32]byte, _weights []*big.Int) (*types.Transaction, error) {
	return _Voter.Contract.Vote(&_Voter.TransactOpts, tokenId, _modelVote, _weights)
}

// Whitelist is a paid mutator transaction binding the contract method 0xafb40c8e.
//
// Solidity: function whitelist(bytes32 _model) returns()
func (_Voter *VoterTransactor) Whitelist(opts *bind.TransactOpts, _model [32]byte) (*types.Transaction, error) {
	return _Voter.contract.Transact(opts, "whitelist", _model)
}

// Whitelist is a paid mutator transaction binding the contract method 0xafb40c8e.
//
// Solidity: function whitelist(bytes32 _model) returns()
func (_Voter *VoterSession) Whitelist(_model [32]byte) (*types.Transaction, error) {
	return _Voter.Contract.Whitelist(&_Voter.TransactOpts, _model)
}

// Whitelist is a paid mutator transaction binding the contract method 0xafb40c8e.
//
// Solidity: function whitelist(bytes32 _model) returns()
func (_Voter *VoterTransactorSession) Whitelist(_model [32]byte) (*types.Transaction, error) {
	return _Voter.Contract.Whitelist(&_Voter.TransactOpts, _model)
}
