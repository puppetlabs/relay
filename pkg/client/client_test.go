package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/puppetlabs/nebula-cli/pkg/client/api/auth"
	"github.com/puppetlabs/nebula-cli/pkg/client/api/models"
	"github.com/puppetlabs/nebula-cli/pkg/client/api/workflow_revisions"
	runs "github.com/puppetlabs/nebula-cli/pkg/client/api/workflow_runs"
	secrets "github.com/puppetlabs/nebula-cli/pkg/client/api/workflow_secrets"
	"github.com/puppetlabs/nebula-cli/pkg/client/api/workflows"
	"github.com/puppetlabs/nebula-cli/pkg/client/testutil"
	"github.com/puppetlabs/nebula-cli/pkg/config"
	"github.com/stretchr/testify/require"
)

func stringP(s string) *string {
	return &s
}

func fakeLogin(t *testing.T, c *APIClient) {
	f, err := os.OpenFile(c.cfg.TokenPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0750)
	require.NoError(t, err)

	_, err = f.Write([]byte("token"))
	require.NoError(t, err)
}

func makeWorkflowRevisionFixture() *models.WorkflowRevision {
	parameters := make(models.WorkflowParameters)
	parameters["key1"] = models.WorkflowParameter{Default: 1, Description: "Key 1"}
	parameters["key2"] = models.WorkflowParameter{Default: "test", Description: "Key 2"}
	parameters["key3"] = models.WorkflowParameter{Default: true, Description: "Key 3"}

	return &models.WorkflowRevision{
		WorkflowData: models.WorkflowData{Parameters: parameters},
	}
}

func makeWorkflowFixture(name, repository, branch, path string) *models.Workflow {
	now := strfmt.DateTime(time.Now())

	return &models.Workflow{
		WorkflowSummary: models.WorkflowSummary{
			WorkflowIdentifier: models.WorkflowIdentifier{
				Name: models.WorkflowName(name),
			},
		},
		Repository: &repository,
		Branch:     &branch,
		Path:       &path,
		Lifecycle: models.Lifecycle{
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}
}

func makeIntegrationFixture(accountLogin, provider string) *models.Integration {
	now := strfmt.DateTime(time.Now())

	return &models.Integration{
		IntegrationSummary: models.IntegrationSummary{
			Provider: &provider,
		},
		AccountLogin: accountLogin,
		Lifecycle: models.Lifecycle{
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}
}

func makeWorkflowRunFixture(wfm *models.Workflow) *models.WorkflowRun {
	runNum := int64(1)
	status := "pending"

	return &models.WorkflowRun{
		WorkflowRunSummary: models.WorkflowRunSummary{
			WorkflowRunIdentifier: models.WorkflowRunIdentifier{
				RunNumber: models.WorkflowRunNumber(runNum),
				Workflow:  &wfm.WorkflowIdentifier,
			},
			Status: &status,
		},
		Workflow: wfm,
	}
}

func withAPIClient(t *testing.T, routes http.Handler, fn func(c *APIClient)) {
	tmpdir, err := ioutil.TempDir("", "nebula-cli-test")
	require.NoError(t, err)

	defer os.RemoveAll(tmpdir)

	testutil.WithTestServer(routes, func(ts *httptest.Server) {
		c, err := NewAPIClient(&config.Config{
			APIHostAddr: ts.URL,
			TokenPath:   filepath.Join(tmpdir, "auth-token"),
		})

		require.NoError(t, err, "failed to setup api client using mock server")

		fn(c)
	})
}

func TestLogin(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		routes := &testutil.MockRoutes{}
		routes.Add("/auth/sessions", http.StatusCreated, &auth.CreateSessionCreatedBody{
			Token: stringP("mocktoken1234"),
		}, nil)

		withAPIClient(t, routes, func(c *APIClient) {
			require.NoError(t, c.Login(context.Background(), "test@example.com", "password1234"), "login failed")
		})
	})

	t.Run("unauthorized login", func(t *testing.T) {
		routes := &testutil.MockRoutes{}
		routes.Add("/auth/sessions", http.StatusUnauthorized, nil, nil)

		withAPIClient(t, routes, func(c *APIClient) {
			err := c.Login(context.Background(), "test@example.com", "password1234")
			require.Error(t, err, "login did not fail")
		})
	})
}

