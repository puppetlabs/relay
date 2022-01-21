package cmd

import (
	"context"
	"os"

	"github.com/puppetlabs/leg/workdir"
	"github.com/puppetlabs/relay/pkg/cluster"
	"github.com/puppetlabs/relay/pkg/config"
	"github.com/puppetlabs/relay/pkg/dev"
	"github.com/spf13/cobra"
)

const (
	InstallHelmControllerFlag = "install-helm-controller"
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

	cmd.Flags().BoolP(InstallHelmControllerFlag, "", false, "Optional installation of Helm Controller")

	return cmd
}

func doInitDevelopmentEnvironment(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	installHelmController, err := cmd.Flags().GetBool(InstallHelmControllerFlag)
	if err != nil {
		return err
	}

	opts := dev.InitializeOptions{
		InstallHelmController: installHelmController,
	}

	return initDevelopmentEnvironment(ctx, opts)
}

func initDevelopmentEnvironment(ctx context.Context, opts dev.InitializeOptions) error {
	dm, err := NewManager(ctx, DevConfig)
	if err != nil {
		return err
	}

	Dialog.Progress("Initializing cluster; this may take several minutes...")
	if err := dm.Initialize(ctx, opts); err != nil {
		return err
	}

	installerOpts := mapInstallerOptionsFromConfig(Config.InstallerConfig,
		dev.InstallerOptions{
			InstallerImage:         dev.RelayInstallerImage,
			LogServiceImage:        dev.RelayLogServiceImage,
			MetadataAPIImage:       dev.RelayMetadataAPIImage,
			OperatorImage:          dev.RelayOperatorImage,
			OperatorVaultInitImage: dev.RelayOperatorVaultInitImage,
			OperatorWebhookCertificateControllerImage: dev.RelayOperatorWebhookCertificateControllerImage,
		})

	logServiceOpts := mapLogServiceOptionsFromConfig(Config.LogServiceConfig)

	Dialog.Info("Initializing relay-core; this may take several minutes...")

	if err := dm.InitializeRelayCore(ctx, installerOpts, logServiceOpts); err != nil {
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
	dm, err := NewManager(ctx, DevConfig)
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

	dm, err := NewManager(ctx, DevConfig)
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

func NewManager(ctx context.Context, devConfig dev.Config) (*dev.Manager, error) {
	cm := cluster.NewManager(
		cluster.Config{
			WorkDir: devConfig.WorkDir,
		})

	if ok, _ := cm.Exists(ctx); ok {
		return dev.NewManagerFromLocalCluster(ctx, cm, DevConfig)
	}

	return dev.NewManagerFromExternalCluster(ctx)
}

func mapInstallerOptionsFromConfig(installerConfig *config.InstallerConfig, defaultInstallerOpts dev.InstallerOptions) dev.InstallerOptions {
	installerOpts := defaultInstallerOpts
	if Config.InstallerConfig != nil {
		installerOpts.InstallerImage = coalesce(installerConfig.InstallerImage, defaultInstallerOpts.InstallerImage)
		installerOpts.LogServiceImage = coalesce(installerConfig.LogServiceImage, defaultInstallerOpts.LogServiceImage)
		installerOpts.MetadataAPIImage = coalesce(installerConfig.MetadataAPIImage, defaultInstallerOpts.MetadataAPIImage)
		installerOpts.OperatorImage = coalesce(installerConfig.OperatorImage, defaultInstallerOpts.OperatorImage)
		installerOpts.OperatorVaultInitImage = coalesce(installerConfig.OperatorVaultInitImage, defaultInstallerOpts.OperatorVaultInitImage)
		installerOpts.OperatorWebhookCertificateControllerImage = coalesce(installerConfig.OperatorWebhookCertificateControllerImage, defaultInstallerOpts.OperatorWebhookCertificateControllerImage)
	}

	return installerOpts
}

func mapLogServiceOptionsFromConfig(logServiceConfig *config.LogServiceConfig) dev.LogServiceOptions {
	logServiceOpts := dev.LogServiceOptions{}
	if logServiceConfig != nil {
		logServiceOpts = dev.LogServiceOptions{
			Enabled:               true,
			CredentialsKey:        logServiceConfig.CredentialsKey,
			CredentialsSecretName: logServiceConfig.CredentialsSecretName,
			Project:               logServiceConfig.Project,
			Dataset:               logServiceConfig.Dataset,
			Table:                 logServiceConfig.Table,
		}
	}

	return logServiceOpts
}

func coalesce(target string, other string) string {
	if target != "" {
		return target
	}

	return other
}
