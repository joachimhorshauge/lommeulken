package handler

import (
	"github.com/joachimhorshauge/lommeulken/cmd/web/components"
	"log/slog"
	"net/http"

	"github.com/joachimhorshauge/lommeulken/cmd/web/templates/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) {
	navLinks := []components.NavLink{
		{Title: "Bliv Medlem", Url: "signup"},
		{Title: "Nyheder", Url: "news"},
		{Title: "Kalender", Url: "calendar"},
		{Title: "Fangster", Url: "catches"},
		{Title: "Konkurrencen", Url: "competition"},
		{Title: "Hvem er vi?", Url: "about"},
	}
	err := home.Index(navLinks).Render(r.Context(), w)
	if err != nil {
		slog.Error("Error rendering home page", "error", err)
		return
	}
}
