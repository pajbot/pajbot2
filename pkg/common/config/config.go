package config

import (
	"encoding/json"
	"io/ioutil"
)

type RedisConfig struct {
	Host     string
	Password string
	Database int
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

type AuthTwitchConfig struct {
	Bot      TwitchAuthConfig
	User     TwitchAuthConfig
	Streamer TwitchAuthConfig
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

type grpcServiceConfig struct {
	Host string
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
	Redis RedisConfig

	Web WebConfig

	SQL SQLConfig

	Auth authConfig

	Hooks map[string]struct {
		Secret string
	}

	TLSKey  string
	TLSCert string

	GRPCService grpcServiceConfig

	PubSub pubsubConfig

	Pajbot1 Pajbot1Config
}

var defaultConfig = Config{
	Redis: RedisConfig{
		Host: "localhost:6379",
	},
	Web: WebConfig{
		Host:   "localhost:2355",
		Domain: "localhost:2355",
	},
	GRPCService: grpcServiceConfig{
		Host: ":50052",
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
