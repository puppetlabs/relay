package cmd

import (
	"os"
	"path/filepath"

	"github.com/puppetlabs/relay/pkg/cluster"
	"github.com/puppetlabs/relay/pkg/dev"
	"github.com/spf13/cobra"
)

var DevConfig = dev.Config{}

func newDevCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "dev",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()

			err := root.PersistentPreRunE(cmd, args)
			if err != nil {
				return err
			}

			DevConfig = dev.Config{
				DataDir: filepath.Join(Config.DataDir, "dev"),
			}

			return nil
		},
		Short: "Manage the local development environment",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(newClusterCommand())
	cmd.AddCommand(newImageCommand())
	cmd.AddCommand(newKubectlCommand())

	// TODO temporary workflow run command until `relay workflow` is integrated
	// with the dev cluster
	cmd.AddCommand(newDevWorkflowRunCommand())

	return cmd
}

func newDevWorkflowRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow-run",
		Short: "Temporary workflow run command",
		RunE:  doDevWorkflowRun,
	}

	cmd.Flags().StringP("file", "f", "", "Path to Relay workflow file")
	cmd.MarkFlagRequired("file")

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

	ctx := cmd.Context()
	cm := cluster.NewManager(cluster.Config{DataDir: DevConfig.DataDir})

	cl, err := cm.GetClient(ctx, cluster.ClientOptions{Scheme: dev.DefaultScheme})
	if err != nil {
		return err
	}

	dm := dev.NewManager(cm, cl, DevConfig)

	Dialog.Info("Running workflow")

	ws, err := dm.RunWorkflow(ctx, file)
	if err != nil {
		return err
	}

	Dialog.Infof("Monitor step progress with: relay dev kubectl -n %s get pods --watch", ws.WorkflowIdentifier.Name)

	return nil
}
