package api

import (
	"net/http"
)

// ListOrganisations returns a list of organisations
func (c *Client) ListOrganisations() ([]byte, error) {
	path := "/organisations"

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
