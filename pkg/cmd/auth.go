package cmd

import (
	"fmt"
	"os"

	"github.com/cli/browser"
	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/spf13/cobra"
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

func doLogin(cmd *cobra.Command, args []string) error {
	Dialog.Progress("Getting authorization...")

	verificationURI, cterr := Client.CreateToken()

	if cterr != nil {
		return cterr
	}

	fmt.Fprintf(os.Stdout, "Press [Enter] to continue authorization in the web browser: %s\n", verificationURI)
	fmt.Scanln()

	if err := browser.OpenURL(verificationURI); err != nil {
		return errors.NewAuthFailedLoginError().WithCause(fmt.Errorf("error opening the web browser: %w", err))
	}

	Dialog.Info("Successfully stored authorization token.")

	return nil
}

func newLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to Relay",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doLogin,
	}

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
