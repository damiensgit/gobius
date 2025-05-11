package ipfs

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"gobius/config"

	"github.com/ipfs/boxo/files"
	iface "github.com/ipfs/kubo/core/coreiface"
	"github.com/ipfs/kubo/core/coreiface/options"
)

// Interface for IPFS things
type IPFSClient interface {
	PinFilesToIPFS(ctx context.Context, taskid string, filesToAdd []IPFSFile) (string, error)
	PinFileToIPFS(data []byte, filename string) (string, error)
}

// BaseIPFSClient provides common fields but should not be used directly.
// Specific implementations like HttpIPFSClient should be used.
type BaseIPFSClient struct {
	config      config.AppConfig
	api         iface.CoreAPI             // Interface to local Kubo node
	ipfsOptions []options.UnixfsAddOption // Options for local Kubo node adds
}

var defaultIPFSOptions = []options.UnixfsAddOption{
	options.Unixfs.CidVersion(0),
	options.Unixfs.RawLeaves(false),
	options.Unixfs.Pin(true),
}

type IPFSFile struct {
	Name   string // name of file on IPFS
	Path   string
	Buffer *bytes.Buffer
}

// NewBaseIPFSClient initializes common fields. SHOULD NOT BE USED DIRECTLY.
// Use NewHttpIPFSClient instead.
func NewBaseIPFSClient(cfg config.AppConfig, api iface.CoreAPI) (*BaseIPFSClient, error) {
	if api == nil {
		return nil, fmt.Errorf("api is nil, cannot create BaseIPFSClient")
	}

	client := &BaseIPFSClient{
		config:      cfg,
		api:         api,                // API must be provided by concrete implementation
		ipfsOptions: defaultIPFSOptions, // Use default options
	}
	return client, nil
}

// PinFileToIPFS pins a single file using the local Kubo node.
// Note: This implementation uses the BaseIPFSClient's api field.
func (ic *BaseIPFSClient) PinFileToIPFS(data []byte, filename string) (string, error) {
	ctx, cancel := context.WithCancel(context.Background()) // Use background context for now
	defer cancel()

	file := files.NewBytesFile(data)
	node, err := ic.api.Unixfs().Add(ctx, file, ic.ipfsOptions...)
	if err != nil {
		return "", fmt.Errorf("error adding file locally: %w", err)
	}

	localCID := node.RootCid().String()
	return localCID, nil
}

// PinFilesToIPFS pins multiple files as a directory using the local Kubo node.
// Note: This implementation uses the BaseIPFSClient's api field.
func (ic *BaseIPFSClient) PinFilesToIPFS(ctx context.Context, taskid string, filesToAdd []IPFSFile) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("context canceled before local pinning: %w", err)
	}

	mapOfFiles := map[string]files.Node{}
	for _, file := range filesToAdd {
		if file.Buffer == nil {
			continue
		}
		mapOfFiles[file.Name] = files.NewReaderFile(file.Buffer)
	}
	// Check if map is empty after potentially skipping files
	if len(mapOfFiles) == 0 {
		return "", fmt.Errorf("no valid files with non-nil buffers provided")
	}
	mapDirectory := files.NewMapDirectory(mapOfFiles)

	node, err := ic.api.Unixfs().Add(ctx, mapDirectory, ic.ipfsOptions...)
	if err != nil {
		return "", fmt.Errorf("error adding directory locally: %w", err)
	}

	localCID := node.RootCid().String()
	return localCID, nil
}

func encodeVarint(n uint64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	len := binary.PutUvarint(buf, n)
	return buf[:len]
}

// Helper function that takes a byte array and returns the IPFS hash.
// This matches the deployed Arbius Engine v2 contract function of same name on the Nova mainnet.
func GetIPFSHashFast(content []byte) ([]byte, error) {
	if len(content) > 65536 {
		return nil, fmt.Errorf("max content size is 65536 bytes")
	}

	contentLengthVarint := encodeVarint(uint64(len(content)))

	// Concatenate the bytes
	meat := append([]byte{0x08, 0x02, 0x12}, contentLengthVarint...)
	meat = append(meat, content...)
	meat = append(meat, 0x18)
	meat = append(meat, contentLengthVarint...)

	// Calculate the SHA-256 hash
	hash := sha256.Sum256(append([]byte{0x0a}, append(encodeVarint(uint64(len(meat))), meat...)...))

	return append([]byte{0x12, 0x20}, hash[:]...), nil
}
