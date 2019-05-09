package client

import (
	"context"
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
	id := "some-workflow-id"
	repository := "test-repo"
	branch := "test-branch"
	path := "workflow.yaml"
	createdAt := time.Now().String()
	updatedAt := time.Now().String()

	wfm := &models.Workflow{
		ID:         &id,
		Repository: &repository,
		Branch:     &branch,
		Path:       &path,
		CreatedAt:  &createdAt,
		UpdatedAt:  &updatedAt,
	}

	routes := &testutil.MockRoutes{}
	routes.Add("/api/workflows", http.StatusCreated, wfm, nil)

	withAPIClient(t, routes, func(c *APIClient) {
		wf, err := c.CreateWorkflow(context.Background(), "test-repo", "test-branch", "workflow.yaml")
		require.NoError(t, err, "could not create workflow")
		require.Equal(t, *wf.ID, id)
	})
}
