package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"net"

	"encoding/json"

	"errors"
	"strconv"

	"github.com/dankeroni/gotwitch"
	"github.com/gempir/go-twitch-irc"
	"github.com/goware/urlx"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/mysql"
	"github.com/pajlada/go-twitch-pubsub"
	"github.com/pajlada/pajbot2/apirequest"
	"github.com/pajlada/pajbot2/bots"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/common/config"
	"github.com/pajlada/pajbot2/emotes"
	"github.com/pajlada/pajbot2/filter"
	pb "github.com/pajlada/pajbot2/grpc"
	"github.com/pajlada/pajbot2/pajbot1"
	"github.com/pajlada/pajbot2/redismanager"
	"github.com/pajlada/pajbot2/sqlmanager"
	"github.com/pajlada/pajbot2/web"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"mvdan.cc/xurls"
)

const migrationsPath = "file://migrations"

func maxpenis(a, b int) int {
	if a > b {
		return a
	}

	return b
}

type pajbotServer struct{}

func (s *pajbotServer) CheckMessage(ctx context.Context, in *pb.MessageRequest) (*pb.MessageAction, error) {
	action := &pb.MessageAction{}

	// url checker
	matchedURLs := xurls.Relaxed().FindAllString(in.GetMessage(), -1)
	for _, matchedURL := range matchedURLs {

		parsedURL, err := urlx.Parse(matchedURL)
		if err != nil {
			return nil, err
		}

		badURL := true

		hostname := "." + parsedURL.Hostname()

		for _, goodURL := range validURLs {
			if strings.HasSuffix(hostname, goodURL) {
				badURL = false
				break
			}
		}

		if badURL {
			/*
				msg := fmt.Sprintf("%s, that's a bad url 😡 FeelsWeirdMan", in.Source.GetDisplayName())
				sayAction := &pb.Action_SayAction{
					SayAction: &pb.SayAction{
						Message: msg,
					},
				}
				action.Actions = append(action.Actions, &pb.Action{Action: sayAction})

				timeoutAction := &pb.Action_TimeoutAction{
					TimeoutAction: &pb.TimeoutAction{
						Target:   in.Source.LoginName,
						Duration: 5,
						Reason:   "Bad link 😡",
					},
				}
				action.Actions = append(action.Actions, &pb.Action{Action: timeoutAction})
			*/
			break
		}
	}

	if strings.Contains(in.GetMessage(), "LOOOOOL 4HEad") {
		msg := fmt.Sprintf("%s, JUST GET A HOUSE 4House", in.Source.GetDisplayName())
		sayAction := &pb.Action_SayAction{
			SayAction: &pb.SayAction{
				Message: msg,
			},
		}
		action.Actions = append(action.Actions, &pb.Action{Action: sayAction})
	}

	return action, nil
}

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

	TwitchBots   map[string]*bots.TwitchBot
	Redis        *redismanager.RedisManager
	SQL          *sqlmanager.SQLManager
	TwitchPubSub *twitch_pubsub.Client
	GRPCClient   pb.ClientClient

	OldPajbot *pajbot1.Pajbot1

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

	a.TwitchBots = make(map[string]*bots.TwitchBot)
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

// InitializeAPIs initializes various APIs that are needed for pajbot
func (a *Application) InitializeAPIs() error {
	// Twitch APIs
	apirequest.Twitch = gotwitch.New(a.config.Auth.Twitch.User.ClientID)
	apirequest.TwitchBot = gotwitch.New(a.config.Auth.Twitch.Bot.ClientID)
	apirequest.TwitchV3 = gotwitch.NewV3(a.config.Auth.Twitch.User.ClientID)
	apirequest.TwitchBotV3 = gotwitch.NewV3(a.config.Auth.Twitch.Bot.ClientID)

	return nil
}

// LoadExternalEmotes xd
func (a *Application) LoadExternalEmotes() error {
	log.Println("Loading globalemotes...")
	go emotes.LoadGlobalEmotes()
	log.Println("Done!")

	return nil
}

// StartWebServer starts the web server associated to the bot
func (a *Application) StartWebServer() error {
	a.Redis = redismanager.Init(a.config.Redis)
	a.SQL = sqlmanager.Init(a.config.SQL)

	webCfg := &web.Config{
		Bots:  a.TwitchBots,
		Redis: a.Redis,
		SQL:   a.SQL,
	}

	webBoss := web.Init(a.config, webCfg)
	go webBoss.Run()

	return nil
}

func addHeheToMessageText(next bots.Handler) bots.Handler {
	return bots.HandlerFunc(func(bot *bots.TwitchBot, channel string, user twitch.User, message *bots.TwitchMessage) {
		message.Text = message.Text + " hehe"
		next.HandleMessage(bot, channel, user, message)
	})
}

