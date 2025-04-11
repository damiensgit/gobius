package account

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"gobius/client"
	"gobius/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
)

type Account struct {
	privateKey *ecdsa.PrivateKey
	Client     *client.Client
	Address    common.Address
	auth       *bind.TransactOpts
	ctx        context.Context
	nonce      uint64         `json:"-"`
	cacheNonce bool           `json:"-"`
	logger     zerolog.Logger `json:"-"`
	sync.RWMutex
}

const basefeeWiggleMultiplier = 2

// Portions of code taken from offchainlabs/go-ethereum to handle trasactions

func NewAccount(privateKeyHex string, client *client.Client, ctx context.Context, cacheNonce bool, logger zerolog.Logger) (*Account, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, client.ChainID)
	if err != nil {
		log.Fatal("NewKeyedTransactorWithChainID failed with:", err)
	}

	if ctx == nil {
		ctx = context.Background()
	}

	return &Account{
		privateKey: privateKey,
		Client:     client,
		Address:    address,
		auth:       auth,
		ctx:        ctx,
		cacheNonce: cacheNonce,
		nonce:      math.MaxUint64,
		logger:     logger,
	}, nil
}

func (account *Account) estimateGasLimit(opts *bind.TransactOpts, contract *common.Address, input []byte, gasPrice, gasTipCap, gasFeeCap, value *big.Int) (uint64, error) {

	msg := ethereum.CallMsg{
		From:      opts.From,
		To:        contract,
		GasPrice:  gasPrice,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Value:     value,
		Data:      input,
	}
	gasLimit, err := account.Client.Client.EstimateGas(opts.Context, msg)
	if err != nil {
		return 0, err
	}
	// Arbitrum: adjust the estimate
	adjustedLimit := gasLimit * (10000 + opts.GasMargin) / 10000
	if adjustedLimit > gasLimit {
		gasLimit = adjustedLimit
	}
	return gasLimit, nil
}

func (account *Account) createDynamicTx(opts *bind.TransactOpts, contract *common.Address, input []byte, head *types.Header) (*types.Transaction, error) {
	// Normalize value
	value := opts.Value
	if value == nil {
		value = new(big.Int)
	}
	// Estimate TipCap
	gasTipCap := opts.GasTipCap
	if gasTipCap == nil {
		tip, err := account.Client.Client.SuggestGasTipCap(account.ctx)
		if err != nil {
			return nil, err
		}
		gasTipCap = tip
	}
	// Estimate FeeCap
	gasFeeCap := opts.GasFeeCap
	if gasFeeCap == nil {
		gasFeeCap = new(big.Int).Add(
			gasTipCap,
			new(big.Int).Mul(head.BaseFee, big.NewInt(basefeeWiggleMultiplier)),
		)
	}
	if gasFeeCap.Cmp(gasTipCap) < 0 {
		return nil, fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", gasFeeCap, gasTipCap)
	}
	// Estimate GasLimit
	gasLimit := opts.GasLimit
	if opts.GasLimit == 0 {
		var err error
		gasLimit, err = account.estimateGasLimit(opts, contract, input, nil, gasTipCap, gasFeeCap, value)
		if err != nil {
			return nil, err
		}
	}

	nonce := uint64(0)
	if opts.Nonce == nil {
		var err error

		nonce, err = account.PendingNonce()
		if err != nil {
			return nil, err
		}
	} else {
		nonce = opts.Nonce.Uint64()
	}

	// create the transaction
	baseTx := &types.DynamicFeeTx{
		To:        contract,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       gasLimit,
		Value:     value,
		Data:      input,
	}
	return types.NewTx(baseTx), nil
}

func (account *Account) createLegacyTx(opts *bind.TransactOpts, contract *common.Address, input []byte) (*types.Transaction, error) {
	if opts.GasFeeCap != nil || opts.GasTipCap != nil {
		return nil, errors.New("maxFeePerGas or maxPriorityFeePerGas specified but london is not active yet")
	}
	// Normalize value
	value := opts.Value
	if value == nil {
		value = new(big.Int)
	}
	// Estimate GasPrice
	gasPrice := opts.GasPrice
	if gasPrice == nil {
		price, err := account.Client.Client.SuggestGasPrice(account.ctx)
		if err != nil {
			return nil, err
		}
		gasPrice = price
	}
	// Estimate GasLimit
	gasLimit := opts.GasLimit
	if opts.GasLimit == 0 {
		var err error
		gasLimit, err = account.estimateGasLimit(opts, contract, input, gasPrice, nil, nil, value)
		if err != nil {
			return nil, err
		}
	}

	nonce := uint64(0)
	if opts.Nonce == nil {
		var err error

		nonce, err = account.PendingNonce()
		if err != nil {
			return nil, err
		}
	} else {
		nonce = opts.Nonce.Uint64()
	}

	// create the transaction
	baseTx := &types.LegacyTx{
		To:       contract,
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		Value:    value,
		Data:     input,
	}
	return types.NewTx(baseTx), nil
}

