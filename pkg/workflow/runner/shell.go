package runner

import (
	"context"

	"github.com/puppetlabs/nebula/pkg/errors"
)

type ShellSpec struct {
	Kind     string   `yaml:"kind"`
	Commands []string `yaml:"commands"`
}

type Shell struct {
	Name string    `yaml:"name"`
	Spec ShellSpec `yaml:"spec"`
}

func (s Shell) Run(ctx context.Context, variables map[string]string) errors.Error {
	return nil
}
