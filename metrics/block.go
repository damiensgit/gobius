package metrics

import (
	"context"
	"fmt"
	"gobius/bindings/engine"
	"gobius/client"
	"gobius/config"
	"log"
	"math/big"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/olekukonko/tablewriter"
)

type Event struct {
	Count     map[common.Address]float64
	Timestamp time.Time
}

type TableCounts struct {
	Tasks       float64
	Solutions   float64
	Commitments float64
	Claims      float64
	Withdraws   float64
	Deposits    float64
	WithdrawSum float64
	DepositSum  float64
}

type BlockMetrics struct {
	TasksSubmitted     []Event
	SubmittedSolutions []Event
	SignalCommitments  []Event
	Claims             []Event
	Withdraw           []Event
	Deposit            []Event
	WithdrawSum        []Event
	DepositSum         []Event
	SessionTotal       map[common.Address]*TableCounts
	rpcClient          *client.Client
	filter             ethereum.FilterQuery
	engineContract     *engine.Engine
	appConfig          *config.AppConfig
}

// var contractAddress = common.HexToAddress("0xYourContractAddress")
var taskSubmittedEvent common.Hash
var signalCommitmentEvent common.Hash
var solutionSubmittedEvent common.Hash
var solutionClaimedEvent common.Hash
var validatorWithdrawEvent common.Hash
var validatorDepositEvent common.Hash

func init() {

	// TODO: move these out of here!
	event, err := engine.EngineMetaData.GetAbi()
	if err != nil {
		panic("error getting engine abi")
	}

	taskSubmittedEvent = event.Events["TaskSubmitted"].ID
	signalCommitmentEvent = event.Events["SignalCommitment"].ID
	solutionSubmittedEvent = event.Events["SolutionSubmitted"].ID
	solutionClaimedEvent = event.Events["SolutionClaimed"].ID
	validatorWithdrawEvent = event.Events["ValidatorWithdraw"].ID
	validatorDepositEvent = event.Events["ValidatorDeposit"].ID

}

func NewBlockMetrics(rpcClient *client.Client, appConfig *config.AppConfig, engineContract *engine.Engine) *BlockMetrics {

	query := ethereum.FilterQuery{
		Addresses: []common.Address{appConfig.BaseConfig.EngineAddress},
		Topics:    [][]common.Hash{{taskSubmittedEvent, signalCommitmentEvent, solutionSubmittedEvent, solutionClaimedEvent, validatorWithdrawEvent, validatorDepositEvent}},
	}

	bm := &BlockMetrics{
		TasksSubmitted:     []Event{},
		SubmittedSolutions: []Event{},
		SignalCommitments:  []Event{},
		Claims:             []Event{},
		Withdraw:           []Event{},
		Deposit:            []Event{},
		WithdrawSum:        []Event{},
		DepositSum:         []Event{},
		SessionTotal:       map[common.Address]*TableCounts{},
		rpcClient:          rpcClient,
		filter:             query,
		engineContract:     engineContract,
		appConfig:          appConfig,
	}

	go bm.captureTimeSeries(1 * time.Minute)
	go bm.captureTimeSeries(10 * time.Minute)
	go bm.captureTimeSeries(30 * time.Minute)
	go bm.captureTimeSeries(1 * time.Hour)
	go bm.captureTimeSeries(4 * time.Hour)
	go bm.captureTimeSeries(8 * time.Hour)
	go bm.cleanupOldEvents(8*time.Hour + 1*time.Minute)
	go bm.displaySessionStats(5 * time.Minute)

	return bm
}

