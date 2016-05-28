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

func (bot *Bot) CheckForCommand(msg Msg) Action {
	a := &Action{}
	m := strings.Split(msg.Message, " ")
	trigger := strings.ToLower(m[0])
	if trigger == "!xd" {
		a.Response = "pajaSWA"
	}
	if a.Response != "" {
		a.Match = true
	}
	return *a
}
