package commands

import (
	"fmt"
	"strings"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/commandlist"
)

func init() {
	commandlist.Register(pkg.CommandInfo{
		Name:        "Join",
		Description: "xd join lol",
		Maker:       NewJoin,
	})
}

type Join struct {
	Base
}

func NewJoin() pkg.CustomCommand2 {
	return &Join{
		Base: NewBase(),
	}
}

func (c Join) Trigger(botChannel pkg.BotChannel, parts []string, user pkg.User, message pkg.Message, action pkg.Action) {
	if !user.HasGlobalPermission(pkg.PermissionAdmin) {
		botChannel.Mention(user, "you do not have permission to use this command. Admin permission is required")
		return
	}

	if len(parts) < 2 {
		return
	}

	channelName := parts[1]

	if strings.EqualFold(channelName, botChannel.Bot().TwitchAccount().Name()) {
		botChannel.Mention(user, "I cannot join my own channel")
		return
	}

	channelID := botChannel.Bot().GetUserStore().GetID(channelName)
	if channelID == "" {
		botChannel.Mention(user, "no channel with that name exists")
		return
	}

	err := botChannel.Bot().JoinChannel(channelID)
	if err != nil {
		botChannel.Mention(user, err.Error())
		return
	}

	botChannel.Mention(user, fmt.Sprintf("joined channel %s(%s)", channelName, channelID))

}
