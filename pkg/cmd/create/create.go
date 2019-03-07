package create

import (
	"context"

	"github.com/puppetlabs/nebula/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand(r config.CLIRuntime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [options] [command]",
		Short: "Initialize and create workflow resources",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			workflow, err := r.WorkflowLoader().Load()
			if err != nil {
				return err
			}

			variables := make(map[string]string)

			r.Logger().Info("Running workflow stage")
			for _, a := range workflow.Stage.Actions {
				r.Logger().Info("Run() "+a.Name, "action name", a.Name)
				a.Runner().Run(context.Background(), variables)
			}

			return nil
		},
	}

	return cmd
}
