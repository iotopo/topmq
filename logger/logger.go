package logger

import (
	"github.com/iotopo/topmq/config"
	"github.com/sirupsen/logrus"
	"log"
)

var fileLoggerHook *Hook

func Init() {
	logger := logrus.StandardLogger()
	logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: false, FullTimestamp: true, ForceColors: true, TimestampFormat: "2006-01-02 15:04:05.000"})
	log.SetOutput(logger.Writer())

	if config.Conf.Logger.Level != "" {
		if level, err := logrus.ParseLevel(config.Conf.Logger.Level); err != nil {
			logrus.WithError(err).Fatalf("failed parsing LOG_LEVEL")
		} else {
			logger.SetLevel(level)
		}
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	if config.Conf.Logger.FileLog {
		fileLoggerHook = NewHook()
		logrus.AddHook(fileLoggerHook)
	}
}

func Close() {
	if fileLoggerHook != nil {
		fileLoggerHook.Close()
	}
}
