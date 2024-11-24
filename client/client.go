package client

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"reflect"
	"sync"
	"unsafe"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// func connectClient(connectionURL string) (*ethclient.Client, error) {
// 	// Connect to the Ethereum client
// 	client, err := ethclient.Dial(connectionURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return client, nil
// }

type Client struct {
	Client    *ethclient.Client
	RpcClient *rpc.Client
	ChainID   *big.Int
	ctx       context.Context

	mutex                      sync.Mutex
	basefee                    *big.Int
	suggestedGasTipCap         *big.Int
	useEthersMethodForGasPrice bool

	basefeeX    float64
	forcegas    bool
	gasoverride *big.Int
}

func NewClient(connectionURL string, appCtx context.Context, useEthersMethodForGasPrice bool, baseFeeX float64, forcegas bool, gasoverridegwei float64) (*Client, error) {
	// Connect to the Ethereum client
	ethClient, err := ethclient.Dial(connectionURL)
	if err != nil {
		return nil, err
	}

	rpcClient := getRPCClient(ethClient)

	// Get the chain ID
	chainID, err := ethClient.NetworkID(appCtx)
	if err != nil {
		return nil, err
	}

	gasoverride := big.NewInt(1000000000)
	if forcegas {
		gasoverride = MulBigByFloat(gasoverride, gasoverridegwei)
	}

	return &Client{
		Client:                     ethClient,
		RpcClient:                  rpcClient,
		ChainID:                    chainID,
		ctx:                        appCtx,
		useEthersMethodForGasPrice: useEthersMethodForGasPrice,
		basefeeX:                   baseFeeX,
		forcegas:                   forcegas,
		gasoverride:                gasoverride,
	}, nil
}

func getRPCClient(client *ethclient.Client) *rpc.Client {

	clientValue := reflect.ValueOf(client).Elem()
	fieldStruct := clientValue.FieldByName("c")
	clientPointer := reflect.NewAt(fieldStruct.Type(), unsafe.Pointer(fieldStruct.UnsafeAddr())).Elem()
	finalClient, _ := clientPointer.Interface().(*rpc.Client)
	return finalClient
}

func (c *Client) GetGasPrice() (*big.Int, error) {

	suggestedGasPrice, err := c.Client.SuggestGasPrice(c.ctx)
	if err != nil {
		return nil, err
	}

	return suggestedGasPrice, nil
}

func (c *Client) SendSignedTransaction(signedTx *types.Transaction) (*types.Transaction, error) {
	err := c.Client.SendTransaction(c.ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

var FILTER_HIGH_GAS_PRICE = big.NewInt(1500000000 * 10)

// Manually set the basefee (usually from a block header)
// TODO: make these sets/gets thread safe
func (c *Client) SetBaseFee(fee *big.Int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if fee.Cmp(FILTER_HIGH_GAS_PRICE) > 0 {
		return
	}

	c.basefee = fee
}

func (c *Client) GasTipCap() (*big.Int, error) {
	if c.suggestedGasTipCap == nil {
		err := c.UpdateGasTipCap()
		if err != nil {
			return nil, err
		}
	}
	return c.suggestedGasTipCap, nil
}

func (c *Client) GetBaseFee() (*big.Int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.basefee == nil {
		err := c.updateCurrentBasefee()
		if err != nil {
			return nil, err
		}
	}
	return c.basefee, nil
}

func (c *Client) UpdateCurrentBasefee() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.updateCurrentBasefee()
}

func (c *Client) updateCurrentBasefee() error {
	head, err := c.Client.HeaderByNumber(c.ctx, nil)
	if err != nil {
		return err
	}

	c.basefee = head.BaseFee

	return nil
}

func (c *Client) UpdateGasTipCap() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	suggestedGasTipCap, err := c.Client.SuggestGasTipCap(c.ctx)
	if err != nil {
		return err
	}

	c.suggestedGasTipCap = suggestedGasTipCap
	return nil
}

var maxFeePerGas = big.NewInt(1500000000)

func MulBigByFloat(amount1 *big.Int, amount2 float64) *big.Int {

	bigval := new(big.Float)
	bigval.SetFloat64(amount2)
	// Set precision if required.
	bigval.SetPrec(256)

	coin := new(big.Float)
	coin.SetInt(amount1)

	bigval.Mul(bigval, coin)

	result := new(big.Int)
	bigval.Int(result) // store converted number in result

	return result
}

