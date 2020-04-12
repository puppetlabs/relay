package main

import (
	"fmt"

	"github.com/puppetlabs/relay/pkg/config"
	"github.com/spf13/cobra"
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
		Use:   "login",
		Short: "Log in to Relay",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, cfgerr := config.GetConfig(cmd.Flags())

			if cfgerr != nil {
				return cfgerr
			}

			fmt.Printf("Logged in. Config is %+v\n", cfg)

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
			fmt.Println("Logged out")

			return nil
		},
	}

	return cmd
}
