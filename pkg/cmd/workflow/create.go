package workflow

import (
	"github.com/puppetlabs/nebula/pkg/config"
	"github.com/spf13/cobra"
)

func NewWorkflowCommand(r config.CLIRuntime) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Creates workflows defined in workflow.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			r.Logger().Debug("workflow create command")
		},
	}
}
