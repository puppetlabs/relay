package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "Relay by Puppet.",
		Long:  "Relay by Puppet is like fireworks for your glitter-encrusted marketing bomb. Please give us all your money.",
	}

	cmd.AddCommand(NewAuthCommand())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