var badCharacters = []string{
	"\x01",
}

const badCharacterTimeoutLength = 300

func timeoutBadCharacters(next bots.Handler) bots.Handler {
	return bots.HandlerFunc(func(bot *bots.TwitchBot, channel string, user twitch.User, message *bots.TwitchMessage) {
		for _, badCharacter := range badCharacters {
			if strings.Contains(message.Text, badCharacter) {
				bot.Timeout(channel, user, badCharacterTimeoutLength, "Message contains a forbidden character")
				return
			}
		}

		next.HandleMessage(bot, channel, user, message)
	})
}

func parseBTTVEmotes(next bots.Handler) bots.Handler {
	return bots.HandlerFunc(func(bot *bots.TwitchBot, channel string, user twitch.User, message *bots.TwitchMessage) {
		m := strings.Split(message.Text, " ")
		emoteCount := make(map[string]*common.Emote)
		for _, word := range m {
			if emote, ok := emoteCount[word]; ok {
				emote.Count++
			} else if emote, ok := emotes.GlobalEmotes.Bttv[word]; ok {
				emoteCount[word] = &emote
			}
		}

		for _, emote := range emoteCount {
			message.BTTVEmotes = append(message.BTTVEmotes, *emote)
		}

		next.HandleMessage(bot, channel, user, message)
	})
}

func handleCommands(next bots.Handler) bots.Handler {
	return bots.HandlerFunc(func(bot *bots.TwitchBot, channel string, user twitch.User, message *bots.TwitchMessage) {
		if user.UserType == "mod" {
			if strings.HasPrefix(message.Text, "!xd") {
				bot.Reply(channel, user, "XDDDDDDDDDD")
				return
			}

			if strings.HasPrefix(message.Text, "!whisperme") {
				log.Printf("Send whisper!")
				// bot.Say(channel, "XDDDDDDDDDD")
				bot.Whisper("pajlada", "hehe")
				return
			}

			if strings.HasPrefix(message.Text, "!pb2quit") {
				bot.Reply(channel, user, "Quitting...")
				time.AfterFunc(time.Millisecond*500, func() {
					bot.Quit("Quit because pajlada said so")
				})
				return
			}

			if strings.HasPrefix(message.Text, "!subon") {
				if bot.Flags.PermaSubMode {
					bot.Say(channel, "Permanent subscribers mode is already enabled")
					return
				}

				bot.Flags.PermaSubMode = true

				bot.Say(channel, ".subscribers")
				bot.Say(channel, "Permanent subscribers mode has been enabled")
				return
			}

			if strings.HasPrefix(message.Text, "!suboff") {
				if !bot.Flags.PermaSubMode {
					bot.Say(channel, "Permanent subscribers mode is not enabled")
					return
				}

				bot.Flags.PermaSubMode = false

				bot.Say(channel, ".subscribersoff")
				bot.Say(channel, "Permanent subscribers mode has been disabled")
				return
			}
		}

		next.HandleMessage(bot, channel, user, message)
	})
}

func finalMiddleware(bot *bots.TwitchBot, channel string, user twitch.User, message *bots.TwitchMessage) {
	// log.Printf("Found %d BTTV emotes! %#v", len(message.BTTVEmotes), message.BTTVEmotes)
}

type UnicodeRange struct {
	Start rune
	End   rune
}

var unicodeWhitelist []UnicodeRange

func addToWhitelist(start, end rune) {
	unicodeWhitelist = append(unicodeWhitelist, UnicodeRange{start, end})
}

func addRunesToWhitelist(runes ...rune) {
	for _, r := range runes {
		unicodeWhitelist = append(unicodeWhitelist, UnicodeRange{r, r})
	}
}

func simplifyMessageText(next bots.Handler) bots.Handler {
	return bots.HandlerFunc(func(bot *bots.TwitchBot, channel string, user twitch.User, message *bots.TwitchMessage) {
		next.HandleMessage(bot, channel, user, message)
	})
}

