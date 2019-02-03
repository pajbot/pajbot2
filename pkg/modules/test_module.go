package modules

import (
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/eventemitter"
)

func init() {
	Register(testSpec)
}

type test struct {
	botChannel pkg.BotChannel

	server      *server
	connections []*eventemitter.Listener
}

func newTest() pkg.Module {
	return &test{
		server: &_server,
	}
}

var testSpec = &moduleSpec{
	id:    "test",
	name:  "Test",
	maker: newTest,

	enabledByDefault: false,
}

func (m *test) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel
	return nil
}

func (m *test) Disable() error {
	for _, c := range m.connections {
		c.Disconnected = true
	}
	return nil
}

func (m *test) Spec() pkg.ModuleSpec {
	return testSpec
}

func (m *test) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m test) OnWhisper(bot pkg.BotChannel, user pkg.User, message pkg.Message) error {
	return nil
}

func (m test) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	bot.Mention(user, "test module xd")
	return nil
}
