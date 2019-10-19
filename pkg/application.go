package pkg

import (
	"database/sql"
)

// Application is an instance of pajbot2
// It's responsible for initializing all bot accounts (`Bot` class)
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
