package pkg

type BotChannel interface {
	DatabaseID() int64
	ChannelID() string
	ChannelName() string
}
