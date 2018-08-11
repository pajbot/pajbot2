package commands

import (
	"strconv"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/utils"
)

type Rank struct {
}

func (c Rank) Trigger(bot pkg.Sender, parts []string, channel pkg.Channel, source pkg.User, message pkg.Message, action pkg.Action) {
	var potentialTarget string
	targetID := source.GetID()

	if len(parts) >= 2 {
		potentialTarget = utils.FilterUsername(parts[1])
		if potentialTarget != "" {
			potentialTargetID := bot.GetUserStore().GetID(potentialTarget)
			if potentialTargetID != "" {
				targetID = potentialTargetID
			} else {
				potentialTarget = ""
			}
		}
	}

	rank := bot.PointRank(channel, targetID)
	if potentialTarget == "" {
		bot.Mention(channel, source, "you are rank "+strconv.FormatUint(rank, 10)+" in points")
	} else {
		bot.Mention(channel, source, potentialTarget+" is rank "+strconv.FormatUint(rank, 10)+" in points")
	}
}
