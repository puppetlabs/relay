package apply

import (
	"context"

	"github.com/puppetlabs/nebula/pkg/config"
	"github.com/puppetlabs/nebula/pkg/workflow"
	"github.com/spf13/cobra"
)

func NewCommand(r config.CLIRuntime) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "apply [options] [command]",
		Short:                 "Apply and run workflow stages",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var wf workflow.Workflow

			if err := r.WorkflowLoader().Load(&wf); err != nil {
				return err
			}
			r.Logger().Debug("workflow-loaded")

			for _, action := range wf.Actions {
				r.Logger().Info("action-started", "action", action.Name, "kind", action.Kind)
				if err := action.Runner().Run(context.Background(), r, nil); err != nil {
					return err
				}
				r.Logger().Info("action-finished", "action", action.Name, "kind", action.Kind)
			}

			r.Logger().Info("workflow-applied")

			return nil
		},
	}

	return cmd
}
