package apirequest

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/dankeroni/gotwitch"
	"github.com/pajlada/pajbot2/pkg/common/config"
)

// Twitch initialize the gotwitch api
// TODO: Do this in an Init method and use
// the proper oauth token. this will be
// required soon
var Twitch *gotwitch.TwitchAPI

// TwitchBot xD
var TwitchBot *gotwitch.TwitchAPI

type rateLimit struct {
	mutex *sync.RWMutex

	// The rate at which points are added to your bucket. This is the average number of requests per minute you can make over an extended period of time.
	Limit int

	// The number of points you have left to use.
	Remaining int

	// A timestamp of when your bucket is reset to full.
	Reset time.Time
}

func newRateLimit() rateLimit {
	return rateLimit{
		mutex: &sync.RWMutex{},
	}
}

func (l *rateLimit) Update(r *http.Response) {
	limit := r.Header.Get("Ratelimit-Limit")
	remaining := r.Header.Get("Ratelimit-Remaining")
	reset := r.Header.Get("Ratelimit-Reset")

	if limit == "" || remaining == "" || reset == "" {
		return
	}

	nLimit, err := strconv.Atoi(limit)
	if err != nil {
		fmt.Println("Error parsing limit from", limit)
	}
	nRemaining, err := strconv.Atoi(remaining)
	if err != nil {
		fmt.Println("Error parsing remaining from", remaining)
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.Limit = nLimit
	l.Remaining = nRemaining
}

func (l *rateLimit) String() string {
	return fmt.Sprintf("[RateLimit Limit=%d Remaining=%d Reset=%s]", l.Limit, l.Remaining, l.Reset)
}

type TwitchWrapperX struct {
	api *gotwitch.TwitchAPI

	RateLimit rateLimit
}

var TwitchWrapper *TwitchWrapperX

func InitTwitch(cfg *config.Config) (err error) {
	// Twitch APIs
	Twitch = gotwitch.New(cfg.Auth.Twitch.User.ClientID)
	Twitch.Credentials.ClientSecret = cfg.Auth.Twitch.User.ClientSecret
	_, err = Twitch.GetAppAccessTokenSimple()
	// TODO: Refresh the access token every now and then
	if err != nil {
		return
	}

	TwitchBot = gotwitch.New(cfg.Auth.Twitch.Bot.ClientID)
	TwitchBot.Credentials.ClientSecret = cfg.Auth.Twitch.Bot.ClientSecret
	_, err = TwitchBot.GetAppAccessTokenSimple()
	// TODO: Refresh the access token every now and then
	if err != nil {
		return
	}

	TwitchWrapper = &TwitchWrapperX{
		api: Twitch,

		RateLimit: newRateLimit(),
	}

	return
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
