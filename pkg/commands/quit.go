package commands

import (
	"github.com/pajbot/pajbot2/internal/commands/base"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commandlist"
)

func init() {
	commandlist.Register(pkg.CommandInfo{
		Name:        "Quit",
		Description: "quit the bot",
		// FIXME
		// Maker:       NewQuit,
	})
}

type Quit struct {
	base.Command

	bot pkg.BotChannel
}

func NewQuit(bot pkg.BotChannel) pkg.CustomCommand2 {
	c := &Quit{
		Command: base.New(),

		bot: bot,
	}

	return c
}

func (c *Quit) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	user := event.User

	// FIXME: Channel should be part of the message event
	if !user.HasPermission(c.bot.Channel(), pkg.PermissionAdmin) {
		return nil
	}

	// FIXME: this should be an "on done action"
	c.bot.Bot().Quit("hehe")

	return nil
}
