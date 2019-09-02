package api

import (
	"io"
	"io/ioutil"
)

// ListTestCases returns a list of test cases
func (c *Client) ListTestCases(organization string, filter string) (bool, []byte, error) {
	path := "/organisations/" + organization + "/test_cases"

	switch filter {
	case "archived":
		path = path + "/?only=archived"
	case "all":
		path = path + "/?only=all"
	}

	return c.fetch(path)
}

// TestCaseValidate will send a test case definition (JS) to the API
// to validate.
func (c *Client) TestCaseValidate(organization string, fileName string, data io.Reader) (bool, string, error) {
	// TODO how to pass options here?
	//      defining a struct maybe, but where?
	//      finally: add options here
	extraParams := map[string]string{
		"organisation_uid": organization,
	}

	req, err := fileUploadRequest(c.APIEndpoint+"/test_cases/validate", "POST", extraParams, "test_case[javascript_definition]", fileName, "application/javascript", data)
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
	extraParams := map[string]string{
		"test_case[name]": testCaseName,
	}

	req, err := fileUploadRequest(c.APIEndpoint+"/organisations/"+organization+"/test_cases", "POST", extraParams, "test_case[javascript_definition]", fileName, "application/javascript", data)
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
	extraParams := map[string]string{}

	req, err := fileUploadRequest(c.APIEndpoint+"/test_cases/"+testCaseUID, "PATCH", extraParams, "test_case[javascript_definition]", fileName, "application/javascript", data)
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

// DownloadTestCaseDefinition returns the JS definition
// of a given test case
func (c *Client) DownloadTestCaseDefinition(uid string) (bool, []byte, error) {
	path := "/test_cases/" + uid + "/download"

	return c.fetch(path)
}
