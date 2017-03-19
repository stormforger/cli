package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// TestRunCmd is the cobra definition
var TestRunCmd = &cobra.Command{
	Use:   "test-run",
	Short: "Work with and manage test runs",
	Long:  `Work with and manage test runs.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("Cannot be run without subcommand!")
	},
}

func init() {
	RootCmd.AddCommand(TestRunCmd)
}
