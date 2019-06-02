package pkg

import (
	"database/sql"
)

type Application interface {
	UserStore() UserStore
	ChannelStore() ChannelStore
	UserContext() UserContext
	StreamStore() StreamStore
	SQL() *sql.DB
	PubSub() PubSub
	TwitchBots() BotStore
	QuitChannel() chan string
	TwitchAuths() TwitchAuths
	MIMO() MIMO
}

type MIMO interface {
	Subscriber(channelNames ...string) chan interface{}
	Publisher(channelName string) chan interface{}
}
