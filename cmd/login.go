package cmd

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

var (
	noSave        = false
	loginEmail    = ""
	loginPassword = ""

	loginCmd = &cobra.Command{
		Use:   "login [email]",
		Short: "Login to StormForger",
		Long: `Login to StormForger in order to acquire a JWT access token.

	It is discouraged to provide the password via the --password flag. By
	default you are asked to provide the password interactively.`,
		Args: cobra.RangeArgs(0, 1),
		Run:  runLogin,
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

		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		stormforgerConfig := filepath.Join(home, ConfigFilename+".toml")

		if _, err := os.Stat(stormforgerConfig); !noSave && os.IsNotExist(err) {
			var bar = struct {
				JWT string `toml:"jwt"`
			}{JWT: jwt}

			content, err := toml.Marshal(bar)
			if err != nil {
				log.Fatal(err)
			}

			err = ioutil.WriteFile(stormforgerConfig, content, 0644)
			if err != nil {
				log.Fatal(err)
			}
			color.White("\nLogin successful!\n\n")
			color.Red("JWT token stored at %s.\n\n", stormforgerConfig)

			setupConfig()
		} else {
			color.White("\nLogin successful!\n\n")
			color.Red("Found %s. File will not be overridden!\n\n", stormforgerConfig)
			color.White("Add the JWT token to a .stormforger.toml file like this:\n\n")
			color.Green("  echo 'jwt = \"" + jwt + "\"' >> ~/.stormforger.toml")
			color.Green("\n\n")
		}
	}
}

func ensureEmail(args []string) {
	if len(args) == 1 {
		loginEmail = args[0]
		fmt.Printf("Email: %s\n", loginEmail)
	}

	if len(args) == 0 {
		fmt.Printf("No email for login provided\nEmail: ")
		stdInReader := bufio.NewReader(os.Stdin)
		line, _, err := stdInReader.ReadLine()
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		loginEmail = string(line)

		if loginEmail == "" {
			fmt.Println()
			log.Fatal("No email provided")
		}
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

	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Log in with this password")
	loginCmd.Flags().BoolVarP(&noSave, "no-save", "", false, "Don't save acquired JWT")
}
