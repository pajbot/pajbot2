package modules

import (
	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
)

/*
A Module is the base of every handler for commands.
*/
type Module interface {
	Check(bot *bot.Bot, msg *common.Msg, action *bot.Action) error
}
