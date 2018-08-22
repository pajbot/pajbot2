package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"encoding/json"

	"errors"
	"strconv"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"

	_ "github.com/go-sql-driver/mysql" // MySQL Driver

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/dankeroni/gotwitch"
	"github.com/gempir/go-twitch-irc"
	"github.com/pajlada/go-twitch-pubsub"
	"github.com/pajlada/pajbot2/emotes"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/apirequest"
	"github.com/pajlada/pajbot2/pkg/commands"
	"github.com/pajlada/pajbot2/pkg/common/config"
	"github.com/pajlada/pajbot2/pkg/modules"
	pb2twitch "github.com/pajlada/pajbot2/pkg/twitch"
	"github.com/pajlada/pajbot2/pkg/users"
	"github.com/pajlada/pajbot2/web"
)

type channelContext struct {
	// TODO: replace []string with some 5 message long fifo queue
	Channels map[string][]string
}

func NewChannelContext() *channelContext {
	return &channelContext{
		Channels: make(map[string][]string),
	}
}

// Application is the heart of pajbot
// It keeps the functions to initialize, start, and stop pajbot
type Application struct {
	config *config.Config

	TwitchBots   map[string]*pb2twitch.Bot
	Redis        *redis.Pool
	SQL          *sql.DB
	TwitchPubSub *twitch_pubsub.Client

	// key = user ID
	UserContext map[string]*channelContext

	Quit chan string
}

func lol(xd string) *string {
	return &xd
}

func (a *Application) GetUserMessages(channelID, userID string) ([]string, error) {
	if uc, ok := a.UserContext[userID]; ok {
		if cc, ok := uc.Channels[channelID]; ok {
			return cc, nil
		}

		return nil, errors.New("No messages found in this channel for this user")
	}

	return nil, errors.New("No messages found for this user")
}

// NewApplication creates an instance of Application. Generally this should only be done once
func NewApplication() *Application {
	a := Application{}

	a.TwitchBots = make(map[string]*pb2twitch.Bot)
	a.Quit = make(chan string)
	a.UserContext = make(map[string]*channelContext)

	return &a
}

// LoadConfig loads a config file from the given path. The format for the config file is available in the config package
func (a *Application) LoadConfig(path string) error {
	cfg, err := config.LoadConfig(path)
	if err != nil {
		return err
	}

	a.config = cfg

	return nil
}

// RunDatabaseMigrations runs database migrations on the database specified in the config file
func (a *Application) RunDatabaseMigrations() error {
	db, err := sql.Open("mysql", a.config.SQL.DSN)
	if err != nil {
		return err
	}

	defer func() {
		dErr := db.Close()
		if dErr != nil {
			fmt.Println("Error in deferred close:", dErr)
		}
	}()

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
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

func onHTTPError(statusCode int, statusMessage, errorMessage string) {
	fmt.Println("HTTPERROR: ", errorMessage)
}

func onInternalError(err error) {
	fmt.Printf("internal error: %s", err)
}

// InitializeAPIs initializes various APIs that are needed for pajbot
func (a *Application) InitializeAPIs() error {
	// Twitch APIs
	apirequest.Twitch = gotwitch.New(a.config.Auth.Twitch.User.ClientID)
	apirequest.TwitchBot = gotwitch.New(a.config.Auth.Twitch.Bot.ClientID)
	apirequest.TwitchV3 = gotwitch.NewV3(a.config.Auth.Twitch.User.ClientID)
	apirequest.TwitchBotV3 = gotwitch.NewV3(a.config.Auth.Twitch.Bot.ClientID)

	onSuccess := func(data []gotwitch.User) {
		fmt.Printf("%#v\n", data)
	}

	apirequest.Twitch.GetUsersByLogin([]string{"bajlada"}, onSuccess, onHTTPError, onInternalError)

	/*
		apirequest.Twitch.SubscribeFollows("19571641", "http://57552418.ngrok.io/api/callbacks/follow", func() {
			fmt.Println("success")
		}, func() {
			fmt.Println("error")
		})
	*/

	apirequest.Twitch.SubscribeStreams("159849156", "http://57552418.ngrok.io/api/callbacks/streams", func() {
		fmt.Println("streams success")
	}, func() {
		fmt.Println("streams error")
	})

	apirequest.Twitch.SubscribeStreams("11148817", "http://57552418.ngrok.io/api/callbacks/streams", func() {
		fmt.Println("streams success")
	}, func() {
		fmt.Println("streams error")
	})

	return nil
}

// LoadExternalEmotes xd
func (a *Application) LoadExternalEmotes() error {
	fmt.Println("Loading globalemotes...")
	go emotes.LoadGlobalEmotes()
	fmt.Println("Done!")

	return nil
}

func (a *Application) StartRedisClient() error {
	a.Redis = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", a.config.Redis.Host)
			if err != nil {
				log.Fatal("An error occured while connecting to redis: ", err)
				return nil, err
			}
			if a.config.Redis.Database >= 0 {
				_, err = c.Do("SELECT", a.config.Redis.Database)
				if err != nil {
					log.Fatal("Error while selecting redis db:", err)
					return nil, err
				}
			}
			return c, err
		},
	}

	// Ensure that the redis connection works
	conn := a.Redis.Get()
	return conn.Send("PING")
}

