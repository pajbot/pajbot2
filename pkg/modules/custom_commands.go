package modules

import (
	"strings"

	"github.com/pajlada/pajbot2/pkg"
)

type CustomCommands struct {
	botChannel pkg.BotChannel

	server *server

	commands map[string]pkg.CustomCommand
}

func NewCustomCommands() *CustomCommands {
	return &CustomCommands{
		server: &_server,

		commands: make(map[string]pkg.CustomCommand),
	}
}

func (m *CustomCommands) RegisterCommand(aliases []string, command pkg.CustomCommand) {
	for _, alias := range aliases {
		m.commands[alias] = command
	}
}

func (m *CustomCommands) OnWhisper(bot pkg.Sender, source pkg.User, message pkg.Message) error {
	return nil
}

func (m *CustomCommands) OnMessage(bot pkg.Sender, source pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	parts := strings.Split(message.GetText(), " ")
	if len(parts) == 0 {
		return nil
	}

	if command, ok := m.commands[strings.ToLower(parts[0])]; ok {
		command.Trigger(m.botChannel, parts, source, user, message, action)
	}

	return nil
}
