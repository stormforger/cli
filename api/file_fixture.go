package api

import (
	"io"
	"io/ioutil"
	"net/http"
)

// FileFixtureParams represents params BLA TODO
type FileFixtureParams struct {
	Name            string
	Raw             bool
	FieldNames      string
	Delimiter       string
	FirstRowHeaders bool
}

// MoveFileFixture renames a filefixtures
func (c *Client) MoveFileFixture(organization string, fileFixtureUID string, newName string) (bool, string, error) {
	params := map[string]string{"file_fixture[name]": newName}

	req, err := newPatchRequest(c.APIEndpoint+"/file_fixtures/"+organization+"/"+fileFixtureUID, params)
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

// ListFileFixture returns a list of the organizations fixtures
func (c *Client) ListFileFixture(organization string) (bool, []byte, error) {
	path := "/file_fixtures/" + organization + "?only=structured"

	success, response, err := c.fetch(path)

	return success, response, err
}

// PushFileFixture uploads (insert or update) a file fixture
func (c *Client) PushFileFixture(fileName string, data io.Reader, organization string, params *FileFixtureParams) (bool, []byte, error) {
	fixtureType := "structured"
	if params.Raw {
		fixtureType = "raw"
	}
	extraParams := map[string]string{
		"file_fixture[name]": params.Name,
		"file_fixture[type]": fixtureType,
	}

	if params.FirstRowHeaders {
		extraParams["file_fixture[file_fixture_version][first_row_headers]"] = "1"
	}

	if params.Delimiter != "" {
		extraParams["file_fixture[file_fixture_version][delimiter]"] = params.Delimiter
	}

	if params.FieldNames != "" {
		extraParams["file_fixture[file_fixture_version][field_names]"] = params.FieldNames
	}

	req, err := fileUploadRequest(c.APIEndpoint+"/file_fixtures/"+organization, "POST", extraParams, "file_fixture[file_fixture_version][original]", fileName, "application/octet-stream", data)
	if err != nil {
		return false, nil, err
	}

	response, err := c.doRequestRaw(req)
	if err != nil {
		return false, nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, nil, err
	}

	defer close(response.Body)

	if response.StatusCode >= 300 {
		return false, body, nil
	}

	return true, body, nil
}

// DeleteFileFixture deletes a file fixture
func (c *Client) DeleteFileFixture(fileFixtureUID string, organization string) (bool, string, error) {
	req, err := http.NewRequest("DELETE", c.APIEndpoint+"/file_fixtures/"+organization+"/"+fileFixtureUID, nil)
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

	if response.StatusCode >= 300 {
		return false, string(body), nil
	}

	return true, string(body), nil
}

// DownloadFileFixture retrieves the originally uploaded file
func (c *Client) DownloadFileFixture(organization string, fileFixtureUID string, version string) (bool, []byte, error) {
	path := "/file_fixtures/" + organization + "/" + fileFixtureUID + "/download/" + version

	success, response, err := c.fetch(path)

	return success, response, err
}
