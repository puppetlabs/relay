package cluster

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	k3dcluster "github.com/rancher/k3d/v3/pkg/cluster"
	"github.com/rancher/k3d/v3/pkg/runtimes"
	"github.com/rancher/k3d/v3/pkg/types"
)

const (
	DefaultClusterName = "relay-workflows"
	DefaultNetworkName = "relay-workflows-net"
	DefaultWorkerCount = 3
	DefaultK3sVersion  = "latest"
)

type ClusterOptions struct {
	DataDir string
}

func CreateCluster(ctx context.Context, opts ClusterOptions) error {
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
			IsInit:    true,
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
		InitNode:           serverNode,
		Nodes:              nodes,
		CreateClusterOpts:  &types.ClusterCreateOpts{},
		Network:            network,
		ExposeAPI:          exposeAPI,
	}

	if err := k3dcluster.ClusterCreate(ctx, rt, clusterConfig); err != nil {
		return err
	}

	c, err := GetCluster(ctx)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(opts.DataDir, 0700); err != nil {
		return err
	}

	apiconfig, err := k3dcluster.KubeconfigGet(ctx, rt, c)
	if err != nil {
		return err
	}

	err = k3dcluster.KubeconfigWriteToPath(ctx, apiconfig, filepath.Join(opts.DataDir, "kubeconfig"))
	if err != nil {
		return err
	}

	return nil
}

func GetCluster(ctx context.Context) (*types.Cluster, error) {
	rt := runtimes.SelectedRuntime
	clusterConfig := &types.Cluster{
		Name: DefaultClusterName,
	}

	return k3dcluster.ClusterGet(ctx, rt, clusterConfig)
}

func StartCluster(ctx context.Context) error {
	rt := runtimes.SelectedRuntime
	clusterConfig := &types.Cluster{
		Name: DefaultClusterName,
	}

	clusterConfig, err := GetCluster(ctx)
	if err != nil {
		return err
	}

	return k3dcluster.ClusterStart(ctx, rt, clusterConfig, types.ClusterStartOpts{})
}

func StopCluster(ctx context.Context) error {
	rt := runtimes.SelectedRuntime
	clusterConfig := &types.Cluster{
		Name: DefaultClusterName,
	}

	clusterConfig, err := GetCluster(ctx)
	if err != nil {
		return err
	}

	return k3dcluster.ClusterStop(ctx, rt, clusterConfig)
}
