package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/pajlada/pajbot2/boss"
	"github.com/pajlada/pajbot2/common"
)

/*
LoadConfig parses a config file from the given json file at the path
and returns a Config object
*/
func LoadConfig(path string) (*common.Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := &common.Config{}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	// TODO: Use config path from system arguments
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatal("An error occured while loading the config file:", err)
	}

	boss.Init(config)

	quit := make(chan interface{})
	<-quit
}
