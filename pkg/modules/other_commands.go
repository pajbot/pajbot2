package modules

import (
	"github.com/pajbot/pajbot2/internal/commands/getuserid"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commands"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
)

func init() {
	Register("other_commands", func() pkg.ModuleSpec {
		return &Spec{
			id:               "other_commands",
			name:             "Other commands",
			maker:            newOtherCommandsModule,
			enabledByDefault: false,
		}
	})
}

type otherCommandsModule struct {
	mbase.Base

	commands pkg.CommandsManager
}

func newOtherCommandsModule(b *mbase.Base) pkg.Module {
	m := &otherCommandsModule{
		Base: *b,

		commands: commands.NewCommands(),
	}

	// FIXME
	m.Initialize()

	return m
}

func (m *otherCommandsModule) Initialize() {
	m.commands.Register([]string{"!userid"}, &getuserid.Command{})
	m.commands.Register([]string{"!username"}, &commands.GetUserName{})
	m.commands.Register([]string{"!pb2points"}, &commands.GetPoints{})
	m.commands.Register([]string{"!pb2roulette"}, &commands.Roulette{})
	m.commands.Register([]string{"!pb2givepoints"}, &commands.GivePoints{})
	// m.commands.Register([]string{"!pb2addpoints"}, &commands.AddPoints{})
	// m.commands.Register([]string{"!pb2removepoints"}, &commands.RemovePoints{})
	m.commands.Register([]string{"!roffle", "!join"}, commands.NewRaffle())
	m.commands.Register([]string{"!pb2rank"}, &commands.Rank{})
	m.commands.Register([]string{"!pb2simplify"}, &commands.Simplify{})
	// m.commands.Register([]string{"!timemeout"}, &commands.TimeMeOut{})
	m.commands.Register([]string{"!pb2test"}, &commands.Test{})
	m.commands.Register([]string{"!pb2islive"}, commands.IsLive{})
}

func (m *otherCommandsModule) OnMessage(event pkg.MessageEvent) pkg.Actions {
	return m.commands.OnMessage(event)
}
