package modules

import (
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/eventemitter"
)

type welcome struct {
	botChannel pkg.BotChannel

	server      *server
	connections []*eventemitter.Listener
}

func newWelcome() pkg.Module {
	return &welcome{
		server: &_server,
	}
}

var welcomeSpec = &moduleSpec{
	id:    "welcome",
	name:  "Welcome",
	maker: newWelcome,

	enabledByDefault: false,
}

func (m *welcome) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel

	conn, err := m.botChannel.Events().Listen("on_join", func() error {
		go m.botChannel.Say("pb2 joined")
		return nil
	}, 100)
	if err != nil {
		return err
	}

	m.connections = append(m.connections, conn)

	return nil
}

func (m *welcome) Disable() error {
	for _, c := range m.connections {
		c.Disconnected = true
	}
	return nil
}

func (m *welcome) Spec() pkg.ModuleSpec {
	return welcomeSpec
}

func (m *welcome) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m welcome) OnWhisper(bot pkg.BotChannel, user pkg.User, message pkg.Message) error {
	return nil
}

func (m welcome) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	return nil
}