func (bm *BlockMetrics) UpdateBlockMetrics(h *types.Header) {

	bm.filter.FromBlock = h.Number
	bm.filter.ToBlock = h.Number

	logs, err := bm.rpcClient.FilterLogs(context.Background(), bm.filter)
	if err != nil {
		log.Fatalf("Failed to filter logs: %v", err)
	}

	//var taskSubmittedEventCount, solutionSubmittedEventCount, signalCommitmentEventCount, solutionClaimedEventCount = 0, 0, 0, 0

	taskSubmittedEventCountMap := map[common.Address]float64{}
	solutionSubmittedEventCountMap := map[common.Address]float64{}
	signalCommitmentEventCountMap := map[common.Address]float64{}
	solutionClaimedEventCountMap := map[common.Address]float64{}
	validatorWithdrawEventCountMap := map[common.Address]float64{}
	validatorDepositEventCountMap := map[common.Address]float64{}
	validatorWithdrawEventSumMap := map[common.Address]float64{}
	validatorDepositEventSumMap := map[common.Address]float64{}
	seenTxHash := map[common.Hash]common.Address{}

	//var found bool

	for _, currentlog := range logs {
		switch currentlog.Topics[0] {
		case taskSubmittedEvent:
			// parsedLog, err := bm.engineContract.ParseTaskSubmitted(currentlog)
			// if err != nil {
			// 	log.Fatalf("Failed to filter logs: %v", err)
			// }
			//from := parsedLog.Sender
			from, found := seenTxHash[currentlog.TxHash]
			if !found {
				tx, _, err := bm.rpcClient.Client.TransactionByHash(context.Background(), currentlog.TxHash)
				if err != nil {
					log.Fatalf("Failed to get transaction: %v", err)
				}

				from, err = types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
				if err != nil {
					log.Fatalf("Failed to get sender: %v", err)
				}
				seenTxHash[currentlog.TxHash] = from
			}
			taskSubmittedEventCountMap[from]++
		case solutionSubmittedEvent:
			parsedLog, err := bm.engineContract.ParseSolutionSubmitted(currentlog)
			if err != nil {
				log.Fatalf("Failed to filter logs: %v", err)
			}

			solutionSubmittedEventCountMap[parsedLog.Addr]++
		case signalCommitmentEvent:
			parsedLog, err := bm.engineContract.ParseSignalCommitment(currentlog)
			if err != nil {
				log.Fatalf("Failed to filter logs: %v", err)
			}

			signalCommitmentEventCountMap[parsedLog.Addr]++

		case solutionClaimedEvent:
			parsedLog, err := bm.engineContract.ParseSolutionClaimed(currentlog)
			if err != nil {
				log.Fatalf("Failed to filter logs: %v", err)
			}

			solutionClaimedEventCountMap[parsedLog.Addr]++
		case validatorWithdrawEvent:
			parsedLog, err := bm.engineContract.ParseValidatorWithdraw(currentlog)
			if err != nil {
				log.Fatalf("Failed to filter logs: %v", err)
			}
			amountInFloat := bm.appConfig.BaseConfig.BaseToken.ToFloat(parsedLog.Amount)
			validatorWithdrawEventCountMap[parsedLog.Addr]++
			validatorWithdrawEventSumMap[parsedLog.Addr] += amountInFloat

		case validatorDepositEvent:
			parsedLog, err := bm.engineContract.ParseValidatorDeposit(currentlog)
			if err != nil {
				log.Fatalf("Failed to filter logs: %v", err)
			}
			validatorDepositEventCountMap[parsedLog.Validator]++
			amountInFloat := bm.appConfig.BaseConfig.BaseToken.ToFloat(parsedLog.Amount)
			validatorDepositEventSumMap[parsedLog.Validator] += amountInFloat
		}

	}

	for k, v := range taskSubmittedEventCountMap {
		tc, ok := bm.SessionTotal[k]
		if !ok {
			tc = &TableCounts{}
			bm.SessionTotal[k] = tc
		}
		tc.Tasks += v
	}
	for k, v := range solutionSubmittedEventCountMap {
		tc, ok := bm.SessionTotal[k]
		if !ok {
			tc = &TableCounts{}
			bm.SessionTotal[k] = tc
		}
		tc.Solutions += v
	}
	for k, v := range signalCommitmentEventCountMap {
		tc, ok := bm.SessionTotal[k]
		if !ok {
			tc = &TableCounts{}
			bm.SessionTotal[k] = tc
		}
		tc.Commitments += v
	}
	for k, v := range solutionClaimedEventCountMap {
		tc, ok := bm.SessionTotal[k]
		if !ok {
			tc = &TableCounts{}
			bm.SessionTotal[k] = tc
		}
		tc.Claims += v
	}

	bm.TasksSubmitted = append(bm.TasksSubmitted, Event{Count: taskSubmittedEventCountMap, Timestamp: time.Now()})
	bm.SubmittedSolutions = append(bm.SubmittedSolutions, Event{Count: solutionSubmittedEventCountMap, Timestamp: time.Now()})
	bm.SignalCommitments = append(bm.SignalCommitments, Event{Count: signalCommitmentEventCountMap, Timestamp: time.Now()})
	bm.Claims = append(bm.Claims, Event{Count: solutionClaimedEventCountMap, Timestamp: time.Now()})
	bm.Withdraw = append(bm.Withdraw, Event{Count: validatorWithdrawEventCountMap, Timestamp: time.Now()})
	bm.Deposit = append(bm.Deposit, Event{Count: validatorDepositEventCountMap, Timestamp: time.Now()})
	bm.WithdrawSum = append(bm.WithdrawSum, Event{Count: validatorWithdrawEventSumMap, Timestamp: time.Now()})
	bm.DepositSum = append(bm.DepositSum, Event{Count: validatorDepositEventSumMap, Timestamp: time.Now()})
}

