package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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
