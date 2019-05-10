package loader

import (
	"testing"

	"github.com/puppetlabs/nebula/pkg/workflow"
	"github.com/stretchr/testify/require"
)

func TestFilepathLoader(t *testing.T) {
	var l interface{} = NewFilepathLoader("./fixtures/test_filepath_loader/workflow.yaml")

	fpl, ok := l.(Loader)
	require.True(t, ok, "FilepathLoader does not satisfy Loader")

	var w workflow.Workflow

	require.NoError(t, fpl.Load(&w))

	require.Equal(t, "1", w.Version)

	require.Len(t, w.Variables, 2)
	for _, variable := range w.Variables {
		require.NotEmpty(t, variable.Name)
		require.NotEmpty(t, variable.Value)
	}

	require.Len(t, w.Steps, 2)
	require.Empty(t, w.Steps[0].DependsOn)
	require.NotEmpty(t, w.Steps[1].DependsOn)
}
