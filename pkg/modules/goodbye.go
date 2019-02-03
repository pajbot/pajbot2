package modules

import (
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/eventemitter"
)

type goodbye struct {
	botChannel pkg.BotChannel

	server      *server
	connections []*eventemitter.Listener
}

func newGoodbye() pkg.Module {
	return &goodbye{
		server: &_server,
	}
}

var goodbyeSpec = &moduleSpec{
	id:    "goodbye",
	name:  "Goodbye",
	maker: newGoodbye,

	enabledByDefault: false,
}

func (m *goodbye) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel

	conn, err := m.botChannel.Events().Listen("on_quit", func() error {
		go m.botChannel.Say("cya lol")
		return nil
	}, 100)
	if err != nil {
		return err
	}

	m.connections = append(m.connections, conn)

	return nil
}

func (m *goodbye) Disable() error {
	for _, c := range m.connections {
		c.Disconnected = true
	}
	return nil
}

func (m *goodbye) Spec() pkg.ModuleSpec {
	return goodbyeSpec
}

func (m *goodbye) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m goodbye) OnWhisper(bot pkg.BotChannel, user pkg.User, message pkg.Message) error {
	return nil
}

func (m goodbye) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	return nil
}
