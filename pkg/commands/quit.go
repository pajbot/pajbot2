package commands

import (
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commandlist"
)

func init() {
	commandlist.Register(pkg.CommandInfo{
		Name:        "Quit",
		Description: "quit the bot",
		Maker:       NewQuit,
	})
}

type Quit struct {
	Base
}

func NewQuit() pkg.CustomCommand2 {
	c := &Quit{
		Base: NewBase(),
	}

	return c
}

func (c *Quit) Trigger(botChannel pkg.BotChannel, parts []string, user pkg.User, message pkg.Message, action pkg.Action) {
	if !user.HasPermission(botChannel.Channel(), pkg.PermissionAdmin) {
		return
	}

	botChannel.Bot().Quit("hehe")
}
