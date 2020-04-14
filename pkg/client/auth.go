package client

import "github.com/puppetlabs/relay/pkg/errors"

type CreateTokenParameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateTokenResponse struct {
	Token *Token `json:"token"`
}

func (c *Client) CreateToken(email string, password string) (*Token, errors.Error) {

	params := CreateTokenParameters{
		Email:    email,
		Password: password,
	}

	response := CreateTokenResponse{}

	err := c.post("/auth/sessions", nil, params, response)

	if err != nil {
		return nil, err
	}

	return response.Token, nil
}
