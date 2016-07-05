package bot

import (
	"fmt"
	"strings"
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
	Quit     chan string
	ReadChan chan common.Msg
	SendChan chan string
	Channel  string
	Redis    *redismanager.RedisManager
	SQL      *sqlmanager.SQLManager
	Twitter  *pbtwitter.Bot
}

// Channel contains data about the channel
type Channel struct {
	Name       string
	BttvEmotes map[string]common.Emote // channel and global emotes
	Online     bool
	Uptime     time.Time // time when stream went live, time when last stream ended if not online
}

/*
A Bot runs in a single channel and reacts according to its
given commands.
*/
type Bot struct {
	Quit    chan string
	Read    chan common.Msg
	Send    chan string
	Channel Channel
	Redis   *redismanager.RedisManager
	SQL     *sqlmanager.SQLManager
	Modules []Module
	Twitter *pbtwitter.Bot
}

// Bots is a map of bots, keyed by the channel
var Bots = make(map[string]*Bot)

/*
NewBot instansiates a new Bot object with the given Config object
*/
func NewBot(cfg Config, modules []Module) *Bot {
	channel := Channel{
		Name:       cfg.Channel,
		BttvEmotes: make(map[string]common.Emote),
	}
	b := &Bot{
		Quit:    cfg.Quit,
		Read:    cfg.ReadChan,
		Send:    cfg.SendChan,
		Channel: channel,
		Modules: modules,
		Redis:   cfg.Redis,
		SQL:     cfg.SQL,
		Twitter: cfg.Twitter,
	}
	Bots[cfg.Channel] = b
	return b
}

/*
Init starts the bots

Connects to irc, joins channel etc
starts channels?
idk
*/
func (bot *Bot) Init() {
	log.Infof("new bot in %s", bot.Channel)
	go bot.LoadBttvEmotes()
	for {
		select {
		case m := <-bot.Read:
			// log.Infof("#%s %s :%s\n", m.Channel, m.User.Name, m.Text)
			if m.Type != common.MsgSub {
				bot.Redis.GetUser(bot.Channel.Name, &m.User)
			}
			log.Debugf("%s is level %d\n", m.User.Name, m.User.Level)
			go bot.Handle(m)
		case tweet := <-bot.Twitter.Stream:
			bot.SaySafef("PogChamp new tweet from %s: %s", tweet.User.Name, tweet.Text)
		}
	}
}

/*
Say sends a PRIVMSG to the bots given channel
*/
func (bot *Bot) Say(message string) {
	m := fmt.Sprintf("PRIVMSG #%s :%s ", bot.Channel.Name, message)
	bot.Send <- m
}

/*
Sayf sends a formatted PRIVMSG to the bots given channel
*/
func (bot *Bot) Sayf(format string, a ...interface{}) {
	bot.Say(fmt.Sprintf(format, a...))
}

// SayFormat sends a formatted and safe message to the bots channel
func (bot *Bot) SayFormat(line string, msg *common.Msg, a ...interface{}) {
	bot.SaySafef(bot.Format(line, msg), a...)
}

/*
SaySafef sends a formatted PRIVMSG to the bots given channel
*/
func (bot *Bot) SaySafef(format string, a ...interface{}) {
	bot.SaySafe(fmt.Sprintf(format, a...))
}

/*
SaySafe allows only harmless irc commands,
this should be used for commands added by users
*/
func (bot *Bot) SaySafe(message string) {
	if !strings.HasPrefix(message, "/") && !strings.HasPrefix(message, ".") {
		bot.Say(message)
		return
	}
	m := strings.Split(message, " ")
	cmd := m[0][1:] // remove "." or "/"
	switch cmd {
	case "me":
	case "timeout":
	case "unban":
	case "subscribers":
	case "subscribersoff":
	case "emoteonly":
	case "emoteonlyoff":
	default:
		message = " " + message
	}
	bot.Say(message)
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
