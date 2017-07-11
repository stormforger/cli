package api

import (
	"io"
	"io/ioutil"
)

// TestCaseValidate will send a test case definition (JS) to the API
// to validate.
func (c *Client) TestCaseValidate(fileName string, data io.Reader) (bool, string, error) {
	// TODO how to pass options here?
	//      defining a struct maybe, but where?
	//      finally: add options here
	extraParams := map[string]string{}

	req, err := newfileUploadRequest(c.APIEndpoint+"/test_cases/validate", extraParams, "test_case", fileName, data)
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
