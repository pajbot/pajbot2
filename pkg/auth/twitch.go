package auth

import (
	"errors"
	"net/url"

	"github.com/pajbot/pajbot2/pkg/common/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

type TwitchAuths struct {
	twitchBotOauth      *oauth2.Config
	twitchUserOauth     *oauth2.Config
	twitchStreamerOauth *oauth2.Config
}

func NewTwitchAuths(cfg *config.AuthTwitchConfig, webConfig *config.WebConfig) (*TwitchAuths, error) {
	var authConfig *config.TwitchAuthConfig
	var protocol string
	var err error
	ta := &TwitchAuths{}

	if webConfig.Secure {
		protocol = "https"
	} else {
		protocol = "http"
	}

	u := url.URL{
		Scheme: protocol,
		Host:   webConfig.Domain,
	}

	authConfig = &cfg.Bot
	if err = validateAuthConfig("Bot", authConfig); err != nil {
		return nil, err
	}
	u.Path = "/api/auth/twitch/bot/callback"
	ta.twitchBotOauth = &oauth2.Config{
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
		RedirectURL:  u.String(),
		Endpoint:     twitch.Endpoint,
		Scopes: []string{
			"user:edit", // Edit bot account description/profile picture
			"channel:moderate",
			"chat:edit",
			"chat:read",
			"whispers:read",
			"whispers:edit",
			"moderator:manage:announcements", // For sending Announcements
			"user:manage:whispers",           // For sending whispers
			"moderator:manage:banned_users",  // For banning, timing out, and unbanning users
			"moderator:manage:chat_messages", // For deleting messages
			"moderator:manage:chat_settings", // For changing chat settings (e.g. followers mode, unique mode, slow mode)
			"moderator:read:chatters",        // For getting chatters in the channel
		},
	}

	authConfig = &cfg.Streamer
	if err = validateAuthConfig("Streamer", authConfig); err != nil {
		return nil, err
	}
	u.Path = "/api/auth/twitch/streamer/callback"
	ta.twitchStreamerOauth = &oauth2.Config{
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
		RedirectURL:  u.String(),
		Endpoint:     twitch.Endpoint,
		Scopes: []string{
			"channel:read:vips", // For polling the list of VIPs in the channel
			"moderation:read",   // For polling the list of moderators in the channel
		},
	}

	authConfig = &cfg.User
	if err = validateAuthConfig("User", authConfig); err != nil {
		return nil, err
	}
	u.Path = "/api/auth/twitch/user/callback"
	ta.twitchUserOauth = &oauth2.Config{
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
		RedirectURL:  u.String(),
		Endpoint:     twitch.Endpoint,
		Scopes:       []string{},
	}

	return ta, nil
}

func validateAuthConfig(name string, authConfig *config.TwitchAuthConfig) error {
	if authConfig.ClientID == "" {
		return errors.New("Missing required Client ID in " + name + " auth in your config.json file")
	}
	if authConfig.ClientSecret == "" {
		return errors.New("Missing required Client Secret in " + name + " auth in your config.json file")
	}

	return nil
}

func (a *TwitchAuths) Bot() *oauth2.Config {
	return a.twitchBotOauth
}

func (a *TwitchAuths) Streamer() *oauth2.Config {
	return a.twitchStreamerOauth
}

func (a *TwitchAuths) User() *oauth2.Config {
	return a.twitchUserOauth
}
