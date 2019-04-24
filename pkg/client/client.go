package client

import (
	"bytes"
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/puppetlabs/nebula/pkg/client/api"
	authv1 "github.com/puppetlabs/nebula/pkg/client/api/auth_v1"
	"github.com/puppetlabs/nebula/pkg/client/api/models"
	workflowsv1 "github.com/puppetlabs/nebula/pkg/client/api/workflows_v1"
	"github.com/puppetlabs/nebula/pkg/config"
	"github.com/puppetlabs/nebula/pkg/errors"
)

const (
	defaultAPIHostURL = "https://api.nebula.puppet.com"
	defaultTokenFile  = "auth-token"
)

type APIClient struct {
	delegate    *api.Nebula
	cfg         *config.Config
	loadedToken *Token
}

func NewAPIClient(cfg *config.Config) (*APIClient, errors.Error) {
	addr := defaultAPIHostURL

	if cfg.APIHostAddr != "" {
		addr = cfg.APIHostAddr
	}

	host, err := url.Parse(addr)
	if err != nil {
		return nil, errors.NewAPIInvalidHost(addr).WithCause(err)
	}

	transport := api.DefaultTransportConfig()
	transport.Host = host.Host
	transport.Schemes = []string{host.Scheme}

	delegate := api.NewHTTPClientWithConfig(strfmt.Default, transport)

	return &APIClient{
		delegate: delegate,
		cfg:      cfg,
	}, nil
}

func (c *APIClient) Login(ctx context.Context, email string, password string) errors.Error {
	params := authv1.NewCreateSessionParams()
	params.SetBody(&models.NewSession{
		Email:    email,
		Password: password,
	})

	response, err := c.delegate.AuthV1.CreateSession(params)
	if err != nil {
		return errors.NewClientCreateSessionError().WithCause(err)
	}

	token := Token(strings.TrimPrefix(response.Authorization, "Bearer: "))

	if err := c.storeToken(ctx, &token); err != nil {
		return errors.NewClientCreateSessionError().WithCause(err)
	}

	return nil
}

func (c *APIClient) ListWorkflows(ctx context.Context) (*models.IndexWorkflows, errors.Error) {
	auth, err := c.getAuthorization(ctx)
	if err != nil {
		return nil, err
	}

	params := workflowsv1.NewListWorkflowsParams()
	params.Authorization = auth

	response, derr := c.delegate.WorkflowsV1.ListWorkflows(params)
	if derr != nil {
		return nil, errors.NewClientListWorkflowsError().WithCause(derr)
	}

	return &response.Payload, nil
}

func (c *APIClient) storeToken(ctx context.Context, token *Token) errors.Error {
	dest := filepath.Join(c.cfg.CachePath, defaultTokenFile)

	f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0750)
	if err != nil {
		return errors.NewClientTokenStorageError().WithCause(err).Bug()
	}

	defer f.Close()

	if _, err := f.Write([]byte(token.String())); err != nil {
		return errors.NewClientTokenStorageError().WithCause(err).Bug()
	}

	return nil
}

func (c *APIClient) getToken(ctx context.Context) (*Token, errors.Error) {
	if c.loadedToken == nil {
		dest := filepath.Join(c.cfg.CachePath, defaultTokenFile)

		f, err := os.Open(dest)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, errors.NewClientNotLoggedIn()
			}

			return nil, errors.NewClientTokenLoadError().WithCause(err).Bug()
		}

		defer f.Close()

		buf := &bytes.Buffer{}
		if _, err := buf.ReadFrom(f); err != nil {
			return nil, errors.NewClientTokenLoadError().WithCause(err).Bug()
		}

		token := Token(buf.String())

		c.loadedToken = &token

		return &token, nil
	}

	return c.loadedToken, nil
}

func (c *APIClient) getAuthorization(ctx context.Context) (string, errors.Error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return "", err
	}

	return token.Bearer(), nil
}
