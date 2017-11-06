package api

import (
	"io"
	"io/ioutil"
	"log"
)

// Har converts the given HAR archive file into
// a StormForger test case definition
func (c *Client) Har(fileName string, data io.Reader) (string, error) {
	// TODO how to pass options, like --skip-assets here?
	//      defining a struct maybe, but where?
	//      finally: add options here
	extraParams := map[string]string{}

	req, err := fileUploadRequest(c.APIEndpoint+"/har", "POST", extraParams, "har_file", fileName, "application/json", data)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(body), nil
}
