package api

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/jsonapi"
	"github.com/stormforger/cli/api/testrun"
)

// TestRunResources describes infos on a test run
type TestRunResources struct {
	UID          string
	Organisation string
	TestCase     string
	SequenceID   string
}

// TestRunList will list all test runs for a given test case
func (c *Client) TestRunList(testCaseUID string) (bool, []byte, error) {
	path := "/test_cases/" + testCaseUID + "/test_runs"

	status, response, err := c.fetch(path)

	return status, response, err
}

// FetchTestRun will show some basic information on a given
// test run
func (c *Client) FetchTestRun(uid string) (bool, []byte, error) {
	path := "/test_runs/" + uid

	status, response, err := c.fetch(path)

	return status, response, err
}

// TestRunWatch will show some basic information on a given
// test run
func (c *Client) TestRunWatch(uid string) (testrun.TestRun, string, error) {
	path := "/test_runs/" + uid + "/watch"

	req, err := http.NewRequest("GET", c.APIEndpoint+path, nil)
	if err != nil {
		return testrun.TestRun{}, "", err
	}

	response, err := c.doRequestRaw(req)
	if err != nil {
		return testrun.TestRun{}, "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return testrun.TestRun{}, "", err
	}

	if response.StatusCode >= 400 {
		log.Fatal(string(body))
	}

	testRun := new(testrun.TestRun)
	err = jsonapi.UnmarshalPayload(bytes.NewReader(body), testRun)
	if err != nil {
		log.Fatal(err)
	}

	return *testRun, string(body), nil
}

// TestRunCallLog will download the first 10k lines
// of the test run's call log
func (c *Client) TestRunCallLog(pathID string, preview bool) (io.ReadCloser, error) {
	testRun := ExtractTestRunResources(pathID)

	path := "/test_runs/" + testRun.UID + "/call_log"

	if preview {
		path += "?preview=true"
	}

	req, err := http.NewRequest("GET", c.APIEndpoint+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept-Encoding", "gzip")

	response, err := c.doRequestRaw(req)
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

func (c *Client) TestRunDump(pathID string) (io.ReadCloser, error) {
	testRun := ExtractTestRunResources(pathID)

	path := "/test_runs/" + testRun.UID + "/dump"

	req, err := http.NewRequest("GET", c.APIEndpoint+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept-Encoding", "gzip")

	response, err := c.doRequestRaw(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("could not load full dump")
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

// TestRunCreate will send a test case definition (JS) to the API
// to update an existing test case it.
func (c *Client) TestRunCreate(testCaseUID string, title string, notes string) (bool, string, error) {
	payload := bytes.NewBuffer(nil)
	newTestRun := &testrun.TestRun{
		Title: title,
		Notes: notes,
	}
	jsonapi.MarshalOnePayloadEmbedded(payload, newTestRun)

	req, err := http.NewRequest("POST", c.APIEndpoint+"/test_cases/"+testCaseUID+"/test_runs", payload)

	req.Header.Set("Content-Type", "application/json")

	response, err := c.doRequestRaw(req)
	if err != nil {
		return false, "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, "", err
	}

	defer response.Body.Close()

	return response.StatusCode < 400, string(body), nil
}

// TestRunAbort will send a test case definition (JS) to the API
// to update an existing test case it.
func (c *Client) TestRunAbort(testRunUID string) (bool, string, error) {
	req, err := http.NewRequest("POST", c.APIEndpoint+"/test_runs/"+testRunUID+"/abort", nil)

	response, err := c.doRequestRaw(req)
	if err != nil {
		return false, "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, "", err
	}

	defer response.Body.Close()

	return response.StatusCode < 400, string(body), nil
}

// ExtractTestRunResources will try to extract information to the
// given test run based on a "reference".
//
// Currently as "reference" a part of the forge URL is used.
// This contains the organisation, test case and the sequence
// id of the test run. Example: "foo/demo/test_runs/19"
func ExtractTestRunResources(ref string) TestRunResources {
	segments := strings.Split(ref, "/")

	if len(segments) == 4 && segments[2] == "test_runs" {
		return TestRunResources{
			Organisation: segments[0],
			TestCase:     segments[1],
			SequenceID:   segments[3],
		}
	}

	if len(segments) == 3 {
		return TestRunResources{
			Organisation: segments[0],
			TestCase:     segments[1],
			SequenceID:   segments[2],
		}
	}

	if len(segments) == 1 {
		return TestRunResources{
			UID: segments[0],
		}
	}

	return TestRunResources{}
}
