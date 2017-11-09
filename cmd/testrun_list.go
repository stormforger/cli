package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	// testRunListCmd represents the calllog command
	testRunListCmd = &cobra.Command{
		Use:   "list <test-case-ref>",
		Short: "bla",
		Long:  `bla`,
		Run:   testRunList,
	}
)

func init() {
	TestRunCmd.AddCommand(testRunListCmd)
}

func testRunList(cmd *cobra.Command, args []string) {
	client := NewClient()

	_, err := client.TestRunList(args[0])
	if err != nil {
		log.Fatal(err)
	}
}
