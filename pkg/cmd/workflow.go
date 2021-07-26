package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/puppetlabs/relay/pkg/debug"
	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/puppetlabs/relay/pkg/format"
	"github.com/puppetlabs/relay/pkg/util"
	"github.com/spf13/cobra"
)

func newWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage your Relay workflows",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(newSaveWorkflowCommand())
	cmd.AddCommand(newValidateWorkflowFileCommand())
	cmd.AddCommand(newDeleteWorkflowCommand())
	cmd.AddCommand(newRunWorkflowCommand())
	cmd.AddCommand(newListWorkflowsCommand())
	cmd.AddCommand(newDownloadWorkflowCommand())
	cmd.AddCommand(newSecretCommand())

	// Deprecated
	cmd.AddCommand(newAddWorkflowCommand())
	cmd.AddCommand(newReplaceWorkflowCommand())

	return cmd
}

func newSaveWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "save [workflow name]",
		Short: "Save a Relay workflow",
		Args:  cobra.MaximumNArgs(3),
		RunE:  doSaveWorkflow,
	}

	cmd.Flags().StringP("file", "f", "", "Path to Relay workflow file")
	cmd.Flags().BoolP("no-overwrite", "O", false, "Abort instead of overwriting existing workflow")
	cmd.Flags().BoolP("no-create", "C", false, "Abort instead of creating a workflow that does not exist")

	return cmd
}

func newAddWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:        "add [workflow name]",
		Short:      "Add a Relay workflow from a local file",
		Args:       cobra.MaximumNArgs(1),
		RunE:       doSaveWorkflow,
		Deprecated: "Use `save` instead",
	}

	cmd.Flags().StringP("file", "f", "", "Path to Relay workflow file")

	return cmd
}

func newReplaceWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:        "replace [workflow name]",
		Short:      "Replace an existing Relay workflow",
		Args:       cobra.MaximumNArgs(1),
		RunE:       doSaveWorkflow,
		Deprecated: "Use `save` instead",
	}

	cmd.Flags().StringP("file", "f", "", "Path to Relay workflow file")

	return cmd
}

func doValidateWorkflowFile(cmd *cobra.Command, args []string) error {
	filepath, file, err := readFile(cmd)

	if err != nil {
		return err
	}

	Dialog.Info("Validating workflow file " + filepath)

	_, rerr := Client.Validate(file)

	if rerr != nil {
		return rerr
	}

	Dialog.Infof(`Successfully validated workflow file %v`, filepath)

	return nil
}

func newValidateWorkflowFileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a local Relay workflow file",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doValidateWorkflowFile,
	}

	cmd.Flags().StringP("file", "f", "", "Path to Relay workflow file")

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
		Short: "Delete a Relay workflow",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doDeleteWorkflow,
	}

	return cmd
}

func parseParameter(str string) (key, value string) {
	strs := strings.SplitN(str, "=", 2)

	if len(strs) == 2 {
		return strs[0], strs[1]
	}

	debug.Logf("invalid parameter: %s", str)
	return "", ""
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
		debug.Log("The parameters flag is missing on the Cobra command configuration")
		return errors.NewGeneralUnknownError().WithCause(err).Bug()
	}

	// TODO: Same here as above. Could really DRY all this up.
	name, err := getWorkflowName(args)

	if err != nil {
		return err
	}

	Dialog.Progress("Starting your workflow...")

	resp, err := Client.RunWorkflow(name, parseParameters(params))

	if err != nil {
		return err
	}

	link := format.GuiLink(Config, "/workflows/%s/runs/%d/graph", name, resp.Run.RunNumber)
	Dialog.Info(fmt.Sprintf("Your run has started. Monitor its progress here: %s", link))

	return nil
}

func newRunWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "run [workflow name]",
		Short:             "Invoke a Relay workflow",
		Args:              cobra.MaximumNArgs(1),
		RunE:              doRunWorkflow,
		ValidArgsFunction: doListWorkflowsCompletion,
	}

	cmd.Flags().StringArrayP("parameter", "p", []string{}, "Parameters to invoke this workflow run with")

	return cmd
}

func doDownloadWorkflow(cmd *cobra.Command, args []string) error {
	name, err := getWorkflowName(args)

	if err != nil {
		return err
	}

	body, err := Client.DownloadWorkflow(name)

	if err != nil {
		if errors.IsClientResponseNotFound(err) {
			Dialog.Warnf(`No file data found for workflow %v

View more information or update workflow settings at: %v`,
				name,
				format.GuiLink(Config, "/workflows/%v", name),
			)

			return nil
		}
		return err
	}

	filepath, ferr := cmd.Flags().GetString("file")
	if ferr != nil {
		return ferr
	}

	if filepath == "" {
		Dialog.WriteString(body)
	} else {
		if err := ioutil.WriteFile(filepath, []byte(body), 0644); err != nil {
			debug.Logf("failed to write to file %s: %s", filepath, err.Error())
			return err
		}
	}

	return nil
}

func newDownloadWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download [workflow name]",
		Short: "Download a workflow from the service",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doDownloadWorkflow,
	}

	cmd.Flags().StringP("file", "f", "", "Path to write workflow file")

	return cmd
}

func doListWorkflows(cmd *cobra.Command, args []string) error {
	resp, err := Client.ListWorkflows()

	if err != nil {
		debug.Logf("failed to list workflows: %s", err.Error())
		return err
	}

	t := Dialog.Table()

	t.Headers([]string{"Name", "Last Run Number"})

	for _, workflow := range resp.Workflows {
		t.AppendRow([]string{workflow.Name, fmt.Sprintf("%d", workflow.MostRecentRun.RunNumber)})
	}

	t.Flush()

	return nil
}

func doListWorkflowsCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	resp, err := Client.ListWorkflows()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	results := []string{}

	for _, workflow := range resp.Workflows {
		if strings.HasPrefix(workflow.Name, toComplete) {
			results = append(results, workflow.Name)
		}
	}

	return results, cobra.ShellCompDirectiveDefault
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

func readFile(cmd *cobra.Command) (string, string, errors.Error) {
	filepath, err := cmd.Flags().GetString("file")

	if err != nil {
		return "", "", errors.NewGeneralUnknownError().WithCause(err)
	}

	if filepath == "" {
		return "", "", errors.NewWorkflowMissingFileFlagError()
	}

	bytes, err := ioutil.ReadFile(filepath)

	if err != nil {
		return "", "", errors.NewWorkflowWorkflowFileReadError().WithCause(err)
	}

	return filepath, string(bytes), nil
}
