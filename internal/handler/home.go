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
	}
	navMenus := []components.NavMenu{
		{Title: "Kalender", Links: []components.NavLink{
			{Title: "Arrangementer", Url: "calendar"},
			{Title: "Book udstyr", Url: "calendar/booking"}}},
		{Title: "Fangster", Links: []components.NavLink{
			{Title: "Seneste", Url: "catch"},
			{Title: "Klubrekorder", Url: "club-records"},
			{Title: "Opret fangst", Url: "new-catch"}}},
		{Title: "Konkurrencen", Links: []components.NavLink{
			{Title: "Igangværende", Url: "competition/current"},
			{Title: "Tidligere", Url: "competition/previous"},
			{Title: "Regler", Url: "competition/rules"}}},
		{Title: "Hvem er vi?", Links: []components.NavLink{
			{Title: "Om os", Url: "about"},
			{Title: "Vandpleje", Url: "about/water-care"}}},
	}
	err := home.Index(navLinks, navMenus).Render(r.Context(), w)
	if err != nil {
		slog.Error("Error rendering home page", "error", err)
		return
	}
}
