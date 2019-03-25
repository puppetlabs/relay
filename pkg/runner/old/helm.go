package runner

import (
	"context"

	"github.com/puppetlabs/nebula/pkg/errors"
	"github.com/puppetlabs/nebula/pkg/infra/provider/gcp"
	"github.com/puppetlabs/nebula/pkg/infra/provider/kubernetes/helm"
	"gopkg.in/yaml.v2"
)

type HelmValueOverride struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

func overridesMap(overrides []HelmValueOverride) map[string]string {
	m := make(map[string]string)

	for _, o := range overrides {
		m[o.Key] = o.Value
	}

	return m
}

type HelmChartDeploymentSpec struct {
	ResourceID string              `yaml:"resourceID"`
	ChartPath  string              `yaml:"chartPath"`
	ValuesPath string              `yaml:"valuesPath"`
	Release    string              `yaml:"release"`
	Overrides  []HelmValueOverride `yaml:"overrides"`
}

type HelmChartDeployment struct {
	Name string                  `yaml:"name"`
	Spec HelmChartDeploymentSpec `yaml:"spec"`
}

func (h *HelmChartDeployment) Run(ctx context.Context, rid string, r ActionRuntime, variables map[string]string) errors.Error {
	cluster, err := gcp.NewClusterFromResourceID(h.Spec.ResourceID, r.StateManager(), r.Logger())
	if err != nil {
		return err
	}

	ready, err := cluster.LookupRemote(ctx)
	if err != nil {
		return err
	}

	if !ready {
		r.Logger().Warn("cluster-not-ready", "resource-id", h.Spec.ResourceID)
		return nil
	}

	hm := helm.NewHelmManager(cluster.KubeconfigPath(), r.Logger())

	var valuesPaths []string

	if h.Spec.ValuesPath != "" {
		valuesPaths = append(valuesPaths, h.Spec.ValuesPath)
	}

	chart := &helm.Chart{
		Path:           h.Spec.ChartPath,
		ValuesPaths:    valuesPaths,
		ValueOverrides: overridesMap(h.Spec.Overrides),
	}

	if err := hm.DeployChart(ctx, h.Spec.Release, chart, variables); err != nil {
		return err
	}

	return nil
}

func (h *HelmChartDeployment) Decoder() Decoder {
	return &HelmChartDeploymentDecoder{h: h}
}

type HelmChartDeploymentDecoder struct {
	h *HelmChartDeployment
}

func (d *HelmChartDeploymentDecoder) Decode(b []byte) errors.Error {
	if err := yaml.Unmarshal(b, d.h); err != nil {
		return errors.NewWorkflowRunnerDecodeError().WithCause(err).Bug()
	}

	return nil
}
