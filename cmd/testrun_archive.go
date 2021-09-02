package cmd

import "github.com/spf13/cobra"

var (
	testRunArchiveCmd = &cobra.Command{
		Use:              "archive <test-run-ref>",
		Short:            "Mark a test run as archived.",
		Long:             `Mark a test run as archived.`,
		Run:              testRunArchive,
		PersistentPreRun: nil, // TODO

	}
)

func init() {
	TestRunCmd.AddCommand(testRunArchiveCmd)
}

func testRunArchive(cmd *cobra.Command, args []string) {
	// TODO
}
