package cmd

import (
	"context"

	"github.com/puppetlabs/relay/pkg/cluster"
	"github.com/puppetlabs/relay/pkg/dev"
	"github.com/spf13/cobra"
)

func newClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Manage the local workflow execution cluster",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(newCreateClusterCommand())
	cmd.AddCommand(newStartClusterCommand())
	cmd.AddCommand(newStopClusterCommand())
	cmd.AddCommand(newDeleteClusterCommand())

	return cmd
}

func newCreateClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create the local cluster",
		RunE:  doCreateCluster,
	}

	cmd.Flags().IntP("load-balancer-port", "", cluster.DefaultLoadBalancerHostPort, "The port to map from the host to the service load balancer")
	cmd.Flags().IntP("worker-count", "", cluster.DefaultWorkerCount, "The number of worker nodes to create on the cluster")

	return cmd
}

func doCreateCluster(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	lbHostPort, err := cmd.Flags().GetInt("load-balancer-port")
	if err != nil {
		return err
	}

	workerCount, err := cmd.Flags().GetInt("worker-count")
	if err != nil {
		return err
	}

	opts := cluster.CreateOptions{
		LoadBalancerHostPort: lbHostPort,
		WorkerCount:          workerCount,
	}

	return createCluster(ctx, opts)
}

func createCluster(ctx context.Context, opts cluster.CreateOptions) error {
	cm := cluster.NewManager(ClusterConfig)

	if exists, _ := cm.Exists(ctx); !exists {
		Dialog.Info("Creating a new cluster")
		if err := cm.Create(ctx, opts); err != nil {
			return err
		}

		Dialog.Info("Writing kubeconfig")
		if err := cm.WriteKubeconfig(ctx); err != nil {
			return err
		}

		Dialog.Info("Cluster is ready to use")
	} else {
		Dialog.Info("Cluster already exists")
	}

	return nil
}

func newStartClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the local cluster",
		RunE:  doStartCluster,
	}

	return cmd
}

func doStartCluster(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	return startCluster(ctx)
}

func startCluster(ctx context.Context) error {
	cm := cluster.NewManager(ClusterConfig)

	if exists, _ := cm.Exists(ctx); !exists {
		Dialog.Info("Cluster does not exist")
		return nil
	}

	if err := cm.Start(ctx); err != nil {
		return err
	}

	_, err := dev.NewManagerFromLocalCluster(ctx, cm, DevConfig)
	if err != nil {
		return err
	}

	Dialog.Info("Relay dev cluster is ready to use")

	return nil
}

func newStopClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the local cluster",
		RunE:  doStopCluster,
	}

	return cmd
}

func doStopCluster(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cm := cluster.NewManager(ClusterConfig)

	if exists, _ := cm.Exists(ctx); !exists {
		Dialog.Info("Cluster does not exist")
		return nil
	}

	return cm.Stop(ctx)
}

func newDeleteClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete the local cluster",
		RunE:  doDeleteCluster,
	}

	return cmd
}

func doDeleteCluster(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cm := cluster.NewManager(ClusterConfig)

	if exists, _ := cm.Exists(ctx); !exists {
		Dialog.Info("Cluster does not exist")
		return nil
	}

	Dialog.Progress("Deleting cluster; this may take several minutes...")

	dm, err := dev.NewManagerFromLocalCluster(ctx, cm, DevConfig)
	if err != nil {
		return err
	}

	return dm.Delete(ctx)
}
