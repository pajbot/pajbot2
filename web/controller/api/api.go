package api

import (
	"fmt"
	"net/http"

	"github.com/pajlada/pajbot2/pkg/common/config"
	"github.com/pajlada/pajbot2/web/controller/api/auth"
	"github.com/pajlada/pajbot2/web/controller/api/channel"
	"github.com/pajlada/pajbot2/web/router"
)

func apiRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "XD API ROOT")
}

func Load(cfg *config.Config) {
	m := router.Subrouter("/api")

	router.RGet(m, "", apiRoot)
	// m.HandleFunc("", apiRoot)
	// router.Get("/api", apiRoot)

	auth.Load(m, cfg)

	channel.Load(m)
}
