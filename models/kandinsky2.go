package models

import (
	"bytes"
	"encoding/base64"
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
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mr-tron/base58"
	"github.com/rs/zerolog"
)

type Kadinsky2Inner struct {
	Prompt string `json:"prompt"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Seed   uint64 `json:"seed"`
}

type Kandinsky2Prompt struct {
	Input Kadinsky2Inner `json:"input"`
}

type Kandinsky2ModelResponse struct {
	Input  map[string]interface{} `json:"input"`
	Output []string               `json:"output"`
}

type Kandinsky2Model struct {
	Model
	Filters []MiningFilter
	ipfs    ipfs.IPFSClient
	config  *config.AppConfig
	//url     string
	client *http.Client
	logger *zerolog.Logger
}

// Ensure Kandinsky2Model implements the Model interface.
var _ ModelInterface = (*Kandinsky2Model)(nil)

var Kandinsky2ModelTemplate = Model{
	ID:       "",
	Mineable: true,
	Template: map[string]interface{}{
		"meta": map[string]interface{}{
			"title":       "Kandinsky 2",
			"description": "text2img model trained on LAION HighRes and fine-tuned on internal datasets",
			"git":         "https://github.com/kasumi-1/Kandinsky-2/tree/aa5ee2f68a1785b0833d32c27dff097b9b5e8f47",
			"docker":      "r8.im/kasumi-1/kandinsky-2@sha256:373fa540ae197fc89f0679ed835bc4524152956d4f3027580244e10b09d6d3a5",
			"version":     1,
			"input": []map[string]interface{}{
				{
					"variable":    "prompt",
					"type":        "string",
					"required":    true,
					"default":     "",
					"description": "Input prompt",
				},
				{
					"variable":    "width",
					"type":        "int_enum",
					"required":    false,
					"choices":     []int{768, 1024},
					"default":     768,
					"description": "Width of output image.",
				},
				{
					"variable":    "height",
					"type":        "int_enum",
					"required":    false,
					"choices":     []int{768, 1024},
					"default":     768,
					"description": "Height of output image.",
				},
			},
			"output": []map[string]interface{}{
				{
					"filename": "out-1.png",
					"type":     "image",
				},
			},
		},
	},
}

func NewKandinsky2Model(client ipfs.IPFSClient, appConfig *config.AppConfig, logger *zerolog.Logger) *Kandinsky2Model {

	if logger == nil {
		nopLogger := zerolog.Nop()
		logger = &nopLogger
	}

	model, ok := appConfig.BaseConfig.Models["kandinsky2"]
	if !ok {
		return nil
	}

	if model.ID == "" {
		logger.Error().Str("model", "kandinsky2").Msg("kandinsky2 model ID is empty")
		return nil
	}

	http := &http.Client{
		Timeout: time.Second * 30, // TODO: make this a config based setting - set timeout to 30 seconds
	}

	m := &Kandinsky2Model{
		Model:  Kandinsky2ModelTemplate,
		config: appConfig,
		ipfs:   client,
		Filters: []MiningFilter{
			{
				MinFee:  0,
				MinTime: 0,
			},
		},
		client: http,
		logger: logger,
	}
	m.Model.ID = model.ID
	return m
}

func (m *Kandinsky2Model) HydrateInput(preprocessedInput map[string]interface{}, seed uint64) (InputHydrationResult, error) {
	input := make(map[string]interface{})

	// messy but works
	for _, row := range m.Model.Template.(map[string]interface{})["meta"].(map[string]interface{})["input"].([]map[string]interface{}) {
		col, ok := preprocessedInput[row["variable"].(string)]

		if row["required"].(bool) {
			if !ok {
				return nil, fmt.Errorf("input missing required field (%s)", row["variable"])
			}
		}

		if ok {
			switch row["type"].(string) {
			case "string", "string_enum":
				_, ok := col.(string)
				if !ok {
					return nil, fmt.Errorf("input wrong type (%s)", row["variable"])
				}
			case "int", "int_enum", "decimal":
				switch v := col.(type) {
				case int:
					// col is of type int, no further checks needed
				case float64:
					// col is of type float64, check if it's an integer
					if v != float64(int(v)) {
						return nil, fmt.Errorf("input wrong type (%s)", row["variable"])
					}
				default:
					return nil, fmt.Errorf("input wrong type (%s)", row["variable"])
				}
			}

			if row["type"].(string) == "int" || row["type"].(string) == "decimal" {
				if col.(int) < row["min"].(int) || col.(int) > row["max"].(int) {
					return nil, fmt.Errorf("input out of bounds (%s)", row["variable"])
				}
			}

			if row["type"].(string) == "string_enum" || row["type"].(string) == "int_enum" {
				found := false
				switch choices := row["choices"].(type) {
				case []string:
					for _, choice := range choices {
						if choice == col.(string) {
							found = true
							break
						}
					}
				case []int:
					var colInt int
					switch v := col.(type) {
					case int:
						colInt = v
					case float64:
						colInt = int(v)
					default:
						return nil, fmt.Errorf("unexpected type for col (%s)", row["variable"])
					}
					for _, choice := range choices {
						if choice == colInt {
							found = true
							break
						}
					}
				default:
					return nil, fmt.Errorf("unexpected type for choices (%s)", row["variable"])
				}
				if !found {
					return nil, fmt.Errorf("input not in enum (%s)", row["variable"])
				}
			}

			input[row["variable"].(string)] = col
		}

		if !ok {
			input[row["variable"].(string)] = row["default"]
		}
	}

	// This takes our mapped input and converts it to the expected format
	// input json => map[string]interface{} => json => Kadinsky2Prompt
	// This ensure that we have some type safety
	var inner Kadinsky2Inner

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &inner)
	if err != nil {
		return nil, err
	}
	inner.Seed = seed

	result := Kandinsky2Prompt{
		Input: inner,
	}

	return result, nil
}

func (m *Kandinsky2Model) GetID() string {
	return m.Model.ID
}

func (m *Kandinsky2Model) GetFiles(gpu *common.GPU, taskid string, input interface{}) ([]ipfs.IPFSFile, error) {
	// TODO: validate this?
	//url := m.config.ML.Cog[m.Model.ID].URL

	marshaledInput, _ := json.Marshal(input)

	//start := time.Now()
	postResp, err := m.client.Post(gpu.Url, "application/json", bytes.NewBuffer([]byte(marshaledInput)))
	//elapsed := time.Since(start)
	//fmt.Println("GPU Post took:", elapsed)
	if err != nil {
		return nil, err
	}
	defer postResp.Body.Close()

	body, err := io.ReadAll(postResp.Body)
	if err != nil {
		return nil, err
	}

	var resp Kandinsky2ModelResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Output) != 1 {
		return nil, err
	}

	// Remove the "data:image/png;base64," prefix
	resp.Output[0] = strings.TrimPrefix(resp.Output[0], "data:image/png;base64,")

	// Assuming body is a base64 encoded string
	buf, err := base64.StdEncoding.DecodeString(resp.Output[0])
	if err != nil {
		return nil, err
	}

	fileName := fmt.Sprintf("%d.%s.png", gpu.ID, uuid.New().String())

	path := filepath.Join(m.config.CachePath, fileName)
	// err = os.WriteFile(path, buf, 0644)
	// if err != nil {
	// 	return nil, err
	// }

	buffer := bytes.NewBuffer(buf)

	return []ipfs.IPFSFile{{Name: "out-1.png", Path: path, Buffer: buffer}}, nil
}

func (m *Kandinsky2Model) GetCID(gpu *common.GPU, taskid string, input interface{}) ([]byte, error) {
	//start := time.Now()
	paths, err := utils.ExpRetry(func() (interface{}, error) {
		return m.GetFiles(gpu, taskid, input)
	}, 3, 1000)
	//elapsed := time.Since(start)
	//fmt.Println("GPU DIRECT CALL TOOK:", elapsed)
	if err != nil {
		return nil, err //errors.New("cannot get paths")
	}

	//start = time.Now()
	// TODO: calculate cid and pin async
	cid58, err := utils.ExpRetry(func() (interface{}, error) {
		return m.ipfs.PinFilesToIPFS(taskid, paths.([]ipfs.IPFSFile))
	}, 3, 1000)
	//elapsed = time.Since(start)

	//fmt.Println("IPFS DIRECT CALL TOOK:", elapsed)
	if err != nil {
		return nil, errors.New("cannot pin files to retrieve cid")
	}

	cidBytes, err := base58.Decode(cid58.(string))
	if err != nil {
		return nil, err
	}

	//cid := "0x" + hex.EncodeToString(cidBytes)
	return cidBytes, nil
}

func (m *Kandinsky2Model) Validate(gpu *common.GPU, taskid string) error {
	testPrompt := Kandinsky2Prompt{
		Input: Kadinsky2Inner{
			Prompt: "render a cat in the style of kandinsky",
			Height: 768,
			Width:  768,
			Seed:   1337,
		},
	}

	cid, err := m.GetCID(gpu, "startup-test-taskid", testPrompt)
	if err != nil {
		return err
	}

	expected := "0x12200f8c99111abf301ceb8965af7b111c77bcd6e1903c0c713c4b610665dd270be3"
	cidStr := "0x" + hex.EncodeToString(cid)
	if cidStr == expected {
		m.logger.Info().Str("model", m.GetID()).Str("cid", cidStr).Str("expected", expected).Msg("model CID matches expected CID")
	} else {
		m.logger.Error().Str("model", m.GetID()).Str("cid", cidStr).Str("expected", expected).Msg("model CID does not match expected CID")
		return errors.New("model CID does not match expected CID")
	}
	return nil
}
