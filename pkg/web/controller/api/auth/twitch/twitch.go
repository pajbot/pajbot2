package twitch

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/dankeroni/gotwitch"
	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/apirequest"
	"github.com/pajlada/pajbot2/pkg/utils"
	"github.com/pajlada/pajbot2/pkg/web/router"
	"github.com/pajlada/pajbot2/pkg/web/state"
	"golang.org/x/oauth2"
)

type stateData struct {
	redirect string
}

var statesMutex sync.Mutex
var states = make(map[string]*stateData)

func stateExists(state string) bool {
	statesMutex.Lock()
	defer statesMutex.Unlock()

	_, ok := states[state]
	return ok
}

func stateExistsClear(state string) (*stateData, bool) {
	statesMutex.Lock()
	defer statesMutex.Unlock()

	d, ok := states[state]
	if ok {
		delete(states, state)
	}
	return d, ok
}

func makeState(redirectURL string) (string, error) {
	state, err := utils.GenerateRandomString(32)
	if err != nil {
		return "", err
	}

	if stateExists(state) {
		return "", errors.New("state already exists")
	}

	statesMutex.Lock()
	defer statesMutex.Unlock()

	states[state] = &stateData{
		redirect: redirectURL,
	}

	return state, nil
}

type authorizedCallback func(w http.ResponseWriter, r *http.Request,
	self gotwitch.ValidateResponse, oauth2Token *oauth2.Token,
	stateData *stateData)

func initializeOauthRoutes(ctx context.Context, m *mux.Router, config *oauth2.Config, name string, onAuthorized authorizedCallback) {
	router.RGet(m, "/"+name, func(w http.ResponseWriter, r *http.Request) {
		requestStateString, err := makeState(r.FormValue("redirect"))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.Redirect(w, r, config.AuthCodeURL(requestStateString), http.StatusFound)
	})

	router.RGet(m, "/"+name+"/callback",
		func(w http.ResponseWriter, r *http.Request) {
			stateData, exists := stateExistsClear(r.FormValue("state"))
			if !exists {
				http.Error(w, "Invalid OAuth state", http.StatusInternalServerError)
				return
			}

			// Get code
			oauth2Token, err := config.Exchange(ctx, r.URL.Query().Get("code"))
			if err != nil {
				http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
				return
			}

			validateResponse, _, err := apirequest.Twitch.ValidateOAuthTokenSimple(oauth2Token.AccessToken)
			if err != nil {
				http.Error(w, "Error validating token: "+err.Error(), http.StatusInternalServerError)
				return
			}

			onAuthorized(w, r, *validateResponse, oauth2Token, stateData)
		})
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "twitch root")
}

func Load(parent *mux.Router, a pkg.Application) error {
	auths := a.TwitchAuths()
	m := parent.PathPrefix("/twitch").Subrouter()
	ctx := context.Background()

	router.RGet(m, "", root)

	initializeOauthRoutes(ctx, m, auths.Bot(), "bot", onBotAuthenticated)
	initializeOauthRoutes(ctx, m, auths.User(), "user", onUserAuthenticated)
	initializeOauthRoutes(ctx, m, auths.Streamer(), "streamer", onStreamerAuthenticated)

	return nil
}

func onBotAuthenticated(w http.ResponseWriter, r *http.Request, self gotwitch.ValidateResponse, oauth2Token *oauth2.Token, stateData *stateData) {
	const queryF = `
INSERT INTO Bot
	(twitch_userid, twitch_username, twitch_access_token, twitch_refresh_token, twitch_access_token_expiry)
	VALUES (?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE twitch_username=?, twitch_access_token=?, twitch_refresh_token=?, twitch_access_token_expiry=?
	`
	c := state.Context(w, r)
	_, err := c.SQL.Exec(queryF, self.UserID, self.Login, oauth2Token.AccessToken, oauth2Token.RefreshToken, oauth2Token.Expiry,
		self.Login, oauth2Token.AccessToken, oauth2Token.RefreshToken, oauth2Token.Expiry)
	if err != nil {
		w.Write([]byte("Unable to insert bot :rage:"))
		return
	}

	w.Write([]byte("Bot added/updated! Restart the bot for the changes to take effect"))
}

func onUserAuthenticated(w http.ResponseWriter, r *http.Request, self gotwitch.ValidateResponse, oauth2Token *oauth2.Token, stateData *stateData) {
	c := state.Context(w, r)
	twitchUserName := self.Login

	twitchUserID := c.TwitchUserStore.GetID(twitchUserName)

	if twitchUserID == "" {
		// TODO: Fix proper error handling
		return
	}

	const queryF = `
INSERT INTO User
	(twitch_username, twitch_userid)
VALUES (?, ?)
	ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id), twitch_username=?
	`

	res, err := c.SQL.Exec(queryF, twitchUserName, twitchUserID, twitchUserName)
	if err != nil {
		panic(err)
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Last insert ID:", lastInsertID)

	sessionID, err := c.CreateSession(lastInsertID)
	if err != nil {
		panic(err)
	}
	state.SetSessionCookies(w, sessionID, twitchUserName)

	// TODO: Secure the redirect
	if stateData.redirect != "" {
		http.Redirect(w, r, stateData.redirect, http.StatusFound)

	}
}

func onStreamerAuthenticated(w http.ResponseWriter, r *http.Request, self gotwitch.ValidateResponse, oauth2Token *oauth2.Token, stateData *stateData) {
	// fmt.Printf("STREAMER Username: %s - Access token: %s\n", self.Token.UserName, oauth2Token.AccessToken)
}
