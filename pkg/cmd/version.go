package cmd

import (
	"fmt"

	"github.com/puppetlabs/relay/pkg/version"
	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: `Print version`,
		Run: func(cmd *cobra.Command, args []string) {
			v := version.GetVersion()
			if v == "" {
				fmt.Println("could not determine build information")
			} else {
				fmt.Println(v)
			}
		},
	}
}
