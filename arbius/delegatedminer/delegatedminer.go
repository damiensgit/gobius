// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package delegatedvalidator

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

// DelegatedvalidatorMetaData contains all meta data concerning the Delegatedvalidator contract.
var DelegatedvalidatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIArbius\",\"name\":\"_arbius\",\"type\":\"address\"},{\"internalType\":\"contractIBaseToken\",\"name\":\"_baseToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_miner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"commitments_\",\"type\":\"bytes32[]\"}],\"name\":\"bulkSignalCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_taskids\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes[]\",\"name\":\"_cids\",\"type\":\"bytes[]\"}],\"name\":\"bulkSubmitSolution\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"call\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"count_\",\"type\":\"uint256\"}],\"name\":\"cancelValidatorWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_taskids\",\"type\":\"bytes32[]\"}],\"name\":\"claimSolutions\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegatee\",\"type\":\"address\"}],\"name\":\"delegateVoting\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"}],\"name\":\"initiateValidatorWithdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"safeWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIArbius\",\"name\":\"_newAribius\",\"type\":\"address\"}],\"name\":\"setArbius\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newMiner\",\"type\":\"address\"}],\"name\":\"setMiner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitment_\",\"type\":\"bytes32\"}],\"name\":\"signalCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_taskid\",\"type\":\"bytes32\"}],\"name\":\"submitContestation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_tasks\",\"type\":\"uint8\"}],\"name\":\"submitMultipleTasks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_taskid\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_cid\",\"type\":\"bytes\"}],\"name\":\"submitSolution\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_version\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_model\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_input\",\"type\":\"bytes\"}],\"name\":\"submitTask\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"validatorDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"count_\",\"type\":\"uint256\"}],\"name\":\"validatorWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_taskid\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"_yea\",\"type\":\"bool\"}],\"name\":\"voteOnContestation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawETH\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162002b8338038062002b838339818101604052810190620000379190620002de565b82600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555081600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663095ea7b3600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6040518363ffffffff1660e01b81526004016200019a92919062000366565b6020604051808303816000875af1158015620001ba573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001e09190620003d0565b5050505062000402565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200021c82620001ef565b9050919050565b600062000230826200020f565b9050919050565b620002428162000223565b81146200024e57600080fd5b50565b600081519050620002628162000237565b92915050565b600062000275826200020f565b9050919050565b620002878162000268565b81146200029357600080fd5b50565b600081519050620002a7816200027c565b92915050565b620002b8816200020f565b8114620002c457600080fd5b50565b600081519050620002d881620002ad565b92915050565b600080600060608486031215620002fa57620002f9620001ea565b5b60006200030a8682870162000251565b93505060206200031d8682870162000296565b92505060406200033086828701620002c7565b9150509250925092565b62000345816200020f565b82525050565b6000819050919050565b62000360816200034b565b82525050565b60006040820190506200037d60008301856200033a565b6200038c602083018462000355565b9392505050565b60008115159050919050565b620003aa8162000393565b8114620003b657600080fd5b50565b600081519050620003ca816200039f565b92915050565b600060208284031215620003e957620003e8620001ea565b5b6000620003f984828501620003b9565b91505092915050565b61277180620004126000396000f3fe6080604052600436106101145760003560e01c806365d445fb116100a0578063b0fd035b11610064578063b0fd035b14610372578063cbd2422d1461039b578063dc0143cc146103c4578063eaf7c9dc146103ed578063f14210a61461041657610114565b806365d445fb1461029e578063671f8152146102c75780636dbf2fa0146102f0578063765cb83a146103205780639742ca461461034957610114565b80631825c20e116100e75780631825c20e146101d15780634607315c146101fa578063506ea7de1461022357806350ec55e11461024c57806356914caf1461027557610114565b8063021c79621461011957806308745dd1146101425780630a9857371461016b5780631023ad7c146101a8575b600080fd5b34801561012557600080fd5b50610140600480360381019061013b9190611ac8565b61043f565b005b34801561014e57600080fd5b5061016960048036038101906101649190611c27565b610574565b005b34801561017757600080fd5b50610192600480360381019061018d9190611ac8565b6106b3565b60405161019f9190611cd0565b60405180910390f35b3480156101b457600080fd5b506101cf60048036038101906101ca9190611d41565b6107e9565b005b3480156101dd57600080fd5b506101f860048036038101906101f39190611dc6565b610987565b005b34801561020657600080fd5b50610221600480360381019061021c9190611e44565b610aa8565b005b34801561022f57600080fd5b5061024a60048036038101906102459190611e71565b610b7a565b005b34801561025857600080fd5b50610273600480360381019061026e9190611d41565b610c98565b005b34801561028157600080fd5b5061029c60048036038101906102979190611fdf565b610d06565b005b3480156102aa57600080fd5b506102c560048036038101906102c09190612091565b610e27565b005b3480156102d357600080fd5b506102ee60048036038101906102e99190611e71565b610edb565b005b61030a60048036038101906103059190612150565b610ff9565b6040516103179190612243565b60405180910390f35b34801561032c57600080fd5b5061034760048036038101906103429190612265565b61114b565b005b34801561035557600080fd5b50610370600480360381019061036b9190612265565b611269565b005b34801561037e57600080fd5b5061039960048036038101906103949190611ac8565b61133a565b005b3480156103a757600080fd5b506103c260048036038101906103bd9190611ac8565b6115b4565b005b3480156103d057600080fd5b506103eb60048036038101906103e69190612292565b6116d2565b005b3480156103f957600080fd5b50610414600480360381019061040f9190611ac8565b611796565b005b34801561042257600080fd5b5061043d60048036038101906104389190611ac8565b6118d6565b005b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd3330846040518463ffffffff1660e01b815260040161049e939291906122ce565b6020604051808303816000875af11580156104bd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104e1919061231a565b50600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166393a090ec30836040518363ffffffff1660e01b815260040161053f929190612347565b600060405180830381600087803b15801561055957600080fd5b505af115801561056d573d6000803e3d6000fd5b5050505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610602576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105f9906123cd565b60405180910390fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166308745dd18787878787876040518763ffffffff1660e01b815260040161066796959493929190612438565b6020604051808303816000875af1158015610686573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106aa91906124a9565b50505050505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610744576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161073b906123cd565b60405180910390fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630a985737836040518263ffffffff1660e01b815260040161079f9190611cd0565b6020604051808303816000875af11580156107be573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107e291906124eb565b9050919050565b6040517f77286d1700000000000000000000000000000000000000000000000000000000815236600101835b81811015610850576020816004850137600080602485600073dc64a140aa3e981100a9beca4e685f962f0cf6c95af150602081019050610815565b5050506000739fe46736679d2d9a65f0992f2272de9f3c7fa6e073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016108a29190612518565b602060405180830381865afa1580156108bf573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108e391906124eb565b9050600081111561098257739fe46736679d2d9a65f0992f2272de9f3c7fa6e073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb30836040518363ffffffff1660e01b815260040161093d929190612347565b6020604051808303816000875af115801561095c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610980919061231a565b505b505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610a15576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a0c906123cd565b60405180910390fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16631825c20e83836040518363ffffffff1660e01b8152600401610a72929190612542565b600060405180830381600087803b158015610a8c57600080fd5b505af1158015610aa0573d6000803e3d6000fd5b505050505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610b36576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b2d906123cd565b60405180910390fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610c08576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bff906123cd565b60405180910390fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663506ea7de826040518263ffffffff1660e01b8152600401610c63919061256b565b600060405180830381600087803b158015610c7d57600080fd5b505af1158015610c91573d6000803e3d6000fd5b5050505050565b6040517f506ea7de00000000000000000000000000000000000000000000000000000000815236600101835b81811015610cff576020816004850137600080602485600073dc64a140aa3e981100a9beca4e685f962f0cf6c95af150602081019050610cc4565b5050505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610d94576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d8b906123cd565b60405180910390fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166356914caf83836040518363ffffffff1660e01b8152600401610df1929190612586565b600060405180830381600087803b158015610e0b57600080fd5b505af1158015610e1f573d6000803e3d6000fd5b505050505050565b82816020810284016040517f56914caf000000000000000000000000000000000000000000000000000000008152600060048201526040602482015260226044820152602084028801885b81811015610ecf5760208160048501376020840193506020846064850137602084019350602084608485013760208401935060008060a485600073dc64a140aa3e981100a9beca4e685f962f0cf6c95af150602081019050610e72565b50505050505050505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610f69576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f60906123cd565b60405180910390fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663671f8152826040518263ffffffff1660e01b8152600401610fc4919061256b565b600060405180830381600087803b158015610fde57600080fd5b505af1158015610ff2573d6000803e3d6000fd5b5050505050565b606060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611089576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611080906123cd565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff16036110c257600080fd5b6000808673ffffffffffffffffffffffffffffffffffffffff168686866040516110ed9291906125e6565b60006040518083038185875af1925050503d806000811461112a576040519150601f19603f3d011682016040523d82523d6000602084013e61112f565b606091505b50915091508161113e57600080fd5b8092505050949350505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146111d9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016111d0906123cd565b60405180910390fd5b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16635c19a95c826040518263ffffffff1660e01b81526004016112349190612518565b600060405180830381600087803b15801561124e57600080fd5b505af1158015611262573d6000803e3d6000fd5b5050505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146112f7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112ee906123cd565b60405180910390fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146113c8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016113bf906123cd565b60405180910390fd5b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016114259190612518565b602060405180830381865afa158015611442573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061146691906124eb565b9050818110156114ab576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016114a29061264b565b60405180910390fd5b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff16856040518363ffffffff1660e01b815260040161152a929190612347565b6020604051808303816000875af1158015611549573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061156d919061231a565b9050806115af576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016115a6906126b7565b60405180910390fd5b505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611642576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611639906123cd565b60405180910390fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663cbd2422d826040518263ffffffff1660e01b815260040161169d9190611cd0565b600060405180830381600087803b1580156116b757600080fd5b505af11580156116cb573d6000803e3d6000fd5b5050505050565b6040517f08745dd10000000000000000000000000000000000000000000000000000000081526000600482015273f39fd6e51aad88f6f4ce6ab8827279cfffb9226660248201527f2f77e35c9918358c8ac8b1404b7e9b62cb25d971b1041a531b24137f64dd967a60448201526000606482015260a06084820152601060a4820152600060c482015260005b828110156117915760008060e884600073dc64a140aa3e981100a9beca4e685f962f0cf6c95af15060018101905061175e565b505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611824576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161181b906123cd565b60405180910390fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663763253bb8260008054906101000a900473ffffffffffffffffffffffffffffffffffffffff166040518363ffffffff1660e01b81526004016118a19291906126d7565b600060405180830381600087803b1580156118bb57600080fd5b505af11580156118cf573d6000803e3d6000fd5b5050505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611964576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161195b906123cd565b60405180910390fd5b6000479050818110156119ac576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016119a39061264b565b60405180910390fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16836040516119f390612726565b60006040518083038185875af1925050503d8060008114611a30576040519150601f19603f3d011682016040523d82523d6000602084013e611a35565b606091505b5050905080611a79576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611a70906126b7565b60405180910390fd5b505050565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b611aa581611a92565b8114611ab057600080fd5b50565b600081359050611ac281611a9c565b92915050565b600060208284031215611ade57611add611a88565b5b6000611aec84828501611ab3565b91505092915050565b600060ff82169050919050565b611b0b81611af5565b8114611b1657600080fd5b50565b600081359050611b2881611b02565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000611b5982611b2e565b9050919050565b611b6981611b4e565b8114611b7457600080fd5b50565b600081359050611b8681611b60565b92915050565b6000819050919050565b611b9f81611b8c565b8114611baa57600080fd5b50565b600081359050611bbc81611b96565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f840112611be757611be6611bc2565b5b8235905067ffffffffffffffff811115611c0457611c03611bc7565b5b602083019150836001820283011115611c2057611c1f611bcc565b5b9250929050565b60008060008060008060a08789031215611c4457611c43611a88565b5b6000611c5289828a01611b19565b9650506020611c6389828a01611b77565b9550506040611c7489828a01611bad565b9450506060611c8589828a01611ab3565b935050608087013567ffffffffffffffff811115611ca657611ca5611a8d565b5b611cb289828a01611bd1565b92509250509295509295509295565b611cca81611a92565b82525050565b6000602082019050611ce56000830184611cc1565b92915050565b60008083601f840112611d0157611d00611bc2565b5b8235905067ffffffffffffffff811115611d1e57611d1d611bc7565b5b602083019150836020820283011115611d3a57611d39611bcc565b5b9250929050565b60008060208385031215611d5857611d57611a88565b5b600083013567ffffffffffffffff811115611d7657611d75611a8d565b5b611d8285828601611ceb565b92509250509250929050565b60008115159050919050565b611da381611d8e565b8114611dae57600080fd5b50565b600081359050611dc081611d9a565b92915050565b60008060408385031215611ddd57611ddc611a88565b5b6000611deb85828601611bad565b9250506020611dfc85828601611db1565b9150509250929050565b6000611e1182611b4e565b9050919050565b611e2181611e06565b8114611e2c57600080fd5b50565b600081359050611e3e81611e18565b92915050565b600060208284031215611e5a57611e59611a88565b5b6000611e6884828501611e2f565b91505092915050565b600060208284031215611e8757611e86611a88565b5b6000611e9584828501611bad565b91505092915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b611eec82611ea3565b810181811067ffffffffffffffff82111715611f0b57611f0a611eb4565b5b80604052505050565b6000611f1e611a7e565b9050611f2a8282611ee3565b919050565b600067ffffffffffffffff821115611f4a57611f49611eb4565b5b611f5382611ea3565b9050602081019050919050565b82818337600083830152505050565b6000611f82611f7d84611f2f565b611f14565b905082815260208101848484011115611f9e57611f9d611e9e565b5b611fa9848285611f60565b509392505050565b600082601f830112611fc657611fc5611bc2565b5b8135611fd6848260208601611f6f565b91505092915050565b60008060408385031215611ff657611ff5611a88565b5b600061200485828601611bad565b925050602083013567ffffffffffffffff81111561202557612024611a8d565b5b61203185828601611fb1565b9150509250929050565b60008083601f84011261205157612050611bc2565b5b8235905067ffffffffffffffff81111561206e5761206d611bc7565b5b60208301915083602082028301111561208a57612089611bcc565b5b9250929050565b600080600080604085870312156120ab576120aa611a88565b5b600085013567ffffffffffffffff8111156120c9576120c8611a8d565b5b6120d587828801611ceb565b9450945050602085013567ffffffffffffffff8111156120f8576120f7611a8d565b5b6121048782880161203b565b925092505092959194509250565b600061211d82611b2e565b9050919050565b61212d81612112565b811461213857600080fd5b50565b60008135905061214a81612124565b92915050565b6000806000806060858703121561216a57612169611a88565b5b60006121788782880161213b565b945050602061218987828801611ab3565b935050604085013567ffffffffffffffff8111156121aa576121a9611a8d565b5b6121b687828801611bd1565b925092505092959194509250565b600081519050919050565b600082825260208201905092915050565b60005b838110156121fe5780820151818401526020810190506121e3565b60008484015250505050565b6000612215826121c4565b61221f81856121cf565b935061222f8185602086016121e0565b61223881611ea3565b840191505092915050565b6000602082019050818103600083015261225d818461220a565b905092915050565b60006020828403121561227b5761227a611a88565b5b600061228984828501611b77565b91505092915050565b6000602082840312156122a8576122a7611a88565b5b60006122b684828501611b19565b91505092915050565b6122c881611b4e565b82525050565b60006060820190506122e360008301866122bf565b6122f060208301856122bf565b6122fd6040830184611cc1565b949350505050565b60008151905061231481611d9a565b92915050565b6000602082840312156123305761232f611a88565b5b600061233e84828501612305565b91505092915050565b600060408201905061235c60008301856122bf565b6123696020830184611cc1565b9392505050565b600082825260208201905092915050565b7f6f6e6c794d696e65720000000000000000000000000000000000000000000000600082015250565b60006123b7600983612370565b91506123c282612381565b602082019050919050565b600060208201905081810360008301526123e6816123aa565b9050919050565b6123f681611af5565b82525050565b61240581611b8c565b82525050565b600061241783856121cf565b9350612424838584611f60565b61242d83611ea3565b840190509392505050565b600060a08201905061244d60008301896123ed565b61245a60208301886122bf565b61246760408301876123fc565b6124746060830186611cc1565b818103608083015261248781848661240b565b9050979650505050505050565b6000815190506124a381611b96565b92915050565b6000602082840312156124bf576124be611a88565b5b60006124cd84828501612494565b91505092915050565b6000815190506124e581611a9c565b92915050565b60006020828403121561250157612500611a88565b5b600061250f848285016124d6565b91505092915050565b600060208201905061252d60008301846122bf565b92915050565b61253c81611d8e565b82525050565b600060408201905061255760008301856123fc565b6125646020830184612533565b9392505050565b600060208201905061258060008301846123fc565b92915050565b600060408201905061259b60008301856123fc565b81810360208301526125ad818461220a565b90509392505050565b600081905092915050565b60006125cd83856125b6565b93506125da838584611f60565b82840190509392505050565b60006125f38284866125c1565b91508190509392505050565b7f496e73756666696369656e742062616c616e6365000000000000000000000000600082015250565b6000612635601483612370565b9150612640826125ff565b602082019050919050565b6000602082019050818103600083015261266481612628565b9050919050565b7f5472616e73666572206661696c65640000000000000000000000000000000000600082015250565b60006126a1600f83612370565b91506126ac8261266b565b602082019050919050565b600060208201905081810360008301526126d081612694565b9050919050565b60006040820190506126ec6000830185611cc1565b6126f960208301846122bf565b9392505050565b50565b60006127106000836125b6565b915061271b82612700565b600082019050919050565b600061273182612703565b915081905091905056fea26469706673582212204c3bd57410c13e8fea61e9d8de128045fc88f793b3238eca07c7a9fa6eda81da64736f6c63430008150033",
}

