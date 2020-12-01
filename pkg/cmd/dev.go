package cmd

import (
	"os"

	"github.com/puppetlabs/horsehead/v2/workdir"
	"github.com/puppetlabs/relay/pkg/cluster"
	"github.com/puppetlabs/relay/pkg/dev"
	"github.com/spf13/cobra"
)

var DevConfig = dev.Config{}
var ClusterConfig = cluster.Config{}

func newDevCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "dev",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()

			err := root.PersistentPreRunE(cmd, args)
			if err != nil {
				return err
			}

			datadir, err := workdir.NewNamespace([]string{"relay", "dev"}).New(workdir.DirTypeData, workdir.Options{})
			if err != nil {
				return err
			}

			DevConfig = dev.Config{
				WorkDir: datadir,
				Dialog:  Dialog,
			}

			ClusterConfig = cluster.Config{
				WorkDir: datadir,
				Dialog:  Dialog,
			}

			return nil
		},
		Short: "Manage the local development environment",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(newClusterCommand())
	cmd.AddCommand(newImageCommand())
	cmd.AddCommand(newKubectlCommand())
	cmd.AddCommand(newMetadataCommand())

	// TODO temporary workflow commands until `relay workflow` is integrated
	// with the dev cluster
	cmd.AddCommand(newDevWorkflowCommand())

	return cmd
}

// TODO the commands below are essentially duplicates of the primary workflow
// and secret commands. These will eventually be merged with the main commands
// after the experimental phase.

func newDevWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Run Workflow commands against the dev cluster",
	}

	cmd.AddCommand(newDevWorkflowRunCommand())
	cmd.AddCommand(newDevWorkflowSecretCommand())

	return cmd
}

func newDevWorkflowRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a workflow on the dev cluster",
		RunE:  doDevWorkflowRun,
	}

	cmd.Flags().StringP("file", "f", "", "Path to Relay workflow file")
	cmd.MarkFlagRequired("file")

	cmd.Flags().StringArrayP("parameter", "p", []string{}, "Parameters to invoke this workflow run with")

	return cmd
}

func doDevWorkflowRun(cmd *cobra.Command, args []string) error {
	fp, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	file, err := os.Open(fp)
	if err != nil {
		return err
	}

	params, err := cmd.Flags().GetStringArray("parameter")
	if err != nil {
		return err
	}

	ctx := cmd.Context()
	cm := cluster.NewManager(ClusterConfig)

	cl, err := cm.GetClient(ctx, cluster.ClientOptions{Scheme: dev.DefaultScheme})
	if err != nil {
		return err
	}

	dm := dev.NewManager(cm, cl, DevConfig)

	Dialog.Info("Running workflow")

	ws, err := dm.RunWorkflow(ctx, file, parseParameters(params))
	if err != nil {
		return err
	}

	Dialog.Infof("Monitor step progress with: relay dev kubectl -n %s get pods --watch", ws.WorkflowIdentifier.Name)

	return nil
}

func newDevWorkflowSecretCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "Manage workflow secrets",
	}

	cmd.AddCommand(newDevWorkflowSecretSetCommand())

	return cmd
}

func newDevWorkflowSecretSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set a workflow secret",
		Args:  cobra.MaximumNArgs(2),
		RunE:  doDevWorkflowSecretSet,
	}

	cmd.Flags().Bool("value-stdin", false, "accept secret value from stdin")

	return cmd
}

func doDevWorkflowSecretSet(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	cm := cluster.NewManager(ClusterConfig)

	cl, err := cm.GetClient(ctx, cluster.ClientOptions{Scheme: dev.DefaultScheme})
	if err != nil {
		return err
	}

	dm := dev.NewManager(cm, cl, DevConfig)

	sc, err := getSecretValues(cmd, args)
	if err != nil {
		return err
	}

	Dialog.Info("Setting your secret...")

	return dm.SetWorkflowSecret(ctx, sc.workflowName, sc.name, sc.value)
}
