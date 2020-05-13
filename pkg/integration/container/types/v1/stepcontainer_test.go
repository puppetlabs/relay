package v1_test

import (
	"io/ioutil"
	"testing"

	v1 "github.com/puppetlabs/relay/pkg/integration/container/types/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStepContainerValid(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/container_valid.yaml")
	require.NoError(t, err)

	sct, err := v1.NewStepContainerFromString(string(b))
	require.NoError(t, err)

	assert.Equal(t, v1.Version, sct.APIVersion)
	assert.Equal(t, v1.StepContainerKind, sct.Kind)
	assert.NotEmpty(t, sct.Title, "container has no title")
	assert.NotEmpty(t, sct.Description, "container has no description")
}
