package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/stormforger/cli/api"
	"github.com/stormforger/cli/api/filefixture"
	"github.com/stormforger/cli/api/organisation"
	"github.com/stormforger/cli/api/testcase"
	"github.com/stormforger/cli/api/testrun"
	"github.com/stormforger/cli/internal/esbundle"
)

const bundlingHelpInfo = `Bundling
--------
If you use the .mjs file extension, this command will automatically bundle your
JavaScript file using ECMAScript modules. See 'forge test-case build' for more details.
`

// FindFixtureByName fetches a FileFixture from a given Organisation.
func findFixtureByName(client api.Client, orga string, name string) *filefixture.FileFixture {
	success, result, err := client.ListFileFixture(orga)
	if err != nil {
		log.Fatalf("ListFileFixtures failed: %v", err)
	}

	if !success {
		fmt.Fprintln(os.Stderr, "Could not lookup data source!")
		fmt.Fprintln(os.Stderr, string(result))

		os.Exit(1)
	}

	fileFixtures, err := filefixture.Unmarshal(bytes.NewReader(result))
	if err != nil {
		log.Fatalf("Unmarshal failed: %v", err)
	}

	fileFixture := fileFixtures.FindByName(name)
	if fileFixture.ID == "" {
		log.Fatalf("Data source %s not found in organisation %s!", name, orga)
	}

	return &fileFixture
}

// findOrganisationByName fetches a FileFixture from a given Organisation.
func findOrganisationByName(client *api.Client, name string) *organisation.Organisation {
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

// readFromStdinOrReadFromArgument returns the filename and a filereader for fileArg.
// If fileArg matches "-", the defaultFileName and stdin is returned
func readFromStdinOrReadFromArgument(fileArg, defaultFileName string) (fileName string, reader io.Reader, err error) {
	if fileArg == "-" {
		fileName = defaultFileName
		reader = os.Stdin
	} else {
		fileName = filepath.Base(fileArg)
		reader, err = os.OpenFile(fileArg, os.O_RDONLY, 0755)
		if err != nil {
			return fileName, nil, err
		}
	}

	return fileName, reader, err
}

type testCaseFileBundle struct {
	Name    string
	Content io.Reader
	Mapper  esbundle.SourceMapper
}
type testCaseFileBundler struct {
	Defines map[string]string
}

func (bundler testCaseFileBundler) Bundle(arg, defaultFileName string) (testCaseFileBundle, error) {
	fileName, testCaseFile, err := readFromStdinOrReadFromArgument(arg, defaultFileName)
	if err != nil {
		return testCaseFileBundle{}, err
	}

	if arg != "-" && filepath.Ext(arg) == ".mjs" {
		result, err := esbundle.Bundle(arg, bundler.Defines)
		if err != nil {
			return testCaseFileBundle{}, err
		}
		return testCaseFileBundle{Name: fileName, Content: strings.NewReader(result.CompiledContent), Mapper: result.SourceMapper}, nil
	}

	return testCaseFileBundle{Name: fileName, Content: testCaseFile}, err
}

func watchTestRun(testRunUID string, maxWatchTime float64, outputFormat string) {
	client := NewClient()
	started := time.Now()
	first := true
	testStarted := false
	testEnded := false

	for true {
		runningSince := time.Since(started).Seconds()

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

func lookupOrganisationUID(client *api.Client, input string) string {
	organisation := findOrganisationByName(client, input)
	if organisation.ID == "" {
		log.Fatalf("Organisation %s not found", input)
	}

	return organisation.ID
}

// mustLookupTestCase returns the ID of the test-case for the given input or calls log.Fatal().
func mustLookupTestCase(client *api.Client, input string) string {
	s := lookupTestCase(client, input)
	if s == "" {
		log.Fatalf("Test case for query '%s' not found", input)
	}
	return s
}

func lookupTestCase(client *api.Client, input string) string {
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
		var response struct {
			Status  string
			Message string
			Error   string
		}

		log.Println("Could not perform test-run NFR checks...")
		if err := json.Unmarshal(result, &response); err != nil {
			log.Println(string(result))
		} else {
			if response.Status != "" {
				log.Printf(" Status:\t%s\n", response.Status)
			}
			if response.Message != "" {
				log.Printf(" Message:\t%s\n", response.Message)
			}
			if response.Error != "" {
				log.Printf(" Error:\t%s\n", response.Error)
			}
		}
		os.Exit(1)
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
	var success, total int
	for _, item := range items.NfrResults {
		total++ // we count everything, including disable and unavailable here.
		if !item.Disabled {
			actualSubject := ""

			if item.SubjectAvailable {
				if item.Success {
					success++
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

	fmt.Printf("%d/%d checks passed\n", success, total)
	if !anyFails {
		fmt.Printf(green("\nAll checks passed!\n"))
	} else {
		fmt.Printf(red("\nYou have failing checks!\n"))
	}

	return anyFails
}
