package commands

import (
	"fmt"
	"runtime"

	"github.com/pajbot/pajbot2/internal/commands/base"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commandlist"
	"github.com/pajbot/pajbot2/pkg/common"
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
	msg := fmt.Sprintf("pb2 has been running for %s (%s %s", utils.TimeSince(startTime), common.Version(), runtime.Version())
	if common.BuildTime != "" {
		msg += " built " + common.BuildTime
	}
	msg += ")"

	return twitchactions.Mention(event.User, msg)
}
