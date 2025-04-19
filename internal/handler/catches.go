package handler

import (
	"io"
	"log/slog"
	"lommeulken/cmd/web"
	"net/http"
	"strings"
)

func (h *Handler) NewCatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := web.NewCatch().Render(r.Context(), w)
		if err != nil {
			slog.Error("Failed to load new catch template")
			return
		}
		return
	} else if r.Method == http.MethodPost {
		// Parse multipart form (max 10MB)
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			w.Header().Set("HX-Retarget", "#error-message")
			w.Header().Set("HX-Reswap", "innerHTML")
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Get other form values
		// date := r.FormValue("date")
		// species := r.FormValue("make")
		// length := r.FormValue("length_cm")
		// weight := r.FormValue("weight_kg")
		// title := r.FormValue("title")
		// description := r.FormValue("description")

		// Handle file upload
		file, header, err := r.FormFile("image")
		if err == nil { // If file was uploaded
			defer file.Close()

			// Validate file size (5MB max)
			if header.Size > 5<<20 {
				http.Error(w, "File too large (max 5MB)", http.StatusBadRequest)
				return
			}

			// Validate file type
			buff := make([]byte, 512)
			_, err = file.Read(buff)
			if err != nil {
				http.Error(w, "Error reading file", http.StatusInternalServerError)
				return
			}

			filetype := http.DetectContentType(buff)
			if !strings.HasPrefix(filetype, "image/") {
				http.Error(w, "Invalid file type", http.StatusBadRequest)
				return
			}

			// Reset file pointer
			_, err = file.Seek(0, io.SeekStart)
			if err != nil {
				http.Error(w, "Error processing file", http.StatusInternalServerError)
				return
			}

			// Upload to Backblaze B2
			imageURL, err := h.uploadToB2(file, header)
			if err != nil {
				http.Error(w, "Failed to upload image", http.StatusInternalServerError)
				return
			}
			slog.Info("Successfully uploaded image to bucket", "url", imageURL)
		}

		// Save catch data to your database including imageURL
		// ... your database logic here ...

		// Return success response
		w.WriteHeader(http.StatusCreated)
	}
	return
}
