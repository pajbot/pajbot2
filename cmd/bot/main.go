package main

import (
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"

	"github.com/pajlada/pajbot2/pkg/common"
	"github.com/pajlada/pajbot2/pkg/common/config"
)

var buildTime string

var version = flag.Bool("version", false, "Show pajbot2 version")
var configPath = flag.String("config", "./config.json", "")

var validURLs = []string{
	"imgur.com",        // Image host
	"twitter.com",      // Social media
	"twimg.com",        // Twitter image host
	"forsen.tv",        // Bot website
	"pajlada.se",       // Bot creator website
	"pajlada.com",      // Bot creator website
	"pajbot.com",       // Bot website
	"youtube.com",      // Video hosting website
	"youtu.be",         // Youtube short-url
	"prntscr.com",      // Image host
	"prnt.sc",          // prntscr short-url
	"steampowered.com", // Game shop
	"gyazo.com",        // Image host
	"www.com",          // Meme
}

func main() {
	common.BuildTime = buildTime

	flag.Usage = func() {
		helpCmd()
	}
	flag.Parse()
	command := flag.Arg(0)

	if *version {
		fmt.Println(*version)
		os.Exit(0)
	}

	switch command {
	case "check":
		_, err := config.LoadConfig(*configPath)
		if err != nil {
			fmt.Println("An error occured while loading the config file:", err)
			os.Exit(1)
		} else {
			fmt.Println("No errors found in the config file")
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
	_, err := os.Stderr.WriteString(
		`usage: pajbot2 <command> [<args>]
Commands:
   run            Run the bot (Default)
   check          Check the config file for missing fields
   install        Start the installation process (WIP)
   create <name>  Create a migration (WIP)
   newbot         Create a new bot
   linkchannel    Link a channel to a bot ID
`)
	if err != nil {
		log.Fatal(err)
	}
}

func runCmd() {
	application := newApplication()

	err := application.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("An error occured while loading the config file: ", err)
	}

	err = application.InitializeAPIs()
	if err != nil {
		log.Fatal("An error occured while initializing APIs: ", err)
	}

	err = application.InitializeSQL()
	if err != nil {
		log.Fatal("Error starting SQL client:", err)
	}

	err = application.RunDatabaseMigrations()
	if err != nil {
		log.Fatal("An error occured while running database migrations: ", err)
	}

	err = application.InitializeRedis()
	if err != nil {
		log.Fatal("Error starting redis client:", err)
	}

	err = application.InitializeModules()
	if err != nil {
		log.Fatal("Error initializing modules:", err)
	}

	err = application.LoadExternalEmotes()
	if err != nil {
		log.Fatal("An error occured while loading external emotes: ", err)
	}

	err = application.StartWebServer()
	if err != nil {
		log.Fatal("An error occured while starting the web server: ", err)
	}

	err = application.LoadBots()
	if err != nil {
		log.Fatal("An error occured while loading bots: ", err)
	}

	err = application.StartBots()
	if err != nil {
		log.Fatal("An error occured while starting bots: ", err)
	}

	err = application.StartPubSubClient()
	if err != nil {
		fmt.Println("Error starting PubSub Client:", err)
	}

	log.Fatal(application.Run())
}

func installCmd() {
	_, err := os.Stderr.WriteString(
		`"install" not yet implemented
`)
	if err != nil {
		log.Fatal(err)
	}
}

func createCmd() {
	_, err := os.Stderr.WriteString(
		`"create" not yet implemented
`)
	if err != nil {
		log.Fatal(err)
	}
}
