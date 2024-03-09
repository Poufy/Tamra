package services

import (
	"Tamra/internal/pkg/models"

	"github.com/sirupsen/logrus"
)

func NotifyUser(fcmToken string, order *models.Order) error {
	// Here we would send a notification to the user using the fcmToken
	logrus.Info("Sending notification to user: ", fcmToken)
	logrus.Info("Order: ", order)
	return nil
}
