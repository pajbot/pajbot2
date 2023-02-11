package modules

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

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

type pb2Exec struct {
	m *debugModule
}

func (c *pb2Exec) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !event.User.HasPermission(event.Channel, pkg.PermissionModeration) {
		return nil
	}

	if len(parts) < 2 {
		return nil
	}

	command := strings.TrimLeftFunc(strings.ToLower(parts[1]), func(c rune) bool {
		if c == '.' || c == '/' {
			return true
		}
		return false
	})

	getTarget := func() (pkg.User, error) {
		if len(parts) <= 2 {
			return nil, errors.New("no target specified")
		}

		login := strings.ToLower(parts[2])

		return event.UserStore.GetUserByLogin(login)
	}

	parseDuration := func(t string, defaultUnit time.Duration, defaultTime time.Duration) time.Duration {
		if seconds, err := strconv.Atoi(t); err == nil {
			// No suffix = treat as seconds
			return time.Duration(time.Duration(seconds) * defaultUnit)
		}

		d, err := time.ParseDuration(t)
		if err != nil {
			return defaultTime
		}
		return d
	}

	actions := &twitchactions.Actions{}

	switch command {
	case "ban":
		target, err := getTarget()
		if err != nil {
			return twitchactions.DoWhisperf(event.User, "missing target: %s", err)
		}
		reason := strings.Join(parts[3:], " ")
		actions.Ban(target).SetReason(reason)

	case "unban":
		target, err := getTarget()
		if err != nil {
			return twitchactions.DoWhisperf(event.User, "missing target: %s", err)
		}
		actions.Unban(target)

	case "timeout":
		target, err := getTarget()
		if err != nil {
			return twitchactions.DoWhisperf(event.User, "missing target: %s", err)
		}

		var duration time.Duration
		reason := ""
		if len(parts) > 3 {
			// default timeout duration is 10m
			duration = parseDuration(parts[3], time.Second, 10*time.Minute)
			if len(parts) > 4 {
				reason = strings.Join(parts[4:], " ")
			}
		} else {
			duration = 10 * time.Minute
		}

		actions.Timeout(target, duration).SetReason(reason)

	case "untimeout":
		target, err := getTarget()
		if err != nil {
			return twitchactions.DoWhisperf(event.User, "missing target: %s", err)
		}
		actions.Untimeout(target)

	case "delete":
		if len(parts) <= 2 {
			return twitchactions.DoWhisperf(event.User, "missing msg id")
		}

		actions.Delete(parts[2])

	case "subscribers":
		if err := c.m.BotChannel().SetSubscribers(true); err != nil {
			return twitchactions.DoWhisperf(event.User, "Error enabling subscribers mode '%s'", err)
		}

	case "subscribersoff":
		if err := c.m.BotChannel().SetSubscribers(false); err != nil {
			return twitchactions.DoWhisperf(event.User, "Error disabling subscribers mode '%s'", err)
		}

	case "r9kbeta":
		fallthrough
	case "uniquechat":
		if err := c.m.BotChannel().SetUniqueChat(true); err != nil {
			return twitchactions.DoWhisperf(event.User, "Error enabling unique chat '%s'", err)
		}

	case "r9kbetaoff":
		fallthrough
	case "uniquechatoff":
		if err := c.m.BotChannel().SetUniqueChat(false); err != nil {
			return twitchactions.DoWhisperf(event.User, "Error disabling unique chat '%s'", err)
		}

	case "nonmodchatdelay":
		duration := 30 * time.Second
		if len(parts) > 2 {
			duration = parseDuration(parts[2], time.Second, 30*time.Second)
		}

		if err := c.m.BotChannel().SetNonModChatDelay(true, int(duration.Seconds())); err != nil {
			return twitchactions.DoWhisperf(event.User, "Error enabling slow mode '%s'", err)
		}

	case "nonmodchatdelayoff":
		if err := c.m.BotChannel().SetNonModChatDelay(false, 0); err != nil {
			return twitchactions.DoWhisperf(event.User, "Error enabling slow mode '%s'", err)
		}

	case "emoteonly":
		if err := c.m.BotChannel().SetEmoteOnly(true); err != nil {
			return twitchactions.DoWhisperf(event.User, "Error enabling emote only '%s'", err)
		}

	case "emoteonlyoff":
		if err := c.m.BotChannel().SetEmoteOnly(false); err != nil {
			return twitchactions.DoWhisperf(event.User, "Error disabling emote only '%s'", err)
		}

	case "slow":
		duration := 30 * time.Second
		if len(parts) > 2 {
			duration = parseDuration(parts[2], time.Second, 30*time.Second)
		}

		if err := c.m.BotChannel().SetSlowMode(true, int(duration.Seconds())); err != nil {
			return twitchactions.DoWhisperf(event.User, "Error enabling slow mode '%s'", err)
		}

	case "slowoff":
		if err := c.m.BotChannel().SetSlowMode(false, 0); err != nil {
			return twitchactions.DoWhisperf(event.User, "Error disabling slow mode '%s'", err)
		}

	case "followers":
		duration := 0 * time.Second
		if len(parts) > 2 {
			duration = parseDuration(parts[2], time.Minute, 0*time.Second)
		}

		if err := c.m.BotChannel().SetFollowerMode(true, int(math.Floor(duration.Minutes()))); err != nil {
			return twitchactions.DoWhisperf(event.User, "Error enabling follower mode '%s'", err)
		}

	case "followersoff":
		if err := c.m.BotChannel().SetFollowerMode(false, 0); err != nil {
			return twitchactions.DoWhisperf(event.User, "Error disabling follower mode '%s'", err)
		}

	case "clear":
		return twitchactions.DoWhisper(event.User, "clear is no longer implemented, poke pajlada why you would need this")

	default:
		return twitchactions.DoWhisperf(event.User, "You are not allowed to run the command '%s'", command)
	}

	fmt.Printf("pb2exec Executing command for %s: %s\n", event.User.GetName(), strings.Join(parts[1:], " "))

	return actions
}

type debugModule struct {
	mbase.Base

	commands pkg.CommandsManager
}

func newDebugModule(b *mbase.Base) pkg.Module {
	m := &debugModule{
		Base: *b,

		commands: commands.NewCommands(),
	}

	// FIXME
	m.Initialize()

	return m
}

func (m *debugModule) Initialize() {
	m.commands.Register([]string{"!pb2say"}, &pb2Say{})
	m.commands.Register([]string{"!pb2whisper"}, &pb2Whisper{})
	m.commands.Register([]string{"!pb2exec"}, &pb2Exec{m})
}

func (m *debugModule) OnWhisper(event pkg.MessageEvent) pkg.Actions {
	return m.commands.OnMessage(event)
}

func (m *debugModule) OnMessage(event pkg.MessageEvent) pkg.Actions {
	return m.commands.OnMessage(event)
}
