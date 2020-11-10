package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
)

var (
	// testCaseValidateCmd represents the testCaseValidate command
	testCaseValidateCmd = &cobra.Command{
		Use:   "validate <organisation-ref> <test-case-files>",
		Short: "Upload a test case definition JavaScript and validate it",
		Long: `Upload a test case definition JavaScript and validate it.

We do require the organisation in order to validate the test case against
the available resources and limits of that given organisation.

<organisation-ref> is the name or the UID of your organisation
<test-case-files> is one or more file names to validate

Examples
--------
* Validate a test case (with limits of 'acme-inc' organisation)

  forge test-case validate acme-inc cases/checkout_process.js

* Alternatively the test definition can be piped in as well

	cat cases/checkout_process.js | forge test-case validate acme-inc -

* Verify multiple files at once

	forge test-case validate acme-inc ./dist/foo.js ./dist/bar.js ./dist/foobar.js

`,

		Run: runTestCaseValidate,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				log.Fatal("Missing arguments; organisation reference and test case file to validate (or - to read from stdin)")
			}

			stdinUsed := false
			for _, arg := range args {
				if arg == "-" {
					if stdinUsed {
						log.Fatalf("Stdin ('-') provided multiple times")
					}
					stdinUsed = true
				}
			}

			testCaseValidateOpts.Organisation = lookupOrganisationUID(NewClient(), args[0])
			if testCaseValidateOpts.Organisation == "" {
				log.Fatal("Missing organisation")
			}
		},
		ValidArgsFunction: completeOrga,
	}

	testCaseValidateOpts struct {
		Organisation string
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseValidateCmd)
}

func runTestCaseValidate(cmd *cobra.Command, args []string) {
	client := NewClient()

	validationError := false
	for _, arg := range args[1:] {
		argValidationError, err := runTestCaseValidateArg(cmd, client, arg)
		if err != nil {
			log.Fatalf("ERROR: %v for %s\n", err, arg)
		}
		if argValidationError {
			validationError = true
		}
	}

	if validationError {
		os.Exit(1)
	}
}

// runTestCaseValidateArg returns true if there were any validation ERRORS (not warnings)!
func runTestCaseValidateArg(cmd *cobra.Command, client *api.Client, fileOrStdin string) (bool, error) {
	result, err := testCaseFileBundler{}.bundle(fileOrStdin, "test_case.js")
	if err != nil {
		return true, err
	}

	success, message, errValidation := client.TestCaseValidate(testCaseValidateOpts.Organisation, result.Name, result.Content)
	if errValidation != nil {
		return true, errValidation
	}
	// NOTE: We can get success, success with warnings or just straight up validation errors (success=false)
	// see testcase_update.go

	if rootOpts.OutputFormat == "json" {
		// if the user wants json, we don't bother to parse it and just dump it.
		printValidationResultJSON(message)
		return !success, nil
	}

	errorMeta, err := api.ErrorDecoder{SourceMapper: result.Mapper}.UnmarshalErrorMeta(strings.NewReader(message))
	if err != nil {
		return true, err
	}

	printValidationResultHuman(os.Stderr, result.Name, success, errorMeta)

	if len(errorMeta.Errors) == 0 {
		return false, nil
	}
	return true, nil
}
