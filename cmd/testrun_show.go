package cmd

import (
	"bytes"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/testrun"
)

var (
	// testRunShowCmd represents the calllog command
	testRunShowCmd = &cobra.Command{
		Use:   "show <test-run-ref>",
		Short: "Show test run details",
		Long:  `Show test run details.`,
		Run:   testRunShow,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatal("Expect exactly one argument: test run reference!")
			}
		},
	}

	testRunShowOpts struct {
		Type       string
		Full       bool
		OutputFile string
	}
)

func init() {
	TestRunCmd.AddCommand(testRunShowCmd)
}

func testRunShow(cmd *cobra.Command, args []string) {
	client := NewClient()

	result := fetchTestRun(*client, args[0])

	if rootOpts.OutputFormat == "json" {
		fmt.Println(string(result))
		return
	}

	testRun, err := testrun.UnmarshalSingle(bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s (%s, %s)\n", testRun.Scope, testRun.State, testRun.ID)
	if testRun.Title != "" {
		fmt.Printf("Title   %s\n", testRun.Title)
	}
	fmt.Printf("Started %s\n", testRun.StartedAt)
	fmt.Printf("Ended   %s\n", testRun.EndedAt)
	if testRun.Notes != "" {
		fmt.Printf("%s\n", testRun.Notes)
	}
}
