package modules

import (
	"github.com/pajbot/pajbot2/pkg"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
)

func init() {
	Register("twitter", func() pkg.ModuleSpec {
		return &Spec{
			id:    "twitter",
			name:  "Twitter",
			maker: newTwitter,

			enabledByDefault: false,
		}
	})
}

type twitter struct {
	mbase.Base
}

func newTwitter(b mbase.Base) pkg.Module {
	m := &twitter{
		Base: b,
	}

	// FIXME
	m.Initialize()

	return m
}

func (m *twitter) Initialize() {
	recv := m.BotChannel().Bot().Application().MIMO().Subscriber("twitter")
	go func() {
		for raw := range recv {
			message := raw.(string)
			m.BotChannel().Say("got tweet: " + message)
		}
	}()

}
