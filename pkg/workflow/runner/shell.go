package runner

import (
	"context"

	"github.com/puppetlabs/nebula/pkg/errors"
	"github.com/puppetlabs/nebula/pkg/execution"
	"gopkg.in/yaml.v2"
)

type ShellSpec struct {
	Kind     string   `yaml:"kind"`
	Commands []string `yaml:"commands"`
}

type Shell struct {
	Name string    `yaml:"name"`
	Spec ShellSpec `yaml:"spec"`
}

func (s *Shell) Run(ctx context.Context, rid string, r ActionRuntime, variables map[string]string) errors.Error {
	for _, command := range s.Spec.Commands {
		_, err := execution.ExecuteCommand(command, variables, r.Logger())

		if err != nil {
			return errors.NewWorkflowUnknownRuntimeError().WithCause(err)
		}
	}

	return nil
}

func (s *Shell) Decoder() Decoder {
	return &ShellDecoder{s: s}
}

type ShellDecoder struct {
	s *Shell
}

func (d *ShellDecoder) Decode(b []byte) errors.Error {
	if err := yaml.Unmarshal(b, d.s); err != nil {
		return errors.NewWorkflowRunnerDecodeError().WithCause(err).Bug()
	}

	return nil
}
