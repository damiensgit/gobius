package main

import (
	"errors"
	"gobius/bindings/engine"
	"gobius/bindings/voter"
	task "gobius/common"
	"gobius/storage"
	"gobius/utils"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
)

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

// returns true if the task is valid and can be claimed or false if not
func (m *EngineWrapper) CanTaskIdBeClaimed(claim storage.ClaimTask, cooldownTimes map[common.Address]uint64) (bool, error) {

	taskIdStr := claim.ID.String()
	//taskIdByte, err := task.ConvertTaskIdString2Bytes(taskid.ID)

	// TODO: add this check at claim time once to make sure we have enough staked
	// require(
	// 	validators[solutions[taskid_].validator].staked -
	// 		validatorWithdrawPendingAmount[solutions[taskid_].validator] >=
	// 		getValidatorMinimum(),
	// 	"validator min staked too low"
	// );

	contestationDetails, err := m.GetContestation(claim.ID)
	if err != nil {
		m.logger.Error().Err(err).Str("task", taskIdStr).Msg("cloud not get contestation details")
		return false, err
	}

	if contestationDetails.Validator.String() != "0x0000000000000000000000000000000000000000" {
		contestor := contestationDetails.Validator.String()

		m.logger.Warn().Str("task", taskIdStr).Str("contestor", contestor).Str("slashedamount", contestationDetails.SlashAmount.String()).Msg("⚠️ task was contested ⚠️")

		return false, nil
	} else {
		solution, err := m.GetSolution(claim.ID)
		if err != nil {
			m.logger.Error().Err(err).Str("task", taskIdStr).Msg("cloud not get solution details")
			return false, err
		}

		solTime := time.Unix(int64(solution.Blocktime), 0)

		// Check if user can even claim this task - if they lost a contestation then they forfeit all claims in the cooldown period
		// which is the last constestation loss time + min claim solution time + the contestation vote period time
		cooldownTime := cooldownTimes[solution.Validator]
		if solution.Blocktime <= cooldownTime {
			m.logger.Warn().Str("taskid", taskIdStr).Msg("⚠️ claim is lost due to lost contestation cooldown - removing from storage ⚠️")
			return false, nil
		}

		m.logger.Debug().Str("taskid", taskIdStr).Bool("claimed", solution.Claimed).Time("solved", solTime).Str("validator", solution.Validator.String()).Msg("solution information")

		if solution.Claimed {
			return false, nil
		}
	}
	return true, nil
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


		/* replicate this solidiy calcs in go:
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
