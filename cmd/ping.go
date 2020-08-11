package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

type PingCommandResult struct {
	Success  bool            `json:"success"`            // set to true if the ping was successful
	Error    error           `json:"error,omitempty"`    // Available if the ping call failed on a fundamental way (network error, unparsable response etc)
	Response json.RawMessage `json:"response,omitempty"` // Available if we talked to the StormForger API

	Unauthenticated bool      `json:"unauthenticated"`   // True if the ping command was unauthenticated
	Subject         *api.User `json:"subject,omitempty"` // Available if the ping suceeded and was an authenticated ping
}

func runPingCmd(cmd *cobra.Command, args []string) {
	client := NewClient()

	var printPingCommandResult func(io.Writer, PingCommandResult) = PrintHumanPingResult
	if rootOpts.OutputFormat == "json" {
		printPingCommandResult = printJSONPingResult
	}

	var status bool
	var response []byte
	var err error
	var user *api.User

	if pingOpts.Unauthenticated {
		status, response, err = client.PingUnauthenticated()
	} else {
		status, response, err = client.Ping()
		if status {
			var u api.User
			err = json.NewDecoder(bytes.NewReader(response)).Decode(&u)
			user = &u
		}
	}

	result := PingCommandResult{
		Success:         status,
		Error:           err,
		Response:        response,
		Unauthenticated: pingOpts.Unauthenticated,
		Subject:         user,
	}
	printPingCommandResult(os.Stdout, result)
}

func PrintHumanPingResult(w io.Writer, result PingCommandResult) {
	if result.Success {
		printHumanPingSuccess(w, !result.Unauthenticated, result.Subject)
	} else {
		printHumanPingError(w, !result.Unauthenticated, result.Response, result.Error)
	}
}

func printHumanPingSuccess(w io.Writer, authenticated bool, data *api.User) {
	if authenticated {
		if data.AuthenticatedAs.Type == api.UserTypeServiceAccount {
			fmt.Fprintf(w, "PONG! Authenticated as %s (%s)\n", data.AuthenticatedAs.Label, data.AuthenticatedAs.UID)
		} else {
			fmt.Fprintf(w, "PONG! Authenticated as %s\n", data.Mail)
		}
	} else {
		fmt.Fprintln(w, "PONG!")
	}
}

func printHumanPingError(w io.Writer, authenticated bool, response []byte, err error) {
	if authenticated {
		fmt.Fprint(w, "Could not perform authenticated ping!")
	} else {
		fmt.Fprint(w, "Could not perform ping!")
	}
	fmt.Fprintln(w, " Please verify that you can login with these credentials at https://app.stormforger.com!")

	if err != nil {
		fmt.Fprintln(w, "ERROR:", err)
	}
}

func printJSONPingResult(w io.Writer, result PingCommandResult) {
	json.NewEncoder(w).Encode(result)
}
