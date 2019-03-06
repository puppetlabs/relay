package loader

import (
	"io/ioutil"

	"github.com/puppetlabs/nebula/pkg/errors"
	"github.com/puppetlabs/nebula/pkg/workflow"
	"gopkg.in/yaml.v2"
)

type Loader interface {
	Load() (*workflow.Workflow, errors.Error)
}

type FilepathLoader struct {
	path string
}

func (f FilepathLoader) Load() (*workflow.Workflow, errors.Error) {
	b, err := ioutil.ReadFile(f.path)
	if err != nil {
		return nil, errors.NewWorkflowLoaderError().WithCause(err)
	}

	var wf workflow.Workflow

	if err := yaml.Unmarshal(b, &wf); err != nil {
		return nil, errors.NewWorkflowLoaderError().WithCause(err)
	}

	return &wf, nil
}

func NewFilepathLoader(path string) *FilepathLoader {
	return &FilepathLoader{path: path}
}
