package fiber

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/quic-go/quic-go/http3"
)

type QUICFiber struct {
	*fiber.App
}

func NewQUICFiber(
	app *fiber.App,
) *QUICFiber {
	app.Use(func(c *fiber.Ctx) error {
		baseUrl := c.BaseURL()
		addr := strings.TrimPrefix(strings.TrimPrefix(baseUrl, "http://"), "https://")
		c.Set("Alt-Svc", fmt.Sprintf("h3=\"%s\"; ma=60", addr))
		return c.Next()
	})

	return &QUICFiber{
		App: app,
	}
}

func (f *QUICFiber) Listen(addr string, certFile string, keyFile string) error {
	errChan := make(chan error)
	go func() {
		if err := f.App.ListenTLS(addr, certFile, keyFile); err != nil {
			errChan <- err
			return
		}
	}()

	if err := http3.ListenAndServeQUIC(addr, certFile, keyFile, adaptor.FiberApp(f.App)); err != nil {
		return err
	}

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}
