package workflow

import (
	"github.com/puppetlabs/nebula/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand(r config.CLIRuntime) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "workflow [options] [command]",
		Short:                 "Manage workflows",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	cmd.AddCommand(NewWorkflowCommand(r))

	return cmd
}
