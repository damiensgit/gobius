package utils

import (
	"bytes"
	"encoding/hex"
	task "gobius/common"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestGenerateCommitmentB(t *testing.T) {
	address := common.HexToAddress("0x1A320E53A25f518B893F286f3600cc204c181a8E")
	taskid := task.TaskId(common.HexToHash("dc6a147f2cd937a1b290b0dc4eff49a084ad72468fcb36df6d8ddb00c5ff6f7b").Bytes())
	cid, _ := hex.DecodeString("122002e2550d45270ed9c0df80be9a331940d391bcb138e0edfce4d2ff20168d6691")

	expected, _ := hex.DecodeString("b5fa4a5b4a1febe1ac696ace872867f629cfbed3996c128a25f1aafd56b7cb50")

	result, err := GenerateCommitment(address, taskid, cid)

	if err != nil {
		t.Errorf("GenerateCommitmentB() returned an error: %v", err)
	}

	if !bytes.Equal(result[:], expected) {
		t.Errorf("GenerateCommitmentB() returned incorrect result, got: %v, want: %v", result, expected)
	}
}
