package types

import "github.com/google/uuid"

type ContextKey string

var UserContextKey ContextKey = "user"

type AuthenticatedUser struct {
	Email     string
	LoggedIn  bool
	FirstName string
	LastName  string
	AvatarUrl string
	ID        uuid.UUID
}
