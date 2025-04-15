package view

import (
	"context"
	"log/slog"
	"lommeulken/types"
)

func AuthenticatedUser(ctx context.Context) types.AuthenticatedUser {
	user := ctx.Value(types.UserContextKey)
	slog.Info("user", "user", user)
	if u, ok := user.(types.AuthenticatedUser); ok {
		return u
	}
	return types.AuthenticatedUser{
		Email:    "",
		LoggedIn: false,
	}
}
