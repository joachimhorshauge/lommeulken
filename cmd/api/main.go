package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lommeulken/gen/dbstore"
	"lommeulken/internal/handler"
	"lommeulken/internal/middleware"
	"lommeulken/internal/server"

	"github.com/Backblaze/blazer/b2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	slog.Info("Starting application setup...")
	// Initialize Backblaze B2 client
	b2Client, err := b2.NewClient(
		context.Background(),
		os.Getenv("B2_APPLICATION_KEY_ID"),
		os.Getenv("B2_APPLICATION_KEY"),
	)
	if err != nil {
		slog.Error("Failed to initialize B2 client", "Error", err)
	}
	slog.Info("Successfully connected to bucket")

	ctx := context.Background()

	conn, err := pgxpool.New(ctx, os.Getenv("GOOSE_DBSTRING"))
	if err != nil {
		slog.Error("Failed to connect to database", "Error", err)
	}
	defer conn.Close()
	slog.Info("Successfully connected to database")

	queries := dbstore.New(conn)

	handler := handler.NewHandler(
		b2Client,
		os.Getenv("B2_BUCKET_NAME"),
		os.Getenv("B2_BASE_URL"),
		queries,
	)

	middleware := middleware.NewMiddleware(queries)

	server := server.NewServer(handler, middleware)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	slog.Info("Started server", "port", os.Getenv("PORT"))
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