func (a *Application) StartSQLClient() error {
	var err error
	a.SQL, err = sql.Open("mysql", a.config.SQL.DSN)

	return err
}

// StartWebServer starts the web server associated to the bot
func (a *Application) StartWebServer() error {
	webCfg := &web.Config{
		Redis: a.Redis,
		SQL:   a.SQL,
	}

	webBoss := web.Init(a.config, webCfg)
	go webBoss.Run()

	return nil
}

type UnicodeRange struct {
	Start rune
	End   rune
}

func checkModules(next pb2twitch.Handler) pb2twitch.Handler {
	return pb2twitch.HandlerFunc(func(bot *pb2twitch.Bot, channel pkg.Channel, user pkg.User, message *pb2twitch.TwitchMessage, action pkg.Action) {
		modulesStart := time.Now()
		defer func() {
			modulesEnd := time.Now()

			if pkg.VerboseBenchmark {
				fmt.Printf("[% 26s] %s", "Total", modulesEnd.Sub(modulesStart))
			}
		}()

		for _, module := range bot.Modules {
			moduleStart := time.Now()
			var err error
			if channel == nil {
				err = module.OnWhisper(bot, user, message)
			} else {
				err = module.OnMessage(bot, channel, user, message, action)
			}
			moduleEnd := time.Now()
			if pkg.VerboseBenchmark {
				fmt.Printf("[% 26s] %s", module.Name(), moduleEnd.Sub(moduleStart))
			}
			if err != nil {
				fmt.Printf("%s: %s\n", module.Name(), err)
			}
		}

		next.HandleMessage(bot, channel, user, message, action)
	})
}

// LoadBots loads bots from the database
func (a *Application) LoadBots() error {
	db, err := sql.Open("mysql", a.config.SQL.DSN)
	if err != nil {
		return err
	}

	defer func() {
		dErr := db.Close()
		if dErr != nil {
			fmt.Println("Error in deferred close:", dErr)
		}
	}()

	rows, err := db.Query("SELECT `name`, `twitch_access_token` FROM `pb_bot`")
	if err != nil {
		return err
	}

	defer func() {
		dErr := rows.Close()
		if dErr != nil {
			fmt.Println("Error in deferred rows close:", dErr)
		}
	}()

	/*
	 Sorry :( To prevent racism we only allow basic Latin Letters with some exceptions. If you think your message should not have been timed out, please send a link to YOUR chatlogs for the MONTH with a TIMESTAMP of the offending message to "omgscoods@gmail.com" and we'll review it.
	*/

	err = modules.InitServer(a.Redis, a.SQL, a.config.Pajbot1)
	if err != nil {
		return err
	}

	err = users.InitServer(a.SQL)
	if err != nil {
		return err
	}

	for rows.Next() {
		var name string
		var twitchAccessToken string
		if err := rows.Scan(&name, &twitchAccessToken); err != nil {
			return err
		}

		finalHandler := pb2twitch.HandlerFunc(pb2twitch.FinalMiddleware)

		bot := pb2twitch.NewBot(twitch.NewClient(name, "oauth:"+twitchAccessToken))
		bot.Name = name
		bot.QuitChannel = a.Quit

		// Parsing
		bot.AddModule(modules.NewBTTVEmoteParser(&emotes.GlobalEmotes.Bttv))

		// Report module/Admin commands
		bot.AddModule(modules.NewReportModule())

		// Filtering
		bot.AddModule(modules.NewBadCharacterFilter())
		bot.AddModule(modules.NewLatinFilter())
		bot.AddModule(modules.NewPajbot1BanphraseFilter())
		bot.AddModule(modules.NewEmoteFilter(bot))
		bot.AddModule(modules.NewBannedNames())
		bot.AddModule(modules.NewLinkFilter())

		bot.AddModule(modules.NewMessageLengthLimit())

		// Actions
		bot.AddModule(modules.NewActionPerformer())

		// Commands
		bot.AddModule(modules.NewPajbot1Commands(bot))

		customCommands := modules.NewCustomCommands()
		customCommands.RegisterCommand([]string{"!userid"}, &commands.GetUserID{})
		customCommands.RegisterCommand([]string{"!pb2points"}, &commands.GetPoints{})
		customCommands.RegisterCommand([]string{"!pb2roulette"}, &commands.Roulette{})
		customCommands.RegisterCommand([]string{"!pb2givepoints"}, &commands.GivePoints{})
		// customCommands.RegisterCommand([]string{"!pb2addpoints"}, &commands.AddPoints{})
		// customCommands.RegisterCommand([]string{"!pb2removepoints"}, &commands.RemovePoints{})
		customCommands.RegisterCommand([]string{"!roffle", "!join"}, commands.NewRaffle())
		customCommands.RegisterCommand([]string{"!user"}, &commands.User{})
		customCommands.RegisterCommand([]string{"!pb2rank"}, &commands.Rank{})

		bot.AddModule(customCommands)

		bot.AddModule(modules.NewGiveaway(bot))

		// Moderation
		bot.AddModule(modules.NewNuke())

		bot.SetHandler(checkModules(pb2twitch.HandleCommands(finalHandler)))

		a.TwitchBots[name] = bot
	}

	return nil
}

