package auth

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/puppetlabs/relay/pkg/client"
	"github.com/puppetlabs/relay/pkg/config"
	"github.com/puppetlabs/relay/pkg/dialog"
	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/puppetlabs/relay/pkg/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// readLimit is the max bytes allowed from stdin (512 bytes being well above a reasonable
// password length) to avoid reading massive byte arrays into memory and sending them
// to a remote API server.
const readLimit = 512

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage your authentication credentials",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(NewLoginCommand())
	cmd.AddCommand(NewLogoutCommand())

	return cmd
}

type loginParameters struct {
	Password string
	Email    string
}

func getLoginParameters(cmd *cobra.Command, args []string) (*loginParameters, errors.Error) {
	passFromStdin, perr := cmd.Flags().GetBool("password-stdin")

	if perr != nil {
		return nil, errors.NewAuthFailedLoginError().WithCause(perr)
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
	cfg, cfgerr := config.FromFlags(cmd.Flags())

	if cfgerr != nil {
		return cfgerr
	}

	log := dialog.FromConfig(cfg)

	loginParams, lperr := getLoginParameters(cmd, args)

	if lperr != nil {
		return lperr
	}

	log.Info("Logging in...")

	client := client.NewClient(cfg)

	cterr := client.CreateToken(loginParams.Email, loginParams.Password)

	if cterr != nil {
		return cterr
	}

	log.Info("Sucessfully logged in!")

	return nil
}

func NewLoginCommand() *cobra.Command {
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
	cfg, cfgerr := config.FromFlags(cmd.Flags())

	if cfgerr != nil {
		return cfgerr
	}

	log := dialog.FromConfig(cfg)

	log.Info("Logging out...")

	client := client.NewClient(cfg)

	iterr := client.InvalidateToken()

	if iterr != nil {
		return iterr
	}

	log.Info("You have been sucesfully logged out.")

	return nil
}

func NewLogoutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out of Relay",
		RunE:  doLogout,
	}

	return cmd
}
