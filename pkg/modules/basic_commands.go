package modules

import (
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commands"
)

func init() {
	Register("basic_commands", func() pkg.ModuleSpec {
		return &moduleSpec{
			id:               "basic_commands",
			name:             "Basic commands",
			maker:            newBasicCommandsModule,
			enabledByDefault: true,
		}
	})
}

type basicCommandsModule struct {
	base

	commands pkg.CommandsManager
}

func newBasicCommandsModule(b base) pkg.Module {
	m := &basicCommandsModule{
		base: b,

		commands: commands.NewCommands(),
	}

	// FIXME
	m.Initialize()

	return m
}

func (m *basicCommandsModule) Initialize() {
	m.commands.Register([]string{"!pb2ping"}, commands.NewPing())
	m.commands.Register([]string{"!pb2join"}, commands.NewJoin())
	m.commands.Register([]string{"!pb2leave"}, commands.NewLeave())
	m.commands.Register([]string{"!pb2module"}, commands.NewModule())
	m.commands.Register([]string{"!pb2quit"}, commands.NewQuit())

	m.commands.Register([]string{"!user"}, commands.NewUser())
}

func (m *basicCommandsModule) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	return m.commands.OnMessage(bot, user, message, action)
}
