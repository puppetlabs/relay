package cluster

import (
	"context"
	"fmt"

	k3dcluster "github.com/rancher/k3d/v3/pkg/cluster"
	"github.com/rancher/k3d/v3/pkg/runtimes"
	"github.com/rancher/k3d/v3/pkg/types"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

const (
	K3sVersion = "v1.18.6-k3s1"
)

type Client struct {
	APIClient client.Client
	Mapper    meta.RESTMapper
}

// K3dClusterManager wraps rancher's k3d to manage the lifecycle
// of a kubernetes cluster running in docker.
type K3dClusterManager struct{}

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
	rt := runtimes.SelectedRuntime
	k3sImage := fmt.Sprintf("%s:%s", types.DefaultK3sImageRepo, K3sVersion)

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
		Labels: map[string]string{
			"nebula.puppet.com/scheduling.customer-ready": "true",
		},
	}

	nodes := []*types.Node{
		serverNode,
	}

	for i := 0; i < WorkerCount; i++ {
		node := &types.Node{
			Role:  types.AgentRole,
			Image: k3sImage,
			Labels: map[string]string{
				"nebula.puppet.com/scheduling.customer-ready": "true",
			},
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

	if err := k3dcluster.ClusterCreate(ctx, rt, clusterConfig); err != nil {
		return err
	}

	return nil
}

// Start starts the cluster. Attempting to start a cluster that doesn't exist
// results in an error.
//
// Note: There is currently a bug in k3d that causes ClusterStart to hang
// while waiting for the serverlb node if the cluster is already started.
// I filed a ticker here: https://github.com/rancher/k3d/issues/310
// In order to make the `relay dev cluster start` command more idempotent, this
// bug will need to be fixed or worked around.
func (m *K3dClusterManager) Start(ctx context.Context) error {
	rt := runtimes.SelectedRuntime
	clusterConfig := &types.Cluster{
		Name: ClusterName,
	}

	clusterConfig, err := m.get(ctx)
	if err != nil {
		return err
	}

	return k3dcluster.ClusterStart(ctx, rt, clusterConfig, types.ClusterStartOpts{})
}

// Stop stops the cluster. Attempting to stop a cluster that doesn't exist
// results in an error.
func (m *K3dClusterManager) Stop(ctx context.Context) error {
	rt := runtimes.SelectedRuntime
	clusterConfig := &types.Cluster{
		Name: ClusterName,
	}

	clusterConfig, err := m.get(ctx)
	if err != nil {
		return err
	}

	return k3dcluster.ClusterStop(ctx, rt, clusterConfig)
}

// Delete deletes the cluster and all its resources (docker network and volumes included).
// Attempting to delete a cluster that doesn't exist results in an error.
func (m *K3dClusterManager) Delete(ctx context.Context) error {
	rt := runtimes.SelectedRuntime
	clusterConfig := &types.Cluster{
		Name: ClusterName,
	}

	clusterConfig, err := m.get(ctx)
	if err != nil {
		return err
	}

	return k3dcluster.ClusterDelete(ctx, rt, clusterConfig)
}

// GetKubeconfig returns a k8s client-go config for the cluster's context. This can be
// be used to generate the yaml that is often seen on disk and used with kubectl. Attempting
// to get a kubeconfig for a cluster that doesn't exist results in an error.
func (m *K3dClusterManager) GetKubeconfig(ctx context.Context) (*clientcmdapi.Config, error) {
	rt := runtimes.SelectedRuntime

	clusterConfig, err := m.get(ctx)
	if err != nil {
		return nil, err
	}

	return k3dcluster.KubeconfigGet(ctx, rt, clusterConfig)
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
	rt := runtimes.SelectedRuntime
	clusterConfig := &types.Cluster{
		Name: ClusterName,
	}

	return k3dcluster.ClusterGet(ctx, rt, clusterConfig)
}

// NewK3dClusterManager returns a new K3dClusterManager.
func NewK3dClusterManager() *K3dClusterManager {
	return &K3dClusterManager{}
}
