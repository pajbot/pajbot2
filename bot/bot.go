package bot

import (
	"fmt"
	"log"

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
	Redis    *redismanager.Redismanager
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
	Redis   *redismanager.Redismanager
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
	fmt.Printf("new bot in %s\n", bot.Channel)
	for {
		m := <-bot.Read
		fmt.Printf("#%s %s :%s\n", m.Channel, m.User.Name, m.Message)
		if m.Type == common.MsgSub {
			fmt.Printf("%s subbed for %d months in a row\n", m.User.Name, m.Length)
		}
		if m.Type != common.MsgSub {
			bot.Redis.UpdateUser(bot.Channel, &m.User)
			bot.Redis.GetUser(bot.Channel, &m.User)
		}
		log.Printf("%s is level %d\n", m.User.Name, m.User.Level)
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
Ban returns the ban/TO message and reason, perm ban if duration < 0
*/
func (bot *Bot) Ban(user string, duration int, reason string) string {
	if duration < 0 {
		return fmt.Sprintf(".ban %s %s - autoban by pajbot", user, reason)
	}
	return fmt.Sprintf(".timeout %s %d %s - autoban by pajbot", user, duration, reason)
}
