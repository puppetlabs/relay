package login

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/puppetlabs/nebula/pkg/client"
	"github.com/puppetlabs/nebula/pkg/config/runtimefactory"
	"github.com/puppetlabs/nebula/pkg/errors"
	"github.com/puppetlabs/nebula/pkg/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// readLimit is the max bytes allowed from stdin (512 bytes being well above a reasonable
// password length) to avoid reading massive byte arrays into memory and sending them
// to a remote API server.
const readLimit = 512

const defaultServerAddr = "https://api.nebula.puppet.com"

func NewCommand(rt runtimefactory.RuntimeFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "login",
		Short:                 "Authenticate with Nebula",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var reader *bufio.Reader
			var user, pass string

			passFromStdin, err := cmd.Flags().GetBool("password-stdin")
			if err != nil {
				return err
			}

			// if --password-stdin was set and we got a char device registered on
			// os.Stdin, then we can continue to try and read the password from the pipe.
			// Otherwise we will try to ask the user for their email and password via prompts.
			if passFromStdin {
				gotStdin, err := util.PassedStdin()
				if err != nil {
					return err
				}

				if gotStdin {
					buf := bytes.Buffer{}
					reader := &io.LimitedReader{R: os.Stdin, N: readLimit}

					_, err := buf.ReadFrom(reader)
					if err != nil && err != io.EOF {
						return err
					}

					pass = buf.String()

					user, err = cmd.Flags().GetString("email")
					if err != nil {
						return err
					}

					if user == "" {
						return errors.NewClientMissingEmailError("--password-stdin must be used with --email")
					}
				} else {
					return errors.NewClientPasswordError("did not get anything from stdin")
				}
			} else {
				reader = bufio.NewReader(os.Stdin)

				fmt.Print("Email: ")
				user, err = reader.ReadString('\n')
				if err != nil {
					return err
				}

				fmt.Print("Password: ")
				passBytes, err := terminal.ReadPassword(int(syscall.Stdin))
				if err != nil {
					return err
				}

				pass = string(passBytes)
				fmt.Println()
			}

			client, err := client.NewAPIClient(rt.Config())
			if err != nil {
				return err
			}

			// trim all prefixed and suffixed whitespace
			user = strings.TrimSpace(user)
			pass = strings.TrimSpace(pass)

			if err := client.Login(context.Background(), user, pass); err != nil {
				return err
			}

			fmt.Fprintln(rt.IO().Out, "Successfully logged in")

			return nil
		},
	}

	cmd.Flags().StringP("email", "e", "", "Nebula email")
	cmd.Flags().BoolP("password-stdin", "p", false, "accept Nebula password from stdin")

	return cmd
}
