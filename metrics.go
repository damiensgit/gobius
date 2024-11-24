package main

import (
	"context"
	"fmt"
	"gobius/utils"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// TODO: rename this to reflect actual function
type GasMetrics struct {
	TotalGasUsed     *big.Int
	Commitments      *big.Int
	Solutions        *big.Int
	Claims           *big.Int
	Tasks            *big.Int
	lastEthPrice     float64
	lastBasePrice    float64
	lastReward       float64
	lastSellTime     time.Time
	sessionStartTime time.Time
	rewardEMA        *utils.MovingAveragePrice
	basePriceEMA     *utils.MovingAveragePrice
	gasPriceEMA      *utils.MovingAveragePrice
	profitEMA        *utils.MovingAveragePrice
	services         *Services

	// TODO: Add more metrics here like tx counts, etc.
}

func NewMetricsManager(ctx context.Context, d time.Duration) *GasMetrics {

	// Get the services from the context
	services, ok := ctx.Value(servicesKey{}).(*Services)
	if !ok {
		panic("Could not get services from context")
	}

	ema := utils.NewMovingAveragePrice(240, 4*time.Hour)
	basepriceema := utils.NewMovingAveragePrice(240, 4*time.Hour)
	gasema := utils.NewMovingAveragePrice(240, 4*time.Hour)

	samplesPerHr := int(60 * (60.0 / d.Seconds()))
	profitema := utils.NewMovingAveragePrice(samplesPerHr, 1*time.Hour)

	gm := &GasMetrics{
		TotalGasUsed:     big.NewInt(0),
		Commitments:      big.NewInt(0),
		Solutions:        big.NewInt(0),
		Claims:           big.NewInt(0),
		Tasks:            big.NewInt(0),
		lastEthPrice:     0,
		lastBasePrice:    0,
		lastReward:       0,
		lastSellTime:     time.Now(),
		sessionStartTime: time.Now(),
		rewardEMA:        ema,
		basePriceEMA:     basepriceema,
		gasPriceEMA:      gasema,
		profitEMA:        profitema,
		services:         services,
	}

	var err error
	gm.lastBasePrice, gm.lastEthPrice, err = gm.services.Paraswap.GetPrices()
	if err != nil {
		panic(err)
	}

	return gm
}

func (gm *GasMetrics) Start() {
	go gm.updateMetrics(time.Duration(60) * time.Second)
}

func (gm *GasMetrics) String() string {
	return fmt.Sprintf("TotalGasUsed: %s, Commitments: %s, Solutions: %s, Claims: %s", gm.TotalGasUsed.String(), gm.Commitments.String(), gm.Solutions.String(), gm.Claims.String())
}

func (gm *GasMetrics) AddReward(reward float64) {
	if math.IsNaN(reward) {
		gm.services.Logger.Error().Msg("AddReward: reward is NaN")
	}
	gm.rewardEMA.Add(reward)
}

func (gm *GasMetrics) AddBasePrice(price float64) {
	if math.IsNaN(price) {
		gm.services.Logger.Error().Msg("AddReward: reward is NaN")
	}
	gm.basePriceEMA.Add(price)
}

func (gm *GasMetrics) AddBasefee(basefee float64) {
	gm.gasPriceEMA.Add(basefee)
}

func (gm *GasMetrics) PrintPrice() string {
	return gm.rewardEMA.String()
}

func (gm *GasMetrics) updateMetrics(pollingtime time.Duration) {
	ticker := time.NewTicker(pollingtime)
	for range ticker.C {
		var err error
		gm.lastBasePrice, gm.lastEthPrice, err = gm.services.Paraswap.GetPrices()
		if err != nil {
			gm.services.Logger.Error().Err(err).Msg("could not get prices from sushi api!")
			continue
		}

		gm.AddBasePrice(gm.lastBasePrice)

		basefee, err := gm.services.OwnerAccount.Client.GetBaseFee()
		if err != nil {
			gm.services.Logger.Error().Err(err).Msg("could not get basefee!")
			continue
		}

		// convert basefee to gwei

		basefeeingwei := gm.services.Eth.ToFloat(new(big.Int).Mul(basefee, big.NewInt(1e9)))

		if basefeeingwei < 10 {
			gm.AddBasefee(basefeeingwei)
		}

		val, err := gm.services.Engine.Engine.GetReward(nil)
		if err != nil {
			gm.services.Logger.Error().Err(err).Msg("could not get reward!")
			continue
		}

		if val == nil {
			gm.services.Logger.Error().Err(err).Msg("rewards returned as nil!")
			continue
		}

		gm.lastReward = gm.services.Config.BaseConfig.BaseToken.ToFloat(val) * 0.9 // remove 10% for the treasury
		gm.AddReward(gm.lastReward)

		totalCostInUSD := gm.services.Config.BaseConfig.BaseToken.ToFloat(gm.TotalGasUsed) * gm.lastEthPrice
		totalCostInUSDSinceLastSell := gm.services.Config.BaseConfig.BaseToken.ToFloat(gm.GetTotals()) * gm.lastEthPrice

		totalCostInUSDSinceLastSellFmt := fmt.Sprintf("%0.4f$", totalCostInUSDSinceLastSell)
		totalCostInUSDFmt := fmt.Sprintf("%0.4f$", totalCostInUSD)
		ethPriceFmt := fmt.Sprintf("%0.4f$", gm.lastEthPrice)
		basePriceFmt := fmt.Sprintf("%0.4f$", gm.lastBasePrice)

		lastRewardFmt := fmt.Sprintf("%.3g", gm.lastReward)
		basefeeingweiFmt := fmt.Sprintf("%.6g", basefeeingwei)
		gm.services.Logger.Info().Str("eth_price", ethPriceFmt).Str("last_base_price", basePriceFmt).Str("last_reward", lastRewardFmt).Str("gas_cost", totalCostInUSDSinceLastSellFmt).Str("total_gas_cost", totalCostInUSDFmt).Msg("⛏️ gas cost and reward metrics")
		gm.services.Logger.Info().Str("price_ema", gm.basePriceEMA.String()).Str("reward_ema", gm.rewardEMA.String()).Msg("⛏️ ema metrics")
		gm.services.Logger.Info().Str("last_basefee", basefeeingweiFmt).Str("basefee_ema", gm.gasPriceEMA.String()).Msg("⛏️ gas price metrics")

		//Str("since_last_sell", time.Since(gm.lastSellTime).String())

		averageTaskPerPeriod := gm.services.TaskTracker.AverageTasksPerPeriod()

		// Cost
		timeSinceSessionStart := time.Since(gm.sessionStartTime).Minutes()
		if timeSinceSessionStart > 0 {
			averageCostsPerMin := (totalCostInUSD / timeSinceSessionStart)
			averageCostsPerHr := averageCostsPerMin * 60
			averageCostsPerDay := averageCostsPerHr * 24

			averageCostsPerMinFmt := fmt.Sprintf("%0.2f$", averageCostsPerMin)
			averageCostsPerHrFmt := fmt.Sprintf("%0.2f$", averageCostsPerHr)
			averageCostsPerDayFmt := fmt.Sprintf("%0.2f$", averageCostsPerDay)

			tokenIncomePerPeriod := averageTaskPerPeriod * gm.rewardEMA.Average()
			tokenIncomePerHour := tokenIncomePerPeriod * 60 // TODO: dont assume 1 minute period
			tokenIncomePerDay := tokenIncomePerHour * 24

			tokenIncomePerPeriodFmt := fmt.Sprintf("%.4g", tokenIncomePerPeriod)
			tokenIncomePerHourFmt := fmt.Sprintf("%.4g", tokenIncomePerHour)
			tokenIncomePerDayFmt := fmt.Sprintf("%.4g", tokenIncomePerDay)

			incomePerPeriod := tokenIncomePerPeriod * gm.basePriceEMA.Average()
			incomePerHour := incomePerPeriod * 60
			incomePerDay := incomePerHour * 24

			incomePerPeriodFmt := fmt.Sprintf("%0.2f$", incomePerPeriod)
			incomePerHourFmt := fmt.Sprintf("%0.2f$", incomePerHour)
			incomePerDayFmt := fmt.Sprintf("%0.2f$", incomePerDay)

			profitPerPeriod := incomePerPeriod - averageCostsPerMin
			profitPerHour := incomePerHour - averageCostsPerHr
			profitPerDay := incomePerDay - averageCostsPerDay

			profitPerPeriodFmt := fmt.Sprintf("%0.2f$ (%.6gΞ)", profitPerPeriod, profitPerPeriod/gm.lastEthPrice)
			profitPerHourFmt := fmt.Sprintf("%0.2f$ (%.6gΞ)", profitPerHour, profitPerHour/gm.lastEthPrice)
			profitPerDayFmt := fmt.Sprintf("%0.2f$ (%.6gΞ)", profitPerDay, profitPerDay/gm.lastEthPrice)

			gm.services.Logger.Info().Str("per_min", tokenIncomePerPeriodFmt).Str("per_hr", tokenIncomePerHourFmt).Str("per_day", tokenIncomePerDayFmt).Msg("💎 aius income metrics")
			gm.services.Logger.Info().Str("per_min", incomePerPeriodFmt).Str("per_hr", incomePerHourFmt).Str("per_day", incomePerDayFmt).Msg("💰 income metrics")
			gm.services.Logger.Info().Str("per_min", averageCostsPerMinFmt).Str("per_hr", averageCostsPerHrFmt).Str("per_day", averageCostsPerDayFmt).Msg("💸 cost metrics")
			gm.services.Logger.Info().Str("per_min", profitPerPeriodFmt).Str("per_hr", profitPerHourFmt).Str("per_day", profitPerDayFmt).Msg("📈 profit metrics")
		}

		if gm.services.Config.ValidatorConfig.SellInterval > 0 && time.Since(gm.lastSellTime) >= time.Duration(gm.services.Config.ValidatorConfig.SellInterval)*time.Second {
			// get the total spent since we last sold and calculate how much AIUS we need to sell to cover this cost
			// Rules:
			// 1. we want to only start selling when the AIUS balance is over some threshold (to cover any validator deposits we need)
			// 2. the amount we sell needs to leave this threshold in place
			// 3. we want to sell enough to cover the cost of the gas we've used since the last sell

			// get the AIUS balance
			aiusBalanceAsBig, err := gm.services.Basetoken.BalanceOf(nil, gm.services.OwnerAccount.Address)
			if err != nil {
				gm.services.Logger.Err(err).Msg("failed to get balance")
				continue
			}

			// convert AIUS balance to float
			aiusBalance := gm.services.Config.BaseConfig.BaseToken.ToFloat(aiusBalanceAsBig)

			// get the AIUS balance in USD
			aiusBalanceUSD := aiusBalance * gm.lastBasePrice

			aiusToSell := -1.0
			aiusSellMethod := "default"

			// if we are selling all over the threshold, then we sell all the AIUS over the threshold
			if gm.services.Config.ValidatorConfig.SellAllOverThreshold {
				aiusToSell = aiusBalance - gm.services.Config.ValidatorConfig.MinBasetokenThreshold
				if gm.services.Config.ValidatorConfig.SellMaxAmount > 0 && aiusToSell > gm.services.Config.ValidatorConfig.SellMaxAmount {
					aiusToSell = gm.services.Config.ValidatorConfig.SellMaxAmount
				}
				aiusSellMethod = "all over threshold"
			} else {

				if gm.services.Config.ValidatorConfig.SellEthBalanceTarget > 0 {
					// if this value is > 0 e.g. set then we want to ensure the balance of ETH reaches this target
					// so we weant to keep selling all the AIUS until we reach this target (over the min )
					// get the Eth balance
					ethBalance, err := gm.services.OwnerAccount.GetBalance()
					if err != nil {
						gm.services.Logger.Err(err).Str("account", gm.services.OwnerAccount.Address.String()).Msg("could not get eth balance on account")
						continue
					}
					// convert ETH balance to float
					balAsFloat := gm.services.Eth.ToFloat(ethBalance)

					// we want to hit this target so only sell enough to get us there
					if balAsFloat < gm.services.Config.ValidatorConfig.SellEthBalanceTarget {
						aiusToSell = aiusBalance - gm.services.Config.ValidatorConfig.MinBasetokenThreshold
						//ethToBuy := gm.services.Config.ValidatorConfig.SellEthBalanceTarget - balAsFloat
						//aiusToSell = ethToBuy * (gm.lastEthPrice / gm.lastBasePrice)
						if gm.services.Config.ValidatorConfig.SellMaxAmount > 0 && aiusToSell > gm.services.Config.ValidatorConfig.SellMaxAmount {
							aiusToSell = gm.services.Config.ValidatorConfig.SellMaxAmount
						}
						if aiusToSell > aiusBalance-gm.services.Config.ValidatorConfig.MinBasetokenThreshold {
							aiusToSell = aiusBalance - gm.services.Config.ValidatorConfig.MinBasetokenThreshold
						}
						aiusSellMethod = "eth balance target"
					}
				}

				if aiusToSell < 0 {
					// get the amount of AIUS we need to sell
					aiusToSell = (totalCostInUSDSinceLastSell / gm.lastBasePrice)

					aiusToSell += aiusToSell * gm.services.Config.ValidatorConfig.SellBuffer

					// Adjust for the treasury split
					if gm.services.Config.ValidatorConfig.TreasurySplit > 0.0 {
						aiusToSell /= (1 - gm.services.Config.ValidatorConfig.TreasurySplit)
					}

					// Sell n Eth worth of AIUS
					if gm.services.Config.ValidatorConfig.SellProfitInEth > 0.0 {
						aiusToSell += gm.services.Config.ValidatorConfig.SellProfitInEth * (gm.lastEthPrice / gm.lastBasePrice)
					}

					if gm.services.Config.ValidatorConfig.SellMaxAmount > 0 && aiusToSell > gm.services.Config.ValidatorConfig.SellMaxAmount {
						aiusToSell = gm.services.Config.ValidatorConfig.SellMaxAmount
					}

					aiusSellMethod = "cost since last sell"
				}
			}

			gm.services.Logger.Info().
				Str("aiusbalance", fmt.Sprintf("%0.4f", aiusBalance)).
				Str("aiusbalanceUSD", fmt.Sprintf("%0.4f$", aiusBalanceUSD)).
				Str("amount_to_sell", fmt.Sprintf("%0.4f", aiusToSell)).
				Str("min_sell_threshold", fmt.Sprintf("%0.4f", gm.services.Config.ValidatorConfig.MinBasetokenThreshold)).
				Str("sell_method", aiusSellMethod).Msg("autosell checks")

			// check if we have enough AIUS to sell
			if aiusToSell > 0 && aiusToSell >= gm.services.Config.ValidatorConfig.SellMinAmount && aiusBalance-aiusToSell >= gm.services.Config.ValidatorConfig.MinBasetokenThreshold {

				if gm.services.Config.ValidatorConfig.TreasurySplit > 0.0 {
					if gm.services.Config.ValidatorConfig.TreasurySplit < 0 || gm.services.Config.ValidatorConfig.TreasurySplit > 1.0 {
						gm.services.Logger.Warn().Float64("split", gm.services.Config.ValidatorConfig.TreasurySplit).Msg("treasury split value is invalid, must be between 0 and 1")
					}

					aiusToTransfer := aiusToSell * gm.services.Config.ValidatorConfig.TreasurySplit

					if aiusToTransfer > 0.0005 {

						if success := gm.transferBasetokens(gm.services.Config.ValidatorConfig.TreasuryAddress, aiusToTransfer); success {
							aiusToSell -= aiusToTransfer
						}
					}
				}

				gm.services.Logger.Info().Str("amount_to_sell", fmt.Sprintf("%0.4f", aiusToSell)).Msg("💰 selling AIUS")

				tx, err := gm.services.Paraswap.Allowance(aiusToSell)

				if err != nil {
					gm.services.Logger.Error().Err(err).Msg("❌ could not approve AIUS!")
					continue
				}

				if tx != nil {
					gm.services.Logger.Info().Msg("✔️approving AIUS to be sold")

					_, success, _, _ := gm.services.SenderOwnerAccount.WaitForConfirmedTx(gm.services.Logger, tx)

					if !success {
						continue
					}

					gm.services.Logger.Info().Str("txhash", tx.Hash().String()).Msg("✅ allowance increased")
				}

				// sell the AIUS
				tx, err = gm.services.Paraswap.SellAius(aiusToSell)
				if err != nil {
					gm.services.Logger.Error().Err(err).Msg("❌ could not sell AIUS!")
					continue
				}

				_, success, _, _ := gm.services.SenderOwnerAccount.WaitForConfirmedTx(gm.services.Logger, tx)
				if !success {
					continue
				}

				gm.services.Logger.Info().Str("txhash", tx.Hash().String()).Str("amount_sold", fmt.Sprintf("%0.4f", aiusToSell)).Msg("✅ AIUS sold!")

				gm.lastSellTime = time.Now()
				gm.Reset()

			} else {
				gm.services.Logger.Info().Msg("amount to sell is below min amount or will fall balance below threshold")
			}

		}

	}
}

func (tm *GasMetrics) transferBasetokens(to common.Address, amount float64) bool {

	tm.services.Logger.Info().Msgf("💼 transfering %.4g %s to treasury address %s", amount, tm.services.Config.BaseConfig.BaseToken.Symbol, to.String())

	amountAsBig := tm.services.Config.BaseConfig.BaseToken.FromFloat(amount)

	tx, err := tm.services.SenderOwnerAccount.NonceManagerWrapper(5, 425, 1.5, true, func(opts *bind.TransactOpts) (interface{}, error) {
		opts.GasLimit = 0
		return tm.services.Basetoken.Transfer(opts, to, amountAsBig)
	})

	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("error sending transfer")
		return false
	}

	_, success, _, err := tm.services.SenderOwnerAccount.WaitForConfirmedTx(tm.services.Logger, tx)

	if err != nil {
		tm.services.Logger.Error().Err(err).Msg("error waiting for transfer")
	}

	if success {
		tm.services.Logger.Info().Str("txhash", tx.Hash().String()).Str("amount_transfered", fmt.Sprintf("%.4g", amount)).Msg("✅ transfered to treasury!")
	}

	return success
}

