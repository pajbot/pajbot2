package common

// this should make things easier with redis

import (
	"fmt"
	"time"
)

var (
	// BuildTime is the time when the binary was built
	// filled in with ./build.sh (ldflags)
	BuildTime string

	BuildRelease string

	BuildHash string

	BuildBranch string
)

func Version() string {
	if BuildRelease == "git" {
		return fmt.Sprintf("%s@%s", BuildHash, BuildBranch)
	}

	return BuildRelease
}

// GlobalUser will only be used by boss to check if user is admin
// and to decide what channel to send the message to if its a whisper
type GlobalUser struct {
	LastActive time.Time
	Channel    string
	Level      int
}

// MsgType specifies the message's type, for example PRIVMSG or WHISPER
type MsgType uint32

// Various message types which describe what sort of message they are
const (
	MsgPrivmsg MsgType = iota + 1
	MsgWhisper
	MsgSub
	MsgThrowAway
	MsgUnknown
	MsgUserNotice
	MsgReSub
	MsgNotice
	MsgRoomState
	MsgSubsOn
	MsgSubsOff
	MsgSlowOn
	MsgSlowOff
	MsgR9kOn
	MsgR9kOff
	MsgHostOn
	MsgHostOff
	MsgTimeoutSuccess
	MsgReconnect
)
