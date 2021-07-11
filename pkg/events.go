package pkg

type BaseEvent struct {
	UserStore UserStore
}

type MessageEvent struct {
	BaseEvent

	User    User
	Message Message
	Channel ChannelWithStream
}

type EventSubNotificationEvent struct {
	BaseEvent

	Notification TwitchEventSubNotification
}
