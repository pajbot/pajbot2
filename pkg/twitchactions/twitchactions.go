package twitchactions

import (
	"fmt"
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

func Mentionf(user pkg.User, format string, a ...interface{}) *Actions {
	actions := &Actions{}
	actions.Mention(user, fmt.Sprintf(format, a...))
	return actions
}

func Say(content string) *Actions {
	actions := &Actions{}
	actions.Say(content)
	return actions
}

func Sayf(format string, a ...interface{}) *Actions {
	actions := &Actions{}
	actions.Say(fmt.Sprintf(format, a...))
	return actions
}

func DoWhisper(user pkg.User, content string) *Actions {
	actions := &Actions{}
	actions.Whisper(user, content)
	return actions
}

func DoWhisperf(user pkg.User, format string, a ...interface{}) *Actions {
	actions := &Actions{}
	actions.Whisper(user, fmt.Sprintf(format, a...))
	return actions
}

func DoTimeout(user pkg.User, duration time.Duration, reason string) *Actions {
	actions := &Actions{}
	actions.Timeout(user, duration).SetReason(reason)
	return actions
}
