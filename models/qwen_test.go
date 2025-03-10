package models

import (
	"encoding/json"
	"fmt"
	"gobius/common"
	"gobius/config"
	"gobius/ipfs"
	"strings"
	"testing"
)

func Test_QwenTestModel_Config(t *testing.T) {
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

	NewQwenTestModel(ipfsClient, &appConfig, nil)
	t.Errorf("expected panic")
}

// TODO: this test is incomplete
// {"input":{"prompt":"Trump is a....","n":1,"max_length":50,"temperature":0.75,"top_p":1.0,"repetition_penalty":1,"seed":500},"output":["Trump is a....?\nOptions:\n- president\n- prime minister\n- 45th president\n- monarch\n- king\nLet's think through this step by step: King and monarch are titles of monarchy. But the question refers to a"],"id":null,"version":null,"created_at":null,"started_at":"2025-03-07T09:40:15.400981+00:00","completed_at":"2025-03-07T09:40:16.847613+00:00","logs":"","error":null,"status":"succeeded","metrics":{"predict_time":1.446632}}
func Test_QwenTestModel_WithMockIPFS_GetFiles(t *testing.T) {
	appConfig := config.AppConfig{}
	modelID := "0x98617a8cd4a11db63100ad44bea4e5e296aecfd78b2ef06aee3e364c7307f212"
	appConfig.ML.Cog = map[string]config.Cog{
		modelID: {
			URL: []string{"http://150.136.215.9:8000/predictions"},
		},
	}
	ipfsClient, err := ipfs.NewMockIPFSClient(appConfig, true)
	if err != nil {
		t.Fatal(err)
	}

	model := NewQwenTestModel(ipfsClient, &appConfig, nil)

	testPrompt := QwenPrompt{
		Input: QwenInner{
			Prompt: "the answer to the universe is",
			Seed:   500,
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

	if file[0].Name != "out-1.txt" {
		t.Errorf("Name of file %s was incorrect, got: %s, want: %s.", file[0].Name, err, "out-1.txt")
	}

	cid, err := model.GetCID(gpu, "startup-test-taskid", testPrompt)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(cid)

}

func Test_QwenTestModel_HydrateInput(t *testing.T) {
	testcases := []struct {
		inputJson      string
		expectedOutput string
		errMsg         string
	}{
		{`{
			"prompt": "once upon a time"
		  }`, "", "input missing required field (seed)",
		},
		{`{
		  }`, "", "input missing required field",
		},
		{`{
			"prompt": "once upon a time, in a land far far away",
			"teeth": 42,
			"seed": 100
		  }`, "", "input not in enum (teeth)",
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

	for _, tc := range testcases {

		var result map[string]interface{}
		err = json.Unmarshal([]byte(tc.inputJson), &result)
		if err != nil {
			t.Fatal(err)
		}
		newModel := NewQwenTestModel(ipfsClient, &appConfig, nil)

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
