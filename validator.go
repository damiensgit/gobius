package main

import (
	"context"
	"gobius/account"
	"gobius/client"
	task "gobius/common"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type IValidator interface {
	SignalCommitment(validator common.Address, taskId task.TaskId, commitment [32]byte) error
	SubmitSolution(validator common.Address, taskId task.TaskId, cid []byte) error
	SubmitIpfsCid(validator common.Address, taskId task.TaskId, cid []byte) error
	GetNextValidatorAddress() common.Address
	InitiateValidatorWithdraw(validator common.Address, amount float64) error
	ValidatorWithdraw(validator common.Address) error
	CancelValidatorWithdraw(validator common.Address, count int64) error
	BulkClaim(taskIds [][32]byte) (*types.Receipt, error)
	BatchCommitments() error
	BatchSolutions() error
	VoteOnContestation(validator common.Address, taskId task.TaskId, yeah bool) error
	SubmitContestation(validator common.Address, taskId task.TaskId) error
}

type Validator struct {
	services  *Services
	Account   *account.Account
	ratelimit float64
}

func NewValidator(services *Services, ctx context.Context, privateKey string, client *client.Client, ratelimit float64) (*Validator, error) {

	account, err := account.NewAccount(privateKey, client, ctx, services.Config.Blockchain.CacheNonce, services.Logger)
	if err != nil {
		return nil, err
	}
	account.UpdateNonce()

	va := &Validator{
		services:  services,
		Account:   account,
		ratelimit: ratelimit,
	}

	return va, nil
}

func (v *Validator) ValidatorAddress() common.Address {
	return v.Account.Address
}

func (v *Validator) InitiateValidatorWithdraw(amount float64) error {

	var validatorBal *big.Int

	if amount > 0 {
		validatorBal = v.services.Config.BaseConfig.BaseToken.FromFloat(amount)
		v.services.Logger.Info().Float64("withdraw_amount", amount).Msg("initiating validator withdraw")
	} else {
		valInfo, err := v.services.Engine.Engine.Validators(nil, v.Account.Address)
		if err != nil {
			v.services.Logger.Err(err).Msg("failed to get validator info")
			return err
		}
		validatorBal = valInfo.Staked
		bal := v.services.Config.BaseConfig.BaseToken.ToFloat(valInfo.Staked)
		v.services.Logger.Info().Float64("staked", bal).Msg("initiating validator withdraw for staked balance")
	}

	gp, gasFeeCap, gasFeeTip, _ := v.Account.Client.GasPriceOracle(true)
	opts := v.Account.GetOpts(0, gp, gasFeeCap, gasFeeTip)

	tx, err := v.services.Engine.Engine.InitiateValidatorWithdraw(opts, validatorBal)
	if err != nil {
		v.services.Logger.Err(err).Msg("failed to initiate validator withdraw")
		return err
	}

	// Wait for the transaction to be mined
	_, success, _, _ := v.Account.WaitForConfirmedTx(tx)
	if !success {
		return err
	}

	minUnlockTimeBig, err := v.services.Engine.Engine.ExitValidatorMinUnlockTime(nil)
	if err != nil {
		v.services.Logger.Err(err).Msg("error getting validator min unlock time")
		return err
	}

	secondsToUnlock := time.Duration(minUnlockTimeBig.Uint64())

	completeTimein24hrs := time.Now().Add(secondsToUnlock * time.Second)
	v.services.Logger.Info().Str("txhash", tx.Hash().String()).Msgf("initated validator withdraw of %s. complete on or after: %s", validatorBal.String(), completeTimein24hrs.Format(time.DateTime))
	return nil
}

func (v *Validator) VoteOnContestation(taskId task.TaskId, yea bool) error {
	gp, gasFeeCap, gasFeeTip, _ := v.Account.Client.GasPriceOracle(true)
	opts := v.Account.GetOpts(0, gp, gasFeeCap, gasFeeTip)

	tx, err := v.services.Engine.Engine.VoteOnContestation(opts, taskId, yea)
	if err != nil {
		v.services.Logger.Err(err).Msg("failed to vote on contestation")
		return err
	}

	// Wait for the transaction to be mined
	_, success, _, _ := v.Account.WaitForConfirmedTx(tx)
	if !success {
		return err
	}

	v.services.Logger.Info().Str("txhash", tx.Hash().String()).Msg("voted on contestation")
	return nil
}

func (v *Validator) SubmitContestation(taskId task.TaskId) error {
	gp, gasFeeCap, gasFeeTip, _ := v.Account.Client.GasPriceOracle(true)
	opts := v.Account.GetOpts(0, gp, gasFeeCap, gasFeeTip)

	tx, err := v.services.Engine.Engine.SubmitContestation(opts, taskId)
	if err != nil {
		v.services.Logger.Err(err).Msg("failed to submit contestation")
		return err
	}

	// Wait for the transaction to be mined
	_, success, _, _ := v.Account.WaitForConfirmedTx(tx)
	if !success {
		return err
	}

	v.services.Logger.Info().Str("txhash", tx.Hash().String()).Msg("submitted contestation")
	return nil
}

func (v *Validator) CancelValidatorWithdraw(index int64) error {
	// get the baseTokenBalance
	countBig, err := v.services.Engine.Engine.PendingValidatorWithdrawRequestsCount(nil, v.ValidatorAddress())
	if err != nil {
		v.services.Logger.Err(err).Msg("failed to get pending validator withdraw request count")
		return err
	}
	count := countBig.Int64()

	if index < 0 || index > count {
		v.services.Logger.Error().Msgf("invalid index %d (request count: %d)", index, count)
		return nil
	}

	currCountBig := big.NewInt(index)
	request, err := v.services.Engine.Engine.PendingValidatorWithdrawRequests(nil, v.ValidatorAddress(), currCountBig)
	if err != nil {
		v.services.Logger.Err(err).Msg("failed to get pending validator withdraw request info")
		return err
	}

	unlockUnix := request.UnlockTime.Int64()

	if unlockUnix == 0 {
		v.services.Logger.Warn().Msgf("withdraw request #%d is invalid or completed", index)
		return nil
	}

	v.services.Logger.Info().Msgf("cancelling withdraw #%d", index)

	gp, gasFeeCap, gasFeeTip, _ := v.Account.Client.GasPriceOracle(true)

	opts := v.Account.GetOpts(0, gp, gasFeeCap, gasFeeTip)
	tx, err := v.services.Engine.Engine.CancelValidatorWithdraw(opts, currCountBig)
	if err != nil {
		v.services.Logger.Err(err).Msg("failed to cancel validator withdraw")
		return err
	}

	// Wait for the transaction to be mined
	_, success, _, _ := v.Account.WaitForConfirmedTx(tx)
	if !success {
		return err
	}

	v.services.Logger.Info().Str("txhash", tx.Hash().String()).Msgf("validator withdraw #%d completed", index)

	return nil
}

func (v *Validator) ValidatorWithdraw() error {
	// get the baseTokenBalance
	countBig, err := v.services.Engine.Engine.PendingValidatorWithdrawRequestsCount(nil, v.ValidatorAddress())
	if err != nil {
		v.services.Logger.Err(err).Msg("failed to get pending validator withdraw request count")
		return err
	}
	count := countBig.Uint64()

	for i := uint64(0); i <= count; i++ {

		currCountBig := big.NewInt(int64(i))

		request, err := v.services.Engine.Engine.PendingValidatorWithdrawRequests(nil, v.ValidatorAddress(), currCountBig)
		if err != nil {
			v.services.Logger.Err(err).Msg("failed to get pending validator withdraw request info")
			return err
		}

		unlockUnix := request.UnlockTime.Int64()

		if unlockUnix == 0 {
			continue
		}

		t := time.Unix(unlockUnix, 0)
		formatted := t.Format(time.DateTime)
		valAsFromat := v.services.Config.BaseConfig.BaseToken.FormatFixed(request.Amount)
		v.services.Logger.Info().Msgf("withdraw request #%d of %s, has unlock time %s", i, valAsFromat, formatted)

		if time.Now().Unix() > unlockUnix {

			v.services.Logger.Info().Msgf("withdrawing %s aius...", valAsFromat)

			gp, gasFeeCap, gasFeeTip, _ := v.Account.Client.GasPriceOracle(true)

			opts := v.Account.GetOpts(0, gp, gasFeeCap, gasFeeTip)
			tx, err := v.services.Engine.Engine.ValidatorWithdraw(opts, currCountBig, v.ValidatorAddress())
			if err != nil {
				v.services.Logger.Err(err).Msg("failed to complete validator withdraw")
				return err
			}

			// Wait for the transaction to be mined
			_, success, _, _ := v.Account.WaitForConfirmedTx(tx)
			if !success {
				return err
			}

			v.services.Logger.Info().Str("txhash", tx.Hash().String()).Msgf("validator withdraw #%d completed", i)
		} else {
			v.services.Logger.Warn().Msgf("withdraw request #%d unlock time is in future", i)
		}
	}

	return nil
}

func (v *Validator) GetValidatorStakeBuffer() (float64, error) {
	stakedAmount, err := v.services.Engine.GetValidatorStaked(v.ValidatorAddress())
	if err != nil {
		v.services.Logger.Err(err).Err(err).Msg("could get staked amount")
		return 0, err
	}

	// get the validator minimum
	validatorMin, err := v.services.Engine.GetValidatorMinimum()
	if err != nil {
		v.services.Logger.Err(err).Err(err).Msg("could get validator min amount")
		return 0, err
	}

	// get the validator pending withdraw
	validatorPending, err := v.services.Engine.GetValidatorWithdrawPendingAmount(v.ValidatorAddress())
	if err != nil {
		v.services.Logger.Err(err).Err(err).Msg("could get validator min amount")
		return 0, err
	}

	stakedAmountFloat := v.services.Config.BaseConfig.BaseToken.ToFloat(stakedAmount)
	validatorMinFloat := v.services.Config.BaseConfig.BaseToken.ToFloat(validatorMin)
	validatorPendingFloat := v.services.Config.BaseConfig.BaseToken.ToFloat(validatorPending)

	return stakedAmountFloat - validatorPendingFloat - validatorMinFloat, nil
}

// TODO: return error
func (v *Validator) ProcessValidatorStake(baseTokenBalance *big.Int) {

	// get the baseTokenBalance
	validatorBal, err := v.services.Basetoken.BalanceOf(nil, v.ValidatorAddress())
	if err != nil {
		v.services.Logger.Err(err).Msg("failed to get balance")
		return
	}

	validatorBaseBalFloat := v.services.Config.BaseConfig.BaseToken.ToFloat(validatorBal)
	//	validatorBaseBalAsFmt := fmt.Sprintf("%.8g%s", v.services.Config.BaseConfig.BaseToken.ToFloat(validatorBal), v.services.Config.BaseConfig.BaseToken.Symbol)

	stakedAmount, err := v.services.Engine.GetValidatorStaked(v.ValidatorAddress())
	if err != nil {
		v.services.Logger.Err(err).Err(err).Msg("could get staked amount")
		return
	}

	// get the validator minimum
	validatorMin, err := v.services.Engine.GetValidatorMinimum()
	if err != nil {
		v.services.Logger.Err(err).Msg("failed to get validator minimum")
		return
	}

	// get the validator pending withdraw
	validatorPending, err := v.services.Engine.GetValidatorWithdrawPendingAmount(v.ValidatorAddress())
	if err != nil {
		v.services.Logger.Err(err).Err(err).Msg("could get validator pending withdraw amount")
		return
	}

	// moved this to earlier to ensure we have correct allowance EVEN if we dont need to deposit (fixes issue on sepolia where min stake is 0)
	// get the allowance
	allowanceAddress := v.services.Config.BaseConfig.EngineAddress

	allowance, err := v.services.Basetoken.Allowance(nil, v.services.SenderOwnerAccount.Address, allowanceAddress)
	if err != nil {
		v.services.Logger.Err(err).Msg("failed to get allowance")
		return
	}

	// FormatFixed is used because we cant represent the full amount in float64
	v.services.Logger.Debug().Msgf("allowance amount: %s", v.services.Config.BaseConfig.BaseToken.FormatFixed(allowance))

	// check if the allowance is less than the balance
	if allowance.Cmp(baseTokenBalance) < 0 {
		v.services.Logger.Info().Msgf("will need to increase allowance")

		allowanceAmount := new(big.Int).Sub(abi.MaxUint256, allowance)

		opts := v.services.SenderOwnerAccount.GetOpts(0, big.NewInt(1000000000), nil, nil)
		// increase the allowance
		tx, err := v.services.Basetoken.Approve(opts, allowanceAddress, allowanceAmount)
		if err != nil {
			v.services.Logger.Err(err).Msg("failed to approve allowance")
			return
		}
		// Wait for the transaction to be mined
		_, success, _, _ := v.services.SenderOwnerAccount.WaitForConfirmedTx(tx)
		if !success {
			return
		}

		v.services.Logger.Info().Str("txhash", tx.Hash().String()).Msgf("allowance increased")
	}

	stakedAmount.Sub(stakedAmount, validatorPending)
	stakedAmountFloat := v.services.Config.BaseConfig.BaseToken.ToFloat(stakedAmount)
	validatorMinFloat := v.services.Config.BaseConfig.BaseToken.ToFloat(validatorMin)

	v.services.Logger.Info().Float64("minimum", validatorMinFloat).Float64("basetoken_bal", validatorBaseBalFloat).Float64("staked", stakedAmountFloat).Float64("pending_withdraw", v.services.Eth.ToFloat(validatorPending)).Str("address", v.ValidatorAddress().String()).Msg("validator staked amount")

	depositAmount := new(big.Int)

	// If initial stake > 0 and the staked amount is less than the initial stake (say a fresh or emptied validator) then top up the stake to initial stake value
	if v.services.Config.ValidatorConfig.InitialStake > 0 && stakedAmountFloat < v.services.Config.ValidatorConfig.InitialStake {
		if v.services.Config.ValidatorConfig.InitialStake <= validatorMinFloat {
			v.services.Logger.Error().Msg("⚠️ initial stake amount is less than validator minimum")
			return
		}
		depositAmount = v.services.Config.BaseConfig.BaseToken.FromFloat(v.services.Config.ValidatorConfig.InitialStake - stakedAmountFloat)
	} else if v.services.Config.ValidatorConfig.StakeBufferStakeAmount > 0 {

		if stakedAmountFloat-validatorMinFloat >= v.services.Config.ValidatorConfig.StakeBufferStakeAmount {
			v.services.Logger.Info().Msg("✅ staked amount is sufficient")
			return
		} else {
			depositAmount = v.services.Config.BaseConfig.BaseToken.FromFloat(v.services.Config.ValidatorConfig.StakeBufferTopupAmount)
		}

	} else {

		// calculate the minimum with topup buffer
		minWithTopupBuffer := new(big.Int).Mul(validatorMin, big.NewInt(100))
		minWithTopupBuffer = new(big.Int).Div(minWithTopupBuffer, big.NewInt(int64(100-v.services.Config.ValidatorConfig.StakeBufferTopupPercent)))

		minWithTopupBufferAF := v.services.Config.BaseConfig.BaseToken.ToFloat(minWithTopupBuffer)

		v.services.Logger.Debug().Float64("amount", minWithTopupBufferAF).Msgf("minWithTopupBufferAF")

		// check if the staked amount is greater than or equal to the minimum with topup buffer
		if stakedAmount.Cmp(minWithTopupBuffer) >= 0 {
			v.services.Logger.Info().Msg("✅ staked amount is sufficient")
			return
		}

		// calculate the minimum with buffer
		minWithBuffer := new(big.Int).Mul(validatorMin, big.NewInt(100))
		minWithBuffer = new(big.Int).Div(minWithBuffer, big.NewInt(int64(100-v.services.Config.ValidatorConfig.StakeBufferPercent)))

		// calculate the deposit amount
		depositAmount.Sub(minWithBuffer, stakedAmount)

	}
	depositAmountAsFloat := v.services.Config.BaseConfig.BaseToken.ToFloat(depositAmount)
	v.services.Logger.Info().Float64("amount", depositAmountAsFloat).Msgf("deposit amount")

	if depositAmountAsFloat <= 0 {
		v.services.Logger.Error().Msgf("deposit amount is less than zero - check stake values in config")
		return
	}

	// check if the balance is less than the deposit amount
	if baseTokenBalance.Cmp(depositAmount) < 0 {
		v.services.Logger.Error().Msgf("⚠️ balance %g less than deposit amount %g ⚠️",
			v.services.Config.BaseConfig.BaseToken.ToFloat(baseTokenBalance),
			depositAmountAsFloat,
		)
		return
	}

	tx, err := v.services.SenderOwnerAccount.NonceManagerWrapper(3, 425, 1.5, true, func(opts *bind.TransactOpts) (interface{}, error) {
		return v.services.Engine.Engine.ValidatorDeposit(opts, v.ValidatorAddress(), depositAmount)
	})
	if err != nil {
		v.services.Logger.Err(err).Msg("failed to deposit stake")
		return
	}

	// Wait for the transaction to be mined
	_, success, _, _ := v.services.SenderOwnerAccount.WaitForConfirmedTx(tx)
	if !success {
		return
	}

	v.services.Logger.Info().Str("txhash", tx.Hash().String()).Msg("✅ deposited stake")

	stakedAmount, err = v.services.Engine.GetValidatorStaked(v.ValidatorAddress())
	if err != nil {
		v.services.Logger.Err(err).Err(err).Msg("could not get staked amount")
		return
	}

	stakedAmount.Sub(stakedAmount, validatorPending)

	v.services.Logger.Info().Float64("staked", v.services.Eth.ToFloat(stakedAmount)).Str("address", v.ValidatorAddress().String()).Msg("post staked amount")

}

func (v *Validator) CooldownTime(minClaimSolutionTime, minContestationVotePeriodTime uint64) (uint64, error) {

	lastContestationLossTimeBig, err := v.services.Engine.Engine.LastContestationLossTime(nil, v.ValidatorAddress())
	if err != nil {
		return 0, err
	}
	cooldownTime := uint64(0)
	lastContestationLossTime := lastContestationLossTimeBig.Uint64()
	if lastContestationLossTime > 0 {
		cooldownTime = lastContestationLossTime + minClaimSolutionTime + minContestationVotePeriodTime
		v.services.Logger.Debug().Str("address", v.ValidatorAddress().String()).Uint64("lastcontestationlosttime", lastContestationLossTime).Uint64("cooldowntime", cooldownTime).Msg("last contestation time")
	}
	return 0, err
}

// Return the maxmimum theoretical submissions that can be sent since the last submission from this validator
func (v *Validator) MaxSubmissions(blockTime time.Time) (time.Time, int64, error) {
	lastSubmissionBig, err := v.services.Engine.Engine.LastSolutionSubmission(nil, v.ValidatorAddress())
	if err != nil {
		v.services.Logger.Error().Err(err).Msg("LastSolutionSubmission error")
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
