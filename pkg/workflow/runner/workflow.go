package runner

import (
	"context"
	"net/url"

	"github.com/puppetlabs/nebula/pkg/errors"
)

type WorkflowSpec struct {
	Import *url.URL `yaml:"import"`
}

type Workflow struct {
	Name string       `yaml:"name"`
	Spec WorkflowSpec `yaml:"spec"`
}

func (w Workflow) Run(ctx context.Context, variables map[string]string) errors.Error {
	return nil
}
