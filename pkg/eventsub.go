package pkg

import (
	"encoding/json"

	"github.com/nicklaw5/helix/v2"
)

type TwitchEventSubNotification struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Challenge    string                     `json:"challenge"`
	Event        json.RawMessage            `json:"event"`
}
