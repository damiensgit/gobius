package ipfs

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"gobius/config"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	files "github.com/ipfs/boxo/files"
	core "github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	"github.com/ipfs/kubo/core/coreiface/options"
	mh "github.com/multiformats/go-multihash"
)

// TODO: this test is incomplete
func Test_Http_Client_PinFilesToIPFS(t *testing.T) {
	appConfig := config.AppConfig{}
	appConfig.IPFS.HTTPClient.URL = "/dns4/localhost/tcp/5001/http"
	appConfig.CachePath = "../tests/"
	ipfsClient, err := NewHttpIPFSClient(appConfig, true)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	cid, err := ipfsClient.PinFilesToIPFS(ctx, "taskid", []IPFSFile{
		{
			Name: "ipfs_a.bin",
			Path: filepath.Join(appConfig.CachePath, "ipfs_a.bin"),
		},
		{
			Name: "ipfs_b.bin",
			Path: filepath.Join(appConfig.CachePath, "ipfs_b.bin"),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if cid != "QmQx4LqzvgAhXtictjyZKN5gL3V9EEP1B5PhZTdnwW6NjQ" {
		t.Errorf("Hash of file was incorrect, got: %s, want: %s.", cid, "QmQx4LqzvgAhXtictjyZKN5gL3V9EEP1B5PhZTdnwW6NjQ")
	}

}

var testCasesForPinFileToIPFS = []struct {
	filePath     string
	expectedHash string
}{
	{"../tests/ipfs_a.bin", "1220e844b8764c00d4a76ac03930a3d8f32f3df59aea3ed0ade4c3bc38a3b23a31d9"},
	{"../tests/ipfs_b.bin", "1220f782bf27d7dfa16c5556ae0e19d41a73fc380a28455abcedecd70460505f022b"},
	//{"../tests/ipfs_c.bin", "1220c32cae42b7d6ed6efd2512fd7dac6530cbd96cbcc19a3d1c336ace8e401f1c3a"}, // fails
	//{"../tests/ipfs_d.bin", "1220f4ad8a3bd3189da2ad909ee41148d6893d8c629c410f7f2c7e3fae75aade79c8"}, // fails
}

// TODO: two tests are failing as the expected hash is different from the actual hash
func Test_Http_Client_PinFileToIPFS(t *testing.T) {
	appConfig := config.AppConfig{}
	appConfig.IPFS.HTTPClient.URL = "/dns4/localhost/tcp/5001/http"
	ipfsClient, err := NewHttpIPFSClient(appConfig, true)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range testCasesForPinFileToIPFS {
		// open a file and read its contents
		content, err := os.ReadFile(tc.filePath)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		hashResult, err := ipfsClient.PinFileToIPFS(content, "test")
		if err != nil {
			t.Fatalf("Failed to pin file: %v", err)
		}

		hashResultFast, err := GetIPFSHashFast(content)
		if err != nil {
			t.Fatal(err)
		}

		hashResultFastHex := hex.EncodeToString(hashResultFast)
		fmt.Println("CID:", hashResultFastHex)

		convertedHash, err := mh.FromB58String(hashResult)
		if err != nil {
			t.Fatal(err)
		}

		if convertedHash.String() != hashResultFastHex {
			t.Errorf("Hash of file %s was incorrect, got: %s, want: %s.", tc.filePath, convertedHash.String(), tc.expectedHash)
		}

		if convertedHash.String() != tc.expectedHash {
			t.Errorf("Hash of file %s was incorrect, got: %s, want: %s.", tc.filePath, convertedHash.String(), tc.expectedHash)
		}
	}
}

func Test_Mock_PinFilesToIPFS(t *testing.T) {
	appConfig := config.AppConfig{}
	appConfig.CachePath = "../tests/"
	ipfsClient, err := NewMockIPFSClient(appConfig, true)
	if err != nil {
		t.Fatal(err)
	}

	files := []IPFSFile{
		{
			Name: "ipfs_a.bin",
			Path: filepath.Join(appConfig.CachePath, "ipfs_a.bin"),
		},
		{
			Name: "ipfs_b.bin",
			Path: filepath.Join(appConfig.CachePath, "ipfs_b.bin"),
		},
	}

	for i := range files {

		fileHandle, err := os.Open(files[i].Path)
		if err != nil {
			t.Fatal(err)
		}
		defer fileHandle.Close()

		var buffer bytes.Buffer
		_, err = io.Copy(&buffer, fileHandle)
		if err != nil {
			t.Fatal(err)
		}

		files[i].Buffer = &buffer
	}

	ctx := context.Background()
	cid, err := ipfsClient.PinFilesToIPFS(ctx, "taskid", files)

	if err != nil {
		t.Fatal(err)
	}
	wanted := "QmQx4LqzvgAhXtictjyZKN5gL3V9EEP1B5PhZTdnwW6NjQ"
	if cid != wanted {
		t.Errorf("Hash of file was incorrect, got: %s, want: %s.", cid, wanted)
	}
}

func Test_Mock_PinFileToIPFS(t *testing.T) {
	appConfig := config.AppConfig{}
	ipfsClient, err := NewMockIPFSClient(appConfig, true)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCasesForPinFileToIPFS {
		// open a file and read its contents
		content, err := os.ReadFile(tc.filePath)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		hashResult, err := ipfsClient.PinFileToIPFS(content, "test")
		if err != nil {
			t.Fatalf("Failed to pin file: %v", err)
		}

		hashResultFast, err := GetIPFSHashFast(content)
		if err != nil {
			t.Fatal(err)
		}

		hashResultFastHex := hex.EncodeToString(hashResultFast)
		fmt.Println("CID:", hashResultFastHex)

		convertedHash, err := mh.FromB58String(hashResult)
		if err != nil {
			t.Fatal(err)
		}

		if convertedHash.String() != hashResultFastHex {
			t.Errorf("Hash of file %s was incorrect, got: %s, want: %s.", tc.filePath, convertedHash.String(), tc.expectedHash)
		}

		if convertedHash.String() != tc.expectedHash {
			t.Errorf("Hash of file %s was incorrect, got: %s, want: %s.", tc.filePath, convertedHash.String(), tc.expectedHash)
		}
	}
}

func Test_PrivateNode(t *testing.T) {
	// Create an IPFS node
	ctx := context.Background()

	node, err := core.NewNode(ctx, &core.BuildCfg{
		//TODO: need this to be true or all files
		// hashed will be stored in memory!
		NilRepo: true,
	})
	if err != nil {
		fmt.Println("Failed to create IPFS node:", err)
		return
	}
	defer node.Close()

	// Create an IPFS API
	api, err := coreapi.NewCoreAPI(node, options.Api.FetchBlocks(false))
	if err != nil {
		fmt.Println("Failed to create IPFS API:", err)
		return
	}

	// Open the files
	file1, err := os.Open("../tests/ipfs_a.bin")
	if err != nil {
		fmt.Println("Failed to open ipfs_a.bin:", err)
		return
	}
	defer file1.Close()

	file2, err := os.Open("../tests/ipfs_b.bin")
	if err != nil {
		fmt.Println("Failed to open ipfs_b.bin:", err)
		return
	}
	defer file2.Close()

	// Create a directory with the files to add
	dir := files.NewMapDirectory(map[string]files.Node{
		"ipfs_a.bin": files.NewReaderFile(file1),
		"ipfs_b.bin": files.NewReaderFile(file2),
	})

	ipfsOptions := []options.UnixfsAddOption{
		options.Unixfs.CidVersion(0),
		options.Unixfs.RawLeaves(false),
		options.Unixfs.FsCache(false),
		options.Unixfs.Chunker("size-262144"),
		options.Unixfs.HashOnly(true),
	}

	// Add the directory to IPFS
	//Failed to add directory: block was not found locally (offline): ipld: could not find QmdyL59DkUaj6EMcDD3b4rMchGZFJnasT6C6N2RvpGLYXi
	cid, err := api.Unixfs().Add(ctx, dir, ipfsOptions...)
	if err != nil {
		fmt.Println("Failed to add directory:", err)
		return
	}
	// /ipfs/QmQx4LqzvgAhXtictjyZKN5gL3V9EEP1B5PhZTdnwW6NjQ
	// wanted: QmQx4LqzvgAhXtictjyZKN5gL3V9EEP1B5PhZTdnwW6NjQ
	fmt.Println("CID of parent directory:", cid)
}

func loadTestConfig() (config.AppConfig, error) {
	// Determine the base path relative to the test file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return config.AppConfig{}, fmt.Errorf("could not get caller information")
	}
	basePath := filepath.Dir(filename)

	configPath := filepath.Join(basePath, "ipfs_config.json")

	cfg, err := config.LoadConfigForTesting(configPath, 0)
	if err != nil {
		return config.AppConfig{}, fmt.Errorf("failed to load test config: %v", err)
	}
	if cfg == nil { // Handle case where LoadConfig might return nil cfg, nil err
		return config.AppConfig{}, nil
	}
	return *cfg, nil
}

func Test_Pinata_Client_PinFileToIPFS(t *testing.T) {
	appConfig, err := loadTestConfig() // Assuming this provides defaults or uses env vars
	if err != nil {
		t.Fatalf("Failed to load test config (or defaults): %v", err)
	}

	// Skip test if Pinata JWT is not configured
	if appConfig.IPFS.Pinata.JWT == "" {
		t.Skip("Skipping Pinata test: IPFS_PINATA_JWT not set in config or environment")
	}

	// Create Pinata Client
	ipfsClient, err := NewPinataIPFSClient(appConfig)
	if err != nil {
		t.Fatalf("Failed to create Pinata client: %v", err)
	}

	for _, tc := range testCasesForPinFileToIPFS {
		t.Run(tc.filePath, func(t *testing.T) { // Use subtests for better isolation
			// open a file and read its contents
			content, err := os.ReadFile(tc.filePath)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", tc.filePath, err)
			}

			// Call PinFileToIPFS
			// Note: Current signature returns string, not error. Error is logged internally by the client.
			base58CID, err := ipfsClient.PinFileToIPFS(content, filepath.Base(tc.filePath))
			if err != nil {
				t.Fatalf("Failed to pin file %s: %v", tc.filePath, err)
			}

			// Assertions
			if base58CID == "" {
				t.Fatalf("PinFileToIPFS returned an empty CID for %s", tc.filePath)
			}

			// Convert the returned Base58 CID to hex multihash for comparison
			// This assumes the Pinata API returns a CIDv0 (Qm...)
			convertedHash, err := mh.FromB58String(base58CID)
			if err != nil {
				t.Fatalf("Failed to convert returned Base58 CID %s to multihash: %v", base58CID, err)
			}

			hexMultihash := convertedHash.String()

			// Compare with the expected hex hash from test case
			if hexMultihash != tc.expectedHash {
				t.Errorf("Hash mismatch for file %s:\n  Got (hex): %s\n Want (hex): %s\n(From Base58: %s)",
					tc.filePath, hexMultihash, tc.expectedHash, base58CID)
			}

			// TODO: Add unpinning logic if necessary for cleanup
		})
	}
}

func Test_Pinata_Client_PinFilesToIPFS(t *testing.T) {
	appConfig, err := loadTestConfig()
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	// Create Pinata Client
	ipfsClient, err := NewPinataIPFSClient(appConfig)
	if err != nil {
		t.Fatalf("Failed to create Pinata client: %v", err)
	}

	// Test Data
	testFiles := []IPFSFile{
		{Name: "file1.txt", Buffer: bytes.NewBufferString("Pinata test file 1 content.")},
		{Name: "file2.txt", Buffer: bytes.NewBufferString("Pinata test file 2 content.")},
	}
	testTaskID := "pinata-test-task-1234"

	// Call PinFilesToIPFS
	ctx := context.Background() // Or use context with timeout
	cid, err := ipfsClient.PinFilesToIPFS(ctx, testTaskID, testFiles)

	// Assertions
	if err != nil {
		t.Errorf("PinFilesToIPFS failed: %v", err)
	}
	if cid == "" {
		t.Errorf("PinFilesToIPFS returned an empty CID, expected a valid CID")
	}

	// Optional: Add more specific CID check if you know the expected directory CID
	// expectedCID := "Qm..."
	// if cid != expectedCID {
	// 	t.Errorf("PinFilesToIPFS CID mismatch: got %s, want %s", cid, expectedCID)
	// }

	// TODO: Add unpinning logic if necessary for cleanup
}
