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
	"github.com/puppetlabs/relay/pkg/util/confirm"
	"github.com/spf13/cobra"
)

func NewWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage your relay workflows",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(NewAddWorkflowCommand())
	cmd.AddCommand(NewReplaceWorkflowCommand())
	cmd.AddCommand(NewDeleteWorkflowCommand())

	return cmd
}

func NewAddWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [workflow name]",
		Short: "Add a relay workflow from a local file",
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
			log.Info(fmt.Sprint("Successfully created workflow ", workflow.Workflow.Name))

			return nil
		},
	}

	cmd.Flags().StringP("file", "f", "", "Path to relay workflow file.")

	return cmd
}

func NewReplaceWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replace [workflow name]",
		Short: "Replace the yaml definition of a relay workflow",
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

			log.Info(fmt.Sprint("Replacing workflow ", workflowName))

			client := client.NewClient(cfg)

			_, rerr := client.CreateRevision(workflowName, file)

			if rerr != nil {
				return rerr
			}

			// TODO: JSON and Text formatters for workflow and revision objects
			log.Info(fmt.Sprint("Successfully replaced workflow ", workflowName))

			return nil
		},
	}

	cmd.Flags().StringP("file", "f", "", "Path to relay workflow file.")

	return cmd
}

func NewDeleteWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [workflow name]",
		Short: "Delete a relay workflow",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, cfgerr := config.GetConfig(cmd.Flags())

			if cfgerr != nil {
				return cfgerr
			}

			workflowName, nerr := getWorkflowName(args)

			if nerr != nil {
				return nerr
			}

			proceed, cerr := confirm.Confirm("Are you sure you want to delete this workflow?", cfg)

			if cerr != nil {
				return cerr
			}

			if !proceed {
				return nil
			}

			log := dialog.NewDialog(cfg)

			log.Info("Deleting workflow...")

			client := client.NewClient(cfg)

			_, err := client.DeleteWorkflow(workflowName)

			if err != nil {
				return err
			}

			// TODO: log response object in json mode
			log.Info("Workflow successfully deleted")

			return nil
		},
	}

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

	name := strings.TrimSpace(namePrompt)

	if name == "" {
		return "", errors.NewWorkflowMissingNameError()
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
