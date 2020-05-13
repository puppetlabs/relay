package def_test

import (
	"strings"
	"testing"

	"github.com/puppetlabs/relay/pkg/integration/container/def"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAdditionalImage(t *testing.T) {
	c, err := def.NewFromFilePath("testdata/additional_image.yaml")
	require.NoError(t, err)

	assert.NotNil(t, c.FileRef, "file reference not resolved")
	assert.Equal(t, "testdata", c.Container.Name)
	assert.Equal(t, "Jira resolve", c.Container.Title)
	assert.Equal(t, "A task that can update the state of a Jira ticket.", strings.TrimSpace(c.Container.Description))
	assert.Equal(t, "v1", c.Container.SDKVersion)
	assert.NotEmpty(t, c.Container.Images, "container has no images")
	assert.NotEmpty(t, c.Container.Settings, "container has no settings")
	assert.NotNil(t, c.Container.Images["base"], "container has no base image")
	assert.NotNil(t, c.Container.Images["bonus"], "container has no bonus image")
	assert.NotEmpty(t, c.Container.Settings["Image"].Description, "container has no description for the `Image` setting")
}

func TestNewExplicitName(t *testing.T) {
	c, err := def.NewFromFilePath("testdata/explicit_name.yaml")
	require.NoError(t, err)

	assert.NotNil(t, c.FileRef, "file reference not resolved")
	assert.Equal(t, "jira-resolve", c.Container.Name)
	assert.Equal(t, "v1", c.Container.SDKVersion)
}
