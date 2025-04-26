package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"lommeulken/cmd/web"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type fishSpecies struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type ImageData struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	IsPrimary bool   `json:"is_primary"`
	CreatedAt string `json:"created_at"`
}

func (h *Handler) HandleListSpecies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	species := []fishSpecies{
		{Value: "Aborre", Label: "Aborre"},
		{Value: "Brasen", Label: "Brasen"},
		{Value: "Bækørred", Label: "Bækørred"},
		{Value: "Fladfisk", Label: "Fladfisk"},
		{Value: "Gedde", Label: "Gedde"},
		{Value: "Græskarpe", Label: "Græskarpe"},
		{Value: "Havbars", Label: "Havbars"},
		{Value: "Havørred Kysten", Label: "Havørred Kysten"},
		{Value: "Havørred Åen", Label: "Havørred Åen"},
		{Value: "Helt", Label: "Helt"},
		{Value: "Hornfisk", Label: "Hornfisk"},
		{Value: "Laks", Label: "Laks"},
		{Value: "Makrel", Label: "Makrel"},
		{Value: "Multe", Label: "Multe"},
		{Value: "Pighvar/Slethvar", Label: "Pighvar/Slethvar"},
		{Value: "Put & Take ørred", Label: "Put & Take ørred"},
		{Value: "Rimte", Label: "Rimte"},
		{Value: "Sandart", Label: "Sandart"},
		{Value: "Skalle", Label: "Skalle"},
		{Value: "Skælkarpe", Label: "Skælkarpe"},
		{Value: "Spejlkarpe", Label: "Spejlkarpe"},
		{Value: "Suder", Label: "Suder"},
		{Value: "Søørred", Label: "Søørred"},
		{Value: "Torsk", Label: "Torsk"},
		{Value: "Uden for kategori", Label: "Uden for kategori"},
		{Value: "Ulk", Label: "Ulk"},
	}

	speciesJson, err := json.Marshal(species)
	if err != nil {
		slog.Error("Failed to marshal list of users", "Error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(speciesJson)
}

func (h *Handler) HandleCatchPageIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse catch ID
	catchId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid catch ID", http.StatusNotFound)
		return
	}

	// Fetch data
	data, err := h.queries.GetPostByIDWithImages(ctx, catchId)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch post", "error", err, "id", catchId)
		http.Error(w, "Failed to fetch catch data", http.StatusInternalServerError)
		return
	}

	user, err := h.queries.GetUserByID(ctx, data.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch user", "error", err, "id", data.UserID)
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		return
	}

	// Parse images JSON
	var images []ImageData
	if err := json.Unmarshal(data.Images, &images); err != nil {
		slog.ErrorContext(ctx, "Failed to parse images JSON", "error", err, "id", catchId)
		http.Error(w, "Failed to process images", http.StatusInternalServerError)
		return
	}

	if len(images) == 0 {
		slog.ErrorContext(ctx, "No images found for post", "id", catchId)
		http.Error(w, "Catch has no images", http.StatusInternalServerError)
		return
	}

	// Find primary image or use first one
	var displayImage string
	for _, img := range images {
		if img.IsPrimary {
			displayImage = img.URL
			break
		}
	}
	if displayImage == "" {
		displayImage = images[0].URL
	}

	// Prepare view data
	catchData := web.CatchData{
		Title:       data.Title,
		Description: data.Description.String,
		Species:     data.Species,
		Length:      formatLength(data.LengthCm),
		Weight:      formatWeight(data.WeightKg),
		Image:       displayImage,
		User:        fmt.Sprintf("%s %s", user.FirstName.String, user.LastName.String),
	}

	// Render page
	if err := web.CatchPageIndex(catchData).Render(ctx, w); err != nil {
		slog.ErrorContext(ctx, "Error rendering catch page", "error", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func formatLength(length pgtype.Int4) string {
	if !length.Valid {
		return "N/A"
	}
	return fmt.Sprintf("%d cm", length.Int32)
}

func formatWeight(weight pgtype.Float8) string {
	if !weight.Valid {
		return "N/A"
	}
	return fmt.Sprintf("%.1f kg", weight.Float64)
}
