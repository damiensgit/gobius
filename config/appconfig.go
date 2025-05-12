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
	DBPath                 string            `json:"db_path"`
	LogPath                string            `json:"log_path"`
	LogLevel               int               `json:"log_level"`
	CachePath              string            `json:"cache_path"`
	CheckCommitment        bool              `json:"check_commitment"`
	DryRunMode             bool              `json:"dry_run_mode"`
	EvilMode               bool              `json:"evil_mode"`     // for testing purposes
	EvilModeMinTime        int               `json:"evil_mode_min"` // for testing purposes
	EvilModeRandInt        int               `json:"evil_mode_int"` // for testing purposes
	NumWorkersPerGPU       int               `json:"num_workers_per_gpu"`
	PriceOracleContract    ethcommon.Address `json:"price_oracle_contract"`
	PopTaskRandom          bool              `json:"pop_task_random"`
	VerificationSampleRate int               `json:"verification_sample_rate"`

	Miner SolverConfig `json:"solver"`

	TelegramBot TelegramBot `json:"telegram"` // not used

	ValidatorConfig ValidatorConfig `json:"validator_config"`

	BatchTasks BatchTasks `json:"batchtasks"`
	Strategies Strategies `json:"strategies"`
	Blockchain Blockchain `json:"blockchain"`
	Claim      Claimer    `json:"claim"`
	ML         ML         `json:"ml"`
	IPFS       IPFS       `json:"ipfs"`
	BaseConfig BaseConfig `json:"baseconfig"`

	ParaswapCacheTTL string `json:"paraswap_cache_ttl"`
	ParaswapTimeout  string `json:"paraswap_timeout"`
}

type TelegramBot struct {
	Enabled bool   `json:"enabled"`
	Token   string `json:"token"`
	ChatID  string `json:"chat_id"`
}

type ValidatorConfig struct {
	// InitialStake specifies the minimum amount to stake for a new validator (in tokens)
	// If > 0 and validator's stake is below this value, stake will be topped up to this amount
	// takes precedence over stake_buffer_amount and stake_buffer_percent
	// only useful for initial stake, not for topups e.g. testing purposes
	// it is recommended to set this to 0 and use the stake_buffer_amount and stake_buffer_percent instead
	InitialStake float64 `json:"initial_stake"`

	// StakeBufferStakeAmount is an absolute stake amount (in tokens) to maintain above minimum
	// If > 0, this overrides the percentage-based buffer mechanism
	// System will top up to InitialStake or ensure StakeBufferStakeAmount above minimum
	StakeBufferStakeAmount float64 `json:"stake_buffer_amount"`
	// StakeBufferTopupAmount is the amount to add when topping up with StakeBufferStakeAmount
	StakeBufferTopupAmount float64 `json:"stake_buffer_topup_amount"`

	// StakeBufferPercent specifies minimum buffer as percentage above required stake
	// For example, 10 means maintain at least 10% more than minimum stake
	StakeBufferPercent int `json:"stake_buffer_percent"`
	// StakeBufferTopupPercent specifies when to trigger a top-up as percentage above minimum
	// For example, if StakeBufferPercent=10 and StakeBufferTopupPercent=5,
	// top-up occurs when stake falls below 5% buffer and restores to 10% buffer
	// making these values the same will lead to constant topups so keep some distance between them
	StakeBufferTopupPercent int `json:"stake_buffer_topup_percent"`

	// StakeCheck specifies if we should check the validator's stake, stake_check_interval is the polling interval
	StakeCheck bool `json:"stake_check"`
	// Interval for checking the validator's stake (e.g. "120s", "1m", "5m", "1h", etc.)
	StakeCheckInterval string `json:"stake_check_interval"`

	// If the validator has less than this amount of eth, we send log warnings
	EthLowThreshold float64 `json:"eth_low_threshold"`

	// Minimum balance to leave on the validator account at all times regardless of the sell settings
	// Use this to ensure we have enough balance to cover validator min stake topups in future
	MinBasetokenThreshold float64 `json:"min_basetoken_threshold"`

	// Auto sell settings (TODO: break out into own struct)
	// Auto sell check interval, in seconds
	// by default, auto sell will sell enough aius to cover the gas used since the last sell
	SellInterval int `json:"sell_interval"`
	// Extra buffer to sell, e.g. 1.5 means sell 1.5 times the amount of tokens
	SellBuffer float64 `json:"sell_buffer"`
	// An additional amount of AIUS in Eth terms to sell
	SellProfitInEth float64 `json:"sell_profit_in_eth"`
	// if we should sell all tokens if the token balance is over the threshold
	SellAllOverThreshold bool `json:"sell_all_over_threshold"`
	// Minimum amount of tokens to sell
	SellMinAmount float64 `json:"sell_min_amount"`
	// Maximum amount of tokens to sell
	SellMaxAmount float64 `json:"sell_max_amount"`
	// If this value is > 0 then we want to ensure the balance of ETH reaches this target
	// so we weant to keep selling all the AIUS until we reach this target (over the min threshold)
	SellEthBalanceTarget float64 `json:"sell_eth_bal_target"`

	// TreasurySplit is the percentage of the rewards to send to the treasury
	TreasurySplit float64 `json:"treasury_split"`
	// TreasuryAddress is the address of the treasury
	TreasuryAddress ethcommon.Address `json:"treasury_address"`

	// PrivateKeys is a list of 1 or more private keys to use for each validator
	PrivateKeys []string `json:"private_keys"`
}

