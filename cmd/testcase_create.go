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
		Use:   "create",
		Short: "Create a new test case",
		Long:  `Create a new test case.`,
		Run:   runTestCaseCreate,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if testCaseCreateOpts.Organisation == "" {
				log.Fatal("Missing organization flag")
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

	testCaseCreateCmd.PersistentFlags().StringVarP(&testCaseCreateOpts.Organisation, "organization", "o", "", "Name of the organization")
	testCaseCreateCmd.PersistentFlags().StringVarP(&testCaseCreateOpts.Name, "name", "n", "", "Name of the new test case")
}

func runTestCaseCreate(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		fileName, testCaseFile, err := readFromStdinOrReadFirstArgument(args, "test_case.js")
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

		success, message, errValidation := client.TestCaseCreate(testCaseCreateOpts.Organisation, testCaseName, fileName, testCaseFile)
		if errValidation != nil {
			log.Fatal(errValidation)
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
