package apirequest

import (
	"fmt"
	"time"

	"github.com/dankeroni/gotwitch/v2"
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
	// fmt.Println("Subscriptions:", TwitchWrapper.WebhookSubscriptions)

	return nil
}

func (w *TwitchWrapperX) API() *gotwitch.TwitchAPI {
	return w.api
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
	// TODO: RATE LIMITING XD
	data, err := w.api.Helix().WebhookSubscribe(callbackURL, topic, userID, leaseTime, w.cfg.Secret)
	if err != nil {
		return err
	}

	fmt.Println("Response after subscribing:", string(data))

	return nil
}

func (w *TwitchWrapperX) GetUsersByLogin(in []string) (data []gotwitch.User, err error) {
	// TODO: RATE LIMITING XD
	return w.api.Helix().GetUsers(gotwitch.NewGetUsersParameters().SetUserLogins(in))
}

func (w *TwitchWrapperX) GetUsersByID(in []string) (data []gotwitch.User, err error) {
	// TODO: RATE LIMITING XD
	return w.api.Helix().GetUsers(gotwitch.NewGetUsersParameters().SetUserIDs(in))
}

func (w *TwitchWrapperX) GetStreams(userIDs, userLogins []string) (data []gotwitch.HelixStream, err error) {
	// TODO: RATE LIMITING XD
	return w.api.Helix().
		GetStreams(gotwitch.NewGetStreamsParameters().
			SetUserIDs(userIDs).
			SetUserLogins(userLogins))
}

func (w *TwitchWrapperX) GetWebhookSubscriptions(after, first string) (data gotwitch.WebhookSubscriptionsResponse, err error) {
	// TODO: RATE LIMITING XD
	return w.api.Helix().
		GetWebhookSubscriptions(after, first)
}
