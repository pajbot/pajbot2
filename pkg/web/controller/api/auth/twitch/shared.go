package twitch

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/dankeroni/gotwitch"
	"github.com/gorilla/mux"
	"github.com/pajbot/pajbot2/pkg/apirequest"
	"github.com/pajbot/pajbot2/pkg/web/router"
	"github.com/pajbot/utils"
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

			validateResponse, err := apirequest.Twitch.ID().Authenticate(oauth2Token.AccessToken).Validate()
			if err != nil {
				http.Error(w, "Error validating token: "+err.Error(), http.StatusInternalServerError)
				return
			}

			onAuthorized(w, r, validateResponse, oauth2Token, stateData)
		})
}
