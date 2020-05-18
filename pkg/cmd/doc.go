package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func newDocCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doc",
		Short: "generate docs for relay",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(newGenerateCommand())

	return cmd
}

func newGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate documentation",
		Args:  cobra.NoArgs,
		RunE:  genDocs,
	}

	cmd.Flags().StringP("format", "f", "markdown", "format of output docs: man or markdown")
	cmd.Flags().StringP("target", "t", "./docs/md", "target directory for output")

	return cmd
}

func genDocs(cmd *cobra.Command, args []string) error {
	docFormat, nil := cmd.Flags().GetString("format")
	targetDir, nil := cmd.Flags().GetString("target")

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		Dialog.Errorf(`Failed to create target directory %s: %s`, targetDir, err.Error())
		return err
	}

	rootcmd := getCmd()

	if docFormat == "markdown" {
		doc.GenMarkdownTree(rootcmd, targetDir)
		Dialog.Infof(`Generated markdown tree in %s`, targetDir)
	} else if docFormat == "man" {
		header := &doc.GenManHeader{
			Title:   "relay",
			Section: "1",
		}
		doc.GenManTree(rootcmd, header, targetDir)
		Dialog.Infof(`Generated man pages in %s`, targetDir)
	} else {
		Dialog.Errorf(`Unknown documentation format %s specified`, docFormat)
		return nil
	}

	return nil

}