// DelegatedvalidatorABI is the input ABI used to generate the binding from.
// Deprecated: Use DelegatedvalidatorMetaData.ABI instead.
var DelegatedvalidatorABI = DelegatedvalidatorMetaData.ABI

// DelegatedvalidatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DelegatedvalidatorMetaData.Bin instead.
var DelegatedvalidatorBin = DelegatedvalidatorMetaData.Bin

// DeployDelegatedvalidator deploys a new Ethereum contract, binding an instance of Delegatedvalidator to it.
func DeployDelegatedvalidator(auth *bind.TransactOpts, backend bind.ContractBackend, _arbius common.Address, _baseToken common.Address, _miner common.Address) (common.Address, *types.Transaction, *Delegatedvalidator, error) {
	parsed, err := DelegatedvalidatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DelegatedvalidatorBin), backend, _arbius, _baseToken, _miner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Delegatedvalidator{DelegatedvalidatorCaller: DelegatedvalidatorCaller{contract: contract}, DelegatedvalidatorTransactor: DelegatedvalidatorTransactor{contract: contract}, DelegatedvalidatorFilterer: DelegatedvalidatorFilterer{contract: contract}}, nil
}

// Delegatedvalidator is an auto generated Go binding around an Ethereum contract.
type Delegatedvalidator struct {
	DelegatedvalidatorCaller     // Read-only binding to the contract
	DelegatedvalidatorTransactor // Write-only binding to the contract
	DelegatedvalidatorFilterer   // Log filterer for contract events
}

// DelegatedvalidatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type DelegatedvalidatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegatedvalidatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DelegatedvalidatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegatedvalidatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DelegatedvalidatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegatedvalidatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DelegatedvalidatorSession struct {
	Contract     *Delegatedvalidator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// DelegatedvalidatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DelegatedvalidatorCallerSession struct {
	Contract *DelegatedvalidatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// DelegatedvalidatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DelegatedvalidatorTransactorSession struct {
	Contract     *DelegatedvalidatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// DelegatedvalidatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type DelegatedvalidatorRaw struct {
	Contract *Delegatedvalidator // Generic contract binding to access the raw methods on
}

// DelegatedvalidatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DelegatedvalidatorCallerRaw struct {
	Contract *DelegatedvalidatorCaller // Generic read-only contract binding to access the raw methods on
}

// DelegatedvalidatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DelegatedvalidatorTransactorRaw struct {
	Contract *DelegatedvalidatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDelegatedvalidator creates a new instance of Delegatedvalidator, bound to a specific deployed contract.
func NewDelegatedvalidator(address common.Address, backend bind.ContractBackend) (*Delegatedvalidator, error) {
	contract, err := bindDelegatedvalidator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Delegatedvalidator{DelegatedvalidatorCaller: DelegatedvalidatorCaller{contract: contract}, DelegatedvalidatorTransactor: DelegatedvalidatorTransactor{contract: contract}, DelegatedvalidatorFilterer: DelegatedvalidatorFilterer{contract: contract}}, nil
}

