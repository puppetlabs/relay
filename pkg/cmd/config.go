package cmd

import (
	"fmt"
	"strings"

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
		Short: "Clear stored authentication data for the current context",
		Args:  cobra.ExactArgs(0),
		RunE:  doConfigAuthClear,
	}

	cmd.Flags().StringP("type", "t", "",
		fmt.Sprintf("Authentication type (%s)", strings.Join(config.AuthTokenTypesAsString(), "|")))

	return cmd
}

func doConfigAuthClear(cmd *cobra.Command, args []string) error {
	authTokenType, err := cmd.Flags().GetString("type")
	if err != nil {
		return err
	}

	context := Config.CurrentContext
	cfg := &config.Config{
		ContextConfig: map[string]*config.ContextConfig{
			context: {
				Auth: &config.AuthConfig{
					Tokens: map[config.AuthTokenType]string{},
				},
			},
		},
	}

	for _, tokenType := range config.AuthTokenTypes() {
		if authTokenType == "" || authTokenType == tokenType.String() {
			cfg.ContextConfig[context].Auth.Tokens[tokenType] = ""
		}
	}

	config.WriteConfig(cfg, cmd.Flags())

	return nil
}
