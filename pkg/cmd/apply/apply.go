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

			stageName, err := cmd.Flags().GetString("stage")
			if err != nil {
				return err
			}

			stage, err := wf.GetStage(stageName)
			if err != nil {
				return err
			}

			r.Logger().Info("running stage", "stage", stage.Name)

			variables := make(map[string]string)

			for _, v := range wf.Variables {
				variables[v.Name] = v.Value
			}

			for _, action := range stage.Actions {
				r.Logger().Info("action-started", "action", action.Name, "kind", action.Kind)
				if err := action.Runner().Run(context.Background(), r, variables); err != nil {
					return err
				}
				r.Logger().Info("action-finished", "action", action.Name, "kind", action.Kind)
			}

			r.Logger().Info("workflow-applied")

			return nil
		},
	}

	cmd.Flags().StringP("stage", "s", "", "name of the stage to run actions for")

	return cmd
}
