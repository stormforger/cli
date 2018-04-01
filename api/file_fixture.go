package api

import (
	"io"
	"net/http"
)

// FileFixtureParams represents params BLA TODO
type FileFixtureParams struct {
	Name       string
	Type       string
	FieldNames string
	Delimiter  string
}

// MoveFileFixture renames a filefixtures
func (c *Client) MoveFileFixture(organization string, fileFixtureUID string, newName string) (string, error) {
	params := map[string]string{"file_fixture[name]": newName}

	req, err := newPatchRequest(c.APIEndpoint+"/file_fixtures/"+organization+"/"+fileFixtureUID, params)
	if err != nil {
		return "", err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// ListFileFixture returns a list of the organizations fixtures
func (c *Client) ListFileFixture(organization string) ([]byte, error) {
	path := "/file_fixtures/" + organization + "?only=structured"

	_, response, err := c.fetch(path)

	return response, err
}

// PushFileFixture uploads (insert or update) a file fixture
func (c *Client) PushFileFixture(fileName string, data io.Reader, organization string, params *FileFixtureParams) (string, error) {
	extraParams := map[string]string{
		"file_fixture[name]": params.Name,
		"file_fixture[type]": params.Type,
	}

	if params.Delimiter != "" {
		extraParams["file_fixture[file_fixture_version][delimiter]"] = params.Delimiter
	}

	if params.FieldNames != "" {
		extraParams["file_fixture[file_fixture_version][field_names]"] = params.FieldNames
	}

	req, err := fileUploadRequest(c.APIEndpoint+"/file_fixtures/"+organization, "POST", extraParams, "file_fixture[file_fixture_version][original]", fileName, "application/octet-stream", data)
	if err != nil {
		return "", err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// DeleteFileFixture deletes a file fixture
func (c *Client) DeleteFileFixture(fileFixtureUID string, organization string) (string, error) {
	req, err := http.NewRequest("DELETE", c.APIEndpoint+"/file_fixtures/"+organization+"/"+fileFixtureUID, nil)
	if err != nil {
		return "", err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// DownloadFileFixture retrieves the originally uploaded file
func (c *Client) DownloadFileFixture(organization string, fileFixtureUID string, version string) ([]byte, error) {
	path := "/file_fixtures/" + organization + "/" + fileFixtureUID + "/download/" + version

	_, response, err := c.fetch(path)

	return response, err
}
