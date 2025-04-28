package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
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
		os.Exit(1)
	}
	slog.Info("Successfully connected to bucket")

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("BLUEPRINT_DB_HOST"),
		5432,
		os.Getenv("BLUEPRINT_DB_USERNAME"),
		os.Getenv("BLUEPRINT_DB_PASSWORD"),
		os.Getenv("BLUEPRINT_DB_DATABASE"),
	)

	// Parse the complete configuration
	dbConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		slog.Error("Failed to parse database config",
			"error", err,
			"connectionString", connStr)
		os.Exit(1)
	}

	// TCP settings
	dbConfig.ConnConfig.DialFunc = (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 5 * time.Minute,
	}).DialContext

	// Connection with retries
	var pool *pgxpool.Pool
	maxAttempts := 5
	for i := range maxAttempts {
		pool, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
		if err == nil {
			// Verify connection
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = pool.Ping(ctx)
			if err == nil {
				break
			}
		}

		if i < maxAttempts-1 {
			slog.Warn("Database connection failed, retrying...",
				"attempt", i+1,
				"error", err,
				"connectionDetails", fmt.Sprintf("%s@%s:%d/%s",
					dbConfig.ConnConfig.User,
					dbConfig.ConnConfig.Host,
					dbConfig.ConnConfig.Port,
					dbConfig.ConnConfig.Database))
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}
	defer pool.Close()
	slog.Info("Successfully connected to database")

	queries := dbstore.New(pool)

	handler := handler.NewHandler(
		b2Client,
		os.Getenv("B2_BUCKET_NAME"),
		os.Getenv("B2_BASE_URL"),
		queries,
	)

	middleware := middleware.NewMiddleware(queries)
	server := server.NewServer(handler, middleware)

	done := make(chan bool, 1)
	go gracefulShutdown(server, done)

	slog.Info("Started server", "port", os.Getenv("PORT"))
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Server error", "Error", err)
		os.Exit(1)
	}

	<-done
	log.Println("Graceful shutdown complete.")
}
