package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/puppetlabs/relay/pkg/debug"
	"github.com/puppetlabs/relay/pkg/dev"
	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/spf13/cobra"
)

func newMetadataCommand() *cobra.Command {

	// TODO Add help about usage, idling and direct execution
	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "Run a mock metadata service for entrypoint testing",
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		RunE:               doRunMetadata,
	}

	cmd.Flags().StringP("input", "i", "", "Path to metadata mock file")
	cmd.MarkFlagRequired("input")

	cmd.Flags().StringP("run", "r", "1", "Run ID of step to serve")

	cmd.Flags().StringP("step", "s", "default", "Step name to serve")

	return cmd
}

func doRunMetadata(cmd *cobra.Command, subcommand []string) error {
	input, err := cmd.Flags().GetString("input")
	if err != nil {
		debug.Log("The input flag is missing on the Cobra command configuration")
		return errors.NewGeneralUnknownError().WithCause(err).Bug()
	}
	runID, err := cmd.Flags().GetString("run")
	if err != nil {
		debug.Log("The run flag is missing on the Cobra command configuration")
		return errors.NewGeneralUnknownError().WithCause(err).Bug()
	}
	stepName, err := cmd.Flags().GetString("step")
	if err != nil {
		debug.Log("The step flag is missing on the Cobra command configuration")
		return errors.NewGeneralUnknownError().WithCause(err).Bug()
	}

	ctx := cmd.Context()
	m := dev.NewMetadataAPIManager(DevConfig)

	url, err := m.InitializeMetadataApi(ctx, dev.MetadataMockOptions{
		RunID: runID,
		StepName: stepName,
		Input: input,
	})
	if err != nil {
		return err
	}

	if len(subcommand) == 0 {
		Dialog.Infof("No command was supplied, idling metadataservice with access via: %s", url)
		for { }
	} else {
		// TODO this may not be strictly correct, if the subcommand is fancy
		command := exec.Command(subcommand[0], subcommand[1:]...)
		command.Env = os.Environ()
		command.Env = append(command.Env, fmt.Sprintf("METADATA_API_URL=%s", url))
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		err = command.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
