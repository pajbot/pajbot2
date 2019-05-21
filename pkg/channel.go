package pkg

type Channel interface {
	GetName() string
	GetID() string
}

type ChannelStore interface {
	TwitchChannel(channelID string) Channel
	RegisterTwitchChannel(channel Channel)
}
