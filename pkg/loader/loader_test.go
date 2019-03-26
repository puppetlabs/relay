package loader

import (
	"testing"

	"github.com/puppetlabs/nebula/pkg/plan/types"
	"github.com/stretchr/testify/require"
)

func TestFilepathLoader(t *testing.T) {
	var l interface{} = NewFilepathLoader("./fixtures/test_filepath_loader/plan.yaml")

	fpl, ok := l.(Loader)
	require.True(t, ok, "FilepathLoader does not satisfy Loader")

	var p types.Plan

	require.NoError(t, fpl.Load(&p))

	require.Equal(t, "nebula-plan-test", p.Name)
	require.Equal(t, "1", p.Version)

	require.Len(t, p.Variables, 2)
	for _, variable := range p.Variables {
		require.NotEmpty(t, variable.Name)
		require.NotEmpty(t, variable.Value)
	}

	require.Len(t, p.Actions, 5)
}
