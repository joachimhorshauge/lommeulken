package view

import (
	"context"
	"lommeulken/types"
)

func AuthenticatedUser(ctx context.Context) types.AuthenticatedUser {
	user := ctx.Value(types.UserContextKey)
	if u, ok := user.(types.AuthenticatedUser); ok {
		return u
	}
	return types.AuthenticatedUser{
		Email:    "",
		LoggedIn: false,
	}
}
