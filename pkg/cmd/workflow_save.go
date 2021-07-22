package cmd

import (
	"fmt"

	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/puppetlabs/relay/pkg/format"
	"github.com/puppetlabs/relay/pkg/model"
	"github.com/spf13/cobra"
)

func doSaveWorkflow(cmd *cobra.Command, args []string) error {
	workflowName, err := getWorkflowName(args)
	if err != nil {
		return err
	}

	if err := checkFile(cmd); err != nil {
		return err
	}

	workflow, gerr := getOrCreateWorkflow(cmd, workflowName)
	if gerr != nil {
		return gerr
	}

	info, err := updateWorkflowRevision(cmd, workflow)
	if err != nil {
		return err
	}

	Dialog.Infof(`%s

View more information or update workflow settings at: %v`,
		info,
		format.GuiLink(Config, "/workflows/%v", workflow.Workflow.Name),
	)

	return nil
}

func checkFile(cmd *cobra.Command) error {
	_, _, err := readFile(cmd)
	if err != nil {
		return err
	}
	return nil
}

func getOrCreateWorkflow(cmd *cobra.Command, workflowName string) (*model.WorkflowEntity, error) {
	Dialog.Progress("Checking for workflow " + workflowName)

	workflow, err := Client.GetWorkflow(workflowName)
	if err != nil {
		if !errors.IsClientResponseNotFound(err) {
			return nil, err
		}

		if cmd.Name() == "replace" {
			return nil, errors.NewWorkflowDoesNotExistError()
		}

		if f := cmd.Flags().Lookup("no-create"); f != nil {
			if noCreate, err := cmd.Flags().GetBool("no-create"); err != nil {
				return nil, err
			} else if noCreate {
				return nil, errors.NewWorkflowDoesNotExistError()
			}
		}

		Dialog.Progress("Creating workflow " + workflowName)
		workflow, err = Client.CreateWorkflow(workflowName)
		if err != nil {
			return nil, err
		}
	} else {
		if cmd.Name() == "add" {
			return nil, errors.NewWorkflowAlreadyExistsError()
		}

		if f := cmd.Flags().Lookup("no-overwrite"); f != nil {
			if noOverwrite, err := cmd.Flags().GetBool("no-overwrite"); err != nil {
				return nil, err
			} else if noOverwrite {
				return nil, errors.NewWorkflowAlreadyExistsError()
			}
		}
	}

	return workflow, err
}

func updateWorkflowRevision(cmd *cobra.Command, workflow *model.WorkflowEntity) (string, errors.Error) {
	filePath, revisionContent, err := readFile(cmd)
	if err != nil {
		return "", err
	}

	info := fmt.Sprintf("Successfully saved workflow %v with file %s.", workflow.Workflow.Name, filePath)

	latestRevision, err := Client.GetLatestRevision(workflow.Workflow.Name)
	if err != nil && !errors.IsClientResponseNotFound(err) {
		return "", err
	}

	if latestRevision == nil || latestRevision.Revision.Raw != revisionContent {
		revision, err := Client.CreateRevision(workflow.Workflow.Name, revisionContent)
		if err != nil {
			Dialog.Warnf(`When uploading the file %s, we encountered the following errors:

	%s

	`,
				filePath,
				format.Error(err, cmd),
			)

			info = fmt.Sprintf("Attempted to save workflow %v, but the file content contained errors.", workflow.Workflow.Name)
		} else {
			wr := model.NewWorkflowRevision(workflow.Workflow, revision.Revision)
			wr.Output(Config)
		}
	}

	return info, nil
}
