package module

import "github.com/pajlada/pajbot2/bot"

/*
A Module is the base of every handler for commands.
*/
type Module interface {
	Check(bot *bot.Bot, msg *bot.Msg, action *bot.Action) error
}
