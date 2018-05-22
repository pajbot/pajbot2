package pkg

import twitch "github.com/gempir/go-twitch-irc"

type Channel interface {
	Say(string, string)
	Timeout(string, twitch.User, int, string)
}
