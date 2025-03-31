package models

import (
	"encoding/json"
	"fmt"
	"gobius/common"
	"gobius/config"
	"gobius/ipfs"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func Test_Kandinsky2Model_Config(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("caught panic: %v", r)
		}
	}()

	appConfig := config.AppConfig{}
	appConfig.ML.Cog = map[string]config.Cog{
		"0x11111111111111": {
			URL: []string{"http://howmanycaloriesindust.com"},
		},
	}
	ipfsClient, err := ipfs.NewMockIPFSClient(appConfig, true)
	if err != nil {
		t.Fatal(err)
	}

	logger := zerolog.Nop()

	NewKandinsky2Model(ipfsClient, &appConfig, logger)
	t.Errorf("expected panic")
}

// TODO: this test is incomplete
func Test_Kandinsky2Model_WithMockIPFS_GetFiles(t *testing.T) {
	appConfig := config.AppConfig{}
	modelID := "0x98617a8cd4a11db63100ad44bea4e5e296aecfd78b2ef06aee3e364c7307f212"
	appConfig.ML.Cog = map[string]config.Cog{
		modelID: {
			URL: []string{"http://localhost:41108/predictions"},
		},
	}
	ipfsClient, err := ipfs.NewMockIPFSClient(appConfig, true)
	if err != nil {
		t.Fatal(err)
	}

	logger := zerolog.Nop()

	model := NewKandinsky2Model(ipfsClient, &appConfig, logger)

	testPrompt := Kandinsky2Prompt{
		Input: Kadinsky2Inner{
			Prompt: "arbius test cat",
			Seed:   1337,
		},
	}

	//	input := `"input": {"prompt": "arbius test cat", "seed": 1337}`
	gpu := common.NewGPU(0, appConfig.ML.Cog[modelID].URL[0])

	file, err := model.GetFiles(gpu, "startup-test-taskid", testPrompt)
	if err != nil {
		t.Fatal(err)
	}
	if len(file) != 1 {
		t.Errorf("Length of returned array was incorrect, got: %d, want: %d.", len(file), 1)
	}

	if file[0].Name != "out-1.png" {
		t.Errorf("Name of file %s was incorrect, got: %s, want: %s.", file, err, "out-1.png")
	}

}

func Test_Kandinsky2Model_HydrateInput(t *testing.T) {
	testcases := []struct {
		inputJson      string
		expectedOutput string
		errMsg         string
	}{
		// missing width & height
		{`{
			"prompt": "A sea of blue and green, with a hint of red.",
			"seed": 96696969420
		  }`, "", "",
		},
		{`{
		  }`, "", "input missing required field",
		},
		{`{
			"prompt": "A sea of blue and green, with a hint of red.",
			"width": 900,
			"height": 100
		  }`, "", "input not in enum (width)",
		},
	}

	appConfig := config.AppConfig{}
	appConfig.ML.Cog = map[string]config.Cog{
		"0x98617a8cd4a11db63100ad44bea4e5e296aecfd78b2ef06aee3e364c7307f212": {
			URL: []string{"http://hell/predictions"},
		},
	}
	ipfsClient, err := ipfs.NewMockIPFSClient(appConfig, true)
	if err != nil {
		t.Fatal(err)
	}

	logger := zerolog.Nop()

	for _, tc := range testcases {

		var result map[string]interface{}
		err = json.Unmarshal([]byte(tc.inputJson), &result)
		if err != nil {
			t.Fatal(err)
		}
		newModel := NewKandinsky2Model(ipfsClient, &appConfig, logger)

		output, err := newModel.HydrateInput(result, 1337)
		if err != nil {
			if tc.errMsg != "" {
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Fatal(err)
				}
				continue
			} else {
				t.Fatal(err)
			}
		}
		outputStr, err := json.Marshal(output)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(outputStr))
	}
}
