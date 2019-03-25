package loader

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/puppetlabs/nebula/pkg/errors"
	"github.com/puppetlabs/nebula/pkg/plan/types"
	"gopkg.in/yaml.v2"
)

type Loader interface {
	Load(*types.Plan) errors.Error
}

type FilepathLoader struct {
	path string
}

func (f FilepathLoader) Load(p *types.Plan) errors.Error {
	b, err := ioutil.ReadFile(f.path)
	if err != nil {
		return errors.NewPlanLoaderError().WithCause(err)
	}

	if err := yaml.Unmarshal(b, p); err != nil {
		return errors.NewPlanLoaderError().WithCause(err)
	}

	return nil
}

func NewFilepathLoader(path string) *FilepathLoader {
	return &FilepathLoader{path: path}
}

type ImpliedPlanFileLoader struct{}

func (i ImpliedPlanFileLoader) Load(p *types.Plan) errors.Error {
	impliedPath := filepath.Join(".", "plan.yaml")

	if _, err := os.Stat(impliedPath); err != nil {
		if os.IsNotExist(err) {
			return errors.NewPlanFileNotFound(impliedPath)
		}

		return errors.NewPlanLoaderError().WithCause(err).Bug()
	}

	delegate := NewFilepathLoader(impliedPath)

	return delegate.Load(p)
}
