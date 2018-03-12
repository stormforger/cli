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
		Use:     "list <organization-ref>",
		Aliases: []string{"ls"},
		Short:   "List test case for a given organization",
		Run:     runTestCaseList,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				log.Fatal("Too many arguments")
			}

			if len(args) < 1 {
				log.Fatal("Missing organization")
			}

			testCaseListOpts.Organisation = lookupOrganisationUID(*NewClient(), args[0])
			if testCaseListOpts.Organisation == "" {
				log.Fatal("Missing organization")
			}
		},
	}

	testCaseListOpts struct {
		Organisation string
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
