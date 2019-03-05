package cmd

import (
	"github.com/puppetlabs/nebula/pkg/cmd/app"
	"github.com/puppetlabs/nebula/pkg/cmd/infra"
	"github.com/puppetlabs/nebula/pkg/cmd/workflow"
	"github.com/puppetlabs/nebula/pkg/config"
	"github.com/spf13/cobra"
)

func NewRootCommand() (*cobra.Command, error) {
	c := &cobra.Command{
		Use:   "nebula",
		Short: "Nebula workflow management cli",
	}

	r, err := config.NewCLIRuntime()
	if err != nil {
		return nil, err
	}

	c.AddCommand(infra.NewCommand(r))
	c.AddCommand(app.NewCommand(r))
	c.AddCommand(workflow.NewCommand(r))

	return c, nil
}
