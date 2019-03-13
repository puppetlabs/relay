package runner

import (
	"context"

	"github.com/puppetlabs/nebula/pkg/errors"
	"github.com/puppetlabs/nebula/pkg/infra/provider/gcp"
	"github.com/puppetlabs/nebula/pkg/infra/provider/kubernetes/helm"
	"gopkg.in/yaml.v2"
)

type HelmDeploySpec struct {
	ResourceID string `yaml:"resourceID"`
	ChartPath  string `yaml:"chartPath"`
	ValuesPath string `yaml:"valuesPath"`
	Release    string `yaml:"release"`
}

type HelmDeploy struct {
	Name string         `yaml:"name"`
	Spec HelmDeploySpec `yaml:"spec"`
}

func (h *HelmDeploy) Run(ctx context.Context, rid string, r ActionRuntime, variables map[string]string) errors.Error {
	cluster, err := gcp.NewClusterFromResourceID(h.Spec.ResourceID, r.StateManager(), r.Logger())
	if err != nil {
		return err
	}

	if err := cluster.LookupRemote(ctx); err != nil {
		return err
	}

	hm := helm.NewHelmManager(cluster.KubeconfigPath(), r.Logger())

	chart := &helm.Chart{
		Path:        h.Spec.ChartPath,
		ValuesPaths: []string{h.Spec.ValuesPath},
	}

	if err := hm.DeployChart(ctx, h.Spec.Release, chart); err != nil {
		return err
	}

	return nil
}

func (h *HelmDeploy) Decoder() Decoder {
	return &HelmDeployDecoder{h: h}
}

type HelmDeployDecoder struct {
	h *HelmDeploy
}

func (d *HelmDeployDecoder) Decode(b []byte) errors.Error {
	if err := yaml.Unmarshal(b, d.h); err != nil {
		return errors.NewWorkflowRunnerDecodeError().WithCause(err).Bug()
	}

	return nil
}
