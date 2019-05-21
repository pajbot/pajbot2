package state

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/utils"
)

var (
	sqlClient          *sql.DB
	twitchUserStore    pkg.UserStore
	twitchChannelStore pkg.ChannelStore
	pubSub             pkg.PubSub
	application        pkg.Application

	mutex = &sync.RWMutex{}

	sessionStore = &SessionStore{}
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

func StoreTwitchChannelStore(twitchChannelStore_ pkg.ChannelStore) {
	mutex.Lock()
	twitchChannelStore = twitchChannelStore_
	mutex.Unlock()
}

func StorePubSub(pubSub_ pkg.PubSub) {
	mutex.Lock()
	pubSub = pubSub_
	mutex.Unlock()
}

func StoreApplication(application_ pkg.Application) {
	mutex.Lock()
	application = application_
	mutex.Unlock()
}

type Session struct {
	ID     string
	UserID uint64

	TwitchUserID   string
	TwitchUserName string
}

type State struct {
	SQL                *sql.DB
	TwitchUserStore    pkg.UserStore
	TwitchChannelStore pkg.ChannelStore
	PubSub             pkg.PubSub
	Application        pkg.Application
	Session            *Session
	SessionID          *string

	// Channel is filled in if the user has provided a channel_id argument with the request.
	// Filled in automatically with the Context function
	Channel pkg.Channel
}

func (s *State) CreateSession(userID int64) (sessionID string, err error) {
	// TODO: Switch this out for a proper randomizer or something KKona
	sessionID, err = utils.GenerateRandomString(64)
	if err != nil {
		return
	}

	const queryF = `
INSERT INTO
	UserSession
(id, user_id)
	VALUES (?, ?)`

	// TODO: Make sure the exec didn't error
	_, err = sqlClient.Exec(queryF, sessionID, userID)
	if err != nil {
		return
	}

	return
}

func Context(w http.ResponseWriter, r *http.Request) State {
	mutex.RLock()
	state := State{
		SQL:                sqlClient,
		TwitchUserStore:    twitchUserStore,
		TwitchChannelStore: twitchChannelStore,
		PubSub:             pubSub,
		Application:        application,
	}
	mutex.RUnlock()

	// Authorization via header api key (not implemented yet)
	credentials := r.Header.Get("Authorization")
	if credentials != "" {
		fmt.Println("Credentials:", credentials)
	}

	state.SessionID = getCookie(r, SessionIDCookie)
	if state.SessionID != nil && *state.SessionID != "" {
		state.Session = sessionStore.Get(state.SQL, *state.SessionID)
	}

	// Figure out channel context if this request has one
	vars := mux.Vars(r)

	if channelID, ok := vars["channel_id"]; ok {
		state.Channel = state.TwitchChannelStore.TwitchChannel(channelID)
	}

	return state
}
