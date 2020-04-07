package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "Relay by puppet. Blah-de-blah-blah give me all your money",
	}

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
