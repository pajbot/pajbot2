package bot

import "time"

type BotConfig struct {
	Readchan chan Msg
	Sendchan chan string
	Channel  string
}

type Command struct {
	Id       int
	Trigger  []string
	Cooldown map[string]time.Time
	Response string
	Level    int
}

type User struct {
	Id          int
	Name        string
	DisplayName string
	Level       int
	Points      int
	LastSeen    time.Time
	Banned      bool
}

type Msg struct {
	Color       string
	Displayname string
	Emotes      []Emote
	Mod         bool
	Subscriber  bool
	Turbo       bool
	Usertype    string
	Username    string
	Channel     string
	Message     string
	MessageType string
	Me          bool
	Length      int
}

type Emote struct {
	EmoteType string   // twitch / bttv / ffz
	Id        string   // twitchID / bttv hash
	Name      string   // Kappa KKona ...
	Pos       []string // [2-5, 7-12]
	Count     int
}

// level system
// 0 banned
// <50 completely ignored (other channel bots)
// < 100 not responding to commands
// 100 pleb
// 250 sub
// >250 ban immune
// 500 mod
// 1000 manager/channel admin
// 1500 broadcaster
// 2000 admin

type Bot struct {
	Read    chan Msg
	Send    chan string
	Channel string
}
