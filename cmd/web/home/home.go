package home

import (
	"log"
	"net/http"
)

func HomeWebHandler(w http.ResponseWriter, r *http.Request) {
	component := Index()
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in HelloWebHandler: %e", err)
	}
}