// StartBots starts bots that were loaded from the LoadBots method
func (a *Application) StartBots() error {
	for _, bot := range a.TwitchBots {
		go func(bot *pb2twitch.Bot) {
			if bot.Name != "snusbot" {
				// continue
			}

			bot.OnNewWhisper(bot.HandleWhisper)

			bot.OnNewMessage(bot.HandleMessage)

			bot.OnNewRoomstateMessage(bot.HandleRoomstateMessage)

			if bot.Name == "snusbot" {
				bot.Join("forsen")
			}

			if bot.Name == "botnextdoor" {
				bot.Join("nymn")
			}

			if bot.Name == "pajbot" {
				bot.Join("krakenbul")
				bot.Join("nani")
				bot.Join("pajlada")
				bot.Join("narwhal_dave")
				// err := bot.ConnectToPointServer()
				// if err != nil {
				// 	log.Fatal(err)
				// }
				// bot.StartChatterPoller()
			}

			bot.Join(bot.Name)

			fmt.Printf("Connecting... %#v", bot)
			err := bot.Connect()
			if err != nil {
				log.Fatal(err)
			}
		}(bot)
	}

	return nil
}

func (a *Application) StartPubSubClient() error {
	cfg := &a.config.PubSub
	a.TwitchPubSub = twitch_pubsub.NewClient()

	err := a.TwitchPubSub.Connect()
	if err != nil {
		return err
	}

	if cfg.ChannelID == "" || cfg.UserID == "" || cfg.UserToken == "" {
		return errors.New("Missing PubSub configuration stuff")
	}

	moderationTopic := fmt.Sprintf("chat_moderator_actions.%s.%s", cfg.UserID, cfg.ChannelID)
	fmt.Println("Moderation topic:", moderationTopic)
	a.TwitchPubSub.Listen(moderationTopic, cfg.UserToken, func(bytes []byte) error {
		msg := twitch_pubsub.Message{}
		err := json.Unmarshal(bytes, &msg)
		if err != nil {
			return err
		}

		timeoutData := twitch_pubsub.TimeoutData{}
		err = json.Unmarshal([]byte(msg.Data.Message), &timeoutData)
		if err != nil {
			return err
		}

		const ActionUnknown = 0
		const ActionTimeout = 1
		const ActionBan = 2
		const ActionUnban = 3
		duration := 0

		content := fmt.Sprintf("lol %+v", timeoutData.Data)
		fmt.Println(content)
		var actionContext *string
		action := 0
		reason := ""
		const queryF = "INSERT INTO `ModerationAction` (ChannelID, UserID, Action, Duration, TargetID, Reason, Context) VALUES (?, ?, ?, ?, ?, ?, ?);"
		switch timeoutData.Data.ModerationAction {
		case "timeout":
			action = ActionTimeout
			content = fmt.Sprintf("%s timed out %s for %s seconds", timeoutData.Data.CreatedBy, timeoutData.Data.Arguments[0], timeoutData.Data.Arguments[1])
			duration, _ = strconv.Atoi(timeoutData.Data.Arguments[1])
			if len(timeoutData.Data.Arguments[2]) > 0 {
				reason = timeoutData.Data.Arguments[2]
				content += " for reason: \"" + reason + "\""
			}
			msgs, err := a.GetUserMessages(cfg.ChannelID, timeoutData.Data.TargetUserID)
			if err == nil {
				actionContext = lol(strings.Join(msgs, "\n"))
			}

		case "ban":
			action = ActionBan
			content = fmt.Sprintf("%s banned %s", timeoutData.Data.CreatedBy, timeoutData.Data.Arguments[0])
			if len(timeoutData.Data.Arguments[1]) > 0 {
				reason = timeoutData.Data.Arguments[1]
				content += " for reason: \"" + reason + "\""
			}
			msgs, err := a.GetUserMessages(cfg.ChannelID, timeoutData.Data.TargetUserID)
			if err == nil {
				actionContext = lol(strings.Join(msgs, "\n"))
			}

		case "unban", "untimeout":
			action = ActionUnban
			content = fmt.Sprintf("%s unbanned %s", timeoutData.Data.CreatedBy, timeoutData.Data.Arguments[0])
		}

		if action != 0 {
			stmt, err := a.SQL.Prepare(queryF)
			if err != nil {
				return err
			}
			_, err = stmt.Exec(cfg.ChannelID, timeoutData.Data.CreatedByUserID, action, duration, timeoutData.Data.TargetUserID, reason, actionContext)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

// Run blocks the current thread, waiting for something to put an exit string into the Quit channel
func (a *Application) Run() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		a.Quit <- "Quitting due to SIGTERM/SIGINT"
	}()

	quitString := <-a.Quit

	return fmt.Errorf(quitString)
}
