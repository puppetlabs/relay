package main

import (
	"encoding/json"
	"fmt"

	"github.com/puppetlabs/relay/pkg/client"
	"github.com/puppetlabs/relay/pkg/config"
	"github.com/puppetlabs/relay/pkg/dialog"
	"github.com/spf13/cobra"
)

func NewWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage your relay workflows",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(NewAddWorkflowCommand())

	return cmd
}

func NewAddWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [workflow]",
		Short: "Add Relay workflow",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, cfgerr := config.GetConfig(cmd.Flags())

			if cfgerr != nil {
				return cfgerr
			}

			log := dialog.NewDialog(cfg)

			log.Info("Creating your workflow...")

			client := client.NewClient(cfg)

			workflow, cwerr := client.CreateWorkflow(args[0])

			if cwerr != nil {
				return cwerr
			}

			jsonBytes, _ := json.MarshalIndent(workflow, "", "  ")

			fmt.Println(string(jsonBytes))

			return nil
		},
	}

	return cmd
}
