package swagger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Options struct {
	App *fiber.App
}

func UseSwagger(opts *Options) {
	opts.App.Get("/swagger/*", swagger.HandlerDefault)
}
