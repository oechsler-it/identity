package modules

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/oechsler-it/identity/modules/permission"
	"github.com/oechsler-it/identity/modules/session"
	"github.com/oechsler-it/identity/modules/token"
	"github.com/oechsler-it/identity/modules/user"
)

type Options struct {
	App        *fiber.App
	User       *user.Options
	Session    *session.Options
	Permission *permission.Options
	Token      *token.Options
}

func UseModules(opts *Options) {
	session.UseSession(opts.Session)
	user.UseUser(opts.User)
	permission.UsePermission(opts.Permission)
	token.UseToken(opts.Token)
}

var WireModules = wire.NewSet(
	wire.Struct(new(Options), "*"),
	user.WireUser,
	session.WireSession,
	permission.WirePermission,
	token.WireToken,
)
