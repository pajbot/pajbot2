package modules

import (
	"strings"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commands"
)

func init() {
	Register("other_commands", func() pkg.ModuleSpec {
		return &moduleSpec{
			id:               "other_commands",
			name:             "Other commands",
			maker:            newOtherCommandsModule,
			enabledByDefault: false,
		}
	})
}

type otherCommandsModule struct {
	base

	commands map[string]pkg.CustomCommand
}

func newOtherCommandsModule(b base) pkg.Module {
	m := &otherCommandsModule{
		base: b,

		commands: make(map[string]pkg.CustomCommand),
	}

	// FIXME
	m.Initialize()

	return m
}

func (m *otherCommandsModule) registerCommand(aliases []string, command pkg.CustomCommand) {
	for _, alias := range aliases {
		m.commands[alias] = command
	}
}

func (m *otherCommandsModule) Initialize() {
	m.registerCommand([]string{"!userid"}, &commands.GetUserID{})
	m.registerCommand([]string{"!username"}, &commands.GetUserName{})
	m.registerCommand([]string{"!pb2points"}, &commands.GetPoints{})
	m.registerCommand([]string{"!pb2roulette"}, &commands.Roulette{})
	m.registerCommand([]string{"!pb2givepoints"}, &commands.GivePoints{})
	// m.registerCommand([]string{"!pb2addpoints"}, &commands.AddPoints{})
	// m.registerCommand([]string{"!pb2removepoints"}, &commands.RemovePoints{})
	m.registerCommand([]string{"!roffle", "!join"}, commands.NewRaffle())
	m.registerCommand([]string{"!pb2rank"}, &commands.Rank{})
	m.registerCommand([]string{"!pb2simplify"}, &commands.Simplify{})
	// m.registerCommand([]string{"!timemeout"}, &commands.TimeMeOut{})
	m.registerCommand([]string{"!pb2test"}, &commands.Test{})
	m.registerCommand([]string{"!pb2islive"}, commands.IsLive{})
}

func (m *otherCommandsModule) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	parts := strings.Split(message.GetText(), " ")
	if len(parts) == 0 {
		return nil
	}

	if command, ok := m.commands[strings.ToLower(parts[0])]; ok {
		command.Trigger(m.bot, parts, bot.Channel(), user, message, action)
	}

	return nil
}
