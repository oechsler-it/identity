package modules

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/oechsler-it/identity/modules/session"
	"github.com/oechsler-it/identity/modules/user"
)

type Options struct {
	App     *fiber.App
	User    *user.Options
	Session *session.Options
}

func UseModules(opts *Options) {
	user.UseUser(opts.User)
	session.UseSession(opts.Session)
}

var WireModules = wire.NewSet(
	wire.Struct(new(Options), "*"),
	user.WireUser,
	session.WireSession,
)
