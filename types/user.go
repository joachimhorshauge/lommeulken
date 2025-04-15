package types

type ContextKey string

var UserContextKey ContextKey = "user"

type AuthenticatedUser struct {
	Email    string
	LoggedIn bool
}
