package models

import (
	"context"
	"errors"
	"strings"

	"gobius/common"
	"gobius/config"
	"gobius/ipfs"

	"github.com/rs/zerolog"
)

type Model struct {
	ID       string
	Template any
	Mineable bool
	Filters  []MiningFilter
}

type InputHydrationResult any

type ModelInterface interface {
	GetFiles(ctx context.Context, gpu *common.GPU, taskid string, input any) ([]ipfs.IPFSFile, error)
	GetCID(ctx context.Context, gpu *common.GPU, taskid string, input any) ([]byte, error)
	GetID() string
	HydrateInput(preprocessedInput map[string]any, seed uint64) (InputHydrationResult, error)
	Validate(gpu *common.GPU, taskid string) error
}

var ErrResourceBusy = errors.New("resource busy")

type MiningFilter struct {
	MinFee  int
	MinTime int
}

// ModelFactory manages model registration and retrieval
type ModelFactory struct {
	registeredModels map[string]ModelInterface
}

// NewModelFactory creates a new model factory
func NewModelFactory() *ModelFactory {
	return &ModelFactory{
		registeredModels: make(map[string]ModelInterface),
	}
}

// RegisterModel adds a model to the factory
func (mf *ModelFactory) RegisterModel(model ModelInterface) {
	modelID := model.GetID()
	mf.registeredModels[modelID] = model
}

// GetModel retrieves a model by ID
func (mf *ModelFactory) GetModel(id string) ModelInterface {
	if !strings.HasPrefix(id, "0x") {
		id = "0x" + id
	}
	return mf.registeredModels[id]
}

// GetAllModels returns all registered models
func (mf *ModelFactory) GetAllModels() []ModelInterface {
	models := make([]ModelInterface, 0, len(mf.registeredModels))
	for _, model := range mf.registeredModels {
		models = append(models, model)
	}
	return models
}

// Global factory instance
var ModelRegistry *ModelFactory

// InitModelRegistry initializes the model registry with available models
func InitModelRegistry(client ipfs.IPFSClient, config *config.AppConfig, logger zerolog.Logger) {
	ModelRegistry = NewModelFactory()

	// Register Qwen Mainnet
	modelQwenMainnet := NewQwenMainnetModel(client, config, logger)
	if modelQwenMainnet != nil {
		ModelRegistry.RegisterModel(modelQwenMainnet)
	}

	// Register Qwen Testnet
	modelQwenTest := NewQwenTestModel(client, config, logger)
	if modelQwenTest != nil {
		ModelRegistry.RegisterModel(modelQwenTest)
	}

	// Register Kandinsky2
	modelKandinsky2 := NewKandinsky2Model(client, config, logger)
	if modelKandinsky2 != nil {
		ModelRegistry.RegisterModel(modelKandinsky2)
	}

	// Register Metabaron-Uncensored-8B
	modelMetabaron := NewMetabaronModel(client, config, logger)
	if modelMetabaron != nil {
		ModelRegistry.RegisterModel(modelMetabaron)
	}

}
