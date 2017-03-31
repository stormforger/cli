package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	// calllogCmd represents the calllog command
	calllogCmd = &cobra.Command{
		Use:   "logs <test-run-ref>",
		Short: "Fetch the the test runs call log (request log)",
		Long: `Will fetch the the test runs call log (request log).

By default, you will get the first 10k lines. Using --full you
will download the entire request log.

The call log contains:
  * time (epoch in seconds)
  * HTTP verb
  * HTTP host
  * request path
  * HTTP Status Code
  * response size (in Bytes)
  * duration (in ms)
  * request tag`,
		Run:              runTestRunLogsOptions,
		PersistentPreRun: ensureTestRunLogsOptions,
	}

	logOpts struct {
		Type       string
		Full       bool
		OutputFile string
	}
)

func init() {
	TestRunCmd.AddCommand(calllogCmd)

	calllogCmd.Flags().StringVar(&logOpts.Type, "type", "request", "type of logs")
	calllogCmd.Flags().BoolVarP(&logOpts.Full, "full", "f", false, "download full logs")
	calllogCmd.Flags().StringVar(&logOpts.OutputFile, "output", "-", "save logs to file or '-' for stdout")
}

func ensureTestRunLogsOptions(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Expecting exactly one argument: Test Run Reference")
	}

	if logOpts.Type != "request" {
		log.Fatal(fmt.Sprintf("Unsupported log type %s", logOpts.Type))
	}
}

func runTestRunLogsOptions(cmd *cobra.Command, args []string) {
	client := NewClient()

	reader, err := client.TestRunCallLog(args[0], !logOpts.Full)
	if err != nil {
		log.Fatal(err)
	}

	if logOpts.OutputFile == "-" {
		io.Copy(os.Stdout, reader)
	} else {
		file, err := os.Create(logOpts.OutputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		defer reader.Close()

		io.Copy(file, reader)
	}
}
