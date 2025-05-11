package main

import (
	"context"
	"errors"
	"gobius/account"
	"gobius/bindings/basetoken"
	"gobius/client"
	task "gobius/common"
	"gobius/config"
	"gobius/erc20"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
)

//	type IValidator interface {
//		SignalCommitment(validator common.Address, taskId task.TaskId, commitment [32]byte) error
//		SubmitSolution(validator common.Address, taskId task.TaskId, cid []byte) error
//		SubmitIpfsCid(validator common.Address, taskId task.TaskId, cid []byte) error
//		GetNextValidatorAddress() common.Address
//		InitiateValidatorWithdraw(validator common.Address, amount float64) error
//		ValidatorWithdraw(validator common.Address) error
//		CancelValidatorWithdraw(validator common.Address, count int64) error
//		BulkClaim(taskIds [][32]byte) (*types.Receipt, error)
//		BatchCommitments() error
//		BatchSolutions() error
//		VoteOnContestation(validator common.Address, taskId task.TaskId, yeah bool) error
//		SubmitContestation(validator common.Address, taskId task.TaskId) error

type Validators struct {
	validators []*Validator
	index      int
	mu         sync.Mutex
}

var (
	ValidatorNotFound = errors.New("validator not found")
)

type Validator struct {
	Account   *account.Account
	ratelimit float64
	logger    zerolog.Logger
	config    *config.AppConfig
	engine    *EngineWrapper
	basetoken *basetoken.BaseToken
}

var Eth = erc20.NewTokenERC20(common.HexToAddress("0x0"), 18, "ETH", "ETH")

func NewValidator(config *config.AppConfig, engine *EngineWrapper, baseToken *basetoken.BaseToken, logger zerolog.Logger, ctx context.Context, privateKey string, client *client.Client, ratelimit float64) (*Validator, error) {

	account, err := account.NewAccount(privateKey, client, ctx, config.Blockchain.CacheNonce, logger)
	if err != nil {
		return nil, err
	}
	// TODO: this should be done in the account constructor
	account.UpdateNonce()

	validatorLogger := logger.With().Str("validator", account.Address.String()).Logger()

	va := &Validator{
		Account:   account,
		ratelimit: ratelimit,
		logger:    validatorLogger,
		config:    config,
		engine:    engine,
		basetoken: baseToken,
	}

	return va, nil
}

func (v *Validator) ValidatorAddress() common.Address {
	return v.Account.Address
}

// VoteOnContestation
func (v *Validator) VoteOnContestation(taskId task.TaskId, yea bool) error {
	tx, err := v.Account.NonceManagerWrapper(v.config.Miner.ErrorMaxRetries, v.config.Miner.ErrorBackoffTime, v.config.Miner.ErrorBackofMultiplier, false, func(opts *bind.TransactOpts) (interface{}, error) {
		opts.NoSend = true // Prepare transaction but do not send automatically via wrapper
		txToSign, err := v.engine.Engine.VoteOnContestation(opts, taskId, yea)
		if err != nil {
			return nil, err
		}
		return v.Account.SendSignedTransaction(txToSign)
	})

	if err != nil {
		v.logger.Error().Err(err).Msg("error sending contestation vote")
		return err
	}

	if tx == nil {
		err = errors.New("assertion: transaction is nil but no error reported from NonceManagerWrapper")
		v.logger.Error().Err(err).Msg("error sending contestation vote")
		return err
	}

	// Wait for the transaction to be mined
	_, success, _, err := v.Account.WaitForConfirmedTx(tx)

	if err != nil {
		v.logger.Error().Err(err).Msg("error waiting for contestation vote to be mined")
		return err
	}

	if success {
		v.logger.Info().Str("txhash", tx.Hash().String()).Msg("✅ contestation vote successful")
	} else {
		v.logger.Error().Msg("❌ contestation vote failed")
	}

	return nil
}

