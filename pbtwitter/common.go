package pbtwitter

import (
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/pajlada/pajbot2/plog"
	"github.com/pajlada/pajbot2/redismanager"
)

var log = plog.GetLogger()

// TwitterTimeFormat xD
const TwitterTimeFormat = time.RubyDate

// Client twitter client xD
type Client struct {
	StreamClient  *twitter.Client
	Bots          map[string]*Bot
	Rest          *anaconda.TwitterApi
	followedUsers []string // all users that are followed by given account
	doneLoading   bool     // wait to load all followed users until following new ones
	redis         *redismanager.RedisManager
	lastTweets    map[string]*twitter.Tweet
}

// Bot contains the bots followed users
type Bot struct {
	Following []string
	Stream    chan *twitter.Tweet
	Client    *Client
}
