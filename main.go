package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

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

func cleanup() {
	// TODO: Perform cleanups
}

func main() {
	// TODO: Use config path from system arguments
	config, err := LoadConfig("config.json")

	if err != nil {
		log.Fatal("An error occured while loading the config file:", err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		config.Quit <- "Quitting due to SIGTERM/SIGINT"
	}()
	config.Quit = make(chan string)
	go boss.Init(config)
	q := <-config.Quit
	cleanup()
	log.Fatal(q)
}
