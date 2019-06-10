package modules

import (
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commands"
)

type CustomCommands struct {
	botChannel pkg.BotChannel

	server *server

	commands pkg.CommandsManager
}

func NewCustomCommands() *CustomCommands {
	return &CustomCommands{
		server: &_server,

		commands: commands.NewCommands(),
	}
}

// func (m *CustomCommands) RegisterCommand(aliases []string, command pkg.CustomCommand) {
// 	for _, alias := range aliases {
// 		m.commands[alias] = command
// 	}
// }

// FIXME
// func (m *CustomCommands) OnWhisper(bot pkg.Sender, source pkg.User, message pkg.Message) error {
// 	return nil
// }
//
// func (m *CustomCommands) OnMessage(bot pkg.Sender, source pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
// 	parts := strings.Split(message.GetText(), " ")
// 	if len(parts) == 0 {
// 		return nil
// 	}
//
// 	if command, ok := m.commands[strings.ToLower(parts[0])]; ok {
// 		command.Trigger(m.botChannel, parts, source, user, message, action)
// 	}
//
// 	return nil
// }
