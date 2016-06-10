package modules

import (
	"strings"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
)

/*
Command xD
*/
type Command struct {
}

// Ensure the module implements the interface properly
var _ Module = (*Command)(nil)

// Check xD
func (module *Command) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	m := strings.Split(msg.Message, " ")
	trigger := strings.ToLower(m[0])
	if trigger == "!xd" {
		action.Response = "pajaSWA"
		action.Stop = true
	}
	if trigger == "!quit" && msg.User.Name == "nuuls" {
		b.Quit <- "ayy lmao something bad happened xD"
	}
	return nil
}
