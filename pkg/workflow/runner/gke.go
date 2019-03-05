package runner

import (
	"context"

	"github.com/puppetlabs/nebula/pkg/errors"
)

type GKEClusterProvisionerSpec struct {
	Name      string `yaml:"name"`
	ProjectID string `yaml:"projectID"`
	Region    string `yaml:"region"`
}

type GKEClusterProvisioner struct {
	Name string                    `yaml:"name"`
	Spec GKEClusterProvisionerSpec `yaml:"spec"`
}

func (g GKEClusterProvisioner) Run(ctx context.Context, variables map[string]string) errors.Error {
	return nil
}
