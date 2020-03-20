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
		Short: "Fetch call log (request log)",
		Long: `Will fetch the test run's call log (request log).

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
		Run:              runTestRunLogs,
		PersistentPreRun: ensureTestRunLogs,
	}

	logOpts struct {
		Type       string
		Full       bool
		OutputFile string
	}
)

func init() {
	TestRunCmd.AddCommand(calllogCmd)

	calllogCmd.Flags().StringVar(&logOpts.Type, "type", "request", "type of logs (request or user)")
	calllogCmd.Flags().BoolVarP(&logOpts.Full, "full", "f", false, "download full logs")
	calllogCmd.Flags().StringVar(&logOpts.OutputFile, "file", "-", "save logs to file or '-' for stdout")
}

func ensureTestRunLogs(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Expecting exactly one argument: Test Run Reference")
	}

	if logOpts.Type != "request" && logOpts.Type != "user" {
		log.Fatal(fmt.Sprintf("Unsupported log type %s", logOpts.Type))
	}
}

func runTestRunLogs(cmd *cobra.Command, args []string) {
	client := NewClient()

	testRunUID := getTestRunUID(*client, args[0])

	var reader io.ReadCloser
	var err error

	switch logOpts.Type {
	case "request":
		reader, err = client.TestRunCallLog(testRunUID, !logOpts.Full)
		if err != nil {
			log.Fatal(err)
		}
	case "user":
		reader, err = client.TestRunUserLog(testRunUID)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer reader.Close()

	if logOpts.OutputFile == "-" {
		_, err = io.Copy(os.Stdout, reader)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		file, err := os.Create(logOpts.OutputFile)
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(file, reader)
		if err != nil {
			log.Fatal(err)
		}

		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}
