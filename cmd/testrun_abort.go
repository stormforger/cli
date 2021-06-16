package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	// testRunAbortCmd represents the abort command
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

	testRunAbortAllCmd = &cobra.Command{
		Use:   "abort-all",
		Short: "Abort all running tests",
		Long:  "Abort all running tests",
		Run:   testRunAbortAll,
	}
)

func init() {
	TestRunCmd.AddCommand(testRunAbortCmd)
	TestRunCmd.AddCommand(testRunAbortAllCmd)
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

func testRunAbortAll(cmd *cobra.Command, args []string) {
	c := NewClient()
	ok, resp, err := c.TestRunAbortAll()
	if err != nil {
		log.Fatal(err)
	}

	if ok {
		fmt.Println(resp)
		os.Exit(0)
	}

	fmt.Fprintln(os.Stderr, "Could not abort all test runs!")
	fmt.Println(resp)
	os.Exit(1)
}
