package cluster

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/puppetlabs/relay/pkg/dialog"
	k3dcluster "github.com/rancher/k3d/v3/pkg/cluster"
	"github.com/rancher/k3d/v3/pkg/runtimes"
	"github.com/rancher/k3d/v3/pkg/tools"
	"github.com/rancher/k3d/v3/pkg/types"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

const (
	k3sVersion          = "v1.18.6-k3s1"
	k3sLocalStoragePath = "/var/lib/rancher/k3s/storage"
)

var agentArgs = []string{
	"--node-label=nebula.puppet.com/scheduling.customer-ready=true",
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
func (m *K3dClusterManager) Create(ctx context.Context) error {
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

	exposeAPI := types.ExposeAPI{
		Host:   types.DefaultAPIHost,
		HostIP: types.DefaultAPIHost,
		Port:   types.DefaultAPIPort,
	}

	serverNode := &types.Node{
		Role:  types.ServerRole,
		Image: k3sImage,
		ServerOpts: types.ServerOpts{
			ExposeAPI: exposeAPI,
		},
		Volumes: volumes,
	}

	nodes := []*types.Node{
		serverNode,
	}

	for i := 0; i < WorkerCount; i++ {
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

	clusterConfig := &types.Cluster{
		Name:               ClusterName,
		ServerLoadBalancer: &types.Node{Role: types.LoadBalancerRole},
		Nodes:              nodes,
		CreateClusterOpts: &types.ClusterCreateOpts{
			WaitForServer: true,
		},
		Network:   network,
		ExposeAPI: exposeAPI,
	}

	if err := k3dcluster.ClusterCreate(ctx, m.runtime, clusterConfig); err != nil {
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
	if err := k3dcluster.ClusterStart(ctx, m.runtime, clusterConfig, opts); err != nil {
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

	return k3dcluster.ClusterStop(ctx, m.runtime, clusterConfig)
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

	return k3dcluster.ClusterDelete(ctx, m.runtime, clusterConfig)
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

	if err := tools.ImageImportIntoClusterMulti(ctx, m.runtime, images, clusterConfig, types.ImageImportOpts{}); err != nil {
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

	return k3dcluster.KubeconfigGet(ctx, m.runtime, clusterConfig)
}

// WriteKubeconfig takes a path and writes the cluster's kubeconfig file to it. Attempting
// to write a kubeconfig for a cluster that doesn't exist results in an error.
func (m *K3dClusterManager) WriteKubeconfig(ctx context.Context, path string) error {
	apiConfig, err := m.GetKubeconfig(ctx)
	if err != nil {
		return err
	}

	return k3dcluster.KubeconfigWriteToPath(ctx, apiConfig, path)
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

	return k3dcluster.ClusterGet(ctx, m.runtime, clusterConfig)
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
