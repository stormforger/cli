package api

// ListOrganisations returns a list of organisations
func (c *Client) ListOrganisations() (bool, []byte, error) {
	path := "/organisations"

	return c.fetch(path)
}
