package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"

	_ "github.com/mattes/migrate/driver/mysql"
	"github.com/mattes/migrate/migrate"
	"github.com/pajlada/pajbot2/boss"
	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/common/config"
	"github.com/pajlada/pajbot2/helper"
	"github.com/pajlada/pajbot2/plog"
	"github.com/pajlada/pajbot2/sqlmanager"
	"github.com/pajlada/pajbot2/web"
)

var log = plog.GetLogger()

func cleanup() {
	// TODO: Perform cleanups
}

var buildTime string

var version = flag.Bool("version", false, "Show pajbot2 version")
var configPath = flag.String("config", "./config.json", "")

func main() {
	plog.InitLogging()

	common.BuildTime = buildTime

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
		_, err := config.LoadConfig(*configPath)
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

	case "newbot":
		newbotCmd()

	case "linkchannel":
		linkchannelCmd()

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
   newbot         Create a new bot
   linkchannel    Link a channel to a bot ID
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
	config, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("An error occured while loading the config file:", err)
	}

	// Run database migrations
	allErrors, ok := migrate.UpSync("mysql://"+config.SQLDSN, "./migrations")
	if !ok {
		log.Debug("An error occured while trying to run database migrations")
		for _, err := range allErrors {
			log.Debug(err)
		}
		os.Exit(1)
	}

	// Start web server

	go func() {
		log.Error(http.ListenAndServe(":11223", nil))
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		config.Quit <- "Quitting due to SIGTERM/SIGINT"
	}()
	config.Quit = make(chan string)
	b := boss.Init(config)
	var bots []map[string]*bot.Bot
	for _, ircConnection := range b.IRCConnections {
		bots = append(bots, ircConnection.Bots)
	}
	webCfg := &web.Config{
		Bots:  bots,
		Redis: b.Redis,
		SQL:   b.SQL,
	}
	webBoss := web.Init(config, webCfg)
	go webBoss.Run()
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

// add a new bot to pb_bot
func newbotCmd() {
	config, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("An error occured while loading the config file:", err)
	}

	sql := sqlmanager.Init(config)

	reader := bufio.NewReader(os.Stdin)

	var name string
	var accessToken string
	var refreshToken string

	fmt.Println("Enter proper values for the incoming questions to create a new bot in the pb_bot table")

	fmt.Print("Bot name: ")
	name = helper.ReadArg(reader)
	fmt.Print("Bot access token: ")
	accessToken = helper.ReadArg(reader)
	fmt.Print("Bot refresh token: ")
	refreshToken = helper.ReadArg(reader)

	fmt.Println("Creating a new bot with the given credentials")

	common.CreateBotAccount(sql.Session, name, accessToken, refreshToken)
}

// Link a pb_channel to a pb_bot
func linkchannelCmd() {
	config, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("An error occured while loading the config file:", err)
	}

	sql := sqlmanager.Init(config)

	reader := bufio.NewReader(os.Stdin)

	var name string
	var channelName string

	fmt.Println("Enter proper values for the incoming questions to create a new bot in the pb_bot table")

	fmt.Print("Bot name: ")
	name = helper.ReadArg(reader)

	b, err := common.GetBotAccount(sql.Session, name)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Print("Channel name: ")
	channelName = helper.ReadArg(reader)

	c, err := common.GetChannel(sql.Session, channelName)
	if err != nil {
		fmt.Println("No channel with the name " + channelName)
		return
	}

	c.SQLSetBotID(sql, b.ID)

	fmt.Printf("Linked channel %s to bot %s\n", channelName, name)
}
