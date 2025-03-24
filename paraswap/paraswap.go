package paraswap

import (
	"encoding/json"
	"fmt"
	"gobius/account"
	"gobius/bindings/basetoken"
	"gobius/erc20"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
)

const (
	paraswapAPI      = "https://api.paraswap.io"
	chainID          = "42161"                                      // Arbitrum One
	ETH_ADDRESS      = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE" // special address for ETH
	USDC_ADDRESS     = "0xaf88d065e77c8cC2239327C5EDb3A432268e5831" // USDC
	AUGUSTUS_ADDRESS = "0x6a000f20005980200259b80c5102003040001068" // Augustus V6.2 Router
	AIUS_ADDRESS     = "0x4a24B101728e07A52053c13FB4dB2BcF490CAbc3" // AIUS
)

var (
	usdcToken = erc20.NewTokenERC20(common.HexToAddress(USDC_ADDRESS), 6, "USDC", "USDC")
	ethToken  = erc20.NewTokenERC20(common.HexToAddress(ETH_ADDRESS), 18, "ETH", "ETH")
	aiusToken = erc20.NewTokenERC20(common.HexToAddress(AIUS_ADDRESS), 18, "AIUS", "AIUS")
)

type ParaswapManager struct {
	account           *account.Account
	baseToken         *erc20.TokenERC20
	baseTokenContract *basetoken.BaseToken
	logger            *zerolog.Logger
}

type PriceResponse struct {
	PriceRoute struct {
		BlockNumber        int           `json:"blockNumber"`
		Network            int           `json:"network"`
		SrcToken           string        `json:"srcToken"`
		SrcDecimals        int           `json:"srcDecimals"`
		SrcAmount          string        `json:"srcAmount"`
		DestToken          string        `json:"destToken"`
		DestDecimals       int           `json:"destDecimals"`
		DestAmount         string        `json:"destAmount"`
		DestUSD            string        `json:"destUSD"`
		BestRoute          []interface{} `json:"bestRoute"`
		GasCost            string        `json:"gasCost"`
		GasCostUSD         string        `json:"gasCostUSD"`
		Side               string        `json:"side"`
		TokenTransferProxy string        `json:"tokenTransferProxy"`
		ContractAddress    string        `json:"contractAddress"`
	} `json:"priceRoute"`
}

type TransactionResponse struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data"`
	GasPrice string `json:"gasPrice"`
	Gas      string `json:"gas"`
	ChainID  int    `json:"chainId"`
}

func NewParaswapManager(account *account.Account, baseTokenContract *basetoken.BaseToken, baseToken *erc20.TokenERC20, logger *zerolog.Logger) *ParaswapManager {
	return &ParaswapManager{
		account:           account,
		baseToken:         baseToken,
		baseTokenContract: baseTokenContract,
		logger:            logger,
	}
}

func (p *ParaswapManager) GetPrices() (aiusPrice float64, ethPrice float64, err error) {
	// Get AIUS price using 1 AIUS as input
	oneAius := big.NewInt(1000000000000000000) // 1 AIUS
	aiusQuote, err := p.GetPrice(aiusToken, ethToken, oneAius)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get aius price: %v", err)
	}

	// Get ETH price using 1 ETH as input
	oneEth := big.NewInt(1000000000000000000) // 1 ETH
	ethQuote, err := p.GetPrice(ethToken, usdcToken, oneEth)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get eth price: %v", err)
	}

	// Parse USD prices from destUSD field
	aiusPrice, err = strconv.ParseFloat(aiusQuote.PriceRoute.DestUSD, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse aius usd price: %v", err)
	}

	ethPrice, err = strconv.ParseFloat(ethQuote.PriceRoute.DestUSD, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse eth usd price: %v", err)
	}

	return aiusPrice, ethPrice, nil
}

// TODO: make this a context bound http call via NewRequestWithContext
func (p *ParaswapManager) GetPrice(srcToken, destToken *erc20.TokenERC20, amount *big.Int) (*PriceResponse, error) {
	params := url.Values{}
	params.Add("srcToken", srcToken.Address.String())
	params.Add("destToken", destToken.Address.String())
	params.Add("amount", amount.String())
	params.Add("srcDecimals", strconv.Itoa(int(srcToken.Decimals)))
	params.Add("destDecimals", strconv.Itoa(int(destToken.Decimals)))
	params.Add("side", "SELL")
	params.Add("network", chainID)
	params.Add("version", "6.2")
	params.Add("slippage", "100") // 1% slippage

	url := fmt.Sprintf("%s/prices?%s", paraswapAPI, params.Encode())

	p.logger.Debug().
		Str("url", url).
		Msg("Paraswap API request")

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get price: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var priceResp PriceResponse
	if err := json.Unmarshal(body, &priceResp); err != nil {
		return nil, fmt.Errorf("failed to parse price response: %v", err)
	}

	return &priceResp, nil
}

