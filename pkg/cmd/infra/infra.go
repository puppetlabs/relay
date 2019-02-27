package infra

import (
	"github.com/puppetlabs/nebula/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand(r config.CLIRuntime) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "infra [options] [command]",
		Short:                 "Manage infrastructure and environment state",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	cmd.AddCommand(NewCreateCommand(r))

	return cmd
}
