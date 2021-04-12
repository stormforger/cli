package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// ListServiceAccounts returns a list of organisations
func (c *Client) ListServiceAccounts(org string) (bool, []byte, error) {
	path := fmt.Sprintf("/organisations/%s/service_accounts/", org)

	return c.fetch(path)
}

func (c *Client) CreateServiceAccount(org, token_label string) (bool, string, error) {
	path := fmt.Sprintf("/organisations/%s/service_accounts/", org)
	params := url.Values{}
	params.Set("user[token_label]", token_label)

	req, err := http.NewRequest(http.MethodPost, c.APIEndpoint+path, strings.NewReader(params.Encode()))
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
