package commands

import (
	"strings"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/users"
	"github.com/pajlada/pajbot2/pkg/utils"
)

type User struct {
}

func (c *User) Trigger(bot pkg.Sender, parts []string, channel pkg.Channel, source pkg.User, message pkg.Message, action pkg.Action) {
	if len(parts) < 2 {
		return
	}

	target := utils.FilterUsername(parts[1])
	if target == "" {
		// Invalid username
		return
	}

	subCommand := "print"
	if len(parts) >= 3 {
		subCommand = strings.ToLower(parts[2])
	}

	if subCommand == "print" {
		bot.Mention(channel, source, "print user "+target+" lol")
		return
	}

	if subCommand == "toggle_permission" {
		if !source.HasChannelPermission(channel, pkg.PermissionAdmin) && !source.HasGlobalPermission(pkg.PermissionAdmin) {
			bot.Mention(channel, source, "you do not have permission to use this command")
			return
		}

		if len(parts) < 4 {
			bot.Mention(channel, source, "usage: !user USERNAME toggle_permission PERMISSION")
			return
		}

		permission := pkg.GetPermissionBit(strings.ToLower(parts[3]))
		if permission == pkg.PermissionNone {
			bot.Mention(channel, source, "invalid permission")
			return
		}

		userID := bot.GetUserStore().GetID(target)

		oldPermission, err := users.GetUserChannelPermissions(userID, channel.GetID())
		if err != nil {
			bot.Mention(channel, source, "error getting permission: "+err.Error())
			return
		}

		err = users.SetUserChannelPermissions(bot.GetUserStore().GetID(target), channel.GetID(), oldPermission^permission)
		if err != nil {
			bot.Mention(channel, source, "error setting permission: "+err.Error())
		}

		return
	}

	bot.Mention(channel, source, "unhandled subcommand in user command: '"+subCommand+"'")
}
