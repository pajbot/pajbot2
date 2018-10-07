package web

import (
	"context"
	"errors"
	"net/http"
	"sync"

	oidc "github.com/coreos/go-oidc"
	"github.com/dankeroni/gotwitch"
	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/pkg/apirequest"
	"github.com/pajlada/pajbot2/pkg/common"
	"github.com/pajlada/pajbot2/pkg/common/config"
	"github.com/pajlada/pajbot2/pkg/utils"
	"golang.org/x/oauth2"
)

var (
	twitchBotOauth      = &oauth2.Config{}
	twitchUserOauth     = &oauth2.Config{}
	twitchStreamerOauth = &oauth2.Config{}
)

type nonceData struct {
	str      string
	state    string
	redirect string
}

var nonces = make(map[string]nonceData)

var statesMutex sync.Mutex
var states = make(map[string]bool)

func stateExists(state string) bool {
	statesMutex.Lock()
	defer statesMutex.Unlock()

	_, ok := states[state]
	return ok
}

func stateExistsClear(state string) bool {
	statesMutex.Lock()
	defer statesMutex.Unlock()

	_, ok := states[state]
	if ok {
		delete(states, state)
	}
	return ok
}

func makeState() (string, error) {
	state, err := utils.GenerateRandomString(32)
	if err != nil {
		return "", err
	}

	if stateExists(state) {
		return "", errors.New("state already exists")
	}

	statesMutex.Lock()
	defer statesMutex.Unlock()

	states[state] = true

	return state, nil
}

type authorizedCallback func(w http.ResponseWriter, r *http.Request, self gotwitch.Self, oauth2Token *oauth2.Token, nonce *nonceData)

func testpenis(ctx context.Context, provider *oidc.Provider, m *mux.Router, config *oauth2.Config, appConfig *config.TwitchAuthConfig, name string, onAuthorized authorizedCallback) error {
	var nonceEnabledVerifier *oidc.IDTokenVerifier
	useOIDC := false

	if provider != nil {
		oidcConfig := &oidc.Config{
			ClientID:        appConfig.ClientID,
			SkipExpiryCheck: true,
		}

		// Use the nonce source to create a custom ID Token verifier.
		nonceEnabledVerifier = provider.Verifier(oidcConfig)
		useOIDC = true
	}

	m.HandleFunc("/auth/twitch/"+name, func(w http.ResponseWriter, r *http.Request) {
		nonce, err := utils.GenerateRandomString(32)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		nonces[nonce] = nonceData{
			str:      nonce,
			state:    "pending",
			redirect: r.FormValue("redirect"),
		}

		requestStateString, err := makeState()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		url := config.AuthCodeURL(requestStateString, oidc.Nonce(nonce))
		http.Redirect(w, r, url, http.StatusFound)
	})

	m.HandleFunc("/auth/twitch/"+name+"/callback",
		func(w http.ResponseWriter, r *http.Request) {
			if !stateExistsClear(r.FormValue("state")) {
				http.Error(w, "Invalid OAuth state", http.StatusInternalServerError)
				return
			}

			// Get code
			oauth2Token, err := config.Exchange(ctx, r.URL.Query().Get("code"))
			if err != nil {
				http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
				return
			}

			var nonceData *nonceData

			if useOIDC {
				// Extract id_token
				rawIDToken, ok := oauth2Token.Extra("id_token").(string)
				if !ok {
					http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
					return
				}

				// Verify the ID Token signature and nonce.
				idToken, err := nonceEnabledVerifier.Verify(ctx, rawIDToken)
				if err != nil {
					http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
					return
				}

				nd, ok := nonces[idToken.Nonce]
				if !ok || nd.state != "pending" {
					http.Error(w, "Invalid ID Token nonce", http.StatusInternalServerError)
					return
				}

				nonceData = &nd
			}

			onSuccess := func(data gotwitch.Self) {
				if !data.Identified || !data.Token.Valid {
					http.Error(w, "Token invalid (Twitch end)", http.StatusInternalServerError)
					return
				}

				onAuthorized(w, r, data, oauth2Token, nonceData)
			}

			apirequest.TwitchV3.GetSelf(oauth2Token.AccessToken, onSuccess, onHTTPError, onInternalError)

			// w.Write([]byte("penis"))
		})

	return nil
}

func twitchAuthInit(m *mux.Router, appConfig *config.AuthTwitchConfig) (err error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://id.twitch.tv/oauth2")
	if err != nil {
		panic(err)
	}

	err = testpenis(ctx, nil, m, twitchBotOauth, &appConfig.Bot, "bot", func(w http.ResponseWriter, r *http.Request, self gotwitch.Self, oauth2Token *oauth2.Token, nonce *nonceData) {
		err := common.CreateBot(sqlClient, self.Token.UserName, oauth2Token.AccessToken, oauth2Token.RefreshToken)
		if err != nil {
			// TODO: Handle gracefully
			panic(err)
		}
	})
	if err != nil {
		return
	}
	err = testpenis(ctx, provider, m, twitchUserOauth, &appConfig.User, "user", func(w http.ResponseWriter, r *http.Request, self gotwitch.Self, oauth2Token *oauth2.Token, nonce *nonceData) {
		name := self.Token.UserName
		id := twitchUserStore.GetID(name)

		const queryF = `
INSERT INTO User
	(twitch_username, twitch_userid, twitch_nonce)
VALUES (?, ?, ?)
	ON DUPLICATE KEY UPDATE twitch_username=?, twitch_nonce=?
	`

		res, err := sqlClient.Exec(queryF, name, id, nonce.str, name, nonce.str)
		if err != nil {
			panic(err)
		}

		affectedRowsCount, err := res.RowsAffected()
		if err != nil {
			panic(err)
		}

		if affectedRowsCount == 1 {
			// User inserted
		} else {
			// User updated
		}

		// TODO: Secure the redirect
		if nonce.redirect != "" {
			http.Redirect(w, r, nonce.redirect+"#nonce="+nonce.str+";user_id="+id, http.StatusFound)
		}
	})
	if err != nil {
		return
	}
	err = testpenis(ctx, nil, m, twitchStreamerOauth, &appConfig.Streamer, "streamer", func(w http.ResponseWriter, r *http.Request, self gotwitch.Self, oauth2Token *oauth2.Token, nonce *nonceData) {
		// fmt.Printf("STREAMER Username: %s - Access token: %s\n", self.Token.UserName, oauth2Token.AccessToken)
	})
	if err != nil {
		return
	}

	return
}
