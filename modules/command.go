package modules

import (
	"strings"

	"github.com/pajlada/pajbot2/bot"
)

/*
Command xD
*/
type Command struct {
}

// Ensure the module implements the interface properly
var _ Module = (*Command)(nil)

// Check xD
func (module *Command) Check(b *bot.Bot, msg *bot.Msg, action *bot.Action) error {
	m := strings.Split(msg.Message, " ")
	trigger := strings.ToLower(m[0])
	if trigger == "!xd" {
		action.Response = "pajaSWA"
		action.Stop = true
	}
	return nil
}
