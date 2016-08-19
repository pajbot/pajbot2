package pbtwitter

import (
	"net/url"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/pajlada/pajbot2/common/config"
	"github.com/pajlada/pajbot2/redismanager"
)

// Init logs into twitter and starts the stream
func Init(cfg *config.Config, redis *redismanager.RedisManager) *Client {
	// streaming client
	twitterCfg := oauth1.NewConfig(cfg.Auth.Twitter.ConsumerKey, cfg.Auth.Twitter.ConsumerSecret)
	token := oauth1.NewToken(cfg.Auth.Twitter.AccessToken, cfg.Auth.Twitter.AccessSecret)
	httpClient := twitterCfg.Client(oauth1.NoContext, token)
	streamClient := twitter.NewClient(httpClient)
	// rest api client
	anaconda.SetConsumerKey(cfg.Auth.Twitter.ConsumerKey)
	anaconda.SetConsumerSecret(cfg.Auth.Twitter.ConsumerSecret)
	rest := anaconda.NewTwitterApi(cfg.Auth.Twitter.AccessToken, cfg.Auth.Twitter.AccessSecret)

	c := &Client{
		StreamClient: streamClient,
		Bots:         make(map[string]*Bot),
		Rest:         rest,
		redis:        redis,
		lastTweets:   make(map[string]*twitter.Tweet),
	}
	go c.loadAllFollowed()
	go c.stream()
	return c
}

// TODO: store this in redis to avoid rate limits
func (c *Client) loadAllFollowed() {
	// try loading from redis
	all, err := c.redis.LoadTwitterFollows()
	if err == nil {
		c.followedUsers = all
		c.doneLoading = true
		log.Debug("loaded twitter follows form redis")
		return
	}
	log.Error(err)
	v := url.Values{}
	v.Add("count", "200")
	pages := c.Rest.GetFriendsListAll(v)
	for page := range pages {
		if page.Error != nil {
			log.Error(page.Error)
		}
		for _, user := range page.Friends {
			all = append(all, strings.ToLower(user.ScreenName))
			log.Debug(user.ScreenName)
		}
	}
	c.followedUsers = all
	c.doneLoading = true
	c.redis.SaveTwitterFollows(all)
}
