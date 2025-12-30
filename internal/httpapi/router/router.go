package router

import (
	"net/http"
	"splitwise-clone/internal/domain/auth"
	"splitwise-clone/internal/httpapi/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	_ "splitwise-clone/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// Router holds the HTTP router and its dependencies
type Router struct {
	chi.Router
	db *pgxpool.Pool
}

// NewRouter creates and configures a new HTTP router
func NewRouter(db *pgxpool.Pool) *Router {
	r := chi.NewRouter()

	// Middleware setup (order matters!)
	r.Use(TraceIDMiddleware)  // First: inject trace ID into context
	r.Use(RecoveryMiddleware) // Second: catch panics (needs trace ID)
	r.Use(LoggingMiddleware)  // Third: log requests (needs trace ID)

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "http://127.0.0.1:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router := &Router{
		Router: r,
		db:     db,
	}

	router.setupRoutes()

	return router
}

// setupRoutes configures all application routes
func (router *Router) setupRoutes() {
	// Initialize repositories
	authRepo := auth.NewRepository(router.db)

	// Initialize services
	authService := auth.NewService(authRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // The url pointing to API definition
	))

	// Health check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Error().Err(err).Msg("Failed to write health check response")
		}
	})

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		// Auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", authHandler.SignUp)
			r.Post("/login", authHandler.Login)
			// r.Post("/logout", authHandler.Logout)   // For future implementation
		})

		// Protected routes (will need authentication middleware)
		// r.Group(func(r chi.Router) {
		// 	r.Use(AuthMiddleware) // Add JWT middleware here
		// 	r.Get("/users/me", userHandler.GetCurrentUser)
		// })
	})

	// Log registered routes
	log.Info().Msg("Routes registered successfully")
}