var addressToName map[common.Address]string = map[common.Address]string{}

var mu sync.Mutex

func getName(address common.Address) string {
	who, found := addressToName[address]
	if found {
		return who
	}
	return "unknown:" + address.String()
}

func (bm *BlockMetrics) captureTimeSeries(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		mu.Lock() // Lock the mutex before writing to the console

		totalTasksSubmitted := countEventsInLastInterval(bm.TasksSubmitted, interval)
		totalSubmittedSolutions := countEventsInLastInterval(bm.SubmittedSolutions, interval)
		totalSignalCommitments := countEventsInLastInterval(bm.SignalCommitments, interval)
		totalClaims := countEventsInLastInterval(bm.Claims, interval)
		totalWithdraws := countEventsInLastInterval(bm.Withdraw, interval)
		totalDeposits := countEventsInLastInterval(bm.Deposit, interval)
		totalWithdrawSum := countEventsInLastInterval(bm.WithdrawSum, interval)
		totalDepositSum := countEventsInLastInterval(bm.DepositSum, interval)

		allWithdraws := countEventsAll(bm.Withdraw)
		allDeposits := countEventsAll(bm.Deposit)
		allWithdrawSum := countEventsAll(bm.WithdrawSum)
		allDepositSum := countEventsAll(bm.DepositSum)

		counts := make(map[string]*TableCounts)
		countTotals := make(map[string]*TableCounts)

		for who, count := range totalTasksSubmitted {
			username := getName(who)
			if _, ok := counts[username]; !ok {
				counts[username] = &TableCounts{}
			}
			counts[username].Tasks += count
		}
		for who, count := range totalSubmittedSolutions {
			username := getName(who)

			if _, ok := counts[username]; !ok {
				counts[username] = &TableCounts{}
			}
			counts[username].Solutions += count
		}
		for who, count := range totalSignalCommitments {
			username := getName(who)

			if _, ok := counts[username]; !ok {
				counts[username] = &TableCounts{}
			}
			counts[username].Commitments = count
		}
		for who, count := range totalClaims {
			username := getName(who)

			if _, ok := counts[username]; !ok {
				counts[username] = &TableCounts{}
			}
			counts[username].Claims += count
		}

		for who, count := range totalWithdraws {
			username := getName(who)

			if _, ok := counts[username]; !ok {
				counts[username] = &TableCounts{}
			}
			counts[username].Withdraws += count
		}

		for who, count := range totalDeposits {
			username := getName(who)

			if _, ok := counts[username]; !ok {
				counts[username] = &TableCounts{}
			}
			counts[username].Deposits += count
		}

		for who, count := range totalWithdrawSum {
			username := getName(who)

			if _, ok := counts[username]; !ok {
				counts[username] = &TableCounts{}
			}
			counts[username].WithdrawSum += count
		}

		for who, count := range totalDepositSum {
			username := getName(who)

			if _, ok := counts[username]; !ok {
				counts[username] = &TableCounts{}
			}
			counts[username].DepositSum += count
		}

		for who, count := range allWithdraws {
			username := getName(who)

			if _, ok := countTotals[username]; !ok {
				countTotals[username] = &TableCounts{}
			}
			countTotals[username].Withdraws += count
		}

		for who, count := range allDeposits {
			username := getName(who)

			if _, ok := countTotals[username]; !ok {
				countTotals[username] = &TableCounts{}
			}
			countTotals[username].Deposits += count
		}

		for who, count := range allWithdrawSum {
			username := getName(who)

			if _, ok := countTotals[username]; !ok {
				countTotals[username] = &TableCounts{}
			}
			countTotals[username].WithdrawSum += count
		}

		for who, count := range allDepositSum {
			username := getName(who)

			if _, ok := countTotals[username]; !ok {
				countTotals[username] = &TableCounts{}
			}
			countTotals[username].DepositSum += count
		}

		data := [][]string{}

		commitmentFactor := 129_000.0 / 27_900.0
		claimFactor := 129_000.0 / 47_300.0

		for who, count := range counts {
			rating := count.Tasks + count.Commitments/commitmentFactor + count.Solutions + count.Claims/claimFactor
			data = append(data, []string{who, fmt.Sprintf("%d", int(count.Tasks)), fmt.Sprintf("%d", int(count.Commitments)), fmt.Sprintf("%d", int(count.Solutions)), fmt.Sprintf("%d", int(count.Claims)), fmt.Sprintf("%d", int(rating)), fmt.Sprintf("%d", int(count.Withdraws)), fmt.Sprintf("%d", int(count.Deposits)), fmt.Sprintf("%0.2f", count.WithdrawSum), fmt.Sprintf("%0.2f", count.DepositSum)})
		}

		dataTotals := [][]string{}
		for who, count := range countTotals {
			dataTotals = append(dataTotals, []string{who, fmt.Sprintf("%d", int(count.Withdraws)), fmt.Sprintf("%d", int(count.Deposits)), fmt.Sprintf("%0.2f", count.WithdrawSum), fmt.Sprintf("%0.2f", count.DepositSum)})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetCaption(true, fmt.Sprintf("Events in last %v", interval))
		table.SetHeader([]string{"User", "Tasks", "Commitments", "Solutions", "Claims", "Rating", "Withdraws", "Deposits", "WithdrawSum", "DepositSum"})
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
		table.SetAlignment(tablewriter.ALIGN_RIGHT)
		table.AppendBulk(data)
		table.Render()

		tableTotals := tablewriter.NewWriter(os.Stdout)
		tableTotals.SetCaption(true, "Validator deposit/withdraw totals")
		tableTotals.SetHeader([]string{"User", "Withdraws", "Deposits", "WithdrawSum", "DepositSum"})
		tableTotals.SetAutoWrapText(false)
		tableTotals.SetAutoFormatHeaders(true)
		tableTotals.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
		tableTotals.SetAlignment(tablewriter.ALIGN_RIGHT)
		tableTotals.AppendBulk(dataTotals)
		tableTotals.Render()

		mu.Unlock() // Unlock the mutex after writing to the console

	}
}

