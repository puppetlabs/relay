package version

import (
	"fmt"

	"github.com/puppetlabs/nebula/pkg/config/runtimefactory"
	"github.com/puppetlabs/nebula/pkg/version"
	"github.com/spf13/cobra"
)

func NewCommand(rt runtimefactory.RuntimeFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version and build information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(rt.IO().Out, "Nebula version %s\n", version.Version)
		},
	}

	return cmd
}
