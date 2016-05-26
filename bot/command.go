package bot

import (
	"strings"
	"time"
)

/*
Command describes a command object.
A command object will contain a list of trigger words.
The trigger words must be at the beginning of the message to trigger.
*/
type Command struct {
	ID       int
	Trigger  []string
	Cooldown map[string]time.Time
	Response string
	Level    int
}

/*
Handle attempts to handle the given message
*/
func (bot *Bot) Handle(msg Msg) error {
	m := strings.Split(msg.Message, " ")
	trigger := strings.ToLower(m[0])
	if trigger == "!xd" {
		bot.Say("pajaSWA")
	}
	return nil
}