func (bm *BlockMetrics) displaySessionStats(interval time.Duration) {
	ticker := time.NewTicker(interval)
	sessionStart := time.Now()
	commitmentFactor := 129_000.0 / 27_900.0
	claimFactor := 129_000.0 / 47_300.0

	for range ticker.C {
		mu.Lock() // Lock the mutex before writing to the console

		data := [][]string{}
		countTotals := make(map[string]*TableCounts)

		for who, count := range bm.SessionTotal {
			username := getName(who)

			if _, ok := countTotals[username]; !ok {
				countTotals[username] = &TableCounts{}
			}
			countTotals[username].Claims += count.Claims
			countTotals[username].Solutions += count.Solutions
			countTotals[username].Commitments += count.Commitments
			countTotals[username].Tasks += count.Tasks
		}

		type UserCount struct {
			Who    string
			Counts *TableCounts
			Rating int
		}

		// Convert the map to a slice
		userCounts := make([]UserCount, 0, len(countTotals))
		for who, count := range countTotals {
			rating := count.Tasks + count.Commitments/commitmentFactor + count.Solutions + count.Claims/claimFactor
			userCounts = append(userCounts, UserCount{who, count, int(rating)})
		}

		// Sort the slice by rating
		sort.Slice(userCounts, func(i, j int) bool {
			return userCounts[i].Rating > userCounts[j].Rating
		})

		for _, userCount := range userCounts {
			data = append(data, []string{userCount.Who, fmt.Sprintf("%d", int(userCount.Counts.Tasks)), fmt.Sprintf("%d", int(userCount.Counts.Commitments)), fmt.Sprintf("%d", int(userCount.Counts.Solutions)), fmt.Sprintf("%d", int(userCount.Counts.Claims)), fmt.Sprintf("%d", int(userCount.Rating))})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetCaption(true, fmt.Sprintf("Session totals since: %s", time.Since(sessionStart)))
		table.SetHeader([]string{"User", "Tasks", "Commitments", "Solutions", "Claims", "Rating"})
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
		table.SetAlignment(tablewriter.ALIGN_RIGHT)
		table.AppendBulk(data)
		table.Render()

		mu.Unlock() // Unlock the mutex after writing to the console

	}
}

func countEventsInLastInterval(events []Event, interval time.Duration) map[common.Address]float64 {
	now := time.Now()

	countByAddress := make(map[common.Address]float64)

	for i := len(events) - 1; i >= 0; i-- {
		if now.Sub(events[i].Timestamp) <= interval {
			//totalCount += events[i].Count
			for address, v := range events[i].Count {
				countByAddress[address] += v
			}
		} else {
			break
		}
	}

	return countByAddress
}

func countEventsAll(events []Event) map[common.Address]float64 {

	countByAddress := make(map[common.Address]float64)

	for i := len(events) - 1; i >= 0; i-- {
		//totalCount += events[i].Count
		for address, v := range events[i].Count {
			countByAddress[address] += v
		}
	}

	return countByAddress
}

func (bm *BlockMetrics) cleanupOldEvents(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		mu.Lock()
		now := time.Now()
		cleanup := func(events []Event) []Event {
			for i, event := range events {
				if now.Sub(event.Timestamp) <= interval {
					return events[i:]
				}
			}
			return []Event{}
		}

		bm.TasksSubmitted = cleanup(bm.TasksSubmitted)
		bm.SubmittedSolutions = cleanup(bm.SubmittedSolutions)
		bm.Claims = cleanup(bm.Claims)
		bm.SignalCommitments = cleanup(bm.SignalCommitments)
		bm.Withdraw = cleanup(bm.Withdraw)
		bm.Deposit = cleanup(bm.Deposit)
		bm.WithdrawSum = cleanup(bm.WithdrawSum)
		bm.DepositSum = cleanup(bm.DepositSum)
		mu.Unlock()
	}

}