func (v *Validator) SubmitContestation(taskId task.TaskId) error {
	tx, err := v.Account.NonceManagerWrapper(v.config.Miner.ErrorMaxRetries, v.config.Miner.ErrorBackoffTime, v.config.Miner.ErrorBackofMultiplier, false, func(opts *bind.TransactOpts) (interface{}, error) {
		opts.NoSend = true // Prepare transaction but do not send automatically via wrapper
		txToSign, err := v.engine.Engine.SubmitContestation(opts, taskId)
		if err != nil {
			return nil, err
		}
		return v.Account.SendSignedTransaction(txToSign)
	})

	if err != nil {
		v.logger.Error().Err(err).Msg("error sending contestation")
		return err
	}

	// Wait for the transaction to be mined
	_, success, _, err := v.Account.WaitForConfirmedTx(tx)

	if err != nil {
		v.logger.Error().Err(err).Msg("error waiting for contestation to be mined")
		return err
	}

	if success {
		v.logger.Info().Str("txhash", tx.Hash().String()).Msg("✅ contestation successful")
	} else {
		v.logger.Error().Msg("❌ contestation failed")
	}

	return nil
}

func (v *Validator) InitiateValidatorWithdraw(amount float64) error {

	var validatorBal *big.Int

	if amount > 0 {
		validatorBal = v.config.BaseConfig.BaseToken.FromFloat(amount)
		v.logger.Info().Float64("withdraw_amount", amount).Msg("initiating validator withdraw")
	} else {
		valInfo, err := v.engine.Engine.Validators(nil, v.Account.Address)
		if err != nil {
			v.logger.Err(err).Msg("failed to get validator info")
			return err
		}
		validatorBal = valInfo.Staked
		bal := v.config.BaseConfig.BaseToken.ToFloat(valInfo.Staked)
		v.logger.Info().Float64("staked", bal).Msg("initiating validator withdraw for staked balance")
	}

	gp, gasFeeCap, gasFeeTip, _ := v.Account.Client.GasPriceOracle(true)
	opts := v.Account.GetOpts(0, gp, gasFeeCap, gasFeeTip)

	tx, err := v.engine.Engine.InitiateValidatorWithdraw(opts, validatorBal)
	if err != nil {
		v.logger.Err(err).Msg("failed to initiate validator withdraw")
		return err
	}

	// Wait for the transaction to be mined
	_, success, _, _ := v.Account.WaitForConfirmedTx(tx)
	if !success {
		return err
	}

	minUnlockTimeBig, err := v.engine.Engine.ExitValidatorMinUnlockTime(nil)
	if err != nil {
		v.logger.Err(err).Msg("error getting validator min unlock time")
		return err
	}

	secondsToUnlock := time.Duration(minUnlockTimeBig.Uint64())

	completeTimein24hrs := time.Now().Add(secondsToUnlock * time.Second)
	v.logger.Info().Str("txhash", tx.Hash().String()).Msgf("initated validator withdraw of %s. complete on or after: %s", validatorBal.String(), completeTimein24hrs.Format(time.DateTime))
	return nil
}

