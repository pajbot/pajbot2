package modules

import (
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/eventemitter"
)

func init() {
	Register(twitterSpec)
}

type twitter struct {
	botChannel pkg.BotChannel

	server      *server
	connections []*eventemitter.Listener
}

func newTwitter() pkg.Module {
	return &twitter{
		server: &_server,
	}
}

var twitterSpec = &moduleSpec{
	id:    "twitter",
	name:  "Twitter",
	maker: newTwitter,

	enabledByDefault: false,
}

func (m *twitter) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel
	recv := m.botChannel.Bot().Application().MIMO().Subscriber("twitter")
	go func() {
		for raw := range recv {
			message := raw.(string)
			m.botChannel.Say("got tweet: " + message)
		}
	}()
	return nil
}

func (m *twitter) Disable() error {
	for _, c := range m.connections {
		c.Disconnected = true
	}
	return nil
}

func (m *twitter) Spec() pkg.ModuleSpec {
	return twitterSpec
}

func (m *twitter) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m twitter) OnWhisper(bot pkg.BotChannel, user pkg.User, message pkg.Message) error {
	return nil
}

func (m twitter) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	return nil
}
