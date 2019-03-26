package executor

import (
	"context"

	"github.com/puppetlabs/nebula/pkg/errors"
	"github.com/puppetlabs/nebula/pkg/execution"
	"github.com/puppetlabs/nebula/pkg/execution/docker"
	"github.com/puppetlabs/nebula/pkg/plan/types"
)

type ActionExecutor interface {
	Kind() string
	ScheduleAction(context.Context, execution.ExecutorRuntime, *types.Action, map[string]string) errors.Error
}

type RegistryCredentials struct {
	Registry string
	User     string
	Pass     string
}

// TODO just returns docker for now and the registry creds are passed in weird
func NewExecutor(creds RegistryCredentials) (ActionExecutor, error) {
	return docker.NewExecutor(docker.ExecutorOptions{
		Registry:     creds.Registry,
		RegistryUser: creds.User,
		RegistryPass: creds.Pass,
	})
}
