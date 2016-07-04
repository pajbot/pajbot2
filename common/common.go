package common

// this should make things easier with redis

import "time"

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

/*
The Config contains all the data required to connect
to the twitch IRC servers
*/
type Config struct {
	Pass       string  `json:"pass"`
	Nick       string  `json:"nick"`
	BrokerHost *string `json:"broker_host"`
	BrokerPass *string `json:"broker_pass"`
	Silent     bool    `json:"silent"`

	RedisHost     string `json:"redis_host"`
	RedisPassword string `json:"redis_password"`
	RedisDatabase int    `json:"redis_database"`

	WebHost   string `json:"web_host"`
	WebDomain string `json:"web_domain"`

	SQLDSN string `json:"sql_dsn"`

	TLSKey  string `json:"tls_key"`
	TLSCert string `json:"tls_cert"`

	Channels []string `json:"channels"`

	Quit chan string

	ToWeb   chan map[string]interface{}
	FromWeb chan map[string]interface{}
}
