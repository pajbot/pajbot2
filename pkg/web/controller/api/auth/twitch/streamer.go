package twitch

import (
	"net/http"

	"github.com/nicklaw5/helix/v2"
	"golang.org/x/oauth2"
)

func onStreamerAuthenticated(
	w http.ResponseWriter, r *http.Request,
	self *helix.ValidateTokenResponse, oauth2Token *oauth2.Token, stateData *stateData) {
}
