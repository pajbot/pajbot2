package common

// this should make things easier with redis

import "time"

type User struct {
	ID          int
	Name        string
	DisplayName string
	Color       string //do we even want colors ?
	Mod         bool
	Sub         bool
	Turbo       bool
	Type        string // admin , staff etc
	Level       int
	Points      int
	LastSeen    time.Time // should this be time.Time or int/float?
	LastActive  time.Time
}

/*
Msg contains all the information about an IRC message.
This included already-parsed ircv3-tags and the User object
*/
type Msg struct {
	User    User
	Message string
	Channel string
	Type    string // PRIVMSG , WHISPER, (SUB?)
	Length  int    // will be set by a module or length of resub
	Me      bool
	Emotes  []Emote
}

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
	Pass       string `json:"pass"`
	Nick       string `json:"nick"`
	BrokerPort string `json:"broker_port"`

	RedisHost     string `json:"redis_host"`
	RedisPassword string `json:"redis_password"`

	TLSKey  string `json:"tls_key"`
	TLSCert string `json:"tls_cert"`

	Channels []string `json:"channels"`

	Quit chan string

	ToWeb   chan map[string]interface{}
	FromWeb chan map[string]interface{}
}
