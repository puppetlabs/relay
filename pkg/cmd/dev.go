package cmd

import (
	"path/filepath"

	"github.com/puppetlabs/relay/pkg/cluster"
	"github.com/puppetlabs/relay/pkg/dev"
	"github.com/spf13/cobra"
)

func newDevCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
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

	return cmd
}

func doDevWorkflowRun(cmd *cobra.Command, args []string) error {
	file, ferr := readFile(cmd)
	if ferr != nil {
		return ferr
	}

	ctx := cmd.Context()
	cm := cluster.NewManager()

	cl, err := cm.GetClient(ctx, cluster.ClientOptions{Scheme: dev.DefaultScheme})
	if err != nil {
		return err
	}

	dm := dev.NewManager(cm, cl, dev.Options{DataDir: filepath.Join(Config.DataDir, "dev")})

	Dialog.Info("Running workflow")

	return dm.RunWorkflow(ctx, []byte(file))
}
