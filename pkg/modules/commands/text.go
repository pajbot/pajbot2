package mcommands

import (
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

type textResponseCmd struct {
	response string
}

func newTextResponseCommand(response string) *textResponseCmd {
	return &textResponseCmd{
		response: response,
	}
}

func (c *textResponseCmd) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	return twitchactions.Say(c.response)
}
