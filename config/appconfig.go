package config

import (
	"encoding/json"
	"errors"
	"gobius/common"
	"gobius/erc20"
	"math/big"
	"os"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/joho/godotenv"
)

const (
	ConfigPath = "config"
	LogLevel   = "log.level"
	LogFormat  = "log.format"
)

type AppConfig struct {
	DBPath                string            `json:"db_path"`
	LogPath               string            `json:"log_path"`
	LogLevel              int               `json:"log_level"`
	CachePath             string            `json:"cache_path"`
	ListenToTaskSubmitted bool              `json:"listen_to_task_submitted"`
	CheckCommitment       bool              `json:"check_commitment"`
	EvilMode              bool              `json:"evil_mode"`     // for testing purposes
	EvilModeMinTime       int               `json:"evil_mode_min"` // for testing purposes
	EvilModeRandInt       int               `json:"evil_mode_int"` // for testing purposes
	WhitelistTasks        bool              `json:"whitelist_tasks"`
	StealTasks            bool              `json:"steal_tasks"` // not the same as snipe strategy, this just steals tasks from a set of owners defined in snipe.Targets
	NumWorkersPerGPU      int               `json:"num_workers_per_gpu"`
	PriceOracleContract   ethcommon.Address `json:"price_oracle_contract"`

	Miner SolverConfig `json:"solver"`

	TelegramBot TelegramBot `json:"telegram"`

	ValidatorConfig ValidatorConfig `json:"validator_config"`

	BatchTasks BatchTasks `json:"batchtasks"`
	Strategies Strategies `json:"strategies"`
	Blockchain Blockchain `json:"blockchain"`
	Redis      Redis      `json:"redis"`
	RPC        RPC        `json:"rpc"`
	Claim      Claimer    `json:"claim"`
	ML         ML         `json:"ml"`
	IPFS       IPFS       `json:"ipfs"`
	BaseConfig BaseConfig `json:"baseconfig"`
}

type TelegramBot struct {
	Enabled bool   `json:"enabled"`
	Token   string `json:"token"`
	ChatID  string `json:"chat_id"`
}

type ValidatorConfig struct {
	InitialStake            float64           `json:"initial_stake"`             // initial stake to use for the validator
	StakeBufferStakeAmount  float64           `json:"stake_buffer_amount"`       // min buffer amount between validator stake and min stake
	StakeBufferTopupAmount  float64           `json:"stake_buffer_topup_amount"` // amount to top up the stake buffer
	StakeBufferPercent      int               `json:"stake_buffer_percent"`
	StakeBufferTopupPercent int               `json:"stake_buffer_topup_percent"`
	StakeCheck              bool              `json:"stake_check"`             // check if the validator has enough stake
	StakeCheckInterval      string            `json:"stake_check_interval"`    // how often to check the stake
	EthLowThreshold         float64           `json:"eth_low_threshold"`       // if the validator has less than this amount of eth, we send alerts (TODO: implement)
	MinBasetokenThreshold   float64           `json:"min_basetoken_threshold"` // min balance to leave on the validator account
	SellInterval            int               `json:"sell_interval"`           // how often to check if we should sell in seconds
	SellBuffer              float64           `json:"sell_buffer"`             // multiplier for the amount of tokens to sell (e.g. 1.5 means sell 1.5 times the amount of tokens)
	SellProfitInEth         float64           `json:"sell_profit_in_eth"`      // sell this additional amount of AIUS in Eth terms
	SellAllOverThreshold    bool              `json:"sell_all_over_threshold"` // sell all tokens if the balance is over the threshold
	SellMinAmount           float64           `json:"sell_min_amount"`         // minimum amount of tokens to sell
	SellMaxAmount           float64           `json:"sell_max_amount"`         // max amount of tokens to sell
	SellEthBalanceTarget    float64           `json:"sell_eth_bal_target"`     // sell all over threshold if ETH balance is below this target
	TreasurySplit           float64           `json:"treasury_split"`          // split the rewards between the validator and the treasury
	TreasuryAddress         ethcommon.Address `json:"treasury_address"`        // address of the treasury
	PrivateKeys             []string          `json:"private_keys"`            // list of 1 or more validators to use for v3
}

