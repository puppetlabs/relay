package runner

import (
	"context"

	container "cloud.google.com/go/container/apiv1"
	"github.com/puppetlabs/nebula/pkg/errors"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"gopkg.in/yaml.v2"
)

type GKEClusterProvisionerSpec struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	ProjectID   string `yaml:"projectID"`
	Region      string `yaml:"region"`
	Nodes       int    `yaml:"nodes"`
}

type GKEClusterProvisioner struct {
	Name string                     `yaml:"name"`
	Spec *GKEClusterProvisionerSpec `yaml:"spec"`
}

func (g *GKEClusterProvisioner) Run(ctx context.Context, r ActionRuntime, variables map[string]string) errors.Error {
	manager, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		return errors.NewWorkflowUnknownRuntimeError().WithCause(err).Bug()
	}

	req := &containerpb.CreateClusterRequest{
		ProjectId: g.Spec.ProjectID,
		Cluster: &containerpb.Cluster{
			Name:             g.Spec.Name,
			InitialNodeCount: int32(g.Spec.Nodes),
			Description:      g.Spec.Description,
		},
	}

	resp, err := manager.CreateCluster(ctx, req)
	if err != nil {
		return errors.NewWorkflowUnknownRuntimeError().WithCause(err).Bug()
	}

	r.Logger().Info("created cluster", "status", resp.Status)

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
