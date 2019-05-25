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

func (c *Commands) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	match, parts := c.Match(message.GetText())
	if match != nil {
		command := match.(pkg.CustomCommand2)
		if command.HasCooldown(user) {
			return nil
		}

		command.Trigger(bot, parts, user, message, action)

		command.AddCooldown(user)
	}

	return nil
}
