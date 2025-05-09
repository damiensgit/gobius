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

// BulkTasksTaskSolutionWithContestationInfo is an auto generated low-level Go binding around an user-defined struct.
type BulkTasksTaskSolutionWithContestationInfo struct {
	TaskId                       [32]byte
	SolutionValidator            common.Address
	SolutionBlocktime            uint64
	SolutionClaimed              bool
	SolutionExists               bool
	ContestationValidator        common.Address
	ContestationBlocktime        uint64
	ContestationFinishStartIndex uint32
	ContestationSlashAmount      *big.Int
	ContestationExists           bool
}

// IArbiusEngineEngineContestation is an auto generated low-level Go binding around an user-defined struct.
type IArbiusEngineEngineContestation struct {
	Validator        common.Address
	Blocktime        uint64
	FinishStartIndex uint32
	SlashAmount      *big.Int
}

// IArbiusEngineEngineSolution is an auto generated low-level Go binding around an user-defined struct.
type IArbiusEngineEngineSolution struct {
	Validator common.Address
	Blocktime uint64
	Claimed   bool
}

// BulkTasksMetaData contains all meta data concerning the BulkTasks contract.
var BulkTasksMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"basetoken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"engine_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"basetoken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"bulkSignalCommitment\",\"inputs\":[{\"name\":\"commitments_\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimSolutions\",\"inputs\":[{\"name\":\"_taskids\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"engine\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIArbiusEngine\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBulkCombinedTaskInfo\",\"inputs\":[{\"name\":\"taskIds\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structBulkTasks.TaskSolutionWithContestationInfo[]\",\"components\":[{\"name\":\"taskId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"solutionValidator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"solutionBlocktime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"solutionClaimed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"solutionExists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"contestationValidator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"contestationBlocktime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"contestationFinishStartIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"contestationSlashAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"contestationExists\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBulkContestations\",\"inputs\":[{\"name\":\"taskids_\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structIArbiusEngine.EngineContestation[]\",\"components\":[{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"blocktime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"finish_start_index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"slashAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitments\",\"inputs\":[{\"name\":\"commitments_\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSolutions\",\"inputs\":[{\"name\":\"taskids_\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structIArbiusEngine.EngineSolution[]\",\"components\":[{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"blocktime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"claimed\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516109c03803806109c083398101604081905261002f916100ec565b600080546001600160a01b038481166001600160a01b031992831681179093556001805491851691909216811790915560405163095ea7b360e01b81526004810191909152600019602482015263095ea7b3906044016020604051808303816000875af11580156100a4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100c8919061011f565b505050610148565b80516001600160a01b03811681146100e757600080fd5b919050565b600080604083850312156100ff57600080fd5b610108836100d0565b9150610116602084016100d0565b90509250929050565b60006020828403121561013157600080fd5b8151801515811461014157600080fd5b9392505050565b610869806101576000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806302546ffc146100675780631023ad7c146100975780631fbeeddb146100ac57806350ec55e1146100cc57806380452ec5146100df578063c9d4623f146100ff575b600080fd5b60005461007a906001600160a01b031681565b6040516001600160a01b0390911681526020015b60405180910390f35b6100aa6100a53660046104a4565b610112565b005b6100bf6100ba3660046104a4565b6101b1565b60405161008e919061053d565b6100aa6100da3660046104a4565b6102f4565b6100f26100ed3660046104a4565b61038d565b60405161008e91906105f7565b60015461007a906001600160a01b031681565b8060005b818110156101ab576001546001600160a01b03166377286d178585848181106101415761014161063b565b905060200201356040518263ffffffff1660e01b815260040161016691815260200190565b600060405180830381600087803b15801561018057600080fd5b505af1158015610194573d6000803e3d6000fd5b5050505080806101a390610651565b915050610116565b50505050565b60608160008167ffffffffffffffff8111156101cf576101cf610678565b60405190808252806020026020018201604052801561022057816020015b60408051608081018252600080825260208083018290529282015260608082015282526000199092019101816101ed5790505b50905060005b828110156102eb576001546001600160a01b03166375c705098787848181106102515761025161063b565b905060200201356040518263ffffffff1660e01b815260040161027691815260200190565b600060405180830381865afa158015610293573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526102bb919081019061075c565b8282815181106102cd576102cd61063b565b602002602001018190525080806102e390610651565b915050610226565b50949350505050565b8060005b818110156101ab576001546001600160a01b031663506ea7de8585848181106103235761032361063b565b905060200201356040518263ffffffff1660e01b815260040161034891815260200190565b600060405180830381600087803b15801561036257600080fd5b505af1158015610376573d6000803e3d6000fd5b50505050808061038590610651565b9150506102f8565b60608160008167ffffffffffffffff8111156103ab576103ab610678565b6040519080825280602002602001820160405280156103d4578160200160208202803683370190505b50905060005b828110156102eb576001546001600160a01b031663839df9458787848181106104055761040561063b565b905060200201356040518263ffffffff1660e01b815260040161042a91815260200190565b602060405180830381865afa158015610447573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061046b9190610811565b67ffffffffffffffff168282815181106104875761048761063b565b60209081029190910101528061049c81610651565b9150506103da565b600080602083850312156104b757600080fd5b823567ffffffffffffffff808211156104cf57600080fd5b818501915085601f8301126104e357600080fd5b8135818111156104f257600080fd5b8660208260051b850101111561050757600080fd5b60209290920196919550909350505050565b60005b8381101561053457818101518382015260200161051c565b50506000910152565b60006020808301818452808551808352604092508286019150828160051b87010184880160005b838110156105e957888303603f19018552815180516001600160a01b031684528781015167ffffffffffffffff168885015286810151151587850152606090810151608091850182905280519185018290529060a0906105c981838801858d01610519565b96890196601f01601f191694909401909301925090860190600101610564565b509098975050505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561062f57835183529284019291840191600101610613565b50909695505050505050565b634e487b7160e01b600052603260045260246000fd5b60006001820161067157634e487b7160e01b600052601160045260246000fd5b5060010190565b634e487b7160e01b600052604160045260246000fd5b6040516080810167ffffffffffffffff811182821017156106b1576106b1610678565b60405290565b805167ffffffffffffffff811681146106cf57600080fd5b919050565b600082601f8301126106e557600080fd5b815167ffffffffffffffff8082111561070057610700610678565b604051601f8301601f19908116603f0116810190828211818310171561072857610728610678565b8160405283815286602085880101111561074157600080fd5b610752846020830160208901610519565b9695505050505050565b60006020828403121561076e57600080fd5b815167ffffffffffffffff8082111561078657600080fd5b908301906080828603121561079a57600080fd5b6107a261068e565b82516001600160a01b03811681146107b957600080fd5b81526107c7602084016106b7565b6020820152604083015180151581146107df57600080fd5b60408201526060830151828111156107f657600080fd5b610802878286016106d4565b60608301525095945050505050565b60006020828403121561082357600080fd5b61082c826106b7565b939250505056fea26469706673582212200f0a9f0f31d96731e2451ccf8b52e61cb29cdcec3e6292d52f15f8644089ee2864736f6c63430008130033",
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

