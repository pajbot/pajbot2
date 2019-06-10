package modules

import (
	"strings"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commands"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
	"github.com/pajbot/pajbot2/pkg/users"
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

func (c pb2Say) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !event.User.HasPermission(event.Channel, pkg.PermissionAdmin) {
		return nil
	}

	if len(parts) < 2 {
		return nil
	}

	return twitchactions.Say(strings.Join(parts[1:], " "))
}

type pb2Whisper struct {
}

func (c pb2Whisper) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !event.User.HasPermission(event.Channel, pkg.PermissionAdmin) {
		return nil
	}

	if len(parts) < 3 {
		return nil
	}

	username := utils.FilterUsername(parts[1])
	if username == "" {
		// Invalid username
		return nil
	}

	targetUser := users.NewSimpleTwitchUser("", username)

	return twitchactions.DoWhisper(targetUser, strings.Join(parts[2:], " "))
}

type debugModule struct {
	base

	commands pkg.CommandsManager
}

func newDebugModule(b base) pkg.Module {
	m := &debugModule{
		base: b,

		commands: commands.NewCommands(),
	}

	// FIXME
	m.Initialize()

	return m
}

func (m *debugModule) Initialize() {
	m.commands.Register([]string{"!pb2say"}, &pb2Say{})
	m.commands.Register([]string{"!pb2whisper"}, &pb2Whisper{})
}

func (m *debugModule) OnWhisper(event pkg.MessageEvent) pkg.Actions {
	return m.commands.OnMessage(event)
}

func (m *debugModule) OnMessage(event pkg.MessageEvent) pkg.Actions {
	return m.commands.OnMessage(event)
}
