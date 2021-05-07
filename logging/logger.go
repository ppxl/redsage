package logging

import (
	"github.com/sirupsen/logrus"
)

var (
	loggerInstance = logrus.New()
)

type AppLogger struct {
	logger *logrus.Logger
}

func Init(logLevel logrus.Level) error {
	loggerInstance := logrus.New()
	loggerInstance.SetLevel(logLevel)

	return nil
}

func Logger() *logrus.Logger {
	return loggerInstance
}
