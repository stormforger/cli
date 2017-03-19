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
		Type    string
		Preview bool
	}
)

func init() {
	TestRunCmd.AddCommand(calllogCmd)

	calllogCmd.Flags().StringVar(&logOpts.Type, "type", "request", "type of logs")
	calllogCmd.Flags().BoolVarP(&logOpts.Preview, "preview", "p", true, "Preview of logs")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("Missing argument: Test Run Reference")
	}

	if len(args) > 2 {
		log.Fatal("Too many arguments")
	}

	if logOpts.Type != "request" {
		log.Fatal(fmt.Sprintf("Unsupported log type %s", logOpts.Type))
	}

	client := NewClient()

	reader, err := client.TestRunCallLog(args[0], logOpts.Preview)
	if err != nil {
		log.Fatal(err)
	}

	if len(args) == 1 {
		io.Copy(os.Stdout, reader)
	} else {
		file, err := os.Create(args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		defer reader.Close()

		io.Copy(file, reader)
	}
}
