package models

import (
	"gobius/common"
	"gobius/config"
	"gobius/ipfs"
)

type Model struct {
	ID       string
	Template interface{}
	Mineable bool
	Filters  []MiningFilter
}

// TODO: add context support to GetFiles and GetCID
type ModelInterface interface {
	GetFiles(gpu *common.GPU, taskid string, input interface{}) ([]ipfs.IPFSFile, error)
	GetCID(gpu *common.GPU, taskid string, input interface{}) ([]byte, error)
	GetID() string
	HydrateInput(preprocessedInput map[string]interface{}) (InputHydrationResult, error)
}

type MiningFilter struct {
	MinFee  int
	MinTime int
}

type EnabledModel struct {
	Models []ModelInterface
}

func (em *EnabledModel) AddModel(m ModelInterface) {
	em.Models = append(em.Models, m)
}

func (em *EnabledModel) FindModel(id string) ModelInterface {
	// add 0x to the id if it's not there
	if id[:2] != "0x" {
		id = "0x" + id
	}
	for _, m := range em.Models {
		if m.GetID() == id {
			return m
		}
	}
	return nil
}

var EnabledModels EnabledModel = EnabledModel{}

// This will panic on config validation if the model is not found in the config
func InitEnabledModels(client ipfs.IPFSClient, config *config.AppConfig) {
	model := NewKandinsky2Model(client, config)
	EnabledModels.AddModel(model)
}
