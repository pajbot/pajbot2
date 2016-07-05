package pbtwitter

import (
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/pajlada/pajbot2/common"
)

type Client struct {
	TC   *twitter.Client
	Bots map[string]*Bot
}

type Bot struct {
	Following []string
	Stream    chan *twitter.Tweet
}

func Init(cfg *common.Config) *Client {
	twitterCfg := oauth1.NewConfig(cfg.TwitterConsumerKey, cfg.TwitterConsumerSecret)
	token := oauth1.NewToken(cfg.TwitterAccessToken, cfg.TwitterAccessSecret)
	httpClient := twitterCfg.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	return &Client{
		TC:   client,
		Bots: make(map[string]*Bot),
	}
}

func (bot *Bot) Follow(user string) {
	bot.Following = append(bot.Following, strings.ToLower(user))
}

func (c *Client) streamToBots(tweet *twitter.Tweet) {
	log.Debug(tweet.Text)
	for _, bot := range c.Bots {
		for _, followedUser := range bot.Following {
			if strings.ToLower(tweet.User.Name) == followedUser {
				bot.Stream <- tweet
			}
		}
	}
}

func (c *Client) Stream() {
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		go c.streamToBots(tweet)
	}
	params := &twitter.StreamUserParams{
		With:          "followings",
		StallWarnings: twitter.Bool(true),
	}
	stream, err := c.TC.Streams.User(params)
	if err != nil {
		log.Fatal(err)
	}
	demux.HandleChan(stream.Messages)
}