type BatchTasks struct {
	Enabled                  bool     `json:"enabled"`
	MinTasksInQueue          int      `json:"min_tasks_in_queue"` //  the number of tasks left in the queue before we start creating new batches
	OnlyTasks                bool     `json:"only_tasks"`         // stop after a task batches are submitted
	BatchMode                string   `json:"batch_mode"`         // "normal", "account"
	BatchSize                int      `json:"batch_size"`
	NumberOfBatches          int      `json:"number_of_batches"` // number of batches to send out at once
	HoardMode                bool     `json:"hoard_mode"`
	HoardModeBatchSize       int      `json:"hoard_mode_batch_size"`
	HoardModeNumberOfBatches int      `json:"hoard_mode_number_of_batches"` // number of batches to send out at once
	HoardMinGasPrice         float64  `json:"hoard_min_gas_price"`          // min gas price we need to see before we start hoarding tasks
	HoardMaxQueueSize        int      `json:"hoard_max_queue_size"`         // max size we let the queue get when hoarding tasks
	PrivateKeys              []string `json:"private_keys"`
}

type BatchConfig struct {
	MinBatchSize    int `json:"min_batch_size"`
	MaxBatchSize    int `json:"max_batch_size"`
	NumberOfBatches int `json:"number_of_batches"` // number of batches to send out at once
}

type SolverConfig struct {
	Enabled                 bool             `json:"enabled"`
	CommitmentsAndSolutions CommitmentOption `json:"commitments_and_solutions"` // none, both, commitments
	CommitmentBatch         BatchConfig      `json:"commitment_batch"`
	SolutionBatch           BatchConfig      `json:"solution_batch"`
	ConcurrentBatches       bool             `json:"concurrent_batches"`
	ProfitMode              string           `json:"profit_mode"`              // profit mode to use for batch operations
	MinProfit               float64          `json:"min_profit"`               // minimum profit in USD to perform batch operations
	MaxProfit               float64          `json:"max_profit"`               // minimum profit in USD to perform batch operations
	PauseStakeBufferLevel   float64          `json:"pause_stake_buffer_level"` // pause submitting commits/solutions if the buffer = (current stake - min stake) is below this level
	UsePolling              bool             `json:"use_polling"`              // use polling to check profit and submit txes
	PollingTime             string           `json:"polling_time"`
	BatchMode               int              `json:"batch_mode"`               // 0: no batch, 1 - normal batching using storage and polling, 2 - batch manually
	NoChecks                bool             `json:"no_checks"`                // no checks on the tasks, just submit them
	ErrorMaxRetries         int              `json:"error_max_retries"`        // max retries for tx errors
	ErrorBackoffTime        float64          `json:"error_backoff"`            // sleep time between retries
	ErrorBackofMultiplier   float64          `json:"error_backoff_multiplier"` // backoff multiplier tx errors
	MetricsSampleRate       string           `json:"metrics_sample_rate"`      // sample rate for metrics
}

type Strategies struct {
	Model    string   `json:"model"`
	Strategy string   `json:"strategy"`
	Automine Automine `json:"automine"`
}

type Automine struct {
	Enabled      bool     `json:"enabled"`
	Version      int      `json:"version"`
	Model        string   `json:"model"`
	Fee          *big.Int `json:"fee"`
	Input        Input    `json:"input"`
	ModelAsBytes [32]byte `json:"-"`
}

type Blockchain struct {
	PrivateKey    string   `json:"private_key"`
	RPCURL        string   `json:"rpc_url"`         // main rpc used for blockchain streaming and sending transactions
	SenderRPCURL  string   `json:"sender_rpc_url"`  // if set, this rpc is used for sending txes
	ClientRPCURLs []string `json:"client_rpc_urls"` // array of urls to used to send same tx
	EthersGas     bool     `json:"use_ethers_gas_oracle"`
	CacheNonce    bool     `json:"cache_nonce"`
	BasefeeX      float64  `json:"basefee_x"` // basefee multiplier
	ForceGas      bool     `json:"gas_override"`
	GasOverride   float64  `json:"gas_override_gwei"`
}