func (c *Client) GasPriceOracle(override bool) (gasPrice *big.Int, GasFeeCap *big.Int, GasFeeTip *big.Int, err error) {

	basefee, err := c.GetBaseFee()
	if err != nil {
		return nil, nil, nil, err
	}

	// override ignores all config settings and other defaults to ensure caller can send regardless
	if override {
		GasFeeCap = new(big.Int).Mul(basefee, big.NewInt(2))
		GasFeeTip = big.NewInt(0)
		return gasPrice, GasFeeCap, GasFeeTip, nil
	}

	if c.forcegas {
		GasFeeCap := new(big.Int).Set(c.gasoverride)
		GasFeeTip = big.NewInt(0)
		return gasPrice, GasFeeCap, GasFeeTip, nil
	}

	// Use the ethers method for gas price calculation
	if c.useEthersMethodForGasPrice {
		GasFeeTip = big.NewInt(1500000000)
		GasFeeCap = new(big.Int).Add(
			GasFeeTip,
			new(big.Int).Mul(basefee, big.NewInt(2)),
		)
		GasFeeTip = big.NewInt(0)
		return nil, GasFeeCap, GasFeeTip, nil
	} else {
		/*c.UpdateGasTipCap()
		suggestedGasTipCap, err := c.GasTipCap()
		if err != nil {
			if strings.Contains(err.Error(), "Method eth_maxPriorityFeePerGas not found") {
				GetGasPrice, err := c.GetGasPrice()
				if err == nil {
					return GetGasPrice, nil, nil, nil
				}
			}
			return nil, nil, nil, err
		}

		GasFeeCap = new(big.Int).Add(
			suggestedGasTipCap,
			new(big.Int).Mul(basefee, big.NewInt(2)),
		)*/

		GasFeeCap = MulBigByFloat(basefee, c.basefeeX)
		//GasFeeCap = new(big.Int).Mul(basefee, big.NewInt(2))

		if GasFeeCap.Cmp(maxFeePerGas) >= 0 {
			log.Printf("GasFeeCap is too high: %s (basefee: %s), setting to maxFeePerGas: %s", GasFeeCap.String(), basefee.String(), maxFeePerGas.String())
			GasFeeCap.Set(maxFeePerGas)
			// 	return nil, nil, nil, fmt.Errorf("maxFeePerGas (%s) < maxPriorityFeePerGas (%s)", GasFeeCap.String(), suggestedGasTipCap.String())
		}
		GasFeeTip = big.NewInt(0)
		//GasFeeTip = new(big.Int).Set(basefee) // new(big.Int).Add(suggestedGasTipCap, basefee)
	}
	return gasPrice, GasFeeCap, GasFeeTip, nil
}

// Seems to do nothing on Nova
func (c *Client) FeeHistory() (*ethereum.FeeHistory, error) {
	rewardPercentiles := []float64{50}
	reward, err := c.Client.FeeHistory(c.ctx, 10, nil, rewardPercentiles)
	if err != nil {
		return nil, err
	}

	return reward, nil
}

func toFilterArg(q ethereum.FilterQuery) (interface{}, error) {
	arg := map[string]interface{}{
		"address": q.Addresses,
		"topics":  q.Topics,
	}
	if q.BlockHash != nil {
		arg["blockHash"] = *q.BlockHash
		if q.FromBlock != nil || q.ToBlock != nil {
			return nil, fmt.Errorf("cannot specify both BlockHash and FromBlock/ToBlock")
		}
	} else {
		if q.FromBlock == nil {
			arg["fromBlock"] = "0x0"
		} else {
			arg["fromBlock"] = toBlockNumArg(q.FromBlock)
		}
		arg["toBlock"] = toBlockNumArg(q.ToBlock)
	}

	return arg, nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	latest := big.NewInt(-1)
	pending := big.NewInt(-2)
	finalized := big.NewInt(-3)
	if number.Cmp(latest) == 0 {
		return "latest"
	}
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	if number.Cmp(finalized) == 0 {
		return "finalized"
	}
	return hexutil.EncodeBig(number)
}

func (c *Client) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	var result []types.Log
	arg, err := toFilterArg(q)
	if err != nil {
		return nil, err
	}
	err = c.RpcClient.CallContext(ctx, &result, "eth_getLogs", arg)
	return result, err
}
