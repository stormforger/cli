package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/testrun"
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

	fileName, file, err := readFromStdinOrReadFromArgument(args, "nfr.yml", 1)
	if err != nil {
		log.Fatal(err)
	}

	status, result, err := client.TestRunNfrCheck(testRunUID, fileName, file)
	if err != nil {
		log.Fatal(err)
	}

	if rootOpts.OutputFormat == "json" {
		fmt.Println(string(result))
		return
	}

	if !status {
		log.Fatalf("Could not perform test run NFR checks...\n%s", result)
	}

	items, err := testrun.UnmarshalNfrResults(bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	redBg := color.New(color.BgRed).Add(color.FgWhite).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()

	checkStatus := ""
	anyFails := false
	for _, item := range items.NfrResults {
		if !item.Disabled {
			actualSubject := ""
			if item.Success {
				checkStatus = green("\u2713")
				actualSubject = fmt.Sprintf("was %s", item.SubjectWithUnit())
			} else {
				anyFails = true
				checkStatus = red("\u2717")
				actualSubject = fmt.Sprintf("but actually was %s", item.SubjectWithUnit())
			}

			filter := ""
			if item.Filter != "null" && item.Filter != "" {
				filter = " (where: " + item.Filter + ")"
			}

			fmt.Printf(
				"%s %s expected to be %s; %s (%s)%s\n",
				checkStatus,
				item.Metric,
				item.ExpectationWithUnit(),
				actualSubject,
				item.Type,
				filter,
			)
		} else {
			fmt.Printf(
				"%s %s %s expected to be %s (%s)\n",
				white("?"),
				redBg("DISABLED"),
				item.Metric,
				item.ExpectationWithUnit(),
				item.Type,
			)
		}
	}

	if anyFails {
		os.Exit(1)
	}
}
