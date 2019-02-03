package modules

import (
	"strings"

	"github.com/pajlada/pajbot2/pkg"
)

type pb2Say struct {
}

func (c pb2Say) Trigger(botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	if !user.HasPermission(botChannel.Channel(), pkg.PermissionAdmin) {
		return
	}

	botChannel.Say(strings.Join(parts[1:], " "))
}

func init() {
	Register(debugModuleSpec)
}

type debugModule struct {
	botChannel pkg.BotChannel

	server *server

	commands map[string]pkg.CustomCommand
}

var debugModuleSpec = &moduleSpec{
	id:               "debug",
	name:             "Debug",
	maker:            newDebugModule,
	enabledByDefault: true,
}

func newDebugModule() pkg.Module {
	return &debugModule{
		server: &_server,

		commands: make(map[string]pkg.CustomCommand),
	}
}

func (m *debugModule) registerCommand(aliases []string, command pkg.CustomCommand) {
	for _, alias := range aliases {
		m.commands[alias] = command
	}
}

func (m *debugModule) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel

	m.registerCommand([]string{"!pb2say"}, &pb2Say{})

	return nil
}

func (m *debugModule) Disable() error {
	return nil
}

func (m *debugModule) Spec() pkg.ModuleSpec {
	return debugModuleSpec
}

func (m *debugModule) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m *debugModule) OnWhisper(bot pkg.BotChannel, source pkg.User, message pkg.Message) error {
	return nil
}

func (m *debugModule) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	parts := strings.Split(message.GetText(), " ")
	if len(parts) == 0 {
		return nil
	}

	if command, ok := m.commands[strings.ToLower(parts[0])]; ok {
		command.Trigger(m.botChannel, parts, bot.Channel(), user, message, action)
	}

	return nil
}
