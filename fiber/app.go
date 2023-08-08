package fiber

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/runtime"
	"github.com/sirupsen/logrus"
)

func NewFiber(
	env *runtime.Env,
	logger *logrus.Logger,
) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Hooks().OnListen(func(data fiber.ListenData) error {
		if fiber.IsChild() {
			return nil
		}

		scheme := "http"
		if data.TLS {
			scheme = "https"
		}
		addr := fmt.Sprintf("%s://%s:%s", scheme, data.Host, data.Port)

		http3 := env.Bool("FIBER_HTTP3", false)
		if scheme == "https" && http3 {
			logger.WithField("address", addr).Info("Fiber is listening (with HTTP/3 adaptor)")
			return nil
		}
		if scheme == "https" {
			logger.WithField("address", addr).Info("Fiber is listening (with TLS)")
			return nil
		}
		logger.WithField("address", addr).Info("Fiber is listening")
		return nil
	})

	return app
}

type Options struct {
	Env     *runtime.Env
	Hooks   *runtime.Hooks
	Logger  *logrus.Logger
	App     *fiber.App
	QuicApp *QUICFiber
}

func UseFiber(opts *Options) {
	addr := opts.Env.String("FIBER_ADDR", ":3000")

	opts.Hooks.OnStart(func(ctx context.Context) error {
		cert := opts.Env.String("FIBER_CERT", "")
		key := opts.Env.String("FIBER_KEY", "")
		http3 := opts.Env.Bool("FIBER_HTTP3", false)

		go func() {
			if http3 && cert != "" && key != "" {
				if err := opts.QuicApp.Listen(addr, cert, key); err != nil {
					opts.Logger.WithError(err).
						Fatal("Fiber failed")
				}
				return
			}
			if cert != "" && key != "" {
				if err := opts.App.ListenTLS(addr, cert, key); err != nil {
					opts.Logger.WithError(err).
						Fatal("Fiber failed")
				}
				return
			}
			if err := opts.App.Listen(addr); err != nil {
				opts.Logger.WithError(err).
					Fatal("Fiber failed")
			}
		}()
		return nil
	})

	opts.Hooks.OnStop(func(ctx context.Context) error {
		return opts.App.Shutdown()
	})
}
