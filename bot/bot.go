package bot

import (
	"fmt"
	"time"

	"github.com/pajlada/pajbot2/redismanager"

	"github.com/pajlada/pajbot2/common"
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
}

/*
A Bot runs in a single channel and reacts according to its
given commands.
*/
type Bot struct {
	Quit    chan string
	Read    chan common.Msg
	Send    chan string
	Channel string
	Redis   *redismanager.RedisManager
	Modules []Module
}

/*
NewBot instansiates a new Bot object with the given Config object
*/
func NewBot(cfg Config, modules []Module) *Bot {
	return &Bot{
		Quit:    cfg.Quit,
		Read:    cfg.ReadChan,
		Send:    cfg.SendChan,
		Channel: cfg.Channel,
		Modules: modules,
		Redis:   cfg.Redis,
	}
}

/*
Init starts the bots

Connects to irc, joins channel etc
starts channels?
idk
*/
func (bot *Bot) Init() {
	log.Infof("new bot in %s", bot.Channel)
	for {
		m := <-bot.Read
		// log.Infof("#%s %s :%s\n", m.Channel, m.User.Name, m.Text)
		if m.Type == common.MsgSub {
			log.Infof("%s subbed for %d months in a row\n", m.User.Name, m.Length)
		}
		if m.Type != common.MsgSub {
			bot.Redis.GetUser(bot.Channel, &m.User)
		}
		log.Debugf("%s is level %d\n", m.User.Name, m.User.Level)
		go bot.Handle(m)
	}
}

/*
Say sends a PRIVMSG to the bots given channel
*/
func (bot *Bot) Say(message string) {
	m := fmt.Sprintf("PRIVMSG #%s :%s ", bot.Channel, message)
	bot.Send <- m
}

/*
Sayf sends a formatted PRIVMSG to the bots given channel
*/
func (bot *Bot) Sayf(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	m := fmt.Sprintf("PRIVMSG #%s :%s ", bot.Channel, message)
	bot.Send <- m
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
