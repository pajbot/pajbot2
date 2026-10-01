package pkg

type BaseEvent struct {
	UserStore UserStore
}

type MessageSource interface {
	isMessageSource()
}

// ChannelMessageSource means the message came from a Twitch PRIVMSG
type ChannelMessageSource struct {
	Channel ChannelWithStream
}

func (ChannelMessageSource) isMessageSource() {}

// WhisperMessageSource means the message came from a Twitch whisper
type WhisperMessageSource struct {
	User User
}

func (WhisperMessageSource) isMessageSource() {}

type MessageEvent struct {
	BaseEvent

	User    User
	Message Message
	Channel ChannelWithStream
	Source  MessageSource
}

type EventSubNotificationEvent struct {
	BaseEvent

	Notification TwitchEventSubNotification
}
