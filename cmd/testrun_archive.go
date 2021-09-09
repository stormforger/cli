package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	testRunArchiveCmd = &cobra.Command{
		Use:     "archive <test-run-ref>",
		Aliases: []string{"ar", "a"},
		Short:   "Mark a test run as archived.",
		Long:    `Mark a test run as archived.`,
		Run:     testRunArchive,
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
	TestRunCmd.AddCommand(testRunArchiveCmd)
}

func testRunArchive(cmd *cobra.Command, args []string) {
	client := NewClient()

	testRunUID := getTestRunUID(*client, args[0])

	success, response, err := client.TestRunArchive(testRunUID)
	if err != nil {
		log.Fatal(err)
	}

	if !success {
		log.Fatalf("Test run could not be archived!\n%s\n", string(response))
	}

}
