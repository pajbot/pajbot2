package webutils

import (
	"net/http"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/users"
	"github.com/pajbot/utils"
	"github.com/pajbot/pajbot2/pkg/web/state"
)

func RequirePermission(w http.ResponseWriter, c state.State, channel pkg.Channel, permission pkg.Permission) bool {
	if c.Session == nil {
		utils.WebWriteError(w, 400, "Not authorized to view this endpoint")
		return false
	}

	user := users.NewSimpleTwitchUser(c.Session.TwitchUserID, c.Session.TwitchUserName)
	if user == nil {
		utils.WebWriteError(w, 400, "Not authorized to view this endpoint")
		return false
	}

	if channel != nil {
		if user.HasPermission(channel, permission) {
			return true
		}
	} else {
		if user.HasGlobalPermission(permission) {
			return true
		}
	}

	utils.WebWriteError(w, 400, "Not authorized to view this endpoint!!!")
	return false
}
