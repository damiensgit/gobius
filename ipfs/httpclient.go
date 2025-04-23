package ipfs

import (
	"gobius/config"

	//blocks "github.com/ipfs/go-block-format"

	"github.com/ipfs/kubo/client/rpc"
	"github.com/ipfs/kubo/core/coreiface/options"
	ma "github.com/multiformats/go-multiaddr"
)

type HttpIPFSClient struct {
	BaseIPFSClient
}

func NewHttpIPFSClient(appConfig config.AppConfig, hashOnly bool) (*HttpIPFSClient, error) {

	newOptions := append([]options.UnixfsAddOption(nil), defaultIPFSOptions...)

	newOptions = append(newOptions, options.Unixfs.HashOnly(hashOnly))

	ma, err := ma.NewMultiaddr(appConfig.IPFS.HTTPClient.URL)
	if err != nil {
		return nil, err
	}

	//api, err := rpc.NewLocalApi()
	api, err := rpc.NewApi(ma)
	if err != nil {
		return nil, err
	}

	return &HttpIPFSClient{
		BaseIPFSClient: BaseIPFSClient{
			config:      appConfig,
			api:         api,
			ipfsOptions: newOptions,
		},
	}, nil
}
