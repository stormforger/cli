package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	// testRunAbortCmd represents the calllog command
	testRunAbortCmd = &cobra.Command{
		Use:   "abort <test-run-id>",
		Short: "Abort the given running test",
		Long:  `Abort the given running test.`,
		Run:   testRunAbort,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatal("Expect exactly one argument: test run reference!")
			}
		},
	}
)

func init() {
	TestRunCmd.AddCommand(testRunAbortCmd)
}

func testRunAbort(cmd *cobra.Command, args []string) {
	client := NewClient()

	testRunUID := getTestRunUID(*client, args[0])

	status, response, err := client.TestRunAbort(testRunUID)
	if err != nil {
		log.Fatal(err)
	}

	if status {
		fmt.Println(response)

		os.Exit(0)
	} else {
		fmt.Fprintln(os.Stderr, "Could not abort test run!")
		fmt.Println(response)

		os.Exit(1)
	}
}