func (v *Validator) CancelValidatorWithdraw(index int64) error {
	// get the baseTokenBalance
	countBig, err := v.engine.Engine.PendingValidatorWithdrawRequestsCount(nil, v.ValidatorAddress())
	if err != nil {
		v.logger.Err(err).Msg("failed to get pending validator withdraw request count")
		return err
	}
	count := countBig.Int64()

	if index < 0 || index > count {
		v.logger.Error().Msgf("invalid index %d (request count: %d)", index, count)
		return nil
	}

	currCountBig := big.NewInt(index)
	request, err := v.engine.Engine.PendingValidatorWithdrawRequests(nil, v.ValidatorAddress(), currCountBig)
	if err != nil {
		v.logger.Err(err).Msg("failed to get pending validator withdraw request info")
		return err
	}

	unlockUnix := request.UnlockTime.Int64()

	if unlockUnix == 0 {
		v.logger.Warn().Msgf("withdraw request #%d is invalid or completed", index)
		return nil
	}

	v.logger.Info().Msgf("cancelling withdraw #%d", index)

	gp, gasFeeCap, gasFeeTip, _ := v.Account.Client.GasPriceOracle(true)

	opts := v.Account.GetOpts(0, gp, gasFeeCap, gasFeeTip)
	tx, err := v.engine.Engine.CancelValidatorWithdraw(opts, currCountBig)
	if err != nil {
		v.logger.Err(err).Msg("failed to cancel validator withdraw")
		return err
	}

	// Wait for the transaction to be mined
	_, success, _, _ := v.Account.WaitForConfirmedTx(tx)
	if !success {
		return err
	}

	v.logger.Info().Str("txhash", tx.Hash().String()).Msgf("validator withdraw #%d completed", index)

	return nil
}

func (v *Validator) ValidatorWithdraw() error {
	// get the baseTokenBalance
	countBig, err := v.engine.Engine.PendingValidatorWithdrawRequestsCount(nil, v.ValidatorAddress())
	if err != nil {
		v.logger.Err(err).Msg("failed to get pending validator withdraw request count")
		return err
	}
	count := countBig.Uint64()

	for i := uint64(0); i <= count; i++ {

		currCountBig := big.NewInt(int64(i))

		request, err := v.engine.Engine.PendingValidatorWithdrawRequests(nil, v.ValidatorAddress(), currCountBig)
		if err != nil {
			v.logger.Err(err).Msg("failed to get pending validator withdraw request info")
			return err
		}

		unlockUnix := request.UnlockTime.Int64()

		if unlockUnix == 0 {
			continue
		}

		t := time.Unix(unlockUnix, 0)
		formatted := t.Format(time.DateTime)
		valAsFromat := v.config.BaseConfig.BaseToken.FormatFixed(request.Amount)
		v.logger.Info().Msgf("withdraw request #%d of %s, has unlock time %s", i, valAsFromat, formatted)

		if time.Now().Unix() > unlockUnix {

			v.logger.Info().Msgf("withdrawing %s aius...", valAsFromat)

			gp, gasFeeCap, gasFeeTip, _ := v.Account.Client.GasPriceOracle(true)

			opts := v.Account.GetOpts(0, gp, gasFeeCap, gasFeeTip)
			tx, err := v.engine.Engine.ValidatorWithdraw(opts, currCountBig, v.ValidatorAddress())
			if err != nil {
				v.logger.Err(err).Msg("failed to complete validator withdraw")
				return err
			}

			// Wait for the transaction to be mined
			_, success, _, _ := v.Account.WaitForConfirmedTx(tx)
			if !success {
				return err
			}

			v.logger.Info().Str("txhash", tx.Hash().String()).Msgf("validator withdraw #%d completed", i)
		} else {
			v.logger.Warn().Msgf("withdraw request #%d unlock time is in future", i)
		}
	}

	return nil
}

func (v *Validator) GetValidatorStakeBuffer() (float64, error) {
	stakedAmount, err := v.engine.GetValidatorStaked(v.ValidatorAddress())
	if err != nil {
		v.logger.Err(err).Err(err).Msg("could get staked amount")
		return 0, err
	}

	// get the validator minimum
	validatorMin, err := v.engine.GetValidatorMinimum()
	if err != nil {
		v.logger.Err(err).Err(err).Msg("could get validator min amount")
		return 0, err
	}

	// get the validator pending withdraw
	validatorPending, err := v.engine.GetValidatorWithdrawPendingAmount(v.ValidatorAddress())
	if err != nil {
		v.logger.Err(err).Err(err).Msg("could get validator min amount")
		return 0, err
	}

	stakedAmountFloat := v.config.BaseConfig.BaseToken.ToFloat(stakedAmount)
	validatorMinFloat := v.config.BaseConfig.BaseToken.ToFloat(validatorMin)
	validatorPendingFloat := v.config.BaseConfig.BaseToken.ToFloat(validatorPending)

	return stakedAmountFloat - validatorPendingFloat - validatorMinFloat, nil
}

