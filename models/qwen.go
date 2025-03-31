package models

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gobius/common"
	"gobius/config"
	"gobius/ipfs"
	"gobius/utils"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/mr-tron/base58"
	"github.com/rs/zerolog"
)

type QwenInner struct {
	Prompt string `json:"prompt"`
	Seed   uint64 `json:"seed"`
}

type QwenPrompt struct {
	Input QwenInner `json:"input"`
}

type QwenModelResponse struct {
	Input  map[string]any `json:"input"`
	Output []string       `json:"output"`
}

type QwenTestModel struct {
	Model
	Filters []MiningFilter
	config  *config.AppConfig
	client  *http.Client
	logger  zerolog.Logger
	ipfs    ipfs.IPFSClient
}

// Ensure QwenTestModel implements the Model interface.
var _ ModelInterface = (*QwenTestModel)(nil)

var QwenTestModelTemplate = Model{
	ID:       "",
	Mineable: true,
	Template: map[string]any{
		"meta": map[string]any{
			"title":       "Qwen",
			"description": "Qwen Sepolia Test Model",
			"version":     1,
			"input": []map[string]any{
				{
					"variable":    "prompt",
					"type":        "string",
					"required":    true,
					"default":     "",
					"description": "Input prompt",
				},
				{
					"variable":    "seed",
					"type":        "int",
					"required":    false,
					"default":     500,
					"description": "Seed for the random number generator.",
				},
			},
			"output": []map[string]any{
				{
					"filename": "out-1.txt",
					"type":     "text",
				},
			},
		},
	},
}

func NewQwenTestModel(client ipfs.IPFSClient, appConfig *config.AppConfig, logger zerolog.Logger) *QwenTestModel {

	model, ok := appConfig.BaseConfig.Models["qwen-test"]
	if !ok {
		return nil
	}

	if model.ID == "" {
		logger.Error().Str("model", "qwen").Msg("qwen model ID is empty")
		return nil
	}

	http := &http.Client{
		Timeout: time.Second * 30,
	}

	m := &QwenTestModel{
		Model:  QwenTestModelTemplate,
		config: appConfig,
		//url:    url[0],
		Filters: []MiningFilter{
			{
				MinFee:  0,
				MinTime: 0,
			},
		},
		ipfs:   client,
		client: http,
		logger: logger,
	}
	// set this from config for now
	m.Model.ID = model.ID
	return m
}

func (m *QwenTestModel) HydrateInput(preprocessedInput map[string]any, seed uint64) (InputHydrationResult, error) {
	input := make(map[string]any)

	// Helper functions for type conversion
	convertToInt := func(val any) (int, error) {
		switch v := val.(type) {
		case int:
			return v, nil
		case float64:
			return int(v), nil
		default:
			return 0, fmt.Errorf("cannot convert %T to int", val)
		}
	}

	convertToFloat := func(val any) (float64, error) {
		switch v := val.(type) {
		case float64:
			return v, nil
		case int:
			return float64(v), nil
		default:
			return 0, fmt.Errorf("cannot convert %T to float64", val)
		}
	}

	// Get template metadata for input validation
	templateMeta, ok := m.Model.Template.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid template format")
	}

	meta, ok := templateMeta["meta"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid meta format in template")
	}

	inputFields, ok := meta["input"].([]map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid input format in template meta")
	}

	// Process each input field according to template rules
	for _, field := range inputFields {
		varName := field["variable"].(string)
		fieldType := field["type"].(string)
		required, _ := field["required"].(bool)

		// Get value from input or use default
		value, exists := preprocessedInput[varName]

		// Check if required field is missing
		if required && !exists {
			return nil, fmt.Errorf("input missing required field (%s)", varName)
		}

		// If value exists, validate it
		if exists {
			// Validate type
			if err := validateType(value, fieldType, varName); err != nil {
				return nil, err
			}

			// Validate range for numeric types
			switch fieldType {
			case "int":
				intVal, err := convertToInt(value)
				if err != nil {
					return nil, err
				}

				// Check min if defined
				if minVal, ok := field["min"].(int); ok {
					if intVal < minVal {
						return nil, fmt.Errorf("input below minimum (%s): value %d, minimum %d",
							varName, intVal, minVal)
					}
				}

				// Check max if defined
				if maxVal, ok := field["max"].(int); ok {
					if intVal > maxVal {
						return nil, fmt.Errorf("input above maximum (%s): value %d, maximum %d",
							varName, intVal, maxVal)
					}
				}

			case "decimal":
				floatVal, err := convertToFloat(value)
				if err != nil {
					return nil, err
				}

				// Check min if defined
				if minVal, ok := field["min"].(float64); ok {
					if floatVal < minVal {
						return nil, fmt.Errorf("input below minimum (%s): value %f, minimum %f",
							varName, floatVal, minVal)
					}
				}

				// Check max if defined
				if maxVal, ok := field["max"].(float64); ok {
					if floatVal > maxVal {
						return nil, fmt.Errorf("input above maximum (%s): value %f, maximum %f",
							varName, floatVal, maxVal)
					}
				}

			case "string_enum", "int_enum":
				if err := validateEnum(value, field, varName, convertToInt); err != nil {
					return nil, err
				}
			}

			input[varName] = value
		} else {
			// Use default value if provided
			input[varName] = field["default"]
		}
	}

	// Convert validated input to the expected QwenInner format
	var inner QwenInner
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	if err := json.Unmarshal(jsonBytes, &inner); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to QwenInner: %w", err)
	}

	// TODO: probable a better way to handle values we need to set
	inner.Seed = seed

	return QwenPrompt{Input: inner}, nil
}

