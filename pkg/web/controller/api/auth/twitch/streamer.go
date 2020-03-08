package twitch

import (
	"net/http"

	"github.com/dankeroni/gotwitch/v2"
	"golang.org/x/oauth2"
)

func onStreamerAuthenticated(
	w http.ResponseWriter, r *http.Request,
	self gotwitch.ValidateResponse, oauth2Token *oauth2.Token, stateData *stateData) {
}
