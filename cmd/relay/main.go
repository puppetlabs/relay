package main

import (
<<<<<<< HEAD
	"log"
=======
	"fmt"
	"os"
>>>>>>> 62e767f... Refactor config package

	"github.com/puppetlabs/relay/pkg/cmd"
)

func main() {
	cmd := &cobra.Command{
		Use:           "relay",
		Short:         "Relay by Puppet.",
		Args:          cobra.MinimumNArgs(1),
		SilenceErrors: true,
		Long: `Relay connects your tools, APIs, and infrastructure 
to automate common tasks through simple event driven workflows.`,
	}

	cmd.PersistentFlags().BoolP("debug", "d", false, "print debugging information")
	cmd.PersistentFlags().StringP("out", "o", "text", "output type: (text|json)")
	// Config flag is hidden for now
	cmd.PersistentFlags().StringP("config", "c", "", "path to config file (default is $HOME.config/relay)")
	cmd.PersistentFlags().MarkHidden("config")

	cmd.AddCommand(NewAuthCommand())

	// TODO: Errawr formatter.
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
