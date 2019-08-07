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
	// testCaseValidateCmd represents the testCaseValidate command
	testCaseValidateCmd = &cobra.Command{
		Use:   "validate <organisation-ref> <test-case-file>",
		Short: "Upload a test case definition JavaScript and validate it",
		Long: `Upload a test case definition JavaScript and validate it.

We do require the organisation in order to validate the test case against
the available resources and limits of that given organisation.

<organisation-ref> is the name or the UID of your organisation.

Examples
--------
* validate a test case (with limits of 'acme-inc' organisation)

  forge test-case validate acme-inc cases/checkout_process.js

* alternatively the test definition can be piped in as well

  cat cases/checkout_process.js | forge test-case validate acme-inc -

`,

		Run: runTestCaseValidate,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 {
				log.Fatal("Too many arguments")
			}

			if len(args) < 2 {
				log.Fatal("Missing arguments; organisation reference and test case file to validate (or - to read from stdin)")
			}

			testCaseValidateOpts.Organisation = lookupOrganisationUID(*NewClient(), args[0])
			if testCaseValidateOpts.Organisation == "" {
				log.Fatal("Missing organisation")
			}
		},
	}

	testCaseValidateOpts struct {
		Organisation string
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseValidateCmd)
}

func runTestCaseValidate(cmd *cobra.Command, args []string) {
	fileName, testCaseFile, err := readTestCaseFromStdinOrReadFromArgument(args, "test_case.js", 1)
	if err != nil {
		log.Fatal(err)
	}

	client := NewClient()

	success, message, errValidation := client.TestCaseValidate(testCaseValidateOpts.Organisation, fileName, testCaseFile)
	if errValidation != nil {
		log.Fatal(errValidation)
	}

	if rootOpts.OutputFormat == "json" {
		fmt.Println(message)

		if success {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	errorMeta, err := api.UnmarshalErrorMeta(strings.NewReader(message))
	if err != nil {
		log.Fatal(err)
	}

	if len(errorMeta.Errors) == 0 {
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "%s\n\n", errorMeta.Message)
	for i, e := range errorMeta.Errors {
		fmt.Fprintf(os.Stderr, "%d) %s: %s\n", i+1, e.Code, e.Title)
		fmt.Fprintf(os.Stderr, "%s\n\n", e.FormattedError)
	}

	os.Exit(1)
}
