package metrics

import (
	"context"
	"fmt"
	"gobius/arbius/engine"
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

var addressToName map[common.Address]string = map[common.Address]string{
	common.HexToAddress("0xAA4Cc0E46843AdC45F45d774aFEfFBC3224EF8E0"): "lachy",
	common.HexToAddress("0x33A40904D3A6E88d2303027254374BE002Dd2Ccc"): "lachy",
	common.HexToAddress("0xa7fCCF1B5AD8A1b73b15c954442217a593A0cDa5"): "lachy",
	common.HexToAddress("0xd730F6b16083C6b8656e668D9bA5B35a2Aa1a400"): "lachy",
	common.HexToAddress("0x5291f3E0A33935D1FC69567F4E9889dFF502b869"): "lachy",
	common.HexToAddress("0xB2e8f8Bf8b1C2B96A6A790870F1569068673368B"): "lachy",
	common.HexToAddress("0xE12F0f8d71a304b5e5d3CB2A109961fA6aDc16Ff"): "lachy",
	common.HexToAddress("0xF2B218dA3421321821378a50c2A1Aa6942c81D65"): "lachy",
	common.HexToAddress("0xc510569A02EC7a16f8ca2393C8714637f21fF5a5"): "ben",
	common.HexToAddress("0xe6B6C55A46A3F1925CfbE2336f2E8Be135e4BCE9"): "ben",
	common.HexToAddress("0x66a340e8B3D07648851b9538Da3c86FDfDf20B0f"): "ben",
	common.HexToAddress("0x7Fa553D30a028aC7D8d8A5D3e86FE8706D0c783D"): "ben",
	common.HexToAddress("0xeAFAfafaD608daBAfdd899cCE00B2de8FEDaAD4a"): "ben",
	common.HexToAddress("0x874e1d4482CAB921354DA1Cd5e446062f7Df0d29"): "ben",
	common.HexToAddress("0x4A7d57CD22a52FAa85ab08b60333c592AA8F30C5"): "ben",
	common.HexToAddress("0x8F84A3C939d6Ae78cF62EBc5980BFc89aDe6cD04"): "ben",
	common.HexToAddress("0xfd27764bb9d32fDdF85F8Bc29d392c28363747e9"): "ben",
	common.HexToAddress("0x84f6723600388a08Bb8F495374bb37DB59214746"): "ben",
	common.HexToAddress("0x6357d0D4EB25A644d79499CE4078E0b88605280e"): "ben",
	common.HexToAddress("0x5E75D07b3A49091F3DB4b011738dfdE2d9ec8cA7"): "DA",
	common.HexToAddress("0xe99d3fE2eFCb9Ad0B662015EC14b291c1FdeC86D"): "DA",
	common.HexToAddress("0xE0Fca82138fb39a482f65abc40aDA053868fd5B8"): "DA-staker",
	common.HexToAddress("0x7e992C5cbd8dED507175EdB49C79b8A2ae076781"): "blur",
	common.HexToAddress("0x7fB88C0D37720477CEB1BBb02e5F3959483d9de5"): "damien",
	common.HexToAddress("0xD5fe431d0d28Cf7460aCdE73a7A37589a0e96C4E"): "damien-staker",
	common.HexToAddress("0xf730e2b5E85C9b97172CCb6Ea6Fe0e45aDDA5B71"): "smash", // contract
	common.HexToAddress("0x57c7F9bd68C356FE3b9eA03982b5d07787DF5Fb9"): "blur",
	common.HexToAddress("0x4301027aD3EafD203bec23A31841c971767af876"): "blur",
	common.HexToAddress("0x56c632E290D569E6A6F79f6ce096AbA2954eDDBE"): "ben-tasks",
	common.HexToAddress("0x3BA6C658CD69021f517fdcF60a34bC3Efc6B8FE5"): "stock-miner-0",
	common.HexToAddress("0x95964Ad337346D910209D03bFFe5029E974DE6d8"): "stock-miner-0",
	common.HexToAddress("0xAc73b83De2827b72610DaF74182eaF0F334a3ecD"): "stealer",
	common.HexToAddress("0x930AeDf400AfE367e3503Ce9eFAA529AeB2C82A4"): "damien-2",
	common.HexToAddress("0x334e640C3785c364A5B4c1bA5E7c3b7441Eff595"): "lachy-2",
	common.HexToAddress("0xb0F8f7cEB58619a125d92B484378F37d259FBfa1"): "stock-miner-1",
	common.HexToAddress("0x52FA3FD30398034ffB0050Ed5e0665eE91Bec76C"): "stock-miner-2",
	common.HexToAddress("0x5e247a439baA1230F7802e678a666449D027Ab60"): "stock-miner-3",
	common.HexToAddress("0x0145d8445789Fcf805889ef8415808fb5843c2D0"): "new-miner-0x0145",
	common.HexToAddress("0x6A1A5b8f32C1C468a092025D8f947dCB8aE9D624"): "new-miner-0x0145",
	common.HexToAddress("0x69C80127AdA04b1A9b97ed416242242beF870001"): "new-miner-0x0145",
	common.HexToAddress("0xfec84a21ab2439226acb0c0a288a95341a6b2c0e"): "ben",
	common.HexToAddress("0x6241812282320299c4Ce99FcaCCf79CA74D0d9fb"): "smash", // sender
	common.HexToAddress("0x5eb4720d94cb49f892f621d28cc479a51a8c71a0"): "lachy",
	common.HexToAddress("0x2e09E7b91AfE36D7f1aC420845d65D76D1c4Df05"): "damien",
	common.HexToAddress("0xbF843cB80d9ff1eB18d6694f2bfCB505246466ba"): "damien",
	common.HexToAddress("0x32a8397f1C0D61032CDF092b41bc26D02Cc39ce0"): "damien",
	common.HexToAddress("0x41Ac36656988126CA1caed836E236E51FbAdDCcA"): "damien",
	common.HexToAddress("0x90b01b39775C6Ce94A8baaE01DC8983bE89dc01E"): "damien",
	common.HexToAddress("0x2A77b3CB2b7023ee5735a2c85173F0Ad992E7a26"): "damien",
	common.HexToAddress("0x32f7DdFE6b082c1f59f5496829a79c2a69E06614"): "damien",
	common.HexToAddress("0xbFAD1fBE6343664cf3048052FEB848Bb9d9685C7"): "damien",
	common.HexToAddress("0xA8A23164569A182FaD570c8f23f88dB79c345182"): "damien",
	common.HexToAddress("0x64C75bFD5Be8b59E8c81c1c3ce0709D5FB6FA875"): "damien",
	common.HexToAddress("0x9B76782De702bb444a43CcF7Cca56CE6DE002174"): "damien",
	common.HexToAddress("0xaEC65ec372e405213CE72d5c6472709c5B193671"): "DA",
	common.HexToAddress("0x0fb5E652FbD64d37D35eD28894059F7431a26FFF"): "DA",
	common.HexToAddress("0x71FF02CECE8E550624553BCbB52c0249244ea065"): "DA",
	common.HexToAddress("0xDAE7b778bB68Be7B3098b3B2Cbe4780fAF112ef9"): "DA",
	common.HexToAddress("0x3900f7f584Dbcb3f816061EA939EaeAdB0aEe851"): "DA",
	common.HexToAddress("0x4cF21CD658E5a52a9672aE09A14395Cd2ce0eE58"): "DA",
	common.HexToAddress("0x2Ef4E59Ced1106c6C8846F320E522F84648b4499"): "DA",
	common.HexToAddress("0x4b8800f95e9bD1f6BEEAB61494ed8d7c5E458003"): "DA",
	common.HexToAddress("0xFB9198430ea69d0b761A89Ff68899e7cD7951320"): "DA",
	common.HexToAddress("0x82Bb0Fd1313316576e849A5ee0725CDFD4171fdB"): "DA",
	common.HexToAddress("0x4bFD141d4d3ecE1AB55f49Cc00FB4ec99823Db8d"): "new-miner-0x4bFD",
	common.HexToAddress("0x527E2fa7E1171b1bDeF87492Bd1F48c7E4304c4D"): "new-miner-0x4bFD",
	common.HexToAddress("0xC22f96Fa584F330Da44d788a833Fc803a7433dBe"): "ben",
	common.HexToAddress("0x61e1c37C2A1931EcfAd2330f645344544B2214Db"): "blur",
	common.HexToAddress("0x06134fccD9eaf597F2BfD084BC76A8Dd5A0ee9FF"): "smash",
	common.HexToAddress("0x10B979cD1230ee2a04D9F65F51148c66F59DD2B4"): "smash", // main distribution wallet/hodl wallet
	common.HexToAddress("0x737eAeE51d50573baed0e7eEa63Cc7f444299686"): "dreamer",
	common.HexToAddress("0x2f1407e02aa9e63f0e403fe23928ab58ad462778"): "dreamer",
	common.HexToAddress("0x10B979cD1230ee2a04D9F65F51148c66F59DD2B4"): "smash-staker",
	common.HexToAddress("0x0f72046b58f2A5fBC58E6A3490b3A342FB05dF36"): "blur-staker",
}

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
