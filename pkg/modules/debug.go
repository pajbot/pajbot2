package modules

import (
	"fmt"
	"strings"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commands"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
	"github.com/pajbot/pajbot2/pkg/users"
	"github.com/pajbot/utils"
)

func init() {
	Register("debug", func() pkg.ModuleSpec {
		return &Spec{
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

var (
	whitelistedCommands = map[string]interface{}{
		"ban":            nil,
		"unban":          nil,
		"timeout":        nil,
		"untimeout":      nil,
		"delete":         nil,
		"subscribers":    nil,
		"subscribersoff": nil,
		"r9kbeta":        nil,
		"r9kbetaoff":     nil,
		"emoteonly":      nil,
		"emoteonlyoff":   nil,
		"slow":           nil,
		"slowoff":        nil,
		"followers":      nil,
		"followersoff":   nil,
		"clear":          nil,
	}
)

type pb2Exec struct {
}

func (c *pb2Exec) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !event.User.HasPermission(event.Channel, pkg.PermissionModeration) {
		return nil
	}

	if len(parts) < 2 {
		return nil
	}

	command := parts[1]
	if !strings.HasPrefix(command, ".") && !strings.HasPrefix(command, "/") {
		command = "." + command
	}

	if _, ok := whitelistedCommands[command[1:]]; !ok {
		return twitchactions.DoWhisperf(event.User, "You are not allowed to run the command '%s'", command[1:])
	}

	fmt.Printf("pb2exec Executing command for %s: %s\n", event.User.GetName(), strings.Join(parts[1:], " "))

	return twitchactions.Sayf("%s %s", command, strings.Join(parts[2:], " "))
}

type debugModule struct {
	mbase.Base

	commands pkg.CommandsManager
}

func newDebugModule(b mbase.Base) pkg.Module {
	m := &debugModule{
		Base: b,

		commands: commands.NewCommands(),
	}

	// FIXME
	m.Initialize()

	return m
}

func (m *debugModule) Initialize() {
	m.commands.Register([]string{"!pb2say"}, &pb2Say{})
	m.commands.Register([]string{"!pb2whisper"}, &pb2Whisper{})
	m.commands.Register([]string{"!pb2exec"}, &pb2Exec{})
}

func (m *debugModule) OnWhisper(event pkg.MessageEvent) pkg.Actions {
	return m.commands.OnMessage(event)
}

func (m *debugModule) OnMessage(event pkg.MessageEvent) pkg.Actions {
	return m.commands.OnMessage(event)
}
