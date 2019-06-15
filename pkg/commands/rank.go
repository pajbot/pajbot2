package commands

import (
	"github.com/pajbot/pajbot2/pkg"
)

type Rank struct {
}

func (c Rank) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	// FIXME: Re-implement (POINTS SYSTEM)
	// var potentialTarget string
	// targetID := event.User.GetID()

	// if len(parts) >= 2 {
	// 	potentialTarget = utils.FilterUsername(parts[1])
	// 	if potentialTarget != "" {
	// 		potentialTargetID := event.UserStore.GetID(potentialTarget)
	// 		if potentialTargetID != "" {
	// 			targetID = potentialTargetID
	// 		} else {
	// 			potentialTarget = ""
	// 		}
	// 	}
	// }

	// rank := botChannel.Bot().PointRank(event.Channel, targetID)
	// if potentialTarget == "" {
	// 	return twitchactions.Mention(event.User, "you are rank "+strconv.FormatUint(rank, 10)+" in points")
	// }

	// return twitchactions.Mention(event.User, potentialTarget+" is rank "+strconv.FormatUint(rank, 10)+" in points")

	return nil
}
