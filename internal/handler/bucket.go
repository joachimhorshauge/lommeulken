package handler

import (
	"context"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/Backblaze/blazer/b2"
	"github.com/google/uuid"
)

// Initialize B2 client (do this once at startup)
func initB2Client(accountID, applicationKey string) (*b2.Client, error) {
	ctx := context.Background()
	return b2.NewClient(ctx, accountID, applicationKey)
}

func (h *Handler) uploadToB2(file io.Reader, header *multipart.FileHeader) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := uuid.New().String() + ext

	// Get bucket
	bucket, err := h.b2Client.Bucket(context.Background(), h.b2BucketName)
	if err != nil {
		return "", err
	}

	// Create new file in bucket
	obj := bucket.Object(filename)
	w := obj.NewWriter(context.Background())

	// Copy file data
	if _, err := io.Copy(w, file); err != nil {
		w.Close()
		return "", err
	}

	// Close writer to complete upload
	if err := w.Close(); err != nil {
		return "", err
	}

	// Get public URL (assuming your bucket is public)
	return h.b2BaseURL + "/file/" + h.b2BucketName + "/" + filename, nil
}
