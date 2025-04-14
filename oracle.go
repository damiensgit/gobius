package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"gobius/bindings/quoter" // Added for OnChainOracle
	"gobius/erc20"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient" // Added for OnChainOracle
	"github.com/rs/zerolog"
)

// PriceOracle defines the interface for getting token prices.
type IPriceOracle interface {
	// GetPrices returns the price of the base token (e.g., AIUS) in USD
	// and the price of ETH in USD.
	GetPrices() (basePrice float64, ethPrice float64, err error)
}

type ILeverOracle interface {
	MinClaimLever() (float64, error)
}

// OnChainOracle implements PriceOracle using an on-chain Quoter contract.
type OnChainOracle struct {
	quoter    *quoter.Quoter
	eth       *erc20.TokenERC20
	aius      *erc20.TokenERC20
	quoterRaw *quoter.QuoterRaw
	logger    zerolog.Logger
	timeout   time.Duration
}

// create a new lever oracle that just returns the claim level from configuration used for testing
type MinClaimLeverOracle struct {
	claimLevel float64
}

func NewMinClaimLeverOracle(claimLevel float64) *MinClaimLeverOracle {
	return &MinClaimLeverOracle{claimLevel: claimLevel}
}

func (o *MinClaimLeverOracle) MinClaimLever() (float64, error) {
	return o.claimLevel, nil
}

// NewOnChainOracle creates a new OnChainOracle.
func NewOnChainOracle(client *ethclient.Client, oracleAddress common.Address, eth, aius *erc20.TokenERC20, logger zerolog.Logger) (*OnChainOracle, error) {
	if oracleAddress == (common.Address{}) {
		return nil, errors.New("onchain oracle selected but QuoterAddress is not configured")
	}

	quoterInstance, err := quoter.NewQuoter(oracleAddress, client)
	if err != nil {
		return nil, err
	}

	quoterRaw := &quoter.QuoterRaw{Contract: quoterInstance}

	return &OnChainOracle{
		quoter:    quoterInstance,
		eth:       eth,
		aius:      aius,
		logger:    logger,
		quoterRaw: quoterRaw,
		timeout:   10 * time.Second,
	}, nil
}

func (o *OnChainOracle) MinClaimLever() (float64, error) {
	// timeout after 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	opts := &bind.CallOpts{Context: ctx}
	lever, err := o.quoter.ProfitLevel(opts)
	if err != nil {
		return 0, err
	}
	pl := o.aius.ToFloat(lever)
	return pl, nil
}

// GetPrices fetches prices from the on-chain Quoter contract.
func (o *OnChainOracle) GetPrices() (float64, float64, error) {
	// timeout after 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()
	opts := &bind.CallOpts{Context: ctx}
	// The function GetAIUSPrice is not 'view' or 'pure' in Solidity,
	// so the generated binding expects TransactOpts.
	// We need to use the underlying contract.Call method to perform a read-only eth_call.
	var out []interface{}

	// Define the return types we expect based on the IQuoter interface (uint, uint)
	// We need pointers to these types for the Call method.
	var basePriceReturn *big.Int
	var ethPriceReturn *big.Int

	err := o.quoterRaw.Call(opts, &out, "GetAIUSPrice")

	if err != nil {
		return 0, 0, fmt.Errorf("failed to call GetAIUSPrice via eth_call: %w", err)
	}

	// Manually decode the output based on the expected return types.
	// Ensure the types match exactly what the Solidity function returns.
	if len(out) >= 2 {
		var ok bool
		basePriceReturn, ok = out[0].(*big.Int)
		if !ok {
			return 0, 0, fmt.Errorf("failed to decode basePriceReturn from contract call result")
		}
		ethPriceReturn, ok = out[1].(*big.Int)
		if !ok {
			return 0, 0, fmt.Errorf("failed to decode ethPriceReturn from contract call result")
		}
	} else {
		return 0, 0, fmt.Errorf("ASSERT: unexpected number of return values from GetAIUSPrice call: expected 2, got %d", len(out))
	}

	basePrice := o.aius.ToFloat(basePriceReturn)
	ethPrice := o.eth.ToFloat(ethPriceReturn)

	return basePrice, ethPrice, nil
}
