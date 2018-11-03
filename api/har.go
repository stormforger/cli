package api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

// Har converts the given HAR archive file into
// a StormForger test case definition
func (c *Client) Har(fileName string, data io.Reader) (string, error) {
	extraParams := map[string]string{}

	input, err := ioutil.ReadAll(data)
	if err != nil {
		log.Fatal(err)
	}

	if !json.Valid(input) {
		return "", fmt.Errorf("given HAR is not valid JSON")
	}

	var buf bytes.Buffer
	harFileGzip := gzip.NewWriter(&buf)
	io.Copy(harFileGzip, bytes.NewReader(input))
	harFileGzip.Close()

	req, err := fileUploadRequest(c.APIEndpoint+"/har", "POST", extraParams, "har_file", fileName, "application/gzip", bytes.NewReader(buf.Bytes()))
	if err != nil {
		log.Fatal(err)
	}

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
