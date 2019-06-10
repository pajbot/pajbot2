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
		Name:        "Leave",
		Description: "xd leave lol",
		// Maker:       NewLeave,
	})
}

type Leave struct {
	base.Command

	bot pkg.BotChannel
}

func NewLeave(bot pkg.BotChannel) pkg.CustomCommand2 {
	return &Leave{
		Command: base.New(),

		bot: bot,
	}
}

func (c Leave) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	user := event.User

	if !user.HasGlobalPermission(pkg.PermissionAdmin) {
		return twitchactions.Mention(user, "you do not have permission to use this command. Admin permission is required")
	}

	if len(parts) < 2 {
		return nil
	}

	channelName := parts[1]

	if strings.EqualFold(channelName, c.bot.Bot().TwitchAccount().Name()) {
		return twitchactions.Mention(user, "I cannot leave my own channel")
	}

	channelID := c.bot.Bot().GetUserStore().GetID(channelName)
	if channelID == "" {
		return twitchactions.Mention(user, "no channel with that name exists")
	}

	err := c.bot.Bot().LeaveChannel(channelID)
	if err != nil {
		return twitchactions.Mention(user, "Error leaving channel: "+err.Error())
	}

	return twitchactions.Mention(user, fmt.Sprintf("left channel %s(%s)", channelName, channelID))
}