// NewDelegatedvalidatorCaller creates a new read-only instance of Delegatedvalidator, bound to a specific deployed contract.
func NewDelegatedvalidatorCaller(address common.Address, caller bind.ContractCaller) (*DelegatedvalidatorCaller, error) {
	contract, err := bindDelegatedvalidator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DelegatedvalidatorCaller{contract: contract}, nil
}

// NewDelegatedvalidatorTransactor creates a new write-only instance of Delegatedvalidator, bound to a specific deployed contract.
func NewDelegatedvalidatorTransactor(address common.Address, transactor bind.ContractTransactor) (*DelegatedvalidatorTransactor, error) {
	contract, err := bindDelegatedvalidator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DelegatedvalidatorTransactor{contract: contract}, nil
}

// NewDelegatedvalidatorFilterer creates a new log filterer instance of Delegatedvalidator, bound to a specific deployed contract.
func NewDelegatedvalidatorFilterer(address common.Address, filterer bind.ContractFilterer) (*DelegatedvalidatorFilterer, error) {
	contract, err := bindDelegatedvalidator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DelegatedvalidatorFilterer{contract: contract}, nil
}

// bindDelegatedvalidator binds a generic wrapper to an already deployed contract.
func bindDelegatedvalidator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DelegatedvalidatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Delegatedvalidator *DelegatedvalidatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Delegatedvalidator.Contract.DelegatedvalidatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Delegatedvalidator *DelegatedvalidatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.DelegatedvalidatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Delegatedvalidator *DelegatedvalidatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.DelegatedvalidatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Delegatedvalidator *DelegatedvalidatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Delegatedvalidator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Delegatedvalidator *DelegatedvalidatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Delegatedvalidator *DelegatedvalidatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.contract.Transact(opts, method, params...)
}

