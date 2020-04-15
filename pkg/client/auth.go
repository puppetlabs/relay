package client

import (
	"github.com/puppetlabs/relay/pkg/errors"
)

type CreateTokenParameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateTokenResponse struct {
	Token *Token `json:"token"`
}

func (c *Client) CreateToken(email string, password string) errors.Error {

	params := CreateTokenParameters{
		Email:    email,
		Password: password,
	}

	response := &CreateTokenResponse{}

	if err := c.post("/auth/sessions", nil, params, response); err != nil {
		return err
	}

	if err := c.storeToken(response.Token); err != nil {
		return errors.NewClientInternalError().WithCause(err).Bug()
	}

	return nil
}

func (c *Client) InvalidateToken() errors.Error {
	type deleteResponse struct {
		Success bool `json:"success"`
	}

	dr := &deleteResponse{}

	if err := c.delete("/auth/sessions", nil, dr); err != nil {
		return err
	}

	if err := c.clearToken(); err != nil {
		return errors.NewClientInternalError().WithCause(err).Bug()
	}

	return nil
}
