package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/stormforger/cli/api"
	"github.com/stormforger/cli/api/filefixture"
	"github.com/stormforger/cli/api/organisation"
	"github.com/stormforger/cli/api/testcase"
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
	fileFixtureListResponse, err := client.ListFileFixture(organization)
	if err != nil {
		log.Fatal(err)
	}

	fileFixtures, err := filefixture.UnmarshalFileFixtures(bytes.NewReader(fileFixtureListResponse))
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
	response, err := client.ListOrganisations()
	if err != nil {
		log.Fatal(err)
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

func watchTestRun(testRunUID string, maxWatchTime float64) {
	client := NewClient()

	started := time.Now()

	for true {
		runningSince := time.Now().Sub(started).Seconds()

		testRun, response, err := client.TestRunWatch(testRunUID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(response)

		if !testRunOkay(&testRun) {
			os.Exit(1)
		}

		if testRunSuccess(&testRun) {
			os.Exit(0)
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

		_, result, err := client.ListTestCases(organisationUID)
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
