package channel

import (
	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/web/controller/api/channel/moderation"
)

func Load(parent *mux.Router) {
	m := parent.PathPrefix(`/channel/{channelID:\w+}`).Subrouter()

	moderation.Load(m)

	// m.HandleFunc(`/channel/{channel:\w+}/{rest:.*}`, APIHandler)
}
