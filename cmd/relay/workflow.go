package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/puppetlabs/relay/pkg/client"
	"github.com/puppetlabs/relay/pkg/config"
	"github.com/puppetlabs/relay/pkg/dialog"
	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/puppetlabs/relay/pkg/format"
	"github.com/puppetlabs/relay/pkg/model"
	"github.com/puppetlabs/relay/pkg/util"
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
		RunE:  addWorkflow,
	}

	cmd.Flags().StringP("file", "f", "", "Path to relay workflow file.")

	return cmd
}

func NewReplaceWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replace [workflow name]",
		Short: "Replace the yaml definition of a relay workflow",
		Args:  cobra.MaximumNArgs(1),
		RunE:  replaceWorkflow,
	}

	cmd.Flags().StringP("file", "f", "", "Path to relay workflow file.")

	return cmd
}

func NewDeleteWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [workflow name]",
		Short: "Delete a relay workflow",
		Args:  cobra.MaximumNArgs(1),
		RunE:  deleteWorkflow,
	}

	return cmd
}

func addWorkflow(cmd *cobra.Command, args []string) error {
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

	revision, rerr := client.CreateRevision(workflow.Workflow.Name, file)

	if rerr != nil {

		// attempt to revert creation of workflow record
		client.DeleteWorkflow(workflow.Workflow.Name)

		return rerr
	}

	wr := model.NewWorkflowRevision(workflow.Workflow, revision.Revision)

	wr.Output(cfg)

	log.Info(
		fmt.Sprintf(
			`Successfully created workflow %v
			
View more information or update workflow settings at %v`,
			workflow.Workflow.Name,
			format.GuiLink(cfg, "/workflow/%v", workflow.Workflow.Name),
		),
	)

	return nil
}

func replaceWorkflow(cmd *cobra.Command, args []string) error {
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

	workflow, werr := client.GetWorkflow(workflowName)

	if werr != nil {
		return werr
	}

	revision, rerr := client.CreateRevision(workflowName, file)

	if rerr != nil {
		return rerr
	}

	wr := model.NewWorkflowRevision(workflow.Workflow, revision.Revision)

	wr.Output(cfg)

	log.Info(
		fmt.Sprintf(
			`Successfully updated workflow %v
			
Updated configuration is visible at %v`,
			workflowName,
			format.GuiLink(cfg, "/workflow/%v", workflowName),
		),
	)

	return nil
}

func deleteWorkflow(cmd *cobra.Command, args []string) error {
	cfg, cfgerr := config.GetConfig(cmd.Flags())

	if cfgerr != nil {
		return cfgerr
	}

	workflowName, nerr := getWorkflowName(args)

	if nerr != nil {
		return nerr
	}

	proceed, cerr := util.Confirm("Are you sure you want to delete this workflow?", cfg)

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

	log.Info("Workflow successfully deleted")

	return nil
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
		return "", errors.NewGeneralUnknownError().WithCause(err)
	}

	if filepath == "" {
		return "", errors.NewWorkflowMissingFileFlagError()
	}

	bytes, err := ioutil.ReadFile(filepath)

	if err != nil {
		return "", errors.NewWorkflowWorkflowFileReadError().WithCause(err)
	}

	return string(bytes), nil
}
