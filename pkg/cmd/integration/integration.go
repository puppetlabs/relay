package integration

import (
	"context"
	"fmt"

	"github.com/jedib0t/go-pretty/table"
	"github.com/puppetlabs/nebula-cli/pkg/client"
	"github.com/puppetlabs/nebula-cli/pkg/config/runtimefactory"
	"github.com/spf13/cobra"
)

func NewCommand(rt runtimefactory.RuntimeFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "integration",
		Short:                 "Manage nebula integrations",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewListCommand(rt))

	return cmd
}

func NewListCommand(rt runtimefactory.RuntimeFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		Short:                 "List integrations",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.NewAPIClient(rt.Config())
			if err != nil {
				return err
			}
			index, err := c.ListIntegrations(context.Background())
			if err != nil {
				return err
			}

			tw := table.NewWriter()

			tw.AppendHeader(table.Row{"ID", "PROVIDER", "ACCOUNT LOGIN"})
			for _, i := range index {
				integrationName := fmt.Sprintf("%s-%s", *i.Provider, i.AccountLogin)

				tw.AppendRow(table.Row{integrationName, *i.Provider, i.AccountLogin})
			}
			_, _ = fmt.Fprintf(rt.IO().Out, "%s\n", tw.Render())

			return nil
		},
	}

	return cmd
}
