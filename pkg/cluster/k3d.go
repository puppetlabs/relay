package cluster

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/docker/go-connections/nat"
	"github.com/puppetlabs/relay/pkg/dialog"
	k3dclient "github.com/rancher/k3d/v4/pkg/client"
	"github.com/rancher/k3d/v4/pkg/config/v1alpha2"
	"github.com/rancher/k3d/v4/pkg/runtimes"
	"github.com/rancher/k3d/v4/pkg/tools"
	"github.com/rancher/k3d/v4/pkg/types"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

const (
	k3sVersion          = "v1.20.6-k3s1"
	k3sLocalStoragePath = "/var/lib/rancher/k3s/storage"
)

var agentArgs = []string{
	"--node-label=nebula.puppet.com/scheduling.customer-ready=true",
}

// Mirror https://github.com/rancher/k3s/blob/master/pkg/agent/templates/registry.go
type Mirror struct {
	Endpoints []string `toml:"endpoint" yaml:"endpoint"`
}

// Registry https://github.com/rancher/k3s/blob/master/pkg/agent/templates/registry.go
type Registry struct {
	Mirrors map[string]Mirror `toml:"mirrors" yaml:"mirrors"`
}

type Client struct {
	APIClient client.Client
	Mapper    meta.RESTMapper
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
	volumes := []string{
		localStorage,
	}

	// If /dev/mapper exists, we'll automatically map it into the cluster
	// controller.
	if _, err := os.Stat("/dev/mapper"); !os.IsNotExist(err) {
		volumes = append(volumes, "/dev/mapper:/dev/mapper:ro")
	}

	// TODO Temporary workaround to ensure image pulls and manifest lookups can function conjointly against the in-cluster registry
	registryConfigPath := path.Join(m.cfg.WorkDir.Path, "registries.yaml")
	volumes = append(volumes, registryConfigPath+":/etc/rancher/k3s/registries.yaml")

	registry := &Registry{
		Mirrors: map[string]Mirror{
			"docker.io": Mirror{
				Endpoints: []string{ImagePassthroughCacheAddr},
			},
			fmt.Sprintf("%s:%d", opts.ImageRegistryName, opts.ImageRegistryPort): Mirror{
				Endpoints: []string{fmt.Sprintf("http://localhost:%d", opts.ImageRegistryPort)},
			},
		},
	}
	m.storeRegistryConfiguration(registryConfigPath, registry)

	exposeAPI := &types.ExposureOpts{
		Host: types.DefaultAPIHost,
		PortMapping: nat.PortMapping{
			Port: types.DefaultAPIPort,
			Binding: nat.PortBinding{
				HostIP:   types.DefaultAPIHost,
				HostPort: types.DefaultAPIPort,
			},
		},
	}

	imageRegistryPort, err := nat.NewPort("tcp", fmt.Sprintf("%d", opts.ImageRegistryPort))
	if err != nil {
		return fmt.Errorf("failed to create cluster: %w", err)
	}
	registryPortMapping := nat.PortMap{}
	registryPortMapping[imageRegistryPort] = []nat.PortBinding{{
		HostIP:   types.DefaultAPIHost,
		HostPort: imageRegistryPort.Port(),
	}}

	serverNode := &types.Node{
		Role:  types.ServerRole,
		Image: k3sImage,
		ServerOpts: types.ServerOpts{
			KubeAPI: exposeAPI,
		},
		Volumes: volumes,
		Ports:   registryPortMapping,
	}

	if opts.WorkerCount <= 0 {
		serverNode.Args = agentArgs
	}

	nodes := []*types.Node{
		serverNode,
	}

	for i := 0; i < opts.WorkerCount; i++ {
		node := &types.Node{
			Role:    types.AgentRole,
			Image:   k3sImage,
			Args:    agentArgs,
			Volumes: volumes,
		}

		nodes = append(nodes, node)
	}

	network := types.ClusterNetwork{
		Name: NetworkName,
	}

	lbHostPort := DefaultLoadBalancerHostPort
	if opts.LoadBalancerHostPort != 0 {
		lbHostPort = opts.LoadBalancerHostPort
	}

	lbPort, err := nat.NewPort("tcp", fmt.Sprintf("%d", DefaultLoadBalancerHostPort))
	if err != nil {
		return fmt.Errorf("failed to create cluster: %w", err)
	}
	lbPortMapping := nat.PortMap{}
	lbPortMapping[lbPort] = []nat.PortBinding{{
		HostIP:   types.DefaultAPIHost,
		HostPort: fmt.Sprintf("%d", lbHostPort),
	}}

	clusterConfig := &v1alpha2.ClusterConfig{
		ClusterCreateOpts: types.ClusterCreateOpts{
			PrepDisableHostIPInjection: true,
			WaitForServer:              true,
			// HACK this is to workaround a k3d bug
			GlobalLabels: make(map[string]string),
		},
		Cluster: types.Cluster{
			Name: ClusterName,
			ServerLoadBalancer: &types.Node{
				Role:  types.LoadBalancerRole,
				Ports: lbPortMapping,
			},
			Nodes:   nodes,
			Network: network,
			KubeAPI: exposeAPI,
		},
	}

	if err := k3dclient.ClusterRun(ctx, m.runtime, clusterConfig); err != nil {
		return fmt.Errorf("failed to create cluster: %w", err)
	}

	return nil
}

