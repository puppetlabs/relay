package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/puppetlabs/relay/pkg/client"
	"github.com/puppetlabs/relay/pkg/config"
	"github.com/puppetlabs/relay/pkg/dialog"
	"github.com/puppetlabs/relay/pkg/errors"
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
		Use:   "add [workflow name]",
		Short: "Add Relay workflow",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, cfgerr := config.GetConfig(cmd.Flags())

			if cfgerr != nil {
				return cfgerr
			}

			file, ferr := readFile(cmd)

			if ferr != nil {
				return ferr
			}

			workflowName, nerr := getWorkflowName(args)

			if nerr != nil {
				return nerr
			}

			log := dialog.NewDialog(cfg)

			log.Info("Creating your workflow...")

			client := client.NewClient(cfg)

			workflow, cwerr := client.CreateWorkflow(workflowName)

			if cwerr != nil {
				return cwerr
			}

			_, rerr := client.CreateRevision(workflow.Workflow.Name, file)

			if rerr != nil {

				// attempt to revert creation of workflow record
				client.DeleteWorkflow(workflow.Workflow.Name)

				return rerr
			}

			// TODO: JSON and Text formatters for workflow and revision objects
			log.Info(fmt.Sprint("Successfully created workflow", workflow.Workflow.Name))

			return nil
		},
	}

	cmd.Flags().StringP("file", "f", "", "Path to relay workflow file.")

	return cmd
}

func getWorkflowName(args []string) (string, errors.Error) {
	if len(args) > 0 {
		return args[0], nil
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Workflow name: ")
	namePrompt, err := reader.ReadString('\n')

	if err != nil {
		return "", errors.NewWorkflowWorkflowNameReadError().WithCause(err)
	}

	return strings.TrimSpace(namePrompt), nil
}

func readFile(cmd *cobra.Command) (string, errors.Error) {
	filepath, err := cmd.Flags().GetString("file")

	if err != nil {
		return "", errors.NewWorkflowWorkflowFileReadError().WithCause(err)
	}

	if filepath == "" {
		return "", errors.NewWorkflowMissingFileFlagError()
	}

	file, err := os.Open(filepath)

	if err != nil {
		return "", errors.NewWorkflowWorkflowFileReadError().WithCause(err)
	}

	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(file); err != nil {
		return "", errors.NewWorkflowWorkflowFileReadError().WithCause(err)
	}

	return buf.String(), nil
}
