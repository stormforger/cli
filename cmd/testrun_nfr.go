package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	testRunNfrCmd = &cobra.Command{
		Use:   "nfr <test-run-ref> <requirements_file>",
		Short: "Check test run against NFR",
		Long:  `Check test run against non-functional requirements.`,
		Run:   testRunNfrRun,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				log.Fatal("Missing arguments; test run reference and NFR requirements file")
			}

			if len(args) > 2 {
				log.Fatal("Too many arguments")
			}
		},
	}
)

func init() {
	TestRunCmd.AddCommand(testRunNfrCmd)
}

func testRunNfrRun(cmd *cobra.Command, args []string) {
	client := NewClient()

	testRunUID := getTestRunUID(*client, args[0])

	fileName, file, err := readFromStdinOrReadFromArgument(args[1], "nfr.yml")
	if err != nil {
		log.Fatal(err)
	}

	runNfrCheck(*client, testRunUID, fileName, file)
}
