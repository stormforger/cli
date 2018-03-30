package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/testrun"
)

var (
	// testRunListCmd represents the calllog command
	testRunListCmd = &cobra.Command{
		Use:   "list <test-case-ref>",
		Short: "List of completed test runs",
		Long:  `List of completed test runs.`,
		Run:   testRunList,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				log.Fatal("Too many arguments")
			}

			if len(args) < 1 {
				log.Fatal("Missing argument: test case reference")
			}

			segments := strings.Split(args[0], "/")

			if len(segments) > 2 {
				log.Fatal("Invalid argument: <test-case-ref> has to be like organisation-name/test-case-name or test-case-uid")
			}
		},
	}
)

func init() {
	TestRunCmd.AddCommand(testRunListCmd)
}

func testRunList(cmd *cobra.Command, args []string) {
	client := NewClient()

	testCaseUID := lookupTestCase(*client, args[0])

	status, result, err := client.TestRunList(testCaseUID)
	if err != nil {
		log.Fatal(err)
	}

	if !status {
		fmt.Fprintln(os.Stderr, "Could not list test runs!")
		fmt.Fprintln(os.Stderr, string(result))

		os.Exit(1)
	}

	if rootOpts.OutputFormat == "json" {
		fmt.Println(string(result))
		return
	}

	items, err := testrun.Unmarshal(bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items.TestRuns {
		if rootOpts.OutputFormat == "human" {
			fmt.Printf("%s (ID: %s)\n", item.Scope, item.ID)
		} else {
			fmt.Printf("%s\n", item.Scope)
		}
	}
}
