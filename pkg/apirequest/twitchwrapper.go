package apirequest

import (
	"fmt"
	"net/http"

	"github.com/dankeroni/gotwitch"
)

type TwitchWrapperX struct {
	api *gotwitch.TwitchAPI

	RateLimit TwitchRateLimit
}

var TwitchWrapper *TwitchWrapperX

func initWrapper() {
	TwitchWrapper = &TwitchWrapperX{
		api: Twitch,

		RateLimit: NewTwitchRateLimit(),
	}
}

func (w *TwitchWrapperX) GetUsersByLogin(in []string) (data []gotwitch.User, err error) {
	var response *http.Response
	data, response, err = w.api.GetUsersByLoginSimple(in)
	if response != nil {
		w.RateLimit.Update(response)
	}
	return
}

func (w *TwitchWrapperX) GetUsersByID(in []string) (data []gotwitch.User, err error) {
	var response *http.Response
	data, response, err = w.api.GetUsersSimple(in)
	if response != nil {
		w.RateLimit.Update(response)
	}
	return
}

func (w *TwitchWrapperX) GetStreams(userIDs, userLogins []string) (data []gotwitch.Stream, err error) {
	var response *http.Response
	data, response, err = w.api.GetStreamsSimple(userIDs, userLogins)
	if response != nil {
		w.RateLimit.Update(response)
		fmt.Println("Executed twitch request:", w.RateLimit.String())
	}
	return
}
