package logger

import (
	"io"
	"time"

	"github.com/pochard/logrotator"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// Hook is a hook that writes logs of specified LogLevels to specified Writer
type Hook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
	formatter logrus.Formatter
	cron      *cron.Cron
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *Hook) Fire(entry *logrus.Entry) error {
	// line, err := entry.Bytes()
	line, err := hook.formatter.Format(entry) // entry.Bytes()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *Hook) Levels() []logrus.Level {
	return hook.LogLevels
}

func (hook *Hook) Close() {
	hook.cron.Stop()
}

func NewHook() *Hook {
	writer, err := logrotator.NewTimeBasedRotator("./logs/%Y%m%d-%H%M.log", 1*time.Hour)
	if err != nil {
		logrus.Fatal("create logrotator writer error:", err)
	}

	cleaner, err := logrotator.NewTimeBasedCleaner("./logs/*.log", 3*24*time.Hour)
	if err != nil {
		logrus.Fatal("create logrotator cleaner error:", err)
	}

	c := cron.New()
	c.AddFunc("5 1 * * *", func() {
		deleted, err := cleaner.Clean()
		if err != nil {
			logrus.Info("%v\n", err)
			return
		}

		for _, d := range deleted {
			logrus.Infof("%s deleted\n", d)
		}
	})
	c.Start()

	hook := &Hook{ // Send logs with level higher than warning to stderr
		Writer:    writer,
		formatter: &logrus.TextFormatter{},
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
			// logrus.DebugLevel,
		},
		cron: c,
	}
	return hook
}
