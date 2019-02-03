package commands

import "github.com/pajlada/pajbot2/pkg"

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

type subCommandFunc func(pkg.BotChannel, userTarget, pkg.Channel, pkg.User, []string) string

type subCommand struct {
	permission pkg.Permission
	cb         subCommandFunc
}

func (c *subCommand) run(botChannel pkg.BotChannel, target userTarget, channel pkg.Channel, user pkg.User, parts []string) string {
	if c.permission != pkg.PermissionNone {
		if !user.HasPermission(channel, c.permission) {
			return "you do not have permission to use this command"
		}

	}

	return c.cb(botChannel, target, channel, user, parts)
}