// BulkSignalCommitment is a paid mutator transaction binding the contract method 0x50ec55e1.
//
// Solidity: function bulkSignalCommitment(bytes32[] commitments_) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) BulkSignalCommitment(opts *bind.TransactOpts, commitments_ [][32]byte) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "bulkSignalCommitment", commitments_)
}

// BulkSignalCommitment is a paid mutator transaction binding the contract method 0x50ec55e1.
//
// Solidity: function bulkSignalCommitment(bytes32[] commitments_) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) BulkSignalCommitment(commitments_ [][32]byte) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.BulkSignalCommitment(&_Delegatedvalidator.TransactOpts, commitments_)
}

// BulkSignalCommitment is a paid mutator transaction binding the contract method 0x50ec55e1.
//
// Solidity: function bulkSignalCommitment(bytes32[] commitments_) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) BulkSignalCommitment(commitments_ [][32]byte) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.BulkSignalCommitment(&_Delegatedvalidator.TransactOpts, commitments_)
}

// BulkSubmitSolution is a paid mutator transaction binding the contract method 0x65d445fb.
//
// Solidity: function bulkSubmitSolution(bytes32[] _taskids, bytes[] _cids) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) BulkSubmitSolution(opts *bind.TransactOpts, _taskids [][32]byte, _cids [][]byte) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "bulkSubmitSolution", _taskids, _cids)
}

// BulkSubmitSolution is a paid mutator transaction binding the contract method 0x65d445fb.
//
// Solidity: function bulkSubmitSolution(bytes32[] _taskids, bytes[] _cids) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) BulkSubmitSolution(_taskids [][32]byte, _cids [][]byte) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.BulkSubmitSolution(&_Delegatedvalidator.TransactOpts, _taskids, _cids)
}

// BulkSubmitSolution is a paid mutator transaction binding the contract method 0x65d445fb.
//
// Solidity: function bulkSubmitSolution(bytes32[] _taskids, bytes[] _cids) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) BulkSubmitSolution(_taskids [][32]byte, _cids [][]byte) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.BulkSubmitSolution(&_Delegatedvalidator.TransactOpts, _taskids, _cids)
}

// Call is a paid mutator transaction binding the contract method 0x6dbf2fa0.
//
// Solidity: function call(address _to, uint256 _value, bytes _data) payable returns(bytes)
func (_Delegatedvalidator *DelegatedvalidatorTransactor) Call(opts *bind.TransactOpts, _to common.Address, _value *big.Int, _data []byte) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "call", _to, _value, _data)
}

