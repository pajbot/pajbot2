package apirequest

import "github.com/dankeroni/gotwitch"

// Twitch initialize the gotwitch api
// TODO: Do this in an Init method and use
// the proper oauth token. this will be
// required soon
var Twitch = &gotwitch.TwitchAPI{}

// TwitchBot xD
var TwitchBot = &gotwitch.TwitchAPI{}

// TwitchV3 xD
var TwitchV3 = &gotwitch.TwitchAPI{}

// TwitchBotV3 xD
var TwitchBotV3 = &gotwitch.TwitchAPI{}
