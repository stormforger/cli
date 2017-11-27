package cmd

import "github.com/spf13/cobra"

// TestRunCmd is the cobra definition
var TestRunCmd = &cobra.Command{
	Use:     "test-run",
	Aliases: []string{"testrun", "tr"},
	Short:   "Work with and manage test runs",
	Long:    `Work with and manage test runs.`,
}

func init() {
	RootCmd.AddCommand(TestRunCmd)
}
