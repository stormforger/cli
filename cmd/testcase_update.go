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

	fileName, testCaseFile, err := readTestCaseFromStdinOrReadFromArgument(args, "test_case.js", 1)
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

	// NOTE: The testcase api endpoint may return an API error with either a 200 or 400.
	//  200 - no errors
	//  200 - with errors field, in case of validation errors where the testcase is still saved
	//  400 - with errors field, if the testcase could not be parsed and saved

	errorMeta, err := api.UnmarshalErrorMeta(strings.NewReader(message))
	if err != nil {
		log.Fatal(err)
	}

	prefix := "INFO"
	if !success {
		prefix = "ERROR"
	} else if len(errorMeta.Errors) > 0 {
		prefix = "WARN"
	}

	fmt.Fprintf(os.Stderr, "%s: %s\n", prefix, errorMeta.Message)
	if len(errorMeta.Errors) > 0 {
		for i, e := range errorMeta.Errors {
			fmt.Fprintf(os.Stderr, "\n%d) %s: %s\n", i+1, e.Code, e.Title)
			fmt.Fprintf(os.Stderr, "%s\n", e.FormattedError)
		}
	}

	if success {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
