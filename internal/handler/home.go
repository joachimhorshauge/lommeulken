package handler

import (
	"context"
	"net/http"

	"github.com/joachimhorshauge/lommeulken/cmd/web/templates/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) {
    home.Index().Render(context.Background(), w)
}
