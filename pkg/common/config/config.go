package config

import (
	"encoding/json"
	"io/ioutil"
)

type AdminConfig struct {
	TwitchUserID string
}

type WebConfig struct {
	Host   string
	Domain string
	Secure bool
}

type SQLConfig struct {
	DSN string
}

type TwitchAuthConfig struct {
	// Twitch OAuth2 ID and Secret (created at twitch.tv/settings/connections)
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type TwitchWebhookConfig struct {
	HostPrefix       string
	Secret           string
	LeaseTimeSeconds int
}

type AuthTwitchConfig struct {
	// Bot contains the client id, secret, and redirect URI for authenticating new bot accounts
	Bot TwitchAuthConfig

	// User contains the client id, secret, and redirect URI for authenticating random users/moderators who want to log into the dashboard
	User TwitchAuthConfig

	// Streamer contains the client id, secret, and reidrect URI for authenticating streamers
	// This will be an extra option after logging in where streamers can choose to give out more permissions (like getting their subscribers)
	Streamer TwitchAuthConfig

	Webhook TwitchWebhookConfig
}

type authTwitterConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

type authConfig struct {
	Twitch AuthTwitchConfig

	Twitter authTwitterConfig
}

type pubsubConfig struct {
	ChannelID string
	UserID    string
	UserToken string
}

type Pajbot1Config struct {
	SQL SQLConfig
}

/*
The Config contains all the data required to connect
to the twitch IRC servers
*/
type Config struct {
	Admin AdminConfig

	Web WebConfig

	SQL SQLConfig

	Auth authConfig

	Hooks map[string]struct {
		Secret string
	}

	TLSKey  string
	TLSCert string

	PubSub pubsubConfig

	Pajbot1 Pajbot1Config
}

var defaultConfig = Config{
	Web: WebConfig{
		Host:   "localhost:2355",
		Domain: "localhost:2355",
	},
	Auth: authConfig{
		Twitch: AuthTwitchConfig{
			Webhook: TwitchWebhookConfig{
				LeaseTimeSeconds: 24 * 3600,
			},
		},
	},
}

/*
LoadConfig parses a config file from the given json file at the path
and returns a Config object
*/
func LoadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := defaultConfig
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
