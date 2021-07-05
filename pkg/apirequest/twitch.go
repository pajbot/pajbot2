package apirequest

import (
	"github.com/dankeroni/gotwitch/v2"
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
	Twitch.SetClientSecret(cfg.Auth.Twitch.User.ClientSecret)
	err = Twitch.Helix().RefreshAppAccessToken()
	if err != nil {
		return
	}

	TwitchBot = gotwitch.New(cfg.Auth.Twitch.Bot.ClientID)
	TwitchBot.SetClientSecret(cfg.Auth.Twitch.Bot.ClientSecret)
	err = TwitchBot.Helix().RefreshAppAccessToken()
	if err != nil {
		return
	}

	err = initWrapper(&cfg.Auth.Twitch.Webhook)

	return
}
