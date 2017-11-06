package cmd

import (
	"bytes"
	"log"

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
			if testCaseListOpts.Organisation == "" {
				log.Fatal("Missing organization flag")
			}
		},
	}

	testCaseListOpts struct {
		Organisation string
		Name         string
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseListCmd)

	testCaseListCmd.PersistentFlags().StringVarP(&testCaseListOpts.Organisation, "organization", "o", "", "Name of the organization")
}

func runTestCaseList(cmd *cobra.Command, args []string) {
	client := NewClient()

	result, err := client.ListTestCases(testCaseListOpts.Organisation)
	if err != nil {
		log.Fatal(err)
	}

	testcase.ShowNames(bytes.NewReader(result))
}
