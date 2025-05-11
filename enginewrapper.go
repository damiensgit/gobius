package main

import (
	"errors"
	"fmt"
	"gobius/bindings/engine"
	"gobius/bindings/voter"
	task "gobius/common"
	"gobius/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
)

// engine wrapper for the engine contract for view functions only
// do not add transaction functions to this wrapper

var (
	// as per the engine contract V5_1, calculation is (reward * modelRate * gaugeMultiplier) / (2 * 1e18 * 1e18);
	denom1e18         = big.NewInt(1e18)
	denom2e36         = new(big.Int).Mul(denom1e18, denom1e18)     // 1e36
	RewardDenominator = new(big.Int).Mul(big.NewInt(2), denom2e36) // 2 * 1e36
)

type EngineWrapper struct {
	Engine *engine.Engine
	Voter  *voter.Voter
	logger zerolog.Logger
}

func NewEngineWrapper(engine *engine.Engine, voter *voter.Voter, logger zerolog.Logger) *EngineWrapper {
	return &EngineWrapper{
		Engine: engine,
		Voter:  voter,
		logger: logger,
	}
}

func (m *EngineWrapper) GetSolution(taskId task.TaskId) (*struct {
	Validator common.Address
	Blocktime uint64
	Claimed   bool
	Cid       []byte
}, error) {

	getSol := func() (interface{}, error) {
		return m.Engine.Solutions(nil, taskId)
	}

	result, err := utils.ExpRetry(m.logger, getSol, 3, 1000)
	if err != nil {
		return nil, err
	}

	res, ok := result.(struct {
		Validator common.Address
		Blocktime uint64
		Claimed   bool
		Cid       []byte
	})

	// Panic here because this is a critical error and there is no recovery
	// we need to update the code to handle this error
	if !ok {
		panic("Result is not of the expected type from engine.Solutions")
	}

	return &res, nil
}

func (m *EngineWrapper) LookupTask(taskId task.TaskId) (*struct {
	Model     [32]byte
	Fee       *big.Int
	Owner     common.Address
	Blocktime uint64
	Version   uint8
	Cid       []byte
}, error) {

	getTask := func() (interface{}, error) {
		return m.Engine.Tasks(nil, taskId)
	}

	result, err := utils.ExpRetry(m.logger, getTask, 3, 1000)
	if err != nil {
		return nil, err
	}

	res, ok := result.(struct {
		Model     [32]byte
		Fee       *big.Int
		Owner     common.Address
		Blocktime uint64
		Version   uint8
		Cid       []byte
	})

	// Panic here because this is a critical error and there is no recovery
	// we need to update the code to handle this error
	if !ok {
		panic("Result is not of the expected type from engine.Tasks")
	}

	m.logger.Debug().Str("owner", res.Owner.String()).Str("cid", common.Bytes2Hex(res.Cid)).Str("model", common.Bytes2Hex(res.Model[:])).Msg("task information")

	return &res, nil
}

func (m *EngineWrapper) GetValidatorWithdrawPendingAmount(validator common.Address) (*big.Int, error) {

	check := func() (interface{}, error) {
		return m.Engine.ValidatorWithdrawPendingAmount(nil, validator)
	}

	result, err := utils.ExpRetry(m.logger, check, 3, 1000)

	if err != nil {
		return nil, err
	}

	return result.(*big.Int), nil
}

func (m *EngineWrapper) GetValidatorStaked(validator common.Address) (*big.Int, error) {

	check := func() (interface{}, error) {
		return m.Engine.Validators(nil, validator)
	}

	result, err := utils.ExpRetry(m.logger, check, 3, 1000)

	if err != nil {
		return nil, err
	}

	validators, ok := result.(struct {
		Staked *big.Int
		Since  *big.Int
		Addr   common.Address
	})
	if !ok {
		return nil, errors.New("result is not the expected type")
	}

	if validators.Staked == nil {
		return nil, errors.New("validator stake is nil")
	}

	return validators.Staked, nil
}

func (m *EngineWrapper) GetValidatorMinimum() (*big.Int, error) {

	check := func() (interface{}, error) {
		return m.Engine.GetValidatorMinimum(nil)
	}

	result, err := utils.ExpRetry(m.logger, check, 3, 1000)

	if err != nil {
		return nil, err
	}

	return result.(*big.Int), nil
}

func (m *EngineWrapper) VersionCheck(minerVersion *big.Int) bool {
	version, err := m.Engine.Version(nil)

	if err != nil {
		m.logger.Error().Err(err).Msg("could not get engine version")
		return false
	}

	if version.Cmp(minerVersion) <= 0 {
		m.logger.Info().Int("version", int(version.Int64())).Msg("miner version is up to date")
		return true
	} else {
		m.logger.Warn().Int("version", int(version.Int64())).Msg("miner version is out of date")
	}

	return false
}

func (m *EngineWrapper) IsPaused() (bool, error) {
	check := func() (interface{}, error) {
		return m.Engine.Paused(nil)
	}

	result, err := utils.ExpRetry(m.logger, check, 3, 1000)

	if err != nil {
		return false, err
	}

	return result.(bool), nil
}

