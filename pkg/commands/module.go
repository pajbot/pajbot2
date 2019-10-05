package commands

import (
	"strings"

	"github.com/pajbot/pajbot2/internal/commands/base"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

type moduleCommand struct {
	base.Command

	subCommands       *subCommands
	defaultSubCommand string
}

func NewModule(bot pkg.BotChannel) pkg.CustomCommand2 {
	u := &moduleCommand{
		Command:           base.New(),
		subCommands:       newSubCommands(),
		defaultSubCommand: "list",
	}

	u.UserCooldown = 0
	u.GlobalCooldown = 0

	u.subCommands.add("list", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(parts []string, event pkg.MessageEvent) pkg.Actions {
			return twitchactions.Mention(event.User, "TODO: list modules")
		},
	})

	u.subCommands.add("enable", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(parts []string, event pkg.MessageEvent) pkg.Actions {
			if len(parts) < 3 {
				return twitchactions.Mention(event.User, "usage: !module enable MODULE_ID")
			}

			moduleID := parts[2]

			err := bot.EnableModule(moduleID)
			if err != nil {
				return twitchactions.Mention(event.User, err.Error())
			}

			return twitchactions.Mentionf(event.User, "Enabled module %s", moduleID)
		},
	})

	u.subCommands.addSC("disable", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(parts []string, event pkg.MessageEvent) pkg.Actions {
			if len(parts) < 3 {
				return twitchactions.Mention(event.User, "usage: !module disable MODULE_ID")
			}

			moduleID := parts[2]

			err := bot.DisableModule(moduleID)
			if err != nil {
				return twitchactions.Mention(event.User, err.Error())
			}

			return twitchactions.Mentionf(event.User, "Disabled module %s", moduleID)
		},
	})

	return u
}

func (c *moduleCommand) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	subCommandName := c.defaultSubCommand
	if len(parts) >= 2 {
		subCommandName = strings.ToLower(parts[1])
	}

	if subCommand, ok := c.subCommands.find(subCommandName); ok {
		return subCommand.run(parts, event)
	}

	return nil
}
