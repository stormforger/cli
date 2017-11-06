package api

import (
	"io"
	"io/ioutil"
	"net/http"
)

// ListTestCases returns a list of test cases
func (c *Client) ListTestCases(organization string) ([]byte, error) {
	path := "/test_cases?organisation_uid=" + organization

	req, err := http.NewRequest("GET", c.APIEndpoint+path, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// CheckExistanceTestCase returns a list of test cases
func (c *Client) CheckExistanceTestCase(organization string, testCaseName string) ([]byte, error) {
	path := "/test_cases/" + testCaseName + "?organisation_uid=" + organization

	req, err := http.NewRequest("GET", c.APIEndpoint+path, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	return body, nil
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

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, "", err
	}

	defer response.Body.Close()

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
		"test_case[name]":  testCaseName,
		"organisation_uid": organization,
	}

	req, err := fileUploadRequest(c.APIEndpoint+"/test_cases", "POST", extraParams, "test_case[javascript_definition]", fileName, "application/javascript", data)
	if err != nil {
		return false, "", err
	}

	response, err := c.doRequestRaw(req)
	if err != nil {
		return false, "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, "", err
	}

	defer response.Body.Close()

	return false, string(body), nil
}

// TestCaseUpdate will send a test case definition (JS) to the API
// to update an existing test case it.
func (c *Client) TestCaseUpdate(organization string, testCaseUid string, fileName string, data io.Reader) (bool, string, error) {
	// TODO how to pass options here?
	//      defining a struct maybe, but where?
	//      finally: add options here
	extraParams := map[string]string{
		"organisation_uid": organization,
	}

	req, err := fileUploadRequest(c.APIEndpoint+"/test_cases/"+testCaseUid, "PATCH", extraParams, "test_case[javascript_definition]", fileName, "application/javascript", data)
	if err != nil {
		return false, "", err
	}

	response, err := c.doRequestRaw(req)
	if err != nil {
		return false, "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, "", err
	}

	defer response.Body.Close()

	return false, string(body), nil
}
