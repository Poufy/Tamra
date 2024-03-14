package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type FirebaseApp struct {
	App *firebase.App
}

func NewFirebaseApp(configJSON string) *FirebaseApp {
	opt := option.WithCredentialsJSON([]byte(configJSON))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil
	}
	return &FirebaseApp{App: app}
}

func (fb FirebaseApp) FetchFirebaseAuthClient() (*auth.Client, error) {
	authClient, err := fb.App.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create a new firebase auth client: %w", err)
	}
	return authClient, nil
}

func (fb FirebaseApp) FetchFirebaseMessagingClient() (*messaging.Client, error) {
	messagingClient, err := fb.App.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create a new firebase messaging client: %w", err)
	}
	return messagingClient, nil
}
