package modules

import (
	"github.com/pajbot/pajbot2/pkg"
)

func init() {
	Register("twitter", func() pkg.ModuleSpec {
		return &moduleSpec{
			id:    "twitter",
			name:  "Twitter",
			maker: newTwitter,

			enabledByDefault: false,
		}
	})
}

type twitter struct {
	base
}

func newTwitter(b base) pkg.Module {
	m := &twitter{
		base: b,
	}

	// FIXME
	m.Initialize()

	return m
}

func (m *twitter) Initialize() {
	recv := m.bot.Bot().Application().MIMO().Subscriber("twitter")
	go func() {
		for raw := range recv {
			message := raw.(string)
			m.bot.Say("got tweet: " + message)
		}
	}()

}
