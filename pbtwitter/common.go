package pbtwitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/pajlada/pajbot2/plog"
)

var log = plog.GetLogger()

// Client twitter client xD
type Client struct {
	StreamClient  *twitter.Client
	Bots          map[string]*Bot
	Rest          *anaconda.TwitterApi
	followedUsers []string // all users that are followed by given account
	doneLoading   bool     // wait to load all followed users until following new ones
}

// Bot contains the bots followed users
type Bot struct {
	Following []string
	Stream    chan *twitter.Tweet
	Client    *Client
}
