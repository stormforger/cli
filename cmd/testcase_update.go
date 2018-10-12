package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
)

var (
	// testCaseUpdateCmd represents the testCaseValidate command
	testCaseUpdateCmd = &cobra.Command{
		Use:   "update <test-case-ref> <test-case-file>",
		Short: "Update an existing test case",
		Long: `Update an existing test case

<test-case-ref> can be 'organisation-name/test-case-name' or 'test-case-uid'.

Examples
--------
* update a test case by file

  forge test-case update acme-inc/checkout cases/checkout_process.js

* alternatively the test definition can be piped in as well

  cat cases/checkout_process.js | forge test-case update acme-inc/checkout -

`,
		Run: runTestCaseUpdate,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				log.Fatal("Missing arguments; test case reference and test case file (or - to read from stdin)")
			}

			if len(args) > 2 {
				log.Fatal("Too many arguments")
			}
		},
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseUpdateCmd)
}

func runTestCaseUpdate(cmd *cobra.Command, args []string) {
	client := NewClient()

	testCaseUID := lookupTestCase(*client, args[0])

	fileName, testCaseFile, err := readFromStdinOrReadFromArgument(args, "test_case.js", 1)
	if err != nil {
		log.Fatal(err)
	}

	success, message, err := client.TestCaseUpdate(testCaseUID, fileName, testCaseFile)
	if err != nil {
		log.Fatal(err)
	}

	if rootOpts.OutputFormat == "json" {
		fmt.Println(message)

		if success {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if success {
		os.Exit(0)
	}

	errorMeta, err := api.UnmarshalErrorMeta(strings.NewReader(message))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stderr, "%s\n\n", errorMeta.Message)
	for _, e := range errorMeta.Errors {
		fmt.Fprintf(os.Stderr, "%s: %s\n", e.Code, e.Title)
		fmt.Fprintln(os.Stderr, e.FormattedError)
	}

	os.Exit(1)
}
