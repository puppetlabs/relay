package client

import (
	"bytes"
	"context"
	"io"
	"net/url"
	"os"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/puppetlabs/horsehead/v2/encoding/transfer"
	"github.com/puppetlabs/nebula-cli/pkg/client/api"
	"github.com/puppetlabs/nebula-cli/pkg/client/api/auth"
	"github.com/puppetlabs/nebula-cli/pkg/client/api/events"
	"github.com/puppetlabs/nebula-cli/pkg/client/api/integrations"
	"github.com/puppetlabs/nebula-cli/pkg/client/api/models"
	"github.com/puppetlabs/nebula-cli/pkg/client/api/workflow_revisions"
	runs "github.com/puppetlabs/nebula-cli/pkg/client/api/workflow_runs"
	secrets "github.com/puppetlabs/nebula-cli/pkg/client/api/workflow_secrets"
	"github.com/puppetlabs/nebula-cli/pkg/client/api/workflows"
	"github.com/puppetlabs/nebula-cli/pkg/config"
	"github.com/puppetlabs/nebula-cli/pkg/errors"
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

	transport := httptransport.New(host.Host, "/", []string{host.Scheme})
	transport.Producers["application/vnd.puppet.nebula.v20200131+json"] = runtime.JSONProducer()
	transport.Consumers["application/vnd.puppet.nebula.v20200131+json"] = runtime.JSONConsumer()

	delegate := api.New(transport, strfmt.Default)

	return &APIClient{
		delegate: delegate,
		cfg:      cfg,
	}, nil
}

func (c *APIClient) Login(ctx context.Context, email string, password string) errors.Error {
	params := auth.NewCreateSessionParams()
	params.SetBody(auth.CreateSessionBody{
		Email:    &email,
		Password: &password,
	})

	response, err := c.delegate.Auth.CreateSession(params)
	if err != nil {
		return errors.NewClientCreateSessionError().WithCause(translateRuntimeError(err))
	}

	token := Token(*response.Payload.Token)

	if err := c.storeToken(ctx, &token); err != nil {
		return errors.NewClientCreateSessionError().WithCause(err)
	}

	return nil
}

func (c *APIClient) ListIntegrations(ctx context.Context) ([]*models.Integration, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := integrations.NewGetIntegrationsParams()

	response, derr := c.delegate.Integrations.GetIntegrations(params, auth)
	if derr != nil {
		return nil, errors.NewClientListIntegrationsError().WithCause(translateRuntimeError(derr))
	}

	return response.Payload.Integrations, nil
}

func (c *APIClient) ListIntegrationRepositoryBranches(ctx context.Context, id, name, owner string, query *string) ([]*models.RepositoryBranch, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := integrations.NewGetIntegrationRepositoryBranchesParamsWithContext(ctx)

	params.IntegrationID = id
	params.IntegrationRepositoryName = name
	params.IntegrationRepositoryOwner = owner
	params.Q = query

	response, derr := c.delegate.Integrations.GetIntegrationRepositoryBranches(params, auth)
	if derr != nil {
		return nil, errors.NewClientListIntegrationsError().WithCause(translateRuntimeError(derr))
	}

	return response.Payload.Branches, nil
}

func (c *APIClient) GetIntegration(ctx context.Context, id string) (*models.Integration, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := integrations.NewGetIntegrationParams()
	params.IntegrationID = id

	response, derr := c.delegate.Integrations.GetIntegration(params, auth)
	if derr != nil {
		return nil, errors.NewClientGetIntegrationError(id).WithCause(translateRuntimeError(derr))
	}

	return response.Payload.Integration, nil
}

func (c *APIClient) ListEventSources(ctx context.Context) ([]*models.EventSource, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := events.NewGetEventSourcesParams()

	response, derr := c.delegate.Events.GetEventSources(params, auth)
	if derr != nil {
		return nil, errors.NewClientListIntegrationsError().WithCause(translateRuntimeError(derr))
	}

	return response.Payload.EventSources, nil
}

func (c *APIClient) GetEventSource(ctx context.Context, id string) (*models.EventSource, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := events.NewGetEventSourceParams()
	params.EventSourceID = id

	response, derr := c.delegate.Events.GetEventSource(params, auth)
	if derr != nil {
		return nil, errors.NewClientGetIntegrationError(id).WithCause(translateRuntimeError(derr))
	}

	return response.Payload.EventSource, nil
}

func (c *APIClient) ListWorkflows(ctx context.Context) ([]*models.Workflow, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := workflows.NewGetWorkflowsParams()

	response, derr := c.delegate.Workflows.GetWorkflows(params, auth)
	if derr != nil {
		return nil, errors.NewClientListWorkflowsError().WithCause(translateRuntimeError(derr))
	}

	return response.Payload.Workflows, nil
}

