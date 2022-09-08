package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/nicklaw5/helix"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/common/config"
	"github.com/pajbot/pajbot2/pkg/web/state"
)

func apiEventsub(cfg *config.TwitchWebhookConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c := state.Context(w, r)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		defer r.Body.Close()

		if !helix.VerifyEventSubNotification(cfg.Secret, r.Header, string(body)) {
			log.Println("No valid signature in subscription message")
			return
		}

		var notification pkg.TwitchEventSubNotification
		err = json.NewDecoder(bytes.NewReader(body)).Decode(&notification)
		if err != nil {
			fmt.Println(err)
			return
		}

		if notification.Challenge != "" {
			w.Write([]byte(notification.Challenge))
			return
		}

		var botChannel pkg.BotChannel
		var channelID string

		switch notification.Subscription.Type {
		case helix.EventSubTypeChannelFollow:
			fallthrough
		case helix.EventSubTypeStreamOnline:
			fallthrough
		case helix.EventSubTypeStreamOffline:
			channelID = notification.Subscription.Condition.BroadcasterUserID
		}

		if channelID != "" {
			for it := c.Application.TwitchBots().Iterate(); it.Next(); {
				bot := it.Value()
				if bot == nil {
					continue
				}

				botChannel = bot.GetBotChannelByID(channelID)
				if botChannel == nil {
					continue
				}

				break
			}
		}

		if botChannel == nil {
			fmt.Println("No bot channel active to handle this request", channelID, notification.Subscription.Type)
			fmt.Println(string(notification.Event))
			// No bot channel active to handle this request
			return
		}

		err = botChannel.HandleEventSubNotification(notification)
		if err != nil {
			fmt.Println("Error handling eventsub notification:", err)
		}

		w.WriteHeader(http.StatusOK)
	}
}
