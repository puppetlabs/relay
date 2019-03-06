package create

import (
	"github.com/kr/pretty"
	"github.com/puppetlabs/nebula/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand(r config.CLIRuntime) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create [options] [command]",
		Short:                 "Initialize and create workflow resources",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			workflow, err := r.WorkflowLoader().Load()
			if err != nil {
				return err
			}

			pretty.Println(workflow)

			return nil
		},
	}

	return cmd
}
