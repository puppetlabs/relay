package runner

import (
	"context"

	"github.com/puppetlabs/nebula/pkg/errors"
	"gopkg.in/yaml.v2"
)

type GKEClusterProvisionerSpec struct {
	Name      string `yaml:"name"`
	ProjectID string `yaml:"projectID"`
	Region    string `yaml:"region"`
}

type GKEClusterProvisioner struct {
	Name string                     `yaml:"name"`
	Spec *GKEClusterProvisionerSpec `yaml:"spec"`
}

func (g *GKEClusterProvisioner) Run(ctx context.Context, variables map[string]string) errors.Error {
	return nil
}

func (g *GKEClusterProvisioner) Decoder() Decoder {
	return &GKEClusterProvisionerDecoder{gcp: g}
}

type GKEClusterProvisionerDecoder struct {
	gcp *GKEClusterProvisioner
}

func (d *GKEClusterProvisionerDecoder) Decode(b []byte) errors.Error {
	if err := yaml.Unmarshal(b, d.gcp); err != nil {
		return errors.NewWorkflowRunnerDecodeError().WithCause(err).Bug()
	}

	return nil
}
