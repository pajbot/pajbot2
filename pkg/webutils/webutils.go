package webutils

import (
	"net/http"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/users"
	"github.com/pajlada/pajbot2/pkg/utils"
	"github.com/pajlada/pajbot2/pkg/web/state"
)

func RequirePermission(w http.ResponseWriter, c state.State, permission pkg.Permission) bool {
	if c.Session == nil {
		utils.WebWriteError(w, 400, "Not authorized to view this endpoint")
		return false
	}

	user := users.NewSimpleTwitchUser(c.Session.TwitchUserID, c.Session.TwitchUserName)
	if user == nil {
		utils.WebWriteError(w, 400, "Not authorized to view this endpoint")
		return false
	}

	if !user.HasGlobalPermission(pkg.PermissionModeration) {
		utils.WebWriteError(w, 400, "Not authorized to view this endpoint!!!")
		return false
	}

	return true
}