func (c *APIClient) GetWorkflow(ctx context.Context, name string) (*models.Workflow, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := workflows.NewGetWorkflowParams()
	params.WorkflowName = name

	response, derr := c.delegate.Workflows.GetWorkflow(params, auth)
	if derr != nil {
		return nil, errors.NewClientGetWorkflowError(name).WithCause(translateRuntimeError(derr))
	}

	return response.Payload.Workflow, nil
}

func (c *APIClient) CreateWorkflow(ctx context.Context, name, description string, content io.ReadCloser) (*models.Workflow, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := workflows.NewCreateWorkflowParams()
	params.SetBody(workflows.CreateWorkflowBody{
		Name:        models.WorkflowName(name),
		Description: description,
	})

	wfresp, werr := c.delegate.Workflows.CreateWorkflow(params, auth)
	if werr != nil {
		return nil, errors.NewClientCreateWorkflowError().WithCause(translateRuntimeError(werr))
	}

	wf := wfresp.Payload.Workflow

	revisionParams := workflow_revisions.NewPostWorkflowRevisionParams()
	revisionParams.SetBody(content)
	revisionParams.WorkflowName = string(wf.Name)

	_, wrerr := c.delegate.WorkflowRevisions.PostWorkflowRevision(revisionParams, auth)
	if wrerr != nil {
		return nil, errors.NewClientCreateWorkflowRevisionError().WithCause(translateRuntimeError(wrerr))
	}

	return wf, nil
}

func (c *APIClient) UpdateWorkflow(ctx context.Context, name, description string, content io.ReadCloser) (*models.Workflow, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	wf, err := c.GetWorkflow(ctx, name)
	if err != nil {
		return nil, err
	}

	if description != wf.Description {
		params := workflows.NewUpdateWorkflowParams()
		params.WorkflowName = name
		params.SetBody(workflows.UpdateWorkflowBody{
			Description: description,
		})

		response, err := c.delegate.Workflows.UpdateWorkflow(params, auth)
		if err != nil {
			return nil, errors.NewClientUpdateWorkflowError(name)
		}

		wf = response.Payload.Workflow
	}

	revisionParams := workflow_revisions.NewPostWorkflowRevisionParams()
	revisionParams.SetBody(content)
	revisionParams.WorkflowName = string(wf.Name)

	_, wrerr := c.delegate.WorkflowRevisions.PostWorkflowRevision(revisionParams, auth)
	if wrerr != nil {
		return nil, errors.NewClientCreateWorkflowRevisionError().WithCause(translateRuntimeError(wrerr))
	}

	return wf, nil
}

func (c *APIClient) RunWorkflow(ctx context.Context, name string, parameters map[string]string) (*models.WorkflowRun, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	wrp := make(models.WorkflowRunParameters, len(parameters))

	for name, value := range parameters {
		ev, err := transfer.EncodeJSON([]byte(value))
		if err != nil {
			return nil, errors.NewClientInvalidWorkflowParameterValueError().WithCause(err).Bug()
		}

		wrp[name] = models.WorkflowRunParameter{Value: ev}
	}

	params := runs.NewRunWorkflowParamsWithContext(ctx)
	params.WorkflowName = name

	params.SetBody(runs.RunWorkflowBody{
		Parameters: wrp,
	})

	resp, werr := c.delegate.WorkflowRuns.RunWorkflow(params, auth)
	if werr != nil {
		return nil, errors.NewClientRunWorkflowError().WithCause(translateRuntimeError(werr))
	}

	return resp.Payload.Run, nil
}

func (c *APIClient) CancelWorkflowRun(ctx context.Context, name string, runNum int64) errors.Error {
	auth := c.getAuthorizationFunc(ctx)

	params := runs.NewPatchWorkflowRunParams()
	params.WorkflowName = name
	params.WorkflowRunNumber = runNum

	params.SetBody(runs.PatchWorkflowRunBody{
		Operation: &models.WorkflowRunOperation{
			Cancel: true,
		},
	})

	_, werr := c.delegate.WorkflowRuns.PatchWorkflowRun(params, auth)
	if werr != nil {
		return errors.NewClientRunWorkflowError().WithCause(translateRuntimeError(werr))
	}

	return nil
}

func (c *APIClient) ListWorkflowRuns(ctx context.Context, name string) ([]*models.WorkflowRunSummary, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := runs.NewGetWorkflowRunsParams()
	params.WorkflowName = name

	resp, werr := c.delegate.WorkflowRuns.GetWorkflowRuns(params, auth)
	if werr != nil {
		return nil, errors.NewClientListWorkflowRunsError().WithCause(translateRuntimeError(werr))
	}

	return resp.Payload.Runs, nil
}

func (c *APIClient) GetLatestWorkflowRevision(ctx context.Context, name string) (*models.WorkflowRevision, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := workflow_revisions.NewGetLatestWorkflowRevisionParams()
	params.WorkflowName = name

	resp, werr := c.delegate.WorkflowRevisions.GetLatestWorkflowRevision(params, auth)
	if werr != nil {
		return nil, errors.NewClientGetWorkflowRevisionError().WithCause(translateRuntimeError(werr))
	}

	return resp.Payload.Revision, nil
}

