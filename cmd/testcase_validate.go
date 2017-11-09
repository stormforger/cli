package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	// testCaseValidateCmd represents the testCaseValidate command
	testCaseValidateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Upload a test case definition JavaScript and validate it",
		Long:  `Upload a test case definition JavaScript and validate it.`,
		Run:   runTestCaseValidate,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if testCaseValidateOpts.Organisation == "" {
				testCaseValidateOpts.Organisation = readOrganisationUIDFromFile()
				if testCaseValidateOpts.Organisation == "" {
					log.Fatal("Missing organization flag")
				}
			}
		},
	}

	testCaseValidateOpts struct {
		Organisation string
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseValidateCmd)

	testCaseValidateCmd.PersistentFlags().StringVarP(&testCaseValidateOpts.Organisation, "organization", "o", "", "Name of the organization")
}

func runTestCaseValidate(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		fileName, testCaseFile, err := readFromStdinOrReadFirstArgument(args, "test_case.js")
		if err != nil {
			log.Fatal(err)
		}

		client := NewClient()

		success, message, errValidation := client.TestCaseValidate(testCaseValidateOpts.Organisation, fileName, testCaseFile)
		if errValidation != nil {
			log.Fatal(errValidation)
		}

		if success {
			fmt.Println("test case ok")
			os.Exit(0)
		}

		printPrettyJSON(message)

		fmt.Println()
		os.Exit(1)
	} else {
		log.Fatal("Missing argument; test case file or - to read from stdin")
	}
}
