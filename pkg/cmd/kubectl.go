package cmd

import (
	"path/filepath"

	"github.com/puppetlabs/relay/pkg/cluster"
	"github.com/puppetlabs/relay/pkg/dev"
	"github.com/spf13/cobra"
)

func newKubectlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubectl",
		Short: "",
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		DisableFlagParsing: true,
		RunE:               doRunKubectl,
	}

	return cmd
}

func doRunKubectl(cmd *cobra.Command, args []string) error {
	cm := cluster.NewManager()
	dm := dev.NewManager(cm, dev.Options{DataDir: filepath.Join(Config.DataDir, "dev")})

	newcmd, err := dm.KubectlCommand()
	if err != nil {
		return err
	}

	newcmd.SetArgs(args)

	return newcmd.Execute()
}