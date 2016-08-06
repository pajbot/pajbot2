package apirequest

import "github.com/dankeroni/gotwitch"

// Twitch initialize the gotwitch api
// TODO: Do this in an Init method and use
// the proper oauth token. this will be
// required soon
var Twitch = gotwitch.New("xD")
