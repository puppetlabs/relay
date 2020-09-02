package cmd

import (
	"github.com/puppetlabs/relay/pkg/cluster"
	"github.com/spf13/cobra"
)

func newImageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Manage container images used inside the cluster",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(newImageImportCommand())

	return cmd
}

func newImageImportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <image:tag>",
		Short: "Imports a container image into the cluster",
		Args:  cobra.MinimumNArgs(1),
		RunE:  doImageImport,
	}

	return cmd
}

func doImageImport(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cm := cluster.NewManager(ClusterConfig)

	return cm.ImportImages(ctx, args[0])
}
