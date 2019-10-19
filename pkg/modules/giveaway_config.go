package modules

import (
	"strings"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

type giveawayCmdConfig struct {
	m *giveaway
}

var (
	giveawayValidKeys = map[string]string{
		"emoteid":   "EmoteID",
		"emotename": "EmoteName",
	}
)

func (c *giveawayCmdConfig) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !event.User.IsModerator() {
		return nil
	}

	if len(parts) <= 1 {
		return nil
	}

	key := strings.ToLower(parts[1])
	var ok bool
	key, ok = giveawayValidKeys[key]
	if !ok {
		return twitchactions.Mentionf(event.User, "%s is not a valid giveaway parameter key", key)
	}

	if len(parts) >= 3 {
		value := parts[2]
		return c.m.SetParameterResponse(key, value, event)
	}

	switch key {
	case "EmoteID":
		return twitchactions.Mentionf(event.User, "emote ID is %s", c.m.emoteID)
	case "EmoteName":
		return twitchactions.Mentionf(event.User, "emote name is %s", c.m.emoteName)
	}

	return twitchactions.Mentionf(event.User, "error: unhandled key '%s' in giveaway_config.go", key)
}
