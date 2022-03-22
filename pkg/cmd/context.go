package cmd

import (
	"github.com/puppetlabs/relay/pkg/config"
	"github.com/spf13/cobra"
)

func newContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Manage Relay context",
		Args:  cobra.ExactArgs(0),
	}

	cmd.AddCommand(newSetContext())
	cmd.AddCommand(newViewContext())

	return cmd
}

func newSetContext() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [context name]",
		Short: "Set current context",
		Args:  cobra.ExactArgs(1),
		RunE:  doSetContext,
	}

	return cmd
}

func newViewContext() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View current context",
		Args:  cobra.ExactArgs(0),
		RunE:  doViewContext,
	}

	return cmd
}

func doSetContext(cmd *cobra.Command, args []string) error {
	cfg := &config.Config{
		CurrentContext: args[0],
	}

	config.WriteConfig(cfg, cmd.Flags())

	return nil
}

func doViewContext(cmd *cobra.Command, args []string) error {
	context := Config.CurrentContext
	Dialog.Infof("Context: %s", context)

	if contextConfig, ok := Config.ContextConfig[context]; ok {
		if contextConfig.Domains != nil {
			Dialog.Infof("API Domain: %s", contextConfig.Domains.APIDomain)
			Dialog.Infof("UI Domain: %s", contextConfig.Domains.UIDomain)
		} else {
			Dialog.Info("No domains found for current context")
		}
	} else {
		Dialog.Info("No context configuration found")
	}

	return nil
}
