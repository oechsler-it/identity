package modules

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/oechsler-it/identity/modules/middleware"
	"github.com/oechsler-it/identity/modules/permission"
	"github.com/oechsler-it/identity/modules/session"
	"github.com/oechsler-it/identity/modules/token"
	"github.com/oechsler-it/identity/modules/user"
)

type Options struct {
	App        *fiber.App
	Token      *token.Options
	Session    *session.Options
	User       *user.Options
	Permission *permission.Options
}

func UseModules(opts *Options) {
	token.UseTokenMiddleware(opts.Token)
	session.UseSessionMiddleware(opts.Session)

	token.UseToken(opts.Token)
	session.UseSession(opts.Session)
	user.UseUser(opts.User)
	permission.UsePermission(opts.Permission)
}

var WireModules = wire.NewSet(
	wire.Struct(new(Options), "*"),
	middleware.WireMiddleware,
	token.WireToken,
	session.WireSession,
	user.WireUser,
	permission.WirePermission,
)
