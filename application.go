package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dankeroni/gotwitch"
	twitch "github.com/gempir/go-twitch-irc"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/mysql"
	"github.com/pajlada/pajbot2/apirequest"
	"github.com/pajlada/pajbot2/common/config"
	"github.com/pajlada/pajbot2/redismanager"
	"github.com/pajlada/pajbot2/sqlmanager"
	"github.com/pajlada/pajbot2/web"
)

const migrationsPath = "file://migrations"

// Application is the heart of pajbot
// It keeps the functions to initialize, start, and stop pajbot
type Application struct {
	config     *config.Config
	TwitchBots map[string]*twitch.Client
	Redis      *redismanager.RedisManager
	SQL        *sqlmanager.SQLManager
}

// NewApplication creates an instance of Application. Generally this should only be done once
func NewApplication() *Application {
	ret := Application{}

	ret.TwitchBots = make(map[string]*twitch.Client)

	return &ret
}

// LoadConfig loads a config file from the given path. The format for the config file is available in the config package
func (a *Application) LoadConfig(path string) error {
	config, err := config.LoadConfig(path)
	if err != nil {
		return err
	}

	a.config = config

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		a.config.Quit <- "Quitting due to SIGTERM/SIGINT"
	}()
	a.config.Quit = make(chan string)

	return nil
}

// RunDatabaseMigrations runs database migrations on the database specified in the config file
func (a *Application) RunDatabaseMigrations() error {
	db, err := sql.Open("mysql", a.config.SQLDSN)
	if err != nil {
		return err
	}

	defer func() {
		dErr := db.Close()
		if dErr != nil {
			log.Println("Error in deferred close:", dErr)
		}
	}()

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "mysql", driver)
	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}

		return err
	}

	return nil
}

// InitializeAPIs initializeds various APIs that are needed for pajbot
func (a *Application) InitializeAPIs() error {
	// Twitch APIs
	apirequest.Twitch = gotwitch.New(a.config.Auth.Twitch.User.ClientID)
	apirequest.TwitchBot = gotwitch.New(a.config.Auth.Twitch.Bot.ClientID)
	apirequest.TwitchV3 = gotwitch.NewV3(a.config.Auth.Twitch.User.ClientID)
	apirequest.TwitchBotV3 = gotwitch.NewV3(a.config.Auth.Twitch.Bot.ClientID)

	return nil
}

func (a *Application) StartWebServer() error {
	a.Redis = redismanager.Init(a.config)
	a.SQL = sqlmanager.Init(a.config)

	webCfg := &web.Config{
		Bots:  a.TwitchBots,
		Redis: a.Redis,
		SQL:   a.SQL,
	}

	webBoss := web.Init(a.config, webCfg)
	go webBoss.Run()

	return nil
}

// LoadBots loads bots from the database
func (a *Application) LoadBots() error {
	db, err := sql.Open("mysql", a.config.SQLDSN)
	if err != nil {
		return err
	}

	defer func() {
		dErr := db.Close()
		if dErr != nil {
			log.Println("Error in deferred close:", dErr)
		}
	}()

	rows, err := db.Query("SELECT `name`, `twitch_access_token` FROM `pb_bot`")
	if err != nil {
		return err
	}

	for rows.Next() {
		var name string
		var twitchAccessToken string
		if err := rows.Scan(&name, &twitchAccessToken); err != nil {
			log.Fatal("Error scanning values: ", err)
		}

		log.Println("Got bot", name)
		a.TwitchBots[name] = twitch.NewClient(name, "oauth:"+twitchAccessToken)
	}

	return nil
}

// StartBots starts bots that were loaded from the LoadBots method
func (a *Application) StartBots() error {
	for _, bot := range a.TwitchBots {
		bot.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
			log.Printf("%s(%s): %s", user.DisplayName, user.Username, message.Text)
			if message.Text == "!xd" && user.Username == "pajlada" {
				bot.Say(channel, "XDDDDDDDDDD")
			}

			if message.Text == "!pb2quit" && user.Username == "pajlada" {
				a.config.Quit <- "Quit because pajlada said so"
			}
		})

		bot.Join("pajlada")

		go func(bot *twitch.Client) {
			log.Println("Connecting...")
			err := bot.Connect()
			if err != nil {
				log.Fatal(err)
			}
		}(bot)
	}

	return nil
}

// Run blocks the current thread, waiting for something to put an exit string into the Quit channel
func (a *Application) Run() error {
	/*
		b := boss.Init(config)
		go bot.LoadGlobalEmotes()
		for _, ircConnection := range b.IRCConnections {
			bots = append(bots, ircConnection.Bots)
		}
		log.Fatal(q)
	*/

	quitString := <-a.config.Quit

	return fmt.Errorf(quitString)
}
