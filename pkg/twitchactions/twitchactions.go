package twitchactions

import (
	"time"

	"github.com/pajbot/pajbot2/pkg"
)

var _ pkg.Actions = &Actions{}

// Actions lists
type Actions struct {
	baseActions
}

/// HELPER FUNCTIONS FOR SINGLE ACTIONS

func Mention(user pkg.User, content string) *Actions {
	actions := &Actions{}
	actions.Mention(user, content)
	return actions
}

func Say(content string) *Actions {
	actions := &Actions{}
	actions.Say(content)
	return actions
}

func DoWhisper(user pkg.User, content string) *Actions {
	actions := &Actions{}
	actions.Whisper(user, content)
	return actions
}

func DoTimeout(user pkg.User, duration time.Duration, reason string) *Actions {
	actions := &Actions{}
	actions.Timeout(user, duration).SetReason(reason)
	return actions
}
