package client

import (
	"net/http"
	"time"

	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/puppetlabs/relay/pkg/model"
)

type createTokenResponse struct {
	Token                   *model.Token `json:"token"`
	UserCode                string       `json:"user_code"`
	VerificationURI         string       `json:"verification_uri"`
	VerificationURIComplete string       `json:"verification_uri_complete"`
	ExpiresAt               time.Time    `json:"expires_at"`
}

type UserDeviceValues struct {
	UserCode                string
	VerificationURI         string
	VerificationURIComplete string
}

func (c *Client) CreateToken() (*UserDeviceValues, errors.Error) {
	response := &createTokenResponse{}
	if err := c.Request(
		WithMethod(http.MethodPost),
		WithPath("/auth/sessions/device"),
		WithResponseInto(response),
	); err != nil {
		return nil, err
	}

	if err := c.storeToken(response.Token); err != nil {
		return nil, errors.NewClientInternalError().WithCause(err)
	}

	return &UserDeviceValues{
		UserCode:                response.UserCode,
		VerificationURI:         response.VerificationURI,
		VerificationURIComplete: response.VerificationURIComplete,
	}, nil
}

func (c *Client) InvalidateToken() errors.Error {
	type deleteResponse struct {
		Success bool `json:"success"`
	}

	dr := &deleteResponse{}

	// Dont propagate error: if existing token is invalid endpoint will 401. Not sure this is
	// good behavior but it's true nonetheless
	c.Request(
		WithMethod(http.MethodDelete),
		WithPath("/auth/sessions"),
		WithResponseInto(dr),
	)

	if err := c.clearToken(); err != nil {
		return errors.NewClientInternalError().WithCause(err)
	}

	return nil
}
