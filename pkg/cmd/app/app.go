package app

import (
	"github.com/puppetlabs/nebula/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand(r config.CLIRuntime) *cobra.Command {
	return &cobra.Command{
		Use:   "app",
		Short: "Manage application versions and deployments",
		Run: func(cmd *cobra.Command, args []string) {

		},
		SuggestFor: []string{"deploy", "up"},
	}
}
