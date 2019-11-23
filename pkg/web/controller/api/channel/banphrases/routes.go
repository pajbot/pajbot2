package banphrases

import (
	"github.com/gorilla/mux"
	"github.com/pajbot/pajbot2/pkg/web/router"
)

// Load routes for /api/channel/:channel_id/banphrases/
func Load(parent *mux.Router) {
	m := parent.PathPrefix("/banphrases").Subrouter()

	router.RGet(m, `/list`, handleList)
}
