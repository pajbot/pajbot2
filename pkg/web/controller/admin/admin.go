package admin

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/users"
	"github.com/pajlada/pajbot2/pkg/web/router"
	"github.com/pajlada/pajbot2/pkg/web/state"
	"github.com/pajlada/pajbot2/pkg/web/views"
)

func Load(a pkg.Application) {
	router.Get("/admin", handleRoot(a))
}

func handleRoot(a pkg.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Ensure user is logged in
		c := state.Context(w, r)
		if c.Session == nil {
			views.Render403(w, r)
			return
		}

		// Step 2: Ensure logged in user has global admin permissions
		user := users.NewSimpleTwitchUser(c.Session.TwitchUserID, c.Session.TwitchUserName)

		if !user.HasGlobalPermission(pkg.PermissionAdmin) {
			views.Render403(w, r)
			return
		}

		type BotInfo struct {
			Name string
		}

		type Extra struct {
			Bots []BotInfo
		}

		extra := &Extra{}

		for it := a.TwitchBots().Iterate(); it.Next(); {
			bot := it.Value()
			botInfo := BotInfo{
				Name: bot.TwitchAccount().Name(),
			}
			extra.Bots = append(extra.Bots, botInfo)
		}

		// TODO: Also fetch all channel admin permissions, and send that along to.
		// That would tell the API what channels it can administrate

		extraBytes, _ := json.Marshal(extra)

		err := views.RenderExtra("admin", w, r, extraBytes)
		if err != nil {
			log.Println("Error rendering admin view:", err)
		}
	}
}
