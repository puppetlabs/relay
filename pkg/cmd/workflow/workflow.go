package workflow

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

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage your relay workflows",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(NewAddWorkflowCommand())
	cmd.AddCommand(NewReplaceWorkflowCommand())
	cmd.AddCommand(NewDeleteWorkflowCommand())
	cmd.AddCommand(NewRunWorkflowCommand())
	cmd.AddCommand(NewListWorkflowsCommand())

	return cmd
}

func doAddWorkflow(cmd *cobra.Command, args []string) error {
	cfg, cfgerr := config.FromFlags(cmd.Flags())

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

	log := dialog.FromConfig(cfg)

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

func NewAddWorkflowCommand() *cobra.Command {
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
	cfg, cfgerr := config.FromFlags(cmd.Flags())

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

	log := dialog.FromConfig(cfg)

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

func NewReplaceWorkflowCommand() *cobra.Command {
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
	cfg, cfgerr := config.FromFlags(cmd.Flags())

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

	log := dialog.FromConfig(cfg)

	log.Info("Deleting workflow...")

	client := client.NewClient(cfg)

	_, err := client.DeleteWorkflow(workflowName)

	if err != nil {
		return err
	}

	log.Info("Workflow successfully deleted")

	return nil
}

func NewDeleteWorkflowCommand() *cobra.Command {
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
	cfg, err := config.FromFlags(cmd.Flags())

	if err != nil {
		return err
	}

	params, err := cmd.Flags().GetStringArray("parameter")

	if err != nil {
		panic(err)
	}

	// TODO: Same here as above. Could really DRY all this up.
	name, err := getWorkflowName(args)

	if err != nil {
		return err
	}

	log := dialog.FromConfig(cfg)

	log.Info("Starting your workflow...")

	client := client.NewClient(cfg)

	resp, err := client.RunWorkflow(name, parseParameters(params))

	if err != nil {
		// TODO: This error should be translated for the user. Right now it just
		// says whatever the service says.
		return err
	}

	link := format.GuiLink(cfg, "/workflows/ec2-reaper/runs/%d/graph", resp.Run.RunNumber)
	log.Info(fmt.Sprintf("Your run has started. Monitor it's progress: %s", link))

	return nil
}

func NewRunWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [workflow name]",
		Short: "Invoke a relay workflow",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doRunWorkflow,
	}

	return cmd
}

func doListWorkflows(cmd *cobra.Command, args []string) error {
	cfg, err := config.FromFlags(cmd.Flags())

	if err != nil {
		return err
	}

	client := client.NewClient(cfg)

	resp, err := client.ListWorkflows()

	if err != nil {
		panic(err)
	}

	log := dialog.FromConfig(cfg)
	t := log.Table()

	t.Headers([]string{"Name", "Last Run Number"})

	for _, workflow := range resp.Workflows {
		t.AppendRow([]string{workflow.Name, fmt.Sprintf("%d", workflow.MostRecentRun.RunNumber)})
	}

	t.WriteTo(os.Stdout)

	return nil
}

func NewListWorkflowsCommand() *cobra.Command {
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
