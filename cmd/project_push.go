package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	projectPushCmd = &cobra.Command{
		Use:   "push <file>",
		Short: "Upload a test case or data source",
		Run:   runProjectPush,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		}}
)

func runProjectPush(cmd *cobra.Command, args []string) {
	log.Println("...")
}

func init() {
	projectCmd.AddCommand(projectPushCmd)
}
