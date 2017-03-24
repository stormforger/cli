package api

import (
	"io/ioutil"
	"log"
)

// Har converts the given HAR archive file into
// a StormForger test case definition
func (c *Client) Har(file string) (string, error) {
	// TODO how to pass options, like --skip-assets here?
	//      defining a struct maybe, but where?
	//      finally: add options here
	extraParams := map[string]string{}

	req, err := newfileUploadRequest(c.APIEndpoint+"/har", extraParams, "har_file", file)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	return string(body), nil
}
