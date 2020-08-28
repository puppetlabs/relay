package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/puppetlabs/relay/pkg/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// readLimit is set to 10kb to support RSA key files and the like.
const readLimit = 10 * 1024

func newAuthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage your authentication credentials",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(newLoginCommand())
	cmd.AddCommand(newLogoutCommand())

	return cmd
}

type loginParameters struct {
	Password string
	Email    string
}

func getLoginParameters(cmd *cobra.Command, args []string) (*loginParameters, errors.Error) {
	passFromStdin, err := cmd.Flags().GetBool("password-stdin")

	if err != nil {
		return nil, errors.NewAuthFailedLoginError().WithCause(err)
	}

	var email string
	var password string

	if len(args) > 0 {
		email = args[0]
	}

	if passFromStdin {
		gotStdin, stdinerr := util.PassedStdin()

		if stdinerr != nil {
			return nil, errors.NewAuthFailedPassFromStdin().WithCause(stdinerr)
		}

		if gotStdin {
			buf := bytes.Buffer{}
			reader := &io.LimitedReader{R: os.Stdin, N: readLimit}

			_, berr := buf.ReadFrom(reader)
			if berr != nil && berr != io.EOF {
				return nil, errors.NewAuthFailedPassFromStdin().WithCause(berr)
			}

			password = buf.String()

			if email == "" {
				return nil, errors.NewAuthMismatchedEmailPassMethods()
			}
		} else {
			return nil, errors.NewAuthFailedNoStdin()
		}
	} else {
		if email == "" {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Email: ")
			promptEmail, eperr := reader.ReadString('\n')

			if eperr != nil {
				return nil, errors.NewAuthFailedLoginError().WithCause(eperr)
			}

			email = strings.TrimSpace(promptEmail)
		}

		fmt.Print("Password: ")
		passBytes, pperr := terminal.ReadPassword(int(syscall.Stdin))
		if pperr != nil {
			return nil, errors.NewAuthFailedLoginError().WithCause(pperr)
		}

		password = string(passBytes)

		// resets to new line after password input
		fmt.Println("")
	}

	return &loginParameters{
		Email:    email,
		Password: strings.TrimSpace(password),
	}, nil
}

func doLogin(cmd *cobra.Command, args []string) error {
	loginParams, lperr := getLoginParameters(cmd, args)

	if lperr != nil {
		return lperr
	}

	Dialog.Progress("Logging in...")

	cterr := Client.CreateToken(loginParams.Email, loginParams.Password)

	if cterr != nil {
		return cterr
	}

	Dialog.Info("Successfully logged in!")

	return nil
}

func newLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [email]",
		Short: "Log in to Relay",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doLogin,
	}

	cmd.Flags().BoolP("password-stdin", "p", false, "accept password from stdin")

	return cmd
}

func doLogout(cmd *cobra.Command, args []string) error {
	Dialog.Progress("Logging out...")

	iterr := Client.InvalidateToken()

	if iterr != nil {
		return iterr
	}

	Dialog.Info("You have been successfully logged out.")

	return nil
}

func newLogoutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out of Relay",
		RunE:  doLogout,
	}

	return cmd
}
