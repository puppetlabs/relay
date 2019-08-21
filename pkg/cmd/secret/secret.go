package secret

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/puppetlabs/nebula-cli/pkg/client"
	"github.com/puppetlabs/nebula-cli/pkg/config/runtimefactory"
	"github.com/puppetlabs/nebula-cli/pkg/errors"
	"github.com/puppetlabs/nebula-cli/pkg/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func NewCommand(rt runtimefactory.RuntimeFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "secret",
		Short:                 "Manage secrets",
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(NewSetCommand(rt))

	return cmd
}

func NewSetCommand(rt runtimefactory.RuntimeFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "set",
		Short:                 "Set the given secret value",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := rt.Config()
			if err != nil {
				return err
			}

			workflow, err := cmd.Flags().GetString("workflow")
			if err != nil {
				return err
			}

			if workflow == "" {
				return errors.NewWorkflowCliFlagError("--workflow", "required")
			}

			key, err := cmd.Flags().GetString("key")
			if err != nil {
				return err
			}

			if key == "" {
				return errors.NewWorkflowCliFlagError("--key", "required")
			}

			var value string

			valueFromStdin, err := cmd.Flags().GetBool("value-stdin")
			if err != nil {
				return err
			}

			if valueFromStdin {
				gotStdin, err := util.PassedStdin()
				if err != nil {
					return err
				}

				if gotStdin {
					buf := bytes.Buffer{}

					_, err := buf.ReadFrom(os.Stdin)
					if err != nil && err != io.EOF {
						return err
					}

					value = buf.String()
				} else {
					return errors.NewWorkflowSecretValueNotSpecifiedError()
				}
			} else {
				fmt.Print("Value: ")
				valueBytes, err := terminal.ReadPassword(int(syscall.Stdin))
				if err != nil {
					return err
				}

				value = string(valueBytes)
				fmt.Println()
			}

			client, err := client.NewAPIClient(cfg)
			if err != nil {
				return err
			}

			// trim all prefixed and suffixed whitespace
			workflow = strings.TrimSpace(workflow)
			key = strings.TrimSpace(key)

			if _, err := client.CreateWorkflowSecret(context.Background(), workflow, key, value); errors.IsClientWorkflowSecretAlreadyExistsError(err) {
				if _, err := client.UpdateWorkflowSecret(context.Background(), workflow, key, value); err != nil {
					return err
				}
			} else if err != nil {
				return err
			}

			fmt.Fprintln(rt.IO().Out, "Successfully set secret")

			return nil
		},
	}

	cmd.Flags().StringP("workflow", "w", "", "the workflow name")
	cmd.Flags().StringP("key", "k", "", "the secret key")
	cmd.Flags().BoolP("value-stdin", "v", false, "accept the value from stdin")

	return cmd
}
