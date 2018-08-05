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

	bot.Mention(channel, source, "unhandled subcommand in user command: '"+subCommand+"'")
}
