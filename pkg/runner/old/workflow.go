package runner

import (
	"context"

	"github.com/puppetlabs/nebula/pkg/errors"
	"gopkg.in/yaml.v2"
)

type WorkflowSpec struct {
	Import string `yaml:"import"`
}

type Workflow struct {
	Name string       `yaml:"name"`
	Spec WorkflowSpec `yaml:"spec"`
}

func (w *Workflow) Run(ctx context.Context, rid string, r ActionRuntime, variables map[string]string) errors.Error {
	return nil
}

func (w *Workflow) Decoder() Decoder {
	return &WorkflowDecoder{w: w}
}

type WorkflowDecoder struct {
	w *Workflow
}

func (d *WorkflowDecoder) Decode(b []byte) errors.Error {
	if err := yaml.Unmarshal(b, d.w); err != nil {
		return errors.NewWorkflowRunnerDecodeError().WithCause(err).Bug()
	}

	return nil
}
