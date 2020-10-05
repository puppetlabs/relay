package cmd

import (
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

	cmd.AddCommand(newStartClusterCommand())
	cmd.AddCommand(newStopClusterCommand())
	cmd.AddCommand(newDeleteClusterCommand())

	return cmd
}

func newStartClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the local cluster that can execute workflows",
		RunE:  doStartCluster,
	}

	cmd.Flags().IntP("load-balancer-port", "", cluster.DefaultLoadBalancerHostPort, "The port to map from the host to the service load balancer")
	cmd.Flags().StringP("image-registry-name", "", cluster.DefaultRegistryName, "The name to use on the host and on the cluster nodes for the container image registry")
	cmd.Flags().IntP("image-registry-port", "", cluster.DefaultRegistryPort, "The port to use on the host and on the cluster nodes for the container image registry")

	return cmd
}

func doStartCluster(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	lbHostPort, err := cmd.Flags().GetInt("load-balancer-port")
	if err != nil {
		return err
	}

	registryName, err := cmd.Flags().GetString("image-registry-name")
	if err != nil {
		return err
	}

	registryPort, err := cmd.Flags().GetInt("image-registry-port")
	if err != nil {
		return err
	}

	cm := cluster.NewManager(ClusterConfig)

	if _, err := cm.Exists(ctx); err != nil {
		Dialog.Info("Creating a new dev cluster")
		opts := cluster.CreateOptions{
			LoadBalancerHostPort: lbHostPort,
			ImageRegistryName:    registryName,
			ImageRegistryPort:    registryPort,
		}
		if err := cm.Create(ctx, opts); err != nil {
			return err
		}

		cl, err := cm.GetClient(ctx, cluster.ClientOptions{Scheme: dev.DefaultScheme})
		if err != nil {
			return err
		}

		dm := dev.NewManager(cm, cl, DevConfig)

		Dialog.Info("Writing kubeconfig")
		if err := dm.WriteKubeconfig(ctx); err != nil {
			return err
		}

		Dialog.Info("Initializing relay-core; this might take a couple minutes...")
		initOpts := dev.InitializeOptions{
			ImageRegistryPort: registryPort,
		}
		if err := dm.InitializeRelayCore(ctx, initOpts); err != nil {
			return err
		}

		Dialog.Infof("Cluster connection can be used with: !Connection {type: kubernetes, name: %s}", dev.RelayClusterConnectionName)
	} else {
		if err := cm.Start(ctx); err != nil {
			return err
		}

		cl, err := cm.GetClient(ctx, cluster.ClientOptions{Scheme: dev.DefaultScheme})
		if err != nil {
			return err
		}

		// dev manager depends on a kubernetes client, which we can't get if a
		// cluster doesn't exist, so we can't create one at the top of this
		// function.
		dm := dev.NewManager(cm, cl, DevConfig)

		if err := dm.StartRelayCore(ctx); err != nil {
			return err
		}
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

	Dialog.Progress("Deleting cluster; this may take a minute...")

	cl, err := cm.GetClient(ctx, cluster.ClientOptions{Scheme: dev.DefaultScheme})
	if err != nil {
		return err
	}

	dm := dev.NewManager(cm, cl, DevConfig)

	return dm.Delete(ctx)
}
