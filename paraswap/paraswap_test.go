package paraswap

import (
	"context"
	"fmt"
	"gobius/account"
	"gobius/arbius/basetoken"
	"gobius/client"
	"gobius/config"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func TestGetSwapV4(t *testing.T) {
	// Create a test logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg, err := config.InitAppConfig("../testinguniswap.json", 0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	appContext := context.Background()

	rpcClient, err := client.NewClient(cfg.Blockchain.RPCURL, appContext, cfg.Blockchain.EthersGas, 0, false, 0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	account, err := account.NewAccount(cfg.Blockchain.PrivateKey, rpcClient, context.Background(), false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	fmt.Println(account.Address.String())

	baseTokenContract, err := basetoken.NewBasetoken(cfg.BaseConfig.BaseTokenAddress, rpcClient.Client)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	paraswapManager := NewParaswapManager(account, baseTokenContract, cfg.BaseConfig.BaseToken, &logger)

	a, b, err := paraswapManager.GetPrices()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	fmt.Printf("Arbius Price: %.4f\n", a)
	fmt.Printf("Ethereum Price: %.4f\n", b)

	allowancetx, err := paraswapManager.Allowance(0.01)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if allowancetx != nil {
		_, success, _, _ := account.WaitForConfirmedTx(&logger, allowancetx)

		if !success {
			t.Errorf("unexpected error: %v", err)
			return
		}
	}

	_, err = paraswapManager.SellAius(0.01)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

}
