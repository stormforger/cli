package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	// testCaseUnArchiveCmd represents the test case unarchive command
	testCaseUnArchiveCmd = &cobra.Command{
		Use:     "unarchive <test-case-ref>",
		Aliases: []string{},
		Short:   "Mark a test case as not archived",
		Long: `Mark the specified test case as not archived"

<test-case-ref> can be 'organisation-name/test-case-name' or 'test-case-uid'.
`,
		Run:               runTestCaseUnArchive,
		PersistentPreRun:  func(cmd *cobra.Command, args []string) {
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
	TestCaseCmd.AddCommand(testCaseUnArchiveCmd)
}

func runTestCaseUnArchive(cmd *cobra.Command, args []string) {
	// TODO
	client := NewClient()

	testCaseUID := mustLookupTestCase(client, args[0])

	success, err := client.TestCaseUnArchive(testCaseUID)
	if err != nil {
		log.Fatal(err)
	}

	if !success {
		log.Fatalf("Test case definition could not be unarchived.\n")
	}
}
