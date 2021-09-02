package cmd

import "github.com/spf13/cobra"

var (
	// testCaseArchiveCmd represents the test case archive command
	testCaseArchiveCmd = &cobra.Command{
		Use:     "archive <test-case-ref>",
		Aliases: []string{},
		Short:   "Mark a test case as archived",
		Long: `Mark the specified test case as archived"

<test-case-ref> can be 'organisation-name/test-case-name' or 'test-case-uid'.
`,
		Run:               runTestCaseArchive,
		PersistentPreRun:  nil, // TODO
		ValidArgsFunction: completeOrgaAndCase,
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseArchiveCmd)
}

func runTestCaseArchive(cmd *cobra.Command, args []string) {
	// TODO
}