// Call is a paid mutator transaction binding the contract method 0x6dbf2fa0.
//
// Solidity: function call(address _to, uint256 _value, bytes _data) payable returns(bytes)
func (_Delegatedvalidator *DelegatedvalidatorSession) Call(_to common.Address, _value *big.Int, _data []byte) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.Call(&_Delegatedvalidator.TransactOpts, _to, _value, _data)
}

// Call is a paid mutator transaction binding the contract method 0x6dbf2fa0.
//
// Solidity: function call(address _to, uint256 _value, bytes _data) payable returns(bytes)
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) Call(_to common.Address, _value *big.Int, _data []byte) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.Call(&_Delegatedvalidator.TransactOpts, _to, _value, _data)
}

// CancelValidatorWithdraw is a paid mutator transaction binding the contract method 0xcbd2422d.
//
// Solidity: function cancelValidatorWithdraw(uint256 count_) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) CancelValidatorWithdraw(opts *bind.TransactOpts, count_ *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "cancelValidatorWithdraw", count_)
}

// CancelValidatorWithdraw is a paid mutator transaction binding the contract method 0xcbd2422d.
//
// Solidity: function cancelValidatorWithdraw(uint256 count_) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) CancelValidatorWithdraw(count_ *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.CancelValidatorWithdraw(&_Delegatedvalidator.TransactOpts, count_)
}

// CancelValidatorWithdraw is a paid mutator transaction binding the contract method 0xcbd2422d.
//
// Solidity: function cancelValidatorWithdraw(uint256 count_) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) CancelValidatorWithdraw(count_ *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.CancelValidatorWithdraw(&_Delegatedvalidator.TransactOpts, count_)
}

// ClaimSolutions is a paid mutator transaction binding the contract method 0x1023ad7c.
//
// Solidity: function claimSolutions(bytes32[] _taskids) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) ClaimSolutions(opts *bind.TransactOpts, _taskids [][32]byte) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "claimSolutions", _taskids)
}

// ClaimSolutions is a paid mutator transaction binding the contract method 0x1023ad7c.
//
// Solidity: function claimSolutions(bytes32[] _taskids) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) ClaimSolutions(_taskids [][32]byte) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.ClaimSolutions(&_Delegatedvalidator.TransactOpts, _taskids)
}

// ClaimSolutions is a paid mutator transaction binding the contract method 0x1023ad7c.
//
// Solidity: function claimSolutions(bytes32[] _taskids) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) ClaimSolutions(_taskids [][32]byte) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.ClaimSolutions(&_Delegatedvalidator.TransactOpts, _taskids)
}

// DelegateVoting is a paid mutator transaction binding the contract method 0x765cb83a.
//
// Solidity: function delegateVoting(address _delegatee) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) DelegateVoting(opts *bind.TransactOpts, _delegatee common.Address) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "delegateVoting", _delegatee)
}

// DelegateVoting is a paid mutator transaction binding the contract method 0x765cb83a.
//
// Solidity: function delegateVoting(address _delegatee) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) DelegateVoting(_delegatee common.Address) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.DelegateVoting(&_Delegatedvalidator.TransactOpts, _delegatee)
}

// DelegateVoting is a paid mutator transaction binding the contract method 0x765cb83a.
//
// Solidity: function delegateVoting(address _delegatee) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) DelegateVoting(_delegatee common.Address) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.DelegateVoting(&_Delegatedvalidator.TransactOpts, _delegatee)
}

// InitiateValidatorWithdraw is a paid mutator transaction binding the contract method 0x0a985737.
//
// Solidity: function initiateValidatorWithdraw(uint256 amount_) returns(uint256)
func (_Delegatedvalidator *DelegatedvalidatorTransactor) InitiateValidatorWithdraw(opts *bind.TransactOpts, amount_ *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "initiateValidatorWithdraw", amount_)
}

// InitiateValidatorWithdraw is a paid mutator transaction binding the contract method 0x0a985737.
//
// Solidity: function initiateValidatorWithdraw(uint256 amount_) returns(uint256)
func (_Delegatedvalidator *DelegatedvalidatorSession) InitiateValidatorWithdraw(amount_ *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.InitiateValidatorWithdraw(&_Delegatedvalidator.TransactOpts, amount_)
}

// InitiateValidatorWithdraw is a paid mutator transaction binding the contract method 0x0a985737.
//
// Solidity: function initiateValidatorWithdraw(uint256 amount_) returns(uint256)
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) InitiateValidatorWithdraw(amount_ *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.InitiateValidatorWithdraw(&_Delegatedvalidator.TransactOpts, amount_)
}

// SafeWithdraw is a paid mutator transaction binding the contract method 0xb0fd035b.
//
// Solidity: function safeWithdraw(uint256 amount) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) SafeWithdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "safeWithdraw", amount)
}

// SafeWithdraw is a paid mutator transaction binding the contract method 0xb0fd035b.
//
// Solidity: function safeWithdraw(uint256 amount) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) SafeWithdraw(amount *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.SafeWithdraw(&_Delegatedvalidator.TransactOpts, amount)
}

