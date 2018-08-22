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

func runPingCmd(cmd *cobra.Command, args []string) {
	client := NewClient()

	var status bool
	var response []byte
	var err error

	if pingOpts.Unauthenticated {
		status, response, err = client.PingUnauthenticated()

		if !status {
			fmt.Println("Could not perform ping! Please verify that you can connect to https://api.stormforger.com")
			if err != nil {
				log.Println(err)
			}
			log.Fatal(string(response))
		} else {
			fmt.Println("PONG!")
		}

	} else {
		status, response, err = client.Ping()

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
}

func init() {
	RootCmd.AddCommand(pingCmd)

	pingCmd.Flags().BoolVarP(&pingOpts.Unauthenticated, "unauthenticated", "", false, "Perform unauthenticated ping")
}
