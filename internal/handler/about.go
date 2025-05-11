package handler

import (
	"log/slog"
	"lommeulken/cmd/web"
	"net/http"
)

func (h *Handler) AboutIndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := web.AboutUsIndex().Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render about us index with errors", "error", err)
	}
	return
}
