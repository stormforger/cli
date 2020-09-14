package cmd

import (
	"strings"
	"testing"

	"github.com/stormforger/cli/api"
	"github.com/stretchr/testify/assert"
)

func TestPingCommandOutput_Unauthenticated(t *testing.T) {
	var w strings.Builder

	// For unauthenticated pings, we have no subject or anything else.
	result := PingCommandResult{
		Success:         true,
		Unauthenticated: true,
	}
	PrintHumanPingResult(&w, result)

	assert.Equal(t, "PONG!\n", w.String())
}

func TestPingCommandOutput_ServiceAccount(t *testing.T) {
	var w strings.Builder

	// For unauthenticated pings, we have no subject or anything else.
	user := api.User{
		Mail: "test-UID@noreplay.stormforger.com",
	}
	user.AuthenticatedAs = &struct {
		UID   string "json:\"uid\""
		Label string "json:\"label\""
		Type  string "json:\"type\""
	}{
		UID:   "myUID",
		Label: "testname",
		Type:  api.UserTypeServiceAccount,
	}

	result := PingCommandResult{
		Success:         true,
		Unauthenticated: false,
		Subject:         &user,
	}
	PrintHumanPingResult(&w, result)

	assert.Equal(t, "PONG! Authenticated as testname (myUID)\n", w.String())
}

func TestPingCommandOutput_User(t *testing.T) {
	var w strings.Builder

	// For unauthenticated pings, we have no subject or anything else.
	user := api.User{
		Mail: "mr.smith@loadtest.party",
	}
	user.AuthenticatedAs = &struct {
		UID   string "json:\"uid\""
		Label string "json:\"label\""
		Type  string "json:\"type\""
	}{
		UID:   "myUID",
		Label: "smith",
		Type:  api.UserTypeUser,
	}

	result := PingCommandResult{
		Success:         true,
		Unauthenticated: false,
		Subject:         &user,
	}
	PrintHumanPingResult(&w, result)

	assert.Equal(t, "PONG! Authenticated as mr.smith@loadtest.party\n", w.String())
}
