package commands

import (
	"github.com/pajbot/commandmatcher"
	"github.com/pajbot/pajbot2/pkg"
)

type Commands struct {
	*commandmatcher.CommandMatcher

	internalCommands map[int64]interface{}
}

func NewCommands() *Commands {
	c := &Commands{
		CommandMatcher: commandmatcher.New(),

		internalCommands: map[int64]interface{}{},
	}

	return c
}

func (c *Commands) Register2(id int64, triggers []string, cmd interface{}) {
	c.CommandMatcher.Register(triggers, cmd)
	c.internalCommands[id] = cmd
}

func (c *Commands) FindByCommandID(id int64) interface{} {
	return c.internalCommands[id]
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
