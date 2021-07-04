package apirequest

import (
	"fmt"
	"log"
	"time"

	"github.com/nicklaw5/helix"
	"github.com/pajbot/pajbot2/pkg/common/config"
)

const WebhookDefaultTime = time.Hour * 24
const TimeToRefresh = time.Hour * 2

type TwitchWrapperX struct {
	helix    *helix.Client
	helixBot *helix.Client
	cfg      *config.TwitchWebhookConfig

	RateLimit TwitchRateLimit
}

var TwitchWrapper *TwitchWrapperX

// initAppAccessToken requests and sets app access token to the provided helix.Client
// and initializes a ticker running every 24 Hours which re-requests and sets app access token
func initAppAccessToken(helixAPI *helix.Client, tokenFetched chan struct{}) {
	response, err := helixAPI.RequestAppAccessToken([]string{})

	if err != nil {
		log.Fatalf("[Helix] Error requesting app access token: %s , \n %s", err.Error(), response.Error)
	}

	log.Printf("[Helix] Requested access token, status: %d, expires in: %d", response.StatusCode, response.Data.ExpiresIn)
	helixAPI.SetAppAccessToken(response.Data.AccessToken)
	close(tokenFetched)

	// initialize the ticker
	ticker := time.NewTicker(24 * time.Hour)

	for range ticker.C {
		response, err := helixAPI.RequestAppAccessToken([]string{})
		if err != nil {
			log.Printf("[Helix] Failed to re-request app access token from ticker, status: %d", response.StatusCode)
			continue
		}
		log.Printf("[Helix] Re-requested access token from ticker, status: %d, expires in: %d", response.StatusCode, response.Data.ExpiresIn)

		helixAPI.SetAppAccessToken(response.Data.AccessToken)
	}
}

func newHelixAPIClient(clientID, clientSecret string) (*helix.Client, chan struct{}, error) {
	apiClient, err := helix.NewClient(&helix.Options{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		return nil, nil, err
	}

	waitForFirstAppAccessToken := make(chan struct{})

	// Initialize methods responsible for refreshing oauth
	go initAppAccessToken(apiClient, waitForFirstAppAccessToken)

	return apiClient, waitForFirstAppAccessToken, nil
}

func initWrapper(cfg *config.AuthTwitchConfig) error {
	helixUser, userChan, err := newHelixAPIClient(cfg.User.ClientID, cfg.User.ClientSecret)
	if err != nil {
		return err
	}

	helixBot, botChan, err := newHelixAPIClient(cfg.Bot.ClientID, cfg.Bot.ClientSecret)
	if err != nil {
		return err
	}

	// Wait for both User and Bot clients to receive their App Access Token
	<-userChan
	<-botChan

	TwitchWrapper = &TwitchWrapperX{
		cfg:      &cfg.Webhook,
		helix:    helixUser,
		helixBot: helixBot,

		RateLimit: NewTwitchRateLimit(),
	}

	return nil
}

func (w *TwitchWrapperX) HelixUser() *helix.Client {
	return w.helix
}

func (w *TwitchWrapperX) HelixBot() *helix.Client {
	return w.helixBot
}

func (w *TwitchWrapperX) GetUsersByLogin(in []string) (data []helix.User, err error) {
	// TODO: RATE LIMITING XD
	resp, err := w.helix.GetUsers(&helix.UsersParams{
		Logins: in,
	})
	if err != nil {
		return nil, err
	}

	return resp.Data.Users, nil
}

func (w *TwitchWrapperX) GetUsersByID(in []string) (data []helix.User, err error) {
	// TODO: RATE LIMITING XD
	resp, err := w.helix.GetUsers(&helix.UsersParams{
		IDs: in,
	})
	if err != nil {
		return nil, err
	}

	return resp.Data.Users, nil
}

func (w *TwitchWrapperX) GetStreams(userIDs, userLogins []string) (data []helix.Stream, err error) {
	// TODO: RATE LIMITING XD
	resp, err := w.helix.GetStreams(&helix.StreamsParams{
		UserIDs:    userIDs,
		UserLogins: userLogins,
	})
	if err != nil {
		return nil, err
	}

	return resp.Data.Streams, nil
}

func (w *TwitchWrapperX) DeleteAllEventSubSubscriptions() error {
	// resp, err := w.helix.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{
	// 	Status: helix.EventSubStatusEnabled,
	// 	// Cursor: resp.Data.Pagination.Cursor,
	// })
	return nil
}

func (w *TwitchWrapperX) DeleteEventSubSubscription(id string) error {
	_, err := w.helix.RemoveEventSubSubscription(id)
	if err != nil {
		return err
	}

	return nil
}

func (w *TwitchWrapperX) EventSubSubscribe(eventType, channelID string) {
	_, err := w.helix.CreateEventSubSubscription(&helix.EventSubSubscription{
		Type:    eventType,
		Version: "1",
		Condition: helix.EventSubCondition{
			BroadcasterUserID: channelID,
		},
		Transport: helix.EventSubTransport{
			Method:   "webhook",
			Callback: w.cfg.HostPrefix + "/api/webhook/callback",
			Secret:   w.cfg.Secret,
		},
	})
	if err != nil {
		fmt.Println("Error subbing:", err)
		return
	}
}

func (w *TwitchWrapperX) GetWebhookSubscriptions() ([]helix.EventSubSubscription, error) {
	resp, err := w.helix.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{
		Status: helix.EventSubStatusEnabled,
	})
	if err != nil {
		fmt.Println("Error getting active subscriptions:", err)
		return nil, err
	}

	return resp.Data.EventSubSubscriptions, nil
}
