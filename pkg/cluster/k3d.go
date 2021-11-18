package cluster

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	k3dclient "github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/config"
	"github.com/rancher/k3d/v5/pkg/config/v1alpha3"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	"github.com/rancher/k3d/v5/pkg/types"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	k3sVersion          = "v1.21.6-k3s1"
	k3sLocalStoragePath = "/var/lib/rancher/k3s/storage"
)

var agentArgs = []string{
	"--node-label=nebula.puppet.com/scheduling.customer-ready=true",
}

// K3dClusterManager wraps rancher's k3d to manage the lifecycle
// of a kubernetes cluster running in docker.
type K3dClusterManager struct {
	runtime runtimes.Runtime
	cfg     Config
}

// Exists checks and reports back if the cluster exists.
func (m *K3dClusterManager) Exists(ctx context.Context) (bool, error) {
	if _, err := m.get(ctx); err != nil {
		return false, err
	}

	return true, nil
}

// Create uses opinionated configuration constants to create a kubernetes cluster
// running inside docker.
func (m *K3dClusterManager) Create(ctx context.Context, opts CreateOptions) error {
	k3sImage := fmt.Sprintf("%s:%s", types.DefaultK3sImageRepo, k3sVersion)

	hostStoragePath := filepath.Join(m.cfg.WorkDir.Path, HostStorageName)
	if err := os.MkdirAll(hostStoragePath, 0700); err != nil {
		return fmt.Errorf("failed to make the host storage directory: %w", err)
	}

	localStorage := fmt.Sprintf("%s:%s",
		hostStoragePath,
		k3sLocalStoragePath)

	volumes := []v1alpha3.VolumeWithNodeFilters{
		{
			Volume: localStorage,
		},
	}

	// If /dev/mapper exists, we'll automatically map it into the cluster
	// controller.
	if _, err := os.Stat("/dev/mapper"); !os.IsNotExist(err) {
		volumes = append(volumes, v1alpha3.VolumeWithNodeFilters{
			Volume: "/dev/mapper:/dev/mapper:ro",
		})
	}

	cfg := v1alpha3.SimpleConfig{
		Name:    ClusterName,
		Network: NetworkName,
		Image:   k3sImage,
		Servers: 1,
		Agents:  opts.WorkerCount,
		ExposeAPI: v1alpha3.SimpleExposureOpts{
			Host:     types.DefaultAPIHost,
			HostIP:   types.DefaultAPIHost,
			HostPort: types.DefaultAPIPort,
		},
		Options: v1alpha3.SimpleConfigOptions{
			K3dOptions: v1alpha3.SimpleConfigOptionsK3d{
				Wait: true,
			},
			KubeconfigOptions: v1alpha3.SimpleConfigOptionsKubeconfig{
				UpdateDefaultKubeconfig: true,
				SwitchCurrentContext:    true,
			},
		},
		Ports: []v1alpha3.PortWithNodeFilters{
			{
				Port: fmt.Sprintf("%d:%d", DefaultLoadBalancerHostPort, DefaultLoadBalancerNodePort),
			},
		},
		Volumes: volumes,
	}

	return m.clusterCreate(ctx, cfg)
}

// Start starts the cluster. Attempting to start a cluster that doesn't exist
// results in an error.
func (m *K3dClusterManager) Start(ctx context.Context) error {
	clusterConfig, err := m.get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cluster config: %w", err)
	}

	opts := types.ClusterStartOpts{WaitForServer: true}
	if err := k3dclient.ClusterStart(ctx, m.runtime, clusterConfig, opts); err != nil {
		return fmt.Errorf("failed to start cluster: %w", err)
	}

	return nil
}

// Stop stops the cluster. Attempting to stop a cluster that doesn't exist
// results in an error.
func (m *K3dClusterManager) Stop(ctx context.Context) error {
	clusterConfig, err := m.get(ctx)
	if err != nil {
		return err
	}

	return k3dclient.ClusterStop(ctx, m.runtime, clusterConfig)
}

