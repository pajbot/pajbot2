package report

import (
	"github.com/gorilla/mux"
	"github.com/pajbot/pajbot2/pkg/web/router"
)

// Load routes for /api/report/
func Load(parent *mux.Router) {
	m := parent.PathPrefix("/report").Subrouter()

	router.RGet(m, `/history`, apiHistory)
}
