package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"lommeulken/internal/database"
	"lommeulken/internal/handler"
	"lommeulken/internal/middleware"
)

type Server struct {
	port       int
	db         database.Service
	handler    *handler.Handler
	middleware *middleware.Middleware
}

func NewServer(h *handler.Handler, m *middleware.Middleware) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	NewServer := &Server{
		port:       port,
		db:         database.New(),
		handler:    h,
		middleware: m,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
