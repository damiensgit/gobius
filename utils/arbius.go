package utils

import (
	task "gobius/common"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// GenerateCommitment generates a commitment hash for a given taskid, cid and address
func GenerateCommitment(address common.Address, taskid task.TaskId, cid []byte) (commitment [32]byte, err error) {
	param1, _ := abi.NewType("address", "", nil)
	param2, _ := abi.NewType("bytes32", "", nil)
	param3, _ := abi.NewType("bytes", "", nil)

	//args := []interface{}{address, taskid, cid}
	encodedData, err := (abi.Arguments{{Type: param1}, {Type: param2}, {Type: param3}}).Pack(address, taskid, cid)

	if err != nil {
		return commitment, err
	}

	hash := crypto.Keccak256(encodedData)
	copy(commitment[:], hash)

	return commitment, nil
}