func (gm *GasMetrics) Reset() {
	//gm.TotalGasUsed.SetInt64(0)
	gm.Commitments.SetInt64(0)
	gm.Solutions.SetInt64(0)
	gm.Claims.SetInt64(0)
	gm.Tasks.SetInt64(0)
}

func (gm *GasMetrics) GetTotals() *big.Int {
	total := new(big.Int).Set(gm.Tasks)
	total.Add(total, gm.Commitments)
	total.Add(total, gm.Solutions)
	total.Add(total, gm.Claims)
	return total
}

func (gm *GasMetrics) AddTotal(gas *big.Int) {
	gm.TotalGasUsed.Add(gm.TotalGasUsed, gas)
}

func (gm *GasMetrics) AddCommitment(gas *big.Int) {
	gm.AddTotal(gas)

	totalCostInUSD := fmt.Sprintf("%0.4f$", gm.services.Config.BaseConfig.BaseToken.ToFloat(gas)*gm.lastEthPrice)
	gm.services.Logger.Info().Str("cost", totalCostInUSD).Msg("batch commitment tx cost in USD")

	gm.Commitments.Add(gm.Commitments, gas)
}

func (gm *GasMetrics) AddSolution(gas *big.Int) {
	gm.AddTotal(gas)

	totalCostInUSD := fmt.Sprintf("%0.4f$", gm.services.Config.BaseConfig.BaseToken.ToFloat(gas)*gm.lastEthPrice)
	gm.services.Logger.Info().Str("cost", totalCostInUSD).Msg("batch solution tx cost")

	gm.Solutions.Add(gm.Solutions, gas)
}

func (gm *GasMetrics) AddClaim(gas *big.Int) {
	gm.AddTotal(gas)

	totalCostInUSD := fmt.Sprintf("%0.4f$", gm.services.Config.BaseConfig.BaseToken.ToFloat(gas)*gm.lastEthPrice)
	gm.services.Logger.Info().Str("cost", totalCostInUSD).Msg("batch claim tx cost")

	gm.Claims.Add(gm.Claims, gas)
}

func (gm *GasMetrics) AddTasks(gas *big.Int) {
	gm.AddTotal(gas)

	totalCostInUSD := fmt.Sprintf("%0.4f$", gm.services.Config.BaseConfig.BaseToken.ToFloat(gas)*gm.lastEthPrice)
	gm.services.Logger.Info().Str("cost", totalCostInUSD).Msg("batch tasks tx cost")

	gm.Tasks.Add(gm.Tasks, gas)
}

func (gm *GasMetrics) ClaimValue(claimCount int) {
	value := gm.lastReward * float64(claimCount) * gm.lastBasePrice
	claimValueFmt := fmt.Sprintf("%0.4f$", value)
	gm.services.Logger.Info().Str("claimvalue", claimValueFmt).Msg("⛏️ batch claim value")
}
