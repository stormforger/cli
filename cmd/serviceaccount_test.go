package cmd

import (
	"testing"

	"github.com/stormforger/cli/api"
	"github.com/stretchr/testify/assert"
)

func TestMainServiceAccountsList(t *testing.T) {
	// Given
	const (
		exampleOrg = "example-inc"
	)

	handler := givenStaticContentHandler(`{
		"data": [
			{
				"type": "service_account",
				"id": "123456",
				"attributes": {
					"token_label": "Hello World"
				}
			}
		]
	}`)
	server := givenHTTPServer(t, handler)

	// When
	client := api.NewClient(server.URL, "<somejwttoken>")
	list, err := MainServiceAccountsList(client, exampleOrg)

	// Then
	assert.NoError(t, err)
	assert.Len(t, list.ServiceAccounts, 1)
	assert.Equal(t, "123456", list.ServiceAccounts[0].UID)
	assert.Equal(t, "Hello World", list.ServiceAccounts[0].TokenLabel)
}

func TestMainServiceAccountsCreate__ReceivedAccessToken(t *testing.T) {
	// TODO: maybe test that this was a POST and we sent the necessary payload?

	// Given
	const (
		exampleOrg = "example-inc"
		tokenLabel = "example-ci-token"
	)

	handler := givenStaticContentHandler(`{
		"data": {
			"type": "service_account",
			"id": "123456",
			"attributes": {
				"token_label": "` + tokenLabel + `",
				"access_token": "SOme.Secret.JWTValue"
			}
		}
	}`)
	server := givenHTTPServer(t, handler)

	// When
	client := api.NewClient(server.URL, "<somejwttoken>")
	sa, err := MainServiceAccountsCreate(client, exampleOrg, tokenLabel)

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, sa)
	assert.Equal(t, "123456", sa.UID)
	assert.NotEmpty(t, sa.AccessToken, "Expected to receive access_token")
}
