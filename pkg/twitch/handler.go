package twitch

import "github.com/pajlada/pajbot2/pkg"

// Handler xD
type Handler interface {
	HandleMessage(*Bot, pkg.Channel, pkg.User, *TwitchMessage, pkg.Action)
}

// HandlerFunc xD
type HandlerFunc func(*Bot, pkg.Channel, pkg.User, *TwitchMessage, pkg.Action)

// HandleMessage ADAPTER xD
func (f HandlerFunc) HandleMessage(bot *Bot, channel pkg.Channel, user pkg.User, message *TwitchMessage, action pkg.Action) {
	f(bot, channel, user, message, action)
}
