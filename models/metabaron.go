package models

import (
	"bytes"
	"context"
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

type MetabaronInner struct {
	Prompt string `json:"prompt"`
	Seed   uint64 `json:"seed"`
}

type MetabaronPrompt struct {
	Input MetabaronInner `json:"input"`
}

type MetabaronModelResponse struct {
	Output string `json:"output"`
}

type MetabaronModel struct {
	Model
	timeoutDuration     time.Duration
	ipfsTimeoutDuration time.Duration
	Filters             []MiningFilter
	config              *config.AppConfig
	client              *http.Client
	logger              zerolog.Logger
	ipfs                ipfs.IPFSClient
}

var _ ModelInterface = (*MetabaronModel)(nil)

var MetabaronModelTemplate = Model{
	ID:       "",
	Mineable: true,
	Template: map[string]any{
		"meta": map[string]any{
			"title":       "Metabaron-Uncensored-8B",
			"description": "A Cog-packaged LLaMA 3.1 8B model fine-tuned for uncensored text generation.",
			"version":     1,
			"input": []map[string]any{
				{
					"variable":    "prompt",
					"type":        "string",
					"required":    true,
					"default":     "",
					"description": "Text prompt",
				},
				{
					"variable":    "seed",
					"type":        "int",
					"required":    false,
					"default":     42,
					"description": "Seed for reproducibility.",
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

func NewMetabaronModel(client ipfs.IPFSClient, appConfig *config.AppConfig, logger zerolog.Logger) *MetabaronModel {
	model, ok := appConfig.BaseConfig.Models["metabaron"]
	if !ok || model.ID == "" {
		logger.Error().Str("model", "metabaron").Msg("Model not configured or missing ID")
		return nil
	}

	cogConfig, ok := appConfig.ML.Cog[model.ID]
	if !ok {
		logger.Error().Str("model", model.ID).Msg("model ID not found in ML.Cog config")
		return nil
	}

	http := &http.Client{
		Transport: &http.Transport{MaxIdleConnsPerHost: 10}, // Use a dedicated transport
	}

	timeout := 120 * time.Second
	ipfsTimeout := 30 * time.Second

	if cogConfig.HttpTimeout != "" {
		if t, err := time.ParseDuration(cogConfig.HttpTimeout); err == nil {
			timeout = t
		}
	}

	if cogConfig.IpfsTimeout != "" {
		if t, err := time.ParseDuration(cogConfig.IpfsTimeout); err == nil {
			ipfsTimeout = t
		}
	}

	m := &MetabaronModel{
		Model:               MetabaronModelTemplate,
		timeoutDuration:     timeout,
		ipfsTimeoutDuration: ipfsTimeout,
		config:              appConfig,
		Filters: []MiningFilter{{
			MinFee:  0,
			MinTime: 0,
		}},
		ipfs:   client,
		client: http,
		logger: logger,
	}

	m.Model.ID = model.ID
	return m
}

func (m *MetabaronModel) HydrateInput(preprocessedInput map[string]any, seed uint64) (InputHydrationResult, error) {
	promptVal, ok := preprocessedInput["prompt"]
	if !ok {
		return nil, errors.New("input missing required field 'prompt'")
	}
	promptStr, ok := promptVal.(string)
	if !ok {
		return nil, fmt.Errorf("input field 'prompt' is not a string, got %T", promptVal)
	}
	input := MetabaronInner{
		Prompt: promptStr, // Use the validated string
		Seed:   seed,
	}
	return MetabaronPrompt{Input: input}, nil

}

func (m *MetabaronModel) GetID() string {
	return m.Model.ID
}

func (m *MetabaronModel) GetFiles(ctx context.Context, gpu *common.GPU, taskid string, input any) ([]ipfs.IPFSFile, error) {
	marshaledInput, _ := json.Marshal(input)

	req, err := http.NewRequestWithContext(ctx, "POST", gpu.Url, bytes.NewBuffer(marshaledInput))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			m.logger.Error().Err(err).Str("task", taskid).Str("gpu", gpu.Url).Msg("model inference request timed out")
			return nil, fmt.Errorf("model inference timed out: %w", err)
		}
		return nil, fmt.Errorf("failed to POST to GPU: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusConflict {
			m.logger.Warn().Str("task", taskid).Str("gpu", gpu.Url).Int("status", resp.StatusCode).Str("body", string(bodyBytes)).Msg("resource busy")
			return nil, fmt.Errorf("resource busy (HTTP 409)")
		}
		return nil, fmt.Errorf("server returned non-200 status: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read model response body: %w", err)
	}

	var result MetabaronModelResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal model response: %w", err)
	}

	if result.Output == "" {
		m.logger.Error().
			Str("taskid", taskid).
			Str("gpu", gpu.Url).
			Str("response", string(body)).
			Msg("Empty output received from GPU")
		return nil, fmt.Errorf("model returned empty text data")
	}

	filename := fmt.Sprintf("%d.%s.txt", gpu.ID, uuid.New().String())
	path := filepath.Join(m.config.CachePath, filename)
	buffer := bytes.NewBufferString(result.Output)

	m.logger.Debug().
		Str("taskid", taskid).
		Str("gpu", gpu.Url).
		Int("buffer_size", buffer.Len()).
		Msg("Created buffer for IPFS upload")

	return []ipfs.IPFSFile{{Name: "output.txt", Path: path, Buffer: buffer}}, nil
}

func (m *MetabaronModel) GetCID(ctx context.Context, gpu *common.GPU, taskid string, input any) ([]byte, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, m.timeoutDuration)
	defer cancel()

	paths, err := utils.ExpRetryWithContext(timeoutCtx, m.logger, func() (any, error) {
		return m.GetFiles(timeoutCtx, gpu, taskid, input)
	}, 3, 1000)
	if err != nil {
		return nil, err
	}

	ipfsCtx, ipfsCancel := context.WithTimeout(ctx, m.ipfsTimeoutDuration)
	defer ipfsCancel()

	cid58, err := utils.ExpRetryWithContext(ipfsCtx, m.logger, func() (any, error) {
		return m.ipfs.PinFilesToIPFS(ipfsCtx, taskid, paths.([]ipfs.IPFSFile))
	}, 3, 1000)
	if err != nil {
		return nil, err
	}

	cidBytes, err := base58.Decode(cid58.(string))
	if err != nil {
		return nil, err
	}

	return cidBytes, nil
}

func (m *MetabaronModel) Validate(gpu *common.GPU, taskid string) error {
	testPrompt := MetabaronPrompt{
		Input: MetabaronInner{
			Prompt: "<|begin_of_text|><|start_header_id|>system<|end_header_id|> You are a helpful assistant.<|eot_id|><|start_header_id|>user<|end_header_id|> what's 4+4?<|eot_id|><|start_header_id|>assistant<|end_header_id|>",
			Seed:   42,
		},
	}

	cid, err := m.GetCID(context.Background(), gpu, "startup-test-taskid", testPrompt)
	if err != nil {
		return err
	}

	cidStr := "0x" + hex.EncodeToString(cid)
	m.logger.Info().Str("model", m.GetID()).Str("cid", cidStr).Msg("model CID after validation")
	return nil
}
