package gorm

import (
	"github.com/oechsler-it/identity/runtime"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Options struct {
	Hooks  *runtime.Hooks
	Env    *runtime.Env
	Logger *logrus.Logger
}

func NewPostgres(opts *Options) *gorm.DB {
	dsn := opts.Env.String("POSTGRES_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
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
