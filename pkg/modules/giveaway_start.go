package modules

import (
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

type giveawayCmdStart struct {
	m *giveaway
}

func newGiveawayCmdStart(m *giveaway) *giveawayCmdStart {
	return &giveawayCmdStart{
		m: m,
	}
}

func (c *giveawayCmdStart) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !event.User.IsModerator() {
		// User does not have permission to start a giveaway
		return nil
	}

	if c.m.stopped() {
		c.m.start()
		return twitchactions.Say("Started giveaway")
	}

	return twitchactions.Say("Giveaway already started")
}
