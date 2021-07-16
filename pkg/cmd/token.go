package cmd

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/puppetlabs/relay/pkg/client/openapi"
	"github.com/puppetlabs/relay/pkg/config"
	"github.com/puppetlabs/relay/pkg/debug"
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
			return errors.New("A token by that name already exists")
		default:
			return err
		}
	}

	if token, ok := t.GetTokenOk(); ok {
		secret := token.UserTokenWithSecret.GetSecret()
		filepath, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		if filepath != "" {
			if err := ioutil.WriteFile(filepath, []byte(secret), 0644); err != nil {
				debug.Logf("failed to write to file %s: %s", filepath, err.Error())
				return err
			}
		}

		use, err := cmd.Flags().GetBool("use")
		if err != nil {
			return err
		}

		if use {
			writeAuthTokenConfig(cmd, *token.UserTokenWithSecret.Secret, config.AuthTokenTypeAPI)
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

		for _, token := range *tokens {
			t.AppendRow([]string{token.UserToken.GetUser().Name, token.UserToken.GetId(), token.UserToken.GetName(), token.UserToken.GetType()})
		}

		t.Flush()
	}

	return nil
}
