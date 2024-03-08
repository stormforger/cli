package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
	"github.com/stormforger/cli/internal/pflagutil"
)

var (
	// testCaseValidateCmd represents the testCaseValidate command
	testCaseValidateCmd = &cobra.Command{
		Use:   "validate <organisation-ref> <test-case-files>",
		Short: "Upload a test case definition JavaScript and validate it",
		Long: fmt.Sprintf(`Upload a test case definition JavaScript and validate it.

We do require the organisation in order to validate the test case against
the available resources and limits of that given organisation.

<organisation-ref> is the name or the UID of your organisation.
<test-case-files> is one or more file names to validate.

%s
`, bundlingHelpInfo),
		Example: `Validate a test case (with limits of 'acme-inc' organisation):

  forge test-case validate acme-inc cases/checkout_process.js

Alternatively the test definition can be piped in as well:

  cat cases/checkout_process.js | forge test-case validate acme-inc -

Verify multiple files at once:

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
		Define       map[string]string
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseValidateCmd)

	testCaseValidateCmd.PersistentFlags().Var(&pflagutil.KeyValueFlag{Map: &testCaseValidateOpts.Define}, "define", "Defines a list of K=V while parsing: debug=false")
}

func runTestCaseValidate(cmd *cobra.Command, args []string) {
	client := NewClient()

	allValid := true
	for _, arg := range args[1:] {
		argValid, err := runTestCaseValidateArg(cmd, client, arg)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "ERROR: %v for %s\n", err, arg)
		}
		allValid = allValid && argValid
	}

	cmdExit(allValid)
}

// runTestCaseValidateArg validates the argument and returns true is the testcase was valid.
func runTestCaseValidateArg(cmd *cobra.Command, client *api.Client, fileOrStdin string) (bool, error) {
	fmt.Fprintf(cmd.ErrOrStderr(), "# FILE: %s\n", fileOrStdin)

	bundler := testCaseFileBundler{Defines: testCaseValidateOpts.Define}
	bundle, err := bundler.Bundle(fileOrStdin, "test_case.js")
	if err != nil {
		return false, err
	}

	success, message, errValidation := client.TestCaseValidate(testCaseValidateOpts.Organisation, bundle.Name, bundle.Content)
	if errValidation != nil {
		return false, errValidation
	}
	// NOTE: We can get success, success with warnings or just straight up validation errors (success=false)
	// see testcase_update.go

	if rootOpts.OutputFormat == "json" {
		// if the user wants json, we don't bother to parse it and just dump it.
		printValidationResultJSON(message)
		return success, nil
	}

	errorMeta, err := api.ErrorDecoder{SourceMapper: bundle.Mapper}.UnmarshalErrorMeta(strings.NewReader(message))
	if err != nil {
		return false, err
	}

	printValidationResultHuman(os.Stderr, success, errorMeta)
	return success, nil
}
