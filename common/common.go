package common

// this should make things easier with redis

import "time"

type User struct {
	ID          int
	Name        string
	Displayname string
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
