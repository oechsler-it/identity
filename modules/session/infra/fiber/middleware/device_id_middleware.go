package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oechsler-it/identity/modules/session/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type DeviceIdMiddleware struct {
	*fiber.App
	// ---
	Logger *logrus.Logger
}

func UseDeviceIdMiddleware(middleware *DeviceIdMiddleware) {
	middleware.Use(middleware.handle)
}

func (e *DeviceIdMiddleware) handle(ctx *fiber.Ctx) error {
	deviceIdCookie := ctx.Cookies("device_id")
	deviceId, err := uuid.FromString(deviceIdCookie)
	if err != nil {
		deviceId = uuid.NewV4()

		e.Logger.WithFields(logrus.Fields{
			"device_id": deviceId.String(),
		}).Info("Registered new device id")
	}

	ctx.Locals("device_id", domain.DeviceId(deviceId))

	ctx.Cookie(&fiber.Cookie{
		Name:    "device_id",
		Value:   deviceId.String(),
		Path:    "/",
		Expires: time.Now().Add(31 * 24 * time.Hour),
	})

	return ctx.Next()
}
