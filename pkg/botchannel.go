package pkg

import "github.com/pajbot/pajbot2/pkg/eventemitter"

type MessageSender interface {
	Say(string)
	Mention(User, string)

	// Moderation
	Timeout(User, int, string)
	Ban(User, string)
}

type BotChannel interface {
	// Implement Channel interface
	GetName() string
	GetID() string

	MessageSender

	DatabaseID() int64
	Channel() Channel
	ChannelID() string
	ChannelName() string

	EnableModule(string) error
	DisableModule(string) error
	GetModule(string) (Module, error)

	// Implement ChannelWithStream interface
	Stream() Stream

	Events() *eventemitter.EventEmitter

	HandleMessage(user User, message Message) error
	HandleEventSubNotification(notification TwitchEventSubNotification) error
	OnModules(cb func(module Module) Actions, stop bool) []Actions

	Bot() Sender
}
