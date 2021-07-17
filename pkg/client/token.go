package client

import (
	"github.com/puppetlabs/relay/pkg/config"
	"github.com/puppetlabs/relay/pkg/model"
)

// getToken reads token from client cache or from path specified on config
func (c *Client) getToken() (*model.Token, error) {
	if c.loadedToken == nil {
		if c.config.ContextConfig != nil {
			context := c.config.CurrentContext

			if contextConfig, ok := c.config.ContextConfig[context]; ok {
				if contextConfig.Auth != nil && contextConfig.Auth.Tokens != nil {
					for _, tokenType := range config.AuthTokenTypes() {
						if value, ok := contextConfig.Auth.Tokens[tokenType]; ok && value != "" {
							token := model.Token(value)
							c.loadedToken = &token
							return c.loadedToken, nil
						}
					}
				}
			}
		}
	}

	return c.loadedToken, nil
}

// clearToken clears the token from the loadedToken cache on client object.
func (c *Client) clearToken() error {
	c.loadedToken = nil

	return nil
}
