package cmd

import (
	"os"

	"github.com/puppetlabs/relay/pkg/client"
	"github.com/puppetlabs/relay/pkg/config"
	"github.com/puppetlabs/relay/pkg/debug"
	"github.com/puppetlabs/relay/pkg/dialog"
	"github.com/puppetlabs/relay/pkg/format"
	"github.com/spf13/cobra"
)

// CommandName is the top level command that we're building
const CommandName = "relay"

// Config is the configuration that our commands should use. We can assume that
// it's been configured accordingly by the time that a command executres.
var Config = config.GetDefaultConfig()

// Client is the client that we should use based on the configuration. If the
// configuration can't be loaded then we can't assume that the client is
// loaded.
var Client = client.NewClient(Config)

// Dialog is the UI to use derrived from the current configuration.
var Dialog = dialog.FromConfig(Config)

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           CommandName,
		Short:         "Relay by Puppet",
		Args:          cobra.MinimumNArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		Long: `Relay connects your tools, APIs, and infrastructure
to automate common tasks through simple, event-driven workflows.

To get started, you'll need a relay.sh account - sign up for free
by following this link: 🔗 https://relay.sh/

Once you've signed up, run this to log in:
▶️   relay auth login

Use the 'workflow' subcommand to interact with workflows:
▶️   relay workflow
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// This turns off usage info in json output mode
			cfg, err := config.FromFlags(cmd.Flags())

			if err != nil {
				// What kind of error could this be? We will abort accordingly.
				return err
			} else if err == nil && cfg.Out == config.OutputTypeJSON {
				cmd.SilenceUsage = true
			}

			// We have a config that we can assume is good to use.
			Config = cfg
			Client = client.NewClient(Config)
			Dialog = dialog.FromConfig(Config)

			return nil
		},
	}

	cmd.PersistentFlags().BoolVarP(&debug.Enabled, "debug", "d", false, "print debugging information")
	cmd.PersistentFlags().BoolP("yes", "y", false, "skip confirmation prompts")
	cmd.PersistentFlags().StringP("out", "o", "text", "output type: (text|json)")

	// allow the user to override the default configuration location if they
	// can find the flag. likely figured out from reading this comment, actually...
	cmd.PersistentFlags().StringP("config", "c", "", "path to config file (default is $HOME.config/relay)")
	cmd.PersistentFlags().MarkHidden("config")

	cmd.AddCommand(newAuthCommand())
	cmd.AddCommand(newWorkflowCommand())

	return cmd
}

// Execute is here so the cmd builder can be called from the go test harness
func Execute() {
	cmd := getCmd()

	if err := cmd.Execute(); err != nil {
		Dialog.Error(format.Error(err, cmd))
		os.Exit(1)
	}
}
