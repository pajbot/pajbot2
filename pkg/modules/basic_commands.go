package modules

import (
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commands"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
)

func init() {
	Register("basic_commands", func() pkg.ModuleSpec {
		return &Spec{
			id:               "basic_commands",
			name:             "Basic commands",
			maker:            newBasicCommandsModule,
			enabledByDefault: true,
		}
	})
}

type basicCommandsModule struct {
	mbase.Base

	commands pkg.CommandsManager
}

func newBasicCommandsModule(b mbase.Base) pkg.Module {
	m := &basicCommandsModule{
		Base: b,

		commands: commands.NewCommands(),
	}

	// FIXME
	m.Initialize()

	return m
}

func (m *basicCommandsModule) Initialize() {
	m.commands.Register([]string{"!pb2ping"}, commands.NewPing())
	m.commands.Register([]string{"!pb2join"}, commands.NewJoin(m.BotChannel()))
	m.commands.Register([]string{"!pb2leave"}, commands.NewLeave(m.BotChannel()))
	m.commands.Register([]string{"!pb2module"}, commands.NewModule(m.BotChannel()))
	m.commands.Register([]string{"!pb2quit"}, commands.NewQuit(m.BotChannel()))

	m.commands.Register([]string{"!user"}, commands.NewUser(m.BotChannel()))
}

func (m *basicCommandsModule) OnMessage(event pkg.MessageEvent) pkg.Actions {
	return m.commands.OnMessage(event)
}