// SafeWithdraw is a paid mutator transaction binding the contract method 0xb0fd035b.
//
// Solidity: function safeWithdraw(uint256 amount) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) SafeWithdraw(amount *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.SafeWithdraw(&_Delegatedvalidator.TransactOpts, amount)
}

// SetArbius is a paid mutator transaction binding the contract method 0x4607315c.
//
// Solidity: function setArbius(address _newAribius) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) SetArbius(opts *bind.TransactOpts, _newAribius common.Address) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "setArbius", _newAribius)
}

// SetArbius is a paid mutator transaction binding the contract method 0x4607315c.
//
// Solidity: function setArbius(address _newAribius) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) SetArbius(_newAribius common.Address) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.SetArbius(&_Delegatedvalidator.TransactOpts, _newAribius)
}

// SetArbius is a paid mutator transaction binding the contract method 0x4607315c.
//
// Solidity: function setArbius(address _newAribius) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) SetArbius(_newAribius common.Address) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.SetArbius(&_Delegatedvalidator.TransactOpts, _newAribius)
}

// SetMiner is a paid mutator transaction binding the contract method 0x9742ca46.
//
// Solidity: function setMiner(address _newMiner) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) SetMiner(opts *bind.TransactOpts, _newMiner common.Address) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "setMiner", _newMiner)
}

// SetMiner is a paid mutator transaction binding the contract method 0x9742ca46.
//
// Solidity: function setMiner(address _newMiner) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) SetMiner(_newMiner common.Address) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.SetMiner(&_Delegatedvalidator.TransactOpts, _newMiner)
}

// SetMiner is a paid mutator transaction binding the contract method 0x9742ca46.
//
// Solidity: function setMiner(address _newMiner) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) SetMiner(_newMiner common.Address) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.SetMiner(&_Delegatedvalidator.TransactOpts, _newMiner)
}

// SignalCommitment is a paid mutator transaction binding the contract method 0x506ea7de.
//
// Solidity: function signalCommitment(bytes32 commitment_) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) SignalCommitment(opts *bind.TransactOpts, commitment_ [32]byte) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "signalCommitment", commitment_)
}

// SignalCommitment is a paid mutator transaction binding the contract method 0x506ea7de.
//
// Solidity: function signalCommitment(bytes32 commitment_) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) SignalCommitment(commitment_ [32]byte) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.SignalCommitment(&_Delegatedvalidator.TransactOpts, commitment_)
}

// SignalCommitment is a paid mutator transaction binding the contract method 0x506ea7de.
//
// Solidity: function signalCommitment(bytes32 commitment_) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) SignalCommitment(commitment_ [32]byte) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.SignalCommitment(&_Delegatedvalidator.TransactOpts, commitment_)
}

// SubmitContestation is a paid mutator transaction binding the contract method 0x671f8152.
//
// Solidity: function submitContestation(bytes32 _taskid) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) SubmitContestation(opts *bind.TransactOpts, _taskid [32]byte) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "submitContestation", _taskid)
}

// SubmitContestation is a paid mutator transaction binding the contract method 0x671f8152.
//
// Solidity: function submitContestation(bytes32 _taskid) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) SubmitContestation(_taskid [32]byte) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.SubmitContestation(&_Delegatedvalidator.TransactOpts, _taskid)
}

// SubmitContestation is a paid mutator transaction binding the contract method 0x671f8152.
//
// Solidity: function submitContestation(bytes32 _taskid) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) SubmitContestation(_taskid [32]byte) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.SubmitContestation(&_Delegatedvalidator.TransactOpts, _taskid)
}

// SubmitMultipleTasks is a paid mutator transaction binding the contract method 0xdc0143cc.
//
// Solidity: function submitMultipleTasks(uint8 _tasks) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) SubmitMultipleTasks(opts *bind.TransactOpts, _tasks uint8) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "submitMultipleTasks", _tasks)
}

// SubmitMultipleTasks is a paid mutator transaction binding the contract method 0xdc0143cc.
//
// Solidity: function submitMultipleTasks(uint8 _tasks) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) SubmitMultipleTasks(_tasks uint8) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.SubmitMultipleTasks(&_Delegatedvalidator.TransactOpts, _tasks)
}

// SubmitMultipleTasks is a paid mutator transaction binding the contract method 0xdc0143cc.
//
// Solidity: function submitMultipleTasks(uint8 _tasks) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) SubmitMultipleTasks(_tasks uint8) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.SubmitMultipleTasks(&_Delegatedvalidator.TransactOpts, _tasks)
}

// SubmitSolution is a paid mutator transaction binding the contract method 0x56914caf.
//
// Solidity: function submitSolution(bytes32 _taskid, bytes _cid) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) SubmitSolution(opts *bind.TransactOpts, _taskid [32]byte, _cid []byte) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "submitSolution", _taskid, _cid)
}

// SubmitSolution is a paid mutator transaction binding the contract method 0x56914caf.
//
// Solidity: function submitSolution(bytes32 _taskid, bytes _cid) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) SubmitSolution(_taskid [32]byte, _cid []byte) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.SubmitSolution(&_Delegatedvalidator.TransactOpts, _taskid, _cid)
}