type BatchTasks struct {
	Enabled                  bool     `json:"enabled"`            // enable batch task creation
	MinTasksInQueue          int      `json:"min_tasks_in_queue"` // the number of tasks left in the queue before we start creating new batches
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
	MaxClaimQueue            int      `json:"max_claim_queue"` // max size we let the claim queue get before disabling task creation
}

type BatchConfig struct {
	MinBatchSize    int `json:"min_batch_size"`    // minimum number of tasks to include in a batch
	MaxBatchSize    int `json:"max_batch_size"`    // maximum number of tasks to include in a batch
	NumberOfBatches int `json:"number_of_batches"` // number of batches to send out at once
}

type SolverConfig struct {
	WaitForTasksOnShutdown          bool             `json:"wait_for_tasks_on_shutdown"` // If true, allow running tasks to finish on shutdown
	CommitmentsAndSolutions         CommitmentOption `json:"commitments_and_solutions"`  // one of: "donothing", "doboth", "docommitments", "dosolutions"
	CommitmentBatch                 BatchConfig      `json:"commitment_batch"`
	SolutionBatch                   BatchConfig      `json:"solution_batch"`
	ConcurrentBatches               bool             `json:"concurrent_batches"`                 // if true, submit multiple batches of commitments and solutions concurrently (requires multiple accounts)
	ProfitMode                      string           `json:"profit_mode"`                        // profit mode to use for batch operations
	MinProfit                       float64          `json:"min_profit"`                         // minimum profit in USD to perform batch operations
	MaxProfit                       float64          `json:"max_profit"`                         // maximum profit in USD to perform batch operations
	PauseStakeBufferLevel           float64          `json:"pause_stake_buffer_level"`           // pause submitting commits/solutions if the buffer = (current stake - min stake) is below this level
	UsePolling                      bool             `json:"use_polling"`                        // use polling to check profit and submit txes if false, new block triggers batching
	PollingTime                     string           `json:"polling_time"`                       // polling interval for profit checks and batching as "1m", "5m", "1h", etc..
	NoChecks                        bool             `json:"no_checks"`                          // perform no onchain checks on the tasks/commitments/solutions, just submit them
	ErrorMaxRetries                 int              `json:"error_max_retries"`                  // max retries for tx errors
	ErrorBackoffTime                float64          `json:"error_backoff"`                      // sleep time between retries
	ErrorBackofMultiplier           float64          `json:"error_backoff_multiplier"`           // backoff multiplier tx errors
	MetricsSampleRate               string           `json:"metrics_sample_rate"`                // sample rate for metrics
	EnableIntrinsicGasCheck         bool             `json:"enable_intrinsic_gas_check"`         // enable intrinsic gas check
	IntrinsicGasBaseline            uint64           `json:"intrinsic_gas_baseline"`             // baseline gas cost for a simple transfer
	IntrinsicGasThresholdMultiplier float64          `json:"intrinsic_gas_threshold_multiplier"` // multiplier for the threshold
	EnableGasEstimationMode         bool             `json:"enable_gas_estimation_mode"`         // enable gas estimation mode
	GasEstimationMargin             uint64           `json:"gas_estimation_margin"`              // gas margin for gas estimation mode
}

