package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"

	_ "github.com/mattes/migrate/driver/mysql"
	"github.com/mattes/migrate/migrate"
	"github.com/pajlada/pajbot2/boss"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/plog"
	pb_websocket "github.com/pajlada/pajbot2/websocket"
)

var log = plog.GetLogger()

/*
LoadConfig parses a config file from the given json file at the path
and returns a Config object
*/
func LoadConfig(path string) (*common.Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := &common.Config{
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

func cleanup() {
	// TODO: Perform cleanups
}

var version = flag.Bool("version", false, "Show pajbot2 version")
var configPath = flag.String("config", "./config.json", "")

func main() {
	plog.InitLogging()

	flag.Usage = func() {
		helpCmd()
	}
	flag.Parse()
	command := flag.Arg(0)

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	switch command {
	case "check":
		_, err := LoadConfig(*configPath)
		if err != nil {
			log.Error("An error occured while loading the config file:", err)
			os.Exit(1)
		} else {
			log.Debug("No errors found in the config file")
			os.Exit(0)
		}

	case "install":
		installCmd()

	case "create":
		createCmd()

	case "help":
		helpCmd()

	default:
		fallthrough
	case "run":
		runCmd()
	}
}

func helpCmd() {
	os.Stderr.WriteString(
		`usage: pajbot2 <command> [<args>]
Commands:
   run            Run the bot (Default)
   check          Check the config file for missing fields
   install        Start the installation process (WIP)
   create <name>  Create a migration (WIP)
`)
}

type msg struct {
	Num int
}

func wsHandler(conn *websocket.Conn) {
	for {
		m := msg{}

		err := conn.ReadJSON(&m)
		if err != nil {
			log.Error("Error reading json.", err)
		}

		log.Infof("Got message: %#v\n", m)

		if err = conn.WriteJSON(m); err != nil {
			log.Error(err)
		}
	}
}

func runCmd() {
	// TODO: Use config path from system arguments
	config, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatal("An error occured while loading the config file:", err)
	}

	// Run database migrations
	log.Debug("Running database migrations")
	allErrors, ok := migrate.UpSync("mysql://"+config.SQLDSN, "./migrations")
	if !ok {
		log.Debug("An error occured while trying to run database migrations")
		for _, err := range allErrors {
			log.Debug(err)
		}
		os.Exit(1)
	}
	log.Debug("Done")

	// Start websocket server
	wsHost := ":2355"
	log.Debugf("Starting websocket server at %s\n", wsHost)
	wsBoss := pb_websocket.Init(wsHost)
	wsBoss.Handler = wsHandler
	go wsBoss.Run()

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

func installCmd() {
	os.Stderr.WriteString(
		`"install" not yet implemented
`)
}

func createCmd() {
	os.Stderr.WriteString(
		`"create" not yet implemented
`)
}
