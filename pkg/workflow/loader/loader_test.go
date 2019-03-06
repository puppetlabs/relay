package loader

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilepathLoader(t *testing.T) {
	var l interface{}
	l = NewFilepathLoader("./fixtures/test_filepath_loader/workflow.yaml")

	fpl, ok := l.(Loader)
	require.True(t, ok, "FilepathLoader does not satisfy Loader")

	wf, err := fpl.Load()
	require.NoError(t, err)

	require.Equal(t, "nebula-workflow-test", wf.Name)
	require.Equal(t, "1", wf.Version)

	require.Len(t, wf.Variables, 2)
	for _, variable := range wf.Variables {
		require.NotEmpty(t, variable.Name)
		require.NotEmpty(t, variable.Value)
	}

	require.Len(t, wf.Actions, 3)
}
