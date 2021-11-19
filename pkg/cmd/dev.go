package cmd

import (
	"context"
	"os"

	"github.com/puppetlabs/leg/workdir"
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
			}

			ClusterConfig = cluster.Config{
				WorkDir: datadir,
			}

			return nil
		},
		Short: "Manage the local development environment",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(newClusterCommand())
	cmd.AddCommand(newInitializeCommand())
	cmd.AddCommand(newMetadataCommand())

	// TODO temporary workflow commands until `relay workflow` is integrated
	// with the dev cluster
	cmd.AddCommand(newDevWorkflowCommand())

	return cmd
}

func newInitializeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "initialize",
		Aliases: []string{"init"},
		Short:   "Initialize the Relay development environment",
		RunE:    doInitDevelopmentEnvironment,
	}

	return cmd
}

func doInitDevelopmentEnvironment(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	opts := cluster.InitializeOptions{}

	return initDevelopmentEnvironment(ctx, opts)
}

func initDevelopmentEnvironment(ctx context.Context, opts cluster.InitializeOptions) error {
	cm := cluster.NewManager(ClusterConfig)

	dm, err := dev.NewManager(ctx, cm, DevConfig)
	if err != nil {
		return err
	}

	logServiceOpts := dev.LogServiceOptions{}
	if Config.LogServiceConfig != nil {
		logServiceOpts = dev.LogServiceOptions{
			Enabled:               true,
			CredentialsSecretName: Config.LogServiceConfig.CredentialsSecretName,
			Project:               Config.LogServiceConfig.Project,
			Dataset:               Config.LogServiceConfig.Dataset,
			Table:                 Config.LogServiceConfig.Table,
		}
	}

	Dialog.Info("Initializing relay-core; this may take several minutes...")

	if err := dm.InitializeRelayCore(ctx, logServiceOpts); err != nil {
		return err
	}

	Dialog.Infof("Cluster connection can be used with: !Connection {type: kubernetes, name: %s}", dev.RelayClusterConnectionName)

	return nil
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

	dm, err := dev.NewManager(ctx, cm, DevConfig)
	if err != nil {
		return err
	}

	Dialog.Infof("Processing workflow file %s", fp)

	wd, err := dm.LoadWorkflow(ctx, file)
	if err != nil {
		return err
	}

	t, err := dm.CreateTenant(ctx, wd.Name)
	if err != nil {
		return err
	}

	wf, err := dm.CreateWorkflow(ctx, wd, t)
	if err != nil {
		return err
	}

	_, err = dm.RunWorkflow(ctx, wf, parseParameters(params))
	if err != nil {
		return err
	}

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
		Use:   "set [workflow name] [secret name]",
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

	dm, err := dev.NewManager(ctx, cm, DevConfig)
	if err != nil {
		return err
	}

	sc, err := getSecretValues(cmd, args)
	if err != nil {
		return err
	}

	Dialog.Infof("Setting secret %s for workflow %s", sc.name, sc.workflowName)

	return dm.SetWorkflowSecret(ctx, sc.workflowName, sc.name, sc.value)
}
