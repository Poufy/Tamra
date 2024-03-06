package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/sirupsen/logrus"

	"google.golang.org/api/option"
)

func NewFirebaseAuth(configFile string) (*auth.Client, error) {
	logrus.Println("Initializing firebase auth with config file", configFile)
	// Load the Firebase configuration from the file
	opt := option.WithCredentialsFile(configFile)

	logrus.Println("file that was loaded", configFile)
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
