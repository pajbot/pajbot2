package command

import (
	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
)

// CustomFunction xD
type CustomFunction func(b *bot.Bot, msg *common.Msg, action *bot.Action)

// A FuncCommand is a command that can trigger a function
type FuncCommand struct {
	BaseCommand
	Function CustomFunction
}

var _ Command = (*FuncCommand)(nil)

/*
IsTriggered returns true with the relevant Command if it finds a match
*/
func (command *FuncCommand) IsTriggered(t string, fullMessage []string, index int) (bool, Command) {
	for _, trigger := range command.Triggers {
		// Matched one of our triggers
		if trigger == t {
			return true, command
		}
	}

	return false, nil
}

// Run xD
func (command *FuncCommand) Run(b *bot.Bot, msg *common.Msg, action *bot.Action) string {
	command.Function(b, msg, action)
	return "TEST"
}

// GetBaseCommand returns the BaseCommand of the current command type
func (command *FuncCommand) GetBaseCommand() *BaseCommand {
	return &command.BaseCommand
}
