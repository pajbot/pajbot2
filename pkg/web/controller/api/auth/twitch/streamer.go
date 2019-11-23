package twitch

import (
	"net/http"

	"github.com/dankeroni/gotwitch"
	"golang.org/x/oauth2"
)

func onStreamerAuthenticated(w http.ResponseWriter, r *http.Request, self gotwitch.ValidateResponse, oauth2Token *oauth2.Token, stateData *stateData) {
	// fmt.Printf("STREAMER Username: %s - Access token: %s\n", self.Token.UserName, oauth2Token.AccessToken)
}
