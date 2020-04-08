package main

import (
	"fmt"

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
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Logged in")
		},
	}

	return cmd
}

func NewLogoutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out of Relay",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Logged out")
		},
	}

	return cmd
}
