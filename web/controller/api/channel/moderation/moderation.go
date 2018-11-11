package moderation

import (
	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/web/router"
)

func Load(parent *mux.Router) {
	m := parent.PathPrefix("/moderation").Subrouter()
	router.RGet(m, `/latest`, apiChannelModerationLatest)
}
