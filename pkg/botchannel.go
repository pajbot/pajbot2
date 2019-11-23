package pkg

import "github.com/pajbot/pajbot2/pkg/eventemitter"

type MessageSender interface {
	Say(string)
	Mention(User, string)

	// Moderation
	Timeout(User, int, string)
	SingleTimeout(User, int, string)
}

type BotChannel interface {
	MessageSender

	DatabaseID() int64
	Channel() Channel
	ChannelID() string
	ChannelName() string

	EnableModule(string) error
	DisableModule(string) error
	GetModule(string) (Module, error)

	Stream() Stream

	Events() *eventemitter.EventEmitter

	HandleMessage(user User, message Message) error
	OnModules(cb func(module Module) Actions, stop bool) []Actions

	Bot() Sender
}