func (c *APIClient) GetWorkflowRun(ctx context.Context, name string, runNum int64) (*models.WorkflowRun, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := runs.NewGetWorkflowRunParams()
	params.WorkflowName = name
	params.WorkflowRunNumber = runNum

	resp, werr := c.delegate.WorkflowRuns.GetWorkflowRun(params, auth)
	if werr != nil {
		return nil, errors.NewClientGetWorkflowRunError().WithCause(translateRuntimeError(werr))
	}

	return resp.Payload.Run, nil
}

func (c *APIClient) GetWorkflowRunStepLog(ctx context.Context, name string, runNum int64, step string, follow bool, writer io.Writer) errors.Error {
	auth := c.getAuthorizationFunc(ctx)

	params := runs.NewGetWorkflowRunStepLogParamsWithContext(ctx)
	params.WorkflowName = name
	params.WorkflowRunNumber = runNum
	params.WorkflowStepName = step
	params.Follow = &follow

	_, _, werr := c.delegate.WorkflowRuns.GetWorkflowRunStepLog(params, auth, writer)
	if werr != nil {
		return errors.NewClientGetWorkflowRunStepLogError().WithCause(translateRuntimeError(werr))
	}

	return nil
}

func (c *APIClient) CreateWorkflowSecret(ctx context.Context, name, key, value string) (*models.WorkflowSecretSummary, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	ev, err := transfer.EncodeJSON([]byte(value))
	if err != nil {
		return nil, errors.NewClientInvalidWorkflowSecretValueError().WithCause(err).Bug()
	}

	params := secrets.NewCreateWorkflowSecretParams()
	params.WorkflowName = name
	params.SetBody(secrets.CreateWorkflowSecretBody{
		Name:  &key,
		Value: ev,
	})

	resp, err := c.delegate.WorkflowSecrets.CreateWorkflowSecret(params, auth)
	if werr := translateRuntimeError(err); werr != nil {
		if werr.Is("napi_secret_create_conflict") {
			return nil, errors.NewClientWorkflowSecretAlreadyExistsError(key)
		}

		return nil, errors.NewClientCreateWorkflowSecretError().WithCause(werr)
	}

	return resp.Payload.Secret, nil
}

func (c *APIClient) UpdateWorkflowSecret(ctx context.Context, name, key, value string) (*models.WorkflowSecretSummary, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	ev, err := transfer.EncodeJSON([]byte(value))
	if err != nil {
		return nil, errors.NewClientInvalidWorkflowSecretValueError().WithCause(err).Bug()
	}

	params := secrets.NewUpdateWorkflowSecretParams()
	params.WorkflowName = name
	params.WorkflowSecretName = key
	params.SetBody(secrets.UpdateWorkflowSecretBody{
		Value: ev,
	})

	resp, werr := c.delegate.WorkflowSecrets.UpdateWorkflowSecret(params, auth)
	if werr != nil {
		return nil, errors.NewClientUpdateWorkflowSecretError().WithCause(translateRuntimeError(werr))
	}

	return resp.Payload.Secret, nil
}

func (c *APIClient) DeleteWorkflowSecret(ctx context.Context, name, key string) errors.Error {
	auth := c.getAuthorizationFunc(ctx)

	params := secrets.NewDeleteWorkflowSecretParams()
	params.WorkflowName = name
	params.WorkflowSecretName = key

	_, werr := c.delegate.WorkflowSecrets.DeleteWorkflowSecret(params, auth)
	if werr != nil {
		return errors.NewClientDeleteWorkflowSecretError().WithCause(translateRuntimeError(werr))
	}

	return nil
}

func (c *APIClient) ListWorkflowSecrets(ctx context.Context, name string) ([]*models.WorkflowSecretSummary, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := secrets.NewListWorkflowSecretsParams()
	params.WorkflowName = name

	resp, werr := c.delegate.WorkflowSecrets.ListWorkflowSecrets(params, auth)
	if werr != nil {
		return nil, errors.NewClientListWorkflowRunsError().WithCause(translateRuntimeError(werr))
	}

	return resp.Payload.Secrets, nil
}

func (c *APIClient) GetWorkflowRevision(ctx context.Context, name string, revisionId string) (*models.WorkflowRevision, errors.Error) {
	auth := c.getAuthorizationFunc(ctx)

	params := workflow_revisions.NewGetWorkflowRevisionParams()
	params.WorkflowName = name
	params.WorkflowRevision = revisionId

	resp, werr := c.delegate.WorkflowRevisions.GetWorkflowRevision(params, auth)
	if werr != nil {
		return nil, errors.NewClientListWorkflowRunsError().WithCause(translateRuntimeError(werr))
	}

	return resp.Payload.Revision, nil
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
