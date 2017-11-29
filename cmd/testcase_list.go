package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/testcase"
)

var (
	// testCaseListCmd represents the testCaseValidate command
	testCaseListCmd = &cobra.Command{
		Use:   "list",
		Short: "List a new test case",
		Long:  `List a new test case.`,
		Run:   runTestCaseList,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) >= 1 {
				testCaseListOpts.Organisation = args[0]
			} else {
				testCaseListOpts.Organisation = ""
			}

			if testCaseListOpts.Organisation == "" {
				testCaseListOpts.Organisation = readOrganisationUIDFromFile()
				if testCaseListOpts.Organisation == "" {
					log.Fatal("Missing organization flag")
				}
			}
		},
	}

	testCaseListOpts struct {
		Organisation string
		Name         string
		JSON         bool
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseListCmd)

	testCaseListCmd.Flags().BoolVarP(&testCaseListOpts.JSON, "json", "", false, "Output machine-readable JSON")
}

func runTestCaseList(cmd *cobra.Command, args []string) {
	client := NewClient()

	status, result, err := client.ListTestCases(testCaseListOpts.Organisation)
	if err != nil {
		log.Fatal(err)
	}

	if status {
		if testCaseListOpts.JSON {
			fmt.Println(string(result))
			return
		}

		testcase.ShowNames(bytes.NewReader(result))
	} else {
		fmt.Fprintln(os.Stderr, "Could not list test cases!")
		fmt.Fprintln(os.Stderr, string(result))

		os.Exit(1)
	}
}
