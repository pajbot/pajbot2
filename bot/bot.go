package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/pbtwitter"
	"github.com/pajlada/pajbot2/redismanager"
	"github.com/pajlada/pajbot2/sqlmanager"
)

/*
Config contains some data about
*/
type Config struct {
	BotAccountID int
	Quit         chan string
	ReadChan     chan common.Msg
	SendChan     chan string
	RawReadChan  chan string
	Join         chan string
	Channel      string
	Redis        *redismanager.RedisManager
	SQL          *sqlmanager.SQLManager
	Twitter      *pbtwitter.Bot
}

/*
A Bot runs in a single channel and reacts according to its
given commands.
*/
type Bot struct {
	BotAccountID int
	Quit         chan string
	Read         chan common.Msg
	Send         chan string

	// IRC Raw read channel in case we want to parse our own messages
	RawRead chan string

	Join    chan string
	Channel common.Channel
	Redis   *redismanager.RedisManager
	SQL     *sqlmanager.SQLManager
	Twitter *pbtwitter.Bot

	// List of all available modules
	AllModules []Module

	// List of enabled modules
	EnabledModules []Module
}

// Bots is a map of bots, keyed by the channel
var Bots = make(map[string]*Bot)

/*
NewBot instansiates a new Bot object with the given Config object
*/
func NewBot(cfg Config) *Bot {
	channel := common.Channel{
		Name: cfg.Channel,
	}
	b := &Bot{
		Quit:    cfg.Quit,
		Read:    cfg.ReadChan,
		Send:    cfg.SendChan,
		RawRead: cfg.RawReadChan,
		Join:    cfg.Join,
		Channel: channel,
		Redis:   cfg.Redis,
		SQL:     cfg.SQL,
		Twitter: cfg.Twitter,
	}
	Bots[cfg.Channel] = b

	return b
}

/*
Init starts the bots

 * Load channel-specific BTTV emotes
 * Load channel-specific FFZ emotes
Connects to irc, joins channel etc
starts channels?
idk
*/
func (bot *Bot) Init() {
	log.Printf("new bot in %s", bot.Channel)
	go bot.LoadBttvEmotes()
	go bot.LoadFFZEmotes()
	go bot.readChat()
	go bot.readTweets()

	bot.Sayf("Joined channel %s, build time: %s", bot.Channel.Name, common.BuildTime)
}

func (bot *Bot) readChat() {
	for m := range bot.Read {
		if m.Type != common.MsgSub {
			bot.Redis.GetUser(bot.Channel.Name, &m.User)
		}
		go bot.Handle(m)
	}
}

func (bot *Bot) readTweets() {
	for tweet := range bot.Twitter.Stream {
		bot.SaySafef("PogChamp new tweet from %s (@%s): %s", tweet.User.Name, tweet.User.ScreenName, tweet.Text)
	}
}

// Timeout sends 2 timeouts with a 500 ms delay
func (bot *Bot) Timeout(user string, duration int, reason string) {
	if duration == 0 {
		return
	}
	m := bot.Ban(user, duration, reason)
	bot.Say(m)
	go func() {
		time.Sleep(500 * time.Millisecond)
		bot.Say(m)
	}()
}

/*
Ban returns the ban/TO message and reason, perm ban if duration < 0
*/
func (bot *Bot) Ban(user string, duration int, reason string) string {
	if duration == -1 {
		return fmt.Sprintf(".ban %s %s - autoban by pajbot", user, reason)
	}
	return fmt.Sprintf(".timeout %s %d %s - autoban by pajbot", user, duration, reason)
}
