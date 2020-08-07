package cmd

import "github.com/spf13/cobra"

func newDevCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Manage the local development environment",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(newClusterCommand())
	cmd.AddCommand(newImageCommand())
	cmd.AddCommand(newKubectlCommand())

	return cmd
}