// SubmitSolution is a paid mutator transaction binding the contract method 0x56914caf.
//
// Solidity: function submitSolution(bytes32 _taskid, bytes _cid) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) SubmitSolution(_taskid [32]byte, _cid []byte) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.SubmitSolution(&_Delegatedvalidator.TransactOpts, _taskid, _cid)
}

// SubmitTask is a paid mutator transaction binding the contract method 0x08745dd1.
//
// Solidity: function submitTask(uint8 _version, address _owner, bytes32 _model, uint256 _fee, bytes _input) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) SubmitTask(opts *bind.TransactOpts, _version uint8, _owner common.Address, _model [32]byte, _fee *big.Int, _input []byte) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "submitTask", _version, _owner, _model, _fee, _input)
}

// SubmitTask is a paid mutator transaction binding the contract method 0x08745dd1.
//
// Solidity: function submitTask(uint8 _version, address _owner, bytes32 _model, uint256 _fee, bytes _input) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) SubmitTask(_version uint8, _owner common.Address, _model [32]byte, _fee *big.Int, _input []byte) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.SubmitTask(&_Delegatedvalidator.TransactOpts, _version, _owner, _model, _fee, _input)
}

// SubmitTask is a paid mutator transaction binding the contract method 0x08745dd1.
//
// Solidity: function submitTask(uint8 _version, address _owner, bytes32 _model, uint256 _fee, bytes _input) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) SubmitTask(_version uint8, _owner common.Address, _model [32]byte, _fee *big.Int, _input []byte) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.SubmitTask(&_Delegatedvalidator.TransactOpts, _version, _owner, _model, _fee, _input)
}

// ValidatorDeposit is a paid mutator transaction binding the contract method 0x021c7962.
//
// Solidity: function validatorDeposit(uint256 _amount) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) ValidatorDeposit(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "validatorDeposit", _amount)
}

// ValidatorDeposit is a paid mutator transaction binding the contract method 0x021c7962.
//
// Solidity: function validatorDeposit(uint256 _amount) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) ValidatorDeposit(_amount *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.ValidatorDeposit(&_Delegatedvalidator.TransactOpts, _amount)
}

// ValidatorDeposit is a paid mutator transaction binding the contract method 0x021c7962.
//
// Solidity: function validatorDeposit(uint256 _amount) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) ValidatorDeposit(_amount *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.ValidatorDeposit(&_Delegatedvalidator.TransactOpts, _amount)
}

// ValidatorWithdraw is a paid mutator transaction binding the contract method 0xeaf7c9dc.
//
// Solidity: function validatorWithdraw(uint256 count_) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) ValidatorWithdraw(opts *bind.TransactOpts, count_ *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "validatorWithdraw", count_)
}

// ValidatorWithdraw is a paid mutator transaction binding the contract method 0xeaf7c9dc.
//
// Solidity: function validatorWithdraw(uint256 count_) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) ValidatorWithdraw(count_ *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.ValidatorWithdraw(&_Delegatedvalidator.TransactOpts, count_)
}

// ValidatorWithdraw is a paid mutator transaction binding the contract method 0xeaf7c9dc.
//
// Solidity: function validatorWithdraw(uint256 count_) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) ValidatorWithdraw(count_ *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.ValidatorWithdraw(&_Delegatedvalidator.TransactOpts, count_)
}

// VoteOnContestation is a paid mutator transaction binding the contract method 0x1825c20e.
//
// Solidity: function voteOnContestation(bytes32 _taskid, bool _yea) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) VoteOnContestation(opts *bind.TransactOpts, _taskid [32]byte, _yea bool) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "voteOnContestation", _taskid, _yea)
}

// VoteOnContestation is a paid mutator transaction binding the contract method 0x1825c20e.
//
// Solidity: function voteOnContestation(bytes32 _taskid, bool _yea) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) VoteOnContestation(_taskid [32]byte, _yea bool) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.VoteOnContestation(&_Delegatedvalidator.TransactOpts, _taskid, _yea)
}

// VoteOnContestation is a paid mutator transaction binding the contract method 0x1825c20e.
//
// Solidity: function voteOnContestation(bytes32 _taskid, bool _yea) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) VoteOnContestation(_taskid [32]byte, _yea bool) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.VoteOnContestation(&_Delegatedvalidator.TransactOpts, _taskid, _yea)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0xf14210a6.
//
// Solidity: function withdrawETH(uint256 amount) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactor) WithdrawETH(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.contract.Transact(opts, "withdrawETH", amount)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0xf14210a6.
//
// Solidity: function withdrawETH(uint256 amount) returns()
func (_Delegatedvalidator *DelegatedvalidatorSession) WithdrawETH(amount *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.WithdrawETH(&_Delegatedvalidator.TransactOpts, amount)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0xf14210a6.
//
// Solidity: function withdrawETH(uint256 amount) returns()
func (_Delegatedvalidator *DelegatedvalidatorTransactorSession) WithdrawETH(amount *big.Int) (*types.Transaction, error) {
	return _Delegatedvalidator.Contract.WithdrawETH(&_Delegatedvalidator.TransactOpts, amount)
}