// TODO: return error
func (v *Validator) ProcessValidatorStake(baseTokenBalance *big.Int) {

	// get the baseTokenBalance
	validatorBal, err := v.basetoken.BalanceOf(nil, v.ValidatorAddress())
	if err != nil {
		v.logger.Err(err).Msg("failed to get balance")
		return
	}

	validatorBaseBalFloat := v.config.BaseConfig.BaseToken.ToFloat(validatorBal)

	stakedAmount, err := v.engine.GetValidatorStaked(v.ValidatorAddress())
	if err != nil {
		v.logger.Err(err).Err(err).Msg("could get staked amount")
		return
	}

	// get the validator minimum
	validatorMin, err := v.engine.GetValidatorMinimum()
	if err != nil {
		v.logger.Err(err).Msg("failed to get validator minimum")
		return
	}

	// get the validator pending withdraw
	validatorPending, err := v.engine.GetValidatorWithdrawPendingAmount(v.ValidatorAddress())
	if err != nil {
		v.logger.Err(err).Err(err).Msg("could get validator pending withdraw amount")
		return
	}

	// moved this to earlier to ensure we have correct allowance EVEN if we dont need to deposit (fixes issue on sepolia where min stake is 0)
	// get the allowance
	allowanceAddress := v.config.BaseConfig.EngineAddress

	allowance, err := v.basetoken.Allowance(nil, v.Account.Address, allowanceAddress)
	if err != nil {
		v.logger.Err(err).Msg("failed to get allowance")
		return
	}

	// FormatFixed is used because we cant represent the full amount in float64
	v.logger.Debug().Msgf("allowance amount: %s", v.config.BaseConfig.BaseToken.FormatFixed(allowance))

	// check if the allowance is less than the balance
	if allowance.Cmp(baseTokenBalance) < 0 {
		v.logger.Info().Msgf("will need to increase allowance")

		allowanceAmount := new(big.Int).Sub(abi.MaxUint256, allowance)

		opts := v.Account.GetOpts(0, nil, nil, nil)
		// increase the allowance
		tx, err := v.basetoken.Approve(opts, allowanceAddress, allowanceAmount)
		if err != nil {
			v.logger.Err(err).Msg("failed to approve allowance")
			return
		}
		// Wait for the transaction to be mined
		_, success, _, _ := v.Account.WaitForConfirmedTx(tx)
		if !success {
			return
		}

		v.logger.Info().Str("txhash", tx.Hash().String()).Msgf("allowance increased")
	}

	stakedAmount.Sub(stakedAmount, validatorPending)
	stakedAmountFloat := v.config.BaseConfig.BaseToken.ToFloat(stakedAmount)
	validatorMinFloat := v.config.BaseConfig.BaseToken.ToFloat(validatorMin)

	v.logger.Info().Float64("minimum", validatorMinFloat).Float64("basetoken_bal", validatorBaseBalFloat).Float64("staked", stakedAmountFloat).Float64("pending_withdraw", Eth.ToFloat(validatorPending)).Msg("validator staked amount")

	depositAmount := new(big.Int)

	// If initial stake > 0 and the staked amount is less than the initial stake (say a fresh or emptied validator) then top up the stake to initial stake value
	if v.config.ValidatorConfig.InitialStake > 0 && stakedAmountFloat < v.config.ValidatorConfig.InitialStake {
		if v.config.ValidatorConfig.InitialStake <= validatorMinFloat {
			v.logger.Error().Msg("⚠️ initial stake amount is less than validator minimum")
			return
		}
		depositAmount = v.config.BaseConfig.BaseToken.FromFloat(v.config.ValidatorConfig.InitialStake - stakedAmountFloat)
	} else if v.config.ValidatorConfig.StakeBufferStakeAmount > 0 {

		if stakedAmountFloat-validatorMinFloat >= v.config.ValidatorConfig.StakeBufferStakeAmount {
			v.logger.Info().Msg("✅ staked amount is sufficient")
			return
		} else {
			depositAmount = v.config.BaseConfig.BaseToken.FromFloat(v.config.ValidatorConfig.StakeBufferTopupAmount)
		}

	} else {

		// calculate the minimum with topup buffer
		minWithTopupBuffer := new(big.Int).Mul(validatorMin, big.NewInt(100))
		minWithTopupBuffer = new(big.Int).Div(minWithTopupBuffer, big.NewInt(int64(100-v.config.ValidatorConfig.StakeBufferTopupPercent)))

		minWithTopupBufferAF := v.config.BaseConfig.BaseToken.ToFloat(minWithTopupBuffer)

		v.logger.Debug().Float64("amount", minWithTopupBufferAF).Msgf("minWithTopupBufferAF")

		// check if the staked amount is greater than or equal to the minimum with topup buffer
		if stakedAmount.Cmp(minWithTopupBuffer) >= 0 {
			v.logger.Info().Msg("✅ staked amount is sufficient")
			return
		}

		// calculate the minimum with buffer
		minWithBuffer := new(big.Int).Mul(validatorMin, big.NewInt(100))
		minWithBuffer = new(big.Int).Div(minWithBuffer, big.NewInt(int64(100-v.config.ValidatorConfig.StakeBufferPercent)))

		// calculate the deposit amount
		depositAmount.Sub(minWithBuffer, stakedAmount)

	}
	depositAmountAsFloat := v.config.BaseConfig.BaseToken.ToFloat(depositAmount)
	v.logger.Info().Float64("amount", depositAmountAsFloat).Msgf("deposit amount")

	if depositAmountAsFloat <= 0 {
		v.logger.Error().Msgf("deposit amount is less than zero - check stake values in config")
		return
	}

	// check if the balance is less than the deposit amount
	if baseTokenBalance.Cmp(depositAmount) < 0 {
		v.logger.Error().Msgf("⚠️ balance %g less than deposit amount %g ⚠️",
			v.config.BaseConfig.BaseToken.ToFloat(baseTokenBalance),
			depositAmountAsFloat,
		)
		return
	}

	tx, err := v.Account.NonceManagerWrapper(3, 425, 1.5, true, func(opts *bind.TransactOpts) (interface{}, error) {
		return v.engine.Engine.ValidatorDeposit(opts, v.ValidatorAddress(), depositAmount)
	})
	if err != nil {
		v.logger.Err(err).Msg("failed to deposit stake")
		return
	}

	// Wait for the transaction to be mined
	_, success, _, _ := v.Account.WaitForConfirmedTx(tx)
	if !success {
		return
	}

	v.logger.Info().Str("txhash", tx.Hash().String()).Msg("✅ deposited stake")

	stakedAmount, err = v.engine.GetValidatorStaked(v.ValidatorAddress())
	if err != nil {
		v.logger.Err(err).Err(err).Msg("could not get staked amount")
		return
	}

	stakedAmount.Sub(stakedAmount, validatorPending)

	v.logger.Info().Float64("staked", Eth.ToFloat(stakedAmount)).Str("address", v.ValidatorAddress().String()).Msg("post staked amount")

}

