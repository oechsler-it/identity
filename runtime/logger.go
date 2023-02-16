package runtime

import (
	"os"
	"time"

	match "github.com/alexpantyukhin/go-pattern-match"
	"github.com/sirupsen/logrus"
)

var LOG_LEVEL = "LOG_LEVEL"

func NewLogger(env *Env) *logrus.Logger {
	logger := logrus.New()

	logger.SetOutput(os.Stdout)

	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
	})

	logger.SetLevel(func() logrus.Level {
		ok, value := match.Match(env.String(LOG_LEVEL)).
			When("debug", logrus.DebugLevel).
			When("warn", logrus.WarnLevel).
			When("error", logrus.ErrorLevel).
			When("fatal", logrus.FatalLevel).
			Result()

		if !ok {
			return logrus.InfoLevel
		}
		return value.(logrus.Level)
	}())

	logger.WithFields(logrus.Fields{
		"environment": env.Get(),
	}).Info("Using environment")

	return logger
}
