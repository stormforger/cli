package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	// testCaseGetCmd represents the testCaseGetCmd command
	testCaseGetCmd = &cobra.Command{
		Use:     "get <test-case-ref> [file]",
		Aliases: []string{"download"},
		Short:   "Download test case definition",
		Long: `Download the JavaScript definition of test case

<test-case-ref> can be 'organisation-name/test-case-name' or 'test-case-uid'.

[file] can be '-' to output to STDOUT (default) or path to file
to write to.
`,
		Run: runTestCaseGet,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("Missing argument: test case reference")
			}

			if len(args) > 2 {
				log.Fatal("Too many arguments")
			}
		},
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseGetCmd)
}

func runTestCaseGet(cmd *cobra.Command, args []string) {
	client := NewClient()

	testCaseUID := mustLookupTestCase(client, args[0])

	success, response, err := client.DownloadTestCaseDefinition(testCaseUID)
	if err != nil {
		log.Fatal(err)
	}

	if !success {
		log.Fatalf("Test case definition could not be downloaded!\n%s\n", string(response))
	}

	if len(args) == 2 && args[1] != "-" {
		err := ioutil.WriteFile(args[1], response, 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_, err := io.Copy(os.Stdout, bytes.NewReader(response))
		if err != nil {
			log.Fatal(err)
		}
	}
}
