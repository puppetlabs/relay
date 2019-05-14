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
	"github.com/stretchr/testify/require"
)

func makeWorkflowFixture(id, repository, branch, path string) *models.Workflow {
	createdAt := time.Now().String()
	updatedAt := time.Now().String()

	return &models.Workflow{
		ID:         &id,
		Repository: &repository,
		Branch:     &branch,
		Path:       &path,
		CreatedAt:  &createdAt,
		UpdatedAt:  &updatedAt,
	}
}

func withAPIClient(t *testing.T, routes http.Handler, fn func(c *APIClient)) {
	tmpdir, err := ioutil.TempDir("", "nebula-cli-test")
	require.NoError(t, err)

	defer func() {
		os.RemoveAll(tmpdir)
	}()

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
	wfm := makeWorkflowFixture("id1234", "repo1", "branch1", "workflow.yaml")

	routes := &testutil.MockRoutes{}
	routes.Add("/api/workflows", http.StatusCreated, wfm, nil)

	withAPIClient(t, routes, func(c *APIClient) {
		wf, err := c.CreateWorkflow(context.Background(), "repo1", "branch1", "workflow.yaml")
		require.NoError(t, err, "could not create workflow")
		require.Equal(t, *wf.ID, "id1234")
		require.Equal(t, *wf.Repository, "repo1")
	})
}

func TestWorkflowList(t *testing.T) {
	wfl := &models.Workflows{}

	for i := 0; i < 10; i++ {
		wfm := makeWorkflowFixture(fmt.Sprintf("id-%d", i), "repo", "branch", fmt.Sprintf("workflow-%d.yaml", i))
		wfl.Items = append(wfl.Items, wfm)
	}

	routes := &testutil.MockRoutes{}
	routes.Add("/api/workflows", http.StatusOK, wfl, nil)

	withAPIClient(t, routes, func(c *APIClient) {
		wfl, err := c.ListWorkflows(context.Background())
		require.NoError(t, err, "could not list workflows")

		for i := 0; i < 10; i++ {
			wf := wfl.Items[i]
			require.Equal(t, *wf.ID, fmt.Sprintf("id-%d", i))
		}
	})
}
