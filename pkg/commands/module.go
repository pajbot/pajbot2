package commands

import (
	"fmt"
	"strings"

	"github.com/pajlada/pajbot2/pkg"
)

type moduleCommand struct {
	Base

	subCommands       *subCommands
	defaultSubCommand string
}

func NewModule() pkg.CustomCommand2 {
	u := &moduleCommand{
		Base:              NewBase(),
		subCommands:       newSubCommands(),
		defaultSubCommand: "list",
	}

	u.Base.UserCooldown = 0
	u.Base.GlobalCooldown = 0

	u.subCommands.add("list", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(botChannel pkg.BotChannel, target userTarget, channel pkg.Channel, user pkg.User, parts []string) string {
			return "list modules"
		},
	})

	u.subCommands.add("enable", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(botChannel pkg.BotChannel, target userTarget, channel pkg.Channel, user pkg.User, parts []string) string {
			if len(parts) < 3 {
				return "usage: !module enable MODULE_ID"
			}

			moduleID := parts[2]

			err := botChannel.EnableModule(moduleID)
			if err != nil {
				return err.Error()
			}

			return fmt.Sprintf("Enabled module %s", moduleID)
		},
	})

	u.subCommands.addSC("disable", &subCommand{
		permission: pkg.PermissionAdmin,
		cb: func(botChannel pkg.BotChannel, target userTarget, channel pkg.Channel, user pkg.User, parts []string) string {
			if len(parts) < 3 {
				return "usage: !module disable MODULE_ID"
			}

			moduleID := parts[2]

			err := botChannel.DisableModule(moduleID)
			if err != nil {
				return err.Error()
			}

			return fmt.Sprintf("Disabled module %s", moduleID)
		},
	})

	return u
}

func (c *moduleCommand) Trigger(botChannel pkg.BotChannel, parts []string, user pkg.User, message pkg.Message, action pkg.Action) {
	subCommandName := c.defaultSubCommand
	if len(parts) >= 2 {
		subCommandName = strings.ToLower(parts[1])
	}

	if subCommand, ok := c.subCommands.find(subCommandName); ok {
		response := subCommand.run(botChannel, userTarget{}, botChannel.Channel(), user, parts)
		if response != "" {
			botChannel.Mention(user, response)
		}
	}
}
