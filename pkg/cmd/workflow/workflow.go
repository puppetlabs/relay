package workflow

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/puppetlabs/nebula/pkg/client"
	"github.com/puppetlabs/nebula/pkg/config/runtimefactory"
	"github.com/puppetlabs/nebula/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCommand(rt runtimefactory.RuntimeFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "workflow",
		Short:                 "Manage nebula workflows",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewListCommand(rt))
	cmd.AddCommand(NewCreateCommand(rt))
	cmd.AddCommand(NewRunCommand(rt))
	cmd.AddCommand(NewListRunsCommand(rt))

	return cmd
}

func NewListCommand(rt runtimefactory.RuntimeFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		Short:                 "List workflows",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := client.NewAPIClient(rt.Config())
			if err != nil {
				return err
			}

			index, err := client.ListWorkflows(context.Background())
			if err != nil {
				return err
			}

			tw := table.NewWriter()

			tw.AppendHeader(table.Row{"NAME", "WORKFLOW"})
			for _, wf := range index.Items {
				p := []string{*wf.Repository, *wf.Branch, *wf.Path}

				tw.AppendRow(table.Row{wf.Name, strings.Join(p, "/")})
			}

			fmt.Fprintf(rt.IO().Out, "%s\n", tw.Render())

			return nil
		},
	}

	return cmd
}

func NewCreateCommand(rt runtimefactory.RuntimeFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create",
		Short:                 "Create workflows",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			if name == "" {
				return errors.NewWorkflowCliFlagError("--name", "required")
			}

			description, err := cmd.Flags().GetString("description")
			if err != nil {
				return err
			}

			repo, err := cmd.Flags().GetString("repository")
			if err != nil {
				return err
			}

			if repo == "" {
				return errors.NewWorkflowCliFlagError("--repository", "required")
			}

			branch, err := cmd.Flags().GetString("branch")
			if err != nil {
				return err
			}

			if branch == "" {
				return errors.NewWorkflowCliFlagError("--branch", "required")
			}

			path, err := cmd.Flags().GetString("filepath")
			if err != nil {
				return err
			}

			if path == "" {
				return errors.NewWorkflowCliFlagError("--filepath", "required")
			}

			client, err := client.NewAPIClient(rt.Config())
			if err != nil {
				return err
			}

			if _, err = client.CreateWorkflow(context.Background(), name, description, repo, branch, path); err != nil {
				return err
			}

			fmt.Fprintln(rt.IO().Out, "Success")

			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "workflow name")
	cmd.Flags().StringP("description", "d", "", "workflow description")
	cmd.Flags().StringP("repository", "r", "", "name of the repository")
	cmd.Flags().StringP("branch", "b", "", "name of the branch")
	cmd.Flags().StringP("filepath", "f", "", "path to the workflow file")

	return cmd
}

func NewRunCommand(rt runtimefactory.RuntimeFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "run",
		Short:                 "Run workflows",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			timeout, err := cmd.Flags().GetDuration("timeout")
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			if name == "" {
				return errors.NewWorkflowCliFlagError("--name", "required")
			}

			client, err := client.NewAPIClient(rt.Config())
			if err != nil {
				return err
			}

			run, err := client.RunWorkflow(ctx, name)
			if err != nil {
				return err
			}

			// TODO: temporary defaults until the API fills out the values
			if run.RunNumber == nil {
				num := int64(1)
				run.RunNumber = &num
			}

			if run.Status == nil {
				status := "pending"
				run.Status = &status
			}

			tw := table.NewWriter()

			tw.AppendHeader(table.Row{"#", "ID", "STATUS"})
			tw.AppendRow(table.Row{fmt.Sprintf("%d", *run.RunNumber), *run.ID, *run.Status})

			fmt.Fprintf(rt.IO().Out, "%s\n", tw.Render())

			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "the workflow name to run against")
	cmd.Flags().DurationP("timeout", "t", 10*time.Minute, "the timeout for a workflow run to start")

	return cmd
}

func NewListRunsCommand(rt runtimefactory.RuntimeFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list-runs",
		Short:                 "List workflow runs",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			if name == "" {
				return errors.NewWorkflowCliFlagError("--name", "required")
			}

			client, err := client.NewAPIClient(rt.Config())
			if err != nil {
				return err
			}

			wrs, err := client.ListWorkflowRuns(context.Background(), name)
			if err != nil {
				return err
			}

			tw := table.NewWriter()
			tw.AppendHeader(table.Row{"#", "ID", "STATUS"})

			for _, run := range wrs.Items {
				// TODO: temporary defaults until the API fills out the values
				if run.RunNumber == nil {
					num := int64(1)
					run.RunNumber = &num
				}

				if run.Status == nil {
					status := "pending"
					run.Status = &status
				}

				tw.AppendRow(table.Row{fmt.Sprintf("%d", *run.RunNumber), *run.ID, *run.Status})
			}

			fmt.Fprintf(rt.IO().Out, "%s\n", tw.Render())

			return nil
		},
	}

	cmd.Flags().StringP("filepath", "f", "", "path to a workflow.yaml")
	cmd.Flags().StringP("name", "n", "", "the workflow name to run against")

	return cmd
}
