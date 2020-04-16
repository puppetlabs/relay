package client

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

type Token string

func (t *Token) Bearer() string {
	return fmt.Sprintf("Bearer %s", t)
}

func (t *Token) String() string {
	return string(*t)
}

// getToken reads token from client cache or from path specified on config
func (c *Client) getToken() (*Token, error) {
	if c.loadedToken == nil {
		f, err := os.Open(c.config.TokenPath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, nil
			}

			return nil, err
		}

		defer f.Close()

		buf := &bytes.Buffer{}
		if _, err := buf.ReadFrom(f); err != nil {
			return nil, err
		}

		token := Token(buf.String())

		c.loadedToken = &token

		return &token, nil
	}

	return c.loadedToken, nil
}

// storeToken Saves token to the token storage location specified by config,
// creating directories as needed
func (c *Client) storeToken(token *Token) error {
	if err := os.MkdirAll(filepath.Dir(c.config.TokenPath), 0750); err != nil {
		return err
	}

	f, err := os.OpenFile(c.config.TokenPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0750)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.Write([]byte(token.String())); err != nil {
		return err
	}

	return nil
}

// clearToken removes token from local storage and from loadedToken cache on client object.
// It does not error if the token does not exist.
func (c *Client) clearToken() error {
	if err := os.Remove(c.config.TokenPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	c.loadedToken = nil

	return nil
}
