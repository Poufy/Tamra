package middleware

import (
	"context"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/sirupsen/logrus"
)

// We pass the firebaseAuth as a parameter to the middleware so we can use it to verify the token and return a handler function that
// Take a http.Handler as a parameter and returns a http.Handler so we can continue the chain of handlers after the middleware.
func AuthMiddleware(firebaseAuth *auth.Client, logger logrus.FieldLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the token from the request header
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// Verify the token.
			// TODO: what is the use of the context.Background()?
			// tokenWithClaims, err := firebaseAuth.VerifyIDToken(context.Background(), token)
			// if err != nil {
			// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
			// 	return
			// }

			// Print all token details
			// fmt.Printf("Token: %v\n", tokenWithClaims)
			logger.Infof("Passed auth token: %s.", token)

			ctx := context.WithValue(r.Context(), "UID", "1234")

			// If the token is valid, we can continue the chain of handlers and pass the context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
