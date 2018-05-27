package bots

import (
	"fmt"

	twitch "github.com/gempir/go-twitch-irc"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/users"
	"github.com/pajlada/pajbot2/redismanager"
)

type botFlags struct {
	PermaSubMode bool
}

// TwitchBot is a wrapper around go-twitch-irc's twitch.Client with a few extra features
type TwitchBot struct {
	*twitch.Client

	Name    string
	handler Handler

	QuitChannel chan string

	Flags botFlags

	Redis *redismanager.RedisManager

	Modules []pkg.Module
}

// TwitchMessage is a wrapper for twitch.Message with some extra stuff
type TwitchMessage struct {
	twitch.Message

	BTTVEmotes []common.Emote
	// TODO: BTTV Emotes

	// TODO: FFZ Emotes

	// TODO: Emojis
}

func (b *TwitchBot) RegisterModules() error {
	for _, m := range b.Modules {
		err := m.Register()
		if err != nil {
			return err
		}
	}

	return nil
}

// Reply will reply to the message in the same way it received the message
// If the message was received in a twitch channel, reply in that twitch channel.
// IF the message was received in a twitch whisper, reply using twitch whispers.
func (b *TwitchBot) Reply(channel string, user pkg.User, message string) {
	if channel == "" {
		b.Client.Whisper(user.GetName(), message)
	} else {
		b.Client.Say(channel, message)
	}
}

func (b *TwitchBot) Say(channel pkg.Channel, message string) {
	b.Client.Say(channel.GetChannel(), message)
}

func (b *TwitchBot) SaySimple(channel string, message string) {
	b.Client.Say(channel, message)
}

func (b *TwitchBot) Timeout(channel pkg.Channel, user pkg.User, duration int, reason string) {
	// Empty string in UserType means a normal user
	if !user.IsModerator() {
		b.Say(channel, fmt.Sprintf(".timeout %s %d %s", user.GetName(), duration, reason))
	}
}

// SetHandler sets the handler to message at the bottom of the list
func (b *TwitchBot) SetHandler(handler Handler) {
	b.handler = handler
}

// HandleMessage goes through all of the bot handlers in the correct order and figures out if anything was triggered
func (b *TwitchBot) HandleMessage(channel string, user twitch.User, message *TwitchMessage) {
	twitchUser := &users.TwitchUser{
		User: user,

		ID: message.Tags["user-id"],
	}

	b.handler.HandleMessage(b, channel, twitchUser, message)
}

// Quit quits the entire application
func (b *TwitchBot) Quit(message string) {
	b.QuitChannel <- message
}
