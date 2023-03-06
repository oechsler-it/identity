package fiber

import (
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"time"
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
	deviceId := ctx.Cookies("device_id")
	deviceIdAsUuid, err := uuid.FromString(deviceId)
	if err != nil {
		deviceIdAsUuid = uuid.NewV4()

		e.Logger.WithFields(logrus.Fields{
			"deviceId": deviceIdAsUuid.String(),
		}).Info("Registered new device id")
	}

	ctx.Locals("device_id", deviceIdAsUuid)

	ctx.Cookie(&fiber.Cookie{
		Name:    "device_id",
		Value:   deviceIdAsUuid.String(),
		Path:    "/",
		Expires: time.Now().Add(31 * 24 * time.Hour),
	})

	return ctx.Next()
}
