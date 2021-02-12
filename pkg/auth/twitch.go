package auth

import (
	"errors"
	"fmt"

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

	authConfig = &cfg.Bot
	if err = validateAuthConfig("Bot", authConfig); err != nil {
		return nil, err
	}
	ta.twitchBotOauth = &oauth2.Config{
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s://%s/api/auth/twitch/bot/callback", protocol, webConfig.Domain),
		Endpoint:     twitch.Endpoint,
		Scopes: []string{
			"user:edit", // Edit bot account description/profile picture
			"channel:moderate",
			"chat:edit",
			"chat:read",
			"whispers:read",
			"whispers:edit",
		},
	}

	authConfig = &cfg.Streamer
	if err = validateAuthConfig("Streamer", authConfig); err != nil {
		return nil, err
	}
	ta.twitchStreamerOauth = &oauth2.Config{
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s://%s/api/auth/twitch/streeamer/callback", protocol, webConfig.Domain),
		Endpoint:     twitch.Endpoint,
		Scopes:       []string{
			// TODO: Figure out what scopes to ask for streamer authentications
		},
	}

	authConfig = &cfg.User
	if err = validateAuthConfig("User", authConfig); err != nil {
		return nil, err
	}
	ta.twitchUserOauth = &oauth2.Config{
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s://%s/api/auth/twitch/user/callback", protocol, webConfig.Domain),
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
