package logger

import "github.com/sirupsen/logrus"

var defaultLogger Logger

type Logger interface {
	logrus.FieldLogger
}

func DefaultLogger() Logger {
	if defaultLogger == nil {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})

		defaultLogger = logrus.StandardLogger()
	}
	return defaultLogger
}
