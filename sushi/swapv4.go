package sushi

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gobius/account"
	"gobius/arbius/basetoken"
	"gobius/erc20"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TODO: this is legacy code for Arbitrum Nova and needs migrating to One

type SwapResponse struct {
	Status             string `json:"status"`
	RouteProcessorAddr string `json:"routeProcessorAddr"`
	RouteProcessorArgs struct {
		TokenIn      string `json:"tokenIn"`
		AmountIn     string `json:"amountIn"`
		TokenOut     string `json:"tokenOut"`
		AmountOutMin string `json:"amountOutMin"`
		To           string `json:"to"`
		RouteCode    string `json:"routeCode"`
	} `json:"routeProcessorArgs"`
}

type SushiManager struct {
	account          *account.Account
	sushiRouter      *bind.BoundContract
	baseToken        *erc20.TokenERC20
	baseTokeContract *basetoken.Basetoken
}

const (
	sushiABI    = `[{"constant":false,"inputs":[{"name":"tokenIn","type":"address"},{"name":"amountIn","type":"uint256"},{"name":"tokenOut","type":"address"},{"name":"amountOutMin","type":"uint256"},{"name":"to","type":"address"},{"name":"route","type":"bytes"}],"name":"processRoute","outputs":[{"name":"amountOut","type":"uint256"}],"payable":true,"stateMutability":"payable","type":"function"}]`
	baseURL     = "https://api.sushi.com"
	sushiRouter = "0xCdBCd51a5E8728E0AF4895ce5771b7d17fF71959"
)

var (
	novaNativeEth = common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")
	novaEth       = common.HexToAddress("0x722E8BdD2ce80A4422E880164f2079488e115365")
)

