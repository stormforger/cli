package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		fileName, testCaseFile, err := readFromStdinOrReadFirstArgument(args, "test_case.js")
		if err != nil {
			log.Fatal(err)
		}

		client := NewClient()

		success, message, errValidation := client.TestCaseValidate(fileName, testCaseFile)
		if errValidation != nil {
			log.Fatal(errValidation)
		}

		if success {
			os.Exit(0)
		}

		var out bytes.Buffer
		err = json.Indent(&out, []byte(message), "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		_, err = out.WriteTo(os.Stdout)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println()
		os.Exit(1)
	} else {
		log.Fatal("Missing argument; test case file or - to read from stdin")
	}
}
