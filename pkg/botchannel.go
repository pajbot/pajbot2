package pkg

import "github.com/pajlada/pajbot2/pkg/eventemitter"

type BotChannel interface {
	DatabaseID() int64
	Channel() Channel
	ChannelID() string
	ChannelName() string

	EnableModule(string) error
	DisableModule(string) error

	Stream() Stream

	Events() *eventemitter.EventEmitter

	Say(string)
	Mention(User, string)

	// Moderation
	Timeout(User, int, string)

	Bot() Sender
}
