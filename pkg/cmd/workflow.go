package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/puppetlabs/relay/pkg/format"
	"github.com/puppetlabs/relay/pkg/model"
	"github.com/puppetlabs/relay/pkg/util"
	"github.com/spf13/cobra"
)

func newWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage your relay workflows",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(newAddWorkflowCommand())
	cmd.AddCommand(newReplaceWorkflowCommand())
	cmd.AddCommand(newDeleteWorkflowCommand())
	cmd.AddCommand(newRunWorkflowCommand())
	cmd.AddCommand(newListWorkflowsCommand())

	return cmd
}

func doAddWorkflow(cmd *cobra.Command, args []string) error {
	file, ferr := readFile(cmd)

	if ferr != nil {
		return ferr
	}

	workflowName, nerr := getWorkflowName(args)

	if nerr != nil {
		return nerr
	}

	Dialog.Progress("Creating your workflow...")

	workflow, cwerr := Client.CreateWorkflow(workflowName)

	if cwerr != nil {
		return cwerr
	}

	revision, rerr := Client.CreateRevision(workflow.Workflow.Name, file)

	if rerr != nil {

		// attempt to revert creation of workflow record
		Client.DeleteWorkflow(workflow.Workflow.Name)

		return rerr
	}

	wr := model.NewWorkflowRevision(workflow.Workflow, revision.Revision)

	wr.Output(Config)

	Dialog.Infof(`Successfully created workflow %v
			
View more information or update workflow settings at %v`,
		workflow.Workflow.Name,
		format.GuiLink(Config, "/workflow/%v", workflow.Workflow.Name),
	)

	return nil
}

func newAddWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [workflow name]",
		Short: "Add a relay workflow from a local file",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doAddWorkflow,
	}

	cmd.Flags().StringP("file", "f", "", "Path to relay workflow file.")

	return cmd
}

func doReplaceWorkflow(cmd *cobra.Command, args []string) error {
	file, err := readFile(cmd)

	if err != nil {
		return err
	}

	workflowName, err := getWorkflowName(args)

	if err != nil {
		return err
	}

	Dialog.Info("Replacing workflow " + workflowName)

	workflow, werr := Client.GetWorkflow(workflowName)

	if werr != nil {
		return werr
	}

	revision, rerr := Client.CreateRevision(workflowName, file)

	if rerr != nil {
		return rerr
	}

	wr := model.NewWorkflowRevision(workflow.Workflow, revision.Revision)

	wr.Output(Config)

	Dialog.Infof(`Successfully updated workflow %v
			
Updated configuration is visible at %v`,
		workflowName,
		format.GuiLink(Config, "/workflow/%v", workflowName),
	)

	return nil
}

func newReplaceWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replace [workflow name]",
		Short: "Replace the configuration of a relay workflow",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doReplaceWorkflow,
	}

	cmd.Flags().StringP("file", "f", "", "Path to relay workflow file.")

	return cmd
}

func doDeleteWorkflow(cmd *cobra.Command, args []string) error {
	workflowName, nerr := getWorkflowName(args)

	if nerr != nil {
		return nerr
	}

	proceed, cerr := util.Confirm("Are you sure you want to delete this workflow?", Config)

	if cerr != nil {
		return cerr
	}

	if !proceed {
		return nil
	}

	Dialog.Progress("Deleting workflow...")

	_, err := Client.DeleteWorkflow(workflowName)

	if err != nil {
		return err
	}

	Dialog.Info("Workflow successfully deleted")

	return nil
}

func newDeleteWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [workflow name]",
		Short: "Delete a relay workflow",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doDeleteWorkflow,
	}

	return cmd
}

func parseParameter(str string) (key, value string) {
	idx := strings.Index(str, "=")

	// TODO: This behavior will basically silently discard parameters that are
	// in the wrong format. Should this, like, panic perhaps? Or notify the user that their
	// parameters are specified incorrectly?
	if idx < 0 {
		return
	}

	key = str[0:idx]
	value = str[idx+1:]
	return
}

func parseParameters(strs []string) map[string]string {
	res := make(map[string]string)

	for _, str := range strs {
		key, val := parseParameter(str)

		// value of empty string could, indeed, be a valid parameter.
		if key != "" {
			res[key] = val
		}
	}

	return res
}

func doRunWorkflow(cmd *cobra.Command, args []string) error {
	params, err := cmd.Flags().GetStringArray("parameter")

	if err != nil {
		panic(err)
	}

	// TODO: Same here as above. Could really DRY all this up.
	name, err := getWorkflowName(args)

	if err != nil {
		return err
	}

	Dialog.Progress("Starting your workflow...")

	resp, err := Client.RunWorkflow(name, parseParameters(params))

	if err != nil {
		// TODO: This error should be translated for the user. Right now it just
		// says whatever the service says.
		return err
	}

	link := format.GuiLink(Config, "/workflows/ec2-reaper/runs/%d/graph", resp.Run.RunNumber)
	Dialog.Info(fmt.Sprintf("Your run has started. Monitor it's progress: %s", link))

	return nil
}

func newRunWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [workflow name]",
		Short: "Invoke a relay workflow",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doRunWorkflow,
	}

	return cmd
}

func doListWorkflows(cmd *cobra.Command, args []string) error {
	resp, err := Client.ListWorkflows()

	if err != nil {
		panic(err)
	}

	t := Dialog.Table()

	t.Headers([]string{"Name", "Last Run Number"})

	for _, workflow := range resp.Workflows {
		t.AppendRow([]string{workflow.Name, fmt.Sprintf("%d", workflow.MostRecentRun.RunNumber)})
	}

	t.Flush()

	return nil
}

func newListWorkflowsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Get a list of all your workflows",
		Args:  cobra.MaximumNArgs(0),
		RunE:  doListWorkflows,
	}

	return cmd
}

// getWorkflowName gets the name of the workflow either from arguments or, if
// none are supplied, reads it from stdin.
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
