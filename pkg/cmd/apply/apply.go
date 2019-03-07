package apply

import (
	"github.com/kr/pretty"
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

			pretty.Println(wf)

			return nil
		},
	}

	return cmd
}
