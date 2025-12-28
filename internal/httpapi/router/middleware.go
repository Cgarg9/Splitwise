package router

import (
	"net/http"
	"runtime/debug"
	"splitwise-clone/internal/logger"
	"time"

	"github.com/google/uuid"
)

// TraceIDMiddleware injects a trace ID into the request context
func TraceIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if trace ID exists in header, otherwise generate new one
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Add trace ID to response header for client tracking
		w.Header().Set("X-Trace-ID", traceID)

		// Add trace ID to context
		ctx := logger.WithTraceID(r.Context(), traceID)

		// Also add request ID (can be same as trace ID or separate)
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = traceID
		}
		ctx = logger.WithRequestID(ctx, requestID)

		// Continue with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggingMiddleware logs HTTP requests with trace ID from context
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture status code
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Get logger from context (includes trace ID)
		log := logger.FromContext(r.Context())

		// Log request start
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Msg("HTTP request started")

		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		// Log request completion
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Int("status", ww.statusCode).
			Dur("duration", duration).
			Int64("duration_ms", duration.Milliseconds()).
			Msg("HTTP request completed")
	})
}

// RecoveryMiddleware recovers from panics and logs them with trace ID
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Get logger from context
				log := logger.FromContext(r.Context())

				log.Error().
					Interface("error", err).
					Str("stack", string(debug.Stack())).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Msg("Panic recovered")

				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal Server Error"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// responseWriter is a custom ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