type Redis struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type RPC struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Input struct {
	Prompt string `json:"prompt"`
	// NOT USED
	// NegativePrompt string `json:"negative_prompt"`
}

type Claimer struct {
	Enabled         bool    `json:"enabled"`
	NumberOfBatches int     `json:"number_of_batches"` // number of batches to send out at once
	MaxClaims       int     `json:"max_claims_per_batch"`
	MinClaims       int     `json:"min_claims_per_batch"`
	Delay           int     `json:"delay"`
	ValidateClaims  bool    `json:"validate_claims"`
	MinReward       float64 `json:"min_reward"`
	MaxGas          float64 `json:"max_claim_gas"`
	SortByCost      bool    `json:"sort_by_cost"`
	// Maximum amount of claims to buffer before submitting them regardless of min reward
	MaxClaimsBuffer int     `json:"max_claims_buffer"`
	ClaimMinReward  float64 `json:"claim_min_reward"` // if reward is this level claim regardless
	//  claim when staked amount approaches stake min level
	ClaimOnApproachMinStake bool    `json:"claim_on_approach"`
	MinStakeBufferLevel     float64 `json:"stake_buffer_level"`
	MinBatchProfit          float64 `json:"min_batch_profit"` // lowest profit to claim a batch
	// hoard claims
	HoardMode bool `json:"hoard_mode"`
	// max amount to hoard
	HoardMaxQueueSize  int      `json:"hoard_max_queue_size"`
	ClaimerPrivateKeys []string `json:"claimer_private_keys"`
}

type ML struct {
	Strategy  string         `json:"strategy"`
	Replicate Replicate      `json:"replicate"`
	Cog       map[string]Cog `json:"cog"`
}

type Replicate struct {
	APIToken string `json:"api_token"`
}

type Cog struct {
	URL []string `json:"url"`
}

type IPFS struct {
	Strategy   string     `json:"strategy"`
	HTTPClient HTTPClient `json:"http_client"`
}

type HTTPClient struct {
	URL string `json:"url"`
}

type CommitmentOption int

const (
	DoNothing CommitmentOption = iota
	DoBoth
	DoCommitmentsOnly
	DoSolutionsOnly
)

func (c CommitmentOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *CommitmentOption) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "donothing":
		*c = DoNothing
	case "doboth":
		*c = DoBoth
	case "docommitments":
		*c = DoCommitmentsOnly
	default:
		return errors.New("invalid CommitmentOption")
	}

	return nil
}

func (c CommitmentOption) String() string {
	switch c {
	case DoNothing:
		return "donothing"
	case DoBoth:
		return "doboth"
	case DoCommitmentsOnly:
		return "docommitments"
	default:
		return "unknown"
	}
}

