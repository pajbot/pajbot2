package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattes/migrate/driver/mysql"
	"github.com/mattes/migrate/migrate"
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
	config := &common.Config{
		RedisDatabase: -1,
	}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func cleanup() {
	// TODO: Perform cleanups
}

var version = flag.Bool("version", false, "Show pajbot2 version")
var configPath = flag.String("config", "./config.json", "")

func main() {
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
   install        Start the installation process (WIP)
   create <name>  Create a migration (WIP)
`)
}

func runCmd() {
	// TODO: Use config path from system arguments
	config, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatal("An error occured while loading the config file:", err)
	}

	// Run database migrations
	log.Println("Running database migrations")
	allErrors, ok := migrate.UpSync("mysql://"+config.SQLDSN, "./migrations")
	if !ok {
		log.Println("An error occured while trying to run database migrations")
		for _, err := range allErrors {
			log.Println(err)
		}
		os.Exit(1)
	}
	log.Println("Done")

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
