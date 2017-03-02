package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping the StormForger API",
	Long:  `Ping the StormForger API and try to authenticate.`,
	Run:   foo,
}

func foo(cmd *cobra.Command, args []string) {
	client := NewClient()

	status, err := client.Ping()

	if !status {
		fmt.Println(err)
		os.Exit(-1)
	} else {
		log.Println("pong!")
	}
}

func init() {
	RootCmd.AddCommand(pingCmd)
}
