package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/dankeroni/gotwitch"
	"github.com/pajlada/pajbot2/apirequest"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/utils"
)

type GetUserID struct {
}

func (c GetUserID) Trigger(bot pkg.Sender, parts []string, channel pkg.Channel, source pkg.User, message pkg.Message, action pkg.Action) {
	usernames := utils.FilterUsernames(parts)

	if len(usernames) == 0 {
		bot.Say(channel, "@"+source.GetName()+", usage: !userid USERNAME (i.e. !userid pajlada)")
		return
	}

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

	apirequest.Twitch.GetUsersByLogin(usernames, onSuccess, onHTTPError, onInternalError)
}

type GetPoints struct {
}

func (c GetPoints) Trigger(bot pkg.Sender, parts []string, channel pkg.Channel, source pkg.User, message pkg.Message, action pkg.Action) {
	points := bot.GetPoints(channel, source)
	bot.Mention(channel, source, "you have "+strconv.FormatUint(points, 10)+" points")
}

type AddPoints struct {
}

func (c AddPoints) Trigger(bot pkg.Sender, parts []string, channel pkg.Channel, source pkg.User, message pkg.Message, action pkg.Action) {
	points := bot.EditPoints(channel, source, rand.Int31n(50))
	bot.Mention(channel, source, "you now have "+strconv.FormatUint(points, 10)+" points")
}

type RemovePoints struct {
}

func (c RemovePoints) Trigger(bot pkg.Sender, parts []string, channel pkg.Channel, source pkg.User, message pkg.Message, action pkg.Action) {
	points := bot.EditPoints(channel, source, -rand.Int31n(50))
	bot.Mention(channel, source, "you now have "+strconv.FormatUint(points, 10)+" points")
}
