package cmd

import (
	"log"

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

	_, err := client.TestRunShow(args[0])
	if err != nil {
		log.Fatal(err)
	}
}
