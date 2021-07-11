package pkg

import (
	"encoding/json"

	"github.com/nicklaw5/helix"
)

type TwitchEventSubNotification struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Challenge    string                     `json:"challenge"`
	Event        json.RawMessage            `json:"event"`
}
