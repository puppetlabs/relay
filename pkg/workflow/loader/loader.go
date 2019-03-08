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

	err = prepareStages(&wf)

	if err != nil {
		return nil, errors.NewWorkflowStageError().WithCause(err)
	}

	return &wf, nil
}

func prepareStages(w *workflow.Workflow) errors.Error {
	actionMap := make(map[string]workflow.Action)

	for _, action := range w.Actions {
		actionMap[action.Name] = action
	}

	for _, step := range w.Stage.Steps {
		// 1. Validate the step is the name of a valid action
		thisAction, ok := actionMap[step]

		if !ok {
			return errors.NewWorkflowNonExistentActionError(step)
		}

		// 2. Add this step to the stage
		w.Stage.Actions = append(w.Stage.Actions, thisAction)
	}

	return nil
}

func NewFilepathLoader(path string) *FilepathLoader {
	return &FilepathLoader{path: path}
}

type ImpliedWorkflowFileLoader struct{}

func (i ImpliedWorkflowFileLoader) Load() (*workflow.Workflow, errors.Error) {
	impliedPath := filepath.Join(".", "workflow.yaml")

	if _, err := os.Stat(impliedPath); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.NewWorkflowFileNotFound(impliedPath)
		}

		return nil, errors.NewWorkflowLoaderError().WithCause(err).Bug()
	}

	delegate := NewFilepathLoader(impliedPath)

	return delegate.Load()
}
