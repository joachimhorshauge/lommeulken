package auth

import (
	"log"
	"net/http"
)

func LoginWebHandler(w http.ResponseWriter, r *http.Request) {
	component := Login()
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in LoginWebHandler: %e", err)
	}
}

func SignupWebHandler(w http.ResponseWriter, r *http.Request) {
	component := SignUp()
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in SignupWebHandler: %e", err)
	}
}
