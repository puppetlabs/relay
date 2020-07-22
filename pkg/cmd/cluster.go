package cmd

import (
	"path/filepath"

	"github.com/puppetlabs/relay/pkg/cluster"
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

	opts := cluster.ClusterOptions{
		DataDir: filepath.Join(Config.DataDir, "cluster"),
	}

	if _, err := cluster.GetCluster(ctx); err != nil {
		if err := cluster.CreateCluster(ctx, opts); err != nil {
			return err
		}
	} else {
		return cluster.StartCluster(ctx)
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

	return cluster.StopCluster(ctx)
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

	return cluster.DeleteCluster(ctx)
}
