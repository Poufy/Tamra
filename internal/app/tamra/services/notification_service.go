package services

import (
	"context"
	"fmt"

	"firebase.google.com/go/messaging"
	"github.com/sirupsen/logrus"
)

type NotificationService interface {
	NotifyUser(fcmToken string, title string, body string) error
}

type NotificationServiceImpl struct {
	logger          logrus.FieldLogger
	messagingClient *messaging.Client
}

func NewNotificationService(logger logrus.FieldLogger, messagingClient *messaging.Client) NotificationService {
	return &NotificationServiceImpl{logger: logger, messagingClient: messagingClient}
}

func (ns NotificationServiceImpl) NotifyUser(fcmToken string, title string, body string) error {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "New Order Received",
			Body:  "You have a new order",
		},
		Token: fcmToken,
	}

	_, err := ns.messagingClient.Send(context.Background(), message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}
