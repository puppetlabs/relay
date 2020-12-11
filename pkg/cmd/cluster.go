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
	cmd.AddCommand(newInitClusterCommand())
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
	cmd.Flags().StringP("image-registry-name", "", cluster.DefaultRegistryName, "The name to use on the host and on the cluster nodes for the container image registry")
	cmd.Flags().IntP("image-registry-port", "", cluster.DefaultRegistryPort, "The port to use on the host and on the cluster nodes for the container image registry")
	cmd.Flags().IntP("worker-count", "", cluster.DefaultWorkerCount, "The number of worker nodes to create on the cluster")

	return cmd
}

func doCreateCluster(cmd *cobra.Command, args []string) error {
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

	workerCount, err := cmd.Flags().GetInt("worker-count")
	if err != nil {
		return err
	}

	opts := cluster.CreateOptions{
		LoadBalancerHostPort: lbHostPort,
		ImageRegistryName:    registryName,
		ImageRegistryPort:    registryPort,
		WorkerCount:          workerCount,
	}

	return createCluster(ctx, opts)
}

func createCluster(ctx context.Context, opts cluster.CreateOptions) error {
	cm := cluster.NewManager(ClusterConfig)

	if exists, _ := cm.Exists(ctx); !exists {
		Dialog.Info("Creating a new dev cluster")
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

		initOpts := dev.InitializeOptions{
			ImageRegistryPort: opts.ImageRegistryPort,
		}

		if err := dm.Initialize(ctx, initOpts); err != nil {
			return err
		}

		Dialog.Info("Relay dev cluster is ready to use")
	} else {
		Dialog.Info("Cluster already exists")
	}

	return nil
}

func newInitClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "initialize",
		Aliases: []string{"init"},
		Short:   "Initialize the local cluster",
		RunE:    doInitCluster,
	}

	return cmd
}

func doInitCluster(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	opts := cluster.InitializeOptions{}

	return initCluster(ctx, opts)
}

func initCluster(ctx context.Context, opts cluster.InitializeOptions) error {
	cm := cluster.NewManager(ClusterConfig)

	if exists, _ := cm.Exists(ctx); !exists {
		Dialog.Info("Cluster does not exist")
		return nil
	}

	cl, err := cm.GetClient(ctx, cluster.ClientOptions{Scheme: dev.DefaultScheme})
	if err != nil {
		return err
	}

	dm := dev.NewManager(cm, cl, DevConfig)

	logServiceOpts := dev.LogServiceOptions{
		Enabled: false,
	}

	if Config.LogServiceConfig != nil {
		logServiceOpts = dev.LogServiceOptions{
			Enabled:               Config.LogServiceConfig.Enabled,
			CredentialsSecretName: Config.LogServiceConfig.CredentialsSecretName,
			Project:               Config.LogServiceConfig.Project,
			Dataset:               Config.LogServiceConfig.Dataset,
			Table:                 Config.LogServiceConfig.Table,
		}
	}

	Dialog.Info("Initializing relay-core; this might take a couple minutes...")

	if err := dm.InitializeRelayCore(ctx, logServiceOpts); err != nil {
		return err
	}

	Dialog.Infof("Cluster connection can be used with: !Connection {type: kubernetes, name: %s}", dev.RelayClusterConnectionName)

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

	Dialog.Progress("Deleting cluster; this may take a minute...")

	cl, err := cm.GetClient(ctx, cluster.ClientOptions{Scheme: dev.DefaultScheme})
	if err != nil {
		return err
	}

	dm := dev.NewManager(cm, cl, DevConfig)

	return dm.Delete(ctx)
}
