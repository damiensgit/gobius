package models

import (
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

// TODO: add context support to GetFiles and GetCID
type ModelInterface interface {
	GetFiles(gpu *common.GPU, taskid string, input any) ([]ipfs.IPFSFile, error)
	GetCID(gpu *common.GPU, taskid string, input any) ([]byte, error)
	GetID() string
	HydrateInput(preprocessedInput map[string]any, seed uint64) (InputHydrationResult, error)
	Validate(gpu *common.GPU, taskid string) error
}

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
	// Normalize ID format (add 0x prefix if missing)
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
func InitModelRegistry(client ipfs.IPFSClient, config *config.AppConfig, logger *zerolog.Logger) {
	ModelRegistry = NewModelFactory()

	// Register available models
	modelQwen := NewQwenTestModel(client, config, logger)
	ModelRegistry.RegisterModel(modelQwen)

	modelKandinsky2 := NewKandinsky2Model(client, config, logger)
	ModelRegistry.RegisterModel(modelKandinsky2)

	// Register additional models here as needed
}