func (v *Validator) CooldownTime(minClaimSolutionTime, minContestationVotePeriodTime uint64) (uint64, error) {

	lastContestationLossTimeBig, err := v.engine.Engine.LastContestationLossTime(nil, v.ValidatorAddress())
	if err != nil {
		return 0, err
	}
	cooldownTime := uint64(0)
	lastContestationLossTime := lastContestationLossTimeBig.Uint64()
	if lastContestationLossTime > 0 {
		cooldownTime = lastContestationLossTime + minClaimSolutionTime + minContestationVotePeriodTime
		v.logger.Debug().Str("address", v.ValidatorAddress().String()).Uint64("lastcontestationlosttime", lastContestationLossTime).Uint64("cooldowntime", cooldownTime).Msg("last contestation time")
	}
	return 0, err
}

// Return the maxmimum theoretical submissions that can be sent since the last submission from this validator
func (v *Validator) MaxSubmissions(blockTime time.Time) (time.Time, int64, error) {
	lastSubmissionBig, err := v.engine.Engine.LastSolutionSubmission(nil, v.ValidatorAddress())
	if err != nil {
		v.logger.Error().Err(err).Msg("LastSolutionSubmission error")
		// Return zero values and the error if fetching fails
		return time.Time{}, 0, err
	}

	lastSubmission := time.Unix(lastSubmissionBig.Int64(), 0)

	diff := blockTime.Sub(lastSubmission)

	// 10 seconds since last submission, rate limit is 2 seconds between so max sols: 10/2=5
	// 10 seconds since last submission, rate limit is 0.5 seconds between so max sols: 10/0.5=20
	maxSubmissions := int64(math.Floor(diff.Seconds() / v.ratelimit))
	return lastSubmission, maxSubmissions, nil
}

