package handler

import (
	"context"
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

	listPostsArgs := dbstore.ListPostsWithImagesParams{
		Limit:  10,
		Offset: 0,
	}
	posts, err := h.queries.ListPostsWithImages(context.Background(), listPostsArgs)
	if err != nil {
		slog.Error("Failed to get 10 latest posts", "msg", err)
		posts = []dbstore.ListPostsWithImagesRow{}
	}

	cardInfoList := postsToCardInfo(posts)

	err = web.HomeIndex(cardInfoList).Render(r.Context(), w)
	if err != nil {
		slog.Error("Error rendering Signup page", "error", err)
		return
	}
	return

}

func postsToCardInfo(posts []dbstore.ListPostsWithImagesRow) []web.CardInfo {
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
		})
	}

	return cards
}
