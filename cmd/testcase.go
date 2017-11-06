package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// TestCaseCmd is the cobra definition
var TestCaseCmd = &cobra.Command{
	Use:     "test-case",
	Aliases: []string{"testcase", "tc"},
	Short:   "Work with and manage test cases",
	Long:    `Work with and manage test cases.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("Cannot be run without subcommand!")
	},
}

func init() {
	RootCmd.AddCommand(TestCaseCmd)
}
