package logger

import (
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init(level string) {
	switch level {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}

	Log.SetFormatter(&logrus.JSONFormatter{})
}

func GetLogger() *logrus.Logger {
	return Log
}
