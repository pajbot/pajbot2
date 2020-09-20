package disconnect

import (
	"github.com/pajbot/pajbot2/pkg"
)

type Command struct {
	Bot pkg.BotChannel
}

func (c *Command) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !event.User.HasPermission(event.Channel, pkg.PermissionAdmin) {
		return nil
	}

	c.Bot.Bot().Disconnect()

	return nil
}