// Delete deletes the cluster and all its resources (docker network and volumes included).
// Attempting to delete a cluster that doesn't exist results in an error.
func (m *K3dClusterManager) Delete(ctx context.Context) error {
	clusterConfig, err := m.get(ctx)
	if err != nil {
		return err
	}

	return k3dclient.ClusterDelete(ctx, m.runtime, clusterConfig, types.ClusterDeleteOpts{})
}

// GetKubeconfig returns a k8s client-go config for the cluster's context. This can be
// be used to generate the yaml that is often seen on disk and used with kubectl. Attempting
// to get a kubeconfig for a cluster that doesn't exist results in an error.
func (m *K3dClusterManager) GetKubeconfig(ctx context.Context) (*clientcmdapi.Config, error) {
	clusterConfig, err := m.get(ctx)
	if err != nil {
		return nil, err
	}

	return k3dclient.KubeconfigGet(ctx, m.runtime, clusterConfig)
}

// WriteKubeconfig takes a path and writes the cluster's kubeconfig file to it. Attempting
// to write a kubeconfig for a cluster that doesn't exist results in an error.
func (m *K3dClusterManager) WriteKubeconfig(ctx context.Context, path string) error {
	clusterConfig, err := m.get(ctx)
	if err != nil {
		return err
	}

	k3dclient.KubeconfigGetWrite(ctx, m.runtime, clusterConfig, "", &k3dclient.WriteKubeConfigOptions{
		OverwriteExisting:    false,
		UpdateCurrentContext: true,
		UpdateExisting:       true,
	})

	apiConfig, err := m.GetKubeconfig(ctx)
	if err != nil {
		return err
	}

	return k3dclient.KubeconfigWriteToPath(ctx, apiConfig, path)
}

func (m *K3dClusterManager) get(ctx context.Context) (*types.Cluster, error) {
	clusterConfig := &types.Cluster{
		Name: ClusterName,
	}

	return k3dclient.ClusterGet(ctx, m.runtime, clusterConfig)
}

// NewK3dClusterManager returns a new K3dClusterManager.
func NewK3dClusterManager(cfg Config) *K3dClusterManager {
	return &K3dClusterManager{
		runtime: runtimes.SelectedRuntime,
		cfg:     cfg,
	}
}

// Functionality based on the official k3d cluster create command
// https://github.com/rancher/k3d/blob/main/cmd/cluster/clusterCreate.go
func (m *K3dClusterManager) clusterCreate(ctx context.Context, cfg v1alpha3.SimpleConfig) error {
	clusterConfig, err := config.TransformSimpleToClusterConfig(ctx, runtimes.SelectedRuntime, cfg)
	if err != nil {
		return err
	}

	clusterConfig, err = config.ProcessClusterConfig(*clusterConfig)
	if err != nil {
		return err
	}

	if err := config.ValidateClusterConfig(ctx, runtimes.SelectedRuntime, *clusterConfig); err != nil {
		return err
	}

	if clusterConfig.KubeconfigOpts.UpdateDefaultKubeconfig {
		clusterConfig.ClusterCreateOpts.WaitForServer = true
	}

	if err := k3dclient.ClusterRun(ctx, m.runtime, clusterConfig); err != nil {
		return err
	}

	if !clusterConfig.KubeconfigOpts.UpdateDefaultKubeconfig && clusterConfig.KubeconfigOpts.SwitchCurrentContext {
		clusterConfig.KubeconfigOpts.SwitchCurrentContext = false
	}

	if clusterConfig.KubeconfigOpts.UpdateDefaultKubeconfig {
		if _, err := k3dclient.KubeconfigGetWrite(ctx, runtimes.SelectedRuntime, &clusterConfig.Cluster, "",
			&k3dclient.WriteKubeConfigOptions{
				UpdateExisting:       true,
				OverwriteExisting:    false,
				UpdateCurrentContext: cfg.Options.KubeconfigOptions.SwitchCurrentContext,
			}); err != nil {
			return err
		}
	}

	return nil
}
