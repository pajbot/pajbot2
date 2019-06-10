package commands

import (
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

type subCommands struct {
	m map[string]*subCommand
}

func newSubCommands() *subCommands {
	return &subCommands{
		m: make(map[string]*subCommand),
	}
}

func (c *subCommands) add(name string, sc *subCommand) {
	c.m[name] = sc
}

func (c *subCommands) find(name string) (*subCommand, bool) {
	a, b := c.m[name]
	return a, b
}

func (c *subCommands) addSC(name string, sc *subCommand) {
	c.add(name, sc)
	c.add(name+"s", sc)
}

type subCommandFunc func(parts []string, event pkg.MessageEvent) pkg.Actions

type subCommand struct {
	permission pkg.Permission
	cb         subCommandFunc
}

func (c *subCommand) run(parts []string, event pkg.MessageEvent) pkg.Actions {
	if c.permission != pkg.PermissionNone {
		if !event.User.HasPermission(event.Channel, c.permission) {
			return twitchactions.Mention(event.User, "you do not have permission to use this command")
		}

	}

	return c.cb(parts, event)
}
