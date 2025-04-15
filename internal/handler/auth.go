package handler

import (
	"context"
	"log/slog"
	"lommeulken/cmd/web"
	"lommeulken/util"
	"net/http"
	"strings"

	sb "lommeulken/internal/supabase"

	"github.com/nedpals/supabase-go"
)

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := web.SignupIndex().Render(r.Context(), w)
		if err != nil {
			slog.Error("Error rendering Signup page", "error", err)
			return
		}
		return
	}

	if r.Method == http.MethodPost {
		credentials := web.UserSignupCredentials{
			Email:           r.PostFormValue("email"),
			Username:        r.PostFormValue("username"),
			Password:        r.PostFormValue("password"),
			ConfirmPassword: r.PostFormValue("confirmPassword"),
		}

		errors := web.SignupErrors{
			EmailInvalid:     util.ValidateEmail(credentials.Email),
			UsernameInvalid:  util.ValidateUsername(credentials.Username),
			PasswordInvalid:  util.ValidatePassword(credentials.Password),
			PasswordMismatch: util.ValidatePasswordMatch(credentials.Password, credentials.ConfirmPassword),
		}

		if errors.EmailInvalid != "" || errors.UsernameInvalid != "" || errors.PasswordInvalid != "" || errors.PasswordMismatch != "" {
			err := web.SignupForm(credentials, errors).Render(r.Context(), w)
			if err != nil {
				slog.Error("Error rendering Signup form with errors", "error", err)
			}
			return
		}

		ctx := context.Background()

		_, err := sb.Client.Auth.SignUp(ctx, supabase.UserCredentials{
			Email:    credentials.Email,
			Password: credentials.Password,
		})

		if err != nil {
			slog.Error("Error creating credentials in Supabase", "error", err)
			http.Error(w, "Failed to create credentials", http.StatusInternalServerError)
			return
		}

		// Redirect or render success page
		w.Header().Add("Hx-Redirect", "/login")
	}
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := web.LoginIndex().Render(r.Context(), w)
		if err != nil {
			slog.Error("Error rendering login page", "error", err)
			return
		}
	} else if r.Method == "POST" {
		credentials := web.UserCredentials{
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}

		ctx := context.Background()
		details, err := sb.Client.Auth.SignIn(ctx, supabase.UserCredentials{
			Email:    credentials.Email,
			Password: credentials.Password,
		})

		if err != nil {
			slog.Error("Failed to login user", "email", credentials.Email, "error", err)
			errors := web.LoginErrors{}

			if strings.Contains(err.Error(), "Email not confirmed") {
				errors.EmailNotVerified = "Du skal verificere din email før du kan logge ind. Tjek din indbakke for en bekræftelsesemail."
			} else {
				errors.InvalidCredentials = "Forkert brugernavn eller kodeord"
			}

			err = web.LoginForm(credentials, errors).Render(r.Context(), w)
			if err != nil {
				slog.Error("Error rendering login page", "error", err)
				return
			}
			return
		}
		cookie := &http.Cookie{
			Value:    details.AccessToken,
			Name:     "at",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		}

		http.SetCookie(w, cookie)

		w.Header().Add("Hx-Redirect", "/")
	}
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Value:    "",
		Name:     "at",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
