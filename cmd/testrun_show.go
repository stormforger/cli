package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// testRunShowCmd represents the calllog command
	testRunShowCmd = &cobra.Command{
		Use:   "show <test-run-ref>",
		Short: "Show test run details",
		Long:  `Show test run details.`,
		Run:   testRunShow,
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

	testRun := lookupTestRun(*client, args[0])

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
