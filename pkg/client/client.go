package client

import (
	"net/http"

	"github.com/puppetlabs/relay/pkg/client/openapi"
	"github.com/puppetlabs/relay/pkg/config"
	"github.com/puppetlabs/relay/pkg/model"
)

const APIVersion = "v20200615"

type Client struct {
	Api *openapi.APIClient

	config      *config.Config
	httpClient  *http.Client
	loadedToken *model.Token
}

func NewClient(config *config.Config) *Client {
	cc := openapi.NewConfiguration()
	cc.Host = config.ContextConfig.APIDomain.Host
	cc.Scheme = config.ContextConfig.APIDomain.Scheme
	cc.Debug = false

	api := openapi.NewAPIClient(cc)

	httpClient := &http.Client{}
	var loadedToken *model.Token = nil

	return &Client{
		api,
		config,
		httpClient,
		loadedToken,
	}
}
