package apirequest

import (
	"github.com/dankeroni/gotwitch"
	"github.com/pajbot/pajbot2/pkg/common/config"
)

// Twitch initialize the gotwitch api
// TODO: Do this in an Init method and use
// the proper oauth token. this will be
// required soon
var Twitch *gotwitch.TwitchAPI

// TwitchBot xD
var TwitchBot *gotwitch.TwitchAPI

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

	err = initWrapper(&cfg.Auth.Twitch.Webhook)

	return
}
