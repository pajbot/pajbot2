package commands

import (
	"fmt"
	"strings"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/commandlist"
)

func init() {
	commandlist.Register(pkg.CommandInfo{
		Name:        "Leave",
		Description: "xd leave lol",
		Maker:       NewLeave,
	})
}

type Leave struct {
	Base
}

func NewLeave() pkg.CustomCommand2 {
	return &Leave{
		Base: NewBase(),
	}
}

func (c Leave) Trigger(botChannel pkg.BotChannel, parts []string, user pkg.User, message pkg.Message, action pkg.Action) {
	if !user.HasGlobalPermission(pkg.PermissionAdmin) {
		botChannel.Mention(user, "you do not have permission to use this command. Admin permission is required")
		return
	}

	if len(parts) < 2 {
		return
	}

	channelName := parts[1]

	if strings.EqualFold(channelName, botChannel.Bot().TwitchAccount().Name()) {
		botChannel.Mention(user, "I cannot leave my own channel")
		return
	}

	channelID := botChannel.Bot().GetUserStore().GetID(channelName)
	if channelID == "" {
		botChannel.Mention(user, "no channel with that name exists")
		return
	}

	err := botChannel.Bot().LeaveChannel(channelID)
	if err != nil {
		botChannel.Mention(user, "Error leaving channel: "+err.Error())
		return
	}

	botChannel.Mention(user, fmt.Sprintf("left channel %s(%s)", channelName, channelID))
}
