package ipfs

import (
	"fmt"
	"gobius/config"

	//blocks "github.com/ipfs/go-block-format"

	"github.com/ipfs/kubo/client/rpc"
	"github.com/ipfs/kubo/core/coreiface/options"
	ma "github.com/multiformats/go-multiaddr"
)

// HttpIPFSClient implements IPFSClient using a remote Kubo node via RPC.
type HttpIPFSClient struct {
	BaseIPFSClient // Embed the base client
}

// NewHttpIPFSClient creates a new client that connects to a Kubo RPC endpoint.
func NewHttpIPFSClient(appConfig config.AppConfig, hashOnly bool) (*HttpIPFSClient, error) {

	multiAddr, err := ma.NewMultiaddr(appConfig.IPFS.HTTPClient.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Kubo RPC multiaddr: %w", err)
	}

	// Connect to the Kubo RPC API
	api, err := rpc.NewApi(multiAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kubo RPC API at %s: %w", multiAddr, err)
	}

	// Prepare specific options for this client instance
	// Start with a copy of the defaults
	clientOptions := append([]options.UnixfsAddOption(nil), defaultIPFSOptions...)
	// Add the hashOnly option if specified
	clientOptions = append(clientOptions, options.Unixfs.HashOnly(hashOnly))

	// Initialize the BaseIPFSClient part
	baseClient, err := NewBaseIPFSClient(appConfig, api) // Pass the connected API
	if err != nil {
		// This shouldn't typically fail with the current NewBaseIPFSClient, but check anyway
		return nil, fmt.Errorf("failed to initialize base IPFS client: %w", err)
	}

	// Override the default options in the base client with our specific ones
	baseClient.ipfsOptions = clientOptions

	// Return the composed client
	return &HttpIPFSClient{
		BaseIPFSClient: *baseClient,
	}, nil
}
