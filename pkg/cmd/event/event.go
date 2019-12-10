package event

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
		Use:                   "event",
		Short:                 "Manage nebula event sources",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewListCommand(rt))

	return cmd
}

func NewListCommand(rt runtimefactory.RuntimeFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		Short:                 "List event sources",
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
			es, err := c.ListEventSources(context.Background())
			if err != nil {
				return err
			}

			tw := table.NewWriter()

			tw.AppendHeader(table.Row{"NAME", "PROVIDER", "LOGIN"})
			for _, e := range es {
				providerMap := e.Provider.(map[string]interface{})
				providerType := providerMap["type"]

				switch providerType {
				case "integration":
					integrationMap := providerMap["integration"].(map[string]interface{})
					tw.AppendRow(table.Row{*e.EventType.Name, integrationMap["provider"], integrationMap["account_login"]})
				}
			}
			_, _ = fmt.Fprintf(rt.IO().Out, "%s\n", tw.Render())

			return nil
		},
	}

	return cmd
}
