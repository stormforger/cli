package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping the StormForger API",
	Long:  `Ping the StormForger API and try to authenticate.`,
	Run:   runPingCmd,
}

func runPingCmd(cmd *cobra.Command, args []string) {
	client := NewClient()

	status, response, err := client.Ping()

	if !status {
		fmt.Println("Could not perform authenticated ping! Please verify that you can login with these credentials at https://app.stormforger.com!")
		if err != nil {
			log.Println(err)
		}
		log.Fatal(string(response))
	} else {
		var data api.User
		if err := json.NewDecoder(bytes.NewReader(response)).Decode(&data); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("PONG! Authenticated as %s\n", data.Mail)
	}
}

func init() {
	RootCmd.AddCommand(pingCmd)
}
