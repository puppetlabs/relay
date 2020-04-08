package integration

import (
	"context"
	"fmt"

	"github.com/jedib0t/go-pretty/table"
	"github.com/puppetlabs/relay/pkg/client"
	"github.com/puppetlabs/relay/pkg/config/runtimefactory"
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
			cfg, err := rt.Config()
			if err != nil {
				return err
			}

			c, err := client.NewAPIClient(cfg)
			if err != nil {
				return err
			}
			integrations, err := c.ListIntegrations(context.Background())
			if err != nil {
				return err
			}

			tw := table.NewWriter()

			tw.AppendHeader(table.Row{"PROVIDER", "LOGIN"})
			for _, integration := range integrations {
				tw.AppendRow(table.Row{*integration.Provider, integration.AccountLogin})
			}

			_, _ = fmt.Fprintf(rt.IO().Out, "%s\n", tw.Render())

			return nil
		},
	}

	return cmd
}
