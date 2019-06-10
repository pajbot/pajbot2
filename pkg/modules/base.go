package modules

import (
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/eventemitter"
)

type base struct {
	spec pkg.ModuleSpec
	bot  pkg.BotChannel

	server *server

	connections []*eventemitter.Listener
}

func newBase(spec pkg.ModuleSpec, bot pkg.BotChannel) base {
	return base{
		spec: spec,
		bot:  bot,

		server: &_server,
	}
}

func (b base) Spec() pkg.ModuleSpec {
	return b.spec
}

func (b base) BotChannel() pkg.BotChannel {
	return b.bot
}

func (b *base) LoadSettings(settingsBytes []byte) error {
	return nil
}

func (b *base) Handle(name string, cb func()) {

}

func (b *base) Disable() error {
	for _, c := range b.connections {
		c.Disconnected = true
	}
	return nil
}

func (b base) OnWhisper(event pkg.MessageEvent) pkg.Actions {
	return nil
}

func (b base) OnMessage(event pkg.MessageEvent) pkg.Actions {
	return nil
}
