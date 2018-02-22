package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	// testRunListCmd represents the calllog command
	testRunListCmd = &cobra.Command{
		Use:   "list <test-case-ref>",
		Short: "List of completed test runs",
		Long:  `List of completed test runs.`,
		Run:   testRunList,
	}
)

func init() {
	TestRunCmd.AddCommand(testRunListCmd)
}

func testRunList(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Expecting exactly one argument: Test Case Reference")
	}

	client := NewClient()

	_, err := client.TestRunList(args[0])
	if err != nil {
		log.Fatal(err)
	}
}