func (m *QwenTestModel) GetID() string {
	return m.Model.ID
}

func (m *QwenTestModel) GetFiles(gpu *common.GPU, taskid string, input any) ([]ipfs.IPFSFile, error) {

	marshaledInput, _ := json.Marshal(input)

	postResp, err := m.client.Post(gpu.Url, "application/json", bytes.NewBuffer([]byte(marshaledInput)))
	if err != nil {
		return nil, err
	}
	defer postResp.Body.Close()

	body, err := io.ReadAll(postResp.Body)
	if err != nil {
		return nil, err
	}

	var resp QwenModelResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Output) != 1 {
		return nil, err
	}

	fileName := fmt.Sprintf("%d.%s.txt", gpu.ID, uuid.New().String())
	path := filepath.Join(m.config.CachePath, fileName)
	buffer := bytes.NewBufferString(resp.Output[0])

	return []ipfs.IPFSFile{{Name: "out-1.txt", Path: path, Buffer: buffer}}, nil
}

func (m *QwenTestModel) GetCID(gpu *common.GPU, taskid string, input any) ([]byte, error) {
	paths, err := utils.ExpRetry(m.logger, func() (any, error) {
		return m.GetFiles(gpu, taskid, input)
	}, 3, 1000)
	if err != nil {
		return nil, err
	}

	cid58, err := utils.ExpRetry(m.logger, func() (any, error) {
		return m.ipfs.PinFilesToIPFS(taskid, paths.([]ipfs.IPFSFile))
	}, 3, 1000)

	if err != nil {
		return nil, errors.New("cannot pin files to retrieve cid")
	}

	cidBytes, err := base58.Decode(cid58.(string))
	if err != nil {
		return nil, err
	}

	return cidBytes, nil
}

func (m *QwenTestModel) Validate(gpu *common.GPU, taskid string) error {

	testPrompt := QwenPrompt{
		Input: QwenInner{
			Prompt: "why is the meaning of life 42?",
			Seed:   1337,
		},
	}

	cid, err := m.GetCID(gpu, "startup-test-taskid", testPrompt)
	if err != nil {
		return err
	}

	expected := "0x122041fa6dba0027cf73097b4c75cae5f16122b4f1cfe85616bdccb6eaea711e7238"
	cidStr := "0x" + hex.EncodeToString(cid)
	if cidStr == expected {
		m.logger.Info().Str("model", m.GetID()).Str("cid", cidStr).Str("expected", expected).Msg("model CID matches expected CID")
	} else {
		m.logger.Error().Str("model", m.GetID()).Str("cid", cidStr).Str("expected", expected).Msg("model CID does not match expected CID")
		return errors.New("model CID does not match expected CID")
	}

	return nil
}
