package cmd

import (
	"github.com/puppetlabs/nebula/pkg/cmd/integration"
	"github.com/puppetlabs/nebula/pkg/cmd/login"
	"github.com/puppetlabs/nebula/pkg/cmd/secret"
	"github.com/puppetlabs/nebula/pkg/cmd/version"
	"github.com/puppetlabs/nebula/pkg/cmd/workflow"
	"github.com/puppetlabs/nebula/pkg/config/runtimefactory"
	"github.com/spf13/cobra"
)

func NewRootCommand() (*cobra.Command, error) {
	r, err := runtimefactory.NewRuntimeFactory()
	if err != nil {
		return nil, err
	}

	c := &cobra.Command{
		Use:   "nebula",
		Short: "Nebula workflow management cli",
		// don't show usage text for every error
		SilenceUsage: true,
		// we want to be able to handle our own errors for display; this allows us to use
		// the CLI display mechanism for errawr.
		SilenceErrors: true,
	}

	c.AddCommand(login.NewCommand(r))
	c.AddCommand(workflow.NewCommand(r))
	c.AddCommand(integration.NewCommand(r))
	c.AddCommand(version.NewCommand(r))
	c.AddCommand(secret.NewCommand(r))

	return c, nil
}
