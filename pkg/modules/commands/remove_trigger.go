package mcommands

import (
	"strconv"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

type cmdRemoveTrigger struct {
	m *commandsModule
}

func newCmdRemoveTrigger(m *commandsModule) *cmdRemoveTrigger {
	return &cmdRemoveTrigger{
		m: m,
	}
}

func (c *cmdRemoveTrigger) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !event.User.IsModerator() {
		return twitchactions.Mention(event.User, "no permission")
	}

	parts = parts[1:]

	if len(parts) < 2 {
		return twitchactions.Mention(event.User, "bad usage")
	}

	commandID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return twitchactions.Mentionf(event.User, "remove trigger error: %s", err)
	}
	newTrigger := parts[1]

	err = c.m.removeTrigger(commandID, newTrigger)
	if err != nil {
		return twitchactions.Mentionf(event.User, "remove error: %s", err)
	}

	return twitchactions.Mentionf(event.User, "remove command trigger: %s", newTrigger)
}
