package commands

import (
	"github.com/pajbot/pajbot2/internal/commands/base"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commandlist"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
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
	base.Command
}

func NewPing() pkg.CustomCommand2 {
	return &Ping{
		Command: base.New(),
	}
}

func (c Ping) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	return twitchactions.Mentionf(event.User, "pb2 has been running for %s", utils.TimeSince(startTime))
}
