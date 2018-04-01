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
)

type testRunResources struct {
	uid          string
	organisation string
	testCase     string
	sequenceID   string
}

// TestRunList is a list of TestRuns, used for index action
type TestRunList struct {
	TestRuns []*TestRun
}

// TestRun represents a single TestRun
type TestRun struct {
	ID           string `jsonapi:"primary,test_runs"`
	Title        string `jsonapi:"attr,title,omitempty"`
	Notes        string `jsonapi:"attr,notes,omitempty"`
	State        string `jsonapi:"attr,state,omitempty"`
	StartedBy    string `jsonapi:"attr,started_by,omitempty"`
	StartedAt    string `jsonapi:"attr,started_at,omitempty"`
	EndedAt      string `jsonapi:"attr,ended_at,omitempty"`
	EstimatedEnd string `jsonapi:"attr,estimated_end,omitempty"`
}

// TestRunList will list all test runs for a given test case
func (c *Client) TestRunList(testCaseUID string) (bool, []byte, error) {
	path := "/test_cases/" + testCaseUID + "/test_runs"

	status, response, err := c.fetch(path)

	return status, response, err
}

// TestRunShow will show some basic information on a given
// test run
func (c *Client) TestRunShow(uid string) (bool, []byte, error) {
	path := "/test_runs/" + uid

	status, response, err := c.fetch(path)

	return status, response, err
}

// TestRunWatch will show some basic information on a given
// test run
func (c *Client) TestRunWatch(uid string) (TestRun, string, error) {
	path := "/test_runs/" + uid + "/watch"

	req, err := http.NewRequest("GET", c.APIEndpoint+path, nil)
	if err != nil {
		return TestRun{}, "", err
	}

	response, err := c.doRequestRaw(req)
	if err != nil {
		return TestRun{}, "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return TestRun{}, "", err
	}

	if response.StatusCode >= 400 {
		log.Fatal(string(body))
	}

	testRun := new(TestRun)
	err = jsonapi.UnmarshalPayload(bytes.NewReader(body), testRun)
	if err != nil {
		log.Fatal(err)
	}

	return *testRun, string(body), nil
}

// TestRunCallLog will download the first 10k lines
// of the test run's call log
func (c *Client) TestRunCallLog(pathID string, preview bool) (io.ReadCloser, error) {
	testRun := extractResources(pathID)

	path := "/test_runs/" + testRun.uid + "/call_log"

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

// TestRunCreate will send a test case definition (JS) to the API
// to update an existing test case it.
func (c *Client) TestRunCreate(testCaseUID string, title string, notes string) (bool, string, error) {
	payload := bytes.NewBuffer(nil)
	newTestRun := &TestRun{
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
