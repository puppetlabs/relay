package client

import (
	"net/http"

	"github.com/puppetlabs/relay/pkg/config"
)

const API_VERSION = "v20200131"

type Client struct {
	config      *config.Config
	httpClient  *http.Client
	loadedToken *Token
}

func NewClient(config *config.Config) *Client {
	httpClient := &http.Client{}
	var loadedToken *Token = nil

	return &Client{
		config,
		httpClient,
		loadedToken,
	}
}
