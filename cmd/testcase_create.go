package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// testCaseCreateCmd represents the testCaseValidate command
	testCaseCreateCmd = &cobra.Command{
		Use:   "create <test-case-ref> <test-case-file>",
		Short: "Create a new test case",
		Long: `Create a new test case.

<test-case-ref> is 'organisation-name/test-case-name'.

Examples
--------
* create a new test case named 'checkout' in the 'acme-inc' organisation

  forge test-case create acme-inc/checkout cases/checkout_process.js

* alternatively the test definition can be piped in as well

  cat cases/checkout_process.js | forge test-case create acme-inc/checkout -

`,
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

			testCaseCreateOpts.Organisation = lookupOrganisationUID(*NewClient(), segments[0])
			if testCaseCreateOpts.Organisation == "" {
				log.Fatal("Missing organization")
			}

			testCaseCreateOpts.Name = segments[1]
			if testCaseCreateOpts.Name == "" {
				log.Fatal("Missing test case name")
			}
		},
	}

	testCaseCreateOpts struct {
		Organisation string
		Name         string
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseCreateCmd)

	testCaseCreateCmd.PersistentFlags().StringVarP(&testCaseCreateOpts.Name, "name", "n", "", "Name of the new test case")
}

func runTestCaseCreate(cmd *cobra.Command, args []string) {
	organizationUID := testCaseCreateOpts.Organisation

	fileName, testCaseFile, err := readFromStdinOrReadFromArgument(args, "test_case.js", 1)
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

	success, message, errValidation := client.TestCaseCreate(organizationUID, testCaseName, fileName, testCaseFile)
	if errValidation != nil {
		log.Fatal(errValidation)
	}

	if success {
		os.Exit(0)
	}

	printPrettyJSON(message)

	fmt.Println()
	os.Exit(1)
}
