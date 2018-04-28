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
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseListCmd)
}

func runTestCaseList(cmd *cobra.Command, args []string) {
	client := NewClient()

	status, result, err := client.ListTestCases(testCaseListOpts.Organisation)
	if err != nil {
		log.Fatal(err)
	}

	if !status {
		fmt.Fprintln(os.Stderr, "Could not list test cases!")
		fmt.Fprintln(os.Stderr, string(result))

		os.Exit(1)
	}

	if rootOpts.OutputFormat == "json" {
		fmt.Println(string(result))
		return
	}

	items, err := testcase.Unmarshal(bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items.TestCases {
		if rootOpts.OutputFormat == "human" {
			fmt.Printf("%s (ID: %s)\n", item.Name, item.ID)
		} else {
			fmt.Printf("%s\n", item.Name)
		}
	}
}
