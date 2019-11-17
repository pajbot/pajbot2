package apirequest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dankeroni/gotwitch"
	"github.com/pajbot/pajbot2/pkg/common/config"
)

const WebhookDefaultTime = time.Hour * 24
const TimeToRefresh = time.Hour * 2

type TwitchWrapperX struct {
	api *gotwitch.TwitchAPI
	cfg *config.TwitchWebhookConfig

	RateLimit TwitchRateLimit

	WebhookSubscriptions []gotwitch.WebhookSubscription
}

var TwitchWrapper *TwitchWrapperX

func initWrapper(cfg *config.TwitchWebhookConfig) error {
	TwitchWrapper = &TwitchWrapperX{
		cfg: cfg,
		api: Twitch,

		RateLimit: NewTwitchRateLimit(),
	}

	subscriptions, err := TwitchWrapper.GetWebhookSubscriptions("", "")
	if err != nil {
		fmt.Println("ERROR GETTING WEBHOOK SUBSCRIPTIONS:", err)
		return err
	}
	// TODO: Every now and then, refresh our list of webhook subscriptions
	// TODO: follow paginations

	TwitchWrapper.WebhookSubscriptions = subscriptions.Data
	fmt.Println("Subscriptions:", TwitchWrapper.WebhookSubscriptions)

	return nil
}

func (w *TwitchWrapperX) WebhookSubscribe(topic gotwitch.WebhookTopic, userID string) error {
	url := topic.URL(userID)
	callbackURL := w.cfg.HostPrefix + "/" + userID + "/" + topic.String()

	for _, subscription := range w.WebhookSubscriptions {
		if subscription.Topic == url &&
			subscription.Callback == callbackURL {
			if subscription.ExpiresAt.Add(-TimeToRefresh).Before(time.Now()) {
				// We are subscribed already, but it's time to refresh our subscription
				break
			}

			// We are already subscribed to this topic with the same callback URL
			return nil
		}
	}

	leaseTime := time.Duration(w.cfg.LeaseTimeSeconds) * time.Second
	// Subscribe!
	data, response, err := w.api.WebhookSubscribeSimple(callbackURL, topic, userID, leaseTime, w.cfg.Secret)
	if err != nil {
		return err
	}
	if response != nil {
		w.RateLimit.Update(response)
	}

	fmt.Println("Response after subscribing:", string(*data))

	return nil
}

func (w *TwitchWrapperX) GetUsersByLogin(in []string) (data []gotwitch.User, err error) {
	var response *http.Response
	data, response, err = w.api.GetUsersByLoginSimple(in)
	if response != nil {
		w.RateLimit.Update(response)
	}
	return
}

func (w *TwitchWrapperX) GetUsersByID(in []string) (data []gotwitch.User, err error) {
	var response *http.Response
	data, response, err = w.api.GetUsersSimple(in)
	if response != nil {
		w.RateLimit.Update(response)
	}
	return
}

func (w *TwitchWrapperX) GetStreams(userIDs, userLogins []string) (data []gotwitch.Stream, err error) {
	var response *http.Response
	data, response, err = w.api.GetStreamsSimple(userIDs, userLogins)
	if response != nil {
		w.RateLimit.Update(response)
		fmt.Println("Executed twitch request:", w.RateLimit.String())
	}
	return
}

func (w *TwitchWrapperX) GetWebhookSubscriptions(after, first string) (data *gotwitch.WebhookSubscriptionsResponse, err error) {
	var response *http.Response
	data, response, err = w.api.GetWebhookSubscriptionsSimple(after, first)
	if response != nil {
		w.RateLimit.Update(response)
		fmt.Println("Executed twitch request:", w.RateLimit.String())
	}
	return
}
