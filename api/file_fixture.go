package api

import (
	"io/ioutil"
	"log"
	"net/http"
)

// FileFixtureParams represents params BLA TODO
type FileFixtureParams struct {
	Name       string
	Type       string
	FieldNames string
	Delimiter  string
}

// ListFileFixture returns a list of the organizations fixtures
func (c *Client) ListFileFixture(organization string) ([]byte, error) {
	path := "/file_fixtures/" + organization

	log.Println(c.APIEndpoint + path)

	req, err := http.NewRequest("GET", c.APIEndpoint+path, nil)
	if err != nil {
		return nil, err
	}

	// TODO how to set these on all requests?
	c.addDefaultHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	return body, nil
}

// PushFileFixture uploads (insert or update) a file fixture
func (c *Client) PushFileFixture(file string, organization string, params *FileFixtureParams) (string, error) {
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

	req, err := newfileUploadRequest(c.APIEndpoint+"/file_fixtures/"+organization, extraParams, "file_fixture[file_fixture_version][original]", file)
	if err != nil {
		return "", err
	}

	// TODO how to set these on all requests?
	c.addDefaultHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	return string(body), nil
}

// DeleteFileFixture deletes a file fixture
func (c *Client) DeleteFileFixture(fileFixtureUID string, organization string) (string, error) {
	req, err := http.NewRequest("DELETE", c.APIEndpoint+"/file_fixtures/"+fileFixtureUID+"?organisation_uid="+organization, nil)
	if err != nil {
		return "", err
	}

	// TODO how to set these on all requests?
	c.addDefaultHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return string(body), nil
}
