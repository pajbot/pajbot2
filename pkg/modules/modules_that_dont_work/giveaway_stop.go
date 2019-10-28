package modules

import (
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

type giveawayCmdStop struct {
	m *giveaway
}

func (c *giveawayCmdStop) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !event.User.IsModerator() {
		// User does not have permission to stop a giveaway
		return nil
	}

	if c.m.started() {
		c.m.stop()
		return twitchactions.Say("Stopped accepting people into the giveaway")
	}

	return twitchactions.Say("Giveaway already stopped")
}