// Start starts the cluster. Attempting to start a cluster that doesn't exist
// results in an error.
func (m *K3dClusterManager) Start(ctx context.Context) error {
	clusterConfig := &types.Cluster{
		Name: ClusterName,
	}

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
	clusterConfig := &types.Cluster{
		Name: ClusterName,
	}

	clusterConfig, err := m.get(ctx)
	if err != nil {
		return err
	}

	return k3dclient.ClusterStop(ctx, m.runtime, clusterConfig)
}

// Delete deletes the cluster and all its resources (docker network and volumes included).
// Attempting to delete a cluster that doesn't exist results in an error.
func (m *K3dClusterManager) Delete(ctx context.Context) error {
	clusterConfig := &types.Cluster{
		Name: ClusterName,
	}

	clusterConfig, err := m.get(ctx)
	if err != nil {
		return err
	}

	return k3dclient.ClusterDelete(ctx, m.runtime, clusterConfig, types.ClusterDeleteOpts{})
}

// ImportImages import's a given image or images from the local container
// runtime into every node in the cluster. It's useful for making sure custom
// built images on the local host machine take priority over remote images in a
// registry.
func (m *K3dClusterManager) ImportImages(ctx context.Context, images ...string) error {
	clusterConfig, err := m.get(ctx)
	if err != nil {
		return fmt.Errorf("failed to lookup cluster: %w", err)
	}

	if err := tools.ImageImportIntoClusterMulti(ctx, m.runtime, images, clusterConfig, types.ImageImportOpts{KeepTar: true}); err != nil {
		return fmt.Errorf("failed to import image: %w", err)
	}

	return nil
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

// GetClient returns a new Client configured with a RESTMapper and k8s api client.
func (m *K3dClusterManager) GetClient(ctx context.Context, opts ClientOptions) (*Client, error) {
	apiConfig, err := m.GetKubeconfig(ctx)
	if err != nil {
		return nil, err
	}

	overrides := &clientcmd.ConfigOverrides{}
	clientConfig := clientcmd.NewDefaultClientConfig(*apiConfig, overrides)

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	c, err := client.New(restConfig, client.Options{
		Scheme: opts.Scheme,
	})
	if err != nil {
		return nil, err
	}

	mapper, err := apiutil.NewDynamicRESTMapper(restConfig)
	if err != nil {
		return nil, err
	}

	return &Client{
		APIClient: c,
		Mapper:    mapper,
	}, nil
}

func (m *K3dClusterManager) get(ctx context.Context) (*types.Cluster, error) {
	clusterConfig := &types.Cluster{
		Name: ClusterName,
	}

	return k3dclient.ClusterGet(ctx, m.runtime, clusterConfig)
}

func (m *K3dClusterManager) storeRegistryConfiguration(path string, registry *Registry) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0750)
	if err != nil {
		return err
	}

	defer f.Close()

	data, err := yaml.Marshal(registry)

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

// NewK3dClusterManager returns a new K3dClusterManager.
func NewK3dClusterManager(cfg Config) *K3dClusterManager {
	log.SetFormatter(&logrusToDialogFormatter{
		dialog: cfg.Dialog,
	})

	return &K3dClusterManager{
		runtime: runtimes.SelectedRuntime,
		cfg:     cfg,
	}
}

type logrusToDialogFormatter struct {
	dialog dialog.Dialog
}

func (f *logrusToDialogFormatter) Format(e *log.Entry) ([]byte, error) {
	switch e.Level {
	case log.DebugLevel, log.InfoLevel, log.WarnLevel:
		f.dialog.Info(e.Message)
	case log.ErrorLevel:
		f.dialog.Error(e.Message)
	}

	return nil, nil
}
