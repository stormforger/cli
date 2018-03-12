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
		Use:   "create <organisation-ref|test-case-ref> <test-case-file>",
		Short: "Create a new test case",
		Long: `Create a new test case.

<test-case-ref> can be 'organisation/test-case'. <organisation-ref> is
either 'organisation name' or the organisation's UID.

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
				log.Fatal("Missing arguments; test case reference and test case file")
			}

			segments := strings.Split(args[0], "/")

			testCaseCreateOpts.Organisation = lookupOrganisationUID(*NewClient(), segments[0])
			if testCaseCreateOpts.Organisation == "" {
				log.Fatal("Missing organization")
			}

			if len(segments) == 2 {
				testCaseCreateOpts.Name = segments[1]
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

	if len(args) > 0 {
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
	} else {
		log.Fatal("Missing argument; test case file or - to read from stdin")
	}
}
