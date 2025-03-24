package handler

import (
	"github.com/joachimhorshauge/lommeulken/cmd/web/components"
	"log/slog"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := components.LoginIndex().Render(r.Context(), w)
		if err != nil {
			slog.Error("Error rendering login page", "error", err)
			return
		}
	} else if r.Method == "POST" {
		user := components.UserCredentials{
			Username: r.PostFormValue("username"),
			Password: r.PostFormValue("password"),
		}
		errors := components.LoginErrors{
			InvalidCredentials: "Forkert brugernavn eller kodeord",
		}
		err := components.LoginForm(user, errors).Render(r.Context(), w)
		if err != nil {
			slog.Error("Error rendering login page", "error", err)
			return
		}
	}

}
