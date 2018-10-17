package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/stormforger/cli/api"
	"github.com/stormforger/cli/api/filefixture"
	"github.com/stormforger/cli/api/organisation"
	"github.com/stormforger/cli/api/testcase"
	"github.com/stormforger/cli/api/testrun"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// FindFixtureByName fetches a FileFixture from a given
// organization.
func findFixtureByName(client api.Client, organization string, name string) *filefixture.FileFixture {
	success, result, err := client.ListFileFixture(organization)
	if err != nil {
		log.Fatal(err)
	}

	if !success {
		fmt.Fprintln(os.Stderr, "Could not lookup data source!")
		fmt.Fprintln(os.Stderr, string(result))

		os.Exit(1)
	}

	fileFixtures, err := filefixture.Unmarshal(bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}

	fileFixture := fileFixtures.FindByName(name)
	if fileFixture.ID == "" {
		log.Fatalf("Data source %s not found in organization %s!", name, organization)
	}

	return &fileFixture
}

// findOrganisationByName fetches a FileFixture from a given
// organization.
func findOrganisationByName(client api.Client, name string) *organisation.Organisation {
	status, response, err := client.ListOrganisations()
	if err != nil {
		log.Fatal(err)
	}
	if !status {
		log.Fatal(string(response))
	}

	organisations, err := organisation.Unmarshal(bytes.NewReader(response))
	if err != nil {
		log.Fatal(err)
	}

	organisation := organisations.FindByNameOrUID(name)
	if organisation.ID == "" {
		log.Fatalf("Organisation %s not found!", name)
	}

	return &organisation
}

func readFromStdinOrReadFromArgument(args []string, defaultFileName string, argPos int) (fileName string, reader io.Reader, err error) {
	fileName = defaultFileName

	argument := args[argPos]

	if argument == "-" {
		reader = os.Stdin
	} else {
		fileName = filepath.Base(argument)
		reader, err = os.OpenFile(argument, os.O_RDONLY, 0755)
		if err != nil {
			return "", nil, err
		}
	}

	return fileName, reader, err
}

func readFromStdinOrReadFirstArgument(args []string, defaultFileName string) (fileName string, reader io.Reader, err error) {
	return readFromStdinOrReadFromArgument(args, defaultFileName, 0)
}

func readTestCaseFromStdinOrReadFromArgument(args []string, defaultFileName string, argPos int) (fileName string, reader io.Reader, err error) {
	fileName, testCaseFile, err := readFromStdinOrReadFromArgument(args, defaultFileName, argPos)
	if err != nil {
		log.Fatal(err)
	}

	basePath := ""
	if f := args[argPos]; f != "-" {
		basePath = filepath.Dir(f)
	} else {
		var err error
		basePath, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		re := regexp.MustCompile("//#include (.+?)$")

		scanner := bufio.NewScanner(testCaseFile)
		for scanner.Scan() {
			line := scanner.Text()

			if match := re.FindStringSubmatch(line); match != nil {
				includeFile := strings.TrimSpace(match[1])

				if !filepath.IsAbs(includeFile) {
					includeFile = filepath.Join(basePath, includeFile)
				}

				f, err := os.Open(includeFile)
				if err != nil {
					log.Fatal(err)
				}

				pw.Write([]byte("// == start include (" + includeFile + ")\n"))
				io.Copy(pw, f)
				pw.Write([]byte("// == end include (" + includeFile + ")\n"))
			} else {
				pw.Write([]byte(line + "\n"))
			}
		}
	}()

	return fileName, pr, err
}

