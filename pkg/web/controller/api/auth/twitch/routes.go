package twitch

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/pajbot/pajbot2/pkg"
)

// Load routes for /api/auth/twitch/
func Load(parent *mux.Router, a pkg.Application) error {
	m := parent.PathPrefix("/twitch").Subrouter()
	auths := a.TwitchAuths()
	ctx := context.Background()

	// Initialize /api/auth/twitch/bot and /api/auth/twitch/bot/callback routes
	initializeOauthRoutes(ctx, m, auths.Bot(), "bot", onBotAuthenticated)

	// Initialize /api/auth/twitch/user and /api/auth/twitch/user/callback routes
	initializeOauthRoutes(ctx, m, auths.User(), "user", onUserAuthenticated)

	// Initialize /api/auth/twitch/streamer and /api/auth/twitch/streamer/callback routes
	initializeOauthRoutes(ctx, m, auths.Streamer(), "streamer", onStreamerAuthenticated)

	return nil
}
