package modules

import (
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

type giveawayCmdStart struct {
	m *giveaway
}

func (c *giveawayCmdStart) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !event.User.IsModerator() {
		// User does not have permission to start a giveaway
		return nil
	}

	if c.m.emoteID == "" {
		return twitchactions.Mention(event.User, "No emote ID/name set. Use '!25config emoteid 98374583' or '!25config emotename NaM' to configure this module")
	}

	if c.m.stopped() {
		c.m.start()
		return twitchactions.Sayf("Started giveaway, type %s to join the giveaway", c.m.emoteName)
	}

	return twitchactions.Say("Giveaway already started")
}
