package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

var (
	loginEmail    = ""
	loginPassword = ""

	loginCmd = &cobra.Command{
		Use:   "login <email>",
		Short: "Login to StormForger",
		Long: `Login to StormForger in order to acquire a JWT access token.

	You can provide the login via argument or --email flag.

	It is discouraged to provide the password via the --password flag. By
	default you are asked to provide the password interactively.`,
		Run: runLogin,
	}
)

func runLogin(cmd *cobra.Command, args []string) {
	client := NewClient()

	ensureEmail(args)
	ensurePassword()

	jwt, err := client.Login(loginEmail, loginPassword)

	if err != nil {
		log.Fatal(err)
	} else {
		color.White("Login successful! Add the JWT token to a .stormforger.toml file like this:\n\n")
		color.Green("  echo 'jwt = \"" + jwt + "\"' >> .stormforger.toml")
		color.Green("\n\n")
	}
}

func ensureEmail(args []string) {
	if loginEmail != "" {
		return
	}

	if len(args) == 1 {
		loginEmail = args[0]
	}

	if len(args) == 0 {
		fmt.Printf("No email for login provided, what is your email? ")
		stdInReader := bufio.NewReader(os.Stdin)
		line, _, err := stdInReader.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		loginEmail = string(line)
	}
}

func ensurePassword() {
	if loginPassword == "" {
		fmt.Printf("Password (will be masked): ")
		pass, err := gopass.GetPasswdMasked()

		if err != nil {
			log.Fatal(err)
		}

		loginPassword = string(pass)
	}
}

func init() {
	RootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&loginEmail, "email", "e", "", "Email for Login")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Log in with this password")
}
