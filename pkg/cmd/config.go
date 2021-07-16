package cmd

import (
	"github.com/puppetlabs/relay/pkg/config"
	"github.com/spf13/cobra"
)

func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage Relay CLI configuration",
		Args:  cobra.ExactArgs(0),
	}

	cmd.AddCommand(newConfigAuthCommand())

	return cmd
}

func newConfigAuthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage Relay CLI authentication configuration",
		Args:  cobra.ExactArgs(0),
	}

	cmd.AddCommand(newConfigAuthClearCommand())

	return cmd
}

func newConfigAuthClearCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear the stored authentication for the current context",
		Args:  cobra.ExactArgs(0),
		RunE:  doConfigAuthClear,
	}

	return cmd
}

func doConfigAuthClear(cmd *cobra.Command, args []string) error {
	cfg := &config.Config{
		ContextConfig: map[string]*config.ContextConfig{
			Config.CurrentContext: {
				Auth: &config.AuthConfig{
					Tokens: map[config.AuthTokenType]string{
						config.AuthTokenTypeAPI:     "",
						config.AuthTokenTypeSession: "",
					},
				},
			},
		},
	}

	config.WriteConfig(cfg, cmd.Flags())

	return nil
}
