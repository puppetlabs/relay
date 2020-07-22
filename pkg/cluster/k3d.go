package cluster

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/puppetlabs/relay/pkg/config"
	k3dcluster "github.com/rancher/k3d/v3/pkg/cluster"
	"github.com/rancher/k3d/v3/pkg/runtimes"
	"github.com/rancher/k3d/v3/pkg/types"
)

type Options struct {
	DataDir string
}

type K3dClusterManager struct {
	opts Options
}

func (m *K3dClusterManager) Exists(ctx context.Context) (bool, error) {
	if _, err := m.get(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (m *K3dClusterManager) Create(ctx context.Context) error {
	rt := runtimes.SelectedRuntime
	k3sImage := fmt.Sprintf("%s:%s", types.DefaultK3sImageRepo, DefaultK3sVersion)

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
	}

	nodes := []*types.Node{
		serverNode,
	}

	for i := 0; i < DefaultWorkerCount; i++ {
		node := &types.Node{
			Role:  types.AgentRole,
			Image: k3sImage,
		}

		nodes = append(nodes, node)
	}

	network := types.ClusterNetwork{
		Name: DefaultNetworkName,
	}

	clusterConfig := &types.Cluster{
		Name:               DefaultClusterName,
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

	c, err := m.get(ctx)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(m.opts.DataDir, 0700); err != nil {
		return err
	}

	apiconfig, err := k3dcluster.KubeconfigGet(ctx, rt, c)
	if err != nil {
		return err
	}

	err = k3dcluster.KubeconfigWriteToPath(ctx, apiconfig, filepath.Join(m.opts.DataDir, "kubeconfig"))
	if err != nil {
		return err
	}

	return nil
}

func (m *K3dClusterManager) Start(ctx context.Context) error {
	rt := runtimes.SelectedRuntime
	clusterConfig := &types.Cluster{
		Name: DefaultClusterName,
	}

	clusterConfig, err := m.get(ctx)
	if err != nil {
		return err
	}

	return k3dcluster.ClusterStart(ctx, rt, clusterConfig, types.ClusterStartOpts{})
}

func (m *K3dClusterManager) Stop(ctx context.Context) error {
	rt := runtimes.SelectedRuntime
	clusterConfig := &types.Cluster{
		Name: DefaultClusterName,
	}

	clusterConfig, err := m.get(ctx)
	if err != nil {
		return err
	}

	return k3dcluster.ClusterStop(ctx, rt, clusterConfig)
}

func (m *K3dClusterManager) Delete(ctx context.Context) error {
	rt := runtimes.SelectedRuntime
	clusterConfig := &types.Cluster{
		Name: DefaultClusterName,
	}

	clusterConfig, err := m.get(ctx)
	if err != nil {
		return err
	}

	return k3dcluster.ClusterDelete(ctx, rt, clusterConfig)
}

func (m *K3dClusterManager) get(ctx context.Context) (*types.Cluster, error) {
	rt := runtimes.SelectedRuntime
	clusterConfig := &types.Cluster{
		Name: DefaultClusterName,
	}

	return k3dcluster.ClusterGet(ctx, rt, clusterConfig)
}

func NewK3dClusterManager(cfg *config.Config) *K3dClusterManager {
	opts := Options{
		DataDir: filepath.Join(cfg.DataDir, "cluster"),
	}

	return &K3dClusterManager{
		opts: opts,
	}
}
