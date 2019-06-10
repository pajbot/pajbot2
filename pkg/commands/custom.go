package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
	"github.com/pajbot/utils"
	normalize "github.com/pajlada/lidl-normalize"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

type GetUserName struct {
}

func (c GetUserName) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	const usage = `usage: !username USERID (i.e. !username 11148817)`
	const noValidUserIDs = `no valid user ids were given`

	userIDs := utils.FilterUserIDs(parts[1:])

	if len(userIDs) == 0 {
		return twitchactions.Mention(event.User, usage)
	}

	names := event.UserStore.GetNames(userIDs)
	var results []string
	for userID, username := range names {
		results = append(results, userID+"="+username)
	}

	if len(results) == 0 {
		return twitchactions.Mention(event.User, noValidUserIDs)
	}

	return twitchactions.Mention(event.User, strings.Join(results, ", "))
}

type GetPoints struct {
}

func (c GetPoints) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	// FIXME Re-implement (POINTS SYSTEM)
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

	// points := botChannel.Bot().GetPoints(channel, targetID)
	// if potentialTarget == "" {
	// 	botChannel.Mention(user, "you have "+strconv.FormatUint(points, 10)+" points")
	// } else {
	// 	botChannel.Mention(user, potentialTarget+" has "+strconv.FormatUint(points, 10)+" points")
	// }

	return nil
}

type AddPoints struct {
}

func (c AddPoints) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	// FIXME Re-implement (POINTS SYSTEM)
	// _, points := botChannel.Bot().AddPoints(channel, user.GetID(), uint64(rand.Int31n(50)))
	// botChannel.Mention(user, "you now have "+strconv.FormatUint(points, 10)+" points")
	return nil
}

type RemovePoints struct {
}

func (c RemovePoints) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	// FIXME Re-implement (POINTS SYSTEM)
	// _, points := botChannel.Bot().RemovePoints(channel, user.GetID(), uint64(rand.Int31n(50)))
	// botChannel.Mention(user, "you now have "+strconv.FormatUint(points, 10)+" points")
	return nil
}

type Roulette struct {
}

func (c Roulette) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	// FIXME Re-implement (POINTS SYSTEM)

	// if len(parts) < 2 {
	// 	botChannel.Mention(user, "usage: !roulette 500 or !roulette all")
	// 	return
	// }

	// var pointsToRoulette uint64

	// if strings.ToLower(parts[1]) == "all" {
	// 	pointsToRoulette = botChannel.Bot().GetPoints(channel, user.GetID())
	// } else {
	// 	var err error
	// 	pointsToRoulette, err = strconv.ParseUint(parts[1], 10, 64)

	// 	if err != nil {
	// 		botChannel.Mention(user, "usage: !roulette 500 or !roulette all")
	// 		return
	// 	}
	// }

	// if pointsToRoulette == 0 {
	// 	botChannel.Mention(user, "you have 0 points, you can't roulette ResidentSleeper")
	// 	return
	// }

	// if result, _ := botChannel.Bot().RemovePoints(channel, user.GetID(), pointsToRoulette); !result {
	// 	botChannel.Mention(user, "you don't have enough points ResidentSleeper")
	// 	return
	// }

	// if rand.Int31n(2) == 0 {
	// 	// loss
	// 	botChannel.Mention(user, "you lost OMEGALUL")
	// } else {
	// 	// win
	// 	// TODO: Check for integer overflow?
	// 	_, newPoints := botChannel.Bot().AddPoints(channel, user.GetID(), pointsToRoulette*2)
	// 	botChannel.Mention(user, "you won PagChomp you now have "+strconv.FormatUint(newPoints, 10)+" points KKona")
	// }
	return nil
}

type GivePoints struct {
}

func (c GivePoints) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	const USAGE = `usage: !givepoints USER POINTS`

	// FIXME Re-implement (POINTS SYSTEM)
	// if len(parts) < 3 {
	// 	botChannel.Mention(user, USAGE)
	// 	return
	// }

	// target := utils.FilterUsername(parts[1])
	// if target == "" {
	// 	// Invalid username
	// 	return
	// }

	// targetID := botChannel.Bot().GetUserStore().GetID(target)
	// if targetID == "" {
	// 	// Invalid username
	// 	return
	// }

	// var pointsToGive uint64

	// if strings.ToLower(parts[2]) == "all" {
	// 	pointsToGive = botChannel.Bot().GetPoints(channel, user.GetID())
	// } else {
	// 	var err error
	// 	pointsToGive, err = strconv.ParseUint(parts[2], 10, 64)

	// 	if err != nil {
	// 		botChannel.Mention(user, USAGE)
	// 		return
	// 	}
	// }

	// if pointsToGive == 0 {
	// 	botChannel.Mention(user, USAGE)
	// 	return
	// }

	// if result, _ := botChannel.Bot().RemovePoints(channel, user.GetID(), pointsToGive); !result {
	// 	botChannel.Mention(user, "you don't have enough points ResidentSleeper")
	// 	return
	// }

	// botChannel.Bot().AddPoints(channel, targetID, pointsToGive)
	// botChannel.Mention(user, "you gave away "+strconv.FormatUint(pointsToGive, 10)+" points to @"+target)
	return nil
}

type Simplify struct {
}

func (c Simplify) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if !event.User.IsModerator() {
		return nil
	}

	if len(parts) <= 1 {
		return nil
	}

	normalizedMessage, err := normalize.Normalize(strings.Join(parts[1:], " "))
	if err != nil {
		return twitchactions.Mention(event.User, fmt.Sprintf("error normalizing string: %s", err.Error()))
	}

	return twitchactions.Mention(event.User, fmt.Sprintf("normalized string: '%s'", normalizedMessage))
}

type TimeMeOut struct {
}

func (c TimeMeOut) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if len(parts) < 2 {
		return nil
	}

	timeoutDuration, err := time.ParseDuration(parts[1])
	if err != nil {
		return twitchactions.Mention(event.User, "invalid duration format. use !timemeout 1s or !timemeout 5m")
	}

	var reason string

	if len(parts) > 2 {
		reason = strings.Join(parts[2:], " ")
	}

	return twitchactions.DoTimeout(event.User, timeoutDuration, reason)
}

type Test struct {
}

func (c Test) Trigger(parts []string, event pkg.MessageEvent) (actions pkg.Actions) {
	if !event.User.IsModerator() {
		return
	}

	if len(parts) <= 1 {
		return
	}

	variations, _, err := utils.MakeVariations(strings.Join(parts[1:], " "), true)
	if err != nil {
		actions = twitchactions.Mention(event.User, err.Error())
		return
	}

	actions = &twitchactions.Actions{}

	for _, variation := range variations {
		actions.Mention(event.User, fmt.Sprintf("variation %s", variation))
	}

	return
}

type IsLive struct {
}

func (c IsLive) Trigger(parts []string, event pkg.MessageEvent) (actions pkg.Actions) {
	if !event.User.IsModerator() {
		return
	}

	// FIXME: Re-implement stream status checker
	// if botChannel.Stream().Status().Live() {
	// 	startedAt := botChannel.Stream().Status().StartedAt()
	// 	return twitchactions.Mention(event.User, fmt.Sprintf("LIVE FOR %s KKona", utils.TimeSince(startedAt)))
	// }

	return twitchactions.Mention(event.User, "offline FeelsBadMan")
}
