package middleware

import (
	"context"
	"log/slog"
	"lommeulken/gen/dbstore"
	"lommeulken/internal/supabase"
	"lommeulken/types"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (m *Middleware) WithUser(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		accessToken, err := r.Cookie("at")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		resp, err := supabase.Client.Auth.User(r.Context(), accessToken.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		id, err := uuid.Parse(resp.ID)
		if err != nil {
			slog.Warn("Failed to parse uuid in middleware")
		}

		profile, err := m.queries.GetUserByID(context.Background(), id)
		if err != nil {
			slog.Warn("Failed to get user with id", "ID", resp.ID)
			_, err = m.queries.CreateUser(context.Background(), dbstore.CreateUserParams{ID: id, Email: pgtype.Text{String: resp.Email, Valid: true}})
			if err != nil {
				slog.Error("Failed to create user  in db with id", "ID", resp.ID)
			}
			slog.Info("Failed to create user  in db with id", "ID", resp.ID)

		}

		user := types.AuthenticatedUser{
			ID:        id,
			Email:     resp.Email,
			LoggedIn:  true,
			FirstName: profile.FirstName.String,
			LastName:  profile.LastName.String,
			AvatarUrl: profile.AvatarUrl.String,
		}

		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(handler)
}
