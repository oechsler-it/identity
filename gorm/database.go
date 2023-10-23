package gorm

import (
	match "github.com/alexpantyukhin/go-pattern-match"
	"github.com/oechsler-it/identity/runtime"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const LOG_LEVEL = "LOG_LEVEL"

type Options struct {
	Hooks  *runtime.Hooks
	Env    *runtime.Env
	Logger *logrus.Logger
}

func NewPostgres(opts *Options) *gorm.DB {
	dsn := opts.Env.String("POSTGRES_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(func() logger.LogLevel {
			ok, log := match.Match(opts.Env.String(LOG_LEVEL)).
				When("debug", logger.Info).
				Result()

			if !ok {
				return logger.Silent
			}
			return log.(logger.LogLevel)
		}()),
		PrepareStmt: true,
	})
	if err != nil {
		opts.Logger.WithError(err).
			Fatal("Failed to connect to database")
	}
	opts.Logger.WithFields(logrus.Fields{
		"name": db.Name(),
	}).Info("Connected to database")
	return db
}
