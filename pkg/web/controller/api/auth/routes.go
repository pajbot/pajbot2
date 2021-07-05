package auth

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/common/config"
	"github.com/pajbot/pajbot2/pkg/web/controller/api/auth/twitch"
)

// Load routes for /api/auth/
func Load(parent *mux.Router, a pkg.Application, cfg *config.Config) {
	m := parent.PathPrefix("/auth").Subrouter()

	// Subroute /api/auth/twitch/
	err := twitch.Load(m, a)
	if err != nil {
		fmt.Println("Error loading /api/auth/twitch:", err)
	}
}
