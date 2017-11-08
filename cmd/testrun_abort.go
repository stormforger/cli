package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	// testRunAbortCmd represents the calllog command
	testRunAbortCmd = &cobra.Command{
		Use:   "abort <test-run-id>",
		Short: "Abort the given running test",
		Long:  `Abort the given running test.`,
		Run:   testRunAbort,
	}
)

func init() {
	TestRunCmd.AddCommand(testRunAbortCmd)
}

func testRunAbort(cmd *cobra.Command, args []string) {
	client := NewClient()

	_, response, err := client.TestRunAbort(args[0])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}
