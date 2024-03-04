package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger(level string) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(lvl)
	}

	return logger
}
