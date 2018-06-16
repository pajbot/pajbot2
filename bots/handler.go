package bots

import "github.com/pajlada/pajbot2/pkg"

// Handler xD
type Handler interface {
	HandleMessage(*TwitchBot, pkg.Channel, pkg.User, *TwitchMessage, pkg.Action)
}

// HandlerFunc xD
type HandlerFunc func(*TwitchBot, pkg.Channel, pkg.User, *TwitchMessage, pkg.Action)

// HandleMessage ADAPTER xD
func (f HandlerFunc) HandleMessage(bot *TwitchBot, channel pkg.Channel, user pkg.User, message *TwitchMessage, action pkg.Action) {
	f(bot, channel, user, message, action)
}
