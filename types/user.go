package types

import "github.com/google/uuid"

type userKey string

const (
	Userkey userKey = "user"
	// ...
)

type AuthenticatedUser struct {
	ID          uuid.UUID
	Email       string
	LoggedIn    bool
	AccessToken string
	Account
}
