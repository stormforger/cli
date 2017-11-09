package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	// testCaseUpdateCmd represents the testCaseValidate command
	testCaseUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a new test case",
		Long:  `Update a new test case.`,
		Run:   runTestCaseUpdate,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if testCaseUpdateOpts.Organisation == "" {
				testCaseUpdateOpts.Organisation = readOrganisationUIDFromFile()
				if testCaseUpdateOpts.Organisation == "" {
					log.Fatal("Missing organization flag")
				}
			}

			if testCaseUpdateOpts.UID == "" {
				log.Fatal("Missing test case UID flag")
			}
		},
	}

	testCaseUpdateOpts struct {
		Organisation string
		UID          string
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseUpdateCmd)

	testCaseUpdateCmd.PersistentFlags().StringVarP(&testCaseUpdateOpts.UID, "uid", "u", "", "UID of the test case")
}

func runTestCaseUpdate(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		fileName, testCaseFile, err := readFromStdinOrReadFirstArgument(args, "test_case.js")
		if err != nil {
			log.Fatal(err)
		}

		client := NewClient()

		success, message, err := client.TestCaseUpdate(testCaseUpdateOpts.Organisation, testCaseUpdateOpts.UID, fileName, testCaseFile)
		if err != nil {
			log.Fatal(err)
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
