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

type FiberOptions struct {
	Hooks  *runtime.Hooks
	Logger *logrus.Logger
	App    *fiber.App
}

func UseFiber(options *FiberOptions) {
	options.Hooks.OnStart(func(ctx context.Context) error {
		go func() {
			err := options.App.Listen(":3000")
			if err != nil {
				options.Logger.WithError(err).
					Fatal("Fiber failed")
			}
		}()
		options.Logger.WithField("address", "http://localhost:3000").
			Info("Fiber is listening")
		return nil
	})

	options.Hooks.OnStop(func(ctx context.Context) error {
		return options.App.Shutdown()
	})
}
