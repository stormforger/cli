package main

import (
	"strings"
	"io/ioutil"

	"github.com/stormforger/cli/api"
)

func main() {
	client := api.NewClient("http://api.stormforger.com", readJwt())

	client.Ping()
}

func readJwt() string {
	jwt, _ := ioutil.ReadFile("./.stormforger_jwt")

	return strings.TrimSpace(string(jwt))
}
