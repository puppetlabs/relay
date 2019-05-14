package loader

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/puppetlabs/nebula/pkg/errors"
	"github.com/puppetlabs/nebula/pkg/workflow"
	"gopkg.in/yaml.v2"
)

type Loader interface {
	Load(*workflow.Workflow) errors.Error
}

type FilepathLoader struct {
	path string
}

func (f FilepathLoader) Load(p *workflow.Workflow) errors.Error {
	b, err := ioutil.ReadFile(f.path)
	if err != nil {
		return errors.NewWorkflowLoaderError().WithCause(err)
	}

	if err := yaml.Unmarshal(b, p); err != nil {
		return errors.NewWorkflowLoaderError().WithCause(err)
	}

	return nil
}

func NewFilepathLoader(path string) *FilepathLoader {
	return &FilepathLoader{path: path}
}

// ImpliedPlanFileLoader is the old plan loader
// this will become useful when we flesh out multi-workflow support
// and a better way of structuring the plans.
type ImpliedPlanFileLoader struct{}

func (i ImpliedPlanFileLoader) Load(p *workflow.Workflow) errors.Error {
	impliedPath := filepath.Join(".", "workflow.yaml")

	if _, err := os.Stat(impliedPath); err != nil {
		if os.IsNotExist(err) {
			return errors.NewWorkflowFileNotFound(impliedPath)
		}

		return errors.NewWorkflowLoaderError().WithCause(err).Bug()
	}

	delegate := NewFilepathLoader(impliedPath)

	return delegate.Load(p)
}
