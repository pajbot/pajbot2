package common

// this should make things easier with redis

import (
	"time"

	"github.com/pajlada/pajbot2/plog"
)

var log = plog.GetLogger()

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
)

/*
Msg contains all the information about an IRC message.
This included already-parsed ircv3-tags and the User object
*/
type Msg struct {
	User    User
	Text    string
	Channel string
	Type    MsgType // PRIVMSG , WHISPER, (SUB?)
	Me      bool
	Emotes  []Emote
	Tags    map[string]string
	Args    []string // needed for bot.Format for now
}

// Emote xD
type Emote struct {
	Name  string
	ID    string
	Type  string // bttv, twitch
	SizeX int    //in px
	SizeY int    //in px
	IsGif bool
	Count int
}
