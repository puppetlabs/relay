package client

import (
	"net/http"

	"github.com/puppetlabs/relay/pkg/config"
	"github.com/puppetlabs/relay/pkg/model"
)

const APIVersion = "v20200615"

type Client struct {
	config      *config.Config
	httpClient  *http.Client
	loadedToken *model.Token
}

func NewClient(config *config.Config) *Client {
	httpClient := &http.Client{}
	var loadedToken *model.Token = nil

	return &Client{
		config,
		httpClient,
		loadedToken,
	}
}
