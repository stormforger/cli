package api

const (
	UserTypeServiceAccount = "service-account"
	UserTypeUser           = "user"
)

// User describes a API user
type User struct {
	Mail            string `json:"user_email"`
	AuthenticatedAs *struct {
		UID   string `json:"uid"`
		Label string `json:"label"`
		Type  string `json:"type"` // "service-account" or "user"
	} `json:"authenticated_as"`
}
