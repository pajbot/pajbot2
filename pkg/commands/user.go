package commands

import (
	"fmt"
	"strings"

	"github.com/pajbot/pajbot2/internal/commands/base"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commandlist"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
	"github.com/pajbot/pajbot2/pkg/users"
	"github.com/pajbot/utils"
)

func init() {
	commandlist.Register(pkg.CommandInfo{
		Name:        "user",
		Description: "do user things",
		// FIXME
		// Maker:       NewUser,
	})
}

func parseUser(bot pkg.BotChannel, content string) pkg.User {
	targetName := utils.FilterUsername(content)
	if targetName == "" {
		return nil
	}

	targetUserID := bot.Bot().GetUserStore().GetID(targetName)
	if targetUserID == "" {
		return nil
	}

	return users.NewSimpleTwitchUser(targetUserID, targetName)
}

func updatePermissions(action, channelID string, user pkg.User, parts []string, event pkg.MessageEvent) pkg.Actions {
	oldPermissions, err := users.GetUserPermissions(user.GetID(), channelID)
	if err != nil {
		return twitchactions.Mention(event.User, "error getting old permissions")
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

	err = users.SetUserPermissions(user.GetID(), channelID, newPermissions)
	if err != nil {
		return twitchactions.Mention(event.User, err.Error())
	}

	return twitchactions.Mention(event.User, fmt.Sprintf("%s %s permissions changed from %b to %b (%s)", user.GetName(), channelName, oldPermissions, newPermissions, action))
}

type User struct {
	base.Command

	subCommands       *subCommands
	defaultSubCommand string
}

func NewUser(bot pkg.BotChannel) pkg.CustomCommand2 {
	u := &User{
		Command:           base.New(),
		subCommands:       newSubCommands(),
		defaultSubCommand: "print",
	}

	u.UserCooldown = 0
	u.GlobalCooldown = 0

	// FIXME
	channel := bot.Channel()

	u.subCommands.add("print", &subCommand{
		permission: pkg.PermissionNone,
		cb: func(parts []string, event pkg.MessageEvent) pkg.Actions {
			target := parseUser(bot, parts[1])
			if target == nil {
				return twitchactions.Mention(event.User, "no valid user found")
			}

			channelPermissions, err := users.GetUserChannelPermissions(target.GetID(), channel.GetID())
			if err != nil {
				return twitchactions.Mention(event.User, "error getting channel permission: "+err.Error())
			}
			globalPermissions, err := users.GetUserGlobalPermissions(target.GetID())
			if err != nil {
				return twitchactions.Mention(event.User, "error getting global permission: "+err.Error())
			}
			permissions := channelPermissions | globalPermissions

			return twitchactions.Mention(event.User, fmt.Sprintf("%s permissions: %b (global: %b, channel: %b)", target.GetName(), permissions, globalPermissions, channelPermissions))
		},
	})

	u.subCommands.addSC("set_global_permission", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(parts []string, event pkg.MessageEvent) pkg.Actions {
			if len(parts) < 4 {
				return twitchactions.Mention(event.User, "usage: !user USERNAME set_global_permissions permission1 permission2")
			}

			target := parseUser(bot, parts[1])
			if target == nil {
				return twitchactions.Mention(event.User, "no valid user found")
			}

			return updatePermissions("set", "global", target, parts[3:], event)
		},
	})

	u.subCommands.addSC("set_channel_permission", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(parts []string, event pkg.MessageEvent) pkg.Actions {
			if len(parts) < 4 {
				return twitchactions.Mention(event.User, "usage: !user USERNAME set_channel_permissions permission1 permission2")
			}

			target := parseUser(bot, parts[1])
			if target == nil {
				return twitchactions.Mention(event.User, "no valid user found")
			}

			return updatePermissions("set", channel.GetID(), target, parts[3:], event)
		},
	})

	u.subCommands.addSC("add_global_permission", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(parts []string, event pkg.MessageEvent) pkg.Actions {
			if len(parts) < 4 {
				return twitchactions.Mention(event.User, "usage: !user USERNAME add_global_permissions permission1 permission2")
			}

			target := parseUser(bot, parts[1])
			if target == nil {
				return twitchactions.Mention(event.User, "no valid user found")
			}

			return updatePermissions("add", "global", target, parts[3:], event)
		},
	})

	u.subCommands.addSC("add_channel_permission", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(parts []string, event pkg.MessageEvent) pkg.Actions {
			if len(parts) < 4 {
				return twitchactions.Mention(event.User, "usage: !user USERNAME add_channel_permissions permission1 permission2")
			}

			target := parseUser(bot, parts[1])
			if target == nil {
				return twitchactions.Mention(event.User, "no valid user found")
			}

			return updatePermissions("add", channel.GetID(), target, parts[3:], event)
		},
	})

	u.subCommands.addSC("remove_global_permission", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(parts []string, event pkg.MessageEvent) pkg.Actions {
			if len(parts) < 4 {
				return twitchactions.Mention(event.User, "usage: !user USERNAME remove_global_permissions permission1 permission2")
			}

			target := parseUser(bot, parts[1])
			if target == nil {
				return twitchactions.Mention(event.User, "no valid user found")
			}

			return updatePermissions("remove", "global", target, parts[3:], event)
		},
	})

	u.subCommands.addSC("remove_channel_permission", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(parts []string, event pkg.MessageEvent) pkg.Actions {
			if len(parts) < 4 {
				return twitchactions.Mention(event.User, "usage: !user USERNAME remove_channel_permissions permission1 permission2")
			}

			target := parseUser(bot, parts[1])
			if target == nil {
				return twitchactions.Mention(event.User, "no valid user found")
			}

			return updatePermissions("remove", channel.GetID(), target, parts[3:], event)
		},
	})

	return u
}

func (c *User) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if len(parts) < 2 {
		return nil
	}

	subCommandName := c.defaultSubCommand
	if len(parts) >= 3 {
		subCommandName = strings.ToLower(parts[2])
	}

	if subCommand, ok := c.subCommands.find(subCommandName); ok {
		return subCommand.run(parts, event)
	}

	return nil
}
