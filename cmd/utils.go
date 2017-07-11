package cmd

import (
	"io"
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
