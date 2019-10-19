package modules

import (
	"strings"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

type giveawayCmdConfig struct {
	m *giveaway
}

func newGiveawayCmdConfig(m *giveaway) *giveawayCmdConfig {
	return &giveawayCmdConfig{
		m: m,
	}
}

var (
	giveawayValidKeys = map[string]string{
		"emoteid":   "emoteID",
		"emotename": "emoteName",
	}
)

func (c *giveawayCmdConfig) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !event.User.IsModerator() {
		// User does not have permission to stop a giveaway
		return nil
	}

	if len(parts) <= 1 {
		return nil
	}

	key := strings.ToLower(parts[0])
	var ok bool
	key, ok = giveawayValidKeys[key]
	if !ok {
		return nil
	}

	if len(parts) >= 2 {
		value := parts[1]
		return c.m.SetParameterResponse(key, value, event)
	}

	switch key {
	case "emoteID":
		return twitchactions.Mentionf(event.User, "emote ID is %s", c.m.emoteID)
	case "emoteName":
		return twitchactions.Mentionf(event.User, "emote name is %s", c.m.emoteName)
	}

	return twitchactions.Mentionf(event.User, "error: unhandled key '%s' in giveaway_config.go", key)
}
