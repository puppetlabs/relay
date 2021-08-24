package client

import (
	"net/http"

	"github.com/puppetlabs/relay-client-go/client/pkg/client/openapi"
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
	if config.ContextConfig != nil {
		context := config.CurrentContext
		if contextConfig, ok := config.ContextConfig[context]; ok {
			if contextConfig.Domains != nil {
				cc.Host = contextConfig.Domains.APIDomain.Host
				cc.Scheme = contextConfig.Domains.APIDomain.Scheme
			}
		}
	}
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
