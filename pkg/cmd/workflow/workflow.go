package workflow

import (
	"context"
	"fmt"

	"github.com/puppetlabs/nebula/pkg/client"
	"github.com/puppetlabs/nebula/pkg/config/runtimefactory"
	"github.com/spf13/cobra"
)

func NewCommand(rt runtimefactory.RuntimeFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "workflow",
		Short:                 "Manage nebula workflows",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewListCommand(rt))

	return cmd
}

func NewListCommand(rt runtimefactory.RuntimeFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		Short:                 "List workflows",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := client.NewAPIClient(rt.Config())
			if err != nil {
				return err
			}

			index, err := client.ListWorkflows(context.Background())
			if err != nil {
				return err
			}

			fmt.Fprintf(rt.IO().Out, "%v\n", index)

			return nil
		},
	}

	return cmd
}
