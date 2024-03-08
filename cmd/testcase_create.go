package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
	"github.com/stormforger/cli/internal/pflagutil"
)

var (

	// testCaseCreateCmd represents the testCaseValidate command
	testCaseCreateCmd = &cobra.Command{
		Use:   "create <test-case-ref> <test-case-file>",
		Short: "Create a new test case",
		Long: fmt.Sprintf(`Create a new test case.

  <test-case-ref> is 'organisation-name/test-case-name'.
  <test-case-file> is a path or - for stdin.

Examples
--------
* Create a new test case named 'checkout' in the 'acme-inc' organisation

  forge test-case create acme-inc/checkout cases/checkout_process.js

* Alternatively the test definition can be piped in as well

  cat cases/checkout_process.js | forge test-case create acme-inc/checkout -

%s
`, bundlingHelpInfo),
		Run: runTestCaseCreate,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 {
				log.Fatal("Too many arguments")
			}

			if len(args) < 2 {
				log.Fatal("Missing arguments; test case reference and test case file (or - to read from stdin)")
			}

			segments := strings.Split(args[0], "/")

			if len(segments) != 2 {
				log.Fatal("Invalid argument: <test-case-ref> has to be like organisation-name/test-case-name")
			}

			testCaseCreateOpts.Organisation = lookupOrganisationUID(NewClient(), segments[0])
			if testCaseCreateOpts.Organisation == "" {
				log.Fatal("Missing organisation")
			}

			testCaseCreateOpts.Name = segments[1]
			if testCaseCreateOpts.Name == "" {
				log.Fatal("Missing test case name")
			}
		},
		ValidArgsFunction: completeOrgaAndCase,
	}

	testCaseCreateOpts struct {
		Organisation string
		Name         string
		Update       bool // update test-case if already exists
		Define       map[string]string
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseCreateCmd)

	testCaseCreateCmd.PersistentFlags().StringVarP(&testCaseCreateOpts.Name, "name", "n", "", "Name of the new test case")
	testCaseCreateCmd.PersistentFlags().BoolVar(&testCaseCreateOpts.Update, "update", false, "Update test-case instead, if it already exists")
	testCaseCreateCmd.PersistentFlags().Var(&pflagutil.KeyValueFlag{Map: &testCaseCreateOpts.Define}, "define", "Defines a list of K=V while parsing: debug=false")
}

func runTestCaseCreate(cmd *cobra.Command, args []string) {
	orgaUID := testCaseCreateOpts.Organisation

	bundler := testCaseFileBundler{Defines: testCaseCreateOpts.Define}
	bundle, err := bundler.Bundle(args[1], "test_case.js")
	if err != nil {
		log.Fatal(err)
	}

	var testCaseName string
	if testCaseCreateOpts.Name != "" {
		testCaseName = testCaseCreateOpts.Name
	} else if args[0] != "-" {
		basename := filepath.Base(args[0])
		testCaseName = strings.TrimSuffix(basename, filepath.Ext(basename))
	} else {
		log.Fatal("Name of test case missing")
		fmt.Println()
		os.Exit(1)
	}

	client := NewClient()

	testcaseUID := lookupTestCase(client, orgaUID+"/"+testCaseName)

	var (
		success       bool
		message       string
		errValidation error
	)
	if testcaseUID != "" && !testCaseCreateOpts.Update {
		printErrorPayloadHuman(os.Stderr, false, api.ErrorPayload{
			Message: "Test-Case already exists.",
		})
		cmdExit(false)
	} else if testcaseUID == "" {
		success, message, errValidation = client.TestCaseCreate(orgaUID, testCaseName, bundle.Name, bundle.Content)
	} else {
		success, message, errValidation = client.TestCaseUpdate(testcaseUID, bundle.Name, bundle.Content)
	}

	if errValidation != nil {
		log.Fatal(errValidation)
	}

	if rootOpts.OutputFormat == "json" {
		printValidationResultJSON(message)
		cmdExit(success)
	}

	errorMeta, err := api.ErrorDecoder{SourceMapper: bundle.Mapper}.UnmarshalErrorMeta(strings.NewReader(message))
	if err != nil {
		log.Fatal(err)
	}

	printErrorPayloadHuman(os.Stderr, success, errorMeta)
	cmdExit(success)
}
