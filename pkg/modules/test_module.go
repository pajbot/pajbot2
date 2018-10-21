package modules

import "github.com/pajlada/pajbot2/pkg"

type test struct {
	botChannel pkg.BotChannel

	server *server
}

func newTest() pkg.Module {
	return &test{
		server: &_server,
	}
}

var testSpec = moduleSpec{
	id:    "test",
	name:  "Test",
	maker: newTest,

	Priority: 0,

	enabledByDefault: false,
}

func (m *test) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel
	return nil
}

func (m *test) Disable() error {
	return nil
}

func (m *test) Spec() pkg.ModuleSpec {
	return &testSpec
}

func (m *test) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m test) OnWhisper(bot pkg.Sender, user pkg.User, message pkg.Message) error {
	return nil
}

func (m test) OnMessage(bot pkg.Sender, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	bot.Mention(channel, user, "test module xd")
	return nil
}
