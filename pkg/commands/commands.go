package commands

import (
	"strings"
	"sync"

	"github.com/pajlada/pajbot2/pkg"
)

type Commands struct {
	commandsMutex sync.Mutex
	commands      map[string]pkg.CustomCommand2
}

func NewCommands() *Commands {
	c := &Commands{
		commands: make(map[string]pkg.CustomCommand2),
	}

	// c.commands["!xdping"] = newPing()

	return c
}

func (c *Commands) Register(aliases []string, command pkg.CustomCommand2) pkg.CustomCommand2 {
	c.commandsMutex.Lock()
	defer c.commandsMutex.Unlock()

	for _, alias := range aliases {
		c.commands[alias] = command
	}

	return command
}

func (c *Commands) Deregister(command pkg.CustomCommand2) {
	c.commandsMutex.Lock()
	defer c.commandsMutex.Unlock()

	var aliasesToRemove []string
	for alias, cmd := range c.commands {
		if cmd == command {
			aliasesToRemove = append(aliasesToRemove, alias)
		}
	}

	for _, alias := range aliasesToRemove {
		delete(c.commands, alias)
	}
}

func (c *Commands) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	parts := strings.Split(message.GetText(), " ")
	if len(parts) == 0 {
		return nil
	}

	if command, ok := c.commands[strings.ToLower(parts[0])]; ok {
		if command.HasCooldown(user) {
			return nil
		}

		command.Trigger(bot, parts, user, message, action)

		command.AddCooldown(user)
	}

	return nil
}
