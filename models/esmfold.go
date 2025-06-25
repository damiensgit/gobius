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
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"os"
	"os/exec"
	"crypto/tls"
	
	"github.com/google/uuid"
	"github.com/mr-tron/base58"
	"github.com/rs/zerolog"
)

type ESMFoldV1Inner struct {
	Prompt string `json:"prompt"`
}

type ESMFoldV1Prompt struct {
	Input ESMFoldV1Inner `json:"input"`
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
	inputFields         []map[string]any
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
			"input": []map[string]any{
				{
					"variable":    "prompt",
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

	httpClient := &http.Client{
        Transport: &http.Transport{
            MaxIdleConnsPerHost: 10,
            ForceAttemptHTTP2:   true,
            TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
        },
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            if len(via) >= 10 {
                return fmt.Errorf("too many redirects")
            }
            logger.Debug().Str("redirect_url", req.URL.String()).Msg("Following redirect")
            return nil
        },
        Timeout: 600 * time.Second,
    }

	timeout := 600 * time.Second
	ipfsTimeout := 30 * time.Second

	cogConfig, ok := appConfig.ML.Cog[model.ID]
	if ok {
		if cogConfig.HttpTimeout != "" {
			parsedTimeout, err := time.ParseDuration(cogConfig.HttpTimeout)
			if err != nil {
				logger.Warn().Err(err).Str("model", model.ID).Str("config_timeout", cogConfig.HttpTimeout).Msg("failed to parse model timeout from cog config, using default 120s")
			} else {
				timeout = parsedTimeout
			}
		} 
		
		if cogConfig.IpfsTimeout != "" {
			parsedIpfsTimeout, err := time.ParseDuration(cogConfig.IpfsTimeout)
			if err != nil {
				logger.Warn().Err(err).Str("model", model.ID).Str("config_ipfs_timeout", cogConfig.IpfsTimeout).Msg("failed to parse IPFS timeout from cog config, using default 30s")
			} else {
				ipfsTimeout = parsedIpfsTimeout
			}
		} 
	}

	templateMeta, ok := ESMFoldV1ModelTemplate.Template.(map[string]any)
	if !ok {
		logger.Error().Str("model", model.ID).Msg("invalid template format")
		return nil
	}

	meta, ok := templateMeta["meta"].(map[string]any)
	if !ok {
		logger.Error().Str("model", model.ID).Msg("invalid meta format in template")
		return nil
	}

	inputFields, ok := meta["input"].([]map[string]any)
	if !ok {
		logger.Error().Str("model", model.ID).Msg("invalid input format in template meta")
		return nil
	}

	m := &ESMFoldV1Model{
		Model: ESMFoldV1ModelTemplate,

		timeoutDuration:     timeout,
		ipfsTimeoutDuration: ipfsTimeout,
		config:              appConfig,
		Filters: []MiningFilter{
			{
				MinFee:  0,
				MinTime: 0,
			},
		},
		ipfs:        client,
		client:      httpClient,
		logger:      logger,
		inputFields: inputFields,
	}
	m.Model.ID = model.ID
	return m
}

func (m *ESMFoldV1Model) HydrateInput(preprocessedInput map[string]any, seed uint64) (InputHydrationResult, error) {
	input := make(map[string]any)

	for _, field := range m.inputFields {
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

	var inner ESMFoldV1Prompt
	inner.Input.Prompt, _ = input["prompt"].(string)
	return inner, nil
}

func (m *ESMFoldV1Model) GetID() string {
	return m.Model.ID
}

func (m *ESMFoldV1Model) GetFiles(ctx context.Context, gpu *common.GPU, taskid string, input interface{}) ([]ipfs.IPFSFile, error) {
    if err := ctx.Err(); err != nil {
        m.logger.Warn().Err(err).Str("task", taskid).Msg("Context canceled before GetFiles execution")
        return nil, err
    }

    inner, ok := input.(ESMFoldV1Prompt)
    if !ok {
        return nil, fmt.Errorf("invalid input type: expected ESMFoldV1Prompt")
    }

    payload := map[string]string{"sequence": inner.Input.Prompt}
    marshaledInput, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal input: %w", err)
    }

    m.logger.Debug().Str("marshaledInput", string(marshaledInput)).Msg("Prepared input payload")
    m.logger.Debug().Str("gpu_url", gpu.Url).Msg("GPU URL")

    endpoint := fmt.Sprintf("%s/predict", strings.TrimSuffix(gpu.Url, "/"))
    m.logger.Debug().Str("endpoint", endpoint).Msg("Constructed endpoint URL")

    tmpFile, err := os.CreateTemp("", "curl_output_*.pdb")
    if err != nil {
        m.logger.Error().Err(err).Str("task", taskid).Msg("Failed to create temporary file")
        return nil, fmt.Errorf("failed to create temporary file: %w", err)
    }
    tmpFileName := tmpFile.Name()
    tmpFile.Close()

    curlArgs := []string{
		"-v",
        "-k",
        "-X", "POST",
        "-H", "Content-Type: application/json",
        "-H", "Accept: chemical/x-pdb",
        "-d", string(marshaledInput),
        "-o", tmpFileName,
        endpoint,
    }

    m.logger.Debug().Str("curl_command", fmt.Sprintf("curl %s", strings.Join(curlArgs, " "))).Msg("Executing curl")

    cmd := exec.CommandContext(ctx, "curl", curlArgs...)
    var stderr bytes.Buffer
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        m.logger.Error().
            Err(err).
            Str("task", taskid).
            Str("gpu", endpoint).
            Str("stderr", stderr.String()).
            Msg("Failed to execute curl")
        os.Remove(tmpFileName)
        return nil, fmt.Errorf("failed to execute curl: %w, stderr: %s", err, stderr.String())
    }

    pdbData, err := os.ReadFile(tmpFileName)
    if err != nil {
        m.logger.Error().Err(err).Str("task", taskid).Msg("Failed to read curl output file")
        os.Remove(tmpFileName)
        return nil, fmt.Errorf("failed to read curl output file: %w", err)
    }
    os.Remove(tmpFileName)

    if len(pdbData) == 0 {
        m.logger.Error().Str("task", taskid).Msg("curl returned empty PDB data")
        return nil, errors.New("curl returned empty PDB data")
    }

    m.logger.Debug().Int("data_length", len(pdbData)).Msg("Received curl response")

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
