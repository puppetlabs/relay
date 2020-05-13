package generator_test

import (
	"strings"
	"testing"

	"github.com/puppetlabs/relay/pkg/integration/container/def"
	"github.com/puppetlabs/relay/pkg/integration/container/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratorFiles(t *testing.T) {
	tpl, err := def.NewTemplateFromFileRef(def.NewFileRef("bash.v1", def.WithFileRefResolver(def.SDKResolver)))
	require.NoError(t, err)

	c := &def.Container{
		Common:      tpl.Common,
		ID:          "abcdef123456",
		Name:        "test",
		Title:       "Test",
		Description: "The test task does the best testing.",
	}
	c.Settings["AdditionalPackages"].Value = []string{"xmlstarlet"}
	c.Settings["AdditionalCommands"].Value = []string{"do\nmy\nbidding"}

	g := generator.New(c, generator.WithRepoNameBase("foo/bar"), generator.WithScriptFilename("foo.sh"))

	fs, err := g.Files()
	require.NoError(t, err)

	assert.Len(t, fs, 2)
	for _, f := range fs {
		if strings.HasPrefix(f.Ref.String(), "Dockerfile") {
			assert.Equal(t, "Dockerfile", f.Ref.String())
		} else {
			assert.Equal(t, "foo.sh", f.Ref.String())
		}
	}
}
