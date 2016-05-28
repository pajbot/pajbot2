package bot

import (
	"fmt"
	"strings"
)

/*
Config contains some data about
*/
type Config struct {
	ReadChan chan Msg
	SendChan chan string
	Channel  string
}

/*
Msg contains all the information about an IRC message.
This included already-parsed ircv3-tags

XXX: Should this contain the user object?
yes we should. That way we can throw away the Subscriber/Mod/Turbo thing from here
*/
type Msg struct {
	Color       string
	Displayname string
	Emotes      []Emote
	Mod         bool
	Subscriber  bool
	Turbo       bool
	Usertype    string
	Username    string
	Channel     string
	Message     string
	MessageType string
	Me          bool
	Length      int
}

/*
A Bot runs in a single channel and reacts according to its
given commands.
*/
type Bot struct {
	Read    chan Msg
	Send    chan string
	Channel string
}

/*
NewBot instansiates a new Bot object with the given Config object
*/
func NewBot(cfg Config) {
	bot := &Bot{
		Read:    cfg.ReadChan,
		Send:    cfg.SendChan,
		Channel: cfg.Channel,
	}
	bot.init()
}

func (bot *Bot) init() {
	fmt.Printf("new bot in %s\n", bot.Channel)
	for {
		m := <-bot.Read
		fmt.Printf("#%s %s :%s\n", m.Channel, m.Username, m.Message)
		if m.MessageType == "sub" {
			fmt.Printf("%s subbed for %d months in a row\n", m.Username, m.Length)
		}
		go bot.Handle(m)
	}
}

/*
Say sends a PRIVMSG to the bots given channel
*/
func (bot *Bot) Say(message string) {
	if !strings.HasPrefix(message, ".") {
		message = ". " + message
	}
	m := fmt.Sprintf("PRIVMSG #%s :%s", bot.Channel, message)
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
