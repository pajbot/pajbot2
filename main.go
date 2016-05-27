package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/nuuls/pajbot2/irc"
)

/*
The Config contains all the data required to connect
to the twitch IRC servers
*/
type Config struct {
	Pass       string `json:"pass"`
	Nick       string `json:"nick"`
	BrokerPort string `json:"broker_port"`

	RedisHost     string `json:"redis_host"`
	RedisPassword string `json:"redis_password"`

	TLSKey  string `json:"tls_key"`
	TLSCert string `json:"tls_cert"`

	Channels []string `json:"channels"`

	ToWeb   chan map[string]interface{}
	FromWeb chan map[string]interface{}
}

func loadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	// TODO: Use config path from system arguments
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal("An error occured while loading the config file:", err)
	}

	irc := irc.Init(config.Pass, config.Nick)

	for _, channel := range config.Channels {
		irc.JoinChannel(channel)
	}

	for {
		time.Sleep(5 * time.Second)
	}
}
