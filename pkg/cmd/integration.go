package cmd

import (
	stderrs "errors"
	"io/ioutil"
	"path/filepath"

	"github.com/puppetlabs/relay/pkg/debug"
	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/puppetlabs/relay/pkg/integration"

	"github.com/spf13/cobra"
)

var ErrIntegrationFileNotFound = stderrs.New("cmd: integration file not found")

func newIntegrationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "integration",
		Short: "Manage Relay integrations",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(newBuildIntegrationCommand())

	return cmd
}

func findIntegrationConfig(path string) (string, error) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return "", err
	}

	// Let's see if we can find the integration.yaml (or integration.yml) file
	// here.
	for _, info := range files {
		if info.Name() == "integration.yaml" || info.Name() == "integration.yml " {
			return filepath.Join(path, info.Name()), nil
		}
	}

	return "", ErrIntegrationFileNotFound
}

func doBuildIntegration(cmd *cobra.Command, args []string) error {
	// This is the path to the integration that we want to build.
	path := pop(args)

	// If this points at nothing then we're going to assume that we're building
	// the current path.
	if path == "" {
		path = "."
	}

	file, err := findIntegrationConfig(path)

	if err != nil {
		if err == ErrIntegrationFileNotFound {
			return errors.NewIntegrationFileNotFound()
		}

		debug.Logf("an error occured while searchin `%s` for integration file: %v", path, err)
		return errors.NewGeneralUnknownError().WithCause(err)
	}

	if err := integration.Build(file); err != nil {
		debug.Logf("an error occured while building `%s`: %v", file, err)
		return errors.NewGeneralUnknownError().WithCause(err)
	}

	return nil
}

func newBuildIntegrationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [integration path]",
		Short: "Build an entire integration",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doBuildIntegration,
	}

	return cmd
}

func pop(args []string) string {
	if len(args) > 0 {
		return args[len(args)-1]
	}

	return ""
}
