package create

import (
	"context"

	"github.com/puppetlabs/nebula/pkg/config"
	"github.com/puppetlabs/nebula/pkg/workflow"
	"github.com/spf13/cobra"
)

func NewCommand(r config.CLIRuntime) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create [options] [command]",
		Short:                 "Initialize and create workflow resources",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var wf workflow.Workflow
			if err := r.WorkflowLoader().Load(&wf); err != nil {
				return err
			}

			for _, action := range wf.Actions {
				if err := action.Runner().Run(context.Background(), r, nil); err != nil {
					return err
				}
			}

			return nil
		},
	}

	return cmd
}
