package client

import (
	"net/http"

	"github.com/puppetlabs/relay/pkg/config"
	"github.com/puppetlabs/relay/pkg/model"
)

const API_VERSION = "v20200131"

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
