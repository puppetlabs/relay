package client

import (
	"bytes"
	"fmt"
	"os"
)

type Token string

func (t *Token) Bearer() string {
	return fmt.Sprintf("Bearer %s", t)
}

func (t *Token) String() string {
	return string(*t)
}

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

func (c *Client) storeToken(token *Token) error {
	if err := os.MkdirAll(c.config.CacheDir, 0750); err != nil {
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

func (c *Client) clearToken() error {
	if err := os.Remove(c.config.TokenPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	c.loadedToken = nil

	return nil
}
