package commands

import (
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commandmatcher"
)

type Commands struct {
	*commandmatcher.CommandMatcher
}

func NewCommands() *Commands {
	c := &Commands{
		CommandMatcher: commandmatcher.NewMatcher(),
	}

	return c
}

func (c *Commands) OnMessage(event pkg.MessageEvent) pkg.Actions {
	message := event.Message
	user := event.User

	match, parts := c.Match(message.GetText())
	if match != nil {
		switch command := match.(type) {
		case pkg.SimpleCommand:
			return command.Trigger(parts, event)

		case pkg.CustomCommand2:
			if command.HasCooldown(user) {
				return nil
			}
			command.AddCooldown(user)
			return command.Trigger(parts, event)
		}
	}

	return nil
}
