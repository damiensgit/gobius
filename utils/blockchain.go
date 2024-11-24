package utils

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetRevertReason(from common.Address, tx *types.Transaction, client *ethclient.Client, blockNo *big.Int) string {

	callMsg := ethereum.CallMsg{
		From:  from,
		To:    tx.To(),
		Gas:   tx.Gas(),
		Value: tx.Value(),
		Data:  tx.Data(),
	}
	if tx.Type() == 2 {
		callMsg.GasFeeCap = tx.GasTipCap()
		callMsg.GasTipCap = tx.GasTipCap()
	} else {
		callMsg.GasPrice = tx.GasPrice()
	}

	response, err := client.CallContract(context.Background(), callMsg, blockNo)
	if err == nil {
		revertreason, err := abi.UnpackRevert(response)

		if err == nil {
			return fmt.Sprintf("REVERTED! %s", revertreason)
		} else {
			return fmt.Sprintf("could not decode reverted reason: %s", err)
		}
	} else {
		return err.Error()
	}
}

var Zero = big.NewInt(0)
var NegativeOne = big.NewInt(-1)

func getMultiplier(decimals int) (*big.Int, error) {
	if decimals >= 0 && decimals <= 256 {

		var base, exponent = big.NewInt(10), big.NewInt(int64(decimals))
		val := base.Exp(base, exponent, nil)

		//	val, _ := new(big.Int).SetString("1"+zeros[:decimals], 10)
		return val, nil
	}

	return nil, fmt.Errorf("invalid decimals: %d", decimals)
}

// FormatFixed formats a fixed-point big.Int number with the given number of decimals.
// returns a string representation of the number.
func FormatFixed(value *big.Int, decimals int) (string, error) {
	multiplier, err := getMultiplier(decimals)
	if err != nil {
		return "", err
	}

	negative := value.Sign() < 0
	if negative {
		value.Neg(value)
	}

	fraction := new(big.Int).Mod(value, multiplier).String()
	for len(fraction) < len(multiplier.String())-1 {
		fraction = "0" + fraction
	}

	fraction = strings.TrimRight(fraction, "0")
	if fraction == "" {
		fraction = "0"
	}

	valueAsStr := ""

	whole := new(big.Int).Div(value, multiplier).String()
	if len(multiplier.String()) == 1 {
		valueAsStr = whole
	} else {
		valueAsStr = whole + "." + fraction
	}

	if negative {
		valueAsStr = "-" + valueAsStr
	}

	return valueAsStr, nil
}
