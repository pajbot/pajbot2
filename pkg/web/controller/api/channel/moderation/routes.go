package moderation

import (
	"github.com/gorilla/mux"
	"github.com/pajbot/pajbot2/pkg/web/router"
)

// Load routes for /api/channel/:channel_id/moderation/
func Load(parent *mux.Router) {
	m := parent.PathPrefix("/moderation").Subrouter()

	router.RGet(m, `/latest`, apiChannelModerationLatest)

	router.RGet(m, `/user`, apiUser).Queries("user_id", `{user_id:[0-9]+}`)
	router.RGet(m, `/user`, apiUser).Queries("user_name", `{user_name:\w+}`)
	router.RGet(m, `/user`, apiUserMissingVariables)
	router.RGet(m, `/check_message`, apiCheckMessage).Queries("message", `{message:.+}`)
	router.RGet(m, `/check_message`, apiCheckMessageMissingVariables)
}
