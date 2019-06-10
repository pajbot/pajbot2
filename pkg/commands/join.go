package commands

import (
	"fmt"
	"strings"

	"github.com/pajbot/pajbot2/internal/commands/base"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commandlist"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

func init() {
	commandlist.Register(pkg.CommandInfo{
		Name:        "Join",
		Description: "xd join lol",
		// FIXME
		// Maker:       NewJoin,
	})
}

type Join struct {
	base.Command

	bot pkg.BotChannel
}

func NewJoin(bot pkg.BotChannel) pkg.CustomCommand2 {
	return &Join{
		Command: base.New(),

		bot: bot,
	}
}

func (c Join) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	user := event.User
	if !user.HasGlobalPermission(pkg.PermissionAdmin) {
		const errorMessage = "you do not have permission to use this command. Admin permission is required"
		return twitchactions.Mention(user, errorMessage)
	}

	if len(parts) < 2 {
		return nil
	}

	channelName := parts[1]

	if strings.EqualFold(channelName, c.bot.Bot().TwitchAccount().Name()) {
		const errorMessage = "I cannot join my own channel"
		return twitchactions.Mention(user, errorMessage)
	}

	channelID := c.bot.Bot().GetUserStore().GetID(channelName)
	if channelID == "" {
		const errorMessage = "no channel with that name exists"
		return twitchactions.Mention(user, errorMessage)
	}

	err := c.bot.Bot().JoinChannel(channelID)
	if err != nil {
		return twitchactions.Mention(user, err.Error())
	}

	return twitchactions.Mention(user, fmt.Sprintf("joined channel %s(%s)", channelName, channelID))
}
