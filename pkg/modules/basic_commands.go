package modules

import (
	"strings"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/commands"
)

type basicCommandsModule struct {
	botChannel pkg.BotChannel

	server *server

	commands map[string]pkg.CustomCommand
}

var basicCommandsModuleSpec = &moduleSpec{
	id:               "basic_commands",
	name:             "Basic commands",
	maker:            newBasicCommandsModule,
	enabledByDefault: true,
}

func newBasicCommandsModule() pkg.Module {
	return &basicCommandsModule{
		server: &_server,

		commands: make(map[string]pkg.CustomCommand),
	}
}

func (m *basicCommandsModule) registerCommand(aliases []string, command pkg.CustomCommand) {
	for _, alias := range aliases {
		m.commands[alias] = command
	}
}

func (m *basicCommandsModule) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel

	m.registerCommand([]string{"!pb2ping"}, &commands.Ping{})
	m.registerCommand([]string{"!pb2join"}, &commands.Join{})
	m.registerCommand([]string{"!pb2leave"}, &commands.Leave{})
	m.registerCommand([]string{"!pb2module"}, commands.NewModule())

	return nil
}

func (m *basicCommandsModule) Disable() error {
	return nil
}

func (m *basicCommandsModule) Spec() pkg.ModuleSpec {
	return basicCommandsModuleSpec
}

func (m *basicCommandsModule) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m *basicCommandsModule) OnWhisper(bot pkg.Sender, source pkg.User, message pkg.Message) error {
	return nil
}

func (m *basicCommandsModule) OnMessage(bot pkg.Sender, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	parts := strings.Split(message.GetText(), " ")
	if len(parts) == 0 {
		return nil
	}

	if command, ok := m.commands[strings.ToLower(parts[0])]; ok {
		command.Trigger(bot, m.botChannel, parts, channel, user, message, action)
	}

	return nil
}