func tryLatinWhitelist(next bots.Handler) bots.Handler {
	return bots.HandlerFunc(func(bot *bots.TwitchBot, channel string, user twitch.User, message *bots.TwitchMessage) {
		if user.UserType == "" || true {
			lol := struct {
				FullMessage   string
				Message       string
				BadCharacters []rune
				Username      string
				Channel       string
				Timestamp     time.Time
			}{
				FullMessage: message.Text,
				Username:    user.Username,
				Channel:     channel,
				Timestamp:   time.Now().UTC(),
			}
			messageRunes := []rune(message.Text)
			transparentSkipRange := bot.TransparentList.Find(messageRunes)
			messageLength := len(messageRunes)
			for i := 0; i < messageLength; {
				if skipLength := transparentSkipRange.ShouldSkip(i); skipLength > 0 {
					i = i + skipLength
					continue
				}

				r := messageRunes[i]
				allowed := false

				for _, allowedRange := range unicodeWhitelist {
					if r >= allowedRange.Start && r <= allowedRange.End {
						allowed = true
						break
					}
				}

				if !allowed {
					if lol.Message == "" {
						lol.Message = message.Text[maxpenis(0, i-2):len(message.Text)]
					}

					alreadySet := false
					for _, bc := range lol.BadCharacters {
						if bc == r {
							alreadySet = true
							break
						}
					}

					if !alreadySet {
						lol.BadCharacters = append(lol.BadCharacters, r)
					}

				}
				i++
			}

			if lol.Message != "" {
				c := bot.Redis.Pool.Get()
				bytes, _ := json.Marshal(&lol)
				c.Do("LPUSH", "karl_kons", bytes)
				c.Close()
				// log.Printf("First bad character: 0x%0x message '%s' from '%s' in '#%s' is disallowed due to our whitelist\n", lol.BadCharacters[0], message.Text, user.Username, channel)
			}
		}

		next.HandleMessage(bot, channel, user, message)
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
			log.Println("Error in deferred close:", dErr)
		}
	}()

	rows, err := db.Query("SELECT `name`, `twitch_access_token` FROM `pb_bot`")
	if err != nil {
		return err
	}

	defer func() {
		dErr := rows.Close()
		if dErr != nil {
			log.Println("Error in deferred rows close:", dErr)
		}
	}()

	/*
	 Sorry :( To prevent racism we only allow basic Latin Letters with some exceptions. If you think your message should not have been timed out, please send a link to YOUR chatlogs for the MONTH with a TIMESTAMP of the offending message to "omgscoods@gmail.com" and we'll review it.
	*/

	addToWhitelist(0x20, 0x7e)       // Basic latin
	addToWhitelist(0x1f600, 0x1f64f) // Emojis
	addToWhitelist(0x1f300, 0x1f5ff) // "Miscellaneous symbols and pictographs". Includes some emojis like 100
	addToWhitelist(0x1f44c, 0x1f44c) // Chatterino?
	addToWhitelist(0x206d, 0x206d)   // Chatterino?
	addToWhitelist(0x2660, 0x2765)   // Chatterino?

	addToWhitelist(0x1f171, 0x1f171) // B emoji
	addToWhitelist(0x1f900, 0x1f9ff) // More emojis

	// Rain
	addToWhitelist(0x30fd, 0x30fd)
	addToWhitelist(0xff40, 0xff40)
	addToWhitelist(0x3001, 0x3001)
	addToWhitelist(0x2602, 0x2602)

	// From Karl
	addToWhitelist(0x1d100, 0x1d1ff)
	addToWhitelist(0x1f680, 0x1f6ff)
	addToWhitelist(0x2600, 0x26ff)
	addToWhitelist(0xfe00, 0xfe0f) // Emoji variation selector 1 to 16
	addToWhitelist(0x2012, 0x2015) // Various dashes
	addToWhitelist(0x3010, 0x3011) // 【 and 】

	for rows.Next() {
		var name string
		var twitchAccessToken string
		if err := rows.Scan(&name, &twitchAccessToken); err != nil {
			return err
		}

		finalHandler := bots.HandlerFunc(finalMiddleware)

		bot := &bots.TwitchBot{
			Client:          twitch.NewClient(name, "oauth:"+twitchAccessToken),
			Name:            name,
			QuitChannel:     a.Quit,
			Redis:           a.Redis,
			TransparentList: filter.NewTransparentList(),
		}

		bot.TransparentList.Add("(/ﾟДﾟ)/")
		bot.TransparentList.Add("(╯°□°）╯︵ ┻━┻")
		bot.TransparentList.Add("(╯°Д°）╯︵/(.□ . )")
		bot.TransparentList.Add("(ノಠ益ಠ)ノ彡┻━┻")
		bot.TransparentList.Add("୧༼ಠ益ಠ༽୨")
		bot.TransparentList.Add("༼ ºل͟º ༽")
		bot.TransparentList.Add("༼つಠ益ಠ༽つ")
		bot.TransparentList.Add("( ° ͜ʖ͡°)╭∩╮")
		bot.TransparentList.Add("ᕙ༼ຈل͜ຈ༽ᕗ")
		bot.TransparentList.Add("ʕ•ᴥ•ʔ")
		bot.TransparentList.Add("༼▀̿ Ĺ̯▀̿༽")
		bot.TransparentList.Add("( ͡° ͜🔴 ͡°)")

		err = bot.TransparentList.Build()
		if err != nil {
			return err
		}

		bot.SetHandler(tryLatinWhitelist(a.OldPajbot.CheckBanphrases(timeoutBadCharacters(addHeheToMessageText(parseBTTVEmotes(handleCommands(finalHandler)))))))

		a.TwitchBots[name] = bot
	}

	return nil
}

