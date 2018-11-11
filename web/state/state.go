package state

import (
	"database/sql"
	"net/http"
	"sync"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/pubsub"
)

var (
	sqlClient       *sql.DB
	twitchUserStore pkg.UserStore
	pubSub          *pubsub.PubSub

	mutex = &sync.RWMutex{}
)

func StoreSQL(sql_ *sql.DB) {
	mutex.Lock()
	sqlClient = sql_
	mutex.Unlock()
}

func StoreTwitchUserStore(twitchUserStore_ pkg.UserStore) {
	mutex.Lock()
	twitchUserStore = twitchUserStore_
	mutex.Unlock()
}

func StorePubSub(pubSub_ *pubsub.PubSub) {
	mutex.Lock()
	pubSub = pubSub_
	mutex.Unlock()
}

type State struct {
	SQL             *sql.DB
	TwitchUserStore pkg.UserStore
	PubSub          *pubsub.PubSub
}

func Context(w http.ResponseWriter, r *http.Request) State {
	mutex.RLock()
	s := State{
		SQL:             sqlClient,
		TwitchUserStore: twitchUserStore,
		PubSub:          pubSub,
	}
	mutex.RUnlock()

	return s
}