func NewAppConfig(testnetType int) AppConfig {

	// TODO: fill with more sensible defaults
	cfg := AppConfig{
		DBPath:                "storage.db",
		LogPath:               "log.txt",
		LogLevel:              0,
		CachePath:             "cache",
		ListenToTaskSubmitted: true,
		CheckCommitment:       true,
		WhitelistTasks:        true,
		StealTasks:            false,
		NumWorkersPerGPU:      1,
		EvilMode:              false,
		EvilModeMinTime:       2000,
		EvilModeRandInt:       1000,

		ValidatorConfig: ValidatorConfig{
			InitialStake:            0,
			StakeBufferStakeAmount:  0,
			StakeBufferTopupAmount:  0,
			StakeBufferPercent:      2,
			StakeBufferTopupPercent: 1,
			StakeCheck:              true,
			StakeCheckInterval:      "120s",
			EthLowThreshold:         0.01,
			MinBasetokenThreshold:   1,
			SellInterval:            600,
			SellBuffer:              0.5,
			SellProfitInEth:         0.01,
			SellAllOverThreshold:    false,
			SellMinAmount:           0.001,
			SellMaxAmount:           0,
			SellEthBalanceTarget:    0,
			TreasurySplit:           0,
			TreasuryAddress:         [20]byte{},
			PrivateKeys:             []string{},
		},
		BatchTasks: BatchTasks{
			Enabled:                  false,
			MinTasksInQueue:          10,
			BatchMode:                "normal",
			BatchSize:                10,
			NumberOfBatches:          1,
			HoardMode:                false,
			HoardModeBatchSize:       10,
			HoardModeNumberOfBatches: 1,
			HoardMinGasPrice:         0.0,
			HoardMaxQueueSize:        1000,
		},
		Miner: SolverConfig{
			Enabled:                 false,
			CommitmentBatch:         BatchConfig{MinBatchSize: 10, MaxBatchSize: 10, NumberOfBatches: 1},
			SolutionBatch:           BatchConfig{MinBatchSize: 10, MaxBatchSize: 10, NumberOfBatches: 1},
			CommitmentsAndSolutions: DoBoth,
			ConcurrentBatches:       false,
			ProfitMode:              "fixed",
			MinProfit:               0,
			MaxProfit:               0,
			PauseStakeBufferLevel:   0,
			UsePolling:              true,
			PollingTime:             "1m",
			BatchMode:               1,
			ErrorMaxRetries:         5,
			ErrorBackoffTime:        425,
			ErrorBackofMultiplier:   1.5,
			MetricsSampleRate:       "10s",
		},
		Blockchain: Blockchain{
			PrivateKey: "",
			RPCURL:     "http://localhost:8545",
			EthersGas:  false,
			CacheNonce: false,
			BasefeeX:   2.0,
		},
		Claim: Claimer{
			Enabled:                 true,
			MaxClaims:               50,
			MinClaims:               10,
			Delay:                   60,
			ValidateClaims:          false,
			MinReward:               0.0,
			MaxClaimsBuffer:         2000,
			ClaimOnApproachMinStake: false,
			MinStakeBufferLevel:     0,
			HoardMode:               false,
			HoardMaxQueueSize:       0,
		},
		Strategies: Strategies{
			Strategy: "nop",
		},
	}

	data := baseConfigJsonData

	switch testnetType {
	case 1:
		data = baseConfigJsonDataLocal
	case 2:
		data = baseConfigJsonDataTestnet
	}

	//var cfg BaseConfig
	err := json.Unmarshal([]byte(data), &cfg.BaseConfig)
	if err != nil {
		panic("failed to unmarshal base config: " + err.Error())
	}

	cfg.BaseConfig.BaseToken = erc20.NewTokenERC20(cfg.BaseConfig.BaseTokenAddress, 18, "AIUS", "AIUS")

	return cfg
}

var (
	ErrNoAutomineModel      = errors.New("automine model not set")
	ErrAutomineModelInvalid = errors.New("automine model could not be converted to [32]byte")
)

// TODO: some basic validation of the values
func InitAppConfig(file string, testnetType int) (*AppConfig, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := NewAppConfig(testnetType)

	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Strategies.Automine.Model == "" {
		return nil, ErrNoAutomineModel
	}

	cfg.Strategies.Automine.ModelAsBytes, err = common.ConvertTaskIdString2Bytes(cfg.Strategies.Automine.Model)
	if err != nil {
		return nil, ErrAutomineModelInvalid
	}

	// load .env file if it exists to set sensible defaults or overrides
	err = godotenv.Load("")
	if err != nil {
		log.Warn("Warning: no .env found: ", err)
		//return nil, err
	}

	rpcURL := os.Getenv("RPC_URL")
	if rpcURL != "" {
		cfg.Blockchain.RPCURL = rpcURL
	}

	return &cfg, nil
}

func (cfg *AppConfig) ExportConfig(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(cfg)
}
