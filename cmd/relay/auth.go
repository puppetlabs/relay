package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/puppetlabs/relay/pkg/client"
	"github.com/puppetlabs/relay/pkg/config"
	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func NewAuthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage your authentication credentials",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(NewLoginCommand())
	cmd.AddCommand(NewLogoutCommand())

	return cmd
}

func NewLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [email]",
		Short: "Log in to Relay",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, cfgerr := config.GetConfig(cmd.Flags())

			if cfgerr != nil {
				return cfgerr
			}

			loginParams, lperr := getLoginParameters(args)

			if lperr != nil {
				return lperr
			}

			client := client.NewClient(cfg)

			cterr := client.CreateToken(loginParams.Email, loginParams.Password)

			if cterr != nil {
				return cterr
			}

			return nil
		},
	}

	return cmd
}

func NewLogoutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out of Relay",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, cfgerr := config.GetConfig(cmd.Flags())

			if cfgerr != nil {
				return cfgerr
			}

			client := client.NewClient(cfg)

			iterr := client.InvalidateToken()

			if iterr != nil {
				return iterr
			}

			return nil
		},
	}

	return cmd
}

type loginParameters struct {
	Password string
	Email    string
}

func getLoginParameters(args []string) (*loginParameters, errors.Error) {
	var email string

	if len(args) > 0 {
		email = args[0]
	}

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

	password := strings.TrimSpace(string(passBytes))

	return &loginParameters{
		Email:    email,
		Password: password,
	}, nil
}
