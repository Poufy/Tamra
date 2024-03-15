package middleware

import (
	"context"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/sirupsen/logrus"
)

// We pass the firebaseAuth as a parameter to the middleware so we can use it to verify the token and return a handler function that
// Take a http.Handler as a parameter and returns a http.Handler so we can continue the chain of handlers after the middleware.
func UserAuthMiddleware(firebaseAuth *auth.Client, logger logrus.FieldLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the token from the request header
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				logger.Warn("No auth token provided.")
				return
			}

			// Verify the token.
			// TODO: what is the use of the context.Background()?
			tokenWithClaims, err := firebaseAuth.VerifyIDToken(r.Context(), token)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// Print all token details
			// fmt.Printf("Token: %v\n", tokenWithClaims)
			logger.Infof("Passed auth token: %+v.", tokenWithClaims)
			ctx := context.WithValue(r.Context(), "UID", tokenWithClaims.UID)

			// If the token is valid, we can continue the chain of handlers and pass the context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RestaurantAuthMiddleware(firebaseAuth *auth.Client, logger logrus.FieldLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the token from the request header
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				logger.Warn("No auth token provided.")
				return
			}

			// Verify the token.
			// TODO: what is the use of the context.Background()?
			// TODO: should this context from the request or should we create it from the context.Background()?
			tokenWithClaims, err := firebaseAuth.VerifyIDToken(r.Context(), token)
			if err != nil {
				logger.Error("Failed to verify token: %v", err)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if the claims contain an email
			// If the claims contain an email, we can assume that the user is a restaurant
			if tokenWithClaims.Claims["email"] == nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			logger.Info("Passed auth token: %+v.", tokenWithClaims)
			// print the email
			logger.Info("Email: %v", tokenWithClaims.Claims["email"])

			ctx := context.WithValue(r.Context(), "UID", tokenWithClaims.UID)

			// If the token is valid, we can continue the chain of handlers and pass the context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
