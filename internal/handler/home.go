package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"lommeulken/cmd/web"
	"lommeulken/gen/dbstore"
	"net/http"
	"time"
)

func (h *Handler) HandleHomeIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := web.HomeIndex().Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render HomeIndex form with errors", "error", err)
	}
	return
}

func PostsToCardInfo(posts []dbstore.ListPostsWithImagesRow) []web.CardInfo {
	cards := []web.CardInfo{}

	for _, post := range posts {
		var images []struct {
			ID        string    `json:"id"`
			URL       string    `json:"url"`
			IsPrimary bool      `json:"is_primary"`
			CreatedAt time.Time `json:"created_at"`
		}

		if err := json.Unmarshal(post.Images, &images); err != nil {
			slog.Warn("Failed to unmarshal post images", "postID", post.ID, "error", err)
		}

		var imageURL string
		for _, img := range images {
			if img.IsPrimary {
				imageURL = img.URL
				break
			}
		}
		if imageURL == "" && len(images) > 0 {
			imageURL = images[0].URL
		}

		cards = append(cards, web.CardInfo{
			ImageUrl: imageURL,
			Species:  post.Species,
			Length:   fmt.Sprintf("%d cm", post.LengthCm.Int32),
			Weight:   fmt.Sprintf("%.2f kg", post.WeightKg.Float64),
			Title:    post.Title,
			PostUrl:  fmt.Sprintf("/catch/%s", post.ID.String()),
		})
	}

	return cards
}