func NewSushiContract(account *account.Account, baseTokeContract *basetoken.Basetoken, baseToken *erc20.TokenERC20) *SushiManager {

	contractAddress := common.HexToAddress(sushiRouter)

	parsedABI, err := abi.JSON(strings.NewReader(sushiABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	contract := bind.NewBoundContract(contractAddress, parsedABI, account.Client.Client, account.Client.Client, account.Client.Client)
	if err != nil {
		log.Fatalf("Failed to create new contract: %v", err)
	}

	return &SushiManager{
		account:          account,
		sushiRouter:      contract,
		baseTokeContract: baseTokeContract,
		baseToken:        baseToken,
	}

}

func (sushi *SushiManager) ProcessRoute(tokenIn common.Address, amountIn *big.Int, tokenOut common.Address, amountOutMin *big.Int, to common.Address, route []byte) (*types.Transaction, error) {
	return sushi.account.NonceManagerWrapper(5, 425, 1.5, true, func(opts *bind.TransactOpts) (interface{}, error) {
		opts.GasLimit = 0
		return sushi.sushiRouter.Transact(opts, "processRoute", tokenIn, amountIn, tokenOut, amountOutMin, to, route)
	})
}

func (sushi *SushiManager) Allowance(balance float64) (*types.Transaction, error) {

	minAllowance := sushi.baseToken.FromFloat(balance)

	contractAddress := common.HexToAddress(sushiRouter)

	allowance, err := sushi.baseTokeContract.Allowance(nil, sushi.account.Address, contractAddress)
	if err != nil {
		return nil, err
	}

	fmt.Printf("AIUS allowance: %s for account %s on contract %s", allowance.String(), sushi.account.Address, contractAddress)

	// check if the allowance is less than the balance
	if allowance.Cmp(minAllowance) < 0 {
		fmt.Printf("approving AIUS")

		allowanceAmount := new(big.Int).Sub(abi.MaxUint256, allowance)

		gp, gasFeeCap, gasFeeTip, _ := sushi.account.Client.GasPriceOracle(true)
		opts := sushi.account.GetOpts(0, gp, gasFeeCap, gasFeeTip)

		// increase the allowance
		tx, err := sushi.baseTokeContract.Approve(opts, contractAddress, allowanceAmount)
		if err != nil {
			return nil, err
		}

		return tx, nil

	}

	return nil, nil
}

// Assumes allowance has been increased for the contract
func (sushi *SushiManager) SellAius(amountToSell float64) (*types.Transaction, error) {

	amountIn := sushi.baseToken.FromFloat(amountToSell)

	gasPrice := 100

	quote, err := sushi.GetSwapV4(sushi.baseToken.Address, novaNativeEth, amountIn, 0.005, gasPrice, sushi.account.Address, true)

	if err != nil {
		return nil, err
	}

	if quote.Status != "Success" {
		return nil, fmt.Errorf("failed to get quote: %s", quote.Status)
	}

	if quote.RouteProcessorAddr != sushiRouter {
		return nil, fmt.Errorf("unexpected router processor address: %s", quote.RouteProcessorAddr)
	}

	//fmt.Printf("Swap v4 router: %s\nAmount In: %0.4f\nAmount Out: %s\n", quote.RouteProcessorAddr, amountToSell, quote.RouteProcessorArgs.AmountOutMin)

	amountOutMin, _ := new(big.Int).SetString(quote.RouteProcessorArgs.AmountOutMin, 10)

	// 0x028AFE4055Ebc86Bd2AFB3940c0095C9aca511d85201ffff019b614cb49880AEE59537fd21D106aed03171438f006c3Db6ef57735B8b62D0bdDa32c94389933d2f5d
	hexString := strings.TrimPrefix(quote.RouteProcessorArgs.RouteCode, "0x")

	rourceCode, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}

	return sushi.ProcessRoute(sushi.baseToken.Address, amountIn, novaNativeEth, amountOutMin, sushi.account.Address, rourceCode)
}

// Assumes allowance has been increased for the contract
func (sushi *SushiManager) BuyAius(amountOfEth float64) (*types.Transaction, error) {

	amountIn := sushi.baseToken.FromFloat(amountOfEth)

	gasPrice := 100

	quote, err := sushi.GetSwapV4(novaNativeEth, sushi.baseToken.Address, amountIn, 0.005, gasPrice, sushi.account.Address, true)

	if err != nil {
		return nil, err
	}

	if quote.Status != "Success" {
		return nil, fmt.Errorf("failed to get quote: %s", quote.Status)
	}

	if quote.RouteProcessorAddr != sushiRouter {
		return nil, fmt.Errorf("unexpected router processor address: %s", quote.RouteProcessorAddr)
	}

	fmt.Printf("Swap v4 router: %s\nAmount In: %0.4f\nAmount Out: %s\n", quote.RouteProcessorAddr, amountOfEth, quote.RouteProcessorArgs.AmountOutMin)

	// amountOutMin, _ := new(big.Int).SetString(quote.RouteProcessorArgs.AmountOutMin, 10)

	// hexString := strings.TrimPrefix(quote.RouteProcessorArgs.RouteCode, "0x")

	// rourceCode, err := hex.DecodeString(hexString)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil //sushi.processRoute(sushi.baseToken.Address, amountIn, novaNativeEth, amountOutMin, sushi.account.Address, rourceCode)
}

func (sushi *SushiManager) GetSwapV4(tokenIn, tokenOut common.Address, amountIn *big.Int, maxPriceImpact float64, gasPrice int, to common.Address, preferSushi bool) (*SwapResponse, error) {
	params := url.Values{}
	params.Add("tokenIn", tokenIn.String())
	params.Add("tokenOut", tokenOut.String())
	params.Add("amount", amountIn.String())
	params.Add("maxPriceImpact", fmt.Sprintf("%g", maxPriceImpact))
	params.Add("gasPrice", strconv.Itoa(gasPrice))
	params.Add("to", to.String())
	params.Add("preferSushi", fmt.Sprintf("%v", preferSushi))

	fullURL := fmt.Sprintf("%s/swap/v4/42170?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get swap data: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var swapResponse SwapResponse
	err = json.Unmarshal(body, &swapResponse)
	if err != nil {
		return nil, err
	}

	return &swapResponse, nil
}

func (sushi *SushiManager) GetTokenPrice(token common.Address) (float64, error) {
	fullURL := fmt.Sprintf("%s/price/v1/42170/%s", baseURL, token.String())

	response, err := http.Get(fullURL)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get token price: %s", response.Status)
	}

	priceData, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var finalPrice float64
	err = json.Unmarshal(priceData, &finalPrice)
	if err != nil {
		return 0, err
	}

	return finalPrice, nil
}

func (sushi *SushiManager) GetPrices() (arbiusPrice, ethereumPrice float64, err error) {

	arbiusPrice, err = sushi.GetTokenPrice(common.HexToAddress("0x8AFE4055Ebc86Bd2AFB3940c0095C9aca511d852"))
	if err != nil {
		return 0, 0, err
	}

	ethereumPrice, err = sushi.GetTokenPrice(novaEth)
	if err != nil {
		return 0, 0, err
	}
	return
}
