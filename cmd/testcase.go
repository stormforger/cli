package cmd

import "github.com/spf13/cobra"

// TestCaseCmd is the cobra definition
var TestCaseCmd = &cobra.Command{
	Use:     "test-case",
	Aliases: []string{"testcase", "tc"},
	Short:   "Work with and manage test cases",
	Long:    `Work with and manage test cases.`,
}

func init() {
	RootCmd.AddCommand(TestCaseCmd)
}
