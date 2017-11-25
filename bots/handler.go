package bots

import twitch "github.com/gempir/go-twitch-irc"

// Handler xD
type Handler interface {
	HandleMessage(*TwitchBot, string, twitch.User, *TwitchMessage)
}

// HandlerFunc xD
type HandlerFunc func(*TwitchBot, string, twitch.User, *TwitchMessage)

// HandleMessage ADAPTER xD
func (f HandlerFunc) HandleMessage(bot *TwitchBot, channel string, user twitch.User, message *TwitchMessage) {
	f(bot, channel, user, message)
}