// GetBulkCombinedTaskInfo is a free data retrieval call binding the contract method 0x3dd871cc.
//
// Solidity: function getBulkCombinedTaskInfo(bytes32[] taskIds) view returns((bytes32,address,uint64,bool,bool,address,uint64,uint32,uint256,bool)[])
func (_BulkTasks *BulkTasksCaller) GetBulkCombinedTaskInfo(opts *bind.CallOpts, taskIds [][32]byte) ([]BulkTasksTaskSolutionWithContestationInfo, error) {
	var out []interface{}
	err := _BulkTasks.contract.Call(opts, &out, "getBulkCombinedTaskInfo", taskIds)

	if err != nil {
		return *new([]BulkTasksTaskSolutionWithContestationInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]BulkTasksTaskSolutionWithContestationInfo)).(*[]BulkTasksTaskSolutionWithContestationInfo)

	return out0, err

}

// GetBulkCombinedTaskInfo is a free data retrieval call binding the contract method 0x3dd871cc.
//
// Solidity: function getBulkCombinedTaskInfo(bytes32[] taskIds) view returns((bytes32,address,uint64,bool,bool,address,uint64,uint32,uint256,bool)[])
func (_BulkTasks *BulkTasksSession) GetBulkCombinedTaskInfo(taskIds [][32]byte) ([]BulkTasksTaskSolutionWithContestationInfo, error) {
	return _BulkTasks.Contract.GetBulkCombinedTaskInfo(&_BulkTasks.CallOpts, taskIds)
}

// GetBulkCombinedTaskInfo is a free data retrieval call binding the contract method 0x3dd871cc.
//
// Solidity: function getBulkCombinedTaskInfo(bytes32[] taskIds) view returns((bytes32,address,uint64,bool,bool,address,uint64,uint32,uint256,bool)[])
func (_BulkTasks *BulkTasksCallerSession) GetBulkCombinedTaskInfo(taskIds [][32]byte) ([]BulkTasksTaskSolutionWithContestationInfo, error) {
	return _BulkTasks.Contract.GetBulkCombinedTaskInfo(&_BulkTasks.CallOpts, taskIds)
}

