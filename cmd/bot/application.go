package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	_ "github.com/lib/pq" // PostgreSQL Driver

	twitch "github.com/gempir/go-twitch-irc/v2"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/apirequest"
	"github.com/pajbot/pajbot2/pkg/auth"
	"github.com/pajbot/pajbot2/pkg/botstore"
	"github.com/pajbot/pajbot2/pkg/channels"
	"github.com/pajbot/pajbot2/pkg/common/config"
	"github.com/pajbot/pajbot2/pkg/emotes"
	"github.com/pajbot/pajbot2/pkg/mimo"
	"github.com/pajbot/pajbot2/pkg/modules"
	"github.com/pajbot/pajbot2/pkg/pubsub"
	"github.com/pajbot/pajbot2/pkg/report"
	pb2twitch "github.com/pajbot/pajbot2/pkg/twitch"
	"github.com/pajbot/pajbot2/pkg/users"
	"github.com/pajbot/pajbot2/pkg/web"
	"github.com/pajbot/pajbot2/pkg/web/controller"
	"github.com/pajbot/pajbot2/pkg/web/state"
	"github.com/pajbot/pajbot2/pkg/web/views"
	"github.com/pajbot/utils"
	twitchpubsub "github.com/pajlada/go-twitch-pubsub"
	"github.com/pajlada/stupidmigration"
)

// Application is the heart of pajbot
// It keeps the functions to initialize, start, and stop pajbot
type Application struct {
	config *config.Config

	mimo pkg.MIMO

	twitchBots   pkg.BotStore
	sqlClient    *sql.DB
	Twitter      *twitter.Client
	TwitchPubSub *twitchpubsub.Client

	ReportHolder *report.Holder

	Quit chan string

	pubSub *pubsub.PubSub

	twitchUserStore    pkg.UserStore
	twitchUserContext  pkg.UserContext
	twitchStreamStore  *StreamStore
	twitchChannelStore pkg.ChannelStore

	// Oauth configs
	twitchAuths *auth.TwitchAuths
}

var _ pkg.PubSubSource = &Application{}
var _ pkg.Application = &Application{}

func (a *Application) MIMO() pkg.MIMO {
	return a.mimo
}

// NewApplication creates an instance of Application. Generally this should only be done once
func newApplication() *Application {
	a := &Application{
		mimo: mimo.New(),

		twitchBots: botstore.New(),
	}

	a.twitchUserStore = NewUserStore()
	state.StoreTwitchUserStore(a.twitchUserStore)
	a.twitchUserContext = NewUserContext()
	a.twitchStreamStore = NewStreamStore()
	a.twitchChannelStore = channels.NewStore()
	state.StoreTwitchChannelStore(a.twitchChannelStore)

	a.Quit = make(chan string)
	a.pubSub = pubsub.New()
	state.StorePubSub(a.pubSub)
	state.StoreApplication(a)

	go a.pubSub.Run()

	return a
}

func (a *Application) UserStore() pkg.UserStore {
	return a.twitchUserStore
}

func (a *Application) ChannelStore() pkg.ChannelStore {
	return a.twitchChannelStore
}

func (a *Application) UserContext() pkg.UserContext {
	return a.twitchUserContext
}

func (a *Application) StreamStore() pkg.StreamStore {
	return a.twitchStreamStore
}

func (a *Application) SQL() *sql.DB {
	return a.sqlClient
}

func (a *Application) PubSub() pkg.PubSub {
	return a.pubSub
}

func (a *Application) TwitchBots() pkg.BotStore {
	return a.twitchBots
}

func (a *Application) QuitChannel() chan string {
	return a.Quit
}

func (a *Application) TwitchAuths() pkg.TwitchAuths {
	return a.twitchAuths
}

func (a *Application) IsApplication() bool {
	return true
}

func (a *Application) Connection() pkg.PubSubConnection {
	return nil
}

