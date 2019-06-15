package getuserid

import (
	"strings"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
	"github.com/pajbot/utils"
)

type Command struct {
}

func (c Command) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	usernames := utils.FilterUsernames(parts[1:])

	if len(usernames) == 0 {
		return twitchactions.Mention(event.User, "usage: !userid USERNAME (i.e. !userid pajlada)")
	}

	userIDs := event.UserStore.GetIDs(usernames)
	var results []string
	for username, userID := range userIDs {
		results = append(results, username+"="+userID)
	}

	if len(results) == 0 {
		return twitchactions.Mention(event.User, "no valid usernames were given")
	}

	return twitchactions.Mention(event.User, strings.Join(results, ", "))
}
