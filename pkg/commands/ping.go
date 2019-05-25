package commands

import (
	"fmt"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commandlist"
	"github.com/pajbot/utils"
)

func init() {
	commandlist.Register(pkg.CommandInfo{
		Name:        "Ping",
		Description: "xd ping lol",
		Maker:       NewPing,
	})
}

type Ping struct {
	Base
}

func NewPing() pkg.CustomCommand2 {
	return &Ping{
		Base: NewBase(),
	}
}

func (c Ping) Trigger(botChannel pkg.BotChannel, parts []string, user pkg.User, message pkg.Message, action pkg.Action) {
	botChannel.Mention(user, fmt.Sprintf("pb2 has been running for %s", utils.TimeSince(startTime)))
}