func (a *Application) AuthenticatedUser() pkg.User {
	return nil
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

// https://dev.twitch.tv/docs/authentication/#scopes
func (a *Application) InitializeOAuth2Configs() (err error) {
	a.twitchAuths, err = auth.NewTwitchAuths(&a.config.Auth.Twitch)

	return
}

// RunDatabaseMigrations runs database migrations on the database specified in the config file
func (a *Application) RunDatabaseMigrations() error {
	sqlClient, err := sql.Open("postgres", a.config.PostgreSQL.DSN)
	if err != nil {
		return err
	}
	err = stupidmigration.Migrate("../../migrations/psql", sqlClient)
	if err != nil {
		fmt.Println("Unable to run SQL migrations", err)
		return err
	}

	return nil
}

func (a *Application) ProvideAdminPermissionsToAdmin() (err error) {
	cfg := a.config.Admin
	if !utils.IsValidUserID(cfg.TwitchUserID) {
		fmt.Println("Warning: No admin user ID specified in the config file. You probably want to do this at least on initial setup")
		return
	}

	oldPermissions, err := users.GetUserPermissions(cfg.TwitchUserID, "global")
	if err != nil {
		return errors.Wrap(err, "get failed")
	}
	newPermissions := oldPermissions | pkg.PermissionAdmin
	err = users.SetUserPermissions(cfg.TwitchUserID, "global", newPermissions)
	if err != nil {
		return errors.Wrap(err, "set failed")
	}

	return
}

// InitializeAPIs initializes various APIs that are needed for pajbot
func (a *Application) InitializeAPIs() (err error) {
	err = apirequest.InitTwitch(a.config)
	if err != nil {
		return
	}

	return
}

// LoadExternalEmotes xd
func (a *Application) LoadExternalEmotes() error {
	go emotes.LoadGlobalEmotes()

	return nil
}

func (a *Application) InitializeSQL() error {
	var err error
	a.sqlClient, err = sql.Open("postgres", a.config.PostgreSQL.DSN)
	if err != nil {
		return err
	}

	err = a.sqlClient.Ping()
	if err != nil {
		return err
	}

	state.StoreSQL(a.sqlClient)

	users.InitServer(a.sqlClient)

	return nil
}

func (a *Application) InitializeModules() (err error) {
	// TODO: move this to init
	a.ReportHolder, err = report.New(a)
	if err != nil {
		return errors.Wrap(err, "initializing report holder")
	}

	a.twitterTest()

	err = a.StartTwitterStream()
	if err != nil {
		fmt.Println("Error starting twitter stream:", err)
	}

	err = modules.InitServer(a, &a.config.Pajbot1, a.ReportHolder)
	if err != nil {
		return errors.Wrap(err, "initializing modules server")
	}

	moduleList := []string{}
	for _, module := range modules.List() {
		moduleList = append(moduleList, module.ID())
	}

	fmt.Println("Available modules:", moduleList)

	return
}

func (a *Application) twitterTest() {
	localConfig := a.config.Auth.Twitter

	config := oauth1.NewConfig(localConfig.ConsumerKey, localConfig.ConsumerSecret)
	token := oauth1.NewToken(localConfig.AccessToken, localConfig.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	user, resp, err := client.Users.Show(&twitter.UserShowParams{
		ScreenName: "pajtest",
	})
	fmt.Println(user)
	fmt.Println(resp)
	fmt.Println(err)
}

func (a *Application) StartTwitterStream() error {
	localConfig := a.config.Auth.Twitter

	// users lookup
	// userLookupParams := &twitter.UserLookupParams{ScreenName: []string{"pajtest"}}
	// users, _, _ := client.Users.Lookup(userLookupParams)
	// fmt.Printf("USERS LOOKUP:\n%+v\n", users)

	if localConfig.ConsumerKey == "" || localConfig.ConsumerSecret == "" || localConfig.AccessToken == "" || localConfig.AccessSecret == "" {
		return errors.New("Missing twitter configuration fields")
	}

	// config := &oauth2.Config{}
	// token := &oauth2.Token{AccessToken: localConfig.AccessToken}
	// httpClient := config.Client(oauth2.NoContext, token)

	config := oauth1.NewConfig(localConfig.ConsumerKey, localConfig.ConsumerSecret)
	token := oauth1.NewToken(localConfig.AccessToken, localConfig.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	params := &twitter.StreamFilterParams{
		Follow:        []string{"81085011"},
		StallWarnings: twitter.Bool(true),
	}

	writer := a.mimo.Publisher("twitter")
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		fmt.Println("got tweet:", tweet.Text)
		writer <- tweet.Text
		fmt.Println("done publishing tweet")
	}
	stream, err := client.Streams.Filter(params)
	if err != nil {
		fmt.Println("Error getting twitter stream:", err)
		return err
	}
	go demux.HandleChan(stream.Messages)
	fmt.Println("Stream:", stream)

	/*
		config := oauth1.NewConfig(localConfig.ConsumerKey, localConfig.ConsumerSecret)
		token := oauth1.NewToken(localConfig.AccessToken, localConfig.AccessSecret)

		httpClient := config.Client(oauth1.NoContext, token)
		client := twitter.NewClient(httpClient)

		demux := twitter.NewSwitchDemux()
		demux.All = func(x interface{}) {
			fmt.Printf("x %#v\n", x)
		}
		demux.StreamDisconnect = func(disconnect *twitter.StreamDisconnect) {
			fmt.Printf("disconnected %#v\n", disconnect)
		}
		demux.Tweet = func(tweet *twitter.Tweet) {
			fmt.Println(tweet.Text)
		}

		demux.Event = func(event *twitter.Event) {
			fmt.Printf("%#v\n", event)
		}

		filterParams := &twitter.StreamFilterParams{
			// Follow:        []string{"81085011"},
			// Track: []string{"cat"},
			// StallWarnings: twitter.Bool(true),
		}
		stream, err := client.Streams.Filter(filterParams)
		if err != nil {
			return err
		}

		fmt.Printf("stream is %#v\n", stream)
		fmt.Printf("messages is is %#v\n", stream.Messages)

		fmt.Println("start handling..")
		for message := range stream.Messages {
			fmt.Printf("got message %#v\n", message)
			// demux.Handle(message)
		}
		_, xd := (<-stream.Messages)
		if xd {
			fmt.Println("channel is not closed")
		}
		fmt.Printf("messages is is %#v\n", stream.Messages)
		fmt.Println("done")
	*/

	return nil
}

// StartWebServer starts the web server associated to the bot
func (a *Application) StartWebServer() error {
	var WSHost string

	if a.config.Web.Secure {
		WSHost = "wss://" + a.config.Web.Domain + "/ws"
	} else {
		WSHost = "ws://" + a.config.Web.Domain + "/ws"
	}

	go web.Run(&a.config.Web)

	controller.LoadRoutes(a, a.config)

	views.Configure(views.Config{
		WSHost: WSHost,
	})

	return nil
}

func (a *Application) loadBot(botConfig botConfig) error {
	bot, err := pb2twitch.NewBot(botConfig.databaseID, botConfig.account, botConfig.tokenSource, a)
	if err != nil {
		return err
	}

	botUserID := botConfig.account.ID()
	bot.OnNewChannelJoined(func(channel pkg.Channel) {
		var token *oauth2.Token
		token, err = bot.GetTokenSource().Token()
		if err != nil {
			fmt.Println("Error renewing token: ", err)
			return
		}
		a.TwitchPubSub.Listen(twitchpubsub.ModerationActionTopic(botUserID, channel.GetID()), token.AccessToken)

		a.twitchChannelStore.RegisterTwitchChannel(channel)
	})

	err = bot.LoadChannels(a.sqlClient)
	if err != nil {
		return err
	}

	a.twitchBots.Add(bot)

	return nil
}

type botConfig struct {
	databaseID  int
	account     *pb2twitch.TwitchAccount
	tokenSource oauth2.TokenSource
}

// LoadBots loads bots from the database
func (a *Application) LoadBots() (err error) {
	const queryF = `SELECT id, twitch_userid, twitch_username, twitch_access_token, twitch_refresh_token, twitch_access_token_expiry FROM bot`
	rows, err := a.sqlClient.Query(queryF) // GOOD
	if err != nil {
		return err
	}

	defer func() {
		dErr := rows.Close()
		if dErr != nil {
			fmt.Println("Error in deferred rows close:", dErr)
		}
	}()

	var botsMutex sync.Mutex
	var bots []botConfig

	var wg sync.WaitGroup

	for rows.Next() {
		var databaseID int
		var acc pb2twitch.TwitchAccount
		var c pb2twitch.BotCredentials
		if err := rows.Scan(&databaseID, &acc.UserID, &acc.UserName, &c.AccessToken, &c.RefreshToken, &c.Expiry); err != nil {
			return err
		}

		wg.Add(1)

		c.AccessToken = strings.TrimPrefix(c.AccessToken, "oauth:")

		go func(databaseID int, acc *pb2twitch.TwitchAccount, c pb2twitch.BotCredentials) {
			defer wg.Done()
			oldToken := &oauth2.Token{
				AccessToken:  c.AccessToken,
				TokenType:    "bearer",
				RefreshToken: c.RefreshToken,
				Expiry:       c.Expiry.Time,
			}

			tokenSource := a.twitchAuths.Bot().TokenSource(context.Background(), oldToken)
			refreshingTokenSource := oauth2.ReuseTokenSource(oldToken, tokenSource)

			botsMutex.Lock()
			defer botsMutex.Unlock()
			bots = append(bots, botConfig{
				databaseID:  databaseID,
				account:     acc,
				tokenSource: refreshingTokenSource,
			})

			token, err := refreshingTokenSource.Token()
			if err != nil {
				fmt.Println("ERROR!!! Error getting token from token source for bot", acc.Name())
				return
			}

			self, _, err := apirequest.TwitchBot.ValidateOAuthTokenSimple(token.AccessToken)
			if err != nil {
				fmt.Println("ERROR!!! Error validating oauth token for", acc.Name(), err)
				return
			}

			if self.UserID != acc.ID() {
				fmt.Println("ERROR!!! User ID for", acc.Name(), acc.ID(), "doesn't match the api response:", self.UserID)
				return
			}
		}(databaseID, &acc, c)
	}

	wg.Wait()

	for _, botConfig := range bots {
		err = a.loadBot(botConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

// StartBots starts bots that were loaded from the LoadBots method
func (a *Application) StartBots() error {
	for it := a.twitchBots.Iterate(); it.Next(); {
		bot := it.Value()
		if bot == nil {
			fmt.Println("nil bot DansGame")
			continue
		}

		pb2bot, ok := bot.(*pb2twitch.Bot)
		if !ok {
			fmt.Println("Unknown bot")
			continue
		}

		go func(bot *pb2twitch.Bot) {
			bot.OnWhisperMessage(bot.HandleWhisper)

			bot.OnPrivateMessage(func(message twitch.PrivateMessage) {
				channelID := message.Tags["room-id"]
				if channelID == "" {
					fmt.Printf("Missing room-id tag in message: %+v\n", message)
					return
				}

				if message.User.ID == bot.TwitchAccount().ID() {
					// Ignore messages from self
					return
				}

				formattedMessage := fmt.Sprintf("[%s] %s: %s", time.Now().Format("15:04:05"), message.User.Name, message.Message)

				// Store message in our twitch message context class
				a.twitchUserContext.AddContext(channelID, message.User.ID, formattedMessage)

				message.Message = strings.TrimPrefix(message.Message, "@"+bot.TwitchAccount().Name()+" ")
				message.Message = strings.TrimPrefix(message.Message, bot.TwitchAccount().Name()+" ")

				// Forward to bot to let its modules work
				bot.HandleMessage(message.Channel, message.User, &message)
			})

			bot.OnRoomStateMessage(bot.HandleRoomstateMessage)

			bot.OnClearChatMessage(func(message twitch.ClearChatMessage) {
				bot.HandleClearChatMessage(&message)
			})

			bot.OnUnsetMessage(func(message twitch.RawMessage) {
				fmt.Println("Unparsed message:", message.Raw)
			})

			// Ensure that the bot has joined its own chat
			bot.JoinChannel(bot.TwitchAccount().ID())

			// TODO: Join some "central control center" like skynetcentral

			// Join all "external" channels
			bot.JoinChannels()

			// err := bot.ConnectToPointServer()
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// bot.StartChatterPoller()

			bot.IsConnected = true
			err := bot.Connect()
			if err != nil {
				if err == twitch.ErrLoginAuthenticationFailed {
					fmt.Printf("%s: Login authentication failed\n", bot.TwitchAccount().Name())
				} else {
					log.Fatal(err)
				}
			}
			bot.IsConnected = false

			// TODO: Check if we want to try to reconnect here (with a delay)
		}(pb2bot)
	}

	go a.twitchStreamStore.Run()

	return nil
}

func (a *Application) StartPubSubClient() error {
	a.TwitchPubSub = twitchpubsub.NewClient(twitchpubsub.DefaultHost)

	a.TwitchPubSub.OnModerationAction(func(channelID string, event *twitchpubsub.ModerationAction) {
		const ActionUnknown = 0
		const ActionTimeout = 1
		const ActionBan = 2
		const ActionUnban = 3
		var actionContext string
		duration := 0

		fullContext := a.twitchUserContext.GetContext(channelID, event.TargetUserID)
		if fullContext != nil {
			actionContext = fullContext[len(fullContext)-1]
		}
		action := 0
		reason := ""
		const queryF = "INSERT INTO moderation_action (channel_id, user_id, action, duration, target_id, reason, context) VALUES ($1, $2, $3, $4, $5, $6, $7);"
		switch event.ModerationAction {
		case "timeout":
			action = ActionTimeout
			duration, _ = strconv.Atoi(event.Arguments[1])
			if len(event.Arguments[2]) > 0 {
				reason = event.Arguments[2]
			}

			e := pkg.PubSubTimeoutEvent{
				Channel: pkg.PubSubUser{
					ID: channelID,
				},
				Target: pkg.PubSubUser{
					ID:   event.TargetUserID,
					Name: event.Arguments[0],
				},
				Source: pkg.PubSubUser{
					ID:   event.CreatedByUserID,
					Name: event.CreatedBy,
				},
				Duration: duration,
				Reason:   reason,
			}

			a.pubSub.Publish(a, "TimeoutEvent", e)

		case "ban":
			action = ActionBan
			if len(event.Arguments[1]) > 0 {
				reason = event.Arguments[1]
			}

			e := pkg.PubSubBanEvent{
				Channel: pkg.PubSubUser{
					ID: channelID,
				},
				Target: pkg.PubSubUser{
					ID:   event.TargetUserID,
					Name: event.Arguments[0],
				},
				Source: pkg.PubSubUser{
					ID:   event.CreatedByUserID,
					Name: event.CreatedBy,
				},
				Reason: reason,
			}

			a.pubSub.Publish(a, "BanEvent", e)

		case "unban", "untimeout":
			action = ActionUnban
		}

		if action != 0 {
			_, err := a.sqlClient.Exec(queryF, channelID, event.CreatedByUserID, action, duration, event.TargetUserID, reason, actionContext) // GOOD
			if err != nil {
				fmt.Println("Error in moderation action callback:", err)
				return
			}
		}
	})

	go a.TwitchPubSub.Start()

	return nil
}

func (a *Application) listenToModeratorActions(userID, channelID, userToken string) error {
	moderationTopic := twitchpubsub.ModerationActionTopic(userID, channelID)
	a.TwitchPubSub.Listen(moderationTopic, userToken)

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