// GetBulkContestations is a free data retrieval call binding the contract method 0x900a8b61.
//
// Solidity: function getBulkContestations(bytes32[] taskids_) view returns((address,uint64,uint32,uint256)[])
func (_BulkTasks *BulkTasksCaller) GetBulkContestations(opts *bind.CallOpts, taskids_ [][32]byte) ([]IArbiusEngineEngineContestation, error) {
	var out []interface{}
	err := _BulkTasks.contract.Call(opts, &out, "getBulkContestations", taskids_)

	if err != nil {
		return *new([]IArbiusEngineEngineContestation), err
	}

	out0 := *abi.ConvertType(out[0], new([]IArbiusEngineEngineContestation)).(*[]IArbiusEngineEngineContestation)

	return out0, err

}

// GetBulkContestations is a free data retrieval call binding the contract method 0x900a8b61.
//
// Solidity: function getBulkContestations(bytes32[] taskids_) view returns((address,uint64,uint32,uint256)[])
func (_BulkTasks *BulkTasksSession) GetBulkContestations(taskids_ [][32]byte) ([]IArbiusEngineEngineContestation, error) {
	return _BulkTasks.Contract.GetBulkContestations(&_BulkTasks.CallOpts, taskids_)
}

// GetBulkContestations is a free data retrieval call binding the contract method 0x900a8b61.
//
// Solidity: function getBulkContestations(bytes32[] taskids_) view returns((address,uint64,uint32,uint256)[])
func (_BulkTasks *BulkTasksCallerSession) GetBulkContestations(taskids_ [][32]byte) ([]IArbiusEngineEngineContestation, error) {
	return _BulkTasks.Contract.GetBulkContestations(&_BulkTasks.CallOpts, taskids_)
}

// GetCommitments is a free data retrieval call binding the contract method 0x80452ec5.
//
// Solidity: function getCommitments(bytes32[] commitments_) view returns(uint256[])
func (_BulkTasks *BulkTasksCaller) GetCommitments(opts *bind.CallOpts, commitments_ [][32]byte) ([]*big.Int, error) {
	var out []interface{}
	err := _BulkTasks.contract.Call(opts, &out, "getCommitments", commitments_)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetCommitments is a free data retrieval call binding the contract method 0x80452ec5.
//
// Solidity: function getCommitments(bytes32[] commitments_) view returns(uint256[])
func (_BulkTasks *BulkTasksSession) GetCommitments(commitments_ [][32]byte) ([]*big.Int, error) {
	return _BulkTasks.Contract.GetCommitments(&_BulkTasks.CallOpts, commitments_)
}

// GetCommitments is a free data retrieval call binding the contract method 0x80452ec5.
//
// Solidity: function getCommitments(bytes32[] commitments_) view returns(uint256[])
func (_BulkTasks *BulkTasksCallerSession) GetCommitments(commitments_ [][32]byte) ([]*big.Int, error) {
	return _BulkTasks.Contract.GetCommitments(&_BulkTasks.CallOpts, commitments_)
}

// GetSolutions is a free data retrieval call binding the contract method 0x1fbeeddb.
//
// Solidity: function getSolutions(bytes32[] taskids_) view returns((address,uint64,bool)[])
func (_BulkTasks *BulkTasksCaller) GetSolutions(opts *bind.CallOpts, taskids_ [][32]byte) ([]IArbiusEngineEngineSolution, error) {
	var out []interface{}
	err := _BulkTasks.contract.Call(opts, &out, "getSolutions", taskids_)

	if err != nil {
		return *new([]IArbiusEngineEngineSolution), err
	}

	out0 := *abi.ConvertType(out[0], new([]IArbiusEngineEngineSolution)).(*[]IArbiusEngineEngineSolution)

	return out0, err

}

// GetSolutions is a free data retrieval call binding the contract method 0x1fbeeddb.
//
// Solidity: function getSolutions(bytes32[] taskids_) view returns((address,uint64,bool)[])
func (_BulkTasks *BulkTasksSession) GetSolutions(taskids_ [][32]byte) ([]IArbiusEngineEngineSolution, error) {
	return _BulkTasks.Contract.GetSolutions(&_BulkTasks.CallOpts, taskids_)
}

// GetSolutions is a free data retrieval call binding the contract method 0x1fbeeddb.
//
// Solidity: function getSolutions(bytes32[] taskids_) view returns((address,uint64,bool)[])
func (_BulkTasks *BulkTasksCallerSession) GetSolutions(taskids_ [][32]byte) ([]IArbiusEngineEngineSolution, error) {
	return _BulkTasks.Contract.GetSolutions(&_BulkTasks.CallOpts, taskids_)
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
