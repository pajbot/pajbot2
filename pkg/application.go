package pkg

import (
	"database/sql"
)

type Application interface {
	UserStore() UserStore
	UserContext() UserContext
	StreamStore() StreamStore
	SQL() *sql.DB
	PubSub() PubSub
	TwitchBots() BotStore
	QuitChannel() chan string
	TwitchAuths() TwitchAuths
}
