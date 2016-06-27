package command

import (
	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
)

// A NestedCommand is a command that can trigger another command
type NestedCommand struct {
	BaseCommand
	Commands []Command
	/* DefaultCommand is the command that would be called if no further argument
	is sent */
	DefaultCommand Command
	/* FallbackCommand is the command that would be called if an argument is sent
	but it matches none of the available Commands */
	FallbackCommand Command
}

var _ Command = (*NestedCommand)(nil)

/*
IsTriggered returns true with the relevant Command if it finds a match
*/
func (command *NestedCommand) IsTriggered(t string, fullMessage []string, index int) (bool, Command) {
	for _, trigger := range command.Triggers {
		// Matched one of our triggers
		if trigger == t {
			// Are there any further arguments?
			if len(fullMessage) > index+1 {
				t = fullMessage[index+1]
				// Check through our Commands
				for _, command := range command.Commands {
					// One of our Commands were triggered, continue with index+1
					if triggered, c := command.IsTriggered(t, fullMessage, index+1); triggered {
						return true, c
					}
				}
				/* A further argument was sent, but it didn't match any of our
				Commands. Fall back */
				if command.FallbackCommand != nil {
					return true, command.FallbackCommand
				}
			} else if command.DefaultCommand != nil {
				// No further argument was sent, use the Default Command if available
				return true, command.DefaultCommand
			}

			return false, nil
		}
	}

	return false, nil
}

// Run xD
func (command *NestedCommand) Run(b *bot.Bot, msg *common.Msg, action *bot.Action) string {
	return ""
}

// GetBaseCommand returns the BaseCommand of the current command type
func (command *NestedCommand) GetBaseCommand() *BaseCommand {
	return &command.BaseCommand
}
