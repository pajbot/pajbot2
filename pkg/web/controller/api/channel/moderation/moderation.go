package moderation

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/pkg/web/router"
)

func root(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Channel moderation root for channel ID '%s'", vars["channelID"])
}

func Load(parent *mux.Router) {
	m := parent.PathPrefix("/moderation").Subrouter()

	router.RGet(m, "", root)

	router.RGet(m, `/latest`, apiChannelModerationLatest)

	router.RGet(m, `/user`, apiUser).Queries("user_id", `{user_id:[0-9]+}`)
	router.RGet(m, `/user`, apiUser).Queries("user_name", `{user_name:\w+}`)
	router.RGet(m, `/user`, apiUserMissingVariables)
	router.RGet(m, `/check_message`, apiCheckMessage).Queries("message", `{message:.+}`)
	router.RGet(m, `/check_message`, apiCheckMessageMissingVariables)
}
