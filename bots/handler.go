package bots

import "github.com/pajlada/pajbot2/pkg"

// Handler xD
type Handler interface {
	HandleMessage(*TwitchBot, string, pkg.User, *TwitchMessage)
}

// HandlerFunc xD
type HandlerFunc func(*TwitchBot, string, pkg.User, *TwitchMessage)

// HandleMessage ADAPTER xD
func (f HandlerFunc) HandleMessage(bot *TwitchBot, channel string, user pkg.User, message *TwitchMessage) {
	f(bot, channel, user, message)
}
