package client

import (
	"bytes"
	"context"
	"net/url"
	"os"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/puppetlabs/nebula/pkg/client/api"
	authv1 "github.com/puppetlabs/nebula/pkg/client/api/auth_v1"
	integrationsv1 "github.com/puppetlabs/nebula/pkg/client/api/integrations_v1"
	"github.com/puppetlabs/nebula/pkg/client/api/models"
	workflowrunsv1 "github.com/puppetlabs/nebula/pkg/client/api/workflow_runs_v1"
	workflowsecretsv1 "github.com/puppetlabs/nebula/pkg/client/api/workflow_secrets_v1"
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
		return nil, errors.NewClientInvalidAPIHost(addr).WithCause(err)
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
	params.SetBody(&models.CreateSessionSubmission{
		Email:    &email,
		Password: &password,
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

func (c *APIClient) ListIntegrations(ctx context.Context) (*models.Integrations, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := integrationsv1.NewListIntegrationsParams()

	response, derr := c.delegate.IntegrationsV1.ListIntegrations(params, auth)
	if derr != nil {
		return nil, errors.NewClientListIntegrationsError().WithCause(derr)
	}

	return response.Payload, nil
}

func (c *APIClient) GetIntegration(ctx context.Context, id string) (*models.Integration, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := integrationsv1.NewGetIntegrationParams()
	params.ID = id

	response, derr := c.delegate.IntegrationsV1.GetIntegration(params, auth)
	if derr != nil {
		return nil, errors.NewClientGetIntegrationError(id).WithCause(derr)
	}

	return response.Payload, nil
}

func (c *APIClient) ListWorkflows(ctx context.Context) (*models.Workflows, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := workflowsv1.NewListWorkflowsParams()

	response, derr := c.delegate.WorkflowsV1.ListWorkflows(params, auth)
	if derr != nil {
		return nil, errors.NewClientListWorkflowsError().WithCause(derr)
	}

	return response.Payload, nil
}

func (c *APIClient) CreateWorkflow(ctx context.Context, name, description, integrationID, repo, branch, path string) (*models.Workflow, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := workflowsv1.NewCreateWorkflowParams()

	// TODO List integrations and get by provider_id-account_name
	params.Body = &models.WorkflowSubmission{
		Name:          models.WorkflowName(name),
		Description:   description,
		IntegrationID: &integrationID,
		Repository:    &repo,
		Branch:        &branch,
		Path:          &path,
	}

	resp, werr := c.delegate.WorkflowsV1.CreateWorkflow(params, auth)
	if werr != nil {
		return nil, errors.NewClientCreateWorkflowError().WithCause(werr)
	}

	return resp.Payload, nil
}

func (c *APIClient) RunWorkflow(ctx context.Context, name string) (*models.WorkflowRun, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := workflowrunsv1.NewCreateWorkflowRunParamsWithContext(ctx)
	params.WorkflowName = name

	resp, werr := c.delegate.WorkflowRunsV1.CreateWorkflowRun(params, auth)
	if werr != nil {
		return nil, errors.NewClientRunWorkflowError().WithCause(werr)
	}

	return resp.Payload, nil
}

func (c *APIClient) ListWorkflowRuns(ctx context.Context, name string) (*models.WorkflowRunSummaries, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := workflowrunsv1.NewListWorkflowRunsParams()
	params.WorkflowName = name

	resp, werr := c.delegate.WorkflowRunsV1.ListWorkflowRuns(params, auth)
	if werr != nil {
		return nil, errors.NewClientRunWorkflowError().WithCause(werr)
	}

	return resp.Payload, nil
}

func (c *APIClient) CreateWorkflowSecret(ctx context.Context, name, key, value string) (*models.SecretSummary, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := workflowsecretsv1.NewCreateWorkflowSecretParams()
	params.WorkflowName = name
	params.Body = &models.CreateSecret{
		Key:   key,
		Value: value,
	}

	resp, werr := c.delegate.WorkflowSecretsV1.CreateWorkflowSecret(params, auth)
	if _, ok := werr.(*workflowsecretsv1.CreateWorkflowSecretConflict); ok {
		return nil, errors.NewClientWorkflowSecretAlreadyExistsError(key)
	} else if werr != nil {
		return nil, errors.NewClientCreateWorkflowSecretError().WithCause(werr)
	}

	return resp.Payload, nil
}

func (c *APIClient) UpdateWorkflowSecret(ctx context.Context, name, key, value string) (*models.SecretSummary, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := workflowsecretsv1.NewUpdateSecretByKeyAndWorkflowIDParams()
	params.WorkflowName = name
	params.SecretKey = key
	params.Body = &models.UpdateSecret{
		Value: value,
	}

	resp, werr := c.delegate.WorkflowSecretsV1.UpdateSecretByKeyAndWorkflowID(params, auth)
	if werr != nil {
		return nil, errors.NewClientUpdateWorkflowSecretError().WithCause(werr)
	}

	return resp.Payload, nil
}

func (c *APIClient) storeToken(ctx context.Context, token *Token) errors.Error {
	f, err := os.OpenFile(c.cfg.TokenPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0750)
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
		f, err := os.Open(c.cfg.TokenPath)
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

func (c *APIClient) getAuthorizationFunc(ctx context.Context) runtime.ClientAuthInfoWriterFunc {
	return func(req runtime.ClientRequest, reg strfmt.Registry) error {
		token, err := c.getToken(ctx)
		if err != nil {
			return err
		}

		return req.SetHeaderParam("Authorization", token.Bearer())
	}
}