func (bm *BlockMetrics) ProcessDepositWithdrawLogs(rpcClient *client.Client, startBlock, endBlock int64, engineContract *engine.Engine) {
	// Define the map to hold the counts
	allWithdraws := make(map[common.Address]float64)
	allDeposits := make(map[common.Address]float64)
	allWithdrawSum := make(map[common.Address]float64)
	allDepositSum := make(map[common.Address]float64)
	uniqueValidators := make(map[common.Address]int)
	allSolutions := make(map[common.Address]float64)
	allClaims := make(map[common.Address]float64)

	log.Printf("Scanning logs from block %d to %d for sols/claims, deposits and withdraws", startBlock, endBlock)

	step := int64(5_000)

	for i := startBlock; i <= endBlock; i += step {

	retry_loop:

		// Define the filter query
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(i),
			ToBlock:   big.NewInt(min(i+step, endBlock)),
			Topics:    [][]common.Hash{{validatorWithdrawEvent, validatorDepositEvent}},
		}
		// Define the filter query
		query2 := ethereum.FilterQuery{
			FromBlock: big.NewInt(i),
			ToBlock:   big.NewInt(min(i+step, endBlock)),
			Topics:    [][]common.Hash{{solutionSubmittedEvent}},
		}
		// Define the filter query
		query3 := ethereum.FilterQuery{
			FromBlock: big.NewInt(i),
			ToBlock:   big.NewInt(min(i+step, endBlock)),
			Topics:    [][]common.Hash{{solutionClaimedEvent}},
		}

		queries := []ethereum.FilterQuery{query, query2, query3}

		alllogs := []types.Log{}
		for _, q := range queries {

			// Get the logs for the block range
			logs, err := rpcClient.FilterLogs(context.Background(), q)
			if err != nil {
				if strings.Contains(err.Error(), "read limit exceeded") || strings.Contains(err.Error(), "unexpected EOF") {
					log.Printf("Read limit exceeded/EOF, retrying with smaller step: %d", step-250)
					step -= 250
					goto retry_loop
				} else {
					log.Fatalf("Failed to filter logs: %v", err)
				}
			}

			alllogs = append(alllogs, logs...)
		}
		if step < 5000 {
			log.Printf("Setting step back to: %d", step+250)
			step += 250
		}

		// Iterate over the logs
		for _, item := range alllogs {
			switch item.Topics[0] {
			case validatorWithdrawEvent:
				parsedLog, err := bm.engineContract.ParseValidatorWithdraw(item)
				if err != nil {
					log.Fatalf("Failed to parse log: %v", err)
				}
				allWithdraws[parsedLog.Addr]++
				amountInFloat := bm.appConfig.BaseConfig.BaseToken.ToFloat(parsedLog.Amount)
				allWithdrawSum[parsedLog.Addr] += amountInFloat
			case validatorDepositEvent:
				parsedLog, err := bm.engineContract.ParseValidatorDeposit(item)
				if err != nil {
					log.Fatalf("Failed to parse log: %v", err)
				}
				allDeposits[parsedLog.Validator]++
				amountInFloat := bm.appConfig.BaseConfig.BaseToken.ToFloat(parsedLog.Amount)
				allDepositSum[parsedLog.Validator] += amountInFloat
			case solutionSubmittedEvent:
				parsedLog, err := bm.engineContract.ParseSolutionSubmitted(item)
				if err != nil {
					log.Fatalf("Failed to parse log: %v", err)
				}
				allSolutions[parsedLog.Addr]++
			case solutionClaimedEvent:
				parsedLog, err := bm.engineContract.ParseSolutionClaimed(item)
				if err != nil {
					log.Fatalf("Failed to parse log: %v", err)
				}
				allClaims[parsedLog.Addr]++
			}
		}

	}

	countTotals := make(map[common.Address]*TableCounts)

	for who, count := range allSolutions {
		uniqueValidators[who] = 1
		if _, ok := countTotals[who]; !ok {
			countTotals[who] = &TableCounts{}
		}
		countTotals[who].Solutions += count
	}

	for who, count := range allClaims {
		uniqueValidators[who] = 1
		if _, ok := countTotals[who]; !ok {
			countTotals[who] = &TableCounts{}
		}
		countTotals[who].Claims += count
	}

	for who, count := range allWithdraws {
		uniqueValidators[who] = 1
		if _, ok := countTotals[who]; !ok {
			countTotals[who] = &TableCounts{}
		}
		countTotals[who].Withdraws += count
	}

	for who, count := range allDeposits {
		uniqueValidators[who] = 1
		if _, ok := countTotals[who]; !ok {
			countTotals[who] = &TableCounts{}
		}
		countTotals[who].Deposits += count
	}

	for who, count := range allWithdrawSum {
		if _, ok := countTotals[who]; !ok {
			countTotals[who] = &TableCounts{}
		}
		countTotals[who].WithdrawSum += count
	}

	for who, count := range allDepositSum {
		if _, ok := countTotals[who]; !ok {
			countTotals[who] = &TableCounts{}
		}
		countTotals[who].DepositSum += count
	}

	optsStart := &bind.CallOpts{
		BlockNumber: big.NewInt(startBlock),
	}
	optsEnd := &bind.CallOpts{
		BlockNumber: big.NewInt(endBlock),
	}

	empty := common.Address{}
	getStake := func(opts *bind.CallOpts, who common.Address) float64 {
		val, err := bm.engineContract.Validators(opts, who)
		if err != nil {
			log.Printf("Failed to get validator info: %v (%s)", err, who.String())
			return 0
		}
		amountInFloat := -.0
		if val.Addr != empty {
			amountInFloat = bm.appConfig.BaseConfig.BaseToken.ToFloat(val.Staked)
		}
		return amountInFloat
	}

	dataTotals := [][]string{}
	for who, count := range countTotals {
		username := getName(who)
		startStake := getStake(optsStart, who)
		endStake := getStake(optsEnd, who)
		dataTotals = append(dataTotals,
			[]string{
				username,
				fmt.Sprintf("%d", int(count.Solutions)),
				fmt.Sprintf("%d", int(count.Claims)),
				fmt.Sprintf("%d", int(count.Withdraws)),
				fmt.Sprintf("%d", int(count.Deposits)),
				fmt.Sprintf("%0.2f", count.WithdrawSum),
				fmt.Sprintf("%0.2f", count.DepositSum),
				fmt.Sprintf("%0.2f", startStake),
				fmt.Sprintf("%0.2f", endStake),
				fmt.Sprintf("%0.2f", count.DepositSum+startStake-count.WithdrawSum),
			})
	}
	tableTotals := tablewriter.NewWriter(os.Stdout)
	tableTotals.SetCaption(true, "Validator deposit/withdraw totals")
	tableTotals.SetHeader([]string{"User", "Solutions", "Claims", "Withdraws", "Deposits", "WithdrawSum", "DepositSum", "StakedAtStartBlk", "StakedAtEndBlk", "StakeChange"})
	tableTotals.SetAutoWrapText(false)
	tableTotals.SetAutoFormatHeaders(true)
	tableTotals.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	tableTotals.SetAlignment(tablewriter.ALIGN_RIGHT)
	tableTotals.AppendBulk(dataTotals)
	tableTotals.Render()

}
