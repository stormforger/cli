package api

import (
	"io"
	"io/ioutil"
	"net/url"
)

// ListTestCases returns a list of test cases
func (c *Client) ListTestCases(organization string, filter string) (bool, []byte, error) {
	path := "/v2/test_cases"
	query := url.Values{}
	query.Set("scope", organization)

	switch filter {
	case "archived":
		query.Set("only", "archived")
	case "all":
		query.Set("only", "all")
	}

	return c.fetchWithQuery(path, query)
}

// TestCaseValidate will send a test case definition (JS) to the API
// to validate.
func (c *Client) TestCaseValidate(organization string, fileName string, data io.Reader) (bool, string, error) {
	// TODO how to pass options here?
	//      defining a struct maybe, but where?
	//      finally: add options here
	extraParams := url.Values{}
	extraParams.Add("scope", organization)

	req, err := fileUploadRequest(c.APIEndpoint+"/v2/test_cases/validate", "POST", extraParams, "test_case[javascript_definition]", fileName, "application/javascript", data)
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

	if response.StatusCode != 200 {
		return false, string(body), nil
	}

	return true, string(body), nil
}

// TestCaseCreate will send a test case definition (JS) to the API
// to create it.
func (c *Client) TestCaseCreate(organization string, testCaseName string, fileName string, data io.Reader) (bool, string, error) {
	// TODO how to pass options here?
	//      defining a struct maybe, but where?
	//      finally: add options here
	extraParams := url.Values{}
	extraParams.Add("test_case[name]", testCaseName)
	extraParams.Set("scope", organization)

	req, err := fileUploadRequest(c.APIEndpoint+"/v2/test_cases", "POST", extraParams, "test_case[javascript_definition]", fileName, "application/javascript", data)
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

	if response.StatusCode != 200 {
		return false, string(body), nil
	}

	return true, string(body), nil
}

// TestCaseUpdate will send a test case definition (JS) to the API
// to update an existing test case it.
func (c *Client) TestCaseUpdate(testCaseUID string, fileName string, data io.Reader) (bool, string, error) {
	// TODO how to pass options here?
	//      defining a struct maybe, but where?
	//      finally: add options here
	extraParams := url.Values{}

	req, err := fileUploadRequest(c.APIEndpoint+"/v2/test_cases/"+testCaseUID, "PATCH", extraParams, "test_case[javascript_definition]", fileName, "application/javascript", data)
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

	if response.StatusCode != 200 {
		return false, string(body), nil
	}

	return true, string(body), nil
}

// TestCaseArchive will mark a test case as archived
func (c *Client) TestCaseArchive(uid string) (bool, []byte, error) {
	path := "/v2/test_cases/" + uid + "/archive"

	return c.put(path, nil)

}

// TestCaseUnArchive will mark a test case as not archived
func (c *Client) TestCaseUnArchive(uid string) (bool, []byte, error) {
	path := "/v2/test_cases/" + uid + "/unarchive"

	return c.put(path, nil)

}

// DownloadTestCaseDefinition returns the JS definition
// of a given test case
func (c *Client) DownloadTestCaseDefinition(uid string) (bool, []byte, error) {
	path := "/test_cases/" + uid + "/download"

	return c.fetch(path)
}
