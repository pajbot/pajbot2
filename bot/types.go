package bot

import "time"

/*
A User is a twitch user.

level system
0 banned
<50 completely ignored (other channel bots)
< 100 not responding to commands
100 pleb
250 sub
>250 ban immune
500 mod
1000 manager/channel admin
1500 broadcaster
2000 admin

TODO: Where do we put this?
*/
type User struct {
	ID          int
	Name        string
	DisplayName string
	Level       int
	Points      int
	LastSeen    time.Time
	Banned      bool
}

/*
Emote is an XXX xD

TODO: Where do we put this?
*/
type Emote struct {
	EmoteType string   // twitch / bttv / ffz
	ID        string   // twitchID / bttv hash
	Name      string   // Kappa KKona ...
	Pos       []string // [2-5, 7-12]
	Count     int
}
