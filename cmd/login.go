package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	loginEmail    = ""
	loginPassword = ""

	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login to StormForger",
		Long:  `Login to StormForger in order to acquire a JWT access token.`,
		Run: func(cmd *cobra.Command, args []string) {
			client := NewClient()

			jwt, err := client.Login(loginEmail, loginPassword)

			if err != nil {
				fmt.Fatal(err)
			} else {
				color.Green("Login successful! Here is your JWT access token:\n\n")
				color.Green("  " + jwt)
				color.Green("\n")
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&loginEmail, "email", "e", "", "Email for Login")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Log in with this password")
}
