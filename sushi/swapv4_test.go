package sushi

import (
	"context"
	"fmt"
	"gobius/account"
	"gobius/arbius/basetoken"
	"gobius/client"
	"gobius/config"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestGetSwapV4(t *testing.T) {

	cfg, err := config.InitAppConfig("../delegated.json", 0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	appContext := context.Background()

	rpcClient, err := client.NewClient(cfg.Blockchain.RPCURL, appContext, cfg.Blockchain.EthersGas, 0, false, 0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	//https://api.sushi.com/swap/v4/42170?tokenIn=0x8AFE4055Ebc86Bd2AFB3940c0095C9aca511d852&tokenOut=0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE&amount=2000000000000000000&maxPriceImpact=0.005&gasPrice=67020000&to=0x6c3Db6ef57735B8b62D0bdDa32c94389933d2f5d&preferSushi=true

	tokenIn := common.HexToAddress("0x8AFE4055Ebc86Bd2AFB3940c0095C9aca511d852")
	tokenOut := common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")
	amount, _ := new(big.Int).SetString("2000000000000000000", 10)
	maxPriceImpact := 0.01
	gasPrice := 1000000000
	to := common.HexToAddress("0x111E4055Ebc86Bd2AFB394010015C1aca511d851")
	preferSushi := true

	expectedResponse := &SwapResponse{
		// Fill in the expected response values here
	}

	account, err := account.NewAccount(cfg.Blockchain.PrivateKey, rpcClient, context.Background(), false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	baseTokenContract, err := basetoken.NewBasetoken(cfg.BaseConfig.BaseTokenAddress, rpcClient.Client)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	sushi := NewSushiContract(account, baseTokenContract, cfg.BaseConfig.BaseToken)

	_, err = sushi.SellAius(0.1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	a, b, err := sushi.GetPrices()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	fmt.Printf("Arbius Price: %.4f\n", a)
	fmt.Printf("Ethereum Price: %.4f\n", b)

	// Call the function under test
	response, err := sushi.GetSwapV4(tokenIn, tokenOut, amount, maxPriceImpact, gasPrice, to, preferSushi)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// Verify the response
	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("unexpected response: got %v, want %v", response, expectedResponse)
	}
}
