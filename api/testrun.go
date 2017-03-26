package api

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"strings"
)

type testRunResources struct {
	uid          string
	organisation string
	testCase     string
	sequenceID   string
}

// TestRunCallLog will download the first 10k lines
// of the test run's call log
func (c *Client) TestRunCallLog(pathID string, preview bool) (io.ReadCloser, error) {
	testRun := extractResources(pathID)

	var path string
	if testRun.uid == "" {
		path = "/beta/test_cases/" + testRun.organisation + "/" + testRun.testCase + "/test_runs/" + testRun.sequenceID
	} else {
		path = "/beta/t/" + testRun.uid
	}

	path += "/call_log"

	if preview {
		path += "?preview=true"
	}

	req, err := http.NewRequest("GET", c.APIEndpoint+path, nil)
	if err != nil {
		return nil, err
	}

	// TODO how to set these on all requests?
	c.addDefaultHeaders(req)

	req.Header.Set("Accept-Encoding", "gzip")

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("could not download call log")
	}

	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}

		defer reader.Close()
	default:
		reader = response.Body
	}

	return reader, nil
}

// extractResources will try to extract information to the
// given test run based on a "reference".
//
// Currently as "reference" a part of the forge URL is used.
// This contains the organisation, test case and the sequence
// id of the test run. Example: "foo/demo/test_runs/19"
func extractResources(ref string) testRunResources {
	segments := strings.Split(ref, "/")

	if len(segments) == 4 && segments[2] == "test_runs" {
		return testRunResources{
			organisation: segments[0],
			testCase:     segments[1],
			sequenceID:   segments[3],
		}
	}

	if len(segments) == 3 {
		return testRunResources{
			organisation: segments[0],
			testCase:     segments[1],
			sequenceID:   segments[2],
		}
	}

	if len(segments) == 1 {
		return testRunResources{
			uid: segments[0],
		}
	}

	return testRunResources{}
}