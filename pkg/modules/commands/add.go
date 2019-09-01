package mcommands

import (
	"strings"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

type addCmd struct {
	m *commandsModule
}

func newAddCmd(m *commandsModule) *addCmd {
	return &addCmd{
		m: m,
	}
}

func (c *addCmd) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !event.User.IsModerator() {
		return twitchactions.Mention(event.User, "no permission")
	}

	parts = parts[1:]

	if len(parts) < 2 {
		return twitchactions.Mention(event.User, "bad usage")
	}

	trigger := parts[0]
	response := strings.Join(parts[1:], " ")

	err := c.m.addToDB(trigger, response)
	if err != nil {
		return twitchactions.Mentionf(event.User, "add error: %s", err)
	}

	return twitchactions.Mentionf(event.User, "added command: %s", trigger)
}
