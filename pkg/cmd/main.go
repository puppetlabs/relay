package cmd

import (
	"os"

	"github.com/puppetlabs/relay/pkg/cmd/auth"
	"github.com/puppetlabs/relay/pkg/cmd/workflow"
	"github.com/puppetlabs/relay/pkg/config"
	"github.com/puppetlabs/relay/pkg/format"
	"github.com/spf13/cobra"
)

func Execute() {
	cmd := &cobra.Command{
		Use:           "relay",
		Short:         "Relay by Puppet",
		Args:          cobra.MinimumNArgs(1),
		SilenceErrors: true,
		Long: `Relay connects your tools, APIs, and infrastructure 
to automate common tasks through simple event driven workflows.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// This turns off usage info in json output mode
			cfg, cfgerr := config.GetConfig(cmd.Flags())

			if cfgerr == nil && cfg.Out == config.OutputTypeJSON {
				cmd.SilenceUsage = true
			}
		},
	}

	cmd.PersistentFlags().BoolP("debug", "d", false, "print debugging information")
	cmd.PersistentFlags().BoolP("yes", "y", false, "skip confirmation prompts")
	cmd.PersistentFlags().StringP("out", "o", "text", "output type: (text|json)")
	// Config flag is hidden for now
	cmd.PersistentFlags().StringP("config", "c", "", "path to config file (default is $HOME.config/relay)")
	cmd.PersistentFlags().MarkHidden("config")

	cmd.AddCommand(auth.NewAuthCommand())
	cmd.AddCommand(workflow.NewCommand())

	if err := cmd.Execute(); err != nil {
		format.FormatError(err, cmd)
		os.Exit(1)
	}
}
