package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
)

// pingCmd represents the ping command
var (
	pingCmd = &cobra.Command{
		Use:   "ping",
		Short: "Ping the StormForger API",
		Long:  `Ping the StormForger API and try to authenticate.`,
		Run:   runPingCmd,
	}

	pingOpts struct {
		Unauthenticated bool
	}
)

func init() {
	RootCmd.AddCommand(pingCmd)

	pingCmd.Flags().BoolVarP(&pingOpts.Unauthenticated, "unauthenticated", "", false, "Perform unauthenticated ping")
}

func runPingCmd(cmd *cobra.Command, args []string) {
	client := NewClient()

	var printPingError func(bool, []byte, error) = printHumanPingError
	var printPingSuccess func(bool, *api.User) = printHumanPingSuccess
	if rootOpts.OutputFormat == "json" {
		printPingError = printJSONPingError
		printPingSuccess = printJSONPingSuccess
	}

	var status bool
	var response []byte
	var err error

	if pingOpts.Unauthenticated {
		status, response, err = client.PingUnauthenticated()
		if !status {
			printPingError(false, response, err)
			os.Exit(1)
		}
		printPingSuccess(false, nil)
	} else {
		status, response, err = client.Ping()
		if !status {
			printPingError(true, response, err)
			os.Exit(1)
		}
		var data api.User
		if err := json.NewDecoder(bytes.NewReader(response)).Decode(&data); err != nil {
			printPingError(true, response, err)
			os.Exit(1)
		}
		printPingSuccess(true, &data)
	}
}

func printHumanPingSuccess(authenticated bool, data *api.User) {
	if authenticated {
		if data.AuthenticatedAs.Type == api.UserTypeServiceAccount {
			fmt.Printf("PONG! Authenticated as %s (%s)\n", data.AuthenticatedAs.Label, data.AuthenticatedAs.UID)
		} else {
			fmt.Printf("PONG! Authenticated as %s\n", data.Mail)
		}
	} else {
		fmt.Println("PONG!")
	}
}

func printHumanPingError(authenticated bool, response []byte, err error) {
	if authenticated {
		fmt.Print("Could not perform authenticated ping!")
	} else {
		fmt.Print("Could not perform ping!")
	}
	fmt.Println(" Please verify that you can login with these credentials at https://app.stormforger.com!")

	if err != nil {
		fmt.Println("ERROR:", err)
	}
}

type pingLogMessage struct {
	Error       bool            `json:"error"`
	Message     string          `json:"message,omitempty"`
	APIResponse json.RawMessage `json:"api_response,omitempty"`
	User        *api.User       `json:"user,omitempty"`
}

func printJSONPingSuccess(authenticated bool, user *api.User) {
	logMessage := pingLogMessage{
		Error:   false,
		Message: "PONG!",
		User:    user,
	}
	json.NewEncoder(os.Stderr).Encode(logMessage)
}

func printJSONPingError(_authenticated bool, jsonResponse []byte, err error) {
	logMessage := pingLogMessage{
		Error:       true,
		APIResponse: jsonResponse,
	}
	if err != nil {
		logMessage.Message = err.Error()
	}
	json.NewEncoder(os.Stderr).Encode(logMessage)
}
