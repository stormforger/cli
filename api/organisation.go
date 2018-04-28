package api

// ListOrganisations returns a list of organisations
func (c *Client) ListOrganisations() ([]byte, error) {
	path := "/organisations"

	_, response, err := c.fetch(path)

	return response, err
}