type Strategies struct {
	Model    string   `json:"model"`
	Strategy string   `json:"strategy"`
	Automine Automine `json:"automine"`
}

type Automine struct {
	Version      int               `json:"version"`
	Model        string            `json:"model"`
	Fee          *big.Int          `json:"fee"`
	Input        Input             `json:"input"`
	Owner        ethcommon.Address `json:"owner"`
	ModelAsBytes [32]byte          `json:"-"`
}

type Blockchain struct {
	PrivateKey    string   `json:"private_key"`
	RPCURL        string   `json:"rpc_url"`               // main rpc used for blockchain streaming and sending transactions
	SenderRPCURL  string   `json:"sender_rpc_url"`        // if set, this rpc is used for sending txes
	ClientRPCURLs []string `json:"client_rpc_urls"`       // array of urls to used to send same tx
	EthersGas     bool     `json:"use_ethers_gas_oracle"` // use same gas oracle as ethers.js
	CacheNonce    bool     `json:"cache_nonce"`           // cache the nonce for the sender account
	BasefeeX      float64  `json:"basefee_x"`             // basefee multiplier
	ForceGas      bool     `json:"gas_override"`          // force the use of a specific gas price
	GasOverride   float64  `json:"gas_override_gwei"`     // gas price to use if gas override is enabled in gwei
}

type Input struct {
	Prompt string `json:"prompt"`
	// NOT USED
	// NegativePrompt string `json:"negative_prompt"`
}

type Claimer struct {
	Enabled         bool    `json:"enabled"`              // solution claimer enabled
	NumberOfBatches int     `json:"number_of_batches"`    // number of batches to send out at once
	MaxClaims       int     `json:"max_claims_per_batch"` // maximum number of claims per batch
	MinClaims       int     `json:"min_claims_per_batch"` // minimum number of claims per batch
	Delay           int     `json:"delay"`                // delay between claims in seconds
	MaxGas          float64 `json:"max_claim_gas"`        // maximum gas to use for a claim
	SortByCost      bool    `json:"sort_by_cost"`         // sort the claims by cost
	// Maximum amount of claims to buffer before submitting them regardless of min reward
	//MaxClaimsBuffer int     `json:"max_claims_buffer"`
	UseLever       bool    `json:"use_lever"`        // use the lever oracle for the claim min level
	ClaimMinReward float64 `json:"claim_min_reward"` // if reward is this level claim regardless
	//  claim when staked amount approaches stake min level
	ClaimOnApproachMinStake bool    `json:"claim_on_approach"`
	MinStakeBufferLevel     float64 `json:"stake_buffer_level"`
	MinBatchProfit          float64 `json:"min_batch_profit"`     // lowest profit to claim a batch
	HoardMode               bool    `json:"hoard_mode"`           // hoard claims when gas price is low
	HoardMaxQueueSize       int     `json:"hoard_max_queue_size"` // max amount to hoard
}

type ML struct {
	Strategy string         `json:"strategy"`
	Cog      map[string]Cog `json:"cog"`
}

type Cog struct {
	// HttpTimeout specifies the duration allowed for a model inference request (e.g., "120s", "5m").
	HttpTimeout string `json:"http_timeout"`
	// IpfsTimeout specifies the duration allowed for IPFS pinning operations (e.g., "30s").
	IpfsTimeout string   `json:"ipfs_timeout"`
	URL         []string `json:"url"`
}