func (m *EngineWrapper) GetContestation(taskId task.TaskId) (*struct {
	Validator        common.Address
	Blocktime        uint64
	FinishStartIndex uint32
	SlashAmount      *big.Int
}, error) {

	getSol := func() (interface{}, error) {
		return m.Engine.Contestations(nil, taskId)
	}

	result, err := utils.ExpRetry(m.logger, getSol, 3, 1000)
	if err != nil {
		return nil, err
	}

	res, ok := result.(struct {
		Validator        common.Address
		Blocktime        uint64
		FinishStartIndex uint32
		SlashAmount      *big.Int
	})

	// Panic here because this is a critical error and there is no recovery
	// we need to update the code to handle this error
	if !ok {
		panic("Result is not of the expected type from engine.Contestations")
	}

	return &res, nil
}

// IsValidatorEligibleToVote checks if a validator can vote on a specific task contestation.
// It returns true if eligible, false otherwise, along with the status code from the contract and any error.
func (ew *EngineWrapper) IsValidatorEligibleToVote(validatorAddr common.Address, taskId task.TaskId) (bool, uint64, error) {
	callOpts := &bind.CallOpts{} // Use default CallOpts

	logger := ew.logger.With().Str("validator", validatorAddr.Hex()).Str("task", taskId.String()).Logger()

	check := func() (interface{}, error) {
		return ew.Engine.ValidatorCanVote(callOpts, validatorAddr, taskId)
	}
	result, err := utils.ExpRetry(logger, check, 3, 1000)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to check ValidatorCanVote during eligibility check")
		return false, 99, err // Indicate error with a distinct code
	}

	voteStatusCode, ok := result.(*big.Int)
	if !ok {
		err := errors.New("ValidatorCanVote returned unexpected type")
		logger.Error().Err(err).Msgf("type was %T", result)
		return false, 98, err // Indicate type error
	}

	if voteStatusCode == nil {
		err := errors.New("ValidatorCanVote returned nil status code")
		logger.Error().Err(err).Msg("ValidatorCanVote returned nil status code")
		return false, 97, err // Indicate nil return
	}

	statusCode := voteStatusCode.Uint64()
	if statusCode == 0 {
		// Validator is eligible
		return true, 0, nil
	}

	// Validator is not eligible, log the reason
	var reason string
	switch statusCode {
	case 1:
		reason = "contestation doesn't exist"
	case 2:
		reason = "voting period ended"
	case 3:
		reason = "already voted"
	case 4:
		reason = "validator never staked"
	case 5:
		reason = "validator staked too long ago"
	case 6:
		reason = "validator staked too recently"
	default:
		reason = fmt.Sprintf("unknown status code: %d", statusCode)
	}
	logger.Debug().Str("reason", reason).Uint64("statusCode", statusCode).Msg("validator cannot vote")
	return false, statusCode, nil
}

// V5:
// get the reward for a specific model - this is used to calculate the reward for a task
// follows the same logic as the engine contract
func (m *EngineWrapper) GetModelReward(modelId [32]byte) (*big.Int, error) {
	reward, err := m.Engine.GetReward(nil)
	if err != nil {
		m.logger.Error().Err(err).Msg("could not get reward!")
		return nil, err
	}

	//0.016461004670253767

	m.logger.Debug().Str("model", common.Bytes2Hex(modelId[:])).Str("reward", reward.String()).Msg("model reward")

	// get model rate from engine
	modelRate, err := m.Engine.Models(nil, modelId)
	if err != nil {
		m.logger.Error().Err(err).Msg("could not get model rate!")
		return nil, err
	}

	// as per engine contract: default to 1e18, so contract still works even if voter is not set
	gaugeMultiplier := big.NewInt(1e18)

	// in reality, this is will always be true on mainnet, but we check anyway as useful for
	// local engine deployment without voter system
	if m.Voter != nil {
		isGauge, err := m.Voter.IsGauge(nil, modelId)
		if err == nil && isGauge {
			gaugeMultiplier, err = m.Voter.GetGaugeMultiplier(nil, modelId)
			if err != nil {
				gaugeMultiplier = big.NewInt(1e18)
			}
		} else {
			gaugeMultiplier = big.NewInt(0)
		}
	}

	// if model rate is 0 and gauge multiplier is 0, return 0
	if modelRate.Rate.Cmp(big.NewInt(0)) > 0 && gaugeMultiplier.Cmp(big.NewInt(0)) > 0 {

		m.logger.Debug().Str("model", common.Bytes2Hex(modelId[:])).Str("rate", modelRate.Rate.String()).Str("gauge", gaugeMultiplier.String()).Msg("model reward")

		// Calculate total reward with gauge multiplier
		totalReward := new(big.Int).Mul(reward, modelRate.Rate)
		totalReward = totalReward.Mul(totalReward, gaugeMultiplier)
		totalReward = totalReward.Div(totalReward, RewardDenominator)

		/* replicate solidity calcs in go:
		if (total > 0) {
			uint256 treasuryReward = total -
				(total * (1e18 - treasuryRewardPercentage)) /
				1e18;

			// v3
			uint256 taskOwnerReward = total -
				(total * (1e18 - taskOwnerRewardPercentage)) /
				1e18;

			uint256 validatorReward = total - treasuryReward - taskOwnerReward; // v3

			baseToken.transfer(treasury, treasuryReward);
			baseToken.transfer(tasks[taskid_].owner, taskOwnerReward); // v3
			baseToken.transfer(
				solutions[taskid_].validator,
				validatorReward
			);

			emit RewardsPaid(total, treasuryReward, taskOwnerReward, validatorReward);
		}*/

		return totalReward, nil
	} else {
		m.logger.Debug().Str("model", common.Bytes2Hex(modelId[:])).Msg("model reward is 0")
		return big.NewInt(0), nil
	}
}
