package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	normalize "github.com/pajlada/lidl-normalize"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/utils"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

type GetUserID struct {
}

func (c GetUserID) Trigger(botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	usernames := utils.FilterUsernames(parts[1:])

	if len(usernames) == 0 {
		botChannel.Mention(user, "usage: !userid USERNAME (i.e. !userid pajlada)")
		return
	}

	userIDs := botChannel.Bot().GetUserStore().GetIDs(usernames)
	var results []string
	for username, userID := range userIDs {
		results = append(results, username+"="+userID)
	}

	if len(results) == 0 {
		botChannel.Mention(user, "no valid usernames were given")
		return
	}

	botChannel.Mention(user, strings.Join(results, ", "))
}

type GetUserName struct {
}

func (c GetUserName) Trigger(botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	userIDs := utils.FilterUserIDs(parts[1:])

	if len(userIDs) == 0 {
		botChannel.Mention(user, "usage: !username USERID (i.e. !username 11148817)")
		return
	}

	names := botChannel.Bot().GetUserStore().GetNames(userIDs)
	var results []string
	for userID, username := range names {
		results = append(results, userID+"="+username)
	}

	if len(results) == 0 {
		botChannel.Mention(user, "no valid user ids were given")
		return
	}

	botChannel.Mention(user, strings.Join(results, ", "))
}

type GetPoints struct {
}

func (c GetPoints) Trigger(botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	var potentialTarget string
	targetID := user.GetID()

	if len(parts) >= 2 {
		potentialTarget = utils.FilterUsername(parts[1])
		if potentialTarget != "" {
			potentialTargetID := botChannel.Bot().GetUserStore().GetID(potentialTarget)
			if potentialTargetID != "" {
				targetID = potentialTargetID
			} else {
				potentialTarget = ""
			}
		}
	}

	points := botChannel.Bot().GetPoints(channel, targetID)
	if potentialTarget == "" {
		botChannel.Mention(user, "you have "+strconv.FormatUint(points, 10)+" points")
	} else {
		botChannel.Mention(user, potentialTarget+" has "+strconv.FormatUint(points, 10)+" points")
	}
}

type AddPoints struct {
}

func (c AddPoints) Trigger(botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	_, points := botChannel.Bot().AddPoints(channel, user.GetID(), uint64(rand.Int31n(50)))
	botChannel.Mention(user, "you now have "+strconv.FormatUint(points, 10)+" points")
}

type RemovePoints struct {
}

func (c RemovePoints) Trigger(botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	_, points := botChannel.Bot().RemovePoints(channel, user.GetID(), uint64(rand.Int31n(50)))
	botChannel.Mention(user, "you now have "+strconv.FormatUint(points, 10)+" points")
}

type Roulette struct {
}

func (c Roulette) Trigger(botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	if len(parts) < 2 {
		botChannel.Mention(user, "usage: !roulette 500 or !roulette all")
		return
	}

	var pointsToRoulette uint64

	if strings.ToLower(parts[1]) == "all" {
		pointsToRoulette = botChannel.Bot().GetPoints(channel, user.GetID())
	} else {
		var err error
		pointsToRoulette, err = strconv.ParseUint(parts[1], 10, 64)

		if err != nil {
			botChannel.Mention(user, "usage: !roulette 500 or !roulette all")
			return
		}
	}

	if pointsToRoulette == 0 {
		botChannel.Mention(user, "you have 0 points, you can't roulette ResidentSleeper")
		return
	}

	if result, _ := botChannel.Bot().RemovePoints(channel, user.GetID(), pointsToRoulette); !result {
		botChannel.Mention(user, "you don't have enough points ResidentSleeper")
		return
	}

	if rand.Int31n(2) == 0 {
		// loss
		botChannel.Mention(user, "you lost OMEGALUL")
	} else {
		// win
		// TODO: Check for integer overflow?
		_, newPoints := botChannel.Bot().AddPoints(channel, user.GetID(), pointsToRoulette*2)
		botChannel.Mention(user, "you won PagChomp you now have "+strconv.FormatUint(newPoints, 10)+" points KKona")
	}
}

type GivePoints struct {
}

func (c GivePoints) Trigger(botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	const USAGE = `usage: !givepoints USER POINTS`
	if len(parts) < 3 {
		botChannel.Mention(user, USAGE)
		return
	}

	target := utils.FilterUsername(parts[1])
	if target == "" {
		// Invalid username
		return
	}

	targetID := botChannel.Bot().GetUserStore().GetID(target)
	if targetID == "" {
		// Invalid username
		return
	}

	var pointsToGive uint64

	if strings.ToLower(parts[2]) == "all" {
		pointsToGive = botChannel.Bot().GetPoints(channel, user.GetID())
	} else {
		var err error
		pointsToGive, err = strconv.ParseUint(parts[2], 10, 64)

		if err != nil {
			botChannel.Mention(user, USAGE)
			return
		}
	}

	if pointsToGive == 0 {
		botChannel.Mention(user, USAGE)
		return
	}

	if result, _ := botChannel.Bot().RemovePoints(channel, user.GetID(), pointsToGive); !result {
		botChannel.Mention(user, "you don't have enough points ResidentSleeper")
		return
	}

	botChannel.Bot().AddPoints(channel, targetID, pointsToGive)
	botChannel.Mention(user, "you gave away "+strconv.FormatUint(pointsToGive, 10)+" points to @"+target)
}

type Simplify struct {
}

func (c Simplify) Trigger(botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	if user.IsModerator() || user.IsBroadcaster(channel) {
		if len(parts) > 1 {
			normalizedMessage, err := normalize.Normalize(strings.Join(parts[1:], " "))
			if err != nil {
				botChannel.Mention(user, fmt.Sprintf("error normalizing string: %s", err.Error()))
				return
			}

			botChannel.Mention(user, fmt.Sprintf("normalized string: '%s'", normalizedMessage))
		}
	}
}

type TimeMeOut struct {
}

func (c TimeMeOut) Trigger(botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	if len(parts) < 2 {
		return
	}

	timeoutDuration, err := time.ParseDuration(parts[1])
	if err != nil {
		botChannel.Mention(user, "invalid duration format. use !timemeout 1s or !timemeout 5m")
		return
	}

	var reason string

	if len(parts) > 2 {
		reason = strings.Join(parts[2:], " ")
	}

	botChannel.Timeout(user, int(timeoutDuration.Seconds()), reason)
}

type Test struct {
}

func (c Test) Trigger(botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	if !user.IsModerator() && !user.IsBroadcaster(channel) {
		return
	}

	if len(parts) <= 1 {
		return
	}

	variations, _, err := utils.MakeVariations(strings.Join(parts[1:], " "), true)
	if err != nil {
		botChannel.Mention(user, err.Error())
		return
	}

	for _, variation := range variations {
		botChannel.Mention(user, fmt.Sprintf("variation %s", variation))
	}
}

type IsLive struct {
}

func (c IsLive) Trigger(botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	if !user.IsModerator() && !user.IsBroadcaster(channel) {
		return
	}

	if botChannel.Stream().Status().Live() {
		startedAt := botChannel.Stream().Status().StartedAt()
		botChannel.Mention(user, fmt.Sprintf("LIVE FOR %s KKona", utils.TimeSince(startedAt)))
	} else {
		botChannel.Mention(user, "offline FeelsBadMan")
	}
}
