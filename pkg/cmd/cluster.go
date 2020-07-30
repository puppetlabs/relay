package cmd

import (
	"path/filepath"

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

	return cmd
}

func doStartCluster(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cm := cluster.NewManager()
	dm := dev.NewManager(cm, dev.Options{DataDir: filepath.Join(Config.DataDir, "dev")})

	if _, err := cm.Exists(ctx); err != nil {
		Dialog.Info("Creating a new dev cluster")
		if err := cm.Create(ctx); err != nil {
			return err
		}

		Dialog.Info("Writing kubeconfig")
		if err := dm.WriteKubeconfig(ctx); err != nil {
			return err
		}

		Dialog.Info("Applying core Relay resources")
		if err := dm.ApplyCoreResources(ctx); err != nil {
			return err
		}
	} else {
		return cm.Start(ctx)
	}

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
	cm := cluster.NewManager()

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
	cm := cluster.NewManager()
	dm := dev.NewManager(cm, dev.Options{DataDir: filepath.Join(Config.DataDir, "dev")})

	if err := cm.Delete(ctx); err != nil {
		return err
	}

	return dm.DeleteDataDir()
}
