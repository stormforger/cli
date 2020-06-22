package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
)

// Har converts the given HAR archive file into
// a StormForger test case definition
func (c *Client) Har(fileName string, data io.Reader) (string, error) {
	extraParams := url.Values{}

	input, err := ioutil.ReadAll(data)
	if err != nil {
		return "", fmt.Errorf("%s: %v", fileName, err)
	}

	if !json.Valid(input) {
		return "", fmt.Errorf("%s: given HAR is not valid JSON", fileName)
	}

	req, err := fileUploadRequest(c.APIEndpoint+"/har", "POST", extraParams, "har_file", fileName, "application/octet-stream", bytes.NewReader(input))
	if err != nil {
		return "", err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
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
