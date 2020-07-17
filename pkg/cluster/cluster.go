package cluster

import (
	"context"
	"fmt"

	k3dcluster "github.com/rancher/k3d/v3/pkg/cluster"
	"github.com/rancher/k3d/v3/pkg/runtimes"
	"github.com/rancher/k3d/v3/pkg/types"
)

const (
	DefaultClusterName = "relay-workflows"
	DefaultNetworkName = "relay-workflows-net"
	DefaultWorkerCount = 3
)

func CreateCluster(ctx context.Context) error {
	rt := runtimes.SelectedRuntime
	k3sImage := fmt.Sprintf("%s:%s", types.DefaultK3sImageRepo, "latest")

	serverNode := &types.Node{
		Role:  types.ServerRole,
		Image: k3sImage,
		ServerOpts: types.ServerOpts{
			IsInit: true,
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

	exposeAPI := types.ExposeAPI{
		Host:   types.DefaultAPIHost,
		HostIP: types.DefaultAPIHost,
		Port:   types.DefaultAPIPort,
	}

	serverNode.ServerOpts.ExposeAPI = exposeAPI

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
