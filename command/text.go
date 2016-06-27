package command

import (
	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
)

// TextCommand xD
type TextCommand struct {
	BaseCommand
	Response string
}

var _ Command = (*TextCommand)(nil)

/*
IsTriggered returns true if the given string `message` would trigger this command,
otherwise return false
*/
func (command *TextCommand) IsTriggered(t string, fullMessage []string, index int) (bool, Command) {
	for _, trigger := range command.Triggers {
		if trigger == t {
			return true, command
		}
	}

	return false, nil
}

// GetResponse xD
func (command *TextCommand) GetResponse() string {
	// TODO: get $(user) variables and shit
	return command.Response
}

// Run is the method that will decide what this sort of command will do forsenE
func (command *TextCommand) Run(b *bot.Bot, msg *common.Msg, action *bot.Action) string {
	return command.GetResponse()
}
