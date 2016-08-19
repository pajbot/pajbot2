package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

/*
The Config contains all the data required to connect
to the twitch IRC servers
*/
type Config struct {
	Pass        string  `json:"pass"`
	Nick        string  `json:"nick"`
	BrokerHost  *string `json:"broker_host"`
	BrokerPass  *string `json:"broker_pass"`
	BrokerLogin string  `json:"broker_login"`
	Silent      bool    `json:"silent"`

	RedisHost     string `json:"redis_host"`
	RedisPassword string `json:"redis_password"`
	RedisDatabase int    `json:"redis_database"`

	WebHost   string `json:"web_host"`
	WebDomain string `json:"web_domain"`

	SQLDSN string `json:"sql_dsn"`

	Auth struct {
		Twitch struct {
			Bot struct {
				// Twitch OAuth2 ID and Secret (created at twitch.tv/settings/connections)
				ClientID     string `json:"client_id"`
				ClientSecret string `json:"client_secret"`
				RedirectURI  string `json:"redirect_uri"`
			} `json:"bot"`
			User struct {
				ClientID     string `json:"client_id"`
				ClientSecret string `json:"client_secret"`
				RedirectURI  string `json:"redirect_uri"`
			} `json:"user"`
		} `json:"twitch"`
	} `json:"auth"`

	TLSKey  string `json:"tls_key"`
	TLSCert string `json:"tls_cert"`

	TwitterConsumerKey    string `json:"twitter_consumer_key"`
	TwitterConsumerSecret string `json:"twitter_consumer_secret"`
	TwitterAccessToken    string `json:"twitter_access_token"`
	TwitterAccessSecret   string `json:"twitter_access_secret"`

	Quit chan string

	ToWeb   chan map[string]interface{}
	FromWeb chan map[string]interface{}
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
	config := &Config{
		RedisDatabase: -1,
	}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	// Check for missing fields
	if config.BrokerHost == nil {
		return nil, errors.New("Missing field broker_host in config file")
	}
	if config.BrokerPass == nil {
		return nil, errors.New("Missing field broker_pass in config file")
	}

	return config, nil
}
