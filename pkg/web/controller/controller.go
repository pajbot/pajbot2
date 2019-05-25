package controller

import (
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/common/config"
	"github.com/pajbot/pajbot2/pkg/web/controller/admin"
	"github.com/pajbot/pajbot2/pkg/web/controller/api"
	"github.com/pajbot/pajbot2/pkg/web/controller/banphrases"
	"github.com/pajbot/pajbot2/pkg/web/controller/channel"
	"github.com/pajbot/pajbot2/pkg/web/controller/dashboard"
	"github.com/pajbot/pajbot2/pkg/web/controller/home"
	"github.com/pajbot/pajbot2/pkg/web/controller/logout"
	"github.com/pajbot/pajbot2/pkg/web/controller/static"
	"github.com/pajbot/pajbot2/pkg/web/controller/ws"
	"github.com/pajbot/pajbot2/pkg/web/router"
)

func LoadRoutes(a pkg.Application, cfg *config.Config) {
	channel.Load(a, cfg)

	dashboard.Load()
	home.Load()
	api.Load(a, cfg)
	static.Load()
	ws.Load()

	// /logout
	logout.Load()

	// /profile
	router.Get("/profile", handleProfile)

	// /banphrases
	banphrases.Load()

	// /admin
	admin.Load(a)

	// /commands
	router.Get("/commands", handleCommands)
}