func printPrettyJSON(message string) {
	prettyJSON := prettyFormatJSON(message)

	_, err := prettyJSON.WriteTo(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func prettyFormatJSON(message string) (out bytes.Buffer) {
	err := json.Indent(&out, []byte(message), "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return out
}

func readOrganisationUIDFromFile() string {
	content, err := ioutil.ReadFile(".stormforger-organisation")
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(content))
}

func watchTestRun(testRunUID string, maxWatchTime float64, outputFormat string) {
	client := NewClient()
	started := time.Now()
	first := true
	testStarted := false
	testEnded := false

	for true {
		runningSince := time.Now().Sub(started).Seconds()

		testRun, response, err := client.TestRunWatch(testRunUID)
		if err != nil {
			log.Fatal(err)
		}

		if first {
			first = false

			if testRunSuccess(&testRun) {
				if outputFormat != "json" {
					fmt.Printf("Test Run %s is finished.\n", testRun.ID)
				}

				return
			}

			if outputFormat == "json" {
				fmt.Println(response)
			} else {
				formattedStartedAt := "not yet"
				formattedEstimatedEnd := "n/a"

				if testRun.StartedAt != "" {
					formattedStartedAt = humanize.Time(parseTime(testRun.StartedAt))
					testStarted = true
				}
				if testRun.EstimatedEnd != "" {
					formattedEstimatedEnd = humanize.Time(parseTime(testRun.EstimatedEnd))
				}

				fmt.Printf("[status] Test Run: %s started %s (est. end %s)\n", testRun.ID, formattedStartedAt, formattedEstimatedEnd)
			}
		}

		if outputFormat == "json" {
			fmt.Println(response)
		} else {
			switch testRun.State {
			case "running":
				if !testStarted {
					formattedStartedAt := "not yet"
					formattedEstimatedEnd := "n/a"

					if testRun.StartedAt != "" {
						formattedStartedAt = humanize.Time(parseTime(testRun.StartedAt))
					}
					if testRun.EstimatedEnd != "" {
						formattedEstimatedEnd = humanize.Time(parseTime(testRun.EstimatedEnd))
					}

					fmt.Printf("[status] Test Run: %s started %s (est. end %s)\n", testRun.ID, formattedStartedAt, formattedEstimatedEnd)

					testStarted = true
				}
				fmt.Printf("[%s] Progress: %d%%\n", testRun.State, testRun.Progress)
			default:
				if testStarted && !testEnded {
					fmt.Printf("[status] Test run ended...\n")
					testEnded = true
				}

				fmt.Printf("[%s]\n", testRun.State)
			}

		}

		if !testRunOkay(&testRun) {
			os.Exit(1)
		}

		if testRunSuccess(&testRun) {
			return
		}

		if int(maxWatchTime) > 0 && int(runningSince) > int(maxWatchTime) {
			os.Exit(2)
		}

		time.Sleep(5 * time.Second)
	}
}

func findFirstNonEmpty(candidates []string) string {
	for _, item := range candidates {
		if item != "" {
			return item
		}
	}

	return ""
}

func lookupOrganisationUID(client api.Client, input string) string {
	organisation := findOrganisationByName(client, input)
	if organisation.ID == "" {
		log.Fatalf("Organisation %s not found", input)
	}

	return organisation.ID
}

func lookupTestCase(client api.Client, input string) string {
	segments := strings.Split(input, "/")
	nameOrUID := input

	if len(segments) == 2 {
		organisationNameOrUID := segments[0]
		nameOrUID = segments[1]

		organisationUID := lookupOrganisationUID(client, organisationNameOrUID)

		_, result, err := client.ListTestCases(organisationUID, "all")
		if err != nil {
			log.Fatal(err)
		}

		testCases, err := testcase.Unmarshal(bytes.NewReader(result))
		if err != nil {
			log.Fatal(err)
		}

		testCase := testCases.FindByNameOrUID(nameOrUID)
		if testCase.ID == "" {
			log.Fatalf("Test case %s not found", nameOrUID)
		}

		return testCase.ID
	}

	return nameOrUID
}

func getTestRunUID(client api.Client, input string) string {
	testRunParts := api.ExtractTestRunResources(input)

	if testRunParts.UID != "" {
		return testRunParts.UID
	} else if testRunParts.Organisation == "" || testRunParts.TestCase == "" {
		log.Fatal("Invalid test run reference provided! Consult with --help to learn more.")
	}

	result := fetchTestRun(client, input)
	testRun, err := testrun.UnmarshalSingle(bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}

	return testRun.ID
}

func fetchTestRun(client api.Client, input string) []byte {
	testRunParts := api.ExtractTestRunResources(input)

	if testRunParts.UID != "" {
		status, response, err := client.FetchTestRun(testRunParts.UID)
		if err != nil {
			log.Fatal(err)
		}
		if !status {
			log.Fatalf("Could not load test run %s", testRunParts.UID)
		}

		return response
	} else if testRunParts.Organisation == "" || testRunParts.TestCase == "" {
		log.Fatal("Invalid test run reference provided! Consult with --help to learn more.")
	}

	status, response, err := client.LookupAndFetchResource("test_run", input)
	if err != nil {
		log.Fatal(err)
	}
	if !status {
		log.Fatalf("Test Run %s not found", input)
	}

	return response
}

func parseTime(subject string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339Nano, subject)
	if err != nil {
		log.Fatal(err)
	}

	return parsedTime
}

func convertToLocalTZ(timeToConvert time.Time) time.Time {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		log.Fatal(err)
	}

	return timeToConvert.In(loc)
}

func runNfrCheck(client api.Client, testRunUID string, fileName string, nfrData io.Reader) {
	status, result, err := client.TestRunNfrCheck(testRunUID, fileName, nfrData)
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

	anyFails := displayNfrResult(items)

	if anyFails {
		os.Exit(1)
	}
}

func displayNfrResult(items testrun.NfrResultList) bool {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	redBg := color.New(color.BgRed).Add(color.FgWhite).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()
	whiteBold := color.New(color.FgWhite, color.Bold).SprintFunc()

	checkStatus := ""
	anyFails := false
	for _, item := range items.NfrResults {
		if !item.Disabled {
			actualSubject := ""

			if item.SubjectAvailable {
				if item.Success {
					checkStatus = green("\u2713")
					actualSubject = fmt.Sprintf("was %s", item.SubjectWithUnit())
				} else {
					anyFails = true
					checkStatus = red("\u2717")
					actualSubject = fmt.Sprintf("but actually was %s", item.SubjectWithUnit())
				}
			} else {
				checkStatus = whiteBold("?")
				actualSubject = fmt.Sprintf("was %s", whiteBold("not available"))
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

	if !anyFails {
		fmt.Printf(green("\nAll checks passed!\n"))
	} else {
		fmt.Printf(red("\nYou have failing checks!\n"))
	}

	return anyFails
}
