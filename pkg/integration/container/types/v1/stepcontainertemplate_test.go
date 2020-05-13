package v1_test

import (
	"fmt"
	"testing"

	"github.com/puppetlabs/relay/pkg/integration/container/asset"
	v1 "github.com/puppetlabs/relay/pkg/integration/container/types/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStepContainerTemplatesValid(t *testing.T) {
	for _, name := range []string{"bash.v1", "go.v1"} {
		t.Run(name, func(t *testing.T) {
			s, err := asset.AssetString(fmt.Sprintf("templates/%s/container.yaml", name))
			require.NoError(t, err)

			sctt, err := v1.NewStepContainerTemplateFromString(s)
			require.NoError(t, err)

			assert.Equal(t, v1.Version, sctt.APIVersion)
			assert.Equal(t, v1.StepContainerTemplateKind, sctt.Kind)
			assert.NotEmpty(t, sctt.Images, "template has no images")
			assert.NotEmpty(t, sctt.Settings, "template has no settings")
		})
	}
}
