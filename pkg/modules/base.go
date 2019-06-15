package modules

import (
	"fmt"

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

	fmt.Println("Loaded parameters:", b.parameters)

	return b
}

func (b base) BotChannel() pkg.BotChannel {
	return b.bot
}

func (b *base) MarshalJSON() ([]byte, error) {
	fmt.Println("BASE MARSHAL JSON")

	return nullBuffer, nil
}

func (b *base) Parameters() map[string]pkg.ModuleParameter {
	return b.parameters
}

func (b *base) LoadSettings(settingsBytes []byte) error {
	if len(b.parameters) == 0 {
		return nil
	}

	values, err := loadParameters(settingsBytes)
	if err != nil {
		return err
	}

	for key, parameter := range b.parameters {
		if value, ok := values[key]; ok {
			parameter.SetInterface(value)
		}
	}

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

func (b *base) setParameter(key string, value string) error {
	// 1. Find parameter spec (This includes type of the parameter)
	param, ok := b.parameters[key]
	if !ok {
		return fmt.Errorf("No parameter found with the key '%s'", key)
	}

	// 2. Parse `value` according to that parameter spec
	if err := param.Parse(value); err != nil {
		return err
	}

	// 3. If the parameter was updated, update any linked values

	return nil
}
