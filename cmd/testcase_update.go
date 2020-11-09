package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
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
		ValidArgsFunction: completeOrgaAndCase,
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseUpdateCmd)
}

func runTestCaseUpdate(cmd *cobra.Command, args []string) {
	client := NewClient()

	testCaseUID := mustLookupTestCase(client, args[0])

	fileName, testCaseFile, mapper, err := readTestCaseBundleFromStdinOrReadFromArgument(args[1], "test_case.js")
	if err != nil {
		log.Fatal(err)
	}

	success, message, err := client.TestCaseUpdate(testCaseUID, fileName, testCaseFile)
	if err != nil {
		log.Fatal(err)
	}

	if rootOpts.OutputFormat == "json" {
		// if the user wants json, we don't bother to parse it and just dump it.
		printValidationResultJSON(message)
		cmdExit(success)
	}

	// NOTE: The testcase api endpoint may return an API error with either a 200 or 400.
	//  200 - no errors
	//  200 - with errors field, in case of validation errors where the testcase is still saved
	//  400 - with errors field, if the testcase could not be parsed and saved

	errorMeta, err := api.ErrorDecoder{SourceMapper: mapper}.UnmarshalErrorMeta(strings.NewReader(message))
	if err != nil {
		log.Fatal(err)
	}

	printValidationResultHuman(os.Stderr, fileName, success, errorMeta)
	cmdExit(success)
}

func cmdExit(success bool) {
	if success {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func printValidationResultJSON(message string) {
	fmt.Println(message)
}

func printValidationResultHuman(fp io.Writer, fileName string, success bool, errorMeta api.ErrorPayload) {
	prefix := "INFO"
	if !success {
		prefix = color.RedString("ERROR")
	} else if len(errorMeta.Errors) > 0 {
		prefix = color.YellowString("WARN")
	}

	fmt.Fprintf(fp, "%s: %s\n", prefix, errorMeta.Message)

	for i, e := range errorMeta.Errors {
		fmt.Fprintf(fp, "\n%d) %s: %s\n", i+1, e.Code, e.Title)
		fmt.Fprintf(fp, "<%s>\n", e.FormattedError)
	}
}
