package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

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

func (c GetUserID) Trigger(bot pkg.Sender, botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	usernames := utils.FilterUsernames(parts[1:])

	if len(usernames) == 0 {
		bot.Mention(channel, user, "usage: !userid USERNAME (i.e. !userid pajlada)")
		return
	}

	userIDs := bot.GetUserStore().GetIDs(usernames)
	var results []string
	for username, userID := range userIDs {
		results = append(results, username+"="+userID)
	}

	if len(results) == 0 {
		bot.Mention(channel, user, "no valid usernames were given")
		return
	}

	bot.Mention(channel, user, strings.Join(results, ", "))
}

type GetUserName struct {
}

func (c GetUserName) Trigger(bot pkg.Sender, botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	userIDs := utils.FilterUserIDs(parts[1:])

	if len(userIDs) == 0 {
		bot.Mention(channel, user, "usage: !username USERID (i.e. !username 11148817)")
		return
	}

	names := bot.GetUserStore().GetNames(userIDs)
	var results []string
	for userID, username := range names {
		results = append(results, userID+"="+username)
	}

	if len(results) == 0 {
		bot.Mention(channel, user, "no valid user ids were given")
		return
	}

	bot.Mention(channel, user, strings.Join(results, ", "))
}

type GetPoints struct {
}

func (c GetPoints) Trigger(bot pkg.Sender, botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	var potentialTarget string
	targetID := user.GetID()

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
		bot.Mention(channel, user, "you have "+strconv.FormatUint(points, 10)+" points")
	} else {
		bot.Mention(channel, user, potentialTarget+" has "+strconv.FormatUint(points, 10)+" points")
	}
}

type AddPoints struct {
}

func (c AddPoints) Trigger(bot pkg.Sender, botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	_, points := bot.AddPoints(channel, user.GetID(), uint64(rand.Int31n(50)))
	bot.Mention(channel, user, "you now have "+strconv.FormatUint(points, 10)+" points")
}

type RemovePoints struct {
}

func (c RemovePoints) Trigger(bot pkg.Sender, botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	_, points := bot.RemovePoints(channel, user.GetID(), uint64(rand.Int31n(50)))
	bot.Mention(channel, user, "you now have "+strconv.FormatUint(points, 10)+" points")
}

type Roulette struct {
}

func (c Roulette) Trigger(bot pkg.Sender, botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	if len(parts) < 2 {
		bot.Mention(channel, user, "usage: !roulette 500 or !roulette all")
		return
	}

	var pointsToRoulette uint64

	if strings.ToLower(parts[1]) == "all" {
		pointsToRoulette = bot.GetPoints(channel, user.GetID())
	} else {
		var err error
		pointsToRoulette, err = strconv.ParseUint(parts[1], 10, 64)

		if err != nil {
			bot.Mention(channel, user, "usage: !roulette 500 or !roulette all")
			return
		}
	}

	if pointsToRoulette == 0 {
		bot.Mention(channel, user, "you have 0 points, you can't roulette ResidentSleeper")
		return
	}

	if result, _ := bot.RemovePoints(channel, user.GetID(), pointsToRoulette); !result {
		bot.Mention(channel, user, "you don't have enough points ResidentSleeper")
		return
	}

	if rand.Int31n(2) == 0 {
		// loss
		bot.Mention(channel, user, "you lost OMEGALUL")
	} else {
		// win
		// TODO: Check for integer overflow?
		_, newPoints := bot.AddPoints(channel, user.GetID(), pointsToRoulette*2)
		bot.Mention(channel, user, "you won PagChomp you now have "+strconv.FormatUint(newPoints, 10)+" points KKona")
	}
}

type GivePoints struct {
}

