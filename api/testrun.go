package api

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/jsonapi"
	"github.com/stormforger/cli/api/testrun"
)

// TestRunLaunchOptions represents a single TestRunLaunchOptions
type TestRunLaunchOptions struct {
	Title                string
	Notes                string
	JavascriptDefinition string

	ClusterRegion         string
	ClusterSizing         string
	DisableGzip           bool
	SkipWait              bool
	DumpTraffic           bool
	SessionValidationMode bool
}

// TestRunResources describes infos on a test run
type TestRunResources struct {
	UID          string
	Organisation string
	TestCase     string
	SequenceID   string
}

// TestRunList will list all test runs for a given test case
func (c *Client) TestRunList(testCaseUID string, filter string) (bool, []byte, error) {
	path := "/test_cases/" + testCaseUID + "/test_runs"

	switch filter {
	case "archived":
		path = path + "/?only=archived"
	case "all":
		path = path + "/?only=all"
	}

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
	defer response.Body.Close()

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

	return c.TestRunLogs(path)
}

// TestRunUserLog will download the user logs
func (c *Client) TestRunUserLog(pathID string) (io.ReadCloser, error) {
	testRun := ExtractTestRunResources(pathID)

	path := "/test_runs/" + testRun.UID + "/user_log"

	return c.TestRunLogs(path)
}

// TestRunLogs will download logs from the given path. It will also
// handle compression accordingly.
func (c *Client) TestRunLogs(path string) (io.ReadCloser, error) {
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
		response.Body.Close()
		if response.StatusCode == 404 {
			return nil, errors.New("test-run or logs not found")
		}
		return nil, errors.New("could not download log")
	}

	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
	default:
		reader = response.Body
	}
	return reader, nil
}

// TestRunDump will fetch a traffic dump for a
// given test run if it is available.
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
		response.Body.Close()
		return nil, err
	}

	if response.StatusCode != 200 {
		response.Body.Close()
		return nil, errors.New("could not load full dump")
	}

	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
	default:
		reader = response.Body
	}

	return reader, nil
}

// TestRunCreate will send a test case definition (JS) to the API
// to update an existing test case it.
func (c *Client) TestRunCreate(testCaseUID string, options TestRunLaunchOptions) (bool, string, error) {
	// Only upload test_configuration_attributes if one of option is non-zero for this
	method := http.MethodPost
	uri := c.APIEndpoint + "/test_cases/" + testCaseUID + "/test_runs"
	payload := url.Values{}

	boolFields := []struct {
		Field     string
		BoolValue bool
	}{
		{"data[test_configuration_attributes][disable_gzip]", options.DisableGzip},
		{"data[test_configuration_attributes][skip_wait]", options.SkipWait},
		{"data[test_configuration_attributes][dump_traffic_full]", options.DumpTraffic},
		{"data[test_configuration_attributes][session_validation_mode]", options.SessionValidationMode},
	}
	for _, f := range boolFields {
		if f.BoolValue {
			payload.Add(f.Field, "true")
		}
	}

	stringFields := []struct {
		Field string
		Value string
	}{
		{"data[test_configuration_attributes][cluster_region]", options.ClusterRegion},
		{"data[test_configuration_attributes][cluster_sizing]", options.ClusterSizing},
	}
	for _, f := range stringFields {
		if f.Value != "" {
			payload.Add(f.Field, f.Value)
		}
	}

	payload.Add("data[attributes][title]", options.Title)
	payload.Add("data[attributes][notes]", options.Notes)

	fmt.Println(payload.Encode())

	// build a multipart request, if we have a javascript_definition to upload
	var req *http.Request
	var err error
	if options.JavascriptDefinition != "" {
		req, err = fileUploadRequest(uri, method, payload, "data[attributes][javascript_definition]", "TODOtest-case.js", "application/javascript", strings.NewReader(options.JavascriptDefinition))
		if err != nil {
			return false, "", err
		}
	} else {
		req, err = http.NewRequest(method, uri, strings.NewReader(payload.Encode()))
		if err != nil {
			return false, "", err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	response, err := c.doRequestRaw(req)
	if err != nil {
		return false, "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, "", err
	}

	defer close(response.Body)

	return response.StatusCode < 400, string(body), nil
}

// TestRunAbort will send a test case definition (JS) to the API
// to update an existing test case it.
func (c *Client) TestRunAbort(testRunUID string) (bool, string, error) {
	req, err := http.NewRequest("POST", c.APIEndpoint+"/test_runs/"+testRunUID+"/abort", nil)
	if err != nil {
		return false, "", err
	}

	response, err := c.doRequestRaw(req)
	if err != nil {
		return false, "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, "", err
	}

	defer close(response.Body)

	return response.StatusCode < 400, string(body), nil
}

// TestRunNfrCheck will upload requirements definition
// and checks if the given test run matches them.
func (c *Client) TestRunNfrCheck(uid string, fileName string, data io.Reader) (bool, []byte, error) {
	extraParams := url.Values{}

	path := "/test_runs/" + uid + "/check_nfr"

	req, err := fileUploadRequest(c.APIEndpoint+path, "POST", extraParams, "nfr_file", fileName, "application/x-yaml", data)
	if err != nil {
		log.Fatal(err)
	}

	response, err := c.doRequestRaw(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, nil, err
	}

	defer close(response.Body)

	return response.StatusCode < 400, body, nil
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
