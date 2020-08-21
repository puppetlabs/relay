package cmd

import (
	"github.com/puppetlabs/relay/pkg/cluster"
	"github.com/puppetlabs/relay/pkg/dev"
	"github.com/spf13/cobra"
)

func newKubectlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubectl",
		Short: "Run kubectl commands against the dev cluster",
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		DisableFlagParsing: true,
		RunE:               doRunKubectl,
	}

	return cmd
}

func doRunKubectl(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cm := cluster.NewManager(cluster.Config{DataDir: DevConfig.DataDir})
	cl, err := cm.GetClient(ctx, cluster.ClientOptions{Scheme: dev.DefaultScheme})
	dm := dev.NewManager(cm, cl, DevConfig)

	newcmd, err := dm.KubectlCommand()
	if err != nil {
		return err
	}

	newcmd.SetArgs(args)

	return newcmd.Execute()
}
