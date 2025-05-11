package server

import (
	"lommeulken/cmd/web"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", s.handler.HandleHomeIndex)

	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", fileServer)
	mux.HandleFunc("/signup", s.handler.HandleSignup)
	mux.HandleFunc("/login", s.handler.HandleLogin)
	mux.HandleFunc("/logout", s.handler.HandleLogout)
	mux.HandleFunc("/catches", s.handler.CatchIndexHandler)
	mux.HandleFunc("/catches/new", s.handler.NewCatchHandler)
	mux.HandleFunc("/catches/cards", s.handler.CatchCards)
	mux.HandleFunc("/catch/{id}", s.handler.HandleCatchPageIndex)
	mux.HandleFunc("/about", s.handler.AboutIndexHandler)

	mux.HandleFunc("/api/users", s.handler.HandleListUsers)
	mux.HandleFunc("/api/species", s.handler.HandleListSpecies)

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(s.middleware.WithUser(s.middleware.WithLogging(mux)))
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}
