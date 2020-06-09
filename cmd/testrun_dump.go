package cmd

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	dumpCmd = &cobra.Command{
		Use:   "dump <test-run-ref>",
		Short: "Fetch traffic dump (if available)",
		Long: `Will fetch the test run's traffic dump file.

If enabled for the given test run, the traffic dump will
contain a "close to the wire" dump of all requests and responses.

The main purpose is for debugging as the traffic dump is
in no way analyzed or processed.
`,
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

	dumpCmd.Flags().StringVar(&dumpOpts.OutputFile, "file", "-", "save traffic dump to file or '-' for stdout")
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
	defer reader.Close()

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
