package admin

import (
	"log"
	"net/http"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/users"
	"github.com/pajlada/pajbot2/pkg/web/router"
	"github.com/pajlada/pajbot2/pkg/web/state"
	"github.com/pajlada/pajbot2/pkg/web/views"
)

func Load() {
	router.Get("/admin", handleRoot)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
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

	// TODO: Also fetch all channel admin permissions, and send that along to.
	// That would tell the API what channels it can administrate

	err := views.Render("admin", w, r)
	if err != nil {
		log.Println("Error rendering admin view:", err)
	}
}
