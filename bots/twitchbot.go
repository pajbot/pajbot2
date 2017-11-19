package bots

import (
	twitch "github.com/gempir/go-twitch-irc"
	"github.com/pajlada/pajbot2/common"
)

// TwitchBot is a wrapper around go-twitch-irc's twitch.Client with a few extra features
type TwitchBot struct {
	*twitch.Client

	handlers []Handler
}

// TwitchMessage is a wrapper for twitch.Message with some extra stuff
type TwitchMessage struct {
	twitch.Message

	BTTVEmotes []common.Emote
	// TODO: BTTV Emotes

	// TODO: FFZ Emotes

	// TODO: Emojis
}

// AddHandler adds a handler to message at the bottom of the list
func (b *TwitchBot) AddHandler(handler Handler) {
	b.handlers = append(b.handlers, handler)
}

// HandleMessage goes through all of the bot handlers in the correct order and figures out if anything was triggered
func (b *TwitchBot) HandleMessage(channel string, user twitch.User, message *TwitchMessage) {
	for _, handler := range b.handlers {
		handler.HandleMessage(channel, user, message)
	}
}
