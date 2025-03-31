package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"gobius/client"
	"gobius/config"
	"os"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type ClaimTask struct {
	ID   string `json:"id"`
	Time int64  `json:"time"`
}

func WriteClaimTasksToJsonFile(data []ClaimTask, filename string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// GenerateCommitment generates a commitment hash for a given taskid, cid and address
func GetTasksFromTaskData(data []byte) (tasks [][32]byte, err error) {
	param1, _ := abi.NewType("bytes32[]", "", nil)

	unpacked, err := (abi.Arguments{{Type: param1}}).Unpack(data[4:])

	if err != nil {
		return nil, err
	}

	return unpacked[0].([][32]byte), nil

}

func TestTaskInput(t *testing.T) {

	//input := `{"prompt":"test"}`

	input := config.Input{
		Prompt: "test",
	}
	outputbytes, err := json.Marshal(input)
	assert.NoError(t, err, "failed to marshal input")

	expectedStr := "7b2270726f6d7074223a2274657374227d"
	expected, err := hex.DecodeString(expectedStr)
	assert.NoError(t, err, "failed to decode")

	if !bytes.Equal(outputbytes, expected) {
		t.Errorf("incorrect result, got: %v, want: %v", outputbytes, expected)
	}
	//{"prompt":"test"}
	//0x7b2270726f6d7074223a2274657374227d
}

func TestValidatorWithdraw(t *testing.T) {

	cfg, err := config.InitAppConfig("testnet.json", 2)
	if err != nil {
		assert.NoError(t, err, "failed to load config file")
	}

	appContext := context.Background()

	rpcClient, err := client.NewClient(cfg.Blockchain.RPCURL, appContext, cfg.Blockchain.EthersGas, 0, false, 0)

	assert.NoError(t, err, "failed to init rpc client")

	// TODO: make generic test logger
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05.000000000"}).Level(zerolog.NoLevel).With().Timestamp().Logger()

	_ = logger
	_ = rpcClient

	var appQuitWG sync.WaitGroup
	validator, err := NewBatchTransactionManager(appServices, appContext, &appQuitWG)

	if err != nil {
		logger.Fatal().Err(err).Msg("could not create transaction manager")
	}

	miner, err := NewMinerEngine(context.Background(), validator, &appQuitWG)

	assert.NoError(t, err, "could not create miner engine")

	validatorAddress := miner.validator.GetNextValidatorAddress()

	err = miner.validator.InitiateValidatorWithdraw(validatorAddress, 0.00001)
	assert.NoError(t, err, "error in initiate validator withdraw")

	err = miner.validator.ValidatorWithdraw(validatorAddress)
	assert.NoError(t, err, "error in  validator withdraw")

}
