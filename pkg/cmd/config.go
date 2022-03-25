package cmd

import (
	"fmt"
	"strconv"
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
	cmd.AddCommand(newConfigGlobalCommand())

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

func newConfigGlobalCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "global",
		Short: "Manage Relay CLI global options",
		Args:  cobra.ExactArgs(0),
	}

	cmd.AddCommand(newConfigDebugFlagCommand())
	cmd.AddCommand(newConfigOutFlagCommand())
	cmd.AddCommand(newConfigYesFlagCommand())

	return cmd
}

func newConfigDebugFlagCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug (true|false)",
		Short: "Set global debug flag",
		Args:  cobra.ExactArgs(1),
		RunE:  doConfigSetDebugFlag,
	}

	return cmd
}

func doConfigSetDebugFlag(cmd *cobra.Command, args []string) error {
	debug, err := strconv.ParseBool(args[0])
	if err != nil {
		return err
	}

	return config.WriteGlobalConfig(&config.Config{
		Debug: debug,
		Out:   Config.Out,
		Yes:   Config.Yes,
	}, cmd.Flags())
}

func newConfigOutFlagCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "out (text|json)",
		Short: "Set global out flag",
		Args:  cobra.ExactArgs(1),
		RunE:  doConfigSetOutFlag,
	}

	return cmd
}
func doConfigSetOutFlag(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case config.OutputTypeJSON.String():
		return config.WriteGlobalConfig(&config.Config{
			Debug: Config.Debug,
			Out:   config.OutputTypeJSON,
			Yes:   Config.Yes,
		}, cmd.Flags())
	case config.OutputTypeText.String():
		return config.WriteGlobalConfig(&config.Config{
			Debug: Config.Debug,
			Out:   config.OutputTypeText,
			Yes:   Config.Yes,
		}, cmd.Flags())
	default:
		return fmt.Errorf("invalid output type: %s", args[0])
	}
}

func newConfigYesFlagCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yes (true|false)",
		Short: "Set global yes flag",
		Args:  cobra.ExactArgs(1),
		RunE:  doConfigSetYesFlag,
	}

	return cmd
}

func doConfigSetYesFlag(cmd *cobra.Command, args []string) error {
	yes, err := strconv.ParseBool(args[0])
	if err != nil {
		return err
	}

	return config.WriteGlobalConfig(&config.Config{
		Debug: Config.Debug,
		Out:   Config.Out,
		Yes:   yes,
	}, cmd.Flags())
}
