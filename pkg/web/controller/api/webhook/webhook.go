package webhook

import (
	"github.com/gorilla/mux"
	"github.com/pajbot/pajbot2/pkg/common/config"
	"github.com/pajbot/pajbot2/pkg/web/router"
)

func Load(parent *mux.Router, cfg *config.Config) {
	m := parent.PathPrefix("/webhook").Subrouter()

	// NEW AND FRESH AND COOL
	router.RPost(m, `/eventsub`, apiEventsub(&cfg.Auth.Twitch.Webhook))

	router.RPost(m, `/github/{channelID:\w+}`, apiGithub(&cfg.Auth.Github.Webhook))
}
