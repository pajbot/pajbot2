package auth

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/common/config"
	"github.com/pajlada/pajbot2/pkg/web/controller/api/auth/twitch"
	"github.com/pajlada/pajbot2/pkg/web/router"
)

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Auth root")
}

func Load(parent *mux.Router, a pkg.Application, cfg *config.Config) {
	m := parent.PathPrefix("/auth").Subrouter()

	router.RGet(m, "", root)

	err := twitch.Load(m, a)
	if err != nil {
		fmt.Println("Error loading /api/auth/twitch:", err)
	}
}
