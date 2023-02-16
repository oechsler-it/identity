package runtime

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

type Runtime struct {
	hooks  *Hooks
	logger *logrus.Logger
}

func NewRuntime(
	hooks *Hooks,
	logger *logrus.Logger,
) *Runtime {
	return &Runtime{
		hooks:  hooks,
		logger: logger,
	}
}

func (r *Runtime) Run() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	if err := r.hooks.start(ctx); err != nil {
		r.logger.WithError(err).Fatal("Failed to run start hooks")
	}

	<-ctx.Done()

	if err := r.hooks.stop(ctx); err != nil {
		r.logger.WithError(err).Fatal("Failed to run stop hooks")
	}
}
