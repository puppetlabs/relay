package workflow

import (
	"context"

	"github.com/puppetlabs/nebula/pkg/errors"
)

type ActionRunner interface {
	Run(ctx context.Context, variables map[string]string) errors.Error
}

type Action struct {
	Name string      `yaml:"name"`
	Kind string      `yaml:"kind"`
	Spec interface{} `yaml:"spec"`
}
