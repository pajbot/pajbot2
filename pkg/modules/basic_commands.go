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
	m.commands.Register([]string{"!pb2join"}, commands.NewJoin(m.bot))
	m.commands.Register([]string{"!pb2leave"}, commands.NewLeave(m.bot))
	m.commands.Register([]string{"!pb2module"}, commands.NewModule(m.bot))
	m.commands.Register([]string{"!pb2quit"}, commands.NewQuit(m.bot))

	m.commands.Register([]string{"!user"}, commands.NewUser(m.bot))
}

func (m *basicCommandsModule) OnMessage(event pkg.MessageEvent) pkg.Actions {
	return m.commands.OnMessage(event)
}
