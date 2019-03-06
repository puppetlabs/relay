package loader

import (
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
)

func TestFilepathLoader(t *testing.T) {
	var l interface{}
	l = NewFilepathLoader("./fixtures/test_filepath_loader/workflow.yaml")

	fpl, ok := l.(Loader)
	require.True(t, ok, "FilepathLoader does not satisfy Loader")

	wf, err := fpl.Load()
	require.NoError(t, err)

	pretty.Println(wf)
}
