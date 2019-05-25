package commands

import (
	"fmt"
	"strings"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commandlist"
	"github.com/pajbot/pajbot2/pkg/users"
	"github.com/pajbot/pajbot2/pkg/utils"
)

func init() {
	commandlist.Register(pkg.CommandInfo{
		Name:        "user",
		Description: "do user things",
		Maker:       NewUser,
	})
}

type userTarget struct {
	id   string
	name string
}

func updatePermissions(action, channelID string, target userTarget, parts []string) string {
	oldPermissions, err := users.GetUserPermissions(target.id, channelID)
	if err != nil {
		return "error getting old permissions"
	}

	channelName := channelID
	if channelID != "global" {
		channelName = "channel"
	}

	permissions := pkg.GetPermissionBits(parts)

	var newPermissions pkg.Permission

	switch action {
	case "set":
		newPermissions = permissions
	case "add":
		newPermissions = oldPermissions | permissions
	case "remove":
		newPermissions = oldPermissions &^ permissions
	}

	err = users.SetUserPermissions(target.id, channelID, newPermissions)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%s %s permissions changed from %b to %b (%s)", target.name, channelName, oldPermissions, newPermissions, action)
}

type User struct {
	Base

	subCommands       *subCommands
	defaultSubCommand string
}

func NewUser() pkg.CustomCommand2 {
	u := &User{
		Base:              NewBase(),
		subCommands:       newSubCommands(),
		defaultSubCommand: "print",
	}

	u.Base.UserCooldown = 0
	u.Base.GlobalCooldown = 0

	u.subCommands.add("print", &subCommand{
		permission: pkg.PermissionNone,
		cb: func(botChannel pkg.BotChannel, target userTarget, channel pkg.Channel, user pkg.User, parts []string) string {
			channelPermissions, err := users.GetUserChannelPermissions(target.id, channel.GetID())
			if err != nil {
				return "error getting channel permission: " + err.Error()
			}
			globalPermissions, err := users.GetUserGlobalPermissions(target.id)
			if err != nil {
				return "error getting global permission: " + err.Error()
			}
			permissions := channelPermissions | globalPermissions

			return fmt.Sprintf("%s permissions: %b (global: %b, channel: %b)", target.name, permissions, globalPermissions, channelPermissions)
		},
	})

	u.subCommands.addSC("set_global_permission", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(botChannel pkg.BotChannel, target userTarget, channel pkg.Channel, user pkg.User, parts []string) string {
			if len(parts) < 4 {
				return "usage: !user USERNAME set_global_permissions permission1 permission2"
			}

			return updatePermissions("set", "global", target, parts[3:])
		},
	})

	u.subCommands.addSC("set_channel_permission", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(botChannel pkg.BotChannel, target userTarget, channel pkg.Channel, user pkg.User, parts []string) string {
			if len(parts) < 4 {
				return "usage: !user USERNAME set_channel_permissions permission1 permission2"
			}

			return updatePermissions("set", channel.GetID(), target, parts[3:])
		},
	})

	u.subCommands.addSC("add_global_permission", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(botChannel pkg.BotChannel, target userTarget, channel pkg.Channel, user pkg.User, parts []string) string {
			if len(parts) < 4 {
				return "usage: !user USERNAME add_global_permissions permission1 permission2"
			}

			return updatePermissions("add", "global", target, parts[3:])
		},
	})

	u.subCommands.addSC("add_channel_permission", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(botChannel pkg.BotChannel, target userTarget, channel pkg.Channel, user pkg.User, parts []string) string {
			if len(parts) < 4 {
				return "usage: !user USERNAME add_channel_permissions permission1 permission2"
			}

			return updatePermissions("add", channel.GetID(), target, parts[3:])
		},
	})

	u.subCommands.addSC("remove_global_permission", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(botChannel pkg.BotChannel, target userTarget, channel pkg.Channel, user pkg.User, parts []string) string {
			if len(parts) < 4 {
				return "usage: !user USERNAME remove_global_permissions permission1 permission2"
			}

			return updatePermissions("remove", "global", target, parts[3:])
		},
	})

	u.subCommands.addSC("remove_channel_permission", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(botChannel pkg.BotChannel, target userTarget, channel pkg.Channel, user pkg.User, parts []string) string {
			if len(parts) < 4 {
				return "usage: !user USERNAME remove_channel_permissions permission1 permission2"
			}

			return updatePermissions("remove", channel.GetID(), target, parts[3:])
		},
	})

	return u
}

func (c *User) Trigger(botChannel pkg.BotChannel, parts []string, user pkg.User, message pkg.Message, action pkg.Action) {
	if len(parts) < 2 {
		return
	}

	targetName := utils.FilterUsername(parts[1])
	if targetName == "" {
		botChannel.Mention(user, "invalid username")
		return
	}

	targetUserID := botChannel.Bot().GetUserStore().GetID(targetName)
	if targetUserID == "" {
		botChannel.Mention(user, "no user with this ID")
		return
	}

	target := userTarget{
		id:   targetUserID,
		name: targetName,
	}

	subCommandName := c.defaultSubCommand
	if len(parts) >= 3 {
		subCommandName = strings.ToLower(parts[2])
	}

	if subCommand, ok := c.subCommands.find(subCommandName); ok {
		response := subCommand.run(botChannel, target, botChannel.Channel(), user, parts)
		if response != "" {
			botChannel.Mention(user, response)
		}
	}
}
