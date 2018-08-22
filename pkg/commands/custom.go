package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/mvdan/xurls"
	normalize "github.com/pajlada/lidl-normalize"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/utils"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

type GetUserID struct {
}

func (c GetUserID) Trigger(bot pkg.Sender, parts []string, channel pkg.Channel, source pkg.User, message pkg.Message, action pkg.Action) {
	usernames := utils.FilterUsernames(parts[1:])

	if len(usernames) == 0 {
		bot.Say(channel, "@"+source.GetName()+", usage: !userid USERNAME (i.e. !userid pajlada)")
		return
	}

	/*
		onHTTPError := func(statusCode int, statusMessage, errorMessage string) {
			bot.Say(channel, "@"+source.GetName()+", an error occursed processing your command ("+errorMessage+", "+statusMessage+")")
		}

		onInternalError := func(err error) {
			bot.Say(channel, "@"+source.GetName()+", an internal error occursed processing your command ("+err.Error()+")")
		}

		onSuccess := func(data []gotwitch.User) {
			if len(data) == 0 {
				bot.Say(channel, "@"+source.GetName()+", no valid usernames were given")
				return
			}
			var results []string
			for _, d := range data {
				results = append(results, d.Login+"="+d.ID)
			}

			bot.Say(channel, "@"+source.GetName()+", "+strings.Join(results, ", "))
			fmt.Printf("%#v\n", data)
		}
	*/

	userIDs := bot.GetUserStore().GetIDs(usernames)
	var results []string
	for username, userID := range userIDs {
		results = append(results, username+"="+userID)
	}

	if len(results) == 0 {
		bot.Mention(channel, source, "no valid usernames were given")
		return
	}

	bot.Mention(channel, source, strings.Join(results, ", "))

	// apirequest.Twitch.GetUsersByLogin(usernames, onSuccess, onHTTPError, onInternalError)
}

type GetPoints struct {
}

func (c GetPoints) Trigger(bot pkg.Sender, parts []string, channel pkg.Channel, source pkg.User, message pkg.Message, action pkg.Action) {
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

	points := bot.GetPoints(channel, targetID)
	if potentialTarget == "" {
		bot.Mention(channel, source, "you have "+strconv.FormatUint(points, 10)+" points")
	} else {
		bot.Mention(channel, source, potentialTarget+" has "+strconv.FormatUint(points, 10)+" points")
	}
}

type AddPoints struct {
}

func (c AddPoints) Trigger(bot pkg.Sender, parts []string, channel pkg.Channel, source pkg.User, message pkg.Message, action pkg.Action) {
	_, points := bot.AddPoints(channel, source.GetID(), uint64(rand.Int31n(50)))
	bot.Mention(channel, source, "you now have "+strconv.FormatUint(points, 10)+" points")
}

type RemovePoints struct {
}

func (c RemovePoints) Trigger(bot pkg.Sender, parts []string, channel pkg.Channel, source pkg.User, message pkg.Message, action pkg.Action) {
	_, points := bot.RemovePoints(channel, source.GetID(), uint64(rand.Int31n(50)))
	bot.Mention(channel, source, "you now have "+strconv.FormatUint(points, 10)+" points")
}

type Roulette struct {
}

func (c Roulette) Trigger(bot pkg.Sender, parts []string, channel pkg.Channel, source pkg.User, message pkg.Message, action pkg.Action) {
	if len(parts) < 2 {
		bot.Mention(channel, source, "usage: !roulette 500 or !roulette all")
		return
	}

	var pointsToRoulette uint64

	if strings.ToLower(parts[1]) == "all" {
		pointsToRoulette = bot.GetPoints(channel, source.GetID())
	} else {
		var err error
		pointsToRoulette, err = strconv.ParseUint(parts[1], 10, 64)

		if err != nil {
			bot.Mention(channel, source, "usage: !roulette 500 or !roulette all")
			return
		}
	}

	if pointsToRoulette == 0 {
		bot.Mention(channel, source, "you have 0 points, you can't roulette ResidentSleeper")
		return
	}

	if result, _ := bot.RemovePoints(channel, source.GetID(), pointsToRoulette); !result {
		bot.Mention(channel, source, "you don't have enough points ResidentSleeper")
		return
	}

	if rand.Int31n(2) == 0 {
		// loss
		bot.Mention(channel, source, "you lost OMEGALUL")
	} else {
		// win
		// TODO: Check for integer overflow?
		_, newPoints := bot.AddPoints(channel, source.GetID(), pointsToRoulette*2)
		bot.Mention(channel, source, "you won PagChomp you now have "+strconv.FormatUint(newPoints, 10)+" points KKona")
	}
}

type GivePoints struct {
}

func (c GivePoints) Trigger(bot pkg.Sender, parts []string, channel pkg.Channel, source pkg.User, message pkg.Message, action pkg.Action) {
	const USAGE = `usage: !givepoints USER POINTS`
	if len(parts) < 3 {
		bot.Mention(channel, source, USAGE)
		return
	}

	target := utils.FilterUsername(parts[1])
	if target == "" {
		// Invalid username
		return
	}

	targetID := bot.GetUserStore().GetID(target)
	if targetID == "" {
		// Invalid username
		return
	}

	var pointsToGive uint64

	if strings.ToLower(parts[2]) == "all" {
		pointsToGive = bot.GetPoints(channel, source.GetID())
	} else {
		var err error
		pointsToGive, err = strconv.ParseUint(parts[2], 10, 64)

		if err != nil {
			bot.Mention(channel, source, USAGE)
			return
		}
	}

	if pointsToGive == 0 {
		bot.Mention(channel, source, USAGE)
		return
	}

	if result, _ := bot.RemovePoints(channel, source.GetID(), pointsToGive); !result {
		bot.Mention(channel, source, "you don't have enough points ResidentSleeper")
		return
	}

	bot.AddPoints(channel, targetID, pointsToGive)
	bot.Mention(channel, source, "you gave away "+strconv.FormatUint(pointsToGive, 10)+" points to @"+target)
}

type Ping struct {
}

func (c Ping) Trigger(bot pkg.Sender, parts []string, channel pkg.Channel, source pkg.User, message pkg.Message, action pkg.Action) {
	bot.Mention(channel, source, fmt.Sprintf("pb2 has been running for %s", time.Since(startTime)))
}

type Simplify struct {
}

func (c Simplify) Trigger(bot pkg.Sender, parts []string, channel pkg.Channel, source pkg.User, message pkg.Message, action pkg.Action) {
	if source.IsModerator() || source.IsBroadcaster(channel) {
		if len(parts) > 1 {
			normalizedMessage, err := normalize.Normalize(strings.Join(parts[1:], " "))
			if err != nil {
				bot.Mention(channel, source, fmt.Sprintf("error normalizing string: %s", err.Error()))
				return
			}

			bot.Mention(channel, source, fmt.Sprintf("normalized string: '%s'", normalizedMessage))
		}
	}
}

type Test struct {
}

func (c Test) Trigger(bot pkg.Sender, parts []string, channel pkg.Channel, source pkg.User, message pkg.Message, action pkg.Action) {
	if !source.IsModerator() && !source.IsBroadcaster(channel) {
		return
	}

	if len(parts) <= 1 {
		return
	}

	links := xurls.Relaxed.FindAllString(strings.Join(parts[1:], " "), -1)

	bot.Mention(channel, source, fmt.Sprintf("found links %s", strings.Join(links, ",")))
}
