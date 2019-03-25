package apply

import (
	"context"
	"strings"

	"github.com/puppetlabs/nebula/pkg/config/runtimefactory"
	"github.com/puppetlabs/nebula/pkg/plan"
	"github.com/puppetlabs/nebula/pkg/plan/types"
	"github.com/spf13/cobra"
)

func NewCommand(rt runtimefactory.RuntimeFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "apply [options] [command]",
		Short:                 "Apply and actions from Nebula workflows",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var p types.Plan

			if err := rt.PlanLoader().Load(&p); err != nil {
				return err
			}

			rt.Logger().Debug("nebula-plan-loaded")

			workflowName, err := cmd.Flags().GetString("workflow")
			if err != nil {
				return err
			}

			env, err := cmd.Flags().GetStringSlice("env")
			if err != nil {
				return err
			}

			for _, v := range env {
				parts := strings.Split(v, "=")

				p.Variables = append(p.Variables, &types.Variable{Name: parts[0], Value: parts[1]})
			}

			runner, err := plan.NewWorkflowRunnerFromName(workflowName, &p, rt.ActionExecutor())
			if err != nil {
				return err
			}

			rt.Logger().Info("running-workflow", "workflow", runner.Name())

			if err := runner.Run(context.Background(), rt); err != nil {
				return err
			}

			rt.Logger().Info("nebula-plan-applied")

			return nil
		},
	}

	cmd.Flags().StringP("workflow", "w", "", "name of the workflow to run actions for")
	cmd.Flags().StringSliceP("env", "e", []string{}, "sets environment variables for actions (e.g. --env=KEY=value,KEY2=value2)")

	return cmd
}
