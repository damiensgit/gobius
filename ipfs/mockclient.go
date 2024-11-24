package ipfs

import (
	"context"
	"gobius/config"

	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	"github.com/ipfs/kubo/core/coreiface/options"
)

type MockIPFSClient struct {
	BaseIPFSClient
}

func NewMockIPFSClient(appConfig config.AppConfig, hashOnly bool) (*MockIPFSClient, error) {

	newOptions := append([]options.UnixfsAddOption(nil), defaultIPFSOptions...)

	newOptions = append(newOptions, options.Unixfs.HashOnly(hashOnly))

	ctx := context.Background()

	node, err := core.NewNode(ctx, &core.BuildCfg{
		NilRepo: true,
	})

	if err != nil {
		return nil, err
	}
	defer node.Close()

	// Create an IPFS API
	api, err := coreapi.NewCoreAPI(node, options.Api.FetchBlocks(false))
	if err != nil {
		return nil, err
	}

	return &MockIPFSClient{
		BaseIPFSClient: BaseIPFSClient{
			config:      appConfig,
			api:         api,
			ipfsOptions: newOptions,
		},
	}, nil
}
