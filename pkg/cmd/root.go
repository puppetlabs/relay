package cmd

import (
	"github.com/puppetlabs/relay/pkg/cmd/event"
	"github.com/puppetlabs/relay/pkg/cmd/integration"
	"github.com/puppetlabs/relay/pkg/cmd/login"
	"github.com/puppetlabs/relay/pkg/cmd/secret"
	"github.com/puppetlabs/relay/pkg/cmd/version"
	"github.com/puppetlabs/relay/pkg/cmd/workflow"
	"github.com/puppetlabs/relay/pkg/config/runtimefactory"
	"github.com/spf13/cobra"
)

func NewRootCommand() (*cobra.Command, error) {
	c := &cobra.Command{
		Use:   "nebula",
		Short: "Nebula workflow management cli",
		// don't show usage text for every error
		SilenceUsage: true,
		// we want to be able to handle our own errors for display; this allows us to use
		// the CLI display mechanism for errawr.
		SilenceErrors: true,
	}

	r := runtimefactory.NewRuntimeFactory(c.Flags())

	c.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.config/nebula/config.yaml)")
	c.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Check if we can load the config file
		if _, err := r.Config(); err != nil {
			return err
		}

		return nil
	}

	c.AddCommand(login.NewCommand(r))
	c.AddCommand(workflow.NewCommand(r))
	c.AddCommand(integration.NewCommand(r))
	c.AddCommand(version.NewCommand(r))
	c.AddCommand(secret.NewCommand(r))
	c.AddCommand(event.NewCommand(r))

	return c, nil
}
