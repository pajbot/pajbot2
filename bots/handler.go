package bots

import twitch "github.com/gempir/go-twitch-irc"

// Handler xD
type Handler interface {
	HandleMessage(string, twitch.User, *TwitchMessage)
}

// HandlerFunc xD
type HandlerFunc func(string, twitch.User, *TwitchMessage)

// HandleMessage ADAPTER xD
func (f HandlerFunc) HandleMessage(channel string, user twitch.User, message *TwitchMessage) {
	f(channel, user, message)
}
