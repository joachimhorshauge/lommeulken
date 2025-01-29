package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const currentPathKey contextKey = "currentPath"

func CurrentPathMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), currentPathKey, r.URL.Path)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCurrentPath(ctx context.Context) string {
	if path, ok := ctx.Value(currentPathKey).(string); ok {
		return path
	}
	return "/"
}
