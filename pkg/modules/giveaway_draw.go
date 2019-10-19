package modules

import (
	"math/rand"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

type giveawayCmdDraw struct {
	m *giveaway
}

func newGiveawayCmdDraw(m *giveaway) *giveawayCmdDraw {
	return &giveawayCmdDraw{
		m: m,
	}
}

func (c *giveawayCmdDraw) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !event.User.IsModerator() {
		// User does not have permission to stop a giveaway
		return nil
	}

	if len(c.m.entrants) == 0 {
		return twitchactions.Say("No one has joined the giveaway")
	}

	winnerIndex := rand.Intn(len(c.m.entrants))
	winnerUsername := c.m.entrants[winnerIndex]

	c.m.entrants = append(c.m.entrants[:winnerIndex], c.m.entrants[winnerIndex+1:]...)

	return twitchactions.Say(winnerUsername + " just won the sub emote giveaway PogChamp")
}
