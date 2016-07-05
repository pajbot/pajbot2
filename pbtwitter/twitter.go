package pbtwitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/pajlada/pajbot2/common"
)

// Init logs into twitter and starts the stream
func Init(cfg *common.Config) *Client {
	// streaming client
	twitterCfg := oauth1.NewConfig(cfg.TwitterConsumerKey, cfg.TwitterConsumerSecret)
	token := oauth1.NewToken(cfg.TwitterAccessToken, cfg.TwitterAccessSecret)
	httpClient := twitterCfg.Client(oauth1.NoContext, token)
	streamClient := twitter.NewClient(httpClient)
	// rest api client
	anaconda.SetConsumerKey(cfg.TwitterConsumerKey)
	anaconda.SetConsumerSecret(cfg.TwitterConsumerSecret)
	rest := anaconda.NewTwitterApi(cfg.TwitterAccessToken, cfg.TwitterAccessSecret)

	c := &Client{
		StreamClient: streamClient,
		Bots:         make(map[string]*Bot),
		Rest:         rest,
	}
	go c.loadAllFollowed()
	go c.stream()
	return c
}

// TODO: store this in redis to avoid rate limits
func (c *Client) loadAllFollowed() {
	var all []string
	pages := c.Rest.GetFriendsListAll(nil)
	for page := range pages {
		if page.Error != nil {
			log.Error(page.Error)
		}
		for _, user := range page.Friends {
			all = append(all, user.Name)
			log.Debug(user.Name)
		}
	}
	c.followedUsers = all
	c.doneLoading = true
}
