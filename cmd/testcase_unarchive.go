package cmd

import "github.com/spf13/cobra"

var (
	// testCaseArchiveCmd represents the test case archive command
	testCaseUnArchiveCmd = &cobra.Command{
		Use: "archive <test-case-ref>",
		Aliases: []string{},
		Short: "Mark a test case as not archived",
		Long: `Mark the specified test case as not archived"

<test-case-ref> can be 'organisation-name/test-case-name' or 'test-case-uid'.
`,
		Run: runTestCaseUnArchive,
		PersistentPreRun: nil, // TODO
		ValidArgsFunction: completeOrgaAndCase,
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseUnArchiveCmd)
}

func runTestCaseUnArchive(cmd *cobra.Command, args []string) {
	// TODO
}
