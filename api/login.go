package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Login acquires a JWT access token for the given email/password
func (c *Client) Login(email string, password string) (string, error) {
	data := map[string]string{"email": email, "password": password}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.APIEndpoint+"/user/token", &body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	c.setUserAgent(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	defer close(resp.Body)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New("login not successful! please verify that you can login with these credentials at https://app.stormforger.com")
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(responseBody, &dat); err != nil {
		return "", err
	}

	return dat["jwt"].(string), nil
}
