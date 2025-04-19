package handler

import (
	"github.com/Backblaze/blazer/b2"
)

type Handler struct {
	// Backblaze B2 client and configuration
	b2Client     *b2.Client
	b2BucketName string
	b2BaseURL    string
}

// NewHandler creates a new Handler instance with all required dependencies
func NewHandler(b2Client *b2.Client, bucketName, baseURL string) *Handler {
	return &Handler{
		b2Client:     b2Client,
		b2BucketName: bucketName,
		b2BaseURL:    baseURL,
		// db:           db,
	}
}
