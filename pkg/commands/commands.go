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

	match, parts := c.Match(message.GetText())
	return c.trigger(match, parts, event)
}

// OnWhisper forwards a whisper event to commands that opt into whisper handling via the CanExecuteWithWhisper function
// Commands that flow through here will have been pre-parsed by bot/botchannel already, so while the user sends e.g. `#forsen !user xD`, the command system only ever sees `!user xD`
func (c *Commands) OnWhisper(event pkg.MessageEvent) pkg.Actions {
	message := event.Message

	match, parts := c.Match(message.GetText())
	command, ok := match.(pkg.WhisperCommand)
	if !ok || !command.CanExecuteWithWhisper() {
		return nil
	}

	return c.trigger(match, parts, event)
}

func (c *Commands) trigger(match interface{}, parts []string, event pkg.MessageEvent) pkg.Actions {
	user := event.User

	if match != nil {
		switch command := match.(type) {
		case pkg.CustomCommand2:
			if command.HasCooldown(user) {
				return nil
			}
			command.AddCooldown(user)
			return command.Trigger(parts, event)

		case pkg.SimpleCommand:
			return command.Trigger(parts, event)
		}
	}

	return nil
}
