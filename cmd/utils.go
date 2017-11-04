package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
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

func printPrettyJson(message string) {
	prettyJson := prettyFormatJson(message)

	_, err := prettyJson.WriteTo(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func prettyFormatJson(message string) (out bytes.Buffer) {
	err := json.Indent(&out, []byte(message), "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return out
}
