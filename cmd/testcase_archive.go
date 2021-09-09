package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	// testCaseArchiveCmd represents the test case archive command
	testCaseArchiveCmd = &cobra.Command{
		Use:     "archive <test-case-ref>",
		Aliases: []string{"ar", "a"},
		Short:   "Archive a test case.",
		Long: `Mark the specified test case as archived"

<test-case-ref> can be 'organisation-name/test-case-name' or 'test-case-uid'.
`,
		Run: runTestCaseArchive,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("Missing argument: test case reference")
			}

			if len(args) > 2 {
				log.Fatal("Too many arguments")
			}
		},
		ValidArgsFunction: completeOrgaAndCase,
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseArchiveCmd)
}

func runTestCaseArchive(cmd *cobra.Command, args []string) {
	client := NewClient()

	testCaseUID := mustLookupTestCase(client, args[0])

	success, response, err := client.TestCaseArchive(testCaseUID)
	if err != nil {
		log.Fatal(err)
	}

	if !success {
		log.Fatalf("Test case definition could not be archived!\n%s\n", string(response))
	}
}
