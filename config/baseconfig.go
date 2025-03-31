package config

import (
	_ "embed"
	"gobius/erc20"

	"github.com/ethereum/go-ethereum/common"
)

type BaseConfig struct {
	BaseTokenAddress    common.Address    `json:"baseTokenAddress"`
	EngineAddress       common.Address    `json:"engineAddress"`
	VoterAddress        common.Address    `json:"voterAddress"`
	VeStakingAddress    common.Address    `json:"veStakingAddress"`
	BulkTasksAddress    common.Address    `json:"bulkTasksAddress"`
	ArbiusRouterAddress common.Address    `json:"arbiusRouterAddress"`
	Models              map[string]Model  `json:"models"`
	BaseToken           *erc20.TokenERC20 `json:"-"`
	TestnetType         int               `json:"-"`
}

type Model struct {
	ID        string                    `json:"id"`
	Mineable  bool                      `json:"mineable"`
	Contracts map[string]common.Address `json:"contracts"`
	Params    ModelParams               `json:"params"`
}

type ModelParams struct {
	Addr common.Address `json:"addr"`
	Fee  string         `json:"fee"`
	Rate string         `json:"rate"`
	Cid  string         `json:"cid"`
}

//go:embed config.json
var baseConfigJsonData string

//go:embed config.local.json
var baseConfigJsonDataLocal string

//go:embed config.testnet.json
var baseConfigJsonDataTestnet string

// func InitBaseConfig() (*BaseConfig, error) {

// 	var cfg BaseConfig
// 	err := json.Unmarshal([]byte(baseConfigJsonData), &cfg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	cfg.BaseToken = erc20.NewTokenERC20(cfg.BaseTokenAddress, 18, "AIUS", "AIUS")

// 	return &cfg, nil
// }
