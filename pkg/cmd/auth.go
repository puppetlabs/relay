package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cli/browser"
	"github.com/eiannone/keyboard"
	"github.com/puppetlabs/relay/pkg/config"
	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/puppetlabs/relay/pkg/util"
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

func negotiateSession(cmd *cobra.Command) error {
	deviceValues, cterr := Client.CreateToken()
	if cterr != nil {
		return cterr
	}

	writeAuthTokenConfig(cmd, deviceValues.Token.String(), config.AuthTokenTypeSession)

	Dialog.Info("Stored authorization token.")

	Dialog.Info(fmt.Sprintf(
		`Your one-time code for activation is:

**%s**
* %s *
**%s**

Press [ENTER] to open %s in a browser or any other key to cancel...`,
		strings.Repeat("*", len(deviceValues.UserCode)),
		deviceValues.UserCode,
		strings.Repeat("*", len(deviceValues.UserCode)),
		deviceValues.VerificationURI,
	))
	_, key, err := keyboard.GetSingleKey()
	if err != nil {
		return errors.NewGeneralUnknownError().WithCause(err)
	}

	if key != keyboard.KeyEnter {
		Dialog.Info("Canceled.")
		return nil
	}

	// The complete url may be empty, depending on the Device Auth Flow implementation.
	var uri string
	if deviceValues.VerificationURIComplete != "" {
		uri = deviceValues.VerificationURIComplete
	} else {
		uri = deviceValues.VerificationURI
	}
	if err := browser.OpenURL(uri); err != nil {
		return errors.NewAuthFailedLoginError().WithCause(fmt.Errorf("error opening the web browser: %w", err))
	}

	return nil
}

func readAuthFromStdin(cmd *cobra.Command) error {
	gotStdin, err := util.PassedStdin()
	if err != nil {
		return err
	}

	if gotStdin {
		token, err := util.ReadStdin(readLimit)
		if err != nil {
			return err
		}

		writeAuthTokenConfig(cmd, string(token), config.AuthTokenTypeAPI)
	}

	return nil
}
func doLogin(cmd *cobra.Command, args []string) error {
	Dialog.Progress("Getting authorization...")

	stdin, err := cmd.Flags().GetBool("stdin")
	if err != nil {
		return err
	}

	if stdin {
		err = readAuthFromStdin(cmd)
		if err != nil {
			return err
		}

		return nil
	}

	file, err := cmd.Flags().GetString("file")
	if err != nil {
		return nil
	}

	if file != "" {
		token, err := ioutil.ReadFile(file)
		if err != nil {
			return nil
		}

		writeAuthTokenConfig(cmd, string(token), config.AuthTokenTypeAPI)

		return nil
	}

	err = negotiateSession(cmd)
	if err != nil {
		return err
	}

	Dialog.Info("Done!")
	return nil
}

func newLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to Relay",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doLogin,
	}

	cmd.Flags().StringP("file", "f", "", "Read authentication credentials from file")
	cmd.Flags().Bool("stdin", false, "Read authentication credentials from stdin")

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

func writeAuthTokenConfig(cmd *cobra.Command, token string, tokenType config.AuthTokenType) {
	if len(token) > 0 {
		cfg := &config.Config{
			ContextConfig: map[string]*config.ContextConfig{
				Config.CurrentContext: {
					Auth: &config.AuthConfig{
						Tokens: map[config.AuthTokenType]string{
							tokenType: strings.TrimSpace(token),
						},
					},
				},
			},
		}

		config.WriteConfig(cfg, cmd.Flags())
	}
}
