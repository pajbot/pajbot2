package api

import (
	"fmt"
	"net/http"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/common/config"
	"github.com/pajbot/pajbot2/pkg/web/controller/api/auth"
	"github.com/pajbot/pajbot2/pkg/web/controller/api/channel"
	"github.com/pajbot/pajbot2/pkg/web/controller/api/report"
	"github.com/pajbot/pajbot2/pkg/web/controller/api/webhook"
	"github.com/pajbot/pajbot2/pkg/web/router"
)

func apiRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "XD API ROOT")
	fmt.Fprintf(w, "\nTODO: Link to api docs")
}

func Load(a pkg.Application, cfg *config.Config) {
	m := router.Subrouter("/api")

	router.RGet(m, "/", apiRoot)
	m.HandleFunc("", apiRoot)

	auth.Load(m, a, cfg)

	channel.Load(m)

	report.Load(m)

	webhook.Load(m, cfg)
}
