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
}

var defaultIPFSOptions = []options.UnixfsAddOption{
	options.Unixfs.CidVersion(0),
	options.Unixfs.RawLeaves(false),
}

type IPFSFile struct {
	Name   string // name of file on IPFS
	Path   string // local path to file to add to IPFS
	Buffer *bytes.Buffer
}

// Note: filename is not used in this function until pinata support is added
func (ic *BaseIPFSClient) PinFileToIPFS(data []byte, filename string) string {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	file := files.NewBytesFile(data)
	test, err := ic.api.Unixfs().Add(ctx, file, ic.ipfsOptions...)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	fmt.Println(test.RootCid().String())
	return test.RootCid().String()
}

// PinFilesToIPFS adds the files to IPFS
// note: taskid is not used in this function until pinata support is added
func (ic *BaseIPFSClient) PinFilesToIPFS(ctx context.Context, taskid string, filesToAdd []IPFSFile) (string, error) {

	// Check if context is already canceled before doing anything
	if err := ctx.Err(); err != nil {
		// Consider logging here if appropriate
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
