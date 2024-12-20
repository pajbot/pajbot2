package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"strconv"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	_ "github.com/lib/pq" // PostgreSQL Driver

	twitch "github.com/gempir/go-twitch-irc/v4"
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

	if err := cfg.Validate(); err != nil {
		return err
	}

	a.config = cfg

	return nil
}

func (a *Application) InitializeOAuth2Configs() (err error) {
	// https://dev.twitch.tv/docs/authentication/#scopes
	a.twitchAuths, err = auth.NewTwitchAuths(&a.config.Auth.Twitch, &a.config.Web)

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

func (a *Application) botOnNewChannelJoined(bot *pb2twitch.Bot, botUserID string) func(pkg.Channel) {
	return func(channel pkg.Channel) {
		accessToken, err := bot.GetAccessToken()
		if err != nil {
			fmt.Println("Error getting access token for", bot.TwitchAccount().Name())
			return
		}

		a.TwitchPubSub.Listen(twitchpubsub.ModerationActionTopic(botUserID, channel.GetID()), accessToken)

		a.twitchChannelStore.RegisterTwitchChannel(channel)
	}
}

func (a *Application) loadBot(botConfig botConfig) error {
	bot, err := pb2twitch.NewBot(botConfig.databaseID, botConfig.account, botConfig.tokenSource, botConfig.token, botConfig.helixClient, a)
	if err != nil {
		return fmt.Errorf("loadBot NewBot: %w", err)
	}

	bot.OnNewChannelJoined(a.botOnNewChannelJoined(bot, botConfig.account.ID()))

	err = bot.LoadChannels(a.sqlClient)
	if err != nil {
		return fmt.Errorf("loadChannels NewBot: %w", err)
	}

	a.twitchBots.Add(bot)

	return nil
}

func appendBotConfig(wg *sync.WaitGroup, botsMutex sync.Locker, bots *[]botConfig, databaseID int, acc pb2twitch.TwitchAccount, c pb2twitch.BotCredentials, oauthConfig *oauth2.Config, sqlClient *sql.DB, clientID string) error {
	defer wg.Done()

	bc, err := newBotConfig(databaseID, &acc, c, oauthConfig, clientID)
	if err != nil {
		return err
	}

	if err := bc.Validate(sqlClient); err != nil {
		fmt.Println("Error validating bot config:", err)
		return err
	}

	botsMutex.Lock()
	defer botsMutex.Unlock()
	*bots = append(*bots, bc)

	return nil
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

	oauthConfig := a.twitchAuths.Bot()

	for rows.Next() {
		var databaseID int
		var acc pb2twitch.TwitchAccount
		var c pb2twitch.BotCredentials
		if err := rows.Scan(&databaseID, &acc.UserID, &acc.UserName, &c.AccessToken, &c.RefreshToken, &c.Expiry); err != nil {
			return err
		}

		wg.Add(1)

		c.AccessToken = strings.TrimPrefix(c.AccessToken, "oauth:")

		go func() {
			err := appendBotConfig(&wg, &botsMutex, &bots, databaseID, acc, c, oauthConfig, a.sqlClient, a.config.Auth.Twitch.Bot.ClientID)
			if err != nil {
				fmt.Println("Error appending bot config:", err)
			}
		}()
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

				// Trim message off any potential Chatterino suffix
				message.Message = strings.TrimSuffix(message.Message, " \U000e0000")

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

			for {
				bot.IsConnected = true
				accessToken, err := bot.GetAccessToken()
				if err != nil {
					fmt.Printf("Unable to refresh token for %s: %s\n", bot.TwitchAccount().Name(), err)
				}
				bot.SetIRCToken("oauth:" + accessToken)
				err = bot.Connect()
				switch err {
				case twitch.ErrLoginAuthenticationFailed:
					fmt.Printf("%s: Login authentication failed\n", bot.TwitchAccount().Name())
				default:
					fmt.Println("Unhandled error from go-twitch-irc Connect:", err)
				}
				bot.IsConnected = false
			}

			// TODO: Check if we want to try to reconnect here (with a delay)
		}(pb2bot)
	}

	go a.twitchStreamStore.Run()

	return nil
}

func (a *Application) StartPubSubClient() error {
	a.TwitchPubSub = twitchpubsub.NewClient(twitchpubsub.DefaultHost)

	a.TwitchPubSub.OnModerationAction(func(channelID string, event *twitchpubsub.ModerationAction) {
		fmt.Println("Got moderation action")
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

	return errors.New(quitString)
}
