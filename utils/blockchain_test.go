package utils

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

func TestGetRevertReason(t *testing.T) {
	// Create a new Ethereum client
	client, err := ethclient.Dial("https://arbitrum-nova.publicnode.com")
	assert.NoError(t, err, "failed to connect to Ethereum client")
	defer client.Close()

	// Set up test data
	from := common.HexToAddress("0x5f5731b1b4ca7a5a569f38db9b97cc0c0547f066")

	tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash("0xad2177d8a35065105f2245d12a301107eb34133be86eb6802cc82407b8384639"))
	assert.NoError(t, err, "failed to get transaction from blockchain")

	// Call the function being tested
	revertReason := GetRevertReason(from, tx, client, nil)

	// Assert the expected result
	assert.NoError(t, err, "failed to get revert reason")
	assert.Equal(t, "execution reverted: solution already submitted", revertReason, "unexpected revert reason")

	tx, _, err = client.TransactionByHash(context.Background(), common.HexToHash("0x1f4fad04b8954d01f7580013932406a4a0f8626e0c86d1e8b94551db0d637fe2"))
	assert.NoError(t, err, "failed to get transaction from blockchain")
	receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash("0x1f4fad04b8954d01f7580013932406a4a0f8626e0c86d1e8b94551db0d637fe2"))
	assert.NoError(t, err, "failed to get receipt from blockchain")

	receipt.BlockNumber.Sub(receipt.BlockNumber, big.NewInt(1))
	revertReason = GetRevertReason(from, tx, client, receipt.BlockNumber)
	// Assert the expected result
	assert.NoError(t, err, "failed to get revert reason")
	assert.Equal(t, "execution reverted: solution already submitted", revertReason, "unexpected revert reason")

}

func TestFormatFixed(t *testing.T) {
	value := big.NewInt(1234567890)
	decimals := 6

	expected := "1234.56789"
	result, err := FormatFixed(value, decimals)
	assert.NoError(t, err, "failed to format fixed value")
	assert.Equal(t, expected, result, "unexpected formatted value")

	value = big.NewInt(-9876543210)
	decimals = 8

	expected = "-98.7654321"
	result, err = FormatFixed(value, decimals)
	assert.NoError(t, err, "failed to format fixed value")
	assert.Equal(t, expected, result, "unexpected formatted value")

	value = big.NewInt(0)
	decimals = 2

	expected = "0.0"
	result, err = FormatFixed(value, decimals)
	assert.NoError(t, err, "failed to format fixed value")
	assert.Equal(t, expected, result, "unexpected formatted value")

	value = new(big.Int).Sub(abi.MaxUint256, big.NewInt(1))
	decimals = 0

	expected = value.String()
	result, err = FormatFixed(value, decimals)
	assert.NoError(t, err, "failed to format fixed value")
	assert.Equal(t, expected, result, "unexpected formatted value")
}
