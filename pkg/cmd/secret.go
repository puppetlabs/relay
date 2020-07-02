package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/puppetlabs/relay/pkg/debug"
	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/puppetlabs/relay/pkg/format"
	"github.com/puppetlabs/relay/pkg/model"
	"github.com/puppetlabs/relay/pkg/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func newSecretCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "Manage your Relay secrets",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(newSetSecretCommand())
	cmd.AddCommand(newListSecretsCommand())
	cmd.AddCommand(newDeleteSecretCommand())

	return cmd
}

func newSetSecretCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [workflow-name] [secret name]",
		Short: "Set a Relay workflow secret",
		Args:  cobra.MaximumNArgs(2),
		RunE:  doSetSecret,
	}

	cmd.Flags().Bool("value-stdin", false, "accept secret value from stdin")

	return cmd
}

func doSetSecret(cmd *cobra.Command, args []string) error {
	workflowName, err := getWorkflowName(args)
	if err != nil {
		return err
	}
	secretName, err := getSecretName(args)
	if err != nil {
		return err
	}
	secretValue, err := getSecretValue(cmd)
	if err != nil {
		return err
	}

	Dialog.Progress("Setting your secret...")

	resp, err := Client.ListWorkflowSecrets(workflowName)
	if err != nil {
		debug.Logf("failed to list workflow secrets: %s", err.Error())
		return err
	}

	exists := func() bool {
		for i := range resp.WorkflowSecrets {
			if resp.WorkflowSecrets[i].Name == secretName {
				return true
			}
		}
		return false
	}()

	var secret *model.WorkflowSecretEntity
	if exists {
		secret, err = Client.UpdateWorkflowSecret(workflowName, secretName, secretValue)
		if err != nil {
			return err
		}
	} else {
		secret, err = Client.CreateWorkflowSecret(workflowName, secretName, secretValue)
		if err != nil {
			return err
		}
	}

	Dialog.Infof(`Successfully set secret %v on workflow %v

View more information or update workflow settings at: %v`,
		secret.Secret.Name,
		workflowName,
		format.GuiLink(Config, "/workflows/%v", workflowName),
	)

	return nil
}

func newDeleteSecretCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [workflow-name] [secret name]",
		Short: "Delete a Relay workflow",
		Args:  cobra.MaximumNArgs(2),
		RunE:  doDeleteSecret,
	}

	return cmd
}

func doDeleteSecret(cmd *cobra.Command, args []string) error {
	workflowName, err := getWorkflowName(args)
	if err != nil {
		return err
	}

	secretName, err := getSecretName(args)
	if err != nil {
		return err
	}

	proceed, err := util.Confirm("Are you sure you want to delete this secret?", Config)
	if err != nil {
		return err
	}
	if !proceed {
		return nil
	}

	Dialog.Progress("Deleting secret...")
	_, err = Client.DeleteWorkflowSecret(workflowName, secretName)
	if err != nil {
		return err
	}
	Dialog.Info("Secret successfully deleted")

	return nil
}

func newListSecretsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [workflow name]",
		Short: "List Relay workflow secrets",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doListSecrets,
	}

	return cmd
}

func doListSecrets(cmd *cobra.Command, args []string) error {
	workflowName, err := getWorkflowName(args)
	if err != nil {
		return err
	}

	resp, err := Client.ListWorkflowSecrets(workflowName)
	if err != nil {
		debug.Logf("failed to list workflow secrets: %s", err.Error())
		return err
	}

	t := Dialog.Table()

	t.Headers([]string{"Name"})

	for _, secret := range resp.WorkflowSecrets {
		t.AppendRow([]string{secret.Name})
	}

	t.Flush()

	return nil

}

// getSecretName gets the name of the secret from the second argument. If
// none are supplied, reads it from stdin.
func getSecretName(args []string) (string, errors.Error) {
	if len(args) > 1 {
		return args[1], nil
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Secret name: ")
	namePrompt, err := reader.ReadString('\n')
	if err != nil {
		return "", errors.NewSecretNameReadError().WithCause(err)
	}

	name := strings.TrimSpace(namePrompt)

	if name == "" {
		return "", errors.NewSecretMissingNameError()
	}

	return strings.TrimSpace(namePrompt), nil
}

// getSecretValue either prompts for the value of the secret with hidden input, or accepts the value from stdin if the
// --value-stdin boolean flag is set
func getSecretValue(cmd *cobra.Command) (string, errors.Error) {
	var value string

	valueFromStdin, err := cmd.Flags().GetBool("value-stdin")
	if err != nil {
		return "", errors.NewGeneralUnknownError().WithCause(err)
	}

	if valueFromStdin {
		gotStdin, err := util.PassedStdin()
		if err != nil {
			return "", errors.NewSecretFailedValueFromStdin().WithCause(err)
		}

		if gotStdin {
			buf := bytes.Buffer{}
			reader := &io.LimitedReader{R: os.Stdin, N: readLimit}

			n, err := buf.ReadFrom(reader)
			if err != nil && err != io.EOF {
				return "", errors.NewSecretFailedValueFromStdin().WithCause(err)
			}
			if n == 0 {
				return "", errors.NewSecretFailedNoStdin()
			}

			value = buf.String()
		} else {
			return "", errors.NewSecretFailedNoStdin()
		}
	} else {
		fmt.Print("Value: ")
		valueBytes, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			return "", errors.NewSecretFailedValueFromStdin().WithCause(err)
		}

		value = string(valueBytes)
		// resets to new line after hidden input
		fmt.Println("")
	}

	return value, nil
}
