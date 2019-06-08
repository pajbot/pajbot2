package modules

import (
	"strings"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/utils"
)

func init() {
	Register("debug", func() pkg.ModuleSpec {
		return &moduleSpec{
			id:               "debug",
			name:             "Debug",
			maker:            newDebugModule,
			enabledByDefault: true,
		}
	})
}

type pb2Say struct {
}

func (c pb2Say) Trigger(botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	if !user.HasPermission(botChannel.Channel(), pkg.PermissionAdmin) {
		return
	}

	if len(parts) < 2 {
		return
	}

	botChannel.Say(strings.Join(parts[1:], " "))
}

type pb2Whisper struct {
}

func (c pb2Whisper) Trigger(botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	if !user.HasPermission(botChannel.Channel(), pkg.PermissionAdmin) {
		return
	}

	if len(parts) < 3 {
		return
	}

	username := utils.FilterUsername(parts[1])
	if username == "" {
		// Invalid username
		return
	}

	botChannel.Bot().Whisper(botChannel.Bot().MakeUser(username), strings.Join(parts[2:], " "))
}

type debugModule struct {
	base

	commands map[string]pkg.CustomCommand
}

func newDebugModule(b base) pkg.Module {
	m := &debugModule{
		base: b,

		commands: make(map[string]pkg.CustomCommand),
	}

	// FIXME
	m.Initialize()

	return m
}

func (m *debugModule) registerCommand(aliases []string, command pkg.CustomCommand) {
	for _, alias := range aliases {
		m.commands[alias] = command
	}
}

func (m *debugModule) Initialize() {
	m.registerCommand([]string{"!pb2say"}, &pb2Say{})
	m.registerCommand([]string{"!pb2whisper"}, &pb2Whisper{})
}

func (m *debugModule) OnWhisper(bot pkg.BotChannel, user pkg.User, message pkg.Message) error {
	parts := strings.Split(message.GetText(), " ")
	if len(parts) == 0 {
		return nil
	}

	if command, ok := m.commands[strings.ToLower(parts[0])]; ok {
		command.Trigger(m.bot, parts, bot.Channel(), user, message, nil)
	}

	return nil
}

func (m *debugModule) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	parts := strings.Split(message.GetText(), " ")
	if len(parts) == 0 {
		return nil
	}

	if command, ok := m.commands[strings.ToLower(parts[0])]; ok {
		command.Trigger(m.bot, parts, bot.Channel(), user, message, action)
	}

	return nil
}
