package models

import (
	"bytes"
	"context"
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
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mr-tron/base58"
	"github.com/rs/zerolog"
)

type ESMFoldV1Input struct {
	Sequence string `json:"sequence"`
}

type ESMFoldV1Model struct {
	Model
	timeoutDuration     time.Duration
	ipfsTimeoutDuration time.Duration
	Filters             []MiningFilter
	config              *config.AppConfig
	client              *http.Client
	logger              zerolog.Logger
	ipfs                ipfs.IPFSClient
}

var _ ModelInterface = (*ESMFoldV1Model)(nil)

var ESMFoldV1ModelTemplate = Model{
	ID:       "",
	Mineable: true,
	Template: map[string]any{
		"meta": map[string]any{
			"title":       "ESMFold v1",
			"description": "Protein folding model from Facebook Research based on ESM architecture",
			"git":         "https://github.com/facebookresearch/esm",
			"docker":      "https://hub.docker.com/r/cody8295/esmfold-arbius",
			"version":     1,
		},
		"input": []map[string]any{
			{
				"variable":    "sequence",
				"type":        "string",
				"required":    true,
				"default":     "",
				"description": "Amino acid sequence to predict protein structure from",
			},
		},
		"output": []map[string]any{
			{
				"filename": "result.pdb",
				"type":     "text",
			},
		},
	},
}

func NewESMFoldV1Model(client ipfs.IPFSClient, appConfig *config.AppConfig, logger zerolog.Logger) *ESMFoldV1Model {
	model, ok := appConfig.BaseConfig.Models["esmfold-v1"]
	if !ok {
		return nil
	}

	if model.ID == "" {
		logger.Error().Str("model", "esmfold-v1").Msg("model ID is empty")
		return nil
	}

	http := &http.Client{
		Transport: &http.Transport{MaxIdleConnsPerHost: 10},
	}

	timeout := 600 * time.Second
	ipfsTimeout := 30 * time.Second

	m := &ESMFoldV1Model{
		Model:               ESMFoldV1ModelTemplate,
		timeoutDuration:     timeout,
		ipfsTimeoutDuration: ipfsTimeout,
		config:              appConfig,
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
	m.Model.ID = model.ID
	return m
}

func (m *ESMFoldV1Model) HydrateInput(preprocessedInput map[string]any, seed uint64) (InputHydrationResult, error) {
	input := make(map[string]any)

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

	for _, field := range inputFields {
		varName := field["variable"].(string)
		fieldType := field["type"].(string)
		required, _ := field["required"].(bool)

		value, exists := preprocessedInput[varName]
		if required && !exists {
			return nil, fmt.Errorf("input missing required field (%s)", varName)
		}

		if exists {
			if err := validateType(value, fieldType, varName); err != nil {
				return nil, err
			}
			input[varName] = value
		} else {
			input[varName] = field["default"]
		}
	}

	var inner ESMFoldV1Input
	inner.Sequence, _ = input["sequence"].(string)
	return inner, nil
}

func (m *ESMFoldV1Model) GetID() string {
	return m.Model.ID
}

func (m *ESMFoldV1Model) GetFiles(ctx context.Context, gpu *common.GPU, taskid string, input any) ([]ipfs.IPFSFile, error) {
	if err := ctx.Err(); err != nil {
		m.logger.Warn().Err(err).Str("task", taskid).Msg("Context canceled before GetFiles execution")
		return nil, err
	}

	inner, ok := input.(ESMFoldV1Input)
	if !ok {
		return nil, fmt.Errorf("invalid input type: expected ESMFoldV1Input")
	}

	payload := map[string]string{"sequence": inner.Sequence}
	marshaledInput, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	endpoint := fmt.Sprintf("%s:8080/predict", strings.TrimSuffix(gpu.Url, "/"))
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(marshaledInput))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			m.logger.Error().Err(err).Str("task", taskid).Str("gpu", endpoint).Msg("model inference request timed out")
			return nil, fmt.Errorf("model inference timed out: %w", err)
		}
		return nil, fmt.Errorf("failed to POST to GPU: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusConflict {
			m.logger.Warn().Str("task", taskid).Str("gpu", endpoint).Int("status", resp.StatusCode).Str("body", string(bodyBytes)).Msg("resource busy")
			return nil, ErrResourceBusy
		}
		return nil, fmt.Errorf("server returned non-200 status: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	pdbData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read model response body: %w", err)
	}

	if len(pdbData) == 0 {
		return nil, errors.New("model returned empty PDB data")
	}

	fileName := fmt.Sprintf("%d.%s.pdb", gpu.ID, uuid.New().String())
	path := filepath.Join(m.config.CachePath, fileName)
	buffer := bytes.NewBuffer(pdbData)

	return []ipfs.IPFSFile{{Name: "result.pdb", Path: path, Buffer: buffer}}, nil
}

func (m *ESMFoldV1Model) GetCID(ctx context.Context, gpu *common.GPU, taskid string, input any) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("parent context canceled before GetCID: %w", err)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, m.timeoutDuration)
	defer cancel()

	paths, err := utils.ExpRetryWithContext(timeoutCtx, m.logger, func() (any, error) {
		return m.GetFiles(timeoutCtx, gpu, taskid, input)
	}, 3, 1000)
	if err != nil {
		if errors.Is(err, ErrResourceBusy) {
			m.logger.Warn().Str("task", taskid).Str("gpu", gpu.Url).Msg("GPU remained busy after retries")
		}
		return nil, err
	}

	ipfsCtx, ipfsCancel := context.WithTimeout(ctx, m.ipfsTimeoutDuration)
	defer ipfsCancel()

	cid58, err := utils.ExpRetryWithContext(ipfsCtx, m.logger, func() (any, error) {
		return m.ipfs.PinFilesToIPFS(ipfsCtx, taskid, paths.([]ipfs.IPFSFile))
	}, 3, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to pin files to IPFS after retries: %w", err)
	}
	cidBytes, err := base58.Decode(cid58.(string))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base58 CID string: %w", err)
	}

	return cidBytes, nil
}

func (m *ESMFoldV1Model) Validate(gpu *common.GPU, taskid string) error {
	// testInput := ESMFoldV1Input{
	// 	Sequence: "MKTAYIAKQRQISFVKSHFSRQLEERLGLIEVQAPILSRVGDGTQDNLSGAEKAVQVKVKALPDAQFKKVVK",
	// }

	// cid, err := m.GetCID(context.Background(), gpu, "startup-test-taskid", testInput)
	// if err != nil {
	// 	return fmt.Errorf("validation failed: %w", err)
	// }

	//cidStr := "0x" + // TODO: Implement Validation here

	return nil
}