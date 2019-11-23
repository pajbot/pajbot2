package twitch

import (
	"fmt"
	"net/http"

	"github.com/dankeroni/gotwitch"
	"github.com/pajbot/pajbot2/pkg/web/state"
	"golang.org/x/oauth2"
)

func onUserAuthenticated(w http.ResponseWriter, r *http.Request, self gotwitch.ValidateResponse, oauth2Token *oauth2.Token, stateData *stateData) {
	c := state.Context(w, r)
	twitchUserName := self.Login

	twitchUserID := c.TwitchUserStore.GetID(twitchUserName)

	if twitchUserID == "" {
		// TODO: Fix proper error handling
		return
	}

	const queryF = `
INSERT INTO "user"
	(twitch_username, twitch_userid)
VALUES ($1, $2)
	ON CONFLICT (twitch_userid) DO UPDATE SET twitch_username=$1
	RETURNING id
	`

	var lastInsertID int64
	row := c.SQL.QueryRow(queryF, twitchUserName, twitchUserID) // GOOD
	err := row.Scan(&lastInsertID)
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
