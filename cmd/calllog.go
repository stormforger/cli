package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// calllogCmd represents the calllog command
var calllogCmd = &cobra.Command{
	Use:   "clog <test-run-ref>",
	Short: "Fetch up to 10k lines of the test runs call log",
	Long: `Will fetch up to 10k lines of the test runs call log.

The call log contains:
  * time (epoch in seconds)
  * HTTP verb
  * HTTP host
  * request path
  * HTTP Status Code
  * response size (in Bytes)
  * duration (in ms)
  * request tag`,
	Run: run,
}

func init() {
	TestRunCmd.AddCommand(calllogCmd)
}

func run(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Expect exactly one argument: Test Run Reference")
	}

	client := NewClient()

	result, err := client.TestRunCallLogPreview(args[0])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(result)
}
