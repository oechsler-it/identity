package modules

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/oechsler-it/identity/modules/permission"
	"github.com/oechsler-it/identity/modules/session"
	"github.com/oechsler-it/identity/modules/user"
)

type Options struct {
	App        *fiber.App
	User       *user.Options
	Session    *session.Options
	Permission *permission.Options
}

func UseModules(opts *Options) {
	user.UseUser(opts.User)
	session.UseSession(opts.Session)
	permission.UsePermission(opts.Permission)
}

var WireModules = wire.NewSet(
	wire.Struct(new(Options), "*"),
	user.WireUser,
	session.WireSession,
	permission.WirePermission,
)