func (c GivePoints) Trigger(bot pkg.Sender, botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	const USAGE = `usage: !givepoints USER POINTS`
	if len(parts) < 3 {
		bot.Mention(channel, user, USAGE)
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
		pointsToGive = bot.GetPoints(channel, user.GetID())
	} else {
		var err error
		pointsToGive, err = strconv.ParseUint(parts[2], 10, 64)

		if err != nil {
			bot.Mention(channel, user, USAGE)
			return
		}
	}

	if pointsToGive == 0 {
		bot.Mention(channel, user, USAGE)
		return
	}

	if result, _ := bot.RemovePoints(channel, user.GetID(), pointsToGive); !result {
		bot.Mention(channel, user, "you don't have enough points ResidentSleeper")
		return
	}

	bot.AddPoints(channel, targetID, pointsToGive)
	bot.Mention(channel, user, "you gave away "+strconv.FormatUint(pointsToGive, 10)+" points to @"+target)
}

type Ping struct {
}

func (c Ping) Trigger(bot pkg.Sender, botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	bot.Mention(channel, user, fmt.Sprintf("pb2 has been running for %s", time.Since(startTime)))
}

type Simplify struct {
}

func (c Simplify) Trigger(bot pkg.Sender, botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	if user.IsModerator() || user.IsBroadcaster(channel) {
		if len(parts) > 1 {
			normalizedMessage, err := normalize.Normalize(strings.Join(parts[1:], " "))
			if err != nil {
				bot.Mention(channel, user, fmt.Sprintf("error normalizing string: %s", err.Error()))
				return
			}

			bot.Mention(channel, user, fmt.Sprintf("normalized string: '%s'", normalizedMessage))
		}
	}
}

type TimeMeOut struct {
}

func (c TimeMeOut) Trigger(bot pkg.Sender, botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	if len(parts) < 2 {
		return
	}

	timeoutDuration, err := time.ParseDuration(parts[1])
	if err != nil {
		bot.Mention(channel, user, "invalid duration format. use !timemeout 1s or !timemeout 5m")
		return
	}

	var reason string

	if len(parts) > 2 {
		reason = strings.Join(parts[2:], " ")
	}

	bot.Timeout(channel, user, int(timeoutDuration.Seconds()), reason)
}

type Join struct {
}

func (c *Join) Trigger(bot pkg.Sender, botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	if !user.HasGlobalPermission(pkg.PermissionAdmin) {
		bot.Mention(channel, user, "you do not have permission to use this command. Admin permission is required")
		return
	}

	if len(parts) < 2 {
		return
	}

	channelName := parts[1]

	if strings.EqualFold(channelName, bot.Name()) {
		bot.Mention(channel, user, "I cannot join my own channel")
		return
	}

	channelID := bot.GetUserStore().GetID(channelName)
	if channelID == "" {
		bot.Mention(channel, user, "no channel with that name exists")
		return
	}

	err := bot.JoinChannel(channelID)
	if err != nil {
		bot.Mention(channel, user, err.Error())
		return
	}

	bot.Mention(channel, user, fmt.Sprintf("joined channel %s(%s)", channelName, channelID))
}

type Leave struct {
}

func (c *Leave) Trigger(bot pkg.Sender, botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	if !user.HasGlobalPermission(pkg.PermissionAdmin) {
		bot.Mention(channel, user, "you do not have permission to use this command. Admin permission is required")
		return
	}

	if len(parts) < 2 {
		return
	}

	channelName := parts[1]

	if strings.EqualFold(channelName, bot.Name()) {
		bot.Mention(channel, user, "I cannot leave my own channel")
		return
	}

	channelID := bot.GetUserStore().GetID(channelName)
	if channelID == "" {
		bot.Mention(channel, user, "no channel with that name exists")
		return
	}

	err := bot.LeaveChannel(channelID)
	if err != nil {
		bot.Mention(channel, user, err.Error())
		return
	}

	bot.Mention(channel, user, fmt.Sprintf("left channel %s(%s)", channelName, channelID))
}

type Test struct {
}

func (c Test) Trigger(bot pkg.Sender, botChannel pkg.BotChannel, parts []string, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) {
	if !user.IsModerator() && !user.IsBroadcaster(channel) {
		return
	}

	if len(parts) <= 1 {
		return
	}

	variations, _, err := utils.MakeVariations(strings.Join(parts[1:], " "), true)
	if err != nil {
		bot.Mention(channel, user, err.Error())
		return
	}

	for _, variation := range variations {
		bot.Mention(channel, user, fmt.Sprintf("variation %s", variation))
	}
}
