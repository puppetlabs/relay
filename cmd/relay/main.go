package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "relay",
		Short: "Relay by Puppet.",
		Args:  cobra.MinimumNArgs(1),
		Long: `Relay connects your tools, APIs, and infrastructure 
to automate common tasks through simple event driven workflows.`,
	}

	cmd.PersistentFlags().BoolP("verbose", "v", false, "print verbose output")
	cmd.PersistentFlags().BoolP("debug", "d", false, "print debugging information")
	cmd.PersistentFlags().StringP("out", "o", "text", "output type: (text|json)")
	// Config flag is hidden for now
	cmd.PersistentFlags().StringP("config", "c", "", "path to config file (default is $HOME.config/relay)")
	cmd.PersistentFlags().MarkHidden("config")

	cmd.AddCommand(NewAuthCommand())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
