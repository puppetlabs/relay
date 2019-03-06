package cmd

import (
	"github.com/puppetlabs/nebula/pkg/cmd/apply"
	"github.com/puppetlabs/nebula/pkg/cmd/create"
	"github.com/puppetlabs/nebula/pkg/config"
	"github.com/puppetlabs/nebula/pkg/workflow/loader"
	"github.com/spf13/cobra"
)

func NewRootCommand() (*cobra.Command, error) {
	r, err := config.NewCLIRuntime()
	if err != nil {
		return nil, err
	}

	c := &cobra.Command{
		Use:   "nebula",
		Short: "Nebula workflow management cli",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			fp, err := cmd.Flags().GetString("filepath")
			if err != nil {
				return err
			}

			if fp != "" {
				r.SetWorkflowLoader(loader.NewFilepathLoader(fp))
			}

			return nil
		},
	}

	c.PersistentFlags().StringP("filepath", "f", "", "optional path to a workflow.yaml")

	c.AddCommand(create.NewCommand(r))
	c.AddCommand(apply.NewCommand(r))

	return c, nil
}
