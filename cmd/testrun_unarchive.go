package cmd

import "github.com/spf13/cobra"

var (
	testRunUnArchiveCmd = &cobra.Command{
		Use:              "unarchive <test-run-ref>",
		Short:            "Mark a test run as not archived.",
		Long:             `Mark a test run as not archived.`,
		Run:              testRunUnArchive,
		PersistentPreRun: nil, // TODO

	}
)

func init() {
	TestRunCmd.AddCommand(testRunUnArchiveCmd)
}

func testRunUnArchive(cmd *cobra.Command, args []string) {
	// TODO
}
