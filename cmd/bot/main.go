package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"

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

	case "fix":
		fixCmd()

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
   fix <number>   Fix issue #NUMBER automatically (or attempt to)
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

func fixCmd() {
	number, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}

	fix(number)
}

func fix(number int) {
	switch number {
	case 1:
		// This will poll the API for the user ID for each bot and fill in the Bot table
		fmt.Println("Fix issue #1: Migration failed in line 0: CREATE UNIQUE INDEX `itwitchuid` ON Bot(twitch_userid); (details: Error 1062: Duplicate entry '' for key 'itwitchuid')")
		fmt.Print("Is this the issue you want to fix? (y/n) ")
		buf := bufio.NewReader(os.Stdin)
		sentence, err := buf.ReadBytes('\n')
		if err != nil {
			log.Fatal(err)
		}
		answer := strings.ToLower(strings.TrimSpace(string(sentence)))

		switch answer {
		case "y":
			// Ensure schema migration is proper
			fmt.Println("Attempting to do fix")

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

			db := application.SQL()
			const queryF = `SELECT version, dirty FROM schema_migrations`
			row := db.QueryRow(queryF)
			var version int64
			var dirty bool
			err = row.Scan(&version, &dirty)
			if err != nil {
				log.Fatal(err)
			}

			if version != 20190118204509 || !dirty {
				log.Fatal("wrong schema versions, this is probably not the fix you need")
			}

			rows, err := db.Query("SELECT id, twitch_userid, name FROM Bot")
			if err != nil {
				log.Fatal(err)
			}
			for rows.Next() {
				var id int64
				var oldUserID string
				var name string
				err = rows.Scan(&id, &oldUserID, &name)
				if err != nil {
					log.Fatal(err)
				}

				if oldUserID != "" {
					fmt.Printf("Skipping %s because it already has an ID set (%s)\n", name, oldUserID)
					continue
				}

				twitchUserID := application.UserStore().GetID(name)
				if twitchUserID == "" {
					fmt.Printf("Unable to get ID for user '%s'\n", name)
					continue
				}

				_, err = db.Exec("UPDATE Bot SET twitch_userid=? WHERE id=?", twitchUserID, id)
				if err != nil {
					log.Println("Error updating user ID:", err)
					continue
				}

				fmt.Println("Updated User ID for", name, "to", twitchUserID)
			}

			fmt.Println("Attempting to revert the schema migration (maybe use migration .Down? here or something)")

			_, err = db.Exec("UPDATE schema_migrations SET version=20190118000448, dirty=0")
			if err != nil {
				log.Fatal("Error reverting schema migration", err)
			}

			fmt.Println("Done! now just run ./bot like normal")

		default:
			fallthrough
		case "n":
			os.Exit(0)
		}
	}
}

func runCmd() {
	application := newApplication()

	err := application.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("An error occured while loading the config file: ", err)
	}

	err = application.InitializeOAuth2Configs()
	if err != nil {
		log.Fatal("An error occured while initializing oauth2 config: ", err)
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

	err = application.ProvideAdminPermissionsToAdmin()
	if err != nil {
		log.Fatal("Error providing admin access to admin:", err)
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
