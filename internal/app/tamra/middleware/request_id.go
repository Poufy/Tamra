package middleware

import (
	"net/http"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// RequestIDMiddleware is a middleware that sets the request ID in the response header in order to be able to trace requests
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the request ID from the context
		requestID := r.Context().Value(chimiddleware.RequestIDKey).(string)

		// Set the request ID in the response header
		w.Header().Set("X-Request-ID", requestID)

		// Continue the chain of handlers
		next.ServeHTTP(w, r)
	})
}
