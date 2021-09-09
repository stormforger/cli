package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	testRunUnArchiveCmd = &cobra.Command{
		Use:     "unarchive <test-run-ref>",
		Aliases: []string{"unar", "ua"},
		Short:   "Mark a test run as not archived.",
		Long:    `Mark a test run as not archived.`,
		Run:     testRunUnArchive,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("Missing argument: test run reference")
			}

			if len(args) > 2 {
				log.Fatal("Too many arguments")
			}
		},
	}
)

func init() {
	TestRunCmd.AddCommand(testRunUnArchiveCmd)
}

func testRunUnArchive(cmd *cobra.Command, args []string) {
	client := NewClient()

	testRunUID := getTestRunUID(*client, args[0])

	success, response, err := client.TestRunUnArchive(testRunUID)
	if err != nil {
		log.Fatal(err)
	}

	if !success {
		log.Fatalf("Test run could not be unarchived!\n%s\n", string(response))
	}
}
