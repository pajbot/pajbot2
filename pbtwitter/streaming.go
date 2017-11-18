package pbtwitter

import (
	"log"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/pajlada/pajbot2/helper"
)

// Follow follows given users timeline stream
func (bot *Bot) Follow(user string) {
	if user == "ALL" {
		bot.Following = bot.Client.followedUsers
		return
	}
	u := strings.ToLower(user)
	for _, usr := range bot.Following {
		if usr == u {
			return
		}
	}
	bot.Following = append(bot.Following, u)
	go bot.Client.Follow(u)
}

func (c *Client) streamToBots(tweet *twitter.Tweet) {
	log.Println(tweet.Text)
	if tweet.RetweetedStatus != nil || tweet.QuotedStatus != nil {
		return
	}
	// cache tweet for lasttweet
	c.lastTweets[strings.ToLower(tweet.User.ScreenName)] = tweet
	for _, bot := range c.Bots {
		for _, followedUser := range bot.Following {
			if strings.ToLower(tweet.User.ScreenName) == followedUser {
				bot.Stream <- tweet
			}
		}
	}
}

// stream starts the stream
func (c *Client) stream() {
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		tweet.Text = helper.RemoveNewlines(tweet.Text)
		go c.streamToBots(tweet)
	}
	params := &twitter.StreamUserParams{
		With:          "followings",
		StallWarnings: twitter.Bool(true),
	}
	stream, err := c.StreamClient.Streams.User(params)
	if err != nil {
		log.Fatal(err)
	}
	demux.HandleChan(stream.Messages)
}
