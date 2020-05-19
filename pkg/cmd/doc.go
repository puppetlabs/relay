package cmd

import (
	"bytes"

	"github.com/spf13/cobra"
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
		Short: "Generate markdown documentation to stdout",
		Args:  cobra.NoArgs,
		RunE:  genDocs,
	}

	return cmd
}

// genOverviewMarkdown makes a single-page 'man' style document
// Much of this is copypasta from doc.GenMarkdownCustom, because
// cobra/doc doesn't provide real formatting customization
func genOverviewMarkdown() (md string, err error) {
	buf := new(bytes.Buffer)
	cmd := getCmd()

	// this is from doc.GenMarkdownCustom
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	name := cmd.CommandPath()

	short := cmd.Short
	long := cmd.Long

	buf.WriteString("## " + name + "\n\n" + short + "\n\n")
	buf.WriteString("### Synopsis\n\n" + long + "\n\n")
	buf.WriteString("### Subcommand Usage\n\n")

	children := cmd.Commands()

	if err := testChildren(children, buf); err != nil {
		return buf.String(), err
	}

	buf.WriteString("### Global flags\n```\n")
	flags := cmd.PersistentFlags()
	buf.WriteString(flags.FlagUsages() + "\n```\n")

	markdown := buf.String()

	return markdown, err

}

// testChildren determines whether this command ought to be documented.
// For brevity, we only want to generate docs for 'leaf' commands, i.e.
// only "relay workflow add", not "relay workflow"
func testChildren(children []*cobra.Command, buf *bytes.Buffer) error {

	for _, child := range children {
		if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := genChildMarkdown(child, buf); err != nil {
			return err
		}
	}

	return nil

}

func genChildMarkdown(cmd *cobra.Command, buf *bytes.Buffer) error {
	if cmd.Runnable() {
		usage := cmd.UseLine()
		buf.WriteString("**`" + usage + "`** -- " + cmd.Short + "\n")
		long := cmd.Long
		if len(long) > 0 {
			buf.WriteString("  " + cmd.Long + "\n")
		}
		flags := cmd.NonInheritedFlags()
		if flags.HasAvailableFlags() {
			buf.WriteString("```\n")
			flags.SetOutput(buf)
			flags.PrintDefaults()
			buf.WriteString("```\n")
		}
		buf.WriteString("\n")
	}

	// Because commands can be nested arbitrarily deep, this recurses into
	// the current command's children and tests them for runnability
	children := cmd.Commands()
	testChildren(children, buf)

	return nil

}

func genDocs(cmd *cobra.Command, args []string) error {

	markdown, err := genOverviewMarkdown()
	if err != nil {
		Dialog.Errorf("problem generating markdown: %s", err.Error)
	}
	Dialog.WriteString(markdown)

	return nil

}