func (p *ParaswapManager) GetTransaction(priceRoute *PriceResponse) (*TransactionResponse, error) {
	data := map[string]interface{}{
		"priceRoute":  priceRoute.PriceRoute,
		"srcToken":    priceRoute.PriceRoute.SrcToken,
		"destToken":   priceRoute.PriceRoute.DestToken,
		"srcAmount":   priceRoute.PriceRoute.SrcAmount,
		"destAmount":  priceRoute.PriceRoute.DestAmount,
		"userAddress": p.account.Address.String(),
		"partner":     "paraswap.io",
		"version":     "6.2",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction data: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/transactions/%s", paraswapAPI, chainID), strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var txResp TransactionResponse
	if err := json.Unmarshal(body, &txResp); err != nil {
		return nil, fmt.Errorf("failed to parse transaction response: %v", err)
	}

	return &txResp, nil
}

func (p *ParaswapManager) executeSwap(srcToken, destToken string, amount *big.Int, srcDecimals, destDecimals int) (*types.Transaction, error) {
	// Build parameters for /swap endpoint
	params := url.Values{}
	params.Add("srcToken", srcToken)
	params.Add("destToken", destToken)
	params.Add("amount", amount.String())
	params.Add("srcDecimals", strconv.Itoa(srcDecimals))
	params.Add("destDecimals", strconv.Itoa(destDecimals))
	params.Add("userAddress", p.account.Address.String())
	params.Add("side", "SELL")
	params.Add("network", chainID)
	params.Add("partner", "paraswap.io")
	params.Add("version", "6.2")
	params.Add("slippage", "100") // 1% slippage (100 basis points)

	url := fmt.Sprintf("%s/swap?%s", paraswapAPI, params.Encode())

	p.logger.Debug().
		Str("url", url).
		Msg("Paraswap API request")

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get swap data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	p.logger.Debug().
		Int("status", resp.StatusCode).
		Str("body", string(body)).
		Msg("Paraswap API response")

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var swapResp struct {
		PriceRoute PriceResponse `json:"priceRoute"`
		TxParams   struct {
			From     string `json:"from"`
			To       string `json:"to"`
			Value    string `json:"value"`
			Data     string `json:"data"`
			GasPrice string `json:"gasPrice"`
			Gas      string `json:"gas"`
			ChainId  int    `json:"chainId"`
		} `json:"txParams"`
	}

	if err := json.Unmarshal(body, &swapResp); err != nil {
		return nil, fmt.Errorf("failed to parse swap response: %v", err)
	}

	value := new(big.Int)
	if swapResp.TxParams.Value != "" {
		value.SetString(swapResp.TxParams.Value, 10)
	}
	contract := common.HexToAddress(swapResp.TxParams.To)

	// Execute the swap
	return p.account.NonceManagerWrapper(5, 425, 1.5, true, func(opts *bind.TransactOpts) (interface{}, error) {
		opts.GasLimit = 0
		opts.Value = value
		return p.account.SendTransactionWithOpts(
			opts,
			&contract,
			common.FromHex(swapResp.TxParams.Data),
		)
	})
}

func (p *ParaswapManager) SellAius(amountToSell float64) (*types.Transaction, error) {
	amountIn := p.baseToken.FromFloat(amountToSell)
	return p.executeSwap(
		p.baseToken.Address.String(),
		ETH_ADDRESS,
		amountIn,
		int(p.baseToken.Decimals),
		18,
	)
}

func (p *ParaswapManager) BuyAius(amountOfEth float64) (*types.Transaction, error) {
	amountIn := p.baseToken.FromFloat(amountOfEth)
	return p.executeSwap(
		ETH_ADDRESS,
		p.baseToken.Address.String(),
		amountIn,
		18,
		int(p.baseToken.Decimals),
	)
}

func (p *ParaswapManager) Allowance(balance float64) (*types.Transaction, error) {
	minAllowance := p.baseToken.FromFloat(balance)
	augustusAddr := common.HexToAddress(AUGUSTUS_ADDRESS)

	allowance, err := p.baseTokenContract.Allowance(nil, p.account.Address, augustusAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get allowance: %v", err)
	}

	p.logger.Info().
		Str("account", p.account.Address.String()).
		Str("router", augustusAddr.String()).
		Str("allowance", allowance.String()).
		Msg("current allowance")

	// Check if the allowance is less than the balance
	if allowance.Cmp(minAllowance) < 0 {
		p.logger.Info().Msg("increasing allowance")

		allowanceAmount := new(big.Int).Sub(abi.MaxUint256, allowance)

		gp, gasFeeCap, gasFeeTip, _ := p.account.Client.GasPriceOracle(true)
		opts := p.account.GetOpts(0, gp, gasFeeCap, gasFeeTip)

		// Increase the allowance
		tx, err := p.baseTokenContract.Approve(opts, augustusAddr, allowanceAmount)
		if err != nil {
			return nil, fmt.Errorf("failed to approve allowance: %v", err)
		}

		return tx, nil
	}

	return nil, nil
}
