package pkg

// Channel is the most barebones way of accessing a Twitch channel
// For a Channel to 'live' we must be able to access its Name (Twitch User Login) and ID (Twitch User ID)
type Channel interface {
	GetName() string
	GetID() string
}

type ChannelWithStream interface {
	Channel

	Stream() Stream
}

type ChannelStore interface {
	TwitchChannel(channelID string) Channel
	RegisterTwitchChannel(channel Channel)
}
