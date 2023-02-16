package fiber

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/runtime"
	"github.com/sirupsen/logrus"
)

func NewFiber() *fiber.App {
	return fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
}

type Options struct {
	Hooks  *runtime.Hooks
	Logger *logrus.Logger
	App    *fiber.App
}

func UseFiber(opts *Options) {
	opts.Hooks.OnStart(func(ctx context.Context) error {
		go func() {
			err := opts.App.Listen(":3000")
			if err != nil {
				opts.Logger.WithError(err).
					Fatal("Fiber failed")
			}
		}()
		opts.Logger.WithField("address", "http://localhost:3000").
			Info("Fiber is listening")
		return nil
	})

	opts.Hooks.OnStop(func(ctx context.Context) error {
		return opts.App.Shutdown()
	})
}
