package def_test

import (
	"testing"

	"github.com/puppetlabs/relay/pkg/integration/container/def"
	"github.com/stretchr/testify/assert"
)

func TestFileRefLocal(t *testing.T) {
	tt := []struct {
		Ref      *def.FileRef
		Expected bool
	}{
		{def.NewFileRef("foo", def.WithFileRefResolver(def.SDKResolver)), false},
		{def.NewFileRef("foo"), true},
		{def.NewFileRef("foo", def.WithFileRefResolver(def.DefaultResolver)), true},
	}
	for _, test := range tt {
		assert.Equal(t, test.Expected, test.Ref.Local())
	}
}

func TestFileRefDir(t *testing.T) {
	tt := []struct {
		Name           string
		Ref            *def.FileRef
		ExpectedString string
	}{
		{
			"without directory, SDK resolver",
			def.NewFileRef("foo", def.WithFileRefResolver(def.SDKResolver)),
			"sdk:templates",
		},
		{
			"without directory, local resolver",
			def.NewFileRef("foo"),
			".",
		},
		{
			"in root directory, local resolver",
			def.NewFileRef("/foo"),
			"/",
		},
		{
			"in child directory, local resolver",
			def.NewFileRef("foo/bar"),
			"foo",
		},
	}
	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.ExpectedString, test.Ref.Dir().String())
		})
	}
}
