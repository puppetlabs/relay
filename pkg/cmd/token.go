package cmd

import (
	"errors"
	"net/http"
	"os"

	"github.com/puppetlabs/relay-client-go/client/pkg/client/openapi"
	"github.com/puppetlabs/relay/pkg/config"
	"github.com/spf13/cobra"
)

func newTokensCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Manage API tokens",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(newCreateToken())
	cmd.AddCommand(newListTokens())
	cmd.AddCommand(newRevokeToken())

	return cmd
}

func newCreateToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [token name]",
		Short: "Create API token",
		Args:  cobra.MinimumNArgs(1),
		RunE:  doCreateToken,
	}

	cmd.Flags().BoolP("use", "u", true, "Configure the CLI to use the generated API token")
	cmd.Flags().StringP("file", "f", "", "Write the generated token to the supplied file")

	return cmd
}

func newRevokeToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke [token id]",
		Short: "Revoke API token",
		Args:  cobra.MinimumNArgs(1),
		RunE:  doRevokeToken,
	}

	return cmd
}

func newListTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List API tokens",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doListTokens,
	}

	cmd.Flags().BoolP("all", "a", false, "Show all account tokens")

	return cmd
}

func doCreateToken(cmd *cobra.Command, args []string) error {
	file, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	var f *os.File
	if file != "" {
		f, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}

		defer f.Close()
	}

	Dialog.Progress("Creating token...")

	req := Client.Api.TokensApi.CreateToken(cmd.Context())
	t, resp, err := Client.Api.TokensApi.CreateTokenExecute(
		req.TokenRequest(openapi.TokenRequest{
			Name: args[0],
			Type: "user",
		}),
	)

	if err != nil {
		switch resp.StatusCode {
		case http.StatusConflict:
			// FIXME This is a bit of an assumption, but it is worth adding for overall usability.
			// A few things need to change to ensure an accurate error message is displayed.
			return errors.New("A token by that name already exists")
		default:
			return err
		}
	}

	if token, ok := t.GetTokenOk(); ok {
		secret := token.UserTokenWithSecret.GetSecret()
		if err != nil {
			return err
		}

		if file != "" {
			if _, err := f.Write([]byte(secret)); err != nil {
				return err
			}

			Dialog.Infof("Your token was written to %s\n"+
				"Use this file to authenticate to the Relay CLI by running: relay auth login --file=%s\n", file, file)
		}

		use, err := cmd.Flags().GetBool("use")
		if err != nil {
			return err
		}

		if use {
			writeAuthTokenConfig(cmd, *token.UserTokenWithSecret.Secret, config.AuthTokenTypeAPI)

			Dialog.WriteString("The generated token has been added to the cached credentials\n" +
				"To clear your cached credentials, use: relay config auth clear\n")
		}

		if !use && file == "" {
			Dialog.WriteString(secret)
		}
	}

	return nil
}

func doRevokeToken(cmd *cobra.Command, args []string) error {
	Dialog.Progress("Revoking token...")

	req := Client.Api.TokensApi.DeleteToken(cmd.Context(), args[0])
	_, _, err := Client.Api.TokensApi.DeleteTokenExecute(req)
	if err != nil {
		return err
	}

	return nil
}

func doListTokens(cmd *cobra.Command, args []string) error {
	Dialog.Progress("Listing tokens...")

	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		return err
	}

	req := Client.Api.TokensApi.GetTokens(cmd.Context())
	t, _, err := Client.Api.TokensApi.GetTokensExecute(req.Owned(!all).Valid(true))
	if err != nil {
		return err
	}

	if tokens, ok := t.GetTokensOk(); ok {
		t := Dialog.Table()

		t.Headers([]string{"User", "Id", "Name", "Type"})

		for _, token := range tokens {
			t.AppendRow([]string{token.UserToken.GetUser().Name, token.UserToken.GetId(), token.UserToken.GetName(), token.UserToken.GetType()})
		}

		t.Flush()
	}

	return nil
}
