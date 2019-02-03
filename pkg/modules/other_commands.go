package modules

import (
	"strings"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/commands"
)

type otherCommandsModule struct {
	botChannel pkg.BotChannel

	server *server

	commands map[string]pkg.CustomCommand
}

var otherCommandsModuleSpec = &moduleSpec{
	id:               "other_commands",
	name:             "Other commands",
	maker:            newOtherCommandsModule,
	enabledByDefault: false,
}

func newOtherCommandsModule() pkg.Module {
	return &otherCommandsModule{
		server: &_server,

		commands: make(map[string]pkg.CustomCommand),
	}
}

func (m *otherCommandsModule) registerCommand(aliases []string, command pkg.CustomCommand) {
	for _, alias := range aliases {
		m.commands[alias] = command
	}
}

func (m *otherCommandsModule) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel

	m.registerCommand([]string{"!userid"}, &commands.GetUserID{})
	m.registerCommand([]string{"!username"}, &commands.GetUserName{})
	m.registerCommand([]string{"!pb2points"}, &commands.GetPoints{})
	m.registerCommand([]string{"!pb2roulette"}, &commands.Roulette{})
	m.registerCommand([]string{"!pb2givepoints"}, &commands.GivePoints{})
	// m.registerCommand([]string{"!pb2addpoints"}, &commands.AddPoints{})
	// m.registerCommand([]string{"!pb2removepoints"}, &commands.RemovePoints{})
	m.registerCommand([]string{"!roffle", "!join"}, commands.NewRaffle())
	m.registerCommand([]string{"!user"}, commands.NewUser())
	m.registerCommand([]string{"!pb2rank"}, &commands.Rank{})
	m.registerCommand([]string{"!pb2simplify"}, &commands.Simplify{})
	// m.registerCommand([]string{"!timemeout"}, &commands.TimeMeOut{})
	m.registerCommand([]string{"!pb2test"}, &commands.Test{})
	m.registerCommand([]string{"!pb2islive"}, commands.IsLive{})

	return nil
}

func (m *otherCommandsModule) Disable() error {
	return nil
}

func (m *otherCommandsModule) Spec() pkg.ModuleSpec {
	return otherCommandsModuleSpec
}

func (m *otherCommandsModule) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m *otherCommandsModule) OnWhisper(bot pkg.BotChannel, source pkg.User, message pkg.Message) error {
	return nil
}

func (m *otherCommandsModule) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	parts := strings.Split(message.GetText(), " ")
	if len(parts) == 0 {
		return nil
	}

	if command, ok := m.commands[strings.ToLower(parts[0])]; ok {
		command.Trigger(m.botChannel, parts, bot.Channel(), user, message, action)
	}

	return nil
}
