package types

type userKey string

const (
	Userkey userKey = "user"
	// ...
)

type AuthenticatedUser struct {
	Email    string
	LoggedIn bool
}
