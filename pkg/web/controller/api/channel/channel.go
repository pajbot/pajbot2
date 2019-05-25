package channel

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pajbot/pajbot2/pkg/web/controller/api/channel/banphrases"
	"github.com/pajbot/pajbot2/pkg/web/controller/api/channel/moderation"
	"github.com/pajbot/pajbot2/pkg/web/router"
)

func root(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Channel root for channel ID '%s'", vars["channel_id"])
}

func Load(parent *mux.Router) {
	m := parent.PathPrefix(`/channel/{channel_id:\w+}`).Subrouter()

	router.RGet(m, "", root)

	moderation.Load(m)
	banphrases.Load(m)

	// m.HandleFunc(`/channel/{channel:\w+}/{rest:.*}`, APIHandler)
}