// IsEligibleToVote checks if the validator is eligible to vote on a given task's contestation.
func (v *Validator) IsEligibleToVote(ctx context.Context, taskId task.TaskId) (bool, uint64, error) {
	// Note: The context passed here isn't directly used in the EngineWrapper's callOpts for ValidatorCanVote,
	// as it's a view function, but kept for potential future use or consistency.
	if v.engine == nil {
		return false, 99, errors.New("validator has nil engine wrapper")
	}
	return v.engine.IsValidatorEligibleToVote(v.ValidatorAddress(), taskId)
}

// Methods for Validators slice

func (v *Validators) GetValidatorByAddress(addr common.Address) *Validator {
	for _, v := range v.validators {
		if v.ValidatorAddress() == addr {
			return v
		}
	}
	return nil
}

func (v *Validators) InitiateValidatorWithdraw(validator common.Address, amount float64) error {
	val := v.GetValidatorByAddress(validator)
	if val != nil {
		return val.InitiateValidatorWithdraw(amount)
	}
	return ValidatorNotFound
}

func (v *Validators) ValidatorWithdraw(validator common.Address) error {
	val := v.GetValidatorByAddress(validator)
	if val != nil {
		return val.ValidatorWithdraw()
	}
	return ValidatorNotFound
}

func (v *Validators) CancelValidatorWithdraw(validator common.Address, count int64) error {
	val := v.GetValidatorByAddress(validator)
	if val != nil {
		return val.CancelValidatorWithdraw(count)
	}
	return ValidatorNotFound
}

func (v *Validators) VoteOnContestation(validator common.Address, taskId task.TaskId, yeah bool) error {
	val := v.GetValidatorByAddress(validator)
	if val != nil {
		return val.VoteOnContestation(taskId, yeah)
	}
	return ValidatorNotFound
}

func (v *Validators) SubmitContestation(validator common.Address, taskId task.TaskId) error {
	val := v.GetValidatorByAddress(validator)
	if val != nil {
		return val.SubmitContestation(taskId)
	}
	return ValidatorNotFound
}

func (v *Validators) IsAddressValidator(address common.Address) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	for _, v := range v.validators {
		if v.ValidatorAddress() == address {
			return true
		}
	}

	return false
}

func (v *Validators) GetNextValidatorAddress() common.Address {
	v.mu.Lock()
	defer v.mu.Unlock()

	if len(v.validators) == 0 {
		return common.Address{}
	}

	validator := v.validators[v.index]
	v.index++
	if v.index >= len(v.validators) {
		v.index = 0
	}

	return validator.ValidatorAddress()
}
