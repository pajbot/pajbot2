package mcommands

import (
	"strconv"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

type cmdAddTrigger struct {
	m *CommandsModule
}

func newCmdAddTrigger(m *CommandsModule) *cmdAddTrigger {
	return &cmdAddTrigger{
		m: m,
	}
}

func (c *cmdAddTrigger) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !event.User.IsModerator() {
		return twitchactions.Mention(event.User, "no permission")
	}

	parts = parts[1:]

	if len(parts) < 2 {
		return twitchactions.Mention(event.User, "bad usage")
	}

	commandID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return twitchactions.Mentionf(event.User, "add trigger error: %s", err)
	}
	newTrigger := parts[1]

	err = c.m.addTrigger(commandID, newTrigger)
	if err != nil {
		return twitchactions.Mentionf(event.User, "add error: %s", err)
	}

	return twitchactions.Mentionf(event.User, "added command trigger: %s", newTrigger)
}
