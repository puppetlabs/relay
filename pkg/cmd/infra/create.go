package infra

import (
	"github.com/puppetlabs/nebula/pkg/config"
	"github.com/spf13/cobra"
)

func NewCreateCommand(r config.CLIRuntime) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Creates an environment defined in definitions.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			r.Logger().Debug("infra create command")
		},
	}
}
