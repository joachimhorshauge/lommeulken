package handler

import (
	"context"
	"io"
	"log/slog"
	"lommeulken/cmd/view"
	"lommeulken/cmd/web"
	"lommeulken/gen/dbstore"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) CatchIndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := web.CatchesIndex().Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render Catches index form with errors", "error", err)
	}
	return
}

func (h *Handler) CatchCards(w http.ResponseWriter, r *http.Request) {

	filterByUserID := false
	userID := uuid.New()
	filterBySpecies := false
	species := []string{}
	sortColumn := "created_at"
	sortDirection := "DESC"
	offset := 0
	limit := 10

	// TODO: Get limit and offset from query params

	requestUrl, _ := url.Parse(r.URL.String())
	urlParams, _ := url.ParseQuery(requestUrl.RawQuery)

	if urlParams["filter.user"] != nil {
		filterByUserID = true
		userID = uuid.MustParse(urlParams["filter.user"][0])
	}

	if urlParams["filter.species"] != nil {
		filterBySpecies = true
		species = urlParams["filter.species"]
	}

	if urlParams["sort.length"] != nil {
		sortColumn = "length_cm"
		sortDirection = urlParams["sort.length"][0]
	}

	if urlParams["sort.weight"] != nil {
		sortColumn = "weight_kg"
		sortDirection = urlParams["sort.weight"][0]
	}

	if urlParams["sort.dateCaught"] != nil {
		sortColumn = "created_at"
		sortDirection = urlParams["sort.dateCaught"][0]
	}

	listPostsArgs := dbstore.ListPostsWithImagesParams{
		FilterByUserID:  filterByUserID,
		UserID:          userID,
		FilterBySpecies: filterBySpecies,
		Species:         species,
		ResultOffset:    int32(offset),
		ResultLimit:     int32(limit),
		SortColumn:      sortColumn,
		SortDirection:   sortDirection,
	}

	posts, err := h.queries.ListPostsWithImages(context.Background(), listPostsArgs)
	if err != nil {
		slog.Error("Failed to get 10 latest posts", "msg", err)
		posts = []dbstore.ListPostsWithImagesRow{}
	}

	cardInfoList := PostsToCardInfo(posts)

	err = web.CatchCards(cardInfoList).Render(r.Context(), w)
	if err != nil {
		slog.Error("Error rendering Signup page", "error", err)
		return
	}
}

func (h *Handler) NewCatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := web.NewCatch(web.CatchInfo{}, web.CatchErrors{}).Render(r.Context(), w)
		if err != nil {
			slog.Error("Failed to load new catch template")
			return
		}
		return
	} else if r.Method == http.MethodPost {
		errors := web.CatchErrors{}

		userID := view.AuthenticatedUser(r.Context()).ID

		length, err := strconv.Atoi(r.FormValue("length_cm"))
		if err != nil {
			errors.InvalidLength = "Ugyldig længde"
		}

		weight, err := strconv.ParseFloat(r.FormValue("weight_kg"), 64)
		if err != nil {
			weight = 0.0
		}

		parsedDate, err := time.Parse("2006-01-02", r.FormValue("date"))
		if err != nil {
			errors.InvalidDate = "Ugyldig dato"
		}

		species := r.FormValue("species")
		if species == "" {
			errors.NoSpecies = "Ingen valgt art"
		}

		title := r.FormValue("title")
		if title == "" {
			errors.NoTitle = "Ingen overskrift"
		}

		description := r.FormValue("description")
		if description == "" {
			errors.NoDescription = "Ingen Beretning"
		}

		// Parse multipart form (max 10MB)
		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			w.Header().Set("HX-Retarget", "#error-message")
			w.Header().Set("HX-Reswap", "innerHTML")
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		var imgUrl string
		// Handle file upload
		file, header, err := r.FormFile("image")
		if err != nil {
			errors.NoImage = "Mangler billede"
		} else {
			defer file.Close()

			// Validate file size (5MB max)
			if header.Size > 5<<20 {
				errors.FileTooLarge = "Filen er for stor (max 5MB)"
			}

			// Validate file type
			buff := make([]byte, 512)
			_, err = file.Read(buff)
			if err != nil {
				slog.Error("Error reading file", "Error", err)
				return
			}

			filetype := http.DetectContentType(buff)
			if !strings.HasPrefix(filetype, "image/") {
				errors.WrongFileType = "Forkert filtype"
			}

			// Reset file pointer
			_, err = file.Seek(0, io.SeekStart)
			if err != nil {
				slog.Error("Error processing file", "Error", err)
				return
			}

			// Upload to Backblaze B2
			url, err := h.uploadToB2(file, header)
			imgUrl = strings.Replace(url, "file/lommeulken/", "", -1)
			if err != nil {
				slog.Error("Failed to upload image", "Error", err)
				return
			}
			slog.Info("Successfully uploaded image to bucket", "url", imgUrl)
		}

		if (errors == web.CatchErrors{}) {

			params := dbstore.CreatePostParams{
				UserID:      userID,
				TripID:      pgtype.UUID{Valid: false},
				Title:       title,
				Description: pgtype.Text{String: description, Valid: true},
				Species:     species,
				LengthCm:    pgtype.Int4{Int32: int32(length), Valid: true},
				WeightKg:    pgtype.Float8{Float64: weight, Valid: weight != 0.0},
				CatchDate:   pgtype.Timestamptz{Time: parsedDate, Valid: true}}

			post, err := h.queries.CreatePost(context.Background(), params)
			if err != nil {
				slog.Error("Error saving post to database", "Error", err, "post", post)
			}

			imageParams := dbstore.AddPostImageParams{
				PostID:    post.ID,
				Url:       imgUrl,
				IsPrimary: pgtype.Bool{Bool: true, Valid: true},
			}
			postImage, err := h.queries.AddPostImage(context.Background(), imageParams)
			if err != nil {
				slog.Error("Error saving post image to database", "Error", err, "Image", postImage)
			}
			w.Header().Add("Hx-Redirect", "/catches")

		}

		catchInfo := web.CatchInfo{
			Title:       r.FormValue("Title"),
			Description: r.FormValue("description"),
			Date:        r.FormValue("date"),
			Species:     r.FormValue("species"),
			Length:      r.FormValue("length_cm"),
			Weight:      r.FormValue("weight_kg"),
		}

		err = web.NewCatchForm(catchInfo, errors).Render(r.Context(), w)
		if err != nil {
			slog.Error("Error rendering NewCatchForm form with errors", "error", err)
		}
		return

	}
	return
}
