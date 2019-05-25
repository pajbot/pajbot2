package banphrases

import (
	"github.com/gorilla/mux"
	"github.com/pajbot/pajbot2/pkg/web/router"
)

func Load(parent *mux.Router) {
	m := parent.PathPrefix("/banphrases").Subrouter()

	router.RGet(m, `/list`, handleList)
}
