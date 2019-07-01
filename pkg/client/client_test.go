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

	"github.com/puppetlabs/nebula/pkg/client/api/models"
	"github.com/puppetlabs/nebula/pkg/client/testutil"
	"github.com/puppetlabs/nebula/pkg/config"
	"github.com/puppetlabs/nebula/pkg/workflow"
	"github.com/stretchr/testify/require"
)

func fakeLogin(t *testing.T, c *APIClient) {
	f, err := os.OpenFile(c.cfg.TokenPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0750)
	require.NoError(t, err)

	_, err = f.Write([]byte("token"))
	require.NoError(t, err)
}

func makeWorkflowFileFixture() *workflow.Workflow {
	return &workflow.Workflow{}
}

func makeWorkflowFixture(name, repository, branch, path string) *models.Workflow {
	createdAt := time.Now().String()
	updatedAt := time.Now().String()

	return &models.Workflow{
		Name:       models.WorkflowName(name),
		Repository: &repository,
		Branch:     &branch,
		Path:       &path,
		CreatedAt:  &createdAt,
		UpdatedAt:  &updatedAt,
	}
}

func makeWorkflowRunFixture(wfm *models.Workflow) *models.WorkflowRun {
	id := "wfr-1"
	runNum := int64(1)
	status := "pending"

	return &models.WorkflowRun{
		ID:        &id,
		RunNumber: &runNum,
		Status:    &status,
		Workflow:  wfm,
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
		header := make(map[string]string)
		header["Authorization"] = "Bearer: mocktoken1234"

		routes := &testutil.MockRoutes{}
		routes.Add("/auth/sessions", http.StatusOK, &models.GenericSuccess{Success: true}, header)

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
	routes.Add("/api/workflows", http.StatusCreated, wfm, nil)

	withAPIClient(t, routes, func(c *APIClient) {
		fakeLogin(t, c)

		wf, err := c.CreateWorkflow(context.Background(), "name", "description", "repo1", "branch1", "workflow.yaml")
		require.NoError(t, err, "could not create workflow")
		require.Equal(t, wf.Name, models.WorkflowName("name"))
		require.Equal(t, *wf.Repository, "repo1")
	})
}

func TestWorkflowList(t *testing.T) {
	wfl := &models.Workflows{}

	for i := 0; i < 10; i++ {
		wfm := makeWorkflowFixture("name", "repo", "branch", fmt.Sprintf("workflow-%d.yaml", i))
		wfl.Items = append(wfl.Items, wfm)
	}

	routes := &testutil.MockRoutes{}
	routes.Add("/api/workflows", http.StatusOK, wfl, nil)

	withAPIClient(t, routes, func(c *APIClient) {
		fakeLogin(t, c)

		wfl, err := c.ListWorkflows(context.Background())
		require.NoError(t, err, "could not list workflows")

		for i := 0; i < 10; i++ {
			wf := wfl.Items[i]
			require.Equal(t, wf.Name, models.WorkflowName("name"))
		}
	})
}

func TestWorkflowRun(t *testing.T) {
	wfm := makeWorkflowFixture("name", "repo1", "branch1", "workflow.yaml")
	wfrm := makeWorkflowRunFixture(wfm)

	routes := &testutil.MockRoutes{}
	routes.Add("/api/workflows/name/runs", http.StatusCreated, wfrm, nil)

	withAPIClient(t, routes, func(c *APIClient) {
		fakeLogin(t, c)

		wfr, err := c.RunWorkflow(context.Background(), "name")
		require.NoError(t, err, "could not run workflow")
		require.Equal(t, *wfr.Status, "pending")
		require.Equal(t, *wfr.ID, "wfr-1")
		require.Equal(t, wfr.Workflow.Name, models.WorkflowName("name"))
	})
}
