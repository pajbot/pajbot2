package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"

	"github.com/pajlada/pajbot2/pkg/common"
	"github.com/pajlada/pajbot2/pkg/common/config"
	"github.com/pajlada/pajbot2/pkg/utils"
	"github.com/pajlada/pajbot2/sqlmanager"
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
			log.Println("An error occured while loading the config file:", err)
			os.Exit(1)
		} else {
			log.Println("No errors found in the config file")
			os.Exit(0)
		}

	case "install":
		installCmd()

	case "create":
		createCmd()

	case "newbot":
		newbotCmd()

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
	application := NewApplication()

	err := application.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("An error occured while loading the config file: ", err)
	}

	err = application.RunDatabaseMigrations()
	if err != nil {
		log.Fatal("An error occured while running database migrations: ", err)
	}

	err = application.StartRedisClient()
	if err != nil {
		log.Fatal("Error starting redis client:", err)
	}

	err = application.InitializeAPIs()
	if err != nil {
		log.Fatal("An error occured while initializing APIs: ", err)
	}

	err = application.LoadExternalEmotes()
	if err != nil {
		log.Fatal("An error occured while loading external emotes: ", err)
	}

	err = application.StartWebServer()
	if err != nil {
		log.Fatal("An error occured while starting the web server: ", err)
	}

	/*
		err = application.LoadOldPajbot()
		if err != nil {
			log.Fatal("An error occured while loading old pajbot: ", err)
		}
	*/

	err = application.LoadBots()
	if err != nil {
		log.Fatal("An error occured while loading bots: ", err)
	}

	/*
		err = application.StartContextBot()
		if err != nil {
			log.Fatal("An error occured while starting context bot: ", err)
		}
	*/

	err = application.StartBots()
	if err != nil {
		log.Fatal("An error occured while starting bots: ", err)
	}

	err = application.StartPubSubClient()
	if err != nil {
		log.Fatal("An error occured while starting pubsub client", err)
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

// add a new bot to pb_bot
func newbotCmd() {
	config, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("An error occured while loading the config file:", err)
	}

	sql := sqlmanager.Init(config.SQL)

	reader := bufio.NewReader(os.Stdin)

	var name string
	var accessToken string
	var refreshToken string

	fmt.Println("Enter proper values for the incoming questions to create a new bot in the pb_bot table")

	fmt.Print("Bot name: ")
	name = utils.ReadArg(reader)
	fmt.Print("Bot access token: ")
	accessToken = utils.ReadArg(reader)
	fmt.Print("Bot refresh token: ")
	refreshToken = utils.ReadArg(reader)

	fmt.Println("Creating a new bot with the given credentials")

	err = common.CreateBot(sql.Session, name, accessToken, refreshToken)
	if err != nil {
		log.Fatal(err)
	}
}
