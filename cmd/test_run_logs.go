package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	// calllogCmd represents the calllog command
	calllogCmd = &cobra.Command{
		Use:   "logs <test-run-ref>",
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

	logOpts struct {
		Type string
	}
)

func init() {
	TestRunCmd.AddCommand(calllogCmd)

	calllogCmd.Flags().StringVar(&logOpts.Type, "type", "request", "type of logs")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Expect exactly one argument: Test Run Reference")
	}

	if logOpts.Type != "request" {
		log.Fatal(fmt.Sprintf("Unsupported log type %s", logOpts.Type))
	}

	client := NewClient()

	result, err := client.TestRunCallLogPreview(args[0])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(result)
}
