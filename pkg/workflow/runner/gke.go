package runner

import (
	"context"

	"github.com/puppetlabs/nebula/pkg/errors"
	"github.com/puppetlabs/nebula/pkg/infra/provider/gcp"
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

func (g *GKEClusterProvisioner) Run(ctx context.Context, rid string, r ActionRuntime, variables map[string]string) errors.Error {
	c, err := gcp.NewCluster(rid, r.StateManager(), gcp.ClusterSpec{
		Name:        g.Spec.Name,
		Description: g.Spec.Description,
		Nodes:       int32(g.Spec.Nodes),
		Region:      g.Spec.Region,
		ProjectID:   g.Spec.ProjectID,
	})
	if err != nil {
		return err
	}

	if err := c.LookupRemote(ctx); err != nil {
		return err
	}
	r.Logger().Info("cluster-state-fetched", "name", g.Spec.Name, "status", c.Status)

	if err := c.Sync(ctx); err != nil {
		return err
	}
	r.Logger().Info("cluster-state-synced", "name", g.Spec.Name, "status", c.Status)

	if err := c.SaveState(ctx); err != nil {
		return err
	}
	r.Logger().Info("cluster-resource-saved", "name", g.Spec.Name, "resource-id", rid)

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
