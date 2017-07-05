package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	// testCaseValidateCmd represents the testCaseValidate command
	testCaseValidateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Upload a test case definition JavaScript and validate it",
		Long:  `Upload a test case definition JavaScript and validate it.`,
		Run:   runTestCaseValidate,
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseValidateCmd)
}

func runTestCaseValidate(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		var testCaseFile string

		// FIXME this is the same as in har.go. Can we extract and generalize this?
		if args[0] == "-" {
			fileInput := readFromStdin()
			tmpFile, err := ioutil.TempFile(os.TempDir(), "forge-test-case")
			if err != nil {
				log.Fatal(err)
			}
			defer os.Remove(tmpFile.Name())

			testCaseFile = tmpFile.Name()

			// TODO what/why this compact syntax with ;?
			if _, err := tmpFile.Write([]byte(fileInput)); err != nil {
				log.Fatal(err)
			}

			if err := tmpFile.Close(); err != nil {
				log.Fatal(err)
			}

		} else {
			// FIXME check if file exists here?
			testCaseFile = args[0]
		}

		client := NewClient()

		_, message, errValidation := client.TestCaseValidate(testCaseFile)
		if errValidation != nil {
			log.Fatal(errValidation)
		}

		var out bytes.Buffer
		err := json.Indent(&out, []byte(message), "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		out.WriteTo(os.Stdout)
	} else {
		log.Fatal("Missing argument; test case file or - to read from stdin")
	}
}
