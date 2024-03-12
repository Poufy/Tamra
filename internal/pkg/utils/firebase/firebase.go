package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func NewFirebaseAuth(configJSON string) (*auth.Client, error) {
	// Load the Firebase configuration from the file
	opt := option.WithCredentialsJSON([]byte(configJSON))

	// NewApp will first look for the FIREBASE_CONFIG environment variable.
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new firebase app: %w", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create a new firebase auth client: %w", err)
	}
	return authClient, nil
}
