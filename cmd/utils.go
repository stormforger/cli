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
)

func readFromStdinOrReadFirstArgument(args []string, defaultFileName string) (fileName string, reader io.Reader, err error) {
	fileName = defaultFileName

	if args[0] == "-" {
		reader = os.Stdin
	} else {
		fileName = filepath.Base(args[0])
		reader, err = os.OpenFile(args[0], os.O_RDONLY, 0755)
		if err != nil {
			return "", nil, err
		}
	}

	return fileName, reader, err
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

func watchTestRun(testRunUID string, maxWatchTime int) {
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

		if maxWatchTime > 0 && int(runningSince) > maxWatchTime {
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
