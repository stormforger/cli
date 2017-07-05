package cmd

import (
	"io/ioutil"
	"log"
	"os"
)

func readFromStdin() string {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	return string(bytes)
}
