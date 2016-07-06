package pbtwitter

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/dghubble/go-twitter/twitter"
)

// Follow follows given user on twitter with the given account
func (c *Client) Follow(targetUser string) error {
	for _, u := range c.followedUsers {
		if u == targetUser {
			// return if already following given user
			return nil
		}
	}
	for !c.doneLoading {
		time.Sleep(5 * time.Second)
		log.Debug("waiting to load all users before following new ones")
	}
	user, err := c.Rest.FollowUser(targetUser)
	if err != nil {
		return err
	}
	log.Debug("followed ", user.Name)
	c.followedUsers = append(c.followedUsers, targetUser)
	c.redis.SaveTwitterFollows(c.followedUsers)
	return nil
}

/*
LastTweetString returns the newest tweet from the given user formatted like this:
(@nuulss): xD test LUL | 4 hours 20 minutes ago
*/
func (bot *Bot) LastTweetString(user string) string {
	c := bot.Client
	tTweet, aTweet := c.LastTweet(user)
	if tTweet == nil && aTweet == nil {
		return "no tweet found :("
	}
	if tTweet != nil {
		created, err := time.Parse(TwitterTimeFormat, tTweet.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		ago := time.Since(created)
		mins := int(ago.Minutes())
		m := mins % 60
		hours := (mins - m) / 60
		mins = m
		agoString := fmt.Sprintf("%dh %dm ago", hours, mins)
		return fmt.Sprintf("(@%s): %s | %s", tTweet.User.ScreenName, tTweet.Text, agoString)
	}
	created, err := time.Parse(TwitterTimeFormat, aTweet.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}
	ago := time.Since(created)
	mins := int(ago.Minutes())
	m := mins % 60
	hours := (mins - m) / 60
	mins = m
	agoString := fmt.Sprintf("%dh %dm ago", hours, mins)
	return fmt.Sprintf("(@%s): %s | %s", aTweet.User.ScreenName, aTweet.Text, agoString)

}

/*
LastTweet returns the newest tweet from given user, *twitter.Tweet if it was cached
or *anaconda.Tweet if it wasnt
*/
func (c *Client) LastTweet(user string) (*twitter.Tweet, *anaconda.Tweet) {
	user = strings.ToLower(user)
	if tweet, ok := c.lastTweets[user]; ok {
		log.Debug("found cached tweet")
		return tweet, nil
	}
	v := url.Values{}
	v.Add("count", "1")
	v.Add("screen_name", user)
	v.Add("exclude_replies", "true")
	v.Add("include_rts", "false")
	tweets, err := c.Rest.GetUserTimeline(v)
	if err != nil {
		log.Error(err)
		return nil, nil
	}
	tweet := tweets[0]
	return nil, &tweet
}
