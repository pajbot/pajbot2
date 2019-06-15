package modules

import (
	"log"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/eventemitter"
)

type base struct {
	spec pkg.ModuleSpec
	bot  pkg.BotChannel

	server *server

	connections []*eventemitter.Listener

	parameters map[string]pkg.ModuleParameter
}

func newBase(spec pkg.ModuleSpec, bot pkg.BotChannel) base {
	b := base{
		spec: spec,
		bot:  bot,

		server: &_server,

		parameters: make(map[string]pkg.ModuleParameter),
	}

	for key, value := range spec.Parameters() {
		b.parameters[key] = value()
	}

	return b
}

func (b base) BotChannel() pkg.BotChannel {
	return b.bot
}

func (b *base) LoadSettings(settingsBytes []byte) error {
	if len(settingsBytes) == 0 {
		return nil
	}
	log.Println("Got settings data, but the module doesn't overwrite the LoadSettings function.")
	log.Println("Module ID:", b.ID())
	log.Println("Settings data:", string(settingsBytes))
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

func (b base) ID() string {
	return b.spec.ID()
}

func (b base) Type() pkg.ModuleType {
	return b.spec.Type()
}

func (b base) Priority() int {
	return b.spec.Priority()
}