// transact executes an actual transaction invocation, first deriving any missing
// authorization fields, and then scheduling the transaction for execution.
func (account *Account) SendTransactionWithOpts(opts *bind.TransactOpts, contract *common.Address, input []byte) (*types.Transaction, error) {
	if opts.GasPrice != nil && (opts.GasFeeCap != nil || opts.GasTipCap != nil) {
		return nil, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	}
	// Create the transaction
	var (
		rawTx *types.Transaction
		err   error
	)
	if opts.GasPrice != nil {
		rawTx, err = account.createLegacyTx(opts, contract, input)
	} else if opts.GasFeeCap != nil && opts.GasTipCap != nil {
		rawTx, err = account.createDynamicTx(opts, contract, input, nil)
	} else {
		// Only query for basefee if gasPrice not specified
		if head, errHead := account.Client.Client.HeaderByNumber(opts.Context, nil); errHead != nil {
			return nil, errHead
		} else if head.BaseFee != nil {
			rawTx, err = account.createDynamicTx(opts, contract, input, head)
		} else {
			// Chain is not London ready -> use legacy transaction
			rawTx, err = account.createLegacyTx(opts, contract, input)
		}
	}
	if err != nil {
		return nil, err
	}
	// Sign the transaction and schedule it for execution
	if opts.Signer == nil {
		return nil, errors.New("no signer to authorize the transaction with")
	}
	signedTx, err := opts.Signer(opts.From, rawTx)
	if err != nil {
		return nil, err
	}
	if opts.NoSend {
		return signedTx, nil
	}
	err = account.Client.Client.SendTransaction(account.ctx, signedTx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func (account *Account) SendSignedTransaction(signedTx *types.Transaction) (*types.Transaction, error) {
	err := account.Client.Client.SendTransaction(account.ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (account *Account) GetBalance() (*big.Int, error) {
	// get the Ether balance
	etherBalance, err := account.Client.Client.BalanceAt(account.ctx, account.Address, nil)
	if err != nil {
		return nil, err
	}

	return etherBalance, nil
}

func (account *Account) Nonce() uint64 {
	return atomic.LoadUint64(&account.nonce)
}

func (account *Account) CacheNonce() bool {
	return account.cacheNonce
}

func (account *Account) SetNonce(newValue uint64) {
	atomic.StoreUint64(&account.nonce, newValue)
}

func (account *Account) IncNonce() {
	atomic.AddUint64(&account.nonce, 1)
}

func (account *Account) DecNonce() {
	atomic.AddUint64(&account.nonce, ^uint64(0))
}

// returns a new auth object with the nonce set to the current value
func (account *Account) GetOptsWithoutNonceInc(limit uint64, gasPrice, gasFeeCap, gasTipCap *big.Int) *bind.TransactOpts {
	account.Lock()
	defer account.Unlock()

	auth := new(bind.TransactOpts)
	*auth = *account.auth
	auth.Context = account.ctx

	// if nonce is at max, then use nil to indicate that the nonce should be determined by the client
	if account.cacheNonce {
		nonce := account.Nonce()
		auth.Nonce = new(big.Int).SetUint64(nonce)
	} else {
		auth.Nonce = nil
	}

	//auth.Value = big.NewInt(0)
	auth.GasLimit = limit
	if gasPrice != nil {
		auth.GasPrice = gasPrice
	} else if gasFeeCap != nil && gasTipCap != nil {
		auth.GasPrice = nil
		auth.GasFeeCap = gasFeeCap
		auth.GasTipCap = gasTipCap
	}

	return auth
}

// returns a new auth object with the nonce set to the current value
func (account *Account) GetOpts(limit uint64, gasPrice, gasFeeCap, gasTipCap *big.Int) *bind.TransactOpts {
	auth := account.GetOptsWithoutNonceInc(limit, gasPrice, gasFeeCap, gasTipCap)

	if account.cacheNonce {
		account.IncNonce()
	}

	return auth
}

func (account *Account) PendingNonce() (uint64, error) {
	nonce, err := account.Client.Client.PendingNonceAt(account.ctx, account.Address)
	if err != nil {
		return 0, err
	}

	return nonce, nil
}

func (account *Account) UpdateNonce() error {

	account.Lock()
	defer account.Unlock()

	// if account.cacheNonce && account.Nonce() != math.MaxUint64 {
	// 	return nil
	// }

	//start := time.Now()
	nonce, err := account.PendingNonce()

	//duration := time.Since(start)
	//zap.S().Debug("Using nonce:", noncet, "(took: ", duration, ")")
	if err != nil {
		//zap.S().Error("UpdateNonce: couldn't fetch pending nonce for account", err)
		return err
	}

	account.SetNonce(nonce)

	return nil
}

// TODO: move logger to account type
// TODO: use app context
func (account *Account) WaitForConfirmedTx(tx *types.Transaction) (receipt *types.Receipt, success bool, revertReason string, err error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	receipt, err = bind.WaitMined(ctxTimeout, account.Client.Client, tx)
	if err != nil {
		account.logger.Error().Err(err).Str("txhash", tx.Hash().String()).Msg("❌ error waiting for transaction to be mined")
		return nil, false, "", err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		revertReason = utils.GetRevertReason(account.Address, tx, account.Client.Client, receipt.BlockNumber)
		account.logger.Error().
			Str("reason", revertReason).
			Str("txhash", tx.Hash().String()).
			Str("	", tx.GasFeeCap().String()).
			Str("gasfeetip", tx.GasTipCap().String()).
			Msg("❌ transaction was not successful")
		return receipt, false, revertReason, nil
	}
	return receipt, true, "", nil
}

func (account *Account) NonceManagerWrapperWithContext(ctx context.Context, opts *bind.TransactOpts, tries int, base, backoffMultiplier float64, override bool, fn func(opts *bind.TransactOpts) (interface{}, error)) (*types.Transaction, error) {
	if opts == nil {
		gasPrice, gasFeeCap, gasFeeTip, _ := account.Client.GasPriceOracle(override)
		opts = account.GetOpts(0, gasPrice, gasFeeCap, gasFeeTip)
	}
	txFunc := func(nonce uint64) (any, error) {

		if nonce > 0 && account.CacheNonce() {
			// nonceMu.Lock()

			// value, found := cache.Get(nonce)
			// if found {
			// 	nonceContention = value.(uint64) + 1
			// }
			// cache.Add(nonce, nonceContention)

			// nonceMu.Unlock()

			//if opts.Nonce != nil {
			optsNonce := opts.Nonce.Uint64()
			diff := int64(nonce - optsNonce)

			//log.Printf("NONCE HANDLING: STATE: %d OPTS: %d DIFF: %d\n", nonce, optsNonce, diff)
			//}

			//newNonce := nonce

			//optsNonce = nonce
			if diff > 0 {
				// nonce too low just set to state
				optsNonce = nonce
			} else {
				// nonce too high
				optsNonce--
			}
			//account.UpdateNonce()

			opts.Nonce = big.NewInt(int64(optsNonce))
			account.SetNonce(optsNonce + 1)
			// if diff < 0 {
			// 	newNonce

			// }

			// if opts.Nonce.Cmp(big.NewInt(int64(nonce))) <= 0 {
			// 	opts.Nonce = big.NewInt(int64(nonce))
			// } else {
			// 	//account.SetNonce(nonce + 1)
			// 	opts.Nonce.Sub(opts.Nonce, big.NewInt(1))
			// }
			//log.Printf("NONCE HANDLING: AFTER: %d", opts.Nonce.Int64())
			//opts.Nonce = big.NewInt(int64(nonce))
		}
		return fn(opts)
	}
	result, err := utils.ExpRetryWithNonceContext(ctx, account.logger, txFunc, tries, base, backoffMultiplier)
	txResult, ok := result.(*types.Transaction)
	if !ok {
		return nil, errors.New("result is not the expected type")
	}
	if err != nil {
		account.UpdateNonce()
		//		account.DecNonce()
	}
	return txResult, err
}

func (account *Account) NonceManagerWrapper(tries int, base, backoffMultiplier float64, override bool, fn func(opts *bind.TransactOpts) (interface{}, error)) (*types.Transaction, error) {
	return account.NonceManagerWrapperWithContext(context.Background(), nil, tries, base, backoffMultiplier, override, fn)
}

func (account *Account) SendEther(opts *bind.TransactOpts, toAddress common.Address, value *big.Int) (*types.Transaction, error) {

	if opts == nil {
		// just use network defaults (estimate gas, etc)
		opts = account.GetOpts(0, nil, nil, nil)
	}
	// Ensure value is set in the options
	if opts.Value == nil {
		opts.Value = value
	} else if opts.Value.Cmp(value) != 0 {
		account.logger.Warn().Str("opts_value", opts.Value.String()).Str("param_value", value.String()).Msg("opts.Value already set, overriding with parameter value")
		opts.Value = value // Override if already set but different, log a warning
	}

	// Send the transaction
	return account.SendTransactionWithOpts(opts, &toAddress, nil)
}

// ensureContext is a helper method to ensure a context is not nil, even if the
// user specified it as such.
func (account *Account) ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return account.ctx
	}
	return ctx
}
