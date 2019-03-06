package runner

import (
	"context"

	"github.com/puppetlabs/nebula/pkg/errors"
)

type ActionRunner interface {
	Run(ctx context.Context, variables map[string]string) errors.Error
	Decoder() Decoder
}

type Decoder interface {
	Decode(b []byte) errors.Error
}

func NewRunner(kind RunnerKind) (ActionRunner, error) {
	switch kind {
	case RunnerKindGKEClusterProvisioner:
		return &GKEClusterProvisioner{}, nil
	case RunnerKindShell:
		return &Shell{}, nil
	case RunnerKindWorkflow:
		return &Workflow{}, nil
	}

	return nil, errors.NewWorkflowRunnerNotFound(string(kind))
}
