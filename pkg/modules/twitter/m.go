package twitter

import (
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/modules"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
)

func init() {
	modules.Register("twitter", func() pkg.ModuleSpec {
		return modules.NewSpec("twitter", "Twitter", false, newTwitter)
	})
}

type twitter struct {
	mbase.Base
}

func newTwitter(b *mbase.Base) pkg.Module {
	m := &twitter{
		Base: *b,
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