func (a *Application) LoadOldPajbot() error {
	a.OldPajbot = pajbot1.Init(a.config.Pajbot1)

	err := a.OldPajbot.LoadBanphrases()
	if err != nil {
		return err
	}

	log.Println("Loaded", len(a.OldPajbot.EnabledBanphrases), "banphrases")

	return nil
}

func (a *Application) StartContextBot() error {
	contextBot := &bots.TwitchBot{
		Client:      twitch.NewClient("justinfan64932", "oauth:b00b5"),
		QuitChannel: a.Quit,
	}

	contextBot.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		if userID, ok := message.Tags["user-id"]; ok {
			if roomID, ok := message.Tags["room-id"]; ok {
				uc, ok := a.UserContext[userID]
				if !ok {
					uc = NewChannelContext()
					a.UserContext[userID] = uc
				}
				uc.Channels[roomID] = append(uc.Channels[roomID], message.Text)
			}
		}
	})

	contextBot.Join("pajlada")

	go func() {
		contextBot.Connect()
	}()

	return nil
}

type ModeState int

const (
	ModeUnset = iota
	ModeEnabled
	ModeDisabled
)

// StartBots starts bots that were loaded from the LoadBots method
func (a *Application) StartBots() error {
	for _, bot := range a.TwitchBots {
		if bot.Name != "snusbot" {
			continue
		}

		bot.OnNewWhisper(func(user twitch.User, rawMessage twitch.Message) {
			message := bots.TwitchMessage{Message: rawMessage}

			// log.Printf("GOT WHISPER! %s(%s): %s", user.DisplayName, user.Username, message.Text)

			bot.HandleMessage("", user, &message)
		})

		bot.OnNewMessage(func(channel string, user twitch.User, rawMessage twitch.Message) {
			message := bots.TwitchMessage{Message: rawMessage}

			bot.HandleMessage(channel, user, &message)

			// log.Printf("#%s: %s(%s): %s", channel, user.DisplayName, user.Username, message.Text)
		})

		bot.OnNewRoomstateMessage(func(channel string, user twitch.User, rawMessage twitch.Message) {
			subMode := ModeUnset

			if readSubMode, ok := rawMessage.Tags["subs-only"]; ok {
				if readSubMode == "1" {
					log.Println("xd")
					subMode = ModeEnabled
				} else {
					subMode = ModeDisabled
				}
			}

			if subMode != ModeUnset {
				if subMode == ModeEnabled {
					log.Printf("Submode enabled")
				} else {
					log.Printf("Submode disabled")

					if bot.Flags.PermaSubMode {
						bot.Say(channel, "Perma sub mode is enabled. A mod can type !suboff to disable perma sub mode")
						bot.Say(channel, ".subscribers")
					}
				}
			}

			log.Printf("#%s: %#v: %#v", channel, user, rawMessage)
		})

		if bot.Name == "snusbot" {
			bot.Join("forsen")
		}

		bot.Join(bot.Name)

		go func(bot *bots.TwitchBot) {
			log.Println("Connecting...")
			err := bot.Connect()
			if err != nil {
				log.Fatal(err)
			}
		}(bot)
	}

	return nil
}

func (a *Application) StartGRPCService() error {
	// Start GRPC Server on :50052
	lis, err := net.Listen("tcp", a.config.GRPCService.Host)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterMessageCheckerServer(s, &pajbotServer{})
	reflection.Register(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve %v", err)
		}
	}()

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
			stmt, err := a.SQL.Session.Prepare(queryF)
			if err != nil {
				return err
			}
			_, err = stmt.Exec(cfg.ChannelID, timeoutData.Data.CreatedByUserID, action, duration, timeoutData.Data.TargetUserID, reason, actionContext)
			if err != nil {
				return err
			}
		}

		sayAction := &pb.Action_SayAction{
			SayAction: &pb.SayAction{
				Message: content,
			},
		}
		messageAction := &pb.MessageAction{}
		messageAction.Actions = append(messageAction.Actions, &pb.Action{Action: sayAction})

		// a.GRPCClient.PerformActions(context.Background(), messageAction)
		return nil
	})

	return nil
}

func (a *Application) StartGRPCClient() error {
	// Connect to GRPC Client on :50051
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return err
	}

	// defer conn.Close()
	a.GRPCClient = pb.NewClientClient(conn)

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
