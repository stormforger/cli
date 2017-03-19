package api

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type testRunResources struct {
	uid          string
	organisation string
	testCase     string
	sequenceID   string
}

// TestRunCallLogPreview will download the first 10k lines
// of the test run's call log
func (c *Client) TestRunCallLogPreview(pathID string) (string, error) {
	testRun := extractResources(pathID)

	path := ""
	if testRun.uid == "" {
		path = "/test_cases/" + testRun.organisation + "/" + testRun.testCase + "/test_runs/" + testRun.sequenceID
	} else {
		path = "/t/" + testRun.uid
	}

	// pathID looks like foo/demo/test_runs/19
	req, err := http.NewRequest("GET", c.APIEndpoint+path+"/call_log_preview", nil)
	if err != nil {
		return "", err
	}

	// TODO how to set these on all requests?
	req.Header.Add("Authorization", "Bearer "+c.JWT)
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New("could not download call log")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
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
