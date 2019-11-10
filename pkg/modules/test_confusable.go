package modules

import (
	"fmt"
	"strings"

	confusablematcher "github.com/pajbot/pajbot2/internal/ConfusableMatcher-go-interop"
	"github.com/pajbot/pajbot2/internal/commands/base"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commands"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

func init() {
	Register("test_confusable", func() pkg.ModuleSpec {
		return &Spec{
			id:               "test_confusable",
			name:             "Test Confusable",
			enabledByDefault: false,

			maker: newTestConfusable,
		}
	})
}

type testConfusable struct {
	mbase.Base

	commands pkg.CommandsManager

	matcher *confusablematcher.CMHandle
}

func newTestConfusable(b mbase.Base) pkg.Module {
	var inMap []confusablematcher.KeyValue

	xd := confusablematcher.InitConfusableMatcher(inMap, true)

	m := &testConfusable{
		Base: b,

		commands: commands.NewCommands(),

		matcher: &xd,
	}

	m.commands.Register([]string{"!pb2testconfusable"}, newTestConfusableCommand(m.matcher))

	return m
}

func (m *testConfusable) OnMessage(event pkg.MessageEvent) pkg.Actions {
	return m.commands.OnMessage(event)
}

type testConfusableCommand struct {
	base.Command

	matcher *confusablematcher.CMHandle
}

func newTestConfusableCommand(matcher *confusablematcher.CMHandle) pkg.CustomCommand2 {
	return &testConfusableCommand{
		Command: base.New(),
		matcher: matcher,
	}
}

func (c *testConfusableCommand) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !event.User.HasPermission(event.Channel, pkg.PermissionModeration) {
		return nil
	}

	if len(parts) < 3 {
		return nil
	}

	fmt.Println(parts)

	text := strings.Join(parts[1:len(parts)-1], " ")
	fmt.Println(text)
	search := strings.Join(parts[len(parts)-1:], " ")
	fmt.Println(search)

	a, b := confusablematcher.IndexOf(*c.matcher, text, search, true, 0)

	return twitchactions.Mentionf(event.User, "index=%d, len=%d", a, b)
}