type IPFS struct {
	Strategy                  IpfsStrategy `json:"strategy"`        // mock, http_client, pinata_client, mixed_client
	HTTPClient                HTTPClient   `json:"http_client"`     // http client to use for ipfs pinning
	Pinata                    Pinata       `json:"pinata"`          // pinata client to use for ipfs pinning
	IncentiveClaim            bool         `json:"incentive_claim"` // set to true to claim the incentive for pinning ipfs content
	ClaimInterval             string       `json:"claim_interval"`  // how often to claim the incentive for pinning ipfs content
	OracleURL                 string       `json:"oracle_url"`      // oracle url to use for incentive claim
	Timeout                   string       `json:"timeout"`
	UseBulkClaim              bool         `json:"use_bulk_claim"`               // New: Option to use bulk IPFS claims
	BulkClaimBatchSize        int          `json:"bulk_claim_batch_size"`        // New: Batch size for bulk IPFS claims
	MaxSingleClaimsPerRun     int          `json:"max_single_claims_per_run"`    // Max single claims per processing run
	MinAiusIncentiveThreshold float64      `json:"min_aius_incentive_threshold"` // Minimum AIUS value to claim incentive (0 = disabled)
}

type Pinata struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
	JWT       string `json:"jwt"`
	BaseURL   string `json:"base_url"`
}

type HTTPClient struct {
	URL string `json:"url"`
}

type IpfsStrategy int

const (
	MockClient IpfsStrategy = iota
	HttpClient
	PinataClient
	MixedClient
)

func (c IpfsStrategy) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *IpfsStrategy) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "mock":
		*c = MockClient
	case "http_client":
		*c = HttpClient
	case "pinata_client":
		*c = PinataClient
	case "mixed_client":
		*c = MixedClient
	default:
		return errors.New("invalid IpfsStrategy")
	}

	return nil
}

func (c IpfsStrategy) String() string {
	switch c {
	case MockClient:
		return "mock"
	case HttpClient:
		return "http_client"
	case PinataClient:
		return "pinata_client"
	case MixedClient:
		return "mixed_client"
	default:
		return "unknown"
	}
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
		DBPath:                 "storage.db",
		LogPath:                "log.txt",
		LogLevel:               0,
		CachePath:              "cache",
		CheckCommitment:        true,
		NumWorkersPerGPU:       1,
		DryRunMode:             false,
		EvilMode:               false,
		EvilModeMinTime:        2000,
		EvilModeRandInt:        1000,
		PopTaskRandom:          false,
		VerificationSampleRate: 0,

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
			PrivateKeys:              []string{},
			MaxClaimQueue:            0, // 0 disables this check
		},
		Miner: SolverConfig{
			WaitForTasksOnShutdown:  false,
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
			ErrorMaxRetries:         5,
			ErrorBackoffTime:        425,
			ErrorBackofMultiplier:   1.5,
			MetricsSampleRate:       "60s",
			EnableGasEstimationMode: true,
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
			ClaimOnApproachMinStake: false,
			MinStakeBufferLevel:     0,
			HoardMode:               false,
			HoardMaxQueueSize:       0,
		},
		Strategies: Strategies{
			Strategy: "nop",
		},
		IPFS: IPFS{
			Strategy:       MockClient,
			IncentiveClaim: false,
			OracleURL:      "",
			Timeout:        "10s",
		},

		// Default Paraswap settings
		ParaswapCacheTTL: "5m",
		ParaswapTimeout:  "30s",
	}

	data := baseConfigJsonDataMainnet

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

	cfg.BaseConfig.TestnetType = testnetType

	cfg.BaseConfig.BaseToken = erc20.NewTokenERC20(cfg.BaseConfig.BaseTokenAddress, 18, "AIUS", "AIUS")

	return cfg
}

var (
	ErrNoAutomineModel      = errors.New("automine model not set")
	ErrAutomineModelInvalid = errors.New("automine model could not be converted to [32]byte")
	ErrNoAutomineOwner      = errors.New("automine owner not set")
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

	if cfg.Strategies.Automine.Owner == (ethcommon.Address{}) {
		return nil, ErrNoAutomineOwner
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

func LoadConfigForTesting(file string, testnetType int) (*AppConfig, error) {
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

	return &cfg, nil
}
