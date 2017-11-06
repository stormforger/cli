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
				log.Fatal("Missing organization flag")
			}
			if testCaseUpdateOpts.Uid == "" {
				log.Fatal("Missing test case UID flag")
			}
		},
	}

	testCaseUpdateOpts struct {
		Organisation string
		Uid          string
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseUpdateCmd)

	testCaseUpdateCmd.PersistentFlags().StringVarP(&testCaseUpdateOpts.Organisation, "organization", "o", "", "Name of the organization")
	testCaseUpdateCmd.PersistentFlags().StringVarP(&testCaseUpdateOpts.Uid, "uid", "u", "", "UID of the test case")
}

func runTestCaseUpdate(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		fileName, testCaseFile, err := readFromStdinOrReadFirstArgument(args, "test_case.js")
		if err != nil {
			log.Fatal(err)
		}

		client := NewClient()

		success, message, err := client.TestCaseUpdate(testCaseUpdateOpts.Organisation, testCaseUpdateOpts.Uid, fileName, testCaseFile)
		if err != nil {
			log.Fatal(err)
		}

		if success {
			os.Exit(0)
		}

		printPrettyJson(message)

		fmt.Println()
		os.Exit(1)
	} else {
		log.Fatal("Missing argument; test case file or - to read from stdin")
	}
}