func TestWorkflowCreate(t *testing.T) {
	wfm := makeWorkflowFixture("name", "repo1", "branch1", "workflow.yaml")

	routes := &testutil.MockRoutes{}
	routes.Add("/api/workflows", http.StatusCreated, &workflows.CreateWorkflowCreatedBody{
		Workflow: wfm,
	}, nil)

	im := makeIntegrationFixture("test", "github")
	routes.Add("/api/integrations", http.StatusOK, im, nil)

	withAPIClient(t, routes, func(c *APIClient) {
		fakeLogin(t, c)

		wf, err := c.CreateWorkflow(context.Background(), "name", "description", "github-test", "repo1", "branch1", "workflow.yaml")
		require.NoError(t, err, "could not create workflow")
		require.Equal(t, wf.Name, models.WorkflowName("name"))
		require.Equal(t, *wf.Repository, "repo1")
	})
}

func TestWorkflowList(t *testing.T) {
	wfl := &workflows.GetWorkflowsOKBody{}

	for i := 0; i < 10; i++ {
		wfm := makeWorkflowFixture("name", "repo", "branch", fmt.Sprintf("workflow-%d.yaml", i))
		wfl.Workflows = append(wfl.Workflows, wfm)
	}

	routes := &testutil.MockRoutes{}
	routes.Add("/api/workflows", http.StatusOK, wfl, nil)

	withAPIClient(t, routes, func(c *APIClient) {
		fakeLogin(t, c)

		wfl, err := c.ListWorkflows(context.Background())
		require.NoError(t, err, "could not list workflows")

		for i := 0; i < 10; i++ {
			wf := wfl[i]
			require.Equal(t, wf.Name, models.WorkflowName("name"))
		}
	})
}

func TestWorkflowRevision(t *testing.T) {
	wfr := &workflow_revisions.GetLatestWorkflowRevisionOKBody{}
	wfr.Revision = makeWorkflowRevisionFixture()

	routes := &testutil.MockRoutes{}
	routes.Add("/api/workflows/name/revisions/latest", http.StatusOK, wfr, nil)

	withAPIClient(t, routes, func(c *APIClient) {
		fakeLogin(t, c)

		workflowRevision, err := c.GetLatestWorkflowRevision(context.Background(), "name")
		require.NoError(t, err, "could not list workflow parameters")

		for name, parameter := range workflowRevision.Parameters {
			require.EqualValues(t, parameter.Default, wfr.Revision.Parameters[name].Default)
			require.EqualValues(t, parameter.Description, wfr.Revision.Parameters[name].Description)
		}
	})
}

func TestWorkflowRun(t *testing.T) {
	wfm := makeWorkflowFixture("name", "repo1", "branch1", "workflow.yaml")
	wfrm := makeWorkflowRunFixture(wfm)

	routes := &testutil.MockRoutes{}
	routes.Add("/api/workflows/name/runs", http.StatusCreated, &runs.RunWorkflowCreatedBody{
		Run: wfrm,
	}, nil)

	withAPIClient(t, routes, func(c *APIClient) {
		fakeLogin(t, c)

		wfr, err := c.RunWorkflow(context.Background(), "name", nil)
		require.NoError(t, err, "could not run workflow")
		require.Equal(t, *wfr.Status, "pending")
		require.Equal(t, wfr.Workflow.Name, models.WorkflowName("name"))
	})
}

func TestCreateWorkflowSecret(t *testing.T) {
	ssm := &models.WorkflowSecretSummary{Key: "key"}

	routes := &testutil.MockRoutes{}
	routes.Add("/api/workflows/name/secrets", http.StatusCreated, &secrets.CreateWorkflowSecretCreatedBody{
		Secret: ssm,
	}, nil)

	withAPIClient(t, routes, func(c *APIClient) {
		fakeLogin(t, c)

		ssr, err := c.CreateWorkflowSecret(context.Background(), "name", "key", "value")
		require.NoError(t, err, "could not create secret")
		require.Equal(t, ssr.Key, "key")
	})
}

func TestUpdateWorkflowSecret(t *testing.T) {
	ssm := &models.WorkflowSecretSummary{Key: "key"}

	routes := &testutil.MockRoutes{}
	routes.Add("/api/workflows/name/secrets/key", http.StatusOK, &secrets.UpdateWorkflowSecretOKBody{
		Secret: ssm,
	}, nil)

	withAPIClient(t, routes, func(c *APIClient) {
		fakeLogin(t, c)

		ssr, err := c.UpdateWorkflowSecret(context.Background(), "name", "key", "value")
		require.NoError(t, err, "could not update secret")
		require.Equal(t, ssr.Key, "key")
	})
}
