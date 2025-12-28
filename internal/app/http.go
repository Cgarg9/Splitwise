package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"splitwise-clone/internal/httpapi/router"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// HTTPServer represents the HTTP server
type HTTPServer struct {
	server *http.Server
	router *router.Router
}

// NewHTTPServer creates a new HTTP server instance
func NewHTTPServer(port int, db *pgxpool.Pool) *HTTPServer {
	r := router.NewRouter(db)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &HTTPServer{
		server: srv,
		router: r,
	}
}

// Start starts the HTTP server
func (s *HTTPServer) Start() error {
	log.Info().Str("addr", s.server.Addr).Msg("Starting HTTP server")
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}
