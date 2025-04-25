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
	PinFileToIPFS(data []byte, filename string) string
}

type BaseIPFSClient struct {
	config      config.AppConfig
	api         iface.CoreAPI
	ipfsOptions []options.UnixfsAddOption
	pinata      *PinataClient
}

var defaultIPFSOptions = []options.UnixfsAddOption{
	options.Unixfs.CidVersion(0),
	options.Unixfs.RawLeaves(false),
	options.Unixfs.Pin(true),
}

type IPFSFile struct {
	Name   string // name of file on IPFS
	Path   string // local path to file to add to IPFS
	Buffer *bytes.Buffer
}

func NewBaseIPFSClient(cfg config.AppConfig) (*BaseIPFSClient, error) {
	client := &BaseIPFSClient{
		config:      cfg,
		ipfsOptions: defaultIPFSOptions,
	}

	// Initialize Pinata client if enabled
	if cfg.IPFS.Pinata.Enabled {
		client.pinata = NewPinataClient(
			cfg.IPFS.Pinata.APIKey,
			cfg.IPFS.Pinata.APISecret,
			cfg.IPFS.Pinata.JWT,
		)
	}

	return client, nil
}

// Note: filename is not used in this function until pinata support is added
func (ic *BaseIPFSClient) PinFileToIPFS(data []byte, filename string) string {
	// Try Pinata first if enabled in config
	if ic.config.IPFS.Pinata.Enabled && ic.pinata != nil {
		cid := ic.pinata.PinFileToIPFS(data, filename)
		if cid != "" {
			return cid
		}
	}

	// Fall back to local IPFS
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	file := files.NewBytesFile(data)
	test, err := ic.api.Unixfs().Add(ctx, file, ic.ipfsOptions...)
	if err != nil {
		return ""
	}
	return test.RootCid().String()
}

// PinFilesToIPFS adds the files to IPFS
// note: taskid is not used in this function until pinata support is added
func (ic *BaseIPFSClient) PinFilesToIPFS(ctx context.Context, taskid string, filesToAdd []IPFSFile) (string, error) {
	// Try Pinata first if enabled in config
	if ic.config.IPFS.Pinata.Enabled && ic.pinata != nil {
		cid, err := ic.pinata.PinFilesToIPFS(ctx, taskid, filesToAdd)
		if err == nil && cid != "" {
			return cid, nil
		}
	}

	// Fall back to local IPFS
	if err := ctx.Err(); err != nil {
		return "", err
	}

	mapOfFiles := map[string]files.Node{}
	for _, file := range filesToAdd {
		mapOfFiles[file.Name] = files.NewReaderFile(file.Buffer)
	}
	mapDirectory := files.NewMapDirectory(mapOfFiles)

	test, err := ic.api.Unixfs().Add(ctx, mapDirectory, ic.ipfsOptions...)
	if err != nil {
		return "", err
	}
	return test.RootCid().String(), nil
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
