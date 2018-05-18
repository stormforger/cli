package cmd

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	// dumpCmd represents the dump command
	dumpCmd = &cobra.Command{
		Use:   "dump <test-run-ref>",
		Short: "Fetch the the test runs dump file (request/response log)",
		Long: `Will fetch the the test runs dump file (request/response log).

The call log contains: FIXME
	* time (epoch in seconds)
	* HTTP verb
	* HTTP host
	* request path
	* HTTP Status Code
	* response size (in Bytes)
	* duration (in ms)
	* request tag`,
		Run:              runTestRunDumpOptions,
		PersistentPreRun: ensureTestRunDumpOptions,
	}

	dumpOpts struct {
		Type       string
		OutputFile string
	}
)

func init() {
	TestRunCmd.AddCommand(dumpCmd)

	dumpCmd.Flags().StringVar(&dumpOpts.OutputFile, "file", "-", "save logs to file or '-' for stdout")
}

func ensureTestRunDumpOptions(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Expecting exactly one argument: Test Run Reference")
	}
}

func runTestRunDumpOptions(cmd *cobra.Command, args []string) {
	client := NewClient()

	testRunUID := getTestRunUID(*client, args[0])

	reader, err := client.TestRunDump(testRunUID)
	if err != nil {
		log.Fatal(err)
	}

	if dumpOpts.OutputFile == "-" {
		_, err = io.Copy(os.Stdout, reader)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		file, err := os.Create(dumpOpts.OutputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		defer reader.Close()

		_, err = io.Copy(file, reader)
		if err != nil {
			log.Fatal(err)
		}
	}
}
