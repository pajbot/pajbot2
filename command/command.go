package command

import (
	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
)

// Command is the shared interface for all commands
type Command interface {
	IsTriggered(t string, fullMessage []string, index int) (bool, Command)
	Run(b *bot.Bot, msg *common.Msg, action *bot.Action) string
}
