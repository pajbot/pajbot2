package modules

import (
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/commands"
)

type basicCommandsModule struct {
	botChannel pkg.BotChannel

	server *server

	commands pkg.CommandsManager
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

		commands: commands.NewCommands(),
	}
}

func (m *basicCommandsModule) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel

	m.commands.Register([]string{"!pb2ping"}, commands.NewPing())
	m.commands.Register([]string{"!pb2join"}, commands.NewJoin())
	// m.commands.Register([]string{"!pb2leave"}, &commands.Leave{})
	m.commands.Register([]string{"!pb2module"}, commands.NewModule())
	m.commands.Register([]string{"!pb2quit"}, commands.NewQuit())

	m.commands.Register([]string{"!user"}, commands.NewUser())

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

func (m *basicCommandsModule) OnWhisper(bot pkg.BotChannel, source pkg.User, message pkg.Message) error {
	return nil
}

func (m *basicCommandsModule) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	return m.commands.OnMessage(bot, user, message, action)
}
